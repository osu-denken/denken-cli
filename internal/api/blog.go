package api

import (
	"context"
	"net/url"
)

// BlogItem はブログ記事一覧の1要素。
type BlogItem struct {
	Name string         `json:"name"`
	Sha  string         `json:"sha"`
	Size int            `json:"size"`
	Meta map[string]any `json:"meta"`
}

// BlogEntry は本文を含むブログ記事。
type BlogEntry struct {
	BlogItem
	Content string `json:"content"`
}

// BlogList は記事の一覧を取得する。
func (c *Client) BlogList(ctx context.Context) ([]BlogItem, error) {
	var out []BlogItem
	if err := c.GetJSON(ctx, "/v2/blog/list", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BlogGet は指定スラッグの本文とメタデータを取得する。
func (c *Client) BlogGet(ctx context.Context, slug string) (*BlogEntry, error) {
	q := url.Values{"page": {slug}}
	out := &BlogEntry{}
	if err := c.GetJSON(ctx, "/v2/blog/get", q, out); err != nil {
		return nil, err
	}
	return out, nil
}

// BlogUpdate は記事を新規作成または更新する。要認証。
func (c *Client) BlogUpdate(ctx context.Context, slug string, meta map[string]any, content string) error {
	body := map[string]any{"page": slug, "meta": meta, "content": content}
	return c.PostJSON(ctx, "/v2/blog/update", body, nil)
}
