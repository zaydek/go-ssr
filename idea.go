package main

import (
	"fmt"
	"time"
)

var epoch = time.Now()

func fn() func() {
	return func() {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Since(epoch))
	}
}

func main() {
	fmt.Println(time.Since(epoch))
	defer fn()()
}
