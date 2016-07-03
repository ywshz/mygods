package swiss

import (
	"encoding/base64"
	"time"
)

type Config struct {
	NodeName          string
	Backend           string
	BackendMachines   []string
	Tags              map[string]string
	ReconnectInterval time.Duration
	ReconnectTimeout  time.Duration
	EncryptKey        string
	Keyspace          string

	MailHost          string
	MailPort          uint16
	MailUsername      string
	MailPassword      string
	MailFrom          string
}

// This is the default port that we use for Serf communication
const DefaultBindPort int = 8946

func NewConfig() *Config {
	tags := make(map[string]string)
	tags["role"] = "server"
	return &Config{
		NodeName : "Server",
		Backend : "etcd",
		BackendMachines  : []string{"127.0.0.1:2379"},
		Tags                  :tags,
		Keyspace              :"swiss",
	}
}

// EncryptBytes returns the encryption key configured.
func (c *Config) EncryptBytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(c.EncryptKey)
}
