package salutespeech

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/4aykovski/tg-notion-bot/config"
)

var (
	contentTypeOgg        = "audio/ogg;codecs=opus"
	speechRecognizeMethod = "/speech:recognize"
)

type Client struct {
	host     string
	basePath string
	token    string
	client   http.Client
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("can't create salutespeech client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:     config.SalutespeechHost,
		basePath: config.SalutesleepchAPIBasePath,
		token:    token,
		client:   http.Client{},
	}, nil
}

func (c *Client) SpeechRecognizeOgg(fileName string) (text string, err error) {

	f, err := os.Open(filepath.Join(config.VoicesFileDirectory, fileName))
	if err != nil {
		return "", fmt.Errorf("can't recognize speech: %w", err)
	}
	defer f.Close()

	body, err := c.doPostRequest(speechRecognizeMethod, contentTypeOgg, f)
	if err != nil {
		return "", fmt.Errorf("can't recognize speech: %w", err)
	}

	var result Response
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("can't recognize speech: %w", err)
	}

	if result.StatusCode != http.StatusOK {
		return "", fmt.Errorf("can't recognize speech:%w", fmt.Errorf("wrong status code: %d", result.StatusCode))
	}

	return result.Result[0], nil
}

func (c *Client) doPostRequest(method string, contentType string, requestBody io.Reader) (data []byte, err error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), requestBody)
	if err != nil {
		return nil, fmt.Errorf("can't do post request on salutespeech: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do post request on salutespeech: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("can't do post request on salutespeech: %w", fmt.Errorf("wrong status code: %d", res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do post request on salutespeech: %w", err)
	}

	return body, nil
}
