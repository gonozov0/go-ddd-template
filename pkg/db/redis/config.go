package redis

type Config struct {
	TlsEnabled  bool
	Addrs       []string
	Username    string
	Password    string
	TlsRootCert string
}
