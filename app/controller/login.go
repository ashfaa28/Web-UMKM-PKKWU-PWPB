package controller

import (
	"UMKM/app/database"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func LoginChecker(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var hashedPassword string

		// Gunakan prepared statement untuk mencegah SQL Injection
		stmt, err := database.DB.Prepare("SELECT password FROM users WHERE email = ?")
		if err != nil {
			http.Error(w, "Server error, unable to prepare SQL statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow(email).Scan(&hashedPassword)
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

		// Jika berhasil login, redirect ke dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
