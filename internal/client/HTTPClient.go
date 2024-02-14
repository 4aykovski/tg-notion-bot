package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type HTTPClient struct {
	Host     string
	BasePath string
	client   http.Client
}

func NewHTTPClient(host, basePath string) *HTTPClient {
	return &HTTPClient{
		Host:     host,
		BasePath: basePath,
		client:   http.Client{},
	}
}

func (hc *HTTPClient) Do(r *http.Request) ([]byte, error) {
	res, err := hc.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("wrong status code on request to %s", res.Request.URL.String()))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("can't read response body: %w", err))
	}

	return body, nil
}

// CreateRequest return http.Request with given parameters. If you don't need some of the parameters then give nil.
func (hc *HTTPClient) CreateRequest(method string, url string, header http.Header, body io.Reader, query url.Values) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	if header != nil {
		req.Header = header
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func (hc *HTTPClient) GetFullUrl() *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   hc.Host,
		Path:   path.Join(hc.BasePath, "/"),
	}
}

func (hc *HTTPClient) GetUlrWithMethods(method string) *url.URL {
	u := hc.GetFullUrl()
	u.Path = path.Join(u.Path, method)
	return u
}
