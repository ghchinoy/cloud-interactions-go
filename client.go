// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interactions

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StreamCreate initiates a new interaction turn and streams back Server-Sent Events.
func (c *Client) StreamCreate(ctx context.Context, req *InteractionRequest, onEvent func(event, data string) error) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}
	if c.APIKey != "" {
		httpReq.Header.Set("x-goog-api-key", c.APIKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if sherlogURL := resp.Header.Get("X-Goog-Sherlog-Link"); sherlogURL != "" {
		_ = onEvent("sherlog", fmt.Sprintf(`{"url": "%s"}`, sherlogURL))
	} else if sherlogURL := resp.Header.Get("x-goog-sherlog-link"); sherlogURL != "" {
		_ = onEvent("sherlog", fmt.Sprintf(`{"url": "%s"}`, sherlogURL))
	}

	scanner := bufio.NewScanner(resp.Body)
	// Increase scanner buffer to handle large base64 video/audio payloads
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 128*1024*1024)
	var currentEvent string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			dataStr := strings.TrimPrefix(line, "data: ")
			if err := onEvent(currentEvent, dataStr); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

// Client is a Go wrapper for the Gemini/Vertex Interactions REST API.
type Client struct {
	APIKey      string
	BearerToken string
	BaseURL     string
	HTTPClient  *http.Client
}

// NewClient creates a new Interactions API client.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Minute},
	}
}

// WithAPIKey sets an API key.
func (c *Client) WithAPIKey(key string) *Client {
	c.APIKey = key
	return c
}

// WithBearerToken sets a Bearer token (e.g. from gcloud auth).
func (c *Client) WithBearerToken(token string) *Client {
	c.BearerToken = token
	return c
}

// formatURL appends the API key to the URL if needed.
func (c *Client) formatURL(endpoint string) string {
	if c.APIKey != "" {
		if strings.Contains(endpoint, "?") {
			return fmt.Sprintf("%s&key=%s", endpoint, c.APIKey)
		}
		return fmt.Sprintf("%s?key=%s", endpoint, c.APIKey)
	}
	return endpoint
}

// Create initiates a new interaction turn.
func (c *Client) Create(ctx context.Context, req *InteractionRequest) (*InteractionResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}
	if c.APIKey != "" {
		httpReq.Header.Set("x-goog-api-key", c.APIKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	return c.do(httpReq)
}

// Get retrieves the status and result of an existing interaction.
func (c *Client) Get(ctx context.Context, id string) (*InteractionResponse, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, id)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}

	return c.do(httpReq)
}

// Delete removes a stored interaction.
func (c *Client) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/%s", c.BaseURL, id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", c.formatURL(url), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}

	_, err = c.do(httpReq)
	return err
}

// WaitForCompletion polls an interaction until its status is a terminal state.
func (c *Client) WaitForCompletion(ctx context.Context, id string, interval time.Duration) (*InteractionResponse, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			resp, err := c.Get(ctx, id)
			if err != nil {
				return nil, err
			}

			status := strings.ToLower(resp.Status)
			// Keep polling if it's in a known active state
			if status == "working" || status == "in_progress" || status == "pending" {
				continue
			}

			// Return if it's likely a terminal state (completed, failed, etc.)
			return resp, nil
		}
	}
}

func (c *Client) do(req *http.Request) (*InteractionResponse, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if len(body) == 0 {
		return nil, nil
	}

	var interactionResp InteractionResponse
	if err := json.Unmarshal(body, &interactionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &interactionResp, nil
}
