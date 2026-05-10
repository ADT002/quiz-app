package routes

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	entity "quiz-app/internal/domain/entity"
// 	"quiz-app/internal/domain/usecase"
// 	"quiz-app/internal/pkg"
// 	utils "quiz-app/internal/util"
// 	"sort"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type RoutesTest struct {
// 	auth usecase.AuthHandler

// 	redisUseCase usecase.RedisUseCase

// 	testUseCase     usecase.TestUseCase
// 	questionUseCase usecase.QuestionUseCase
// 	classUseCase    usecase.ClassUseCase
// 	submissionUseCase   usecase.SubmissionUseCase
// }

// // NewRouterTest creates a new RoutesTest instance
// func NewRouterTest(testUseCase usecase.TestUseCase, classUseCase usecase.ClassUseCase, questionUseCase usecase.QuestionUseCase, submissionUseCase usecase.SubmissionUseCase, redisUseCase usecase.RedisUseCase, auth usecase.AuthHandler) RoutesTest {
// 	return RoutesTest{
// 		testUseCase:     testUseCase,
// 		classUseCase:    classUseCase,
// 		questionUseCase: questionUseCase,
// 		redisUseCase:    redisUseCase,
// 		submissionUseCase:   submissionUseCase,
// 		auth:            auth,
// 	}
// }

// // InitializeRoutesTests initializes all test-related routes
// func (rt RoutesTest) GetTestRouter(r *Router) {
// 	// Routes for managing tests
// 	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getAllTestFromAuthor))).Methods("GET")
// 	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.createTest))).Methods("POST")
// 	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.updateTest))).Methods("PATCH")
// 	r.Handle("/tests", rt.auth.AuthMiddleware(http.HandlerFunc(rt.deleteTest))).Methods("DELETE")

// 	// Routes for class-specific operations
// 	r.Handle("/tests/class", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getAllTestOfClassByEmail))).Methods("POST")

// 	// Routes for managing questions within tests
// 	r.Handle("/tests/questions", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getQuestionOfTest))).Methods("POST")

// 	// Routes for marking test completion and sending test results
// 	r.Handle("/test/done", rt.auth.AuthMiddleware(http.HandlerFunc(rt.getDoneTest))).Methods("POST")
// }

// func (r *RoutesTest) createTest(w http.ResponseWriter, req *http.Request) {
// 	emailID := req.Context().Value("user_id").(string)
// 	email := req.Context().Value("email").(string)

// 	var test entity.Test

// 	// Generate update fields from the test struct
// 	if err := json.NewDecoder(req.Body).Decode(&test); err != nil {
// 		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	test.UserID = emailID
// 	test.EmailName = email

// 	insertedID, err := r.testUseCase.CreateTest(context.TODO(), &test)

// 	if err != nil {
// 		pkg.SendError(w, "Invalid create test", http.StatusBadRequest)
// 	}

// 	test.ID = insertedID
// 	pkg.SendResponse(w, http.StatusCreated, test)
// }

// func (r *RoutesTest) updateTest(w http.ResponseWriter, req *http.Request) {
// 	emailID, ok := req.Context().Value("user_id").(string)

// 	if !ok {
// 		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
// 		return
// 	}

// 	var testUpdate entity.Test
// 	// Generate update fields from the test struct
// 	if err := json.NewDecoder(req.Body).Decode(&testUpdate); err != nil {
// 		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	fmt.Println(testUpdate)
// 	testUpdate.UserID = emailID
// 	r.testUseCase.UpdateTest(context.TODO(), &testUpdate)
// 	pkg.SendResponse(w, http.StatusOK, testUpdate)
// }

// func (r *RoutesTest) deleteTest(w http.ResponseWriter, req *http.Request) {
// 	emailID := req.Context().Value("user_id").(string)

// 	var testDelete struct {
// 		ID primitive.ObjectID `json:"_id"`
// 	}

// 	if err := json.NewDecoder(req.Body).Decode(&testDelete); err != nil {
// 		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	err := r.testUseCase.DeleteTest(context.TODO(), testDelete.ID, emailID)
// 	if err != nil {
// 		pkg.SendError(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}
// 	pkg.SendResponse(w, http.StatusOK, "")
// }

// // GetAllTestOfClassByEmail retrieves all tests for a specific class by email and class ID.
// func (r *RoutesTest) getAllTestFromAuthor(w http.ResponseWriter, req *http.Request) {
// 	emailID, ok := req.Context().Value("email").(string)
// 	if !ok {
// 		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch tests based on class ID and email
// 	tests, err := r.testUseCase.GetTestsByAuthorEmail(req.Context(), emailID)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to get tests", http.StatusInternalServerError)
// 		return
// 	}

// 	pkg.SendResponse(w, http.StatusOK, tests)
// }

// // GetAllTestOfClassByEmail retrieves all tests for a specific class by email and class ID.
// func (r *RoutesTest) getAllTestOfClassByEmail(w http.ResponseWriter, req *http.Request) {
// 	// Extract email from context
// 	email, ok := req.Context().Value("email").(string)
// 	if !ok {
// 		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
// 		return
// 	}

// 	type classIDRequest struct {
// 		ClassIDs primitive.ObjectID `json:"_id"`
// 	}

// 	// Decode class IDs from request body
// 	var classIDData classIDRequest
// 	if err := json.NewDecoder(req.Body).Decode(&classIDData); err != nil {
// 		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
// 		return
// 	}

// 	// Validate that ClassIDs is not empty
// 	if len(classIDData.ClassIDs) == 0 {
// 		pkg.SendError(w, "Class ID list cannot be empty", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch tests based on class IDs and email
// 	tests, err := r.classUseCase.GetAllTestOfClass(req.Context(), email, classIDData.ClassIDs)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to get tests", http.StatusInternalServerError)
// 		return
// 	}

// 	pkg.SendResponse(w, http.StatusOK, tests)
// }

// func (r *RoutesTest) getQuestionOfTest(w http.ResponseWriter, req *http.Request) {
// 	email, ok := req.Context().Value("email").(string)
// 	emailID, ok := req.Context().Value("user_id").(string)

// 	if !ok {
// 		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
// 		return
// 	}

// 	var test struct {
// 		ClassID   primitive.ObjectID `json:"class_id"`
// 		TestID    primitive.ObjectID `json:"test_id"`
// 		EmailName string             `json:"author_mail"`
// 		IsTest    bool               `json:"is_test"`
// 	}

// 	if err := json.NewDecoder(req.Body).Decode(&test); err != nil {
// 		pkg.SendError(w, "Invalid request data", http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Println(test)

// 	cacheKeyTestInfo := fmt.Sprintf("test_info_%s", test.TestID.Hex())
// 	cacheKeyQuestions := fmt.Sprintf("questions_%s", test.TestID.Hex())

// 	// Check and load cached test info and questions
// 	testInfo, questions, err := r.loadCachedTestData(req.Context(), cacheKeyTestInfo, cacheKeyQuestions, email)
// 	if err == nil && testInfo != nil && questions != nil {
// 		pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": pkg.ShuffleQuestionsAndSubmissions(questions)})
// 		return
// 	}

// 	// Fetch question IDs and test info if not cached
// 	questionIDs, testInfo, err := r.classUseCase.GetQuestionOfTest(req.Context(), test.ClassID, test.TestID, email)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to retrieve test data", http.StatusInternalServerError)
// 		return
// 	}
// 	// Cache test info

// 	var minutesDuration int
// 	if val, ok := testInfo["duration_minutes"].(int32); ok {
// 		minutesDuration = int(val)
// 	} else {
// 		fmt.Println("Fail to get testInfo['duration_minutes']")
// 	}
// 	// Convert minutesDuration directly to time.Duration
// 	duration := time.Duration(minutesDuration+1) * time.Minute
// 	cacheData(r.redisUseCase, req.Context(), cacheKeyTestInfo, testInfo, duration)
// 	// Delete allowed users for security and validate timing
// 	// if !isTestAccessible(testInfo) {
// 	// 	pkg.SendError(w, "TEST IS NOT ALLOWED", http.StatusForbidden)
// 	// 	return
// 	// }

// 	// Fetch questions and cache results
// 	questions, err = r.questionUseCase.GetAllQuestionsOfTest(req.Context(), questionIDs)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to retrieve questions", http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println("ALL QUESTION: ", questions)

// 	// Generate final options and cache them
// 	finalOptionMap := r.generateFinalOptionsMap(questions, test.TestID.Hex())
// 	questionsJSON, err := json.Marshal(finalOptionMap)
// 	if err != nil {
// 		pkg.SendError(w, "Failed to process questions", http.StatusInternalServerError)
// 		return
// 	}
// 	r.redisUseCase.HSet(req.Context(), fmt.Sprintf("questions_id_%s", test.TestID.Hex()), duration, map[string]interface{}{"questions": questionsJSON})

// 	if testInfo["is_test"] == true {
// 		shuffledQuestions := pkg.ShuffleQuestionsAndSubmissions(questions)
// 		cacheData(r.redisUseCase, req.Context(), cacheKeyQuestions, shuffledQuestions, duration)
// 		if submission, err := r.submissionUseCase.GetSubmission(req.Context(), primitive.M{"test_id": test.TestID, "user_id": emailID}); err == nil && len(submission.ListQuestionSubmission) != 0 {
// 			response := primitive.M{"test_info": testInfo, "submission": submission, "questions": questions}
// 			pkg.SendResponse(w, http.StatusOK, response)
// 			return
// 		} else {
// 			r.submissionUseCase.CreateNewSubmission(req.Context(), &entity.TestSubmission{
// 				TestId:  test.TestID,
// 				UserID: emailID,
// 				Email:   email,
// 			})
// 			pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": shuffledQuestions})
// 			return
// 		}
// 	} else {
// 		if submission, err := r.submissionUseCase.GetSubmission(req.Context(), primitive.M{"test_id": test.TestID, "user_id": emailID}); err == nil && len(submission.ListQuestionSubmission) != 0 {
// 			response := primitive.M{"test_info": testInfo, "submission": submission, "questions": questions}
// 			pkg.SendResponse(w, http.StatusOK, response)
// 			return
// 		}
// 		pkg.SendResponse(w, http.StatusOK, primitive.M{"test_info": testInfo, "questions": questions})
// 		return
// 	}
// }

// func (r *RoutesTest) loadCachedTestData(
// 	ctx context.Context,
// 	testInfoKey, questionsKey, email string,
// ) (map[string]interface{}, []primitive.M, error) {

// 	cachedTestInfo, errInfo := r.redisUseCase.Get(ctx, testInfoKey)
// 	cachedQuestions, errQuestions := r.redisUseCase.Get(ctx, questionsKey)

// 	if errInfo != nil || errQuestions != nil {
// 		return nil, nil, fmt.Errorf("cache miss")
// 	}

// 	var testInfo map[string]interface{}
// 	var questions []primitive.M

// 	if err := json.Unmarshal([]byte(cachedTestInfo), &testInfo); err != nil {
// 		return nil, nil, err
// 	}
// 	if err := json.Unmarshal([]byte(cachedQuestions), &questions); err != nil {
// 		return nil, nil, err
// 	}

// 	allowedUsersRaw, exists := testInfo["allowed_users"]
// 	if !exists || allowedUsersRaw == nil {
// 		return nil, nil, fmt.Errorf("allowed_users missing")
// 	}

// 	allowedUsers, ok := allowedUsersRaw.([]interface{})
// 	if !ok {
// 		return nil, nil, fmt.Errorf("allowed_users invalid type")
// 	}

// 	for _, user := range allowedUsers {
// 		if userStr, ok := user.(string); ok && userStr == email {
// 			// ❗ copy map trước khi delete để tránh dirty cache
// 			safeTestInfo := make(map[string]interface{}, len(testInfo))
// 			for k, v := range testInfo {
// 				safeTestInfo[k] = v
// 			}
// 			delete(safeTestInfo, "allowed_users")
// 			return safeTestInfo, questions, nil
// 		}
// 	}

// 	return nil, nil, fmt.Errorf("user not allowed")
// }

// // generateFinalOptionsMap prepares the options map for questions
// func (r *RoutesTest) generateFinalOptionsMap(questions []primitive.M, testID string) map[string]map[string]map[string]interface{} {
// 	finalOptionMap := make(map[string]map[string]map[string]interface{})
// 	for _, question := range questions {
// 		questionType, ok := question["type"].(string)
// 		if !ok {
// 			continue
// 		}
// 		questionID := r.getQuestionID(question)
// 		if questionID == "" {
// 			continue
// 		}

// 		optionMap := map[string]interface{}{questionID: []string{}}
// 		switch questionType {
// 		case "order_question":
// 			optionMap[questionID] = r.handleOrderQuestion(question)
// 		case "single_choice_question":
// 			optionMap[questionID] = r.handleSingleChoiceQuestion(question)
// 		case "multiple_choice_question":
// 			optionMap[questionID] = r.handleMultipleChoiceQuestion(question)
// 		case "fill_in_the_blank":
// 			optionMap[questionID] = r.handleFillInTheBlank(question)
// 		case "match_choice_question":
// 			optionMap[questionID] = r.handleMatchChoiceQuestion(question)
// 		}

// 		if finalOptionMap[testID] == nil {
// 			finalOptionMap[testID] = make(map[string]map[string]interface{})
// 		}
// 		finalOptionMap[testID][questionID] = map[string]interface{}{
// 			"optionMap": optionMap,
// 			"type":      questionType,
// 			"score":     question["score"].(float64),
// 		}
// 	}
// 	return finalOptionMap
// }

// // Helper to retrieve question ID
// func (r *RoutesTest) getQuestionID(question primitive.M) string {
// 	if id, ok := question["_id"].(primitive.ObjectID); ok {
// 		return id.Hex()
// 	} else if idStr, ok := question["_id"].(string); ok {
// 		return idStr
// 	}
// 	return ""
// }

// func cacheData(redisUseCase usecase.RedisUseCase, ctx context.Context, key string, data interface{}, time time.Duration) {
// 	dataJSON, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("Error marshaling data for caching:", err)
// 		return
// 	}
// 	if err := redisUseCase.Set(ctx, key, dataJSON, time); err != nil {
// 		fmt.Println("Error caching data:", err)
// 	}
// }

// func isTestAccessible(testInfo map[string]interface{}) bool {
// 	startTime, errStart := utils.StringToTime(testInfo["start_time"].(string))
// 	endTime, errEnd := utils.StringToTime(testInfo["end_time"].(string))

// 	if errStart != nil || errEnd != nil {
// 		fmt.Println("Error processing test timing")
// 		return false
// 	}
// 	currentTime := time.Now()
// 	return startTime.Before(currentTime) && endTime.After(currentTime)
// }

// // // sendTest processes the test submission
// // func (r *RoutesTest) sendTest(w http.ResponseWriter, req *http.Request) {

// // 	// TODO: Process the test submission based on email and emailID

// // 	pkg.SendResponse(w, http.StatusOK, "Test sent successfully")
// // }

// // getDoneTest handles retrieving a user's done test
// func (r *RoutesTest) getDoneTest(w http.ResponseWriter, req *http.Request) {
// 	emailID, ok := req.Context().Value("user_id").(string)
// 	if !ok {
// 		pkg.SendError(w, "Invalid email ID", http.StatusBadRequest)
// 		return
// 	}

// 	// TODO: Implement fetching done test logic

// 	pkg.SendResponse(w, http.StatusOK, fmt.Sprintf("Done test retrieved for email ID: %s", emailID))
// }

// // Handle question
// // handleOrderQuestion processes "order_question" type questions
// func (r *RoutesTest) handleOrderQuestion(question primitive.M) []string {
// 	optionsValue, exists := question["order_items"]
// 	if !exists {
// 		fmt.Println("Error: Order items key does not exist in the question map")
// 		return nil
// 	}

// 	rawOptions, ok := optionsValue.(primitive.A)
// 	if !ok {
// 		fmt.Println("Error: Order items is not of type primitive.A")
// 		return nil
// 	}

// 	options := make([]map[string]interface{}, 0, len(rawOptions))
// 	for _, opt := range rawOptions {
// 		if option, ok := opt.(primitive.M); ok {
// 			options = append(options, option)
// 		} else {
// 			fmt.Println("Error: Option is not of type map[string]interface{}")
// 		}
// 	}

// 	// Sort options based on the "order" field
// 	sort.Slice(options, func(i, j int) bool {
// 		return options[i]["order"].(int32) < options[j]["order"].(int32)
// 	})

// 	var orderedIDs []string
// 	for _, option := range options {
// 		if id, ok := option["id"].(primitive.ObjectID); ok {
// 			orderedIDs = append(orderedIDs, id.Hex())
// 		} else {
// 			fmt.Println("Error: Option ID is not of type ObjectID")
// 		}
// 	}
// 	return orderedIDs
// }

// // handleSingleChoiceQuestion processes "single_choice_question" type questions
// func (r *RoutesTest) handleSingleChoiceQuestion(question primitive.M) []string {
// 	optionsValue, ok := question["options"].(primitive.A)
// 	if !ok {
// 		fmt.Println("Error: Options is not of type primitive.A")
// 		return nil
// 	}

// 	for _, option := range optionsValue {
// 		if optionMap, ok := option.(primitive.M); ok {

// 			if isCorrect, exists := optionMap["is_correct"].(bool); exists && isCorrect {
// 				if id, idOk := optionMap["id"].(primitive.ObjectID); idOk {
// 					return []string{id.Hex()}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// // handleMultipleChoiceQuestion processes "multiple_choice_question" type questions
// func (r *RoutesTest) handleMultipleChoiceQuestion(question primitive.M) []string {
// 	optionsValue, ok := question["options"].(primitive.A)
// 	if !ok {
// 		fmt.Println("Error: Options is not of type primitive.A")
// 		return nil
// 	}

// 	var correctIDs []string
// 	for _, option := range optionsValue {
// 		if optionMap, ok := option.(primitive.M); ok {
// 			if isCorrect, exists := optionMap["is_correct"].(bool); exists && isCorrect {
// 				if id, idOk := optionMap["id"].(primitive.ObjectID); idOk {
// 					correctIDs = append(correctIDs, id.Hex())
// 				}
// 			}
// 		}
// 	}
// 	return correctIDs
// }

// // handleFillInTheBlank processes "fill_in_the_blank" type questions
// func (r *RoutesTest) handleFillInTheBlank(question primitive.M) []map[string]string {
// 	fillInBlanks, ok := question["fill_in_the_blanks"].(primitive.A)
// 	if !ok {
// 		fmt.Println("Error: fill_in_the_blank is not of type primitive.A")
// 		return nil
// 	}

// 	var fillInData []map[string]string
// 	for _, item := range fillInBlanks {
// 		itemMap, ok := item.(primitive.M)
// 		if !ok {
// 			fmt.Println("Error: item in fill_in_the_blank is not of type map[string]interface{}")
// 			continue
// 		}

// 		id, idOk := itemMap["id"].(primitive.ObjectID)
// 		submission, submissionOk := itemMap["correct_submission"].(string)
// 		if idOk && submissionOk {
// 			fillInData = append(fillInData, map[string]string{
// 				"id":             id.Hex(),
// 				"correct_submission": submission,
// 			})
// 		}
// 	}
// 	return fillInData
// }

// // handleMatchChoiceQuestion processes "match_choice_question"
// func (r *RoutesTest) handleMatchChoiceQuestion(question primitive.M) map[string][]string {
// 	options, ok := question["match_options"].(primitive.A)
// 	if !ok {
// 		fmt.Println("Error: match_options is not primitive.A")
// 		return map[string][]string{}
// 	}

// 	matchMap := make(map[string][]string)

// 	for _, opt := range options {
// 		optionMap, ok := opt.(primitive.M)
// 		if !ok {
// 			continue
// 		}

// 		optionID, ok1 := optionMap["id"].(primitive.ObjectID)
// 		matchID, ok2 := optionMap["match_id"].(string)

// 		if !ok1 || !ok2 {
// 			continue
// 		}

// 		// map: match_item_id -> []option_id
// 		matchMap[matchID] = append(matchMap[matchID], optionID.Hex())
// 	}

// 	return matchMap
// }
