package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type APIClient struct {
    baseURL string
    token   string
    client  *http.Client
}

func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *APIClient) SendMetrics(payload map[string]interface{}) error {

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    fmt.Println("Payload JSON:", string(jsonData))

    req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "GoMetricsCLI")

    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    return nil
}