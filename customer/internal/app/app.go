// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"net/http"

	repository "github.com/Kabanya/YAFDS/customer/internal/repository/postgres"
	"github.com/Kabanya/YAFDS/customer/internal/service"
	"github.com/Kabanya/YAFDS/customer/internal/usecase"
	"github.com/Kabanya/YAFDS/pkg/bootstrap"
	pkgRepo "github.com/Kabanya/YAFDS/pkg/order/repository/postgres"
	pkgService "github.com/Kabanya/YAFDS/pkg/order/service"
	"github.com/Kabanya/YAFDS/pkg/utils"
	"github.com/Kabanya/YAFDS/pkg/wallet"
)

func Run() {
	logger := bootstrap.InitLogger("customer_log_info.txt", "Customer")
	bootstrap.LoadEnv(logger)

	customerDB := bootstrap.OpenPostgres(logger, "CUSTOMER_DB", "customer_db", "database")
	defer customerDB.Close()

	ordersDB := bootstrap.OpenPostgres(logger, "ORDER_DB", "order_db", "orders database")
	defer ordersDB.Close()

	courierDB := bootstrap.OpenPostgres(logger, "COURIER_DB", "courier_db", "courier database")
	defer courierDB.Close()

	userRepository := repository.NewUserRepo(customerDB)
	logger.Println("Initialized user repository")

	ordersRepository := pkgRepo.NewPostgresRepository(ordersDB, customerDB, courierDB)
	logger.Println("Initialized orders repository")

	redisClient := bootstrap.OpenRedis(logger)
	defer redisClient.Close()

	sessionTTL := bootstrap.SessionTTL(logger)
	userService := service.NewUserService(userRepository, redisClient, sessionTTL)
	logger.Println("Initialized user service")

	userUseCase := usecase.NewUserUseCase(userService)
	logger.Println("Initialized user usecase")

	customerOrderRepository := repository.NewPostgresUserRepo(customerDB, courierDB, ordersDB)
	walletClient := wallet.NewWalletClient()

	orderService := pkgService.NewOrderService(ordersRepository)
	customerOrderService := service.NewOrderService(customerOrderRepository)
	orderUseCase := usecase.NewOrderUseCase(orderService, customerOrderService, walletClient)
	logger.Println("Initialized order usecase")

	catalogRepository := repository.NewCatalogRepo(courierDB, customerDB)
	catalogService := service.NewCatalogService(catalogRepository)
	catalogUseCase := usecase.NewCatalogUseCase(catalogService)
	logger.Println("Initialized catalog usecase")

	handler := NewHandler(userUseCase, customerDB)
	logger.Println("Initialized handler")

	orderHandler := NewOrderHandler(orderUseCase)
	// registry endpoints
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("GET /orders", orderHandler.OrdersHandler())
	http.HandleFunc("POST /orders", orderHandler.OrdersHandler())
	http.HandleFunc("GET /orders/{order_id}", orderHandler.GetOrder)
	http.HandleFunc("GET /orders/{order_id}/status", orderHandler.GetOrderStatus)
	http.HandleFunc("PUT /orders/{order_id}/status", orderHandler.UpdateOrderStatus)
	http.HandleFunc("POST /orders/{order_id}/accept", orderHandler.AcceptOrder)
	http.HandleFunc("GET /orders/{order_id}/total", orderHandler.CalculateOrderTotal)
	http.HandleFunc("POST /orders/{order_id}/pay", orderHandler.PayOrder)
	http.HandleFunc("POST /orders/{order_id}/items", orderHandler.AddItemIntoOrder)
	http.HandleFunc("/couriers", NewCourierCatalogHandler(catalogUseCase))
	http.HandleFunc("/restaurants", NewRestaurantCatalogHandler(catalogUseCase))
	http.HandleFunc("/menu", NewRestaurantMenuCatalogHandler(catalogUseCase))

	port := bootstrap.Port("CUSTOMER_PORT", "8091")
	logger.Println("Endpoints registered:")
	logger.Printf("  POST     http://localhost:%s/register - Register user with password", port)
	logger.Printf("  POST     http://localhost:%s/login - Login user with password", port)
	logger.Printf("  POST/GET http://localhost:%s/orders - Create/List orders", port)
	logger.Printf("  POST     http://localhost:%s/orders/{order_id}/pay - Pay for order", port)
	logger.Printf("  POST     http://localhost:%s/orders/{order_id}/accept - Accept order", port)
	logger.Printf("  POST     http://localhost:%s/orders/{order_id}/items - Add order item", port)
	logger.Printf("  GET      http://localhost:%s/couriers - List active couriers", port)
	logger.Printf("  GET      http://localhost:%s/restaurants - List active restaurants", port)
	logger.Printf("  GET      http://localhost:%s/menu?restaurant_id=<uuid> - Show restaurant menu items", port)

	bootstrap.ListenAndServe(logger, port, http.DefaultServeMux)

	logger.Println("Process of customer is finished")
	utils.CloseLogger()
}
