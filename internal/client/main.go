package client

import(
	"time"
	"io"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	restyC := resty.New()
	restyC.SetTimeout(time.Second * 10)

	return &Client{ client: restyC }
}

func (c *Client) Post(u string, body []byte) {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(u)
	if err != nil { return }
	defer resp.RawBody().Close()
	io.ReadAll(resp.RawBody())
}
