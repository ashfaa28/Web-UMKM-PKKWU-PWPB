package controller

import (
	"UMKM/app/store"
	"net/http"
)

func LogOut(w http.ResponseWriter, r *http.Request) {
	// Dapatkan session
	session, _ := store.Store.Get(r, "session-name")

	// Hapus nilai dari session
	session.Values["user_id"] = nil
	session.Values["username"] = nil

	// Simpan perubahan session
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Gagal menghapus session", http.StatusInternalServerError)
		return
	}

	// Hapus cookie session
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Redirect ke halaman login setelah logout
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
