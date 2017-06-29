package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
	)

	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT:
				cancel()
			}
		}
	}()
}
