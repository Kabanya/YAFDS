package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hi from courier process!")
	var n int
	fmt.Printf("How many time call seconds: ")
	fmt.Scan(&n)

	for i := 0; i < n; i++ {
		fmt.Println(i, "- Courier", time.Now())
		time.Sleep(time.Second)
	}
	fmt.Println("Courier process is finished!")
}
