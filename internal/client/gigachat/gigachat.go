package gigachat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
)

const (
	contentTypeJson   = "application/json"
	completionsMethod = "/chat/completions"
)

type Client struct {
	hTTPClient client.HTTPClient
	token      string
	client     http.Client
}

func New(cfg config.GigaChatConfig) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("can't create gigachat client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		hTTPClient: *client.NewHTTPClient(cfg.Host, cfg.APIBasePath),
		token:      cfg.Token,
		client:     http.Client{},
	}, nil
}

func (c *Client) Completions(text string) (string, error) {

	rB := newRequestBody(text)

	jsonRB, err := json.Marshal(rB)
	if err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	u := c.hTTPClient.GetUlrWithMethods(completionsMethod)
	header := http.Header{}
	header.Add("Authorization", "Bearer "+c.token)
	header.Add("Content-Type", contentTypeJson)
	header.Add("Accept", contentTypeJson)
	body := strings.NewReader(string(jsonRB))

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u.String(), header, body, nil)
	if err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	var result responseBody
	if err := json.Unmarshal(res, &result); err != nil {
		return "", fmt.Errorf("can't get completion: %w", err)
	}

	return result.Choices[0].Message.Content, nil
}
