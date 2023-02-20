package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now().Format("20060102150405")
	fmt.Println(currentTime)
}
