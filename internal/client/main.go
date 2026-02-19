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
	restyC.SetTimeout(time.Second * 1)
	restyC.SetDebug(true)
	return &Client{ client: restyC}
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
