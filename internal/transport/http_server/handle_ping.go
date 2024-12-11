package http_server

import (
	"internal/adapters/logger"
	"internal/app"
	"net/http"
)

// HTTP request processing
func HandlePingDBServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)

	//username := r.Context().Value(domain.CtxUsername)
	//logger.Infof("This is Ping. Hello %v!", username) //DEBUG

	hs := app.PingDBServer(r.Context())
	if hs.Err != nil {
		logger.Error(hs.Message)
		http.Error(w, hs.Message, hs.HTTPStatus)
		return
	}

	w.Write([]byte(hs.Message))
}
