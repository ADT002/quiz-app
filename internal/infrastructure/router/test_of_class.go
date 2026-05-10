package routes

import (
	"encoding/json"
	"net/http"
	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesTestOfClass struct {
	auth               auth.AuthHandler
	testOfClassUseCase usecase.TestOfClassUseCase
}

func NewRouterTestOfClass(testOfClassUseCase usecase.TestOfClassUseCase, auth auth.AuthHandler) RoutesTestOfClass {
	return RoutesTestOfClass{
		testOfClassUseCase: testOfClassUseCase,
		auth:               auth,
	}
}

func (rt RoutesTestOfClass) GetTestOfClassRouter(r *Router) {
	r.Handle(
		"/test-of-class",
		rt.auth.AuthWithPerm("TEACHER", rt.getTestsOfClass),
	).Methods("GET")

	r.Handle(
		"/test-of-class",
		rt.auth.AuthWithPerm("TEACHER", rt.createTestOfClass),
	).Methods("POST")

	r.Handle(
		"/test-of-class",
		rt.auth.AuthWithPerm("TEACHER", rt.updateTestOfClass),
	).Methods("PATCH")

	r.Handle(
		"/test-of-class",
		rt.auth.AuthWithPerm("TEACHER", rt.deleteTestOfClass),
	).Methods("DELETE")
}

func (r *RoutesTestOfClass) createTestOfClass(w http.ResponseWriter, req *http.Request) {
	var body entity.CreateTestOfClassRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTest, err := entity.ParseNewTestOfClass(body)
	if err != nil {
		pkg.SendError(w, "Failed to create test of class entity", http.StatusInternalServerError)
		return
	}

	userID, _ := req.Context().Value(constant.CtxUserID).(string)
	email, ok := req.Context().Value(constant.CtxEmail).(string)
	if !ok {
		pkg.SendError(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	newTest.UserID = userID
	newTest.AuthorMail = email

	insertedID, err := r.testOfClassUseCase.CreateTestOfClass(req.Context(), newTest)
	if err != nil {
		pkg.SendError(w, "Create test of class failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusCreated, map[string]interface{}{
		"_id":  insertedID,
		"test": newTest,
	})
}

func (r *RoutesTestOfClass) updateTestOfClass(w http.ResponseWriter, req *http.Request) {
	userID, _ := req.Context().Value(constant.CtxUserID).(string)
	email, _ := req.Context().Value(constant.CtxEmail).(string)

	var body entity.UpdateTestOfClassRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateTest, err := entity.ParseUpdateTestOfClass(body)
	if err != nil {
		pkg.SendError(w, "Failed to parse test of class entity", http.StatusInternalServerError)
		return
	}

	updateTest.UserID = userID
	updateTest.AuthorMail = email

	result, err := r.testOfClassUseCase.UpdateTestOfClass(req.Context(), updateTest)
	if err != nil {
		pkg.SendError(w, "Update test of class failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, result)
}

func (r *RoutesTestOfClass) deleteTestOfClass(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	var payload struct {
		ID primitive.ObjectID `json:"_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := r.testOfClassUseCase.DeleteTestOfClass(req.Context(), userID, payload.ID); err != nil {
		pkg.SendError(w, "Delete test of class failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, map[string]string{
		"message": "Test of class deleted successfully",
	})
}

func (r *RoutesTestOfClass) getTestsOfClass(w http.ResponseWriter, req *http.Request) {
	classIDStr := req.URL.Query().Get("class_id")
	if classIDStr == "" {
		pkg.SendError(w, "class_id query param required", http.StatusBadRequest)
		return
	}

	classID, err := primitive.ObjectIDFromHex(classIDStr)
	if err != nil {
		pkg.SendError(w, "Invalid class_id", http.StatusBadRequest)
		return
	}

	tests, err := r.testOfClassUseCase.GetTestsOfClass(req.Context(), classID)
	if err != nil {
		pkg.SendError(w, "Failed to get tests of class", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, tests)
}
