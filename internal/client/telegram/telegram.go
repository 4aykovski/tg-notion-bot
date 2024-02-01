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
	"github.com/4aykovski/tg-notion-bot/lib/helpers"
	"github.com/cavaliergopher/grab/v3"
)

const (
	getUpdatesMethod  = "/getUpdates"
	sendMessageMethod = "/sendMessage"
	getFileMethod     = "/getFile"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) (*Client, error) {
	if token == "" {
		return nil, helpers.ErrWrapIfNotNil("can't create telegram client", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}, nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't send message", err)
	}

	return nil
}

func (c *Client) FileInfo(fileId string) (*File, error) {
	q := url.Values{}
	q.Add("file_id", fileId)

	jsonRes, err := c.doRequest(getFileMethod, q)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't get file", err)
	}

	res, err := c.fileResponse(jsonRes)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't get file", err)
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

	_, err := grab.Get(config.VoicesFileDirectory, u.String())
	if err != nil {
		return helpers.ErrWrapIfNotNil("can't download file", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't send request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath + method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) fileResponse(data []byte) (*GetFileResponse, error) {
	var r GetFileResponse

	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, helpers.ErrWrapIfNotNil("can't deserialize json response", err)
	}

	return &r, nil
}