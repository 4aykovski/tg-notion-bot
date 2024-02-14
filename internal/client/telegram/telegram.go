package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
	"github.com/cavaliergopher/grab/v3"
)

const (
	getUpdatesMethod  = "/getUpdates"
	sendMessageMethod = "/sendMessage"
	getFileMethod     = "/getFile"
)

type Client struct {
	hTTPClient          client.HTTPClient
	voicesFileDirectory string
}

func New(cfg config.TelegramConfig, voicesFileDir string) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("can't create telegram client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		hTTPClient:          *client.NewHTTPClient(cfg.Host, newBasePath(cfg.Token)),
		voicesFileDirectory: voicesFileDir,
	}, nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	u := c.hTTPClient.GetUlrWithMethods(getUpdatesMethod)

	req, err := c.hTTPClient.CreateRequest(http.MethodGet, u.String(), nil, nil, q)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	var updRes UpdatesResponse

	if err := json.Unmarshal(res, &updRes); err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	return updRes.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	u := c.hTTPClient.GetUlrWithMethods(sendMessageMethod)

	req, err := c.hTTPClient.CreateRequest(http.MethodGet, u.String(), nil, nil, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	_, err = c.hTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *Client) FileInfo(fileId string) (*File, error) {
	q := url.Values{}
	q.Add("file_id", fileId)

	u := c.hTTPClient.GetUlrWithMethods(getFileMethod)

	req, err := c.hTTPClient.CreateRequest(http.MethodGet, u.String(), nil, nil, q)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	res, err := c.hTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	body, err := c.fileResponse(res)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	return &File{
		FileId:   body.Result.FileId,
		FilePath: body.Result.FilePath,
	}, nil
}

func (c *Client) DownloadFile(filePath string) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.hTTPClient.Host,
		Path:   path.Join("file", c.hTTPClient.BasePath, filePath),
	}

	_, err := grab.Get(c.voicesFileDirectory, u.String())
	if err != nil {
		return fmt.Errorf("can't download file: %w", err)
	}

	return nil
}

func (c *Client) fileResponse(data []byte) (*GetFileResponse, error) {
	var r GetFileResponse

	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, fmt.Errorf("can't deserialize json response: %w", err)
	}

	return &r, nil
}
