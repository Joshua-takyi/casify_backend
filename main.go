package main

import (
	"fmt"

	"github.com/joshua/casify/router"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("server starting soon")
	router := router.Router()
	router.Run(":8000")
	fmt.Println("server started")
}