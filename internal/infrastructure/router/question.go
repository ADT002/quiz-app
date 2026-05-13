package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"quiz-app/internal/auth"
	"quiz-app/internal/constant"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/pkg"
	utils "quiz-app/internal/util"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutesQuestion struct {
	auth            auth.AuthHandler
	questionUseCase usecase.QuestionUseCase
	topicUseCase    usecase.TopicUseCase
	levelUseCase    usecase.LevelUseCase
}

func NewRouterQuestion(
	questionUseCase usecase.QuestionUseCase,
	topicUseCase usecase.TopicUseCase,
	levelUseCase usecase.LevelUseCase,
	auth auth.AuthHandler,
) *RoutesQuestion {
	return &RoutesQuestion{
		questionUseCase: questionUseCase,
		topicUseCase:    topicUseCase,
		levelUseCase:    levelUseCase,
		auth:            auth,
	}
}

func (rq *RoutesQuestion) GetQuestionRouter(r *Router) {
	// ===== QUESTIONS (TEACHER) =====
	r.Handle("/questions", rq.auth.AuthWithPerm("TEACHER", rq.createQuestion)).Methods(http.MethodPost)
	r.Handle("/questions", rq.auth.AuthWithPerm("TEACHER", rq.listQuestions)).Methods(http.MethodGet)
	r.Handle("/questions", rq.auth.AuthWithPerm("TEACHER", rq.updateQuestion)).Methods(http.MethodPatch)
	r.Handle("/questions/{id}", rq.auth.AuthWithPerm("TEACHER", rq.deleteQuestion)).Methods(http.MethodDelete)

	// ===== TOPIC (TEACHER) =====
	r.Handle("/topic", rq.auth.AuthWithPerm("TEACHER", rq.createTopic)).Methods(http.MethodPost)
	r.Handle("/topic", rq.auth.AuthWithPerm("TEACHER", rq.getAllTopics)).Methods(http.MethodGet)
	r.Handle("/topic", rq.auth.AuthWithPerm("TEACHER", rq.updateTopic)).Methods(http.MethodPatch)
	r.Handle("/topic", rq.auth.AuthWithPerm("TEACHER", rq.deleteTopic)).Methods(http.MethodDelete)
	r.Handle("/topic/reorder", rq.auth.AuthWithPerm("TEACHER", rq.reorderTopics)).Methods(http.MethodPost)

	// ===== LEVEL (TEACHER) =====
	r.Handle("/level", rq.auth.AuthWithPerm("TEACHER", rq.createLevel)).Methods(http.MethodPost)
	r.Handle("/level", rq.auth.AuthWithPerm("TEACHER", rq.getAllLevels)).Methods(http.MethodGet)
	r.Handle("/level", rq.auth.AuthWithPerm("TEACHER", rq.updateLevel)).Methods(http.MethodPatch)
	r.Handle("/level", rq.auth.AuthWithPerm("TEACHER", rq.deleteLevel)).Methods(http.MethodDelete)
}

/* ── QUESTION ─────────────────────────────────────────────────────────── */

// paginatedQuestions matches CLAUDE.md E#10: { items, total, page, limit }.
type paginatedQuestions struct {
	Items []entity.Question `json:"items"`
	Total int64             `json:"total"`
	Page  int64             `json:"page"`
	Limit int64             `json:"limit"`
}

func (rq *RoutesQuestion) createQuestion(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var q entity.Question
	if err := json.NewDecoder(req.Body).Decode(&q); err != nil {
		pkg.SendError(w, "invalid body", http.StatusBadRequest)
		return
	}

	uuidQuestion, _ := uuid.Parse(userID)
	q.Metadata.User_ID = uuidQuestion.String()
	q.ID = primitive.NewObjectID()
	q.Created_At = time.Now()
	q.Updated_At = time.Now()

	id, err := rq.questionUseCase.CreateQuestion(context.TODO(), q)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidQuestion) {
			pkg.SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		pkg.SendError(w, "create failed", http.StatusInternalServerError)
		return
	}

	q.ID = id
	pkg.SendResponse(w, http.StatusCreated, q)
}

func (rq *RoutesQuestion) listQuestions(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	q := req.URL.Query()
	limit, _ := strconv.ParseInt(q.Get("limit"), 10, 64)
	page, _ := strconv.ParseInt(q.Get("page"), 10, 64)

	filters := usecase.QuestionFilters{
		OwnerID: userID,
		TopicID: q.Get("topic_id"),
		LevelID: q.Get("level_id"),
		Type:    q.Get("type"),
		Search:  q.Get("q"),
		Page:    page,
		Limit:   limit,
	}

	items, total, err := rq.questionUseCase.ListQuestions(context.TODO(), filters)
	if err != nil {
		pkg.SendError(w, "query failed", http.StatusInternalServerError)
		return
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 100 {
		filters.Limit = 20
	}
	if items == nil {
		items = []entity.Question{}
	}
	pkg.SendResponse(w, http.StatusOK, paginatedQuestions{
		Items: items,
		Total: total,
		Page:  filters.Page,
		Limit: filters.Limit,
	})
}

func (rq *RoutesQuestion) updateQuestion(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var q entity.Question
	if err := json.NewDecoder(req.Body).Decode(&q); err != nil {
		pkg.SendError(w, "invalid body", http.StatusBadRequest)
		return
	}
	uuidQuestion, _ := uuid.Parse(userID)
	q.Metadata.User_ID = uuidQuestion.String()
	q.Updated_At = time.Now()

	updated, err := rq.questionUseCase.UpdateQuestion(context.TODO(), q)
	if err != nil {
		fmt.Println(err)

		if errors.Is(err, usecase.ErrInvalidQuestion) {
			pkg.SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		pkg.SendError(w, "update failed", http.StatusInternalServerError)
		return
	}

	pkg.SendResponse(w, http.StatusOK, updated)
}

func (rq *RoutesQuestion) deleteQuestion(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idHex := mux.Vars(req)["id"]
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		pkg.SendError(w, "invalid id", http.StatusBadRequest)
		return
	}

	uuidQuestion, _ := uuid.Parse(userID)
	q := entity.Question{
		ID:       id,
		Metadata: entity.Metadata{User_ID: uuidQuestion.String()},
	}
	force := req.URL.Query().Get("force") == "true"
	usage, err := rq.questionUseCase.DeleteQuestion(req.Context(), q, force)
	if err != nil {
		if errors.Is(err, usecase.ErrQuestionInUse) {
			// 409: tell FE where the question is still referenced so it can
			// show a confirm modal that retries with ?force=true.
			pkg.SendResponse(w, http.StatusConflict, map[string]any{
				"code":    "QUESTION_IN_USE",
				"message": "câu hỏi đang được sử dụng trong bài thi",
				"usage":   usage,
			})
			return
		}
		pkg.SendError(w, "delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/* ── TOPIC ────────────────────────────────────────────────────────────── */

func (rq *RoutesQuestion) createTopic(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var topic entity.Topic
	if err := json.NewDecoder(req.Body).Decode(&topic); err != nil {
		pkg.SendError(w, "invalid body", http.StatusBadRequest)
		return
	}

	topic.UserID = userID
	result, err := rq.topicUseCase.CreateTopic(context.TODO(), topic)
	if err != nil {
		pkg.SendError(w, "create failed", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, result)
}

func (rq *RoutesQuestion) getAllTopics(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := rq.topicUseCase.GetAllTopics(context.TODO(), userID)
	if err != nil {
		pkg.SendError(w, "Query fail", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, result)
}

func (rq *RoutesQuestion) updateTopic(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var topic entity.Topic
	if err := json.NewDecoder(req.Body).Decode(&topic); err != nil {
		pkg.SendError(w, "Topic not decode", http.StatusBadRequest)
		return
	}
	topic.UserID = userID

	result, err := rq.topicUseCase.UpdateTopic(context.TODO(), topic)
	if err != nil {
		pkg.SendError(w, "Update fail", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, result)
}

func (rq *RoutesQuestion) reorderTopics(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload struct {
		OrderedIDs []string `json:"ordered_ids"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		pkg.SendError(w, "invalid body", http.StatusBadRequest)
		return
	}
	if len(payload.OrderedIDs) == 0 {
		pkg.SendError(w, "ordered_ids required", http.StatusBadRequest)
		return
	}
	if err := rq.topicUseCase.ReorderTopics(req.Context(), userID, payload.OrderedIDs); err != nil {
		pkg.SendError(w, "reorder failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rq *RoutesQuestion) deleteTopic(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		pkg.SendError(w, "Topic not decode", http.StatusBadRequest)
		return
	}
	topicID, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		pkg.SendError(w, "Invalid topic id", http.StatusBadRequest)
		return
	}
	topic := entity.Topic{ID: topicID, UserID: userID}
	if err := rq.topicUseCase.DeleteTopic(req.Context(), topic); err != nil {
		if errors.Is(err, usecase.ErrTopicInUse) {
			pkg.SendResponse(w, http.StatusConflict, map[string]any{
				"code":    "TOPIC_IN_USE",
				"message": "chủ đề đang được gán cho câu hỏi, không thể xoá",
			})
			return
		}
		pkg.SendError(w, "Topic not deleted", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, map[string]string{
		"message": "Topic deleted successfully",
	})
}

/* ── LEVEL ────────────────────────────────────────────────────────────── */

func (rq *RoutesQuestion) createLevel(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var level entity.Level
	if err := json.NewDecoder(req.Body).Decode(&level); err != nil {
		pkg.SendError(w, "Level not decode", http.StatusBadRequest)
		return
	}
	level.UserID = userID

	result, err := rq.levelUseCase.CreateLevel(context.TODO(), level)
	if err != nil {
		pkg.SendError(w, "Level not created", http.StatusInternalServerError)
		return
	}
	level.ID = result
	pkg.SendResponse(w, http.StatusOK, level)
}

func (rq *RoutesQuestion) getAllLevels(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	result, err := rq.levelUseCase.GetAllLevels(context.TODO(), userID)
	if err != nil {
		pkg.SendError(w, "Query fail", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, result)
}

func (rq *RoutesQuestion) updateLevel(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var level entity.Level
	if err := json.NewDecoder(req.Body).Decode(&level); err != nil {
		pkg.SendError(w, "Level not decode", http.StatusBadRequest)
		return
	}
	level.UserID = userID

	result, err := rq.levelUseCase.UpdateLevel(context.TODO(), level)
	if err != nil {
		pkg.SendError(w, "Level not update", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, result)
}

func (rq *RoutesQuestion) deleteLevel(w http.ResponseWriter, req *http.Request) {
	userID, ok := utils.GetStringFromCtx(req.Context(), constant.CtxUserID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload struct {
		ID string `json:"_id"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		pkg.SendError(w, "Level not decode", http.StatusBadRequest)
		return
	}
	levelID, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		pkg.SendError(w, "Invalid Level id", http.StatusBadRequest)
		return
	}
	level := entity.Level{ID: levelID, UserID: userID}
	if err := rq.levelUseCase.DeleteLevel(req.Context(), level); err != nil {
		if errors.Is(err, usecase.ErrLevelInUse) {
			pkg.SendResponse(w, http.StatusConflict, map[string]any{
				"code":    "LEVEL_IN_USE",
				"message": "độ khó đang được gán cho câu hỏi, không thể xoá",
			})
			return
		}
		pkg.SendError(w, "Level not deleted", http.StatusInternalServerError)
		return
	}
	pkg.SendResponse(w, http.StatusOK, map[string]string{
		"message": "Level deleted successfully",
	})
}
