package controller

import (
	"UMKM/app/store"
	"database/sql"
	"html/template"
	"net/http"
)

type UserInfo struct {
	Username string
	Email    string
}

func AccInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve username from session
	session, _ := store.Store.Get(r, "session-name")
	var username string
	if val, ok := session.Values["username"]; ok {
		username, _ = val.(string)
	}

	// Fetch user information from the database
	db, err := sql.Open("mysql", "root:2007hadi@tcp(localhost:3306)/angkringan")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var email string
	err = db.QueryRow("SELECT email FROM akun WHERE username = ?", username).Scan(&email)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Create UserInfo struct with fetched data
	userInfo := UserInfo{
		Username: username,
		Email:    email,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("app/views/accInfo.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
