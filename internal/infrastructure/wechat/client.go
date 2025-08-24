package wechat

import (
	"context"
	"fmt"
)

// Client simulates interactions with the WeChat API.
type Client struct {
	// In a real implementation, this would hold configuration
	// like AppID and AppSecret.
}

// NewClient creates a new mock WeChat client.
func NewClient() *Client {
	return &Client{}
}

// CodeToOpenID simulates exchanging a temporary code for a user's openid.
// In a real application, this would make an HTTP request to the WeChat API.
func (c *Client) CodeToOpenID(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("wechat code cannot be empty")
	}
	// For simulation purposes, we just prepend a prefix to the code.
	// A real openid is much longer and more complex.
	return "mock_openid_for_" + code, nil
}
