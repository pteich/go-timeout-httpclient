package timeouthttp

import "crypto/tls"

type Option func(c *Config)

func WithTimeout(timeout int) Option {
	return func(c *Config) {
		c.ConnectTimeout = timeout
		c.RequestTimeout = timeout
		c.KeepAliveTimeout = timeout
	}
}

func WithConnectTimeout(timeout int) Option {
	return func(c *Config) {
		c.ConnectTimeout = timeout
	}
}

func WithRequestTimeout(timeout int) Option {
	return func(c *Config) {
		c.RequestTimeout = timeout
	}
}

func WithKeepAliveTimeout(timeout int) Option {
	return func(c *Config) {
		c.KeepAliveTimeout = timeout
	}
}

func WithMaxIdleConnections(count int) Option {
	return func(c *Config) {
		c.MaxIdleConnectionsPerHost = count
	}
}

func WithTlsConfig(tlsConfig *tls.Config) Option {
	return func(c *Config) {
		c.tlsConfig = tlsConfig
	}
}
