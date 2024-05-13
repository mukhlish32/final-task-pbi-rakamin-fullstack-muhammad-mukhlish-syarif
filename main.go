package main

import (
	"rakamin/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
