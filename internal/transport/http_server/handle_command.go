package http_server

import (
	"internal/adapters/logger"
	"internal/app"
	"net/http"
)

// HTTP request processing
func HandleCommand(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Type", "application/json")

	data, hs := app.ExecuteCommand(r.Context(), r.Body)
	if hs.Err != nil {
		logger.Error("HandleCommand: " + hs.Message)
		http.Error(w, hs.Message, hs.HTTPStatus)
		return
	}
	if data != nil {
		w.Write(*data)
	}
}
