package service

import (
	"fmt"
	"net"
)

type DbConnector struct{}

func (dc *DbConnector) Listen(localPort int) (net.Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", localPort))
	if err != nil {
		return nil, err
	}

	return listener, nil
}
