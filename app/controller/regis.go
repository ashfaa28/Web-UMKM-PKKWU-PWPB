package controller

import (
	"UMKM/app/database"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fp := filepath.Join("app", "views", "register.html")
			tmpl, err := template.ParseFiles(fp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error parsing template: " + err.Error()))
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error executing template: " + err.Error()))
				return
			}
			return
		}

		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")

			// Hash password sebelum disimpan ke database
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Server error, unable to create your account", http.StatusInternalServerError)
				return
			}

			// Simpan username, email, dan hashed password ke database
			_, err = database.DB.Exec("INSERT INTO akun (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
			if err != nil {
				http.Error(w, "Server error, unable to create your account", http.StatusInternalServerError)
				return
			}

			// Redirect ke halaman login setelah pendaftaran berhasil
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
		}
	}
}
