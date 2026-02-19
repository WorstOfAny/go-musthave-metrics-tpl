package client

import(
	"net/http"
	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/logger"
	"time"
	"strings"
	"io"
)

type Client struct {
	l logger.Logger
	client *http.Client
}

func NewClient(l logger.Logger) *Client {
	return &Client{
		l: l,
		client: &http.Client{ Timeout: time.Second * 1 },
	}
}

func (c *Client) Post(u string) {
	resp, err := c.client.Post(u, "text/plain", strings.NewReader(""))
	if err != nil { return }
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
}
