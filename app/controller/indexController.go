package controller

import (
	"database/sql"
	"html/template"
	"net/http"
	"path/filepath"
)

func NewIndexHtml(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fp := filepath.Join("app", "views", "index.html")
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
