// use make logs to get debug_customer.txt into your rep

package logger

import (
	"fmt"
	"log"
	"os"
)

func HelloLog() {
	log.Println("Start of customer logger!")
}

func Init(filename string) error {
	filename = "debug_customer"

	logFile := fmt.Sprintf("%s.txt", filename)

	file, err := os.Create(logFile)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	HelloLog()
	return nil
}

func PrintLog(mess string) {
	log.Println(mess)
}
