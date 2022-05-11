package tests

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client *http.Client
}

func NewTestSuite() *TestSuite {
	return &TestSuite{client: http.DefaultClient}
}

func (s TestSuite) DoRequest(t *testing.T, url string, width, height int) (*http.Response, []byte, error) {

	u := fmt.Sprintf("http://image-previewer/fill/%d/%d/%s", width, height, url)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)
	require.NoError(t, err)

	res, err := s.client.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	return res, b, err
}

func TestFill(t *testing.T) {
	s := NewTestSuite()

	url := "http://nginx:80/gopher.jpg"
	width, height := 333, 666

	// nolint:bodyclose
	res, body, err := s.DoRequest(t, url, width, height)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.True(t, res.Header.Get("Content-Type") == "image/jpeg")

	config, _, err := image.DecodeConfig(bytes.NewReader(body))
	require.NoError(t, err)

	require.Equal(t, config.Width, width)
	require.Equal(t, config.Height, height)
}

func TestServerDoesntExist(t *testing.T) {
	s := NewTestSuite()

	url := "http://not_exist.com/gopher.jpg"
	width, height := 333, 666

	// nolint:bodyclose
	res, _, err := s.DoRequest(t, url, width, height)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
}

func TestNotImage(t *testing.T) {
	s := NewTestSuite()

	url := "http://ngingx:80/text.txt"
	width, height := 333, 666

	// nolint:bodyclose
	res, _, err := s.DoRequest(t, url, width, height)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
}

func TestURLWrongScheme(t *testing.T) {
	s := NewTestSuite()

	url := "ftp://ngingx:80/gopher.jpg"
	width, height := 333, 666

	// nolint:bodyclose
	res, body, err := s.DoRequest(t, url, width, height)
	require.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, res.StatusCode)
	require.True(t, strings.Contains(string(body), "got not supported scheme"))
}
