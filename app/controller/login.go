package controller

import (
	"UMKM/app/middleware"
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func LoginChecker(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			var hashedPassword string
			stmt, err := db.Prepare("SELECT password FROM akun WHERE email = ?")
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

			// Set session cookie
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "some_session_id", // Generate and store actual session ID
				Path:  "/",
			})

			// Check if there is a saved URL to redirect to
			requestedURL, err := r.Cookie(middleware.RequestedURLKey)
			if err == nil && requestedURL != nil {
				http.Redirect(w, r, requestedURL.Value, http.StatusSeeOther)
				return
			}

			// Redirect to a default page if no requested URL is saved
			http.Redirect(w, r, "/menu/order", http.StatusSeeOther)
		} else {
			// Display the login page
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
		}
	}
}
