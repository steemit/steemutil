package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/steemit/steemutil/protocol/api"
)

// ClientOption configures a JsonRpc client created by NewClientWithOptions.
type ClientOption func(*JsonRpc)

// WithHTTPClient sets the http.Client used to send requests, giving callers
// control over timeouts, transport, and connection pooling. A nil client is
// a no-op: requests then use the legacy fallback — a fresh client with a 30s
// timeout per call, with no connection-reuse guarantees. Inject a client
// when you need pooling or a custom transport, e.g.:
//
//	tr := http.DefaultTransport.(*http.Transport).Clone()
//	tr.MaxIdleConnsPerHost = 16
//	client := jsonrpc2.NewClientWithOptions(url, jsonrpc2.WithHTTPClient(
//		&http.Client{Timeout: 15 * time.Second, Transport: tr}))
func WithHTTPClient(c *http.Client) ClientOption {
	return func(j *JsonRpc) {
		if c == nil {
			return
		}
		j.client = c
	}
}

// NewClientWithOptions creates a JsonRpc client configured with opts.
func NewClientWithOptions(url string, opts ...ClientOption) *JsonRpc {
	j := &JsonRpc{
		Url: url,
	}
	for _, opt := range opts {
		opt(j)
	}
	return j
}

// Send sends the built request with context.Background(). It is kept for
// backward compatibility; new callers should prefer SendWithContext.
func (j *JsonRpc) Send() (*api.RpcResultData, error) {
	return j.SendWithContext(context.Background())
}

// SendWithContext sends the built request carrying ctx, so the caller can
// cancel the in-flight HTTP request or bound its duration. The http.Client
// injected via WithHTTPClient is used when set; otherwise a client with a 30s
// timeout is created per call (the pre-option behavior).
func (j *JsonRpc) SendWithContext(ctx context.Context) (*api.RpcResultData, error) {
	bodyReader := bytes.NewReader(j.SendData)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, j.Url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := j.client
	if client == nil {
		// Legacy fallback (pre-option behavior): a fresh 30s-timeout client
		// per call. Note this gives no connection-reuse guarantees across
		// calls — inject a client via WithHTTPClient when you need pooling
		// or a custom transport.
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to response(http code): %v", res.StatusCode)
	}
	result := &api.RpcResultData{}
	err = json.NewDecoder(res.Body).Decode(result)
	return result, err
}
