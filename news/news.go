package news

import "net/http"

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

// NewClient is for create a new client obj working with the News API (where we request)
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}
