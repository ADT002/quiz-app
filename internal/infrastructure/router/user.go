package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"quiz-app/internal/auth"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"
)

type RoutesUser struct {
	auth        auth.AuthHandler
	userUseCase usecase.UserUseCase
}

func NewRouterUser(userUseCase usecase.UserUseCase, auth auth.AuthHandler) RoutesUser {
	return RoutesUser{
		userUseCase: userUseCase,
		auth:        auth,
	}
}

func (rt RoutesUser) GetUserRouter(r *Router) {
	r.Handle(
		"/user/me",
		rt.auth.AuthMiddleware(http.HandlerFunc(rt.getMe)),
	).Methods(http.MethodGet)

	r.Handle(
		"/user/me",
		rt.auth.AuthMiddleware(http.HandlerFunc(rt.updateMe)),
	).Methods(http.MethodPatch)
}

func (r *RoutesUser) getMe(w http.ResponseWriter, req *http.Request) {
	token, ok := bearerToken(req)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := r.userUseCase.GetMe(req.Context(), token)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	pkg.SendResponse(w, http.StatusOK, user)
}

func (r *RoutesUser) updateMe(w http.ResponseWriter, req *http.Request) {
	token, ok := bearerToken(req)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body entity.UpdateMeRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		pkg.SendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := r.userUseCase.UpdateMe(req.Context(), token, body)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	pkg.SendResponse(w, http.StatusOK, user)
}

// bearerToken extracts the raw JWT from the Authorization header.
// AuthMiddleware already validated it, so any failure here is a programming bug.
func bearerToken(req *http.Request) (string, bool) {
	h := req.Header.Get("Authorization")
	if h == "" {
		return "", false
	}
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	return token, token != ""
}

// writeUpstreamError forwards the auth-service status code when available so
// the FE sees the same code (401/403/404/...) as if it called auth-service directly.
func writeUpstreamError(w http.ResponseWriter, err error) {
	if status := usecase.UpstreamStatus(err); status >= 400 {
		pkg.SendError(w, err.Error(), status)
		return
	}
	pkg.SendError(w, err.Error(), http.StatusBadGateway)
}
