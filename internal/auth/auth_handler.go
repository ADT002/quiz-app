package auth

import (
	"context"
	"net/http"
	"strings"

	"quiz-app/internal/constant"
	utils "quiz-app/internal/util"
)

type AuthHandler struct {
	authUseCase *AuthUseCase
}

func NewAuthHandler(authUseCase *AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}
func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format, expected Bearer <token>", http.StatusUnauthorized)
			return
		}

		claims, err := h.authUseCase.Authenticate(r.Context(), tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), constant.CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, constant.CtxEmailID, claims.UserID)
		ctx = context.WithValue(ctx, constant.CtxEmail, claims.AuthorMail)
		ctx = context.WithValue(ctx, constant.CtxPermission, claims.Permission)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AuthHandler) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perms, ok := r.Context().Value(constant.CtxPermission).([]string)
			if !ok {
				http.Error(w, "permission not found", http.StatusForbidden)
				return
			}

			if !utils.HasPermission(perms, permission) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *AuthHandler) AuthWithPerm(
	perm string,
	handler http.HandlerFunc,
) http.Handler {
	return h.AuthMiddleware(
		h.RequirePermission(perm)(handler),
	)
}
