package api

import (
	"context"
	"encoding/json"
)

// SwitchBotValidate はトークンが有効かどうかを確認する。SwitchBotControl 権限が必要。
func (c *Client) SwitchBotValidate(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/switchbot/validate", nil)
}

// SwitchBotList はデバイス一覧を返す。SwitchBotControl 権限が必要。
func (c *Client) SwitchBotList(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/switchbot/list", nil)
}

// SwitchBotLock は施錠する。SwitchBotControl 権限が必要。
func (c *Client) SwitchBotLock(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/switchbot/lock", nil)
}

// SwitchBotUnlock は解錠する。SwitchBotControl 権限が必要。
func (c *Client) SwitchBotUnlock(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/switchbot/unlock", nil)
}
