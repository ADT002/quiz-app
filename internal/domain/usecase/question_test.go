package usecase

import (
	"errors"
	"testing"

	entity "quiz-app/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newID() primitive.ObjectID { return primitive.NewObjectID() }

func mkOption(correct bool) entity.Option {
	return entity.Option{ID: newID(), Text: entity.Text{Text: "x"}, IsCorrect: correct}
}

func TestValidateQuestion_Single(t *testing.T) {
	t.Run("requires exactly 1 correct option", func(t *testing.T) {
		q := &entity.Question{
			Type:    SingleQuestion,
			Options: []entity.Option{mkOption(true), mkOption(false)},
		}
		if err := ValidateQuestion(q); err != nil {
			t.Fatalf("expected ok, got %v", err)
		}
	})
	t.Run("rejects 0 correct", func(t *testing.T) {
		q := &entity.Question{
			Type:    SingleQuestion,
			Options: []entity.Option{mkOption(false), mkOption(false)},
		}
		err := ValidateQuestion(q)
		if !errors.Is(err, ErrInvalidQuestion) {
			t.Fatalf("expected ErrInvalidQuestion, got %v", err)
		}
	})
	t.Run("rejects 2 correct", func(t *testing.T) {
		q := &entity.Question{
			Type:    SingleQuestion,
			Options: []entity.Option{mkOption(true), mkOption(true)},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
	t.Run("rejects < 2 options", func(t *testing.T) {
		q := &entity.Question{
			Type:    SingleQuestion,
			Options: []entity.Option{mkOption(true)},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
}

func TestValidateQuestion_Multiple(t *testing.T) {
	t.Run("accepts ≥ 1 correct", func(t *testing.T) {
		q := &entity.Question{
			Type:    MultiQuestion,
			Options: []entity.Option{mkOption(true), mkOption(true), mkOption(false)},
		}
		if err := ValidateQuestion(q); err != nil {
			t.Fatalf("expected ok, got %v", err)
		}
	})
	t.Run("rejects 0 correct", func(t *testing.T) {
		q := &entity.Question{
			Type:    MultiQuestion,
			Options: []entity.Option{mkOption(false), mkOption(false)},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
}

func TestValidateQuestion_Fill(t *testing.T) {
	t.Run("accepts blanks with correct_submission", func(t *testing.T) {
		q := &entity.Question{
			Type: FillInTheBlankQuestion,
			FillInTheBlanks: []entity.FillInTheBlank{
				{ID: newID(), CorrectSubmission: "Paris"},
			},
		}
		if err := ValidateQuestion(q); err != nil {
			t.Fatalf("expected ok, got %v", err)
		}
	})
	t.Run("rejects empty correct_submission", func(t *testing.T) {
		q := &entity.Question{
			Type: FillInTheBlankQuestion,
			FillInTheBlanks: []entity.FillInTheBlank{
				{ID: newID(), CorrectSubmission: ""},
			},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
	t.Run("rejects whitespace-only", func(t *testing.T) {
		q := &entity.Question{
			Type: FillInTheBlankQuestion,
			FillInTheBlanks: []entity.FillInTheBlank{
				{ID: newID(), CorrectSubmission: "   "},
			},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
	t.Run("rejects no blanks", func(t *testing.T) {
		q := &entity.Question{Type: FillInTheBlankQuestion}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
}

func TestValidateQuestion_Order(t *testing.T) {
	t.Run("accepts unique order values", func(t *testing.T) {
		q := &entity.Question{
			Type: OrderQuestion,
			OrderItems: []entity.OrderItem{
				{ID: newID(), Order: 1},
				{ID: newID(), Order: 2},
				{ID: newID(), Order: 3},
			},
		}
		if err := ValidateQuestion(q); err != nil {
			t.Fatalf("expected ok, got %v", err)
		}
	})
	t.Run("rejects duplicate order", func(t *testing.T) {
		q := &entity.Question{
			Type: OrderQuestion,
			OrderItems: []entity.OrderItem{
				{ID: newID(), Order: 1},
				{ID: newID(), Order: 1},
			},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
	t.Run("rejects < 2 items", func(t *testing.T) {
		q := &entity.Question{
			Type:       OrderQuestion,
			OrderItems: []entity.OrderItem{{ID: newID(), Order: 1}},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
}

func TestValidateQuestion_Match(t *testing.T) {
	itemID := newID()
	t.Run("accepts options referencing valid match_id", func(t *testing.T) {
		q := &entity.Question{
			Type: MatchQuestion,
			MatchItems: []entity.MatchItem{
				{ID: itemID, Text: entity.Text{Text: "Animal"}},
			},
			MatchOptions: []entity.MatchOption{
				{ID: newID(), MatchId: itemID.Hex()},
			},
		}
		if err := ValidateQuestion(q); err != nil {
			t.Fatalf("expected ok, got %v", err)
		}
	})
	t.Run("rejects option with unknown match_id", func(t *testing.T) {
		q := &entity.Question{
			Type:       MatchQuestion,
			MatchItems: []entity.MatchItem{{ID: itemID}},
			MatchOptions: []entity.MatchOption{
				{ID: newID(), MatchId: "deadbeefdeadbeefdeadbeef"},
			},
		}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
	t.Run("rejects empty items or options", func(t *testing.T) {
		q := &entity.Question{Type: MatchQuestion}
		if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
			t.Fatal("expected error")
		}
	})
}

func TestValidateQuestion_UnknownType(t *testing.T) {
	q := &entity.Question{Type: "not_a_real_type"}
	if !errors.Is(ValidateQuestion(q), ErrInvalidQuestion) {
		t.Fatal("expected error")
	}
}

// Regression test: ensure constants match CLAUDE.md D2 + NestJS scoring keys.
// A typo here means quiz-app Go writes questions NestJS cannot score.
func TestQuestionTypeConstants(t *testing.T) {
	cases := []struct{ got, want string }{
		{SingleQuestion, "single"},
		{MultiQuestion, "multiple"},
		{FillInTheBlankQuestion, "fill_in_the_blank"},
		{OrderQuestion, "order_question"},
		{MatchQuestion, "match_choice_question"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("constant %q != canonical %q", c.got, c.want)
		}
	}
}
