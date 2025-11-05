package app

import (
	"fmt"
	"net/http"
)

func Run() {
	service := NewService()
	handler := NewHandler(service)
	fmt.Printf("Link: %s", "http://localhost:8080/courier-hello-world")
	http.HandleFunc("/courier-hello-world", handler.Hello)
	http.ListenAndServe(":8080", nil)
}
