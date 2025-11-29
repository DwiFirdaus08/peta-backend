package main

import (
	"log"

	"backend-peta/routes"
)

func main() {
	
	app := routes.SetupApp()
	
	log.Println("Server berjalan di port 3000")
	app.Listen(":3000")
}