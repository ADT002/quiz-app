package usecase

import (
	"context"
	"errors"
	"fmt"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Canonical question type names. MUST match CLAUDE.md D2 + NestJS scoring keys.
// Vi phạm = NestJS không chấm được câu hỏi quiz-app Go tạo (silent bug).
const (
	SingleQuestion         = "single"
	MultiQuestion          = "multiple"
	OrderQuestion          = "order_question"
	MatchQuestion          = "match_choice_question"
	FillInTheBlankQuestion = "fill_in_the_blank"
)

// ErrInvalidQuestion is returned when validation fails. Handler maps to 400.
var ErrInvalidQuestion = errors.New("invalid question")

type QuestionUseCase struct {
	QuestionRepo repository.QuestionRepository
}

func NewQuestionUseCase(tr repository.QuestionRepository) *QuestionUseCase {
	return &QuestionUseCase{QuestionRepo: tr}
}

// QuestionFilters re-exports the repository type so handlers can keep the
// usecase package as the single import surface.
type QuestionFilters = repository.QuestionFilters

func (uc *QuestionUseCase) CreateQuestion(ctx context.Context, q entity.Question) (primitive.ObjectID, error) {
	ensureItemIDs(&q)
	if err := ValidateQuestion(&q); err != nil {
		return primitive.NilObjectID, err
	}
	return uc.QuestionRepo.CreateQuestion(ctx, q)
}

func (uc *QuestionUseCase) UpdateQuestion(ctx context.Context, q entity.Question) (any, error) {
	ensureItemIDs(&q)
	if err := ValidateQuestion(&q); err != nil {
		return nil, err
	}
	return uc.QuestionRepo.UpdateQuestion(ctx, q)
}

func (uc *QuestionUseCase) ListQuestions(ctx context.Context, f QuestionFilters) ([]entity.Question, int64, error) {
	return uc.QuestionRepo.ListQuestions(ctx, f)
}

func (uc *QuestionUseCase) GetAllQuestionsOfTest(ctx context.Context, ids []primitive.ObjectID) ([]entity.Question, error) {
	return uc.QuestionRepo.GetAllQuestions(ctx, ids)
}

func (uc *QuestionUseCase) DeleteQuestion(ctx context.Context, q entity.Question) error {
	return uc.QuestionRepo.DeleteQuestion(ctx, q)
}

// ValidateQuestion enforces per-type invariants per CLAUDE.md D2.
// Caller is expected to have called ensureItemIDs first so item IDs exist.
func ValidateQuestion(q *entity.Question) error {
	switch q.Type {
	case SingleQuestion:
		correct := 0
		for _, o := range q.Options {
			if o.IsCorrect {
				correct++
			}
		}
		if len(q.Options) < 2 {
			return wrap("single requires ≥ 2 options")
		}
		if correct != 1 {
			return wrap("single requires exactly 1 correct option, got " + strconv.Itoa(correct))
		}

	case MultiQuestion:
		correct := 0
		for _, o := range q.Options {
			if o.IsCorrect {
				correct++
			}
		}
		if len(q.Options) < 2 {
			return wrap("multiple requires ≥ 2 options")
		}
		if correct < 1 {
			return wrap("multiple requires ≥ 1 correct option")
		}

	case FillInTheBlankQuestion:
		if len(q.FillInTheBlanks) == 0 {
			return wrap("fill_in_the_blank requires ≥ 1 blank")
		}
		for i, b := range q.FillInTheBlanks {
			if strings.TrimSpace(b.CorrectSubmission) == "" {
				return wrap(fmt.Sprintf("blank #%d missing correct_submission", i+1))
			}
		}

	case OrderQuestion:
		if len(q.OrderItems) < 2 {
			return wrap("order_question requires ≥ 2 items")
		}
		seen := map[int]bool{}
		for _, it := range q.OrderItems {
			if seen[it.Order] {
				return wrap("order_question items must have unique order values")
			}
			seen[it.Order] = true
		}

	case MatchQuestion:
		if len(q.MatchItems) == 0 || len(q.MatchOptions) == 0 {
			return wrap("match_choice_question requires ≥ 1 item and ≥ 1 option")
		}
		itemIDs := map[string]bool{}
		for _, it := range q.MatchItems {
			itemIDs[it.ID.Hex()] = true
		}
		for _, opt := range q.MatchOptions {
			if !itemIDs[opt.MatchId] {
				return wrap("match_choice_question option references unknown match_id: " + opt.MatchId)
			}
		}

	default:
		return wrap("unknown question type: " + q.Type)
	}
	return nil
}

func wrap(msg string) error {
	return fmt.Errorf("%w: %s", ErrInvalidQuestion, msg)
}

// ensureItemIDs assigns ObjectIDs to any item with NilObjectID. Stable IDs are
// required by CLAUDE.md E#7 — option/blank/item IDs are server-issued and never
// derived from index.
func ensureItemIDs(q *entity.Question) {
	for i := range q.Options {
		if q.Options[i].ID.IsZero() {
			q.Options[i].ID = primitive.NewObjectID()
		}
	}
	for i := range q.FillInTheBlanks {
		if q.FillInTheBlanks[i].ID.IsZero() {
			q.FillInTheBlanks[i].ID = primitive.NewObjectID()
		}
	}
	for i := range q.OrderItems {
		if q.OrderItems[i].ID.IsZero() {
			q.OrderItems[i].ID = primitive.NewObjectID()
		}
	}
	for i := range q.MatchItems {
		if q.MatchItems[i].ID.IsZero() {
			q.MatchItems[i].ID = primitive.NewObjectID()
		}
	}
	for i := range q.MatchOptions {
		if q.MatchOptions[i].ID.IsZero() {
			q.MatchOptions[i].ID = primitive.NewObjectID()
		}
	}
}
