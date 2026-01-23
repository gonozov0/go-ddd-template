package netutils

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"fmt"
)

type ServerInfo struct {
	tlsEnabled bool
	domain     string
	GRPCPort   string
	HTTPPort   string
}

func (s ServerInfo) GetHTTPURL() string {
	httpPrefix := s.getHTTPPrefix()

	if s.HTTPPort == "" {
		return httpPrefix + s.domain
	}

	return "http://" + s.domain + ":" + s.HTTPPort
}

func (s ServerInfo) getHTTPPrefix() string {
	if s.tlsEnabled {
		return "https://"
	}

	return "http://"
}

func (s ServerInfo) MakeGRPCConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		s.getGRPCURL(),
		grpc.WithTransportCredentials(s.getGRPCCreds()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client: %w", err)
	}

	return conn, nil
}

func (s ServerInfo) getGRPCURL() string {
	if s.GRPCPort == "" {
		return s.domain
	}

	return s.domain + ":" + s.GRPCPort
}

func (s ServerInfo) getGRPCCreds() credentials.TransportCredentials {
	if s.tlsEnabled {
		return credentials.NewTLS(&tls.Config{})
	}

	return insecure.NewCredentials()
}

func NewServerInfo(tlsEnabled bool, url, grpcPort, httpPort string) ServerInfo {
	return ServerInfo{
		tlsEnabled: tlsEnabled,
		domain:     url,
		GRPCPort:   grpcPort,
		HTTPPort:   httpPort,
	}
}
