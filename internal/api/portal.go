package api

import (
	"context"
	"encoding/json"
)

// Portal はポータルサイト用の情報をまとめて取得する。
func (c *Client) Portal(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/portal", nil)
}

// GithubInvite は指定メールを GitHub Organization に招待する。MemberManage 権限が必要。
func (c *Client) GithubInvite(ctx context.Context, email string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/github/invite", map[string]string{"email": email})
}

// GithubJoin は部員自身が Organization への招待を受け取る。username は連携済みなら省略可。
func (c *Client) GithubJoin(ctx context.Context, username string) (json.RawMessage, error) {
	body := map[string]string{}
	if username != "" {
		body["username"] = username
	}
	return c.postRaw(ctx, "/github/join", body)
}

// GithubUsername は連携済みの GitHub ログイン名を返す (未連携なら null)。
func (c *Client) GithubUsername(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/github/username", nil)
}

// GithubOAuthStart は認可 URL を返す。BlogEdit 権限が必要。
func (c *Client) GithubOAuthStart(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/github/oauth/start", nil)
}

// DiscordInvite は Discord サーバーへの招待コードを返す。
func (c *Client) DiscordInvite(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/discord/invite", nil)
}
