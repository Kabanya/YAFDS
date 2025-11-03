package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hi from restaurant process!")
	var n int = 360

	for i := 0; i < n; i++ {
		fmt.Println(i, "- Restaurant: ", time.Now())
		time.Sleep(time.Second)
	}
	fmt.Println("Restaurant process is finished!")
}
