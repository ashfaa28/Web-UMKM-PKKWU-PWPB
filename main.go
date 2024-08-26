package main

import (
	"UMKM/app/database"
	"UMKM/app/routes"
	"fmt"
	"net/http"
)

func main() {
	db := database.InitDatabase()

	server := http.NewServeMux()

	routes.MapRoutes(server, db)

	fs := http.FileServer(http.Dir("app/static"))
	server.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server berjalan di http://localhost:8080") // Pindahkan sebelum ListenAndServe

	err := http.ListenAndServe(":8080", server)
	if err != nil {
		fmt.Println("Error saat menjalankan server:", err)
	}
}
