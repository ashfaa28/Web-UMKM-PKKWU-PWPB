package routes

import (
	"UMKM/app/controller"
	"UMKM/app/middleware"
	"database/sql"
	"net/http"
)

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	fs := http.FileServer(http.Dir("app/static"))
	server.Handle("/static/", http.StripPrefix("/static/", fs))

	// Rute tanpa middleware
	server.HandleFunc("/verfiedUserLogin", controller.LoginChecker(db))
	server.HandleFunc("/register", controller.Register(db))
	server.HandleFunc("/", controller.NewIndexHtml(db))

	// Rute dengan middleware
	server.Handle("/order", middleware.RequireLogin(http.HandlerFunc(controller.NewAddOrder(db))))
}
