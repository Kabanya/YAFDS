// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"net/http"

	"github.com/Kabanya/YAFDS/pkg/bootstrap"
	pkgRepo "github.com/Kabanya/YAFDS/pkg/order/repository/postgres"
	pkgService "github.com/Kabanya/YAFDS/pkg/order/service"
	"github.com/Kabanya/YAFDS/pkg/utils"
	repository "github.com/Kabanya/YAFDS/restaurant/internal/repository/postgres"
	"github.com/Kabanya/YAFDS/restaurant/internal/service"
	"github.com/Kabanya/YAFDS/restaurant/internal/usecase"
)

func Run() {
	logger := bootstrap.InitLogger("restaurant_log_info.txt", "restaurant")
	bootstrap.LoadEnv(logger)

	restaurantDB := bootstrap.OpenPostgres(logger, "RESTAURANT_DB", "restaurant_db", "database")
	defer restaurantDB.Close()

	ordersDB := bootstrap.OpenPostgres(logger, "ORDER_DB", "order_db", "orders database")
	defer ordersDB.Close()

	userRepository := repository.NewUser(restaurantDB)
	logger.Println("Initialized user repository")

	restaurantMenuItemsRepo := repository.NewRestaurantMenuItemsRepo(restaurantDB)
	logger.Println("Initialized restaurant menu items repository")

	ordersRepository := repository.NewOrdersRepo(ordersDB, restaurantDB)
	logger.Println("Initialized orders repository")

	redisClient := bootstrap.OpenRedis(logger)
	defer redisClient.Close()

	sessionTTL := bootstrap.SessionTTL(logger)
	userService := service.NewUserService(userRepository, redisClient, sessionTTL)
	logger.Println("Initialized user service")

	restaurantMenuItemsService := service.NewRestaurantMenuItemsService(restaurantMenuItemsRepo)
	logger.Println("Initialized restaurant menu items service")

	pkgOrderRepo := pkgRepo.NewPostgresRepository(ordersDB, nil, nil)
	pkgOrderService := pkgService.NewOrderService(pkgOrderRepo)
	ordersService := service.NewOrderService(ordersRepository, pkgOrderRepo, pkgOrderService)
	logger.Println("Initialized orders service")

	userUseCase := usecase.NewUserUseCase(userService)
	logger.Println("Initialized user usecase")

	restaurantMenuItemsUseCase := usecase.NewRestaurantMenuItemsUseCase(restaurantMenuItemsService)
	logger.Println("Initialized restaurant menu items usecase")

	ordersUseCase := usecase.NewOrderUseCase(ordersService)
	logger.Println("Initialized orders usecase")

	handler := NewHandler(userUseCase, restaurantMenuItemsUseCase, ordersUseCase)
	logger.Println("Initialized handler")

	// registry endpoints
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/orders", handler.ListOrders)
	http.HandleFunc("/menu/show", handler.ShowMenuItems)
	http.HandleFunc("/menu/upload", handler.UploadMenuItem)

	port := bootstrap.Port("RESTAURANT_PORT", "8092")
	logger.Println("Endpoints registered:")
	logger.Printf("  POST http://localhost:%s/register - Register user with password", port)
	logger.Printf("  POST http://localhost:%s/login - Login user with password", port)
	logger.Printf("  GET  http://localhost:%s/orders?restaurant_id=<uuid> - List restaurant orders", port)
	logger.Printf("  GET  http://localhost:%s/menu/show?restaurant_id=<uuid> - Show menu items", port)
	logger.Printf("  POST http://localhost:%s/menu/upload - Upload menu item", port)

	bootstrap.ListenAndServe(logger, port, http.DefaultServeMux)

	logger.Println("Process of restaurant is finished")
	utils.CloseLogger()
}
