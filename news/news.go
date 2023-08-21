package news

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

// Article construct by using https://mholt.github.io/json-to-go/
type Article []struct {
	Source struct {
		ID   any    `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// NewClient is for create a new client obj working with the News API (where we request)
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}
	return &Client{httpClient, key, pageSize}
}

func (c *Client) FetchEverything(query, page string) (*Results, error) {
	// filled & do get request
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en",
		url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	// defer the close
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	// return & Unmarshal Json -> Results
	res := &Results{}
	return res, json.Unmarshal(body, res)
}
