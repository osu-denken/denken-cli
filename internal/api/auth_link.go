package api

import (
	"context"
	"encoding/json"
)

// Google は Google の ID トークンを Firebase トークンに交換してログインする。
func (c *Client) Google(ctx context.Context, credential string) (*AuthResult, error) {
	return c.postAuth(ctx, "/user/google", map[string]string{"credential": credential})
}

// LinkGoogle はログイン中のアカウントに Google アカウントを連携する。
func (c *Client) LinkGoogle(ctx context.Context, credential string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/linkGoogle", map[string]string{"credential": credential})
}

// UnlinkGoogle は Google 連携を解除する。
func (c *Client) UnlinkGoogle(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/unlinkGoogle", nil)
}

// Providers はアカウントに紐づくログイン手段を返す。
func (c *Client) Providers(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/providers", nil)
}

// TotpSetup は TOTP のシークレットと QR を発行する (まだ有効化しない)。
func (c *Client) TotpSetup(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/totp/setup", nil)
}

// TotpEnable はコードを検証して2段階認証を有効化し、リカバリコードを返す。
func (c *Client) TotpEnable(ctx context.Context, code string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/totp/enable", map[string]string{"code": code})
}

// TotpDisable はコード (またはリカバリコード) を検証して2段階認証を解除する。
func (c *Client) TotpDisable(ctx context.Context, code string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/totp/disable", map[string]string{"code": code})
}
