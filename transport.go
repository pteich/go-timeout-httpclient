package timeouthttp

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

type Transport struct {
	ht                   http.Transport
	enableCircuitBreaker bool
	breaker              sync.Map
}

func DefaultPooledTransport(config Config) *Transport {
	return &Transport{
		ht: http.Transport{
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
		},
		enableCircuitBreaker: config.circuitBreaker,
	}
}

func DefaultTransport(config Config) *Transport {
	transport := DefaultPooledTransport(config)
	transport.ht.DisableKeepAlives = true
	transport.ht.MaxIdleConnsPerHost = -1
	return transport
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if !t.enableCircuitBreaker {
		return t.ht.RoundTrip(r)
	}

	cbValue, _ := t.breaker.LoadOrStore(r.Host, gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: r.Host,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}))

	cb := cbValue.(*gobreaker.CircuitBreaker)
	res, err := cb.Execute(func() (interface{}, error) {
		res, err := t.ht.RoundTrip(r)
		if err != nil {
			return nil, err
		}
		if res != nil && res.StatusCode >= http.StatusBadRequest {
			return res, fmt.Errorf("http response error: %d", res.StatusCode)
		}

		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*http.Response), nil
}
