package routes

import (
	"crypto/rand"
	"encoding/json"

	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"time"

	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"
	utils "quiz-app/internal/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type routerClass struct {
	auth         auth.AuthHandler
	classUseCase usecase.ClassUseCase
	redisUseCase usecase.RedisUseCase
}

func NewRouterClass(classUC usecase.ClassUseCase, redisUC usecase.RedisUseCase, auth auth.AuthHandler) routerClass {
	return routerClass{
		auth:         auth,
		classUseCase: classUC,
		redisUseCase: redisUC,
	}
}

// ==== Request Body Structs ====
type classIDRequest struct {
	ID     string `json:"class_id"`
	Minute int    `json:"minute"`
}

type classDeleteRequest struct {
	ID string `json:"_id"`
}

type ResetClassTest struct {
	ClassID string `json:"class_id"`
	TestID  string `json:"test_id"`
}

type joinClassRequest struct {
	ID string `json:"key"`
}

// ==== Helpers ====

// generateInviteCode returns a 6-char alphanumeric uppercase code per CLAUDE.md F3.
// Uses crypto/rand for unpredictability (~ log2(36^6) = 31 bits → enough TTL/24h).
const inviteCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // skip 0/O/1/I — visually ambiguous

func generateRandomKey() (string, error) {
	const n = 6
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = inviteCodeAlphabet[int(buf[i])%len(inviteCodeAlphabet)]
	}
	return string(out), nil
}

// ==== Handlers ====
func (rc routerClass) createClass(w http.ResponseWriter, req *http.Request) {
	email, ok := req.Context().Value(constant.CtxEmail).(string)
	userID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "email not found in context", http.StatusUnauthorized)
		return
	}
	var newClass entity.Class
	if !utils.DecodeJSONBody(w, req, &newClass) {
		return
	}

	// Owner is implicitly a member of their own class — but ONLY in
	// students_accept. Inserting into students_wait was a bug (would surface
	// owner as a pending request to themselves).
	newClass.StudentsAccept = []entity.StudentInfo{{StudentEmail: email, StudentID: userID}}
	newClass.StudentsWait = []entity.StudentInfo{}

	if newClass.TestID == nil {
		newClass.TestID = []primitive.ObjectID{}
	}
	newClass.UserID = userID
	newClass.AuthorMail = email
	result, err := rc.classUseCase.CreateClass(req.Context(), newClass)
	if err != nil {
		pkg.SendError(w, "Failed to create class", http.StatusInternalServerError)
		return
	}
	newClass.ID = result

	pkg.SendResponse(w, http.StatusOK, newClass)
}

func (rc routerClass) getAllClass(w http.ResponseWriter, req *http.Request) {
	userId, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	classes, err := rc.classUseCase.GetAllClass(req.Context(), userId)
	if err != nil {
		pkg.SendError(w, "Failed to retrieve classes", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, classes)
}

func (rc routerClass) updateClass(w http.ResponseWriter, req *http.Request) {
	email, ok := req.Context().Value(constant.CtxEmail).(string)
	userID, ok := req.Context().Value(constant.CtxUserID).(string)

	if !ok {
		pkg.SendError(w, "email not found in context", http.StatusUnauthorized)
		return
	}
	var classToUpdate entity.Class
	if !utils.DecodeJSONBody(w, req, &classToUpdate) {
		return
	}

	classToUpdate.AuthorMail = email
	classToUpdate.UserID = userID
	updatedClass, err := rc.classUseCase.UpdateClass(req.Context(), classToUpdate)
	if err != nil {
		pkg.SendError(w, "Failed to update class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, updatedClass)
}

func (rc routerClass) deleteClass(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "userID not found in context", http.StatusUnauthorized)
		return
	}
	idHex := mux.Vars(req)["id"]
	idClass, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		pkg.SendError(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := rc.classUseCase.DeleteClass(req.Context(), userID, idClass); err != nil {
		pkg.SendError(w, "Error deleting class", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// approveStudent / rejectStudent — F3 missing endpoints. Body: { student_id }.
type studentActionRequest struct {
	StudentID string `json:"student_id"`
}

func (rc routerClass) approveStudent(w http.ResponseWriter, req *http.Request) {
	ownerID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	classID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	if err != nil {
		pkg.SendError(w, "invalid class id", http.StatusBadRequest)
		return
	}
	var body studentActionRequest
	if !utils.DecodeJSONBody(w, req, &body) {
		return
	}
	if body.StudentID == "" {
		pkg.SendError(w, "student_id required", http.StatusBadRequest)
		return
	}
	if err := rc.classUseCase.ApproveStudent(req.Context(), classID, ownerID, body.StudentID); err != nil {
		pkg.SendError(w, "approve failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rc routerClass) rejectStudent(w http.ResponseWriter, req *http.Request) {
	ownerID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	classID, err := primitive.ObjectIDFromHex(mux.Vars(req)["id"])
	if err != nil {
		pkg.SendError(w, "invalid class id", http.StatusBadRequest)
		return
	}
	var body studentActionRequest
	if !utils.DecodeJSONBody(w, req, &body) {
		return
	}
	if body.StudentID == "" {
		pkg.SendError(w, "student_id required", http.StatusBadRequest)
		return
	}
	if err := rc.classUseCase.RejectStudent(req.Context(), classID, ownerID, body.StudentID); err != nil {
		pkg.SendError(w, "reject failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (rc routerClass) resetTestInRedis(w http.ResponseWriter, req *http.Request) {
	var info ResetClassTest
	if !utils.DecodeJSONBody(w, req, &info) {
		return
	}

	testInfoKey := fmt.Sprintf("classTest:%s:%s", info.ClassID, info.TestID)
	questionsFullKey := fmt.Sprintf("classTestQuestionsFull:%s:%s", info.ClassID, info.TestID)

	// Optional: if you also want to clear per-student data
	// examQuestionsKey := fmt.Sprintf("examQuestions:%s:%s", info.TestID, info.UserID)
	// sessionKey := fmt.Sprintf("examSession:%s:%s", info.TestID, info.UserID)

	// Delete all related keys
	if err := rc.redisUseCase.Delete(req.Context(), testInfoKey); err != nil {
		pkg.SendError(w, "Error deleting test info", http.StatusInternalServerError)
		return
	}

	if err := rc.redisUseCase.Delete(req.Context(), questionsFullKey); err != nil {
		pkg.SendError(w, "Error deleting questions full", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(
		w,
		http.StatusOK,
		fmt.Sprintf("RESET OK for Test: %s", info.TestID),
	)
}

func (rc routerClass) joinClass(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value(constant.CtxEmailID).(string)
	email, ok := req.Context().Value(constant.CtxEmail).(string)

	if !ok {
		pkg.SendError(w, "email not found in context", http.StatusUnauthorized)
		return
	}
	var reqBody joinClassRequest
	if !utils.DecodeJSONBody(w, req, &reqBody) {
		return
	}

	dataJSON, err := rc.redisUseCase.Get(req.Context(), reqBody.ID)
	if err != nil {
		pkg.SendError(w, "Error fetching class from Redis", http.StatusInternalServerError)
		return
	}

	var redisData struct {
		ClassID string `json:"class_id"`
		UserID  string `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(dataJSON), &redisData); err != nil {
		pkg.SendError(w, "Error decoding Redis data", http.StatusInternalServerError)
		return
	}

	classID, err := primitive.ObjectIDFromHex(redisData.ClassID)
	if err != nil {
		pkg.SendError(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	err = rc.classUseCase.JoinClass(req.Context(), classID, redisData.UserID, emailID, email)
	if err != nil {
		pkg.SendError(w, "Error joining class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Joined class: %s", reqBody.ID),
	})
}

func (rc routerClass) createCodeClass(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value(constant.CtxEmailID).(string)
	if !ok {
		pkg.SendError(w, "email not found in context", http.StatusUnauthorized)
		return
	}
	var reqBody classIDRequest
	if !utils.DecodeJSONBody(w, req, &reqBody) {
		return
	}

	key, err := generateRandomKey()
	if err != nil {
		pkg.SendError(w, "Failed to generate code", http.StatusInternalServerError)
		return
	}

	data := struct {
		ClassID string `json:"class_id"`
		UserID  string `json:"user_id"`
	}{
		ClassID: reqBody.ID,
		UserID:  emailID,
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		pkg.SendError(w, "Failed to encode data", http.StatusInternalServerError)
		return
	}

	err = rc.redisUseCase.Set(req.Context(), key, string(dataJSON), time.Duration(reqBody.Minute)*time.Minute)
	if err != nil {
		pkg.SendError(w, "Failed to store code", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, key)
}
func (rc routerClass) GetClassRouter(r *Router) {
	// ===== TEACHER =====
	r.Router.Handle(
		"/class",
		rc.auth.AuthWithPerm("TEACHER", rc.getAllClass),
	).Methods("GET")

	r.Router.Handle(
		"/class",
		rc.auth.AuthWithPerm("TEACHER", rc.createClass),
	).Methods("POST")

	r.Router.Handle(
		"/class",
		rc.auth.AuthWithPerm("TEACHER", rc.updateClass),
	).Methods("PATCH")

	r.Router.Handle(
		"/class/{id}",
		rc.auth.AuthWithPerm("TEACHER", rc.deleteClass),
	).Methods("DELETE")

	r.Router.Handle(
		"/class/reset-test",
		rc.auth.AuthWithPerm("TEACHER", rc.resetTestInRedis),
	).Methods("POST")

	r.Router.Handle(
		"/class/code-class",
		rc.auth.AuthWithPerm("TEACHER", rc.createCodeClass),
	).Methods("POST")

	// Pending student moderation (F3 spec).
	r.Router.Handle(
		"/class/{id}/approve",
		rc.auth.AuthWithPerm("TEACHER", rc.approveStudent),
	).Methods("POST")

	r.Router.Handle(
		"/class/{id}/reject",
		rc.auth.AuthWithPerm("TEACHER", rc.rejectStudent),
	).Methods("POST")

	// ===== STUDENT =====

	r.Router.Handle(
		"/class/join-class",
		rc.auth.AuthWithPerm("STUDENT", rc.joinClass),
	).Methods("POST")

}
