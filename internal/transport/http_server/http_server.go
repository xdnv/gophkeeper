package http_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"internal/adapters/cryptor"
	"internal/adapters/logger"

	"internal/app"
)

// http server part
func ServeHTTP() *http.Server {
	mux := chi.NewRouter()
	mux.Use(logger.LoggerMiddleware)
	mux.Use(HandleGZIPRequests)
	mux.Use(cryptor.HandleEncryptedRequests)
	// if app.Sc.CompressReplies {
	// 	mux.Use(middleware.Compress(5, app.Sc.CompressibleContentTypes...))
	// }

	mux.Get("/", HandleIndex)
	mux.Get("/ping", HandlePingDBServer)
	mux.Post("/value/", HandleRequestMetricV2)
	mux.Post("/update/", HandleUpdateMetricV2)
	mux.Post("/updates/", HandleUpdateMetrics)

	// create a server
	srv := &http.Server{Addr: app.Sc.Endpoint, Handler: mux}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logger.Errorf("Listen: %s", err.Error())
		}
	}()

	return srv
}
