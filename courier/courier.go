package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hi from courier process!")

	var n int = 10
	fmt.Println()
	for i := 0; i < n; i++ {
		fmt.Println(i, "- Courier", time.Now())
		time.Sleep(time.Second)
	}
	fmt.Println("Courier process is finished!")
}
