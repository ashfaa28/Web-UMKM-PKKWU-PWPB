package middleware

import (
	"net/http"
)

const RequestedURLKey = "menu/order"

func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah cookie sesi ada dan valid
		session, err := r.Cookie("session_id")
		if err != nil || !isValidSession(session.Value) {
			// Simpan URL yang diminta dalam cookie atau sesi
			http.SetCookie(w, &http.Cookie{
				Name:  RequestedURLKey,
				Value: r.URL.Path,
				Path:  "/menu/order",
			})

			// Redirect ke halaman login
			http.Redirect(w, r, "/verfiedUserLogin", http.StatusSeeOther)
			return
		}

		// Jika sesi valid, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func isValidSession(sessionID string) bool {
	// Implementasikan logika untuk memeriksa apakah sesi valid
	// Misalnya, periksa dalam database atau cache
	return true
}
