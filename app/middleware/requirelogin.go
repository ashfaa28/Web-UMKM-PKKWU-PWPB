package middleware

import (
	"net/http"
)

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah pengguna sudah login
		session, err := r.Cookie("session_id")
		if err != nil || !isValidSession(session.Value) {
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
			return
		}

		// Jika sudah login, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func isValidSession(sessionID string) bool {
	// Implementasikan logika untuk memeriksa apakah sesi valid
	// Misalnya, periksa dalam database atau cache
	return true
}
