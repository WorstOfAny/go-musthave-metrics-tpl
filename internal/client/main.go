package client

import(
	"net/http"
	"time"
	"strings"
	"io"
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
	resp, err := c.client.Post(u, "text/plain", strings.NewReader(""))
	if err != nil { return }
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}
