package salutespeech

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
)

const (
	authUrl               = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	rqUid                 = "e6c60511-cc3e-4c08-91de-94f548969dda"
	speechRecognizeMethod = "/speech:recognize"
	dataScope             = "SALUTE_SPEECH_PERS"
)

var (
	ErrCantRecognizeSpeech            = errors.New("can't recognize speech")
	ErrCantSendSpeechRecognizeRequest = errors.New("can't create speech recognize request")
)

type Client struct {
	hTTPClient          *client.HTTPClient
	token               string
	auth                string
	voicesFileDirectory string
}

func New(cfg config.SalutespeechConfig, voicesFileDir string) (*Client, error) {
	if cfg.Token == "" && cfg.Auth == "" {
		return nil, fmt.Errorf("can't create salutespeech client: %w", client.ErrAuthInfoNotSpecified)
	}

	return &Client{
		hTTPClient:          client.NewHTTPClient(cfg.Host, cfg.APIBasePath),
		token:               cfg.Token,
		auth:                cfg.Auth,
		voicesFileDirectory: voicesFileDir,
	}, nil
}

func (c *Client) SpeechRecognizeOgg(fileName string) (text string, err error) {
	res, err := c.speechRecognizeRequest(fileName)
	if err != nil && !errors.Is(err, client.Err401StatusCode) {
		return "", fmt.Errorf("%w: %w", ErrCantRecognizeSpeech, err)
	} else if errors.Is(err, client.Err401StatusCode) {
		if err := c.updateToken(); err != nil {
			return "", fmt.Errorf("%w: %w", ErrCantRecognizeSpeech, err)
		}

		res, err = c.speechRecognizeRequest(fileName)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrCantRecognizeSpeech, err)
		}
	}

	var result Response
	if err = json.Unmarshal(res, &result); err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantRecognizeSpeech, err)
	}

	if result.StatusCode != http.StatusOK {
		return "", fmt.Errorf("can't recognize speech:%w", fmt.Errorf("wrong status code: %d", result.StatusCode))
	}

	return result.Result[0], nil
}

func (c *Client) speechRecognizeRequest(fileName string) ([]byte, error) {
	f, err := os.Open(filepath.Join(c.voicesFileDirectory, fileName))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantRecognizeSpeech, err)
	}
	defer f.Close()

	u := c.hTTPClient.GetUlrWithMethods(speechRecognizeMethod)
	header := http.Header{
		"Authorization": []string{"Bearer " + c.token},
		"Content-Type":  []string{client.ContentTypeOgg},
	}

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u.String(), header, f, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantSendSpeechRecognizeRequest, err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantSendSpeechRecognizeRequest, err)
	}

	return res, nil
}

func (c *Client) updateToken() error {
	u := authUrl

	header := http.Header{}
	header.Add("Authorization", "Basic "+c.auth)
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

	var result Response
	err = json.Unmarshal(res, &result)
	if err != nil {
		return fmt.Errorf("%w: %w", client.ErrCantUpdateToken, err)
	}

	c.token = result.AccessToken
	return nil
}
