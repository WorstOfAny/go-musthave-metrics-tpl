package client

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"

	"github.com/WorstOfAny/go-musthave-metrics-tpl/internal/service"
)

type Client struct {
	client     *resty.Client
	bufferPool *sync.Pool
}

func NewClient(baseURL string, key string) *Client {
	restyC := resty.New()
	restyC.
		SetBaseURL(baseURL).
		SetTimeout(time.Second * 15).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			if key != "" {
				data, ok := req.Body.([]byte)
				if !ok {
					return fmt.Errorf("failed type assertion")
				}
				result := service.SignToString(data, key)
				req.SetHeader("HashSHA256", result)
			}
			return nil
		}).
		SetRetryCount(3).
		SetRetryMaxWaitTime(15 * time.Second).
		SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
			delay := time.Duration(2*(r.Request.Attempt)-1) * time.Second
			return delay, nil
		}).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				var netErr net.Error
				log.Debug().Err(err).Msg("request err")
				if errors.As(err, &netErr) && netErr.Timeout() {
					return true
				}
				return errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.EPIPE) || errors.Is(err, io.EOF)
			}
			return r.StatusCode() == http.StatusTooManyRequests || r.StatusCode() > 499
		}).
		AddRetryHook(func(r *resty.Response, err error) {
			if err != nil {
				log.Debug().Str("attempt", strconv.Itoa(r.Request.Attempt)).Err(err).Msg("Request attempt")
			}
		})

	cl := &Client{client: restyC}

	cl.bufferPool = &sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}

	return cl
}

func (c *Client) Post(action string, body []byte) error {
	cbody := c.bufferPool.Get().(*bytes.Buffer)
	gw, err := gzip.NewWriterLevel(cbody, gzip.BestCompression)
	if err != nil {
		return err
	}

	gw.Write(body)
	gw.Close()
	log.Info().Str("body", string(body)).Int("len", len(body)).Str("compressed_body", cbody.String()).Int("comprassed_len", len(cbody.Bytes())).Msg("")

	resp, err := c.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(cbody.Bytes()).
		Post(action)

	cbody.Reset()
	c.bufferPool.Put(cbody)

	if err != nil {
		log.Debug().Err(err).Msg("request err")
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
