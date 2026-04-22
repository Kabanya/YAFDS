// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"net/http"

	repository "github.com/Kabanya/YAFDS/courier/internal/repository/postgres"
	"github.com/Kabanya/YAFDS/courier/internal/service"
	"github.com/Kabanya/YAFDS/courier/internal/usecase"

	"github.com/Kabanya/YAFDS/pkg/bootstrap"
	"github.com/Kabanya/YAFDS/pkg/utils"
)

func Run() {
	logger := bootstrap.InitLogger("courier_log_info.txt", "courier")
	bootstrap.LoadEnv(logger)

	courierDB := bootstrap.OpenPostgres(logger, "COURIER_DB", "courier_db", "database")
	defer courierDB.Close()

	ordersDB := bootstrap.OpenPostgres(logger, "ORDER_DB", "order_db", "orders database")
	defer ordersDB.Close()

	userRepository := repository.NewUser(courierDB)
	logger.Println("Initialized user repository")

	logger.Println("TODO: make orders repository normalno")

	redisClient := bootstrap.OpenRedis(logger)
	defer redisClient.Close()

	sessionTTL := bootstrap.SessionTTL(logger)
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

	port := bootstrap.Port("COURIER_PORT", "8090")
	logger.Println("Endpoints registered:")
	logger.Printf("  POST http://localhost:%s/register - Register user with password", port)
	logger.Printf("  POST http://localhost:%s/login - Login user with password", port)
	logger.Printf("  GET  http://localhost:%s/orders - List orders", port)

	bootstrap.ListenAndServe(logger, port, http.DefaultServeMux)

	logger.Println("Process of courier is finished")
	utils.CloseLogger()
}
