package main

import (
	"fmt"
	"github.com/gtrxshock/sentinel-proxy/internal/app"
	"github.com/gtrxshock/sentinel-proxy/internal/app/core"
	"github.com/gtrxshock/sentinel-proxy/internal/app/handler"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func init() {
	var configPath string

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exPath := filepath.Dir(ex)
		configPath = exPath + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "parameters.yaml"
	}

	config, err := core.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger := core.GetLogger()

	logger.Infof("sentinel proxy inited with config: %s", configPath)
	app.NewSentinelProxy(config, logger)
}

func main() {
	proxy := app.GetSentinelProxyInstance()
	defer proxy.Close()
	proxy.Logger.Info("proxy started, pid: ", syscall.Getpid())

	termProxyChan := make(chan bool)
	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go handler.SignalListener(signalChan, termProxyChan)

	go func() {
		err := proxy.Start()
		if err != nil {
			proxy.Logger.Fatal("proxy failed: ", err)
		}
	}()

	<-termProxyChan
	proxy.Logger.Info("proxy received term signal, proxy stopped")
}
