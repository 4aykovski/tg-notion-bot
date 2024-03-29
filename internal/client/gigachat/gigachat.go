package gigachat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
)

const (
	completionsMethod = "/chat/completions"
	authUrl           = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	dataScope         = "GIGACHAT_API_PERS"
	rqUid             = "6f0b1291-c7f3-43c6-bb2e-9f3efb2dc98e"
)

var (
	ErrCantGetCompletion         = errors.New("can't get completion")
	ErrCantSendCompletionRequest = errors.New("can't send completion request")
)

type Client struct {
	hTTPClient *client.HTTPClient
	token      string
	auth       string
}

func New(cfg config.GigaChatConfig) (*Client, error) {
	if cfg.Token == "" && cfg.Auth == "" {
		return nil, fmt.Errorf("gigachat %w: %w", client.ErrCantCreateClient, client.ErrAuthInfoNotSpecified)
	}

	return &Client{
		hTTPClient: client.NewHTTPClient(cfg.Host, cfg.APIBasePath),
		token:      cfg.Token,
		auth:       cfg.Auth,
	}, nil
}

func (c *Client) Completions(text string) (string, error) {

	rB := newRequestBody(text)

	jsonRB, err := json.Marshal(rB)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantGetCompletion, err)
	}

	res, err := c.completionRequest(jsonRB)
	if err != nil && !errors.Is(err, client.Err401StatusCode) {
		return "", fmt.Errorf("%w: %w", ErrCantGetCompletion, err)
	} else if errors.Is(err, client.Err401StatusCode) {
		if err = c.updateToken(); err != nil {
			return "", fmt.Errorf("%w: %w", ErrCantGetCompletion, err)
		}

		res, err = c.completionRequest(jsonRB)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrCantGetCompletion, err)
		}
	}

	var result responseBody
	if err := json.Unmarshal(res, &result); err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantGetCompletion, err)
	}

	return result.Choices[0].Message.Content, nil
}

func (c *Client) completionRequest(jsonBody []byte) ([]byte, error) {
	u := c.hTTPClient.GetUlrWithMethods(completionsMethod)
	header := http.Header{
		"Authorization": []string{"Bearer " + c.token},
		"Content-Type":  []string{client.ContentTypeJson},
		"Accept":        []string{client.ContentTypeJson},
	}

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u.String(), header, strings.NewReader(string(jsonBody)), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantSendCompletionRequest, err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantSendCompletionRequest, err)
	}

	return res, nil
}

func (c *Client) updateToken() error {
	u := authUrl

	header := http.Header{}
	header.Add("Authorization", "Bearer "+c.auth)
	header.Add("Content-Type", client.ContentTypeUrlEncoded)
	header.Add("RqUID", rqUid)
	header.Add("Accept", client.ContentTypeJson)

	body := url.Values{}
	body.Set("scope", dataScope)

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u, header, strings.NewReader(body.Encode()), nil)
	if err != nil {
		return fmt.Errorf("%w: %w", client.ErrCantUpdateToken, err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", client.ErrCantUpdateToken, err)
	}

	var result responseBody
	err = json.Unmarshal(res, &result)
	if err != nil {
		return fmt.Errorf("%w: %w", client.ErrCantUpdateToken, err)
	}

	c.token = result.AccessToken
	return nil
}
