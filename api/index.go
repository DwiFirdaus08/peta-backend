package handler

import (
	"net/http"

	"backend-peta/routes"

	"github.com/gofiber/adaptor/v2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Panggil SetupApp dari package routes
	app := routes.SetupApp()
	
	adaptor.FiberApp(app).ServeHTTP(w, r)
}