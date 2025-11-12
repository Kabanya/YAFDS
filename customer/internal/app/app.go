package app

import (
	logger "customer/pkg"
	"net/http"
)

func Run() {
	logger.Init("debug_customer")

	service := NewService()
	handler := NewHandler(service)

	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)

	logger.PrintLog("Link: http://localhost:8081/register")
	logger.PrintLog("Link: http://localhost:8081/login")
	logger.PrintLog("HandleFunc of customer is started")

	http.ListenAndServe(":8081", nil)
}
