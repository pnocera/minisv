// +build !windows

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func httpInit() {
	srv := httpStart()

	// wait for signal to restart server
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGUSR1)
	for range signalChan {
		log.Println("Restarting http server...")
		srv.Shutdown(context.Background())
		srv = httpStart()
	}
}