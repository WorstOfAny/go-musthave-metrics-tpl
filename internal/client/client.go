package client

import(
	"time"
	"io"
	"github.com/go-resty/resty/v2"
	"compress/gzip"
	"bytes"
	"github.com/rs/zerolog/log"
)

type Client struct {
	client *resty.Client
}

func NewClient() *Client {
	restyC := resty.New()
	restyC.SetTimeout(time.Second * 10)

	return &Client{ client: restyC }
}

func (c *Client) Post(u string, body []byte) error {
	var cbody bytes.Buffer
	gw, err := gzip.NewWriterLevel(&cbody, gzip.BestCompression)
	if err != nil { return err }

	gw.Write(body)
	gw.Close()
	log.Info().Str("body", string(body)).Int("len", len(body)).Str("compressed_body", cbody.String()).Int("comprassed_len", len(cbody.Bytes())).Msg("")

	resp, err := c.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(cbody.Bytes()).
		Post(u)

	if err != nil {
		log.Debug().Err(err).Msg("response err")
		return err
	}
	defer resp.RawResponse.Body.Close()
	
	gr, err := gzip.NewReader(resp.RawResponse.Body)
	if err != nil {
		log.Debug().Err(err).Msg("gzip reader error")
		return err
	}
	respbody, err := io.ReadAll(gr)
	if err != nil {
		log.Debug().Err(err).Msg("read body err")
		return err
	}
	log.Info().Str("resp", string(respbody)).Msg("")

	return nil
}
