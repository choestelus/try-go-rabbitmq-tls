package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Config struct {
	CAPaths  []string `required:"true" envconfig:"CA_PATHS"`
	CertPath string   `required:"true" envconfig:"CERT_PATH"`
	KeyPath  string   `required:"true" envconfig:"KEY_PATH"`
}

func (c Config) TLSConfig() *tls.Config {
	// Load CA
	pool := x509.NewCertPool()
	for _, ca := range c.CAPaths {
		pem, err := ioutil.ReadFile(ca)
		if err != nil {
			panic(err)
		}

		ok := pool.AppendCertsFromPEM(pem)
		if !ok {
			panic("failed to parse CA PEM into cert pool")
		}
	}

	// Load certificate
	cert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		// set of CA that server will trust, which client certificates signed with, useful only for http server.
		// Below is description from official documentation:
		//
		// ClientCAs defines the set of root certificate authorities
		// that servers use if required to verify a client certificate
		// by the policy in ClientAuth.
		ClientCAs: pool,
		// set of CA that client will use when verifying server certificates
		// Below is description from official documentation:
		//
		// RootCAs defines the set of root certificate authorities
		// that clients use when verifying server certificates.
		// If RootCAs is nil, TLS uses the host's root CA set.
		RootCAs: pool,
		// server policy of TLS Client Authentication, only useful for use with http server
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

}

func main() {
	cfg := &Config{}
	envconfig.MustProcess("", cfg)
	tlsConfig := cfg.TLSConfig()
	// rabbitmq server certificate canonical name
	// useful when connecting via localhost or raw ip
	tlsConfig.ServerName = "server.mq.choe"

	conn, err := amqp.DialTLS("amqps://guest:guest@localhost:5671", tlsConfig)
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Info("expecting TLSv1.2")
	switch conn.ConnectionState().Version {
	case tls.VersionSSL30:
		logrus.Panic("got SSLv3.0")
	case tls.VersionTLS11:
		logrus.Panic("got TLSv1.1")
	case tls.VersionTLS12:
		logrus.Info("got TLSv1.2")
	case tls.VersionTLS13:
		logrus.Panic("got TLSv1.3")
	default:
		logrus.Panicf("got unexpected TLS version: %x", conn.ConnectionState().Version)
	}
}
