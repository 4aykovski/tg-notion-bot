package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/cavaliergopher/grab/v3"
)

const (
	getUpdatesMethod  = "/getUpdates"
	sendMessageMethod = "/sendMessage"
	getFileMethod     = "/getFile"
)

type Client struct {
	host                string
	basePath            string
	client              http.Client
	voicesFileDirectory string
}

func New(cfg config.TelegramConfig, voicesFileDir string) (*Client, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("can't create telegram client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:                cfg.Host,
		basePath:            newBasePath(cfg.Token),
		client:              http.Client{},
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

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *Client) FileInfo(fileId string) (*File, error) {
	q := url.Values{}
	q.Add("file_id", fileId)

	jsonRes, err := c.doRequest(getFileMethod, q)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	res, err := c.fileResponse(jsonRes)
	if err != nil {
		return nil, fmt.Errorf("can't get file: %w", err)
	}

	return &File{
		FileId:   res.Result.FileId,
		FilePath: res.Result.FilePath,
	}, nil
}

func (c *Client) DownloadFile(filePath string) error {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("file/" + c.basePath + "/" + filePath),
	}

	_, err := grab.Get(c.voicesFileDirectory, u.String())
	if err != nil {
		return fmt.Errorf("can't download file: %w", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath + method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't send request: %w", err)
	}

	return body, nil
}

func (c *Client) fileResponse(data []byte) (*GetFileResponse, error) {
	var r GetFileResponse

	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, fmt.Errorf("can't deserialize json response: %w", err)
	}

	return &r, nil
}
