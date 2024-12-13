// the main agent module provides agent (metric sender) function
package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// "internal/adapters/cryptor"
	"internal/adapters/logger"
	"internal/app"
	"internal/domain"
	"internal/ports/console"
	// "internal/domain"
	// "github.com/google/uuid"
	// "github.com/shirou/gopsutil/cpu"
	// "github.com/shirou/gopsutil/v3/mem"
)

// var sendJobs chan uuid.UUID

// version descriptor
var version = domain.GetVersion()

// store agent IP address to send in header
//var agentIP string

// // returns first non-loopback local IP address
// func getLocalIP() (string, error) {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "", err
// 	}

// 	var localIP string
// 	for _, addr := range addrs {
// 		// Check if address is IPv4 and not loopback
// 		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
// 			localIP = ipnet.IP.String()
// 			break
// 		}
// 	}

// 	if localIP == "" {
// 		return "", errors.New("no valid local IP address found")
// 	}

// 	return localIP, nil
// }

func main() {
	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// a WaitGroup for the goroutines to tell us they've stopped
	wg := sync.WaitGroup{}

	//Warning! do not run outside function, it will break tests due to flag.Parse()
	app.Cc = app.InitClientConfig()

	wg.Add(1)
	go client(ctx, &wg)

	// listen for ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	logger.Info("client: received ^C - shutting down")

	// tell the goroutines to stop
	logger.Info("client: telling goroutines to stop")
	cancel()

	// and wait for them to reply back
	wg.Wait()
	logger.Info("client: shutdown")
}

func client(ctx context.Context, wg *sync.WaitGroup) {
	//execute to exit wait group
	defer wg.Done()

	// version descriptor
	logger.Infof("Build version: %s", naIfEmpty(version.Version))
	logger.Infof("Build date: %s", naIfEmpty(version.Date))
	logger.Infof("Build commit: %s", naIfEmpty(version.Commit))

	app := console.NewApp(ctx)
	app.Init()
	if err := app.Run(); err != nil {
		panic(err)
	}
	//wg.Done() //while we have no other goroutines

	//commit := fmt.Sprintf("v.%s git:%s (%s)", version.Version, version.Commit, version.Date)
	//logger.Infof("VERSION: %s", naIfEmpty(commit))

	// logger.Infof("client: transport mode %s", ac.TransportMode)
	// logger.Infof("client: using endpoint %s", ac.Endpoint)
	// logger.Infof("client: poll interval %d", ac.PollInterval)
	// logger.Infof("client: report interval %d", ac.ReportInterval)
	// logger.Infof("client: encryption=%v", cryptor.CanEncrypt())
	// logger.Infof("client: signed messaging=%v", signer.IsSignedMessagingEnabled())
	// logger.Infof("client: rate limit=%v", ac.RateLimit)

	// //iter24: send local IP to server
	// localIP, err := getLocalIP()
	// if err != nil {
	// 	logger.Errorf("error getting local IP-address, %s", err.Error())
	// }
	// agentIP = localIP

	//<-ctx.Done()
	logger.Info("client: shutdown requested")

	// shut down gracefully with timeout of 5 seconds max
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("client: stopped")
}

func naIfEmpty(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
