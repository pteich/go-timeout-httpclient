# go-timeout-httpclient
Yet another HTTP Client with configurable timeouts and TLS options.

## Usage
Create a HTTP client with specific timeouts once and re-use it wherever you need it.  

```go

import "github.com/pteich/go-timeout-httpclient"

// Create an HTTP client with a timeout of 5s
httpClient := timeouthttp.New(timeouthttp.WithTimeout(5))
```

## Built-in Circuit Breaker 

Using the option `timeouthttp.WithCircuitBreaker()` automatically enables a circuit breaker
for the requested host. This is implemented using the `http.RoundTripper` interface.

By now there are no further configuration values possible. This means the circuit is considered open
if there are more than 3 requests and a failure ratio > 0.6.

Caveat: A new circuit breaker is created for every new host. There will be a memory issue
if you request a lot of different hosts as there is no cleaning build in for the time being.

