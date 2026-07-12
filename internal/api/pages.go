package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// TerminalGet はトップページのターミナルに表示される welcome.md を取得する。
func (c *Client) TerminalGet(ctx context.Context, page string) (json.RawMessage, error) {
	if page == "" {
		page = "welcome"
	}
	var out json.RawMessage
	if err := c.GetJSON(ctx, "/terminal/get", url.Values{"page": {page}}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TerminalUpdate は welcome.md を更新し、公式サイトの再ビルドを起動する。PageEdit 権限が必要。
func (c *Client) TerminalUpdate(ctx context.Context, content string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/terminal/update", map[string]string{"content": content})
}

// SitePagesList は編集できるファイルの一覧を返す。PageEdit 権限が必要。
func (c *Client) SitePagesList(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/site-pages/list", nil)
}

// SitePagesGet は固定ページファイルの中身を取得する。PageEdit 権限が必要。
func (c *Client) SitePagesGet(ctx context.Context, path string) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.GetJSON(ctx, "/site-pages/get", url.Values{"path": {path}}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SitePagesUpdate は固定ページファイルを更新し、再ビルドを起動する。PageEdit 権限が必要。
func (c *Client) SitePagesUpdate(ctx context.Context, path, content string) (json.RawMessage, error) {
	body := map[string]string{"path": path, "content": content}
	return c.postRaw(ctx, "/site-pages/update", body)
}
