package handler

import (
	"github.com/gtrxshock/sentinel-proxy/internal/app/core"
	"os"
	"syscall"
)

func SignalListener(signalChan chan os.Signal, termChan chan bool) {
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGINT,
			syscall.SIGTERM:
			termChan <- true
			return
		default:
			core.GetLogger().Info("received unknown signal", sig)
		}
	}
}
