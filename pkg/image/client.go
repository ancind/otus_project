package image

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	ErrNotSupportedContentType = errors.New("got not supported content type")
	ErrNotSupportedScheme      = errors.New("got not supported scheme")
	SupportedContentTypes      = []string{"image/jpeg"}
)

type Getter interface {
	Get(ctx context.Context, url string, header http.Header) ([]byte, error)
}

type HTTPGetter struct {
	logger         zerolog.Logger
	transport      http.RoundTripper
	requestTimeout time.Duration
}

func NewImageGetter(l zerolog.Logger, connectTimeout time.Duration, requestTimeout time.Duration) *HTTPGetter {
	return &HTTPGetter{
		logger:         l,
		requestTimeout: requestTimeout,
		transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: connectTimeout,
			}).DialContext,
		},
	}
}

func (f HTTPGetter) Get(ctx context.Context, url string, header http.Header) ([]byte, error) {
	proxyRequest, err := prepareRequest(ctx, url, header)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare request")
	}
	responseBody, err := f.doRequest(proxyRequest)
	if err != nil {
		return nil, errors.Wrap(err, "error making request")
	}
	return responseBody, nil
}

func prepareRequest(ctx context.Context, rawURL string, header http.Header) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy request")
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse image url")
	}
	if parsedURL.Scheme != "http" {
		return nil, ErrNotSupportedScheme
	}
	request.URL = parsedURL
	request.Header = header
	return request, nil
}

func (f *HTTPGetter) doRequest(request *http.Request) ([]byte, error) {
	client := http.Client{
		Timeout:   f.requestTimeout,
		Transport: f.transport,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request")
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			f.logger.Error().Err(errClose).Msg("failed to close body")
		}
	}()

	if !isSupported(resp.Header.Get("Content-type")) {
		return nil, ErrNotSupportedContentType
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	return buff, nil
}

func isSupported(s string) bool {
	for _, val := range SupportedContentTypes {
		if strings.EqualFold(val, s) {
			return true
		}
	}
	return false
}
