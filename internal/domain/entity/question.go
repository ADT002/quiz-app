package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question represents the main question structure.
//
// File references use file_id (ObjectID) per CLAUDE.md E#17 — never store
// raw filenames or owner email here.
type Question struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Type            string             `json:"type" bson:"type,omitempty"` // see usecase/question.go for canonical type names
	QuestionContent QuestionContent    `json:"question_content,omitempty" bson:"question_content,omitempty"`
	Metadata        Metadata           `json:"metadata,omitempty" bson:"metadata,omitempty"`
	Tags            []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	Suggestion      []string           `json:"suggestion,omitempty" bson:"suggestion,omitempty"`
	Score           float32            `json:"score,omitempty" bson:"score,omitempty"`

	QuestionLevel primitive.ObjectID `json:"level" bson:"level"`
	QuestionTopic primitive.ObjectID `json:"topic" bson:"topic"`

	Options         []Option         `json:"options,omitempty" bson:"options,omitempty"`
	FillInTheBlanks []FillInTheBlank `json:"fill_in_the_blanks,omitempty" bson:"fill_in_the_blanks,omitempty"`
	OrderItems      []OrderItem      `json:"order_items,omitempty" bson:"order_items,omitempty"`
	MatchItems      []MatchItem      `json:"match_items,omitempty" bson:"match_items,omitempty"`
	MatchOptions    []MatchOption    `json:"match_options,omitempty" bson:"match_options,omitempty"`

	// Files attached to the question body (illustrations, audio prompt, …).
	// Each entry is a reference to the Files collection by file_id.
	Files []FileRef `json:"files,omitempty" bson:"files,omitempty"`

	Created_At time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Updated_At time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Metadata struct {
	User_ID string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

// FileRef is a typed reference into the `files` collection.
// Stable across rename — see CLAUDE.md E#17.
type FileRef struct {
	FileID primitive.ObjectID `json:"file_id" bson:"file_id"`
}

type MatchItem struct {
	ID   primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text Text               `json:"text" bson:"text,omitempty"`
}

type MatchOption struct {
	ID      primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text    Text               `json:"text" bson:"text,omitempty"`
	MatchId string             `json:"match_id" bson:"match_id,omitempty"`
}

// FillInTheBlank — text_before/text_after wrap rich text for LaTeX support.
// `correct_submission` is server-only (stripped before serving to student).
type FillInTheBlank struct {
	ID                primitive.ObjectID `json:"id" bson:"id,omitempty"`
	TextBefore        Text               `json:"text_before" bson:"text_before,omitempty"`
	TextAfter         Text               `json:"text_after" bson:"text_after,omitempty"`
	Blank             string             `json:"blank" bson:"blank,omitempty"`
	CorrectSubmission string             `json:"correct_submission" bson:"correct_submission,omitempty"`
}

// Option for single/multiple choice. `is_correct` is server-only — stripped
// before serving to student. Optional `file_id` references an image/audio.
type Option struct {
	ID        primitive.ObjectID  `json:"id,omitempty" bson:"id,omitempty"`
	Text      Text                `json:"text,omitempty" bson:"text,omitempty"`
	FileID    *primitive.ObjectID `json:"file_id,omitempty" bson:"file_id,omitempty"`
	IsCorrect bool                `json:"is_correct,omitempty" bson:"is_correct,omitempty"`
}

type OrderItem struct {
	ID    primitive.ObjectID `json:"id" bson:"id,omitempty"`
	Text  Text               `json:"text" bson:"text,omitempty"`
	Order int                `json:"order" bson:"order,omitempty"` // server-only, stripped from student view
}

// QuestionContent — body text + optional file (image/audio/video) attached.
// Replaces legacy `file_url` (filename string) per E#17.
type QuestionContent struct {
	Text   Text                `json:"content" bson:"content,omitempty"`
	FileID *primitive.ObjectID `json:"file_id,omitempty" bson:"file_id,omitempty"`
}

type Text struct {
	IsMath bool   `json:"is_math" bson:"is_math,omitempty"`
	Text   string `json:"text" bson:"text,omitempty"`
}
