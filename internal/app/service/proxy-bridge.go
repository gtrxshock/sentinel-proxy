package service

import (
	"github.com/gtrxshock/sentinel-proxy/internal/app/core"
	"io"
	"net"
)

type ProxyBridge struct{}

func (pb *ProxyBridge) Proxy(clientConn, redisConn net.Conn) {
	go pb.proxyConnection(clientConn, redisConn)
	go pb.proxyConnection(redisConn, clientConn)
}

func (pb *ProxyBridge) proxyConnection(destConn net.Conn, srcConn net.Conn) {
	_, err := io.Copy(destConn, srcConn)
	if err != nil {
		core.GetLogger().Debug("proxy connection closed")
	}

	_ = destConn.Close()
	_ = srcConn.Close()
}
