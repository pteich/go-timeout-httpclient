# go-timeout-httpclient
Yet another HTTP Client with configurable timeouts and TLS options.

## Usage
```go

import "github.com/pteich/go-timeout-httpclient"

httpClient := timeouthttp.New(timeouthttp.WithTimeout(5))
```
