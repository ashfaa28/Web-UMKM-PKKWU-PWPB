package routes

import (
	"UMKM/app/controller"
	"database/sql"
	"net/http"
)

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	server.HandleFunc("/menu/order", controller.NewAddOrder(db))
	server.HandleFunc("/verfiedUserLogin", controller.LoginChecker(db))
	server.HandleFunc("/register", controller.Register(db))
	server.HandleFunc("/", controller.NewIndexHtml(db))
}
