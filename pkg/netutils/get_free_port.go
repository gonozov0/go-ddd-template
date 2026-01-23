package netutils

import (
	"net"
	"strconv"

	"fmt"
)

func GetFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", fmt.Errorf("failed to resolve tcp addr: %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen tcp: %w", err)
	}
	defer l.Close()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return "", fmt.Errorf("failed to cast to TCPAddr: %s", l.Addr())
	}

	return strconv.Itoa(addr.Port), nil
}
