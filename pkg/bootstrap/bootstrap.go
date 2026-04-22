package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Kabanya/YAFDS/pkg/middleware"
	"github.com/Kabanya/YAFDS/pkg/utils"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func InitLogger(logFile string, serviceName string) *log.Logger {
	if err := utils.InitFileLogger(logFile); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logger, err := utils.Logger()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logger.Printf("%s service started", serviceName)
	return logger
}

func LoadEnv(logger *log.Logger) {
	if err := utils.LoadEnv(utils.PathToEnv); err != nil {
		logger.Printf("Failed to load .env file: %v", err)
		panic(err)
	}
}

func OpenPostgres(logger *log.Logger, envName string, defaultDBName string, label string) *sql.DB {
	dbName := os.Getenv(envName)
	if dbName == "" {
		dbName = defaultDBName
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Printf("Failed to open %s: %v", label, err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.Printf("Failed to ping %s: %v", label, err)
		panic(err)
	}

	logger.Printf("Successfully connected to %s", label)
	return db
}

func OpenRedis(logger *log.Logger) *redis.Client {
	redisDB := 0
	if redisDBStr := os.Getenv("REDIS_DB"); redisDBStr != "" {
		parsed, err := strconv.Atoi(redisDBStr)
		if err == nil {
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

	logger.Println("Successfully connected to Redis")
	return redisClient
}

func SessionTTL(logger *log.Logger) time.Duration {
	sessionTTL := utils.TimeTtl30Minutes
	ttlStr := os.Getenv("SESSION_TTL")
	if ttlStr == "" {
		return sessionTTL
	}

	parsed, err := time.ParseDuration(ttlStr)
	if err != nil {
		seconds, parseIntErr := strconv.ParseInt(ttlStr, 10, 64)
		if parseIntErr != nil {
			logger.Printf("Invalid SESSION_TTL '%s', using default %v: %v", ttlStr, sessionTTL, parseIntErr)
			return sessionTTL
		}
		parsed = time.Duration(seconds) * time.Second
	}

	if parsed <= 0 {
		logger.Printf("SESSION_TTL must be positive, using default %v", utils.TimeTtl30Minutes)
		return utils.TimeTtl30Minutes
	}

	return parsed
}

func Port(portEnv string, defaultPort string) string {
	port := os.Getenv(portEnv)
	if port == "" {
		return defaultPort
	}

	return port
}

func ListenAndServe(logger *log.Logger, port string, handler http.Handler) {
	addr := ":" + port
	logger.Printf("Starting HTTP server on %s", addr)

	handlerWithCORS := middleware.CORSMiddleware(handler)
	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		logger.Printf("Server error: %v", err)
	}
}
