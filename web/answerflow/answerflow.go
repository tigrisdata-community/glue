package answerflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tigrisdata-community/glue/web/useragent"
)

type Client struct {
	APIKey  string
	BaseURL string
}

type CreateSolutionRequest struct {
	SolutionID string `json:"solutionId"`
}

type CreateSolutionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func New(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://www.answeroverflow.com",
	}
}

// CreateSolution marks a message as the solution for a thread.
func (c *Client) CreateSolution(messageID string, solutionID string) (*CreateSolutionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/messages/%s", c.BaseURL, messageID)

	body := CreateSolutionRequest{
		SolutionID: solutionID,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	var resp CreateSolutionResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}
