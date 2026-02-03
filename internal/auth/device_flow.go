package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/CromulentConsulting/ygm-cli/internal/api"
)

// DeviceFlow handles the OAuth device flow authentication
type DeviceFlow struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewDeviceFlow creates a new device flow handler
func NewDeviceFlow(baseURL string) *DeviceFlow {
	return &DeviceFlow{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RequestDeviceCode requests a new device code to start the flow
func (d *DeviceFlow) RequestDeviceCode() (*api.DeviceCodeResponse, error) {
	resp, err := d.HTTPClient.Post(d.BaseURL+"/oauth/device/codes", "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to request device code (status %d): %s", resp.StatusCode, string(body))
	}

	var result api.DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse device code response: %w", err)
	}

	return &result, nil
}

// PollForToken polls the token endpoint until authorized or expired
func (d *DeviceFlow) PollForToken(deviceCode string, interval int, tokenName string) (*api.TokenResponse, error) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		token, err := d.checkToken(deviceCode, tokenName)
		if err != nil {
			// Check if it's a pending error (continue polling)
			if isPendingError(err) {
				continue
			}
			return nil, err
		}

		if token != nil {
			return token, nil
		}
	}
}

func (d *DeviceFlow) checkToken(deviceCode, tokenName string) (*api.TokenResponse, error) {
	data := url.Values{}
	data.Set("device_code", deviceCode)
	if tokenName != "" {
		data.Set("token_name", tokenName)
	}

	req, err := http.NewRequest("POST", d.BaseURL+"/oauth/device/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to check token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		var errResp api.TokenErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			switch errResp.Error {
			case "authorization_pending":
				return nil, &PendingError{}
			case "access_denied":
				return nil, fmt.Errorf("authorization denied by user")
			case "expired_token":
				return nil, fmt.Errorf("device code expired")
			default:
				return nil, fmt.Errorf("authorization error: %s", errResp.Error)
			}
		}
		return nil, fmt.Errorf("authorization failed: %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed (status %d): %s", resp.StatusCode, string(body))
	}

	var token api.TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &token, nil
}

// PendingError indicates authorization is still pending
type PendingError struct{}

func (e *PendingError) Error() string {
	return "authorization pending"
}

func isPendingError(err error) bool {
	_, ok := err.(*PendingError)
	return ok
}

// OpenBrowser opens the default browser to the given URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}
