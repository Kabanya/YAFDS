// 0. на go структру конфига, который будет парсить переменные среды.
// Три разных комплекта перемных среды. .env1
// docker - единица деплоя

package app

import (
	"customer/internal/repository"
	"customer/internal/service"
	logger "customer/pkg"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func Run() {
	logger.Init("debug_customer")

	var db *sql.DB

	userRepo := repository.NewUser(db)

	userService1 := service.NewUser1(userRepo)
	userService2 := service.NewUser2(userRepo)

	userService1.Save(uuid.Max, "", "", "")
	userService2.Load("")

	// handler := NewHandler(service)

	// http.HandleFunc("/register", handler.Register)
	// http.HandleFunc("/login", handler.Login)

	// logger.PrintLog("Link: http://localhost:8081/register")
	// logger.PrintLog("Link: http://localhost:8081/login")
	// logger.PrintLog("HandleFunc of customer is started")

	http.ListenAndServe(":8081", nil)
}
