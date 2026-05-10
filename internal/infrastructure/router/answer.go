package routes

import (
	"fmt"
	"net/http"
	"quiz-app/internal/auth"
	"quiz-app/internal/domain/usecase"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Question struct {
	ID              string       `json:"_id"`
	CreatedAt       string       `json:"created_at"`
	Options         []Option     `json:"options"`
	QuestionContent QuestionText `json:"question_content"`
	Score           float64      `json:"score"`
	Tags            []string     `json:"tags"`
	Type            string       `json:"type"`
	UpdatedAt       string       `json:"updated_at"`
}

type Metadata struct {
	Author string `json:"author"`
}

type Option struct {
	ID        string `json:"id"`
	IsCorrect bool   `json:"is_correct"`
	Text      string `json:"text"`
	Match     string `json:"match"`
	FileID    string `json:"file_id,omitempty"`
	Order     int    `json:"order"`
}

type QuestionText struct {
	Text string `json:"text"`
}

type QuestionsFromRedis []Question

type RouterSubmission struct {
	auth *auth.AuthHandler

	testTemplateUseCase *usecase.TestTemplateUseCase
	submissionUseCase   *usecase.SubmissionUseCase
	questionUseCase     *usecase.QuestionUseCase
	redisUseCase        *usecase.RedisUseCase
}

func NewRouterSubmission(s *usecase.SubmissionUseCase, t *usecase.TestTemplateUseCase, q *usecase.QuestionUseCase, r *usecase.RedisUseCase, auth *auth.AuthHandler) RouterSubmission {
	return RouterSubmission{
		auth:                auth,
		submissionUseCase:   s,
		testTemplateUseCase: t,
		redisUseCase:        r,
		questionUseCase:     q,
	}
}

func (rc *RouterSubmission) exportSubmissionPDF(
	w http.ResponseWriter,
	req *http.Request,
) {
	ctx := req.Context()
	classIDStr := req.URL.Query().Get("class_id")
	testIDStr := req.URL.Query().Get("test_id")

	if classIDStr == "" || testIDStr == "" {
		http.Error(w, "missing class_id or test_id", http.StatusBadRequest)
		return
	}

	classID, err := primitive.ObjectIDFromHex(classIDStr)
	if err != nil {
		http.Error(w, "invalid class_id", http.StatusBadRequest)
		return
	}

	testID, err := primitive.ObjectIDFromHex(testIDStr)
	if err != nil {
		http.Error(w, "invalid test_id", http.StatusBadRequest)
		return
	}

	pdfBytes, fileName, err := rc.submissionUseCase.ExportPDF(
		ctx,
		classID,
		testID,
	)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}

func (rc RouterSubmission) GetSubmissionRouter(r *Router) {
	r.Router.Handle(
		"/submissions/export/pdf",
		rc.auth.AuthWithPerm("TEACHER", rc.exportSubmissionPDF),
	).Methods("GET")
}
