package postgres

import (
	"errors"
	"fmt"
	"net"
)

type ConnectionData struct {
	Hosts    []string
	Database string
	User     string
	Password string
	Port     string
	SSL      bool
}

func NewConnectionData(hosts []string, dbname, user, password, port string, ssl bool) (ConnectionData, error) {
	if len(hosts) == 0 {
		return ConnectionData{}, errors.New("no host found")
	}
	if port == "" {
		port = "5432"
	}
	return ConnectionData{
		Database: dbname,
		User:     user,
		Password: password,
		Port:     port,
		SSL:      ssl,
		Hosts:    hosts,
	}, nil
}

func (c ConnectionData) String(host string) string {
	connStr := fmt.Sprintf(
		`host=%v port=%s dbname=%s user=%s password=%s`,
		host,
		c.Port,
		c.Database,
		c.User,
		c.Password,
	)
	if c.SSL {
		connStr += " sslmode=verify-full"
	}
	return connStr
}

func (c ConnectionData) URL(host string) string {
	sslMode := "disable"
	if c.SSL {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.User,
		c.Password,
		net.JoinHostPort(host, c.Port),
		c.Database,
		sslMode,
	)
}
