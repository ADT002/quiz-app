package initialize

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"quiz-app/internal/auth"
	"quiz-app/internal/domain/usecase"
	"quiz-app/internal/infrastructure/persistence/aws"
	persistence "quiz-app/internal/infrastructure/persistence/mongodb"
	redisdb "quiz-app/internal/infrastructure/persistence/redis"
	routes "quiz-app/internal/infrastructure/router"

	"github.com/rs/cors"
)

var (
	dbSevice = "mongodb"
	dbName   = "dbapp"
)

func InitApp() {
	InitDB()
	InitRouter()
}

func InitDB() {
	switch dbSevice {
	case "mongodb":
		persistence.ConnectMongoDB(os.Getenv(persistence.MongoConnectionString))
	default:
	}
}

func InitRouter() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router := routes.NewRouter()

	redisRepo, err := redisdb.GetRedis()
	if err != nil {
		fmt.Println("NOT RUN REDIS")
	}
	redisUseCase := usecase.NewRedisUseCase(redisRepo)

	authUseCase := auth.NewAuthUseCase(os.Getenv("AUTH_SERVICE_URL"))
	authHandler := auth.NewAuthHandler(authUseCase)

	// Repositories
	questionRepo := persistence.NewQuestionMongoRepository()
	topicRepo := persistence.NewTopicMongoRepository()
	levelRepo := persistence.NewLevelMongoRepository()

	classRepo := persistence.NewClassMongoRepository()
	testTemplateRepo := persistence.NewTestTemplateMongoRepository()
	testOfClassRepo := persistence.NewTestOfClassMongoRepository()
	fileRepo := persistence.NewFileMongoRepository()
	submissionRepo := persistence.NewSubmissionMongoRepository()

	// Use cases
	questionUseCase := usecase.NewQuestionUseCase(questionRepo)
	topicUseCase := usecase.NewTopicUseCase(topicRepo)
	levelUseCase := usecase.NewLevelUseCase(levelRepo)

	classUseCase := usecase.NewClassUseCase(classRepo, testOfClassRepo)
	testTemplateUseCase := usecase.NewTestTemplateUseCase(testTemplateRepo)
	testOfClassUseCase := usecase.NewTestOfClassUseCase(testOfClassRepo)
	fileUseCase := usecase.NewFileUseCase(fileRepo)

	pdfService := usecase.NewPDFService()
	submissionUseCase := usecase.NewSubmissionUseCase(submissionRepo, pdfService)
	awsS3UseCase := aws.NewFileAWSRepository("quiz-app-image-storage")

	// Routers
	routes.NewRouterTestTemplate(*testTemplateUseCase, *authHandler).GetTestTemplateRouter(router)
	routes.NewRouterTestOfClass(*testOfClassUseCase, *authHandler).GetTestOfClassRouter(router)
	routes.NewRouterQuestion(*questionUseCase, *topicUseCase, *levelUseCase, *authHandler).GetQuestionRouter(router)
	routes.NewRouterClass(*classUseCase, *redisUseCase, *authHandler).GetClassRouter(router)
	routes.NewRoutesFile(fileUseCase, awsS3UseCase, authHandler).GetRoutesFile(router)
	routes.NewRouterSubmission(submissionUseCase, testTemplateUseCase, questionUseCase, redisUseCase, authHandler).GetSubmissionRouter(router)

	handler := c.Handler(router)
	router.HandleFunc("/", homeHandler).Methods("GET")

	port := "8081"
	log.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {}

// suppress unused variable warning
var _ = dbName
