package api

import (
	"context"
	"encoding/json"
)

// GithubTokenGet は保存済み GitHub PAT の状態を確認する。BlogEdit 権限が必要。
func (c *Client) GithubTokenGet(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.GetJSON(ctx, "/github/token", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GithubTokenSave は GitHub PAT を暗号化して保存する。BlogEdit 権限が必要。
func (c *Client) GithubTokenSave(ctx context.Context, token string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/github/token", map[string]string{"token": token})
}

// GithubTokenDelete は保存済み GitHub PAT を削除する。BlogEdit 権限が必要。
func (c *Client) GithubTokenDelete(ctx context.Context) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.SendRaw(ctx, "DELETE", "/github/token", nil, "", &out); err != nil {
		return nil, err
	}
	return out, nil
}
