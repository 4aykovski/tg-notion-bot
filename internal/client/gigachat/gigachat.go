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
	"github.com/4aykovski/tg-notion-bot/lib/helpers"
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

func New(token string) (*Client, error) {
	if token == "" {
		return nil, helpers.ErrWrapIfNotNil("can't create gigachat client", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:     config.GigaChatHost,
		basePath: config.GigaChatAPIBasePath,
		token:    token,
		client:   http.Client{},
	}, nil
}

func (c *Client) Completions(text string) (result string, err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't get completion", err) }()

	rB := newRequestBody(text)

	body, err := json.Marshal(rB)
	if err != nil {
		return "", err
	}

	r, err := c.doRequest(completionsMethod, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	var respBody responseBody
	if err := json.Unmarshal(r, &respBody); err != nil {
		return "", err
	}

	return respBody.Choices[0].Message.Content, nil
}

func (c *Client) doRequest(method string, requestBody io.Reader) (result []byte, err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't do request to gigachat", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentTypeJson)
	req.Header.Add("Accept", contentTypeJson)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
