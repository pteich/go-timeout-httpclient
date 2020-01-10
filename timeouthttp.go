package timeouthttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"runtime"
	"time"
)

// Config defines timeouts and max idle connections per host if pooled transport is used
// all timeouts are seconds
type Config struct {
	ConnectTimeout            int
	RequestTimeout            int
	KeepAliveTimeout          int
	MaxIdleConnectionsPerHost int
	KeepAlive                 bool
	tlsConfig                 *tls.Config
}

func DefaultPooledTransport(config Config) *http.Transport {
	return &http.Transport{
		TLSClientConfig: config.tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(config.ConnectTimeout) * time.Second,
			KeepAlive: time.Duration(config.KeepAliveTimeout) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(config.ConnectTimeout) * time.Second,
		MaxIdleConnsPerHost:   config.MaxIdleConnectionsPerHost,
		ResponseHeaderTimeout: time.Duration(config.RequestTimeout) * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     !config.KeepAlive,
	}
}

func DefaultTransport(config Config) *http.Transport {
	transport := DefaultPooledTransport(config)
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

func setDefaults(config *Config) {
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = 1
	}

	if config.RequestTimeout == 0 {
		config.RequestTimeout = 2 * config.ConnectTimeout
	}

	if config.MaxIdleConnectionsPerHost == 0 {
		config.MaxIdleConnectionsPerHost = runtime.GOMAXPROCS(0) + 1
	}
}

// New returns a new clean HTTP.Client with variable options set
func New(opts ...Option) *http.Client {

	config := Config{}
	setDefaults(&config)

	for _, opt := range opts {
		opt(&config)
	}

	return &http.Client{
		Transport: DefaultTransport(config),
		Timeout:   time.Duration(config.RequestTimeout) * time.Second,
	}
}

// NewClient returns a new clean HTTP.Client with timeouts (default 1s for connection and request), disabled idle connections
// and disabled keep-alives
func NewClient(config Config) *http.Client {

	setDefaults(&config)

	return &http.Client{
		Transport: DefaultTransport(config),
		Timeout:   time.Duration(config.RequestTimeout) * time.Second,
	}
}

// NewPooledClient returns a new clean HTTP.Client with timeouts  (default 1s for connection and request) and shared transport
// across hosts with keepalive on, you can set the number of idle connections per host
// with Config.MaxIdleConnsPerHost (default 1)
func NewPooledClient(config Config) *http.Client {

	setDefaults(&config)

	return &http.Client{
		Transport: DefaultPooledTransport(config),
		Timeout:   time.Duration(config.RequestTimeout) * time.Second,
	}
}
