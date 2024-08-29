package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

type UserData struct {
	UserID   int
	Username string
}

func LoginChecker(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			renderTemplate(w, nil)
			return
		}

		if r.Method == http.MethodPost {
			email := r.FormValue("email")
			password := r.FormValue("password")

			var userID int
			var username, hashedPassword string

			queryLogin, err := db.Prepare("SELECT username, user_id, password FROM akun WHERE email = ?")
			if err != nil {
				http.Error(w, "Server error, unable to prepare SQL statement", http.StatusInternalServerError)
				return
			}
			defer queryLogin.Close()

			err = queryLogin.QueryRow(email).Scan(&username, &userID, &hashedPassword)
			if err != nil {
				if err == sql.ErrNoRows {
					renderTemplate(w, map[string]string{"ErrorMessage": "Pastikan Email dan Password benar"})
					return
				}
				http.Error(w, "Server error, unable to log in", http.StatusInternalServerError)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
			if err != nil {
				renderTemplate(w, map[string]string{"ErrorMessage": "Pastikan Email dan Password benar"})
				return
			}

			session, _ := store.Get(r, "session-name")
			session.Values["user_id"] = userID
			session.Values["username"] = username
			session.Save(r, w)

			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "some_session_id",
				Path:  "/",
			})

			http.Redirect(w, r, "/order", http.StatusSeeOther)
		}
	}
}

func renderTemplate(w http.ResponseWriter, data interface{}) {
	fp := filepath.Join("app", "views", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}
