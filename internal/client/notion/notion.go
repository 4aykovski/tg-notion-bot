package notion

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

var (
	contentTypeJson  = "application/json"
	createPageMethod = "/pages"
)

type Client struct {
	host     string
	basePath string
	token    string
	client   http.Client
}

type JsonAnswer struct {
	Result  string `json:"result,omitempty"`
	Summary string `json:"summary,omitempty"`
	Name    string `json:"name,omitempty"`
}

func New(token string) (*Client, error) {
	if token == "" {
		return nil, helpers.ErrWrapIfNotNil("can't create notion client", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		host:     config.NotionHost,
		basePath: config.NotionAPIBasePath,
		token:    token,
		client:   http.Client{},
	}, nil
}

func (c *Client) CreateNewPageInDatabase(dbId string, pageData string) (err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't create notion page", err) }()

	var pData JsonAnswer
	if err = json.Unmarshal([]byte(pageData), &pData); err != nil {
		return err
	}

	pageParent := newDatabaseParent(dbId)
	pageProperties := map[string]property{
		"Name": *newTitleProperty(pData.Name),
	}
	pageChildren := []block{
		*newHeading2Block(pData.Summary),
		*newParagraphBlock(pData.Result),
	}
	p := newPage(*pageParent, pageProperties, pageChildren)

	jsonPage, err := json.Marshal(*p)
	if err != nil {
		return err
	}

	body := strings.NewReader(string(jsonPage))
	if err = c.doRequest(createPageMethod, body); err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(method string, body io.Reader) (err error) {
	defer func() { err = helpers.ErrWrapIfNotNil("can't do request to notion", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", contentTypeJson)
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Notion-Version", config.NotionVersion)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code: %d", res.StatusCode)
	}

	return nil
}
