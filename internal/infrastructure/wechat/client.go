package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/moriverse/45-server/internal/infrastructure/config"
)

// Client interacts with the WeChat API.
type Client struct {
	AppID            string
	AppSecret        string
	Host             string
	Code2SessionPath string
	HTTPClient       *http.Client
}

// NewClient creates a new WeChat client.
func NewClient(cfg config.WechatConfig) *Client {
	return &Client{
		AppID:            cfg.AppID,
		AppSecret:        cfg.AppSecret,
		Host:             cfg.Host,
		Code2SessionPath: cfg.Code2SessionPath,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// wechatErrorResponse defines the error structure returned by the WeChat API.
type wechatErrorResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Code2SessionResponse defines the successful response structure for the code2session call.
type Code2SessionResponse struct {
	wechatErrorResponse
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

// CodeToOpenID exchanges a temporary code for a user's openid by calling the WeChat API.
func (c *Client) CodeToOpenID(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", fmt.Errorf("wechat code cannot be empty")
	}

	baseURL, err := url.Parse(c.Host)
	if err != nil {
		return "", fmt.Errorf("failed to parse wechat host: %w", err)
	}
	baseURL.Path = c.Code2SessionPath

	params := url.Values{}
	params.Add("appid", c.AppID)
	params.Add("secret", c.AppSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")
	baseURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create wechat request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call wechat api: %w", err)
	}
	defer resp.Body.Close()

	var result Code2SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode wechat response: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("wechat api error: code=%d, msg=%s", result.ErrCode, result.ErrMsg)
	}

	if result.OpenID == "" {
		return "", fmt.Errorf("openid is empty in wechat response")
	}

	return result.OpenID, nil
}
