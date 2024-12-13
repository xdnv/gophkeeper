// the main server module provides server (metric storage and update) function
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"internal/app"
	"internal/domain"
	"internal/ports/storage"
	"internal/transport/http_server"

	"internal/adapters/cryptor"
	"internal/adapters/logger"
)

// version descriptor
var version = domain.GetVersion()

func main() {
	//sync internal/logger upon exit
	defer logger.Sync()

	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// a WaitGroup for the goroutines to tell us they've stopped
	wg := sync.WaitGroup{}

	//Warning! do not run outside function, it will break tests due to flag.Parse()
	app.Sc = app.InitServerConfig()

	app.Stor = storage.NewUniStorage(&app.Sc)
	defer app.Stor.Close()

	//post-init unistorage actions
	err := app.Stor.Bootstrap()
	if err != nil {
		logger.Fatalf("srv: post-init bootstrap failed, error: %s", err)
	}

	// run `server` in its own goroutine
	wg.Add(1)
	go server(ctx, &wg)

	// listen for ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	logger.Info("srv: received ^C - shutting down")

	// tell the goroutines to stop
	logger.Info("srv: telling goroutines to stop")
	cancel()

	// and wait for them to reply back
	wg.Wait()

	logger.Info("srv: shutdown")
}

func server(ctx context.Context, wg *sync.WaitGroup) {
	//execute to exit wait group
	defer wg.Done()

	// version descriptor
	logger.Infof("Build version: %s", naIfEmpty(version.Version))
	logger.Infof("Build date: %s", naIfEmpty(version.Date))
	logger.Infof("Build commit: %s", naIfEmpty(version.Commit))

	logger.Infof("srv: transport mode %s", app.Sc.TransportMode)
	logger.Infof("srv: using endpoint %s", app.Sc.Endpoint)
	logger.Infof("srv: storage mode = %v", app.Sc.StorageMode)
	logger.Infof("srv: compress replies = %v %v", app.Sc.CompressReplies, app.Sc.CompressibleContentTypes)
	logger.Infof("srv: encryption=%v", cryptor.CanDecrypt())

	switch app.Sc.StorageMode {
	case domain.Database:
		//remove password from log output
		re := regexp.MustCompile(`(password)=(?P<password>\S*)`)
		s := re.ReplaceAllLiteralString(app.Sc.DatabaseDSN, "password=***")
		logger.Infof("srv: DSN %s", s)

	case domain.File:
	}

	//init and run server
	var srv *http.Server = http_server.ServeHTTP()

	<-ctx.Done()
	logger.Info("srv: shutdown requested")

	// shut down gracefully with timeout of 5 seconds max
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful server shutdown
	srv.Shutdown(shutdownCtx) // ignore server error "Err shutting down server : context canceled"

	logger.Info("srv: server stopped")
}

func naIfEmpty(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
