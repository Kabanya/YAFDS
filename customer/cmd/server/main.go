package main

import (
	"customer/internal/app"
	logger "customer/pkg"
)

func main() {
	app.Run()
	logger.PrintLog("Process of customer is finished")
}
