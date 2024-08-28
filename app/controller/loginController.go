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

			var userID int
			var username string
			var hashedPassword string

			queryLogin, err := db.Prepare("SELECT username, user_id, password FROM akun WHERE email = ?")
			if err != nil {
				http.Error(w, "Server error, unable to prepare SQL statement", http.StatusInternalServerError)
				return
			}
			defer queryLogin.Close()

			err = queryLogin.QueryRow(email).Scan(&username, &userID, &hashedPassword)
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

			session, _ := store.Get(r, "session-name")
			session.Values["user_id"] = userID
			session.Values["username"] = username
			session.Save(r, w)

			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "some_session_id", // membuat dan store session dengan ID tertentu (dengan gorrila session)
				Path:  "/",
			})

			http.Redirect(w, r, "/order", http.StatusSeeOther)
		}
	}
}
