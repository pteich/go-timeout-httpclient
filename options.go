package timeouthttp

import (
	"crypto/tls"
	"net/http"
)

type Option func(c *Config)

// WithTimeout sets all timeouts to the provided value in seconds
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

func WithCircuitBreaker() Option {
	return func(c *Config) {
		c.circuitBreaker = true
	}
}

func WithPooledTransport(maxIdleConnectionsPerHost int) Option {
	return func(c *Config) {
		c.pooledTransport = true
		c.MaxIdleConnectionsPerHost = maxIdleConnectionsPerHost
	}
}

func WithTransport(t http.RoundTripper) Option {
	return func(c *Config) {
		c.transport = t
	}
}
