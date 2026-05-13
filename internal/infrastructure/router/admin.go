package routes

import (
	"encoding/json"
	"net/http"

	"quiz-app/internal/auth"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"

	"quiz-app/internal/pkg"

	"github.com/gorilla/mux"
)

type RouterAdmin struct {
	auth         auth.AuthHandler
	AdminUseCase usecase.UserUseCase
}

func NewRouterAdmin(userUC usecase.UserUseCase, auth auth.AuthHandler) *RouterAdmin {
	return &RouterAdmin{
		auth:         auth,
		AdminUseCase: userUC,
	}
}

func (rc RouterAdmin) getAllUser(w http.ResponseWriter, req *http.Request) {
	token, ok := bearerToken(req)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := req.URL.Query()
	page := usecase.ParseIntDefault(q.Get("page"), 1)
	limit := usecase.ParseIntDefault(q.Get("limit"), 20)

	users, err := rc.AdminUseCase.ListUsers(req.Context(), token, page, limit)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	pkg.SendResponse(w, http.StatusOK, users)
}

func (rc RouterAdmin) updateUser(w http.ResponseWriter, req *http.Request) {
	token, ok := bearerToken(req)
	if !ok {
		pkg.SendError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := mux.Vars(req)["id"]
	if id == "" {
		pkg.SendError(w, "missing user id", http.StatusBadRequest)
		return
	}

	var body entity.AdminUpdateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		pkg.SendError(w, "invalid body", http.StatusBadRequest)
		return
	}

	user, err := rc.AdminUseCase.AdminUpdateUser(req.Context(), token, id, body)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	pkg.SendResponse(w, http.StatusOK, user)
}

func (rc RouterAdmin) GetAdminRouter(r *Router) {
	r.Router.Handle(
		"/user",
		rc.auth.AuthMiddleware(
			rc.auth.RequirePermission("ADMIN")(
				http.HandlerFunc(rc.getAllUser),
			),
		),
	).Methods(http.MethodGet)

	r.Router.Handle(
		"/user/{id}",
		rc.auth.AuthMiddleware(
			rc.auth.RequirePermission("ADMIN")(
				http.HandlerFunc(rc.updateUser),
			),
		),
	).Methods(http.MethodPatch)
}
