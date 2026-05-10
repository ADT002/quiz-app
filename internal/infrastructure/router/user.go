package routes

import (
	"encoding/json"
	"net/http"
	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"
)

type RoutesUser struct {
	auth        auth.AuthHandler
	userUseCase usecase.UserUseCase
}

// NewRouterUser creates a new RoutesUser instance
func NewRouterUser(userUseCase usecase.UserUseCase, auth auth.AuthHandler) RoutesUser {
	return RoutesUser{
		userUseCase: userUseCase,
		auth:        auth,
	}
}

func (rt RoutesUser) GetUserRouter(r *Router) {
	r.Handle(
		"/user/me",
		rt.auth.AuthMiddleware(http.HandlerFunc(rt.updateUser)),
	).Methods(http.MethodPut)

	r.Handle(
		"/user/me",
		rt.auth.AuthMiddleware(http.HandlerFunc(rt.getUser)),
	).Methods(http.MethodGet)
}

func (r *RoutesUser) getUser(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value(constant.CtxEmailID).(string)
	if !ok || emailID == "" {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := r.userUseCase.GetUser(req.Context(), emailID)
	if err != nil {
		pkg.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user = entity.User{
		FirstName:    result.FirstName,
		LastName:     result.LastName,
		EmailAddress: result.EmailAddress,
	}

	pkg.SendResponse(w, http.StatusOK, user)
}

func (r *RoutesUser) updateUser(w http.ResponseWriter, req *http.Request) {
	emailID, ok := req.Context().Value(constant.CtxEmailID).(string)
	if !ok || emailID == "" {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody entity.UpdateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		pkg.SendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := r.userUseCase.UpdateUserPartial(
		req.Context(),
		emailID,
		reqBody,
	)
	if err != nil {
		pkg.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, result)
}
