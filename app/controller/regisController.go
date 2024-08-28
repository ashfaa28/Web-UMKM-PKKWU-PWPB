package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB) http.HandlerFunc {
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

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Server error, unable to create your account", http.StatusInternalServerError)
				return
			}

			stmt, err := db.Prepare("INSERT INTO akun (username, email, password) VALUES (?, ?, ?)")
			if err != nil {
				http.Error(w, "Server error, unable to prepare SQL statement", http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(username, email, hashedPassword)
			if err != nil {
				http.Error(w, "Server error, unable to create your account", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
		}
	}
}
