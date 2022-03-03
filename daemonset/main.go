package main

import (
	"fmt"
	"time"
)

func main() {

	for {
		fmt.Printf("I am here, I exist")
		time.Sleep(time.Second * 15)
	}
}
