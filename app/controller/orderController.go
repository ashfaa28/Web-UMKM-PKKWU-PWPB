package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
)

// Pastikan Anda sudah mendeklarasikan store di bagian atas file atau di tempat yang sesuai

func NewAddOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil sesi
		session, _ := store.Get(r, "session-name")

		// Periksa apakah sesi valid
		userID, userIDok := session.Values["user_id"].(int)
		username, usernameok := session.Values["username"].(string)
		if !userIDok || !usernameok {
			// Jika sesi tidak valid, redirect ke halaman login
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
			return
		}

		// Buat data untuk template
		data := UserData{
			UserID:   userID,
			Username: username,
		}

		// Parsing dan eksekusi template
		fp := filepath.Join("app", "views", "order.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
