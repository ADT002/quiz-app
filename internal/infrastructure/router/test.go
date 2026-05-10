package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesTestTemplate struct {
	auth                auth.AuthHandler
	testTemplateUseCase usecase.TestTemplateUseCase
}

func NewRouterTestTemplate(testTemplateUseCase usecase.TestTemplateUseCase, auth auth.AuthHandler) RoutesTestTemplate {
	return RoutesTestTemplate{
		testTemplateUseCase: testTemplateUseCase,
		auth:                auth,
	}
}

func (rt RoutesTestTemplate) GetTestTemplateRouter(r *Router) {
	r.Handle(
		"/test-templates",
		rt.auth.AuthWithPerm("TEACHER", rt.getTestTemplates),
	).Methods("GET")

	r.Handle(
		"/test-templates",
		rt.auth.AuthWithPerm("TEACHER", rt.createTestTemplate),
	).Methods("POST")

	r.Handle(
		"/test-templates",
		rt.auth.AuthWithPerm("TEACHER", rt.updateTestTemplate),
	).Methods("PATCH")

	r.Handle(
		"/test-templates",
		rt.auth.AuthWithPerm("TEACHER", rt.deleteTestTemplate),
	).Methods("DELETE")
}

func (r *RoutesTestTemplate) createTestTemplate(w http.ResponseWriter, req *http.Request) {
	var body entity.CreateTestRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTest, err := entity.ParseNewTest(body)
	if err != nil {
		fmt.Println(err)
		pkg.SendError(w, "Failed to create test template entity", http.StatusInternalServerError)
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

	insertedID, err := r.testTemplateUseCase.CreateTestTemplate(req.Context(), newTest)
	if err != nil {
		pkg.SendError(w, "Create test template failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusCreated, map[string]interface{}{
		"_id":  insertedID,
		"test": newTest,
	})
}

func (r *RoutesTestTemplate) updateTestTemplate(w http.ResponseWriter, req *http.Request) {
	userID, _ := req.Context().Value(constant.CtxUserID).(string)
	email, _ := req.Context().Value(constant.CtxEmail).(string)

	var body entity.UpdateTestRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		fmt.Println(err)
		pkg.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateTest, err := entity.ParseUpdateTest(body)
	if err != nil {
		pkg.SendError(w, "Failed to parse test template entity", http.StatusInternalServerError)
		return
	}

	updateTest.UserID = userID
	updateTest.AuthorMail = email

	result, err := r.testTemplateUseCase.UpdateTestTemplate(req.Context(), updateTest)
	if err != nil {
		pkg.SendError(w, "Update test template failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, result)
}

func (r *RoutesTestTemplate) deleteTestTemplate(w http.ResponseWriter, req *http.Request) {
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

	if err := r.testTemplateUseCase.DeleteTestTemplate(req.Context(), userID, payload.ID); err != nil {
		pkg.SendError(w, "Delete test template failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, map[string]string{
		"message": "Test template deleted successfully",
	})
}

func (r *RoutesTestTemplate) getTestTemplates(w http.ResponseWriter, req *http.Request) {
	userID, ok := req.Context().Value(constant.CtxUserID).(string)
	if !ok {
		pkg.SendError(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	tests, err := r.testTemplateUseCase.GetTestTemplates(req.Context(), userID)
	if err != nil {
		pkg.SendError(w, "Failed to get test templates", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, tests)
}
