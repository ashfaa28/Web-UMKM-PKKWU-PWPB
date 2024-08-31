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

	server.HandleFunc("/", controller.NewIndexHtml(db))
	server.HandleFunc("/verfiedUserLogin", controller.LoginChecker(db))
	server.HandleFunc("/register", controller.Register(db))
	server.HandleFunc("/logOut", controller.LogOut)
	http.HandleFunc("/checkOut", controller.CheckoutHandler(db))
	server.Handle("/akun", middleware.RequireLogin(http.HandlerFunc(controller.AccInfoHandler)))
	server.Handle("/order", middleware.RequireLogin(http.HandlerFunc(controller.NewAddOrder(db))))
}
