// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"context"
	"customer/internal/repository"
	"customer/internal/service"
	"customer/internal/usecase"
	"customer/pkg"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

	// Connection to db
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

	redisDB := 0
	if redisDBStr := os.Getenv("REDIS_DB"); redisDBStr != "" {
		if parsed, err := strconv.Atoi(redisDBStr); err == nil {
			redisDB = parsed
		} else {
			logger.Printf("Invalid REDIS_DB value '%s', defaulting to 0: %v", redisDBStr, err)
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		logger.Printf("Failed to connect to Redis: %v", err)
		panic(err)
	}
	defer redisClient.Close()
	logger.Println("Successfully connected to Redis")

	sessionTTL := pkg.TimeTtl30Minutes
	if ttlStr := os.Getenv("SESSION_TTL"); ttlStr != "" {
		var parsed time.Duration
		if d, err := time.ParseDuration(ttlStr); err == nil {
			parsed = d
		} else if sec, err := strconv.ParseInt(ttlStr, 10, 64); err == nil {
			parsed = time.Duration(sec) * time.Second
		} else {
			logger.Printf("Invalid SESSION_TTL '%s', using default %v: %v", ttlStr, sessionTTL, err)
			parsed = 0
		}
		if parsed > 0 {
			sessionTTL = parsed
		}
	}
	if sessionTTL <= 0 {
		logger.Printf("SESSION_TTL must be positive, using default %v", pkg.TimeTtl30Minutes)
		sessionTTL = pkg.TimeTtl30Minutes
	}

	userService := service.NewUserService(userRepository, redisClient, sessionTTL)
	logger.Println("Initialized user service")

	userUseCase := usecase.NewUserUseCase(userService)
	logger.Println("Initialized user usecase")

	handler := NewHandler(userUseCase)
	logger.Println("Initialized handler")

	// registry endpoints
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)

	logger.Println("Endpoints registered:")
	logger.Println("  POST http://localhost:8081/register - Register user with password")
	logger.Println("  POST http://localhost:8081/login - Login user with password")
	logger.Println("Starting HTTP server on :8081")

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		logger.Printf("Server error: %v", err)
	}

	logger.Println("Process of customer is finished")
	pkg.CloseLogger()
}
