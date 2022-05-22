package timeouthttp

import (
	"net/http"
)

type Client struct {
	http.Client
}
