package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/infrastructure/persistence/aws"
	"quiz-app/internal/pkg"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	maxImageSize    = 5 * 1024 * 1024  // 5 MB
	maxAudioSize    = 10 * 1024 * 1024 // 10 MB
	maxVideoSize    = 50 * 1024 * 1024 // 50 MB
	uploadURLExpire = 15 * time.Minute
	signedURLExpire = 15 * time.Minute
)

var allowedMIMEPrefixes = []string{"image/", "audio/", "video/", "application/pdf"}

type RoutesFile struct {
	auth         *auth.AuthHandler
	fileUseCase  *usecase.FileUseCase
	awsS3UseCase *aws.FileAWSRepository
}

func NewRoutesFile(
	uc *usecase.FileUseCase,
	awsS3UseCase *aws.FileAWSRepository,
	auth *auth.AuthHandler,
) *RoutesFile {
	return &RoutesFile{
		fileUseCase:  uc,
		awsS3UseCase: awsS3UseCase,
		auth:         auth,
	}
}

// GetRoutesFile registers file endpoints. Privacy invariants per CLAUDE.md E#15-17:
//   - S3 keys are opaque (file_id-based), never contain user_id/email.
//   - Download endpoint takes file_id (not author_email + filename).
//   - Responses to non-owners NEVER include owner_id, owner email, or s3_key.
func (rf *RoutesFile) GetRoutesFile(r *Router) {
	r.Router.Handle("/files",
		rf.auth.AuthMiddleware(http.HandlerFunc(rf.presignUpload)),
	).Methods(http.MethodPost)

	r.Router.Handle("/files",
		rf.auth.AuthMiddleware(http.HandlerFunc(rf.listMyFiles)),
	).Methods(http.MethodGet)

	r.Router.Handle("/files/{file_id}/url",
		rf.auth.AuthMiddleware(http.HandlerFunc(rf.signedURL)),
	).Methods(http.MethodGet)

	r.Router.Handle("/files/{file_id}",
		rf.auth.AuthMiddleware(http.HandlerFunc(rf.deleteFile)),
	).Methods(http.MethodDelete)
}

/* ── POST /files — owner upload presign ──────────────────────────────────── */

type presignUploadReq struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type presignUploadResp struct {
	FileID      string `json:"file_id"`
	UploadURL   string `json:"upload_url"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	ExpiresIn   int    `json:"expires_in"`
}

func (rf *RoutesFile) presignUpload(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := r.Context().Value(constant.CtxEmailID).(string)
	if !ok || ownerID == "" {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req presignUploadReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := validateUpload(req); err != nil {
		pkg.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Idempotency: same (owner, filename) returns existing record's upload URL
	// rather than creating a duplicate. Also enforced via unique index.
	if existing, err := rf.fileUseCase.FindByOwnerAndFilename(r.Context(), ownerID, req.Filename); err == nil && existing != nil {
		uploadURL, err := rf.awsS3UseCase.PresignUploadURL(r.Context(), existing.S3Key, existing.ContentType, uploadURLExpire)
		if err != nil {
			pkg.SendError(w, "Presign failed", http.StatusInternalServerError)
			return
		}
		pkg.SendResponse(w, http.StatusOK, presignUploadResp{
			FileID:      existing.ID.Hex(),
			UploadURL:   uploadURL,
			Filename:    existing.Filename,
			ContentType: existing.ContentType,
			Size:        existing.Size,
			ExpiresIn:   int(uploadURLExpire.Seconds()),
		})
		return
	}

	file := entity.NewFile(ownerID, req.Filename, req.ContentType, req.Size)
	file.ID = primitive.NewObjectID()
	file.S3Key = entity.BuildS3Key(file.ID, req.Filename)

	if _, err := rf.fileUseCase.Create(r.Context(), file); err != nil {
		pkg.SendError(w, "Create file failed", http.StatusInternalServerError)
		return
	}

	uploadURL, err := rf.awsS3UseCase.PresignUploadURL(r.Context(), file.S3Key, file.ContentType, uploadURLExpire)
	if err != nil {
		_ = rf.fileUseCase.DeleteOwned(r.Context(), file.ID, ownerID)
		pkg.SendError(w, "Presign failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusCreated, presignUploadResp{
		FileID:      file.ID.Hex(),
		UploadURL:   uploadURL,
		Filename:    file.Filename,
		ContentType: file.ContentType,
		Size:        file.Size,
		ExpiresIn:   int(uploadURLExpire.Seconds()),
	})
}

func validateUpload(req presignUploadReq) error {
	if strings.TrimSpace(req.Filename) == "" {
		return errors.New("filename required")
	}
	if req.ContentType == "" {
		return errors.New("content_type required")
	}
	allowed := false
	for _, p := range allowedMIMEPrefixes {
		if strings.HasPrefix(req.ContentType, p) {
			allowed = true
			break
		}
	}
	if !allowed {
		return errors.New("content_type not allowed")
	}
	if req.Size <= 0 {
		return errors.New("size must be > 0")
	}
	switch {
	case strings.HasPrefix(req.ContentType, "image/") && req.Size > maxImageSize:
		return errors.New("image too large (max 5MB)")
	case strings.HasPrefix(req.ContentType, "audio/") && req.Size > maxAudioSize:
		return errors.New("audio too large (max 10MB)")
	case strings.HasPrefix(req.ContentType, "video/") && req.Size > maxVideoSize:
		return errors.New("video too large (max 50MB)")
	}
	return nil
}

/* ── GET /files — owner list (paginated) ─────────────────────────────────── */

type fileListItem struct {
	FileID      string `json:"file_id"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	CreatedAt   string `json:"created_at"`
}

type paginatedFiles struct {
	Items []fileListItem `json:"items"`
	Total int64          `json:"total"`
	Page  int64          `json:"page"`
	Limit int64          `json:"limit"`
}

func (rf *RoutesFile) listMyFiles(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := r.Context().Value(constant.CtxEmailID).(string)
	if !ok || ownerID == "" {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	page, _ := strconv.ParseInt(q.Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(q.Get("limit"), 10, 64)

	mimePrefix := ""
	switch q.Get("type") {
	case "image":
		mimePrefix = "image/"
	case "audio":
		mimePrefix = "audio/"
	case "video":
		mimePrefix = "video/"
	}

	files, total, err := rf.fileUseCase.ListByOwner(r.Context(), ownerID, mimePrefix, page, limit)
	if err != nil {
		pkg.SendError(w, "List files failed", http.StatusInternalServerError)
		return
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	items := make([]fileListItem, 0, len(files))
	for _, f := range files {
		items = append(items, fileListItem{
			FileID:      f.ID.Hex(),
			Filename:    f.Filename,
			ContentType: f.ContentType,
			Size:        f.Size,
			CreatedAt:   f.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	pkg.SendResponse(w, http.StatusOK, paginatedFiles{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

/* ── GET /files/{file_id}/url — signed download URL (any auth user) ────── */

type signedURLResp struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	ExpiresAt   string `json:"expires_at"`
}

// signedURL returns a short-lived signed download URL for the given file_id.
// Any authenticated user can fetch — file_id is non-guessable (24-char hex).
// Response contains NO owner identifier.
func (rf *RoutesFile) signedURL(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value(constant.CtxEmailID).(string); !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idHex := mux.Vars(r)["file_id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		pkg.SendError(w, "invalid file_id", http.StatusBadRequest)
		return
	}

	file, err := rf.fileUseCase.FindByID(r.Context(), id)
	if err != nil {
		pkg.SendError(w, "lookup failed", http.StatusInternalServerError)
		return
	}
	if file == nil {
		pkg.SendError(w, "file not found", http.StatusNotFound)
		return
	}

	url, err := rf.awsS3UseCase.PresignDownloadURL(r.Context(), file.S3Key, signedURLExpire)
	if err != nil {
		pkg.SendError(w, "presign failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, signedURLResp{
		URL:         url,
		ContentType: file.ContentType,
		ExpiresAt:   time.Now().UTC().Add(signedURLExpire).Format(time.RFC3339),
	})
}

/* ── DELETE /files/{file_id} — owner only ────────────────────────────────── */

func (rf *RoutesFile) deleteFile(w http.ResponseWriter, r *http.Request) {
	ownerID, ok := r.Context().Value(constant.CtxEmailID).(string)
	if !ok || ownerID == "" {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idHex := mux.Vars(r)["file_id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		pkg.SendError(w, "invalid file_id", http.StatusBadRequest)
		return
	}

	file, err := rf.fileUseCase.FindByID(r.Context(), id)
	if err != nil {
		pkg.SendError(w, "lookup failed", http.StatusInternalServerError)
		return
	}
	// Same response for not-found and not-owned — avoids ownership oracle.
	if file == nil || file.OwnerID != ownerID {
		pkg.SendError(w, "file not found", http.StatusNotFound)
		return
	}

	if err := rf.fileUseCase.DeleteOwned(r.Context(), id, ownerID); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			pkg.SendError(w, "file not found", http.StatusNotFound)
			return
		}
		pkg.SendError(w, "delete failed", http.StatusInternalServerError)
		return
	}
	// Best-effort S3 delete. Orphan object reclaimed by lifecycle rule.
	_ = rf.awsS3UseCase.DeleteObject(r.Context(), file.S3Key)

	w.WriteHeader(http.StatusNoContent)
}
