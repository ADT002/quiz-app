package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatrixExam struct {
	Topic    string `json:"topic" bson:"topic"`
	Level    string `json:"level" bson:"level"`
	Quantity int    `json:"quantity" bson:"quantity"`
}

// ================= BASE STRUCT =================

type BaseTest struct {
	TestName string `json:"test_name" bson:"test_name"`
	Descript string `json:"descript,omitempty" bson:"descript,omitempty"`

	QuestionIDs []string     `json:"question_ids" bson:"question_ids"`
	TestScore   float32      `json:"test_score" bson:"test_score"`
	Matrix      []MatrixExam `json:"matrix_exam" bson:"matrix_exam"`

	IsTest bool `json:"is_test" bson:"is_test"`

	StartTime       string `json:"start_time" bson:"start_time"`
	EndTime         string `json:"end_time" bson:"end_time"`
	DurationMinutes int    `json:"duration_minutes" bson:"duration_minutes"`

	Tags []string `json:"tags,omitempty" bson:"tags,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

// ================= MAIN ENTITIES =================

// TestOfClassUserSubmit là bản tóm tắt bài đã nộp (đọc từ collection `submissions` do NestJS ghi).
// Không lưu trong document `test_of_class`; chỉ gắn khi GET list cho giáo viên.
type TestOfClassUserSubmit struct {
	UserEmail    string  `json:"user_email"`
	Score        float64 `json:"score"`
	EmailID      string  `json:"email_id"`
	SubmissionID string  `json:"submission_id,omitempty"`
}

type TestOfClass struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID     string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ClassID    primitive.ObjectID `json:"class_id,omitempty" bson:"class_id,omitempty"`
	AuthorMail string             `json:"author_mail" bson:"author_mail"`

	BaseTest `bson:",inline"`

	UserSubmit []TestOfClassUserSubmit `json:"user_submit,omitempty" bson:"-"`
}

// FinishedSubmissionRow — một dòng đọc từ Mongo `submissions` (camelCase từ Nest/Mongoose).
type FinishedSubmissionRow struct {
	TestOfClassID primitive.ObjectID `bson:"test_of_class_id"`
	SubmissionID  primitive.ObjectID `bson:"_id"`
	StudentID     string             `bson:"student_id"`
	Email         string             `bson:"email"`
	Score         float64            `bson:"score"`
	Status        string             `bson:"status"`
}

type TestTemplete struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID     string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AuthorMail string             `json:"author_mail" bson:"author_mail"`

	BaseTest `bson:",inline"`
}

// ================= USER SUBMIT =================

type UserSubmit struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	AuthorMail string             `json:"author_mail,omitempty" bson:"author_mail,omitempty"`
	ClassID    primitive.ObjectID `json:"class_id,omitempty" bson:"class_id,omitempty"`
	TestID     primitive.ObjectID `json:"test_id" bson:"test_id"`

	UserID  string `json:"user_id" bson:"user_id"`
	EmailID string `json:"email_id,omitempty" bson:"email_id,omitempty"`
	Score   int8   `json:"score,omitempty" bson:"score,omitempty"`
}

// ================= REQUEST: TestTemplete =================

type CreateTestRequest struct {
	UserID   primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	TestName string             `json:"test_name" bson:"test_name"`

	Descript string   `json:"descript,omitempty" bson:"descript,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`

	IsTest bool `json:"is_test" bson:"is_test"`

	Matrix      []MatrixExam `json:"matrix_exam" bson:"matrix_exam"`
	QuestionIDs []string     `json:"question_ids" bson:"question_ids"`

	StartTime       string `json:"start_time" bson:"start_time"`
	EndTime         string `json:"end_time" bson:"end_time"`
	DurationMinutes int    `json:"duration_minutes" bson:"duration_minutes"`

	TestScore float32 `json:"test_score" bson:"test_score"`
}

type UpdateTestRequest struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	UserID   string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	TestName string             `json:"test_name" bson:"test_name"`

	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`
	Descript string   `json:"descript,omitempty" bson:"descript,omitempty"`

	QuestionIDs []string     `json:"question_ids" bson:"question_ids"`
	Matrix      []MatrixExam `json:"matrix_exam" bson:"matrix_exam"`

	IsTest bool `json:"is_test" bson:"is_test"`

	StartTime       string `json:"start_time" bson:"start_time"`
	EndTime         string `json:"end_time" bson:"end_time"`
	DurationMinutes int    `json:"duration_minutes" bson:"duration_minutes"`

	TestScore float32 `json:"test_score" bson:"test_score"`
}

// ================= REQUEST: TestOfClass =================

type CreateTestOfClassRequest struct {
	ClassID  primitive.ObjectID `json:"class_id" bson:"class_id"`
	TestName string             `json:"test_name" bson:"test_name"`

	Descript string   `json:"descript,omitempty" bson:"descript,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`

	IsTest bool `json:"is_test" bson:"is_test"`

	Matrix      []MatrixExam `json:"matrix_exam" bson:"matrix_exam"`
	QuestionIDs []string     `json:"question_ids" bson:"question_ids"`

	StartTime       string `json:"start_time" bson:"start_time"`
	EndTime         string `json:"end_time" bson:"end_time"`
	DurationMinutes int    `json:"duration_minutes" bson:"duration_minutes"`

	TestScore float32 `json:"test_score" bson:"test_score"`
}

type UpdateTestOfClassRequest struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	ClassID primitive.ObjectID `json:"class_id" bson:"class_id"`
	UserID  string             `json:"user_id,omitempty" bson:"user_id,omitempty"`

	TestName string   `json:"test_name" bson:"test_name"`
	Descript string   `json:"descript,omitempty" bson:"descript,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`

	QuestionIDs []string     `json:"question_ids" bson:"question_ids"`
	Matrix      []MatrixExam `json:"matrix_exam" bson:"matrix_exam"`

	IsTest bool `json:"is_test" bson:"is_test"`

	StartTime       string `json:"start_time" bson:"start_time"`
	EndTime         string `json:"end_time" bson:"end_time"`
	DurationMinutes int    `json:"duration_minutes" bson:"duration_minutes"`

	TestScore float32 `json:"test_score" bson:"test_score"`
}

// ================= CORE LOGIC =================

func validateQuestionIDs(ids []string) error {
	for _, questionID := range ids {
		questionUUID, err := primitive.ObjectIDFromHex(questionID)
		if err != nil || questionUUID == primitive.NilObjectID {
			return errors.New("invalid Question ID: " + questionID)
		}
	}
	return nil
}

func NewTest(test TestTemplete) (TestTemplete, error) {
	if test.ID == primitive.NilObjectID {
		return TestTemplete{}, errors.New("invalid Class ID")
	}

	if err := validateQuestionIDs(test.QuestionIDs); err != nil {
		return TestTemplete{}, err
	}

	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()

	if test.ID == primitive.NilObjectID {
		test.ID = primitive.NewObjectID()
	}

	return test, nil
}

func NewTestOfClass(test TestOfClass) (TestOfClass, error) {
	if test.ID == primitive.NilObjectID {
		test.ID = primitive.NewObjectID()
	}

	if test.ClassID == primitive.NilObjectID {
		return TestOfClass{}, errors.New("invalid Class ID")
	}

	if err := validateQuestionIDs(test.QuestionIDs); err != nil {
		return TestOfClass{}, err
	}

	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()

	return test, nil
}

// ================= PARSE: TestTemplete =================

func ParseNewTest(req CreateTestRequest, ids ...primitive.ObjectID) (TestTemplete, error) {
	testEntity := TestTemplete{
		BaseTest: BaseTest{
			TestName:        req.TestName,
			Descript:        req.Descript,
			QuestionIDs:     req.QuestionIDs,
			IsTest:          req.IsTest,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationMinutes: req.DurationMinutes,
			Tags:            req.Tags,
			Matrix:          req.Matrix,
			TestScore:       req.TestScore,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	if len(ids) > 0 {
		testEntity.ID = ids[0]
	} else {
		testEntity.ID = primitive.NewObjectID()
	}

	return NewTest(testEntity)
}

func ParseUpdateTest(req UpdateTestRequest) (TestTemplete, error) {
	testEntity := TestTemplete{
		ID: req.ID,
		BaseTest: BaseTest{
			TestName:        req.TestName,
			Descript:        req.Descript,
			QuestionIDs:     req.QuestionIDs,
			IsTest:          req.IsTest,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationMinutes: req.DurationMinutes,
			Tags:            req.Tags,
			Matrix:          req.Matrix,
			TestScore:       req.TestScore,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	return NewTest(testEntity)
}

// ================= PARSE: TestOfClass =================

func ParseNewTestOfClass(req CreateTestOfClassRequest, ids ...primitive.ObjectID) (TestOfClass, error) {
	testEntity := TestOfClass{
		ClassID: req.ClassID,
		BaseTest: BaseTest{
			TestName:        req.TestName,
			Descript:        req.Descript,
			QuestionIDs:     req.QuestionIDs,
			IsTest:          req.IsTest,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationMinutes: req.DurationMinutes,
			Tags:            req.Tags,
			Matrix:          req.Matrix,
			TestScore:       req.TestScore,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	if len(ids) > 0 {
		testEntity.ID = ids[0]
	} else {
		testEntity.ID = primitive.NewObjectID()
	}

	return NewTestOfClass(testEntity)
}

func ParseUpdateTestOfClass(req UpdateTestOfClassRequest) (TestOfClass, error) {
	testEntity := TestOfClass{
		ID:      req.ID,
		ClassID: req.ClassID,
		BaseTest: BaseTest{
			TestName:        req.TestName,
			Descript:        req.Descript,
			QuestionIDs:     req.QuestionIDs,
			IsTest:          req.IsTest,
			StartTime:       req.StartTime,
			EndTime:         req.EndTime,
			DurationMinutes: req.DurationMinutes,
			Tags:            req.Tags,
			Matrix:          req.Matrix,
			TestScore:       req.TestScore,
			UpdatedAt:       time.Now(),
		},
	}

	return NewTestOfClass(testEntity)
}
