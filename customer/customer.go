package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hi from customer process!")
	var n int = 11

	for i := 0; i < n; i++ {
		fmt.Println(i, "- Customer", time.Now())
		time.Sleep(time.Second)
	}
	fmt.Println("Courier process is finished!")
}
