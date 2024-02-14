package salutespeech

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
)

var (
	contentTypeOgg        = "audio/ogg;codecs=opus"
	speechRecognizeMethod = "/speech:recognize"
)

type Client struct {
	hTTPClient          client.HTTPClient
	token               string
	voicesFileDirectory string
}

func New(cfg config.SalutespeechConfig, voicesFileDir string) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("can't create salutespeech client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		hTTPClient:          *client.NewHTTPClient(cfg.Host, cfg.APIBasePath),
		token:               cfg.Token,
		voicesFileDirectory: voicesFileDir,
	}, nil
}

func (c *Client) SpeechRecognizeOgg(fileName string) (text string, err error) {

	f, err := os.Open(filepath.Join(c.voicesFileDirectory, fileName))
	if err != nil {
		return "", fmt.Errorf("can't recognize speech: %w", err)
	}
	defer f.Close()

	u := c.hTTPClient.GetUlrWithMethods(speechRecognizeMethod)
	header := http.Header{}
	header.Add("Authorization", "Bearer "+c.token)
	header.Add("Content-Type", contentTypeOgg)

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u.String(), header, f, nil)
	if err != nil {
		return "", fmt.Errorf("can't recognize speech: %w", err)
	}

	body, err := c.hTTPClient.Do(req)
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
