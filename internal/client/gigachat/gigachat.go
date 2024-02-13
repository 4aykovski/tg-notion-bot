package gigachat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
)

const (
	contentTypeJson   = "application/json"
	completionsMethod = "/chat/completions"
)

type Client struct {
	host     string
	basePath string
	token    string
	client   http.Client
}

func New(cfg config.GigaChatConfig) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("can't create gigachat client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:     cfg.Host,
		basePath: cfg.APIBasePath,
		token:    cfg.Token,
		client:   http.Client{},
	}, nil
}

func (c *Client) Completions(text string) (result string, err error) {

	rB := newRequestBody(text)

	body, err := json.Marshal(rB)
	if err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	r, err := c.doRequest(completionsMethod, strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	var respBody responseBody
	if err := json.Unmarshal(r, &respBody); err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	return respBody.Choices[0].Message.Content, nil
}

func (c *Client) doRequest(method string, requestBody io.Reader) (result []byte, err error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, fmt.Errorf("can't do request to gigachat: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentTypeJson)
	req.Header.Add("Accept", contentTypeJson)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request to gigachat: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't do request to gigachat: %w", fmt.Errorf("wrong status code: %d", res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request to gigachat: %w", err)
	}

	return body, nil
}
