package main

import (
	"customer/internal/app"
	"customer/pkg"
)

func main() {
	app.Run()
	logger, _ := pkg.Logger()
	logger.Println("Process of customer is finished")
	pkg.CloseLogger()
}
