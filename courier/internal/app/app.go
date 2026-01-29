// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"courier/internal/repository"
	"courier/internal/service"
	"courier/internal/usecase"

	pkgHandlers "github.com/Kabanya/YAFDS/pkg/app"
	"github.com/Kabanya/YAFDS/pkg/common/utils"
	pkg_repository "github.com/Kabanya/YAFDS/pkg/repository"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func Run() {
	utils.InitFileLogger("courier_log_info.txt")
	logger, err := utils.Logger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	logger.Println("courier service started")

	// Load environment variables from .env
	err = utils.LoadEnv(utils.PathToEnv)
	if err != nil {
		logger.Printf("Failed to load .env file: %v", err)
		panic(err)
	}

	// Connection to db
	dbName := os.Getenv("COURIER_DB")
	if dbName == "" {
		dbName = "courier_db"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), dbName)
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

	ordersDBName := os.Getenv("ORDER_DB")
	if ordersDBName == "" {
		ordersDBName = "order_db"
	}
	ordersConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ordersDBName)
	ordersDB, err := sql.Open("postgres", ordersConnStr)
	if err != nil {
		logger.Printf("Failed to open orders database: %v", err)
		panic(err)
	}
	defer ordersDB.Close()

	if err := ordersDB.Ping(); err != nil {
		logger.Printf("Failed to ping orders database: %v", err)
		panic(err)
	}
	logger.Println("Successfully connected to orders database")

	userRepository := repository.NewUser(db)
	logger.Println("Initialized user repository")

	ordersRepository := pkg_repository.NewPostgresRepository(ordersDB, db, db)
	logger.Println("Initialized orders repository")

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

	sessionTTL := utils.TimeTtl30Minutes
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
		logger.Printf("SESSION_TTL must be positive, using default %v", utils.TimeTtl30Minutes)
		sessionTTL = utils.TimeTtl30Minutes
	}

	userService := service.NewUserService(userRepository, redisClient, sessionTTL)
	logger.Println("Initialized user service")

	userUseCase := usecase.NewUserUseCase(userService)
	logger.Println("Initialized user usecase")

	handler := NewHandler(userUseCase)
	logger.Println("Initialized handler")

	// registry endpoints
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/orders", pkgHandlers.NewListHandler(ordersRepository))

	port := os.Getenv("COURIER_PORT")
	if port == "" {
		port = "8090"
	}
	addr := ":" + port
	logger.Println("Endpoints registered:")
	logger.Printf("  POST http://localhost:%s/register - Register user with password", port)
	logger.Printf("  POST http://localhost:%s/login - Login user with password", port)
	logger.Printf("  GET  http://localhost:%s/orders - List orders", port)
	logger.Printf("Starting HTTP server on %s", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Printf("Server error: %v", err)
	}

	logger.Println("Process of courier is finished")
	utils.CloseLogger()
}
