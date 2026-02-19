package client

import(
	"net/http"
	"time"
	"io"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{ Timeout: time.Second * 1 },
	}
}

func (c *Client) Post(u string) {
	resp, err := c.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "text/plain").
		Post(u)
	if err != nil { return }
	defer resp.RawBody().Close()
	io.ReadAll(resp.RawBody())
}
