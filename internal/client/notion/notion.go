package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/4aykovski/tg-notion-bot/config"
	"github.com/4aykovski/tg-notion-bot/internal/client"
)

const (
	contentTypeJson  = "application/json"
	createPageMethod = "/pages"
)

type Client struct {
	hTTPClient    client.HTTPClient
	notionVersion string
	token         string
}

type JsonAnswer struct {
	Result  string `json:"result,omitempty"`
	Summary string `json:"summary,omitempty"`
	Name    string `json:"name,omitempty"`
}

func New(cfg config.NotionConfig) (*Client, error) {
	if cfg.IntegrationToken == "" {
		return nil, fmt.Errorf("can't create notion client: %w", fmt.Errorf("token wasn't specified"))
	}
	return &Client{
		hTTPClient:    *client.NewHTTPClient(cfg.Host, cfg.APIBasePath),
		notionVersion: cfg.Version,
		token:         cfg.IntegrationToken,
	}, nil
}

func (c *Client) CreateNewPageInDatabase(dbId string, pageData string) (err error) {

	var pData JsonAnswer
	if err = json.Unmarshal([]byte(pageData), &pData); err != nil {
		return fmt.Errorf("can't create notion page: %w", err)
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
		return fmt.Errorf("can't create notion page: %w", err)
	}

	body := strings.NewReader(string(jsonPage))
	u := c.hTTPClient.GetUlrWithMethods(createPageMethod)
	header := http.Header{}
	header.Add("Content-Type", contentTypeJson)
	header.Add("Authorization", "Bearer "+c.token)
	header.Add("Notion-Version", c.notionVersion)

	req, err := c.hTTPClient.CreateRequest(http.MethodPost, u.String(), header, body, nil)
	if err != nil {
		return fmt.Errorf("can't create notion page: %w", err)
	}

	_, err = c.hTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("can't create notion page: %w", err)
	}

	return nil
}
