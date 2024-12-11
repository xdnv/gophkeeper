package http_server

import (
	"internal/adapters/logger"
	"internal/app"
	"net/http"
)

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	data, hs := app.RegisterNewUser(r.Body)
	if hs.Err != nil {
		logger.Error("handleRegistration: " + hs.Message)
		http.Error(w, hs.Message, hs.HTTPStatus)
		return
	}
	w.WriteHeader(http.StatusOK)
	if data != nil {
		w.Write(*data)
	}
}
