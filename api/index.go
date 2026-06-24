package handler

import (
	"net/http"

	"sipecut/config"
	"sipecut/routes"
	"sipecut/storage"
)

var app http.Handler

func init() {
	// Initialize Database, Storage, and Router once
	config.ConnectDatabase()
	storage.InitStorage()
	app = routes.SetupRouter()
}

// Handler handles Vercel Serverless requests
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
