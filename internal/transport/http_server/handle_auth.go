package http_server

import (
	"internal/adapters/logger"
	"internal/app"
	"internal/domain"
	"net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var md domain.UserLoginMetadata
	md.IP = GetClientIP(r)

	data, hs := app.LoginUser(r.Body, &md)
	if hs.Err != nil {
		logger.Error("handleLogin: " + hs.Message)
		http.Error(w, hs.Message, hs.HTTPStatus)
		return
	}
	if data != nil {
		w.Write(*data)
	}
}

// func handleLogout(w http.ResponseWriter, r *http.Request) {
// 	http.SetCookie(w, &http.Cookie{
// 		Name:   "session",
// 		Value:  "",
// 		Path:   "/",
// 		MaxAge: -1,
// 	})
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
