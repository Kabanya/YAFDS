// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"customer/internal/repository"
	"customer/internal/service"
	"customer/internal/usecase"
	"customer/pkg"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func Run() {
	pkg.InitFileLogger("customer_log_info.txt")
	logger, err := pkg.Logger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	logger.Println("Customer service started")

	// Load environment variables from .env
	err = pkg.LoadEnv(pkg.PathToEnv)
	if err != nil {
		logger.Printf("Failed to load .env file: %v", err)
		panic(err)
	}

	// Подключение к базе данных
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("CUSTOMER_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Printf("Failed to open database: %v", err)
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Printf("Failed to ping database: %v", err)
		panic(err)
	}
	logger.Println("Successfully connected to database")

	userRepository := repository.NewUser(db)
	logger.Println("Initialized user repository")

	userService := service.NewUserService(userRepository)
	logger.Println("Initialized user service")

	userUseCase := usecase.NewUserUseCase(userService)
	logger.Println("Initialized user usecase")

	handler := NewHandler(userUseCase)
	logger.Println("Initialized handler")

	// registry endpoints
	http.HandleFunc("/save", handler.SaveUser)
	http.HandleFunc("/load", handler.LoadUser)

	logger.Println("Endpoints registered:")
	logger.Println("  POST http://localhost:8081/save - Save user")
	logger.Println("  POST http://localhost:8081/load - Load user")
	logger.Println("Starting HTTP server on :8081")

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		logger.Printf("Server error: %v", err)
	}

	logger.Println("Customer service stopped")
}
