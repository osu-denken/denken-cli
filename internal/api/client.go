package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client は Web API への HTTP クライアント。
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// New はクライアントを生成する。token は空でもよい (未認証エンドポイント用)。
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// WithToken は ID トークンを差し替えたクライアントを返す。
func (c *Client) WithToken(token string) *Client {
	return &Client{baseURL: c.baseURL, token: token, http: c.http}
}

// GetJSON は GET リクエストを送り、レスポンスを out に格納する。
func (c *Client) GetJSON(ctx context.Context, path string, query url.Values, out any) error {
	full := path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	req, err := c.newRequest(ctx, http.MethodGet, full, nil, "")
	if err != nil {
		return err
	}
	return c.do(req, out)
}

// PostJSON は JSON ボディで POST し、レスポンスを out に格納する。
// body が nil の場合は空ボディで送る。
func (c *Client) PostJSON(ctx context.Context, path string, body, out any) error {
	var reader io.Reader
	contentType := ""
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
		contentType = "application/json"
	}
	req, err := c.newRequest(ctx, http.MethodPost, path, reader, contentType)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

// GetText は GET リクエストを送り、レスポンスボディを文字列として返す (非 JSON 用)。
func (c *Client) GetText(ctx context.Context, path string) (string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil, "")
	if err != nil {
		return "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", &APIError{Status: resp.StatusCode, Body: strings.TrimSpace(string(data))}
	}
	return strings.TrimSpace(string(data)), nil
}

// SendRaw は任意メソッド・任意ボディを送る低レベル API (multipart 等)。
func (c *Client) SendRaw(ctx context.Context, method, path string, body io.Reader, contentType string, out any) error {
	req, err := c.newRequest(ctx, method, path, body, contentType)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return &APIError{Status: resp.StatusCode, Body: strings.TrimSpace(string(data))}
	}
	if out == nil || len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	return json.Unmarshal(data, out)
}
