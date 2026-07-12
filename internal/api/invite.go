package api

import (
	"context"
	"encoding/json"
)

// InviteValidate は招待コードが有効かどうかを検証する。
func (c *Client) InviteValidate(ctx context.Context, code string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/invite/validate", map[string]string{"code": code})
}

// InviteCreate は新しい招待コードを生成する。InviteCodeCreate 権限が必要。
func (c *Client) InviteCreate(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/invite/create", nil)
}

// InviteDelete は指定した招待コードを無効化する。
func (c *Client) InviteDelete(ctx context.Context, code string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/invite/delete", map[string]string{"code": code})
}
