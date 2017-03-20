# go-timeout-httpclient
Yet another HTTP Client with configurable timeouts

## Usage
```go

import "github.com/pteich/go-timeout-httpclient"

httpClient := timeouthttp.NewClient(timeouthttp.Config{
    RequestTimeout: 5,
    ConnectTimeout: 5,
})
```
