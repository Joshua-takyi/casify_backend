package main

import (
	"github.com/joshua/casify/router"
)

func main() {
	router := router.Router()
	router.Run(":8000")
}
