package http_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	if app.Sc.CompressReplies {
		mux.Use(middleware.Compress(5, app.Sc.CompressibleContentTypes...))
	}

	mux.Post("/login", handleLogin)
	mux.Post("/register", handleRegistration)
	mux.Post("/command", HandleJWTAuth(http.HandlerFunc(HandleCommand)))
	mux.Post("/ping", HandleJWTAuth(http.HandlerFunc(HandlePingDBServer)))
	//mux.Post("/logout", handleLogout)

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
