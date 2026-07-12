package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/osu-denken/denken-cli/internal/config"
)

// errNotLoggedIn は認証が必要なコマンドで未ログインのときに返す。
var errNotLoggedIn = errors.New("ログインしていません。`denken-cli login` を実行してください")

// appContext は各コマンドが共有する実行コンテキスト。
type appContext struct {
	cfg   *config.Config
	out   io.Writer
	token string // --token による上書き (空なら設定を使う)
}

// client は現在の設定に基づく API クライアントを返す。
func (a *appContext) client() *api.Client {
	token := a.token
	if token == "" {
		token = a.cfg.Session.IDToken
	}
	return api.New(a.cfg.BaseURL, token)
}

// requireAuth は認証済みであることを保証し、必要ならトークンを自動更新する。
func (a *appContext) requireAuth(ctx context.Context) error {
	if a.token != "" {
		return nil
	}
	if !a.cfg.LoggedIn() {
		return errNotLoggedIn
	}
	if a.cfg.Session.Expired() {
		return a.refresh(ctx)
	}
	return nil
}

// refresh はリフレッシュトークンで ID トークンを更新し、設定を保存する。
func (a *appContext) refresh(ctx context.Context) error {
	rt := a.cfg.Session.RefreshToken
	if rt == "" {
		return errNotLoggedIn
	}
	res, err := api.New(a.cfg.BaseURL, "").Refresh(ctx, rt)
	if err != nil {
		return err
	}
	a.applyAuth(res)
	return config.Save(a.cfg)
}

// applyAuth は認証結果をセッションに反映する (トークンが空なら変更しない)。
func (a *appContext) applyAuth(res *api.AuthResult) {
	if res.IDToken != "" {
		a.cfg.Session.IDToken = res.IDToken
	}
	if res.RefreshToken != "" {
		a.cfg.Session.RefreshToken = res.RefreshToken
	}
	if res.Email != "" {
		a.cfg.Session.Email = res.Email
	}
	if at := res.ExpiresAt(); !at.IsZero() {
		a.cfg.Session.ExpiresAt = at
	}
}

// printJSON は生の JSON を整形して出力する。
func (a *appContext) printJSON(raw json.RawMessage) error {
	if len(bytes.TrimSpace(raw)) == 0 {
		fmt.Fprintln(a.out, "(空のレスポンス)")
		return nil
	}
	buf := &bytes.Buffer{}
	if err := json.Indent(buf, raw, "", "  "); err != nil {
		fmt.Fprintln(a.out, string(raw))
		return nil
	}
	fmt.Fprintln(a.out, buf.String())
	return nil
}

// newContext はタイムアウト付きのコンテキストを返す。
func newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 60*time.Second)
}
