package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/soheilrt/checkrr/config"
	"io"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	queuePath           = "/api/v3/queue"
	queueDeleteBulkPath = "/api/v3/queue/bulk"
)

type ClientRR struct {
	key     string
	host    string
	options config.Options
}

func NewClientRR(host, apiKey string, options config.Options) *ClientRR {
	return &ClientRR{
		key:     apiKey,
		host:    host,
		options: options,
	}
}

func (c *ClientRR) FetchDownloads() ([]Download, error) {
	var downloads []Download

	for page := 1; ; page++ {
		log.Debugf("Fetching download page %d...\n", page)
		pageDownloads, totalRecords, err := c.fetchDownloadPage(page)
		if err != nil {
			return nil, err
		}
		downloads = append(downloads, pageDownloads...)

		if len(downloads) >= totalRecords {
			break
		}
	}

	return downloads, nil
}

func (c *ClientRR) fetchDownloadPage(page int) ([]Download, int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", c.host+queuePath, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-Api-Key", c.key)

	// Get query values, add "page" param, and set the Host's RawQuery
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = q.Encode() // Important: this actually sets the query string

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("Unexpected Status Code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, 0, err
	}

	return response.Records, response.TotalRecords, nil
}

func (c *ClientRR) DeleteFromQueue(ids []int) error {
	// Create a new HTTP client
	client := &http.Client{}

	// Prepare the request
	req, err := http.NewRequest("DELETE", c.host+queueDeleteBulkPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Api-Key", c.key)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("removeFromClient", strconv.FormatBool(!c.options.KeepInClient))
	q.Add("blocklist", strconv.FormatBool(c.options.BlockList))
	req.URL.RawQuery = q.Encode() // Important: this actually sets the query string

	// Convert IDs to string and marshal the body
	body, err := json.Marshal(struct {
		Ids []int `json:"ids"`
	}{Ids: ids})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Attach the body to the request
	req.Body = io.NopCloser(bytes.NewReader(body))

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for unsuccessful HTTP status codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
