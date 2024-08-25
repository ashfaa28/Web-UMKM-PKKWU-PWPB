package controller

import (
	"UMKM/app/database"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func LoginChecker(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fp := filepath.Join("app", "views", "login.html")
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

		// Proses login jika method POST
		email := r.FormValue("email")
		password := r.FormValue("password")

		var hashedPassword string
		err := database.DB.QueryRow("SELECT password FROM akun WHERE email=?", email).Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error, unable to log in", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Email or password is incorrect", http.StatusUnauthorized)
			return
		}

		// Jika berhasil login, redirect atau tampilkan pesan sukses
		http.Redirect(w, r, "/orderList", http.StatusSeeOther)
	}
}
