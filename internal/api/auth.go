package api

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
)

// AuthResult は認証系エンドポイントが返すトークン情報。
// 2段階認証が有効な場合は MFARequired が true になりトークンは空になる。
type AuthResult struct {
	IDToken         string      `json:"idToken"`
	RefreshToken    string      `json:"refreshToken"`
	ExpiresIn       json.Number `json:"expiresIn"` // API は文字列・数値どちらでも返しうる
	LocalID         string      `json:"localId"`
	Email           string      `json:"email"`
	MFARequired     bool        `json:"mfaRequired"`
	MFAPendingToken string      `json:"mfaPendingToken"`
}

// ExpiresAt は expiresIn (秒) を絶対時刻に変換する。値が無ければゼロ値。
func (r *AuthResult) ExpiresAt() time.Time {
	secs, err := strconv.Atoi(r.ExpiresIn.String())
	if err != nil || secs == 0 {
		return time.Time{}
	}
	return time.Now().Add(time.Duration(secs) * time.Second)
}

// Exists は指定メールのユーザーが存在するかを確認する。
func (c *Client) Exists(ctx context.Context, email string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/exists", map[string]string{"email": email})
}

// Register は新規ユーザーを登録する。passphrase は招待コード。
func (c *Client) Register(ctx context.Context, email, password, passphrase string) (*AuthResult, error) {
	body := map[string]string{"email": email, "password": password, "passphrase": passphrase}
	return c.postAuth(ctx, "/user/register", body)
}

// Login はログインして ID トークンを取得する。
func (c *Client) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	body := map[string]string{"email": email, "password": password}
	return c.postAuth(ctx, "/user/login", body)
}

// LoginTotp は mfaPendingToken と6桁コードを検証してトークンを得る。
func (c *Client) LoginTotp(ctx context.Context, pendingToken, code string) (*AuthResult, error) {
	body := map[string]string{"mfaPendingToken": pendingToken, "code": code}
	return c.postAuth(ctx, "/user/loginTotp", body)
}

// Refresh はリフレッシュトークンから新しい ID トークンを取得する。
func (c *Client) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	body := map[string]string{"refreshToken": refreshToken}
	return c.postAuth(ctx, "/user/refresh", body)
}

// Info は認証済みユーザーの詳細情報を取得する。
func (c *Client) Info(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/info", nil)
}

// UpdateUser はユーザー情報を更新する。変更したい項目のみ body に含める。
func (c *Client) UpdateUser(ctx context.Context, body map[string]string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/update", body)
}

// ResetPassword はパスワードリセットメールを送信する。
func (c *Client) ResetPassword(ctx context.Context, email string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/resetPassword", map[string]string{"email": email})
}

// VerifyEmail はログイン中のユーザーへ確認メールを再送する。
func (c *Client) VerifyEmail(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/user/verifyEmail", nil)
}

func (c *Client) postAuth(ctx context.Context, path string, body any) (*AuthResult, error) {
	out := &AuthResult{}
	if err := c.PostJSON(ctx, path, body, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) postRaw(ctx context.Context, path string, body any) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.PostJSON(ctx, path, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}
