package postgres

import (
	"fmt"
	"net"
	"strings"
)

type Config struct {
	Hosts       []string
	Database    string
	User        string
	Password    string
	Port        string
	SSL         bool
	SSLRootCert string
}

func (c Config) GetConnStr(host string) string {
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
		if c.SSLRootCert != "" {
			connStr += " sslrootcert=" + c.SSLRootCert
		}
	}

	return connStr
}

func (c Config) GetConnURL(host string) string {
	params := make([]string, 0, 2)

	sslMode := "disable"
	if c.SSL {
		sslMode = "verify-full"

		if c.SSLRootCert != "" {
			params = append(params, "sslrootcert="+c.SSLRootCert)
		}
	}

	params = append(params, "sslmode="+sslMode)

	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?%s",
		c.User,
		c.Password,
		net.JoinHostPort(host, c.Port),
		c.Database,
		strings.Join(params, "&"),
	)
}
