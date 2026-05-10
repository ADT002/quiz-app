package persistence

import (
	"context"
	entity "quiz-app/internal/domain/entity"
	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubmissionMongoRepository struct {
	Collection *mongo.Collection
}

func NewSubmissionMongoRepository() repository.SubmissionRepository {
	client := GetMongoClient()
	return &SubmissionMongoRepository{
		Collection: client.Database("dbapp").Collection("classes"),
	}
}

func (ar *SubmissionMongoRepository) ExportData(
	ctx context.Context,
	classID primitive.ObjectID,
	testID primitive.ObjectID,
) (*entity.ClassExport, error) {

	pipeline := mongo.Pipeline{

		// 1. Match class
		{{Key: "$match", Value: bson.M{"_id": classID}}},

		// 2. Unwind test
		{{Key: "$unwind", Value: "$test"}},

		// 3. Match test
		{{Key: "$match", Value: bson.M{"test.id": testID}}},

		// 4. Collect submitted user_ids
		{{Key: "$addFields", Value: bson.M{
			"submitted_user_ids": bson.M{
				"$map": bson.M{
					"input": "$test.user_submit",
					"as":    "u",
					"in":    "$$u.user_id",
				},
			},
		}}},

		// 4.1 Map students_accept -> user_ids
		{{Key: "$addFields", Value: bson.M{
			"students_accept_user_ids": bson.M{
				"$map": bson.M{
					"input": "$students_accept",
					"as":    "s",
					"in":    "$$s.user_id",
				},
			},
		}}},

		// 5. Compute unsubmitted user_ids
		{{Key: "$addFields", Value: bson.M{
			"unsubmitted_user_ids": bson.M{
				"$setDifference": bson.A{
					"$students_accept_user_ids",
					"$submitted_user_ids",
				},
			},
		}}},

		// ================= SUBMITTED =================

		{{Key: "$unwind", Value: bson.M{
			"path":                       "$test.user_submit",
			"preserveNullAndEmptyArrays": true,
		}}},

		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "test.user_submit.user_id",
			"foreignField": "user_id",
			"as":           "user",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$user",
			"preserveNullAndEmptyArrays": true,
		}}},

		{{Key: "$lookup", Value: bson.M{
			"from": "submissions",
			"let": bson.M{
				"class_id": "$_id",
				"test_id":  "$test.id",
				"user_id":  "$test.user_submit.user_id",
			},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$class_id", "$$class_id"}},
								bson.M{"$eq": bson.A{"$test_id", "$$test_id"}},
								bson.M{"$eq": bson.A{"$user_id", "$$user_id"}},
							},
						},
					},
				},
			},
			"as": "submission",
		}}},
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$submission",
			"preserveNullAndEmptyArrays": true,
		}}},

		// ================= GROUP =================

		{{Key: "$group", Value: bson.M{
			"_id": nil,

			"class": bson.M{"$first": bson.M{
				"id":         "$_id",
				"class_name": "$class_name",
			}},

			"test": bson.M{"$first": bson.M{
				"_id":              "$test.id",
				"test_name":        "$test.test_name",
				"descript":         "$test.descript",
				"start_time":       "$test.start_time",
				"end_time":         "$test.end_time",
				"duration_minutes": "$test.duration_minutes",
				"user_id":          "$test.user_id",
				"author_mail":      "$test.author_mail",
				"test_score":       "$test.test_score",
			}},

			"submitted_users": bson.M{
				"$push": bson.M{
					"$cond": bson.A{
						bson.M{"$ne": bson.A{"$user", nil}},
						bson.M{
							"email":      "$user.email",
							"user_id":    "$user.user_id",
							"first_name": "$user.first_name",
							"last_name":  "$user.last_name",
							"score":      bson.M{"$ifNull": bson.A{"$submission.score", 0}},
							"start_time": "$submission.start_time",
							"end_time":   "$submission.end_time",
							"submitted":  true,
						},
						"$$REMOVE",
					},
				},
			},

			"unsubmitted_user_ids": bson.M{"$first": "$unsubmitted_user_ids"},
		}}},

		// 12. Lookup unsubmitted users
		{{Key: "$lookup", Value: bson.M{
			"from": "users",
			"let":  bson.M{"ids": "$unsubmitted_user_ids"},
			"pipeline": bson.A{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$in": bson.A{"$user_id", "$$ids"},
						},
					},
				},
			},
			"as": "unsubmitted_users",
		}}},

		// 13. Normalize unsubmitted users
		{{Key: "$addFields", Value: bson.M{
			"unsubmitted_users": bson.M{
				"$map": bson.M{
					"input": "$unsubmitted_users",
					"as":    "u",
					"in": bson.M{
						"email":      "$$u.email",
						"user_id":    "$$u.user_id",
						"first_name": "$$u.first_name",
						"last_name":  "$$u.last_name",
						"score":      0,
						"submitted":  false,
					},
				},
			},
		}}},

		// 14. Merge users
		{{Key: "$addFields", Value: bson.M{
			"users": bson.M{
				"$concatArrays": bson.A{
					"$submitted_users",
					"$unsubmitted_users",
				},
			},
		}}},

		// 15. Final shape
		{{Key: "$project", Value: bson.M{
			"_id":   0,
			"class": 1,
			"test":  1,
			"users": 1,
		}}},
	}

	cursor, err := ar.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var raw []struct {
		Class struct {
			ID        primitive.ObjectID `bson:"id"`
			ClassName string             `bson:"class_name"`
		} `bson:"class"`

		Test entity.ExportTest `bson:"test"`

		Users []entity.ExportUserSubmission `bson:"users"`
	}

	if err := cursor.All(ctx, &raw); err != nil {
		return nil, err
	}

	if len(raw) == 0 {
		return nil, nil
	}

	result := &entity.ClassExport{
		ClassID:   raw[0].Class.ID,
		ClassName: raw[0].Class.ClassName,
		Test:      raw[0].Test,
	}
	result.Test.Users = raw[0].Users

	return result, nil
}
