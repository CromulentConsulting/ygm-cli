package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is an HTTP client for the YGM API
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.HTTPClient.Do(req)
}

// GetBrand fetches the active brand DNA
func (c *Client) GetBrand() (*BrandDNA, error) {
	resp, err := c.doRequest("GET", "/api/v1/brand", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No brand DNA
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var brand BrandDNA
	if err := json.NewDecoder(resp.Body).Decode(&brand); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &brand, nil
}

// GetBrandVersions fetches all brand DNA versions
func (c *Client) GetBrandVersions() ([]BrandDNA, error) {
	resp, err := c.doRequest("GET", "/api/v1/brand/versions", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result BrandVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Versions, nil
}

// GetTasks fetches tasks with optional filters
func (c *Client) GetTasks(status, platform string) ([]Task, error) {
	path := "/api/v1/tasks"
	query := ""
	if status != "" {
		query += "status=" + status
	}
	if platform != "" {
		if query != "" {
			query += "&"
		}
		query += "platform=" + platform
	}
	if query != "" {
		path += "?" + query
	}

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result TasksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Tasks, nil
}

// GetTask fetches a single task by ID
func (c *Client) GetTask(id int) (*Task, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title             string  `json:"title"`
	Description       string  `json:"description,omitempty"`
	Platform          string  `json:"platform,omitempty"`
	AssetType         string  `json:"asset_type,omitempty"`
	SuggestedPostDate *string `json:"suggested_post_date,omitempty"`
}

// CreateTask creates a new marketing task
func (c *Client) CreateTask(req CreateTaskRequest) (*Task, error) {
	body := map[string]interface{}{
		"task": req,
	}

	resp, err := c.doRequest("POST", "/api/v1/tasks", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.parseError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(id int, req UpdateTaskRequest) (*Task, error) {
	body := map[string]interface{}{
		"task": req,
	}

	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/v1/tasks/%d", id), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// DiscardTask soft-deletes a task
func (c *Client) DiscardTask(id int) (*DiscardResponse, error) {
	resp, err := c.doRequest("DELETE", fmt.Sprintf("/api/v1/tasks/%d", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var result DiscardResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetContext fetches the full context dump
func (c *Client) GetContext() (*ContextResponse, error) {
	resp, err := c.doRequest("GET", "/api/v1/context", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var ctx ContextResponse
	if err := json.NewDecoder(resp.Body).Decode(&ctx); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ctx, nil
}

func (c *Client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
		return fmt.Errorf("API error: %s", errResp.Error)
	}

	return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
}
