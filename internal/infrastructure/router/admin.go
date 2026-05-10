package routes

import (
	"encoding/json"
	"net/http"

	"quiz-app/internal/auth"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"

	"quiz-app/internal/pkg"
)

type RouterAdmin struct {
	auth         auth.AuthHandler
	AdminUseCase usecase.UserUseCase
	redisUseCase usecase.RedisUseCase
}
type UpdatePermissionRequest struct {
	AdminID    string   `json:"userId"`
	Permission []string `json:"permission"`
}

func NewRouterAdmin(UserUC usecase.UserUseCase, redisUC usecase.RedisUseCase, auth auth.AuthHandler) *RouterAdmin {
	return &RouterAdmin{
		auth:         auth,
		AdminUseCase: UserUC,
		redisUseCase: redisUC,
	}
}

func (rc RouterAdmin) getAllUser(w http.ResponseWriter, req *http.Request) {
	users, err := rc.AdminUseCase.GetAllUser(req.Context())

	if err != nil {
		pkg.SendError(w, "Failed to retrieve Admines", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, users)
}

func (rc RouterAdmin) updateUser(w http.ResponseWriter, req *http.Request) {
	var reqBody entity.User
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	updatedAdmin, err := rc.AdminUseCase.UpdateUser(req.Context(), reqBody)
	if err != nil {
		pkg.SendError(w, "Failed to update User", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, updatedAdmin)
}

func (rc RouterAdmin) GetAdminRouter(r *Router) {
	r.Router.Handle(
		"/user",
		rc.auth.AuthMiddleware(
			rc.auth.RequirePermission("ADMIN")(
				http.HandlerFunc(rc.getAllUser),
			),
		),
	).Methods("GET")

	// UPDATE USER - ADMIN
	r.Router.Handle(
		"/user",
		rc.auth.AuthMiddleware(
			rc.auth.RequirePermission("ADMIN")(
				http.HandlerFunc(rc.updateUser),
			),
		),
	).Methods("PATCH")

	// // DELETE USER - ADMIN
	// r.Router.Handle(
	// 	"/user",
	// 	rc.auth.AuthMiddleware(
	// 		rc.auth.RequirePermission("ADMIN")(
	// 			http.HandlerFunc(rc.deleteAdmin),
	// 		),
	// 	),
	// ).Methods("DELETE")
}
