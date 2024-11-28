package client

import (
	"encoding/json"
	"fmt"
	"github.com/Fyefhqdishka/eff-mobile/internal/models"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type ClientInterface interface {
	GetDetails(song string, groupName string) (models.Song, error)
}

type Client struct {
	client  *http.Client
	baseURL string
	log     *slog.Logger
}

func NewClient(baseURL string, log *slog.Logger) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: baseURL,
		log:     log,
	}
}

func (c *Client) GetDetails(song string, groupName string) (models.Song, error) {
	c.log.Debug("client details song=", song, ", group=", groupName)

	url := fmt.Sprintf("http://%s/info?group=%s&song=%s", c.baseURL, url.QueryEscape(groupName), url.QueryEscape(song))
	resp, err := c.client.Get(url)
	if err != nil {
		c.log.Error("failed to make request:", err)
		return models.Song{}, err
	}
	defer resp.Body.Close()

	c.log.Error("server response body:", resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.log.Error("received non-OK response status:", resp.Status)
		return models.Song{}, fmt.Errorf("received non-OK response status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error("failed to read response body:", err)
		return models.Song{}, err
	}
	c.log.Debug("server response body:", string(body))

	if len(body) == 0 {
		c.log.Error("response body is empty")
		return models.Song{}, fmt.Errorf("empty response body")
	}

	var songDetail models.Song
	err = json.Unmarshal(body, &songDetail)
	if err != nil {
		c.log.Error("failed to unmarshal response body:", err)
		return models.Song{}, err
	}

	c.log.Debug("url=", url)

	return songDetail, nil
}
