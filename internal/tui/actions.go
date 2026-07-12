package tui

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/osu-denken/denken-cli/internal/api"
)

// action はメニュー項目1つ分の実行内容。
type action struct {
	title string
	auth  bool
	run   func(context.Context, *api.Client) (json.RawMessage, error)
}

// pingResult は非 JSON の ping を JSON 風に包む簡易ヘルパー。
func pingAction(ctx context.Context, c *api.Client) (json.RawMessage, error) {
	res, err := c.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(`"` + res + `"`), nil
}

// actions は TUI に表示する読み取り系アクション一覧を返す。
func actions() []action {
	return []action{
		{"サーバー稼働確認 (ping)", false, pingAction},
		{"自分の情報 (whoami)", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.Info(ctx) }},
		{"ブログ記事一覧", false, func(ctx context.Context, c *api.Client) (json.RawMessage, error) {
			items, err := c.BlogList(ctx)
			if err != nil {
				return nil, err
			}
			return json.Marshal(items)
		}},
		{"ポータル情報", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.Portal(ctx) }},
		{"部員一覧", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.MembersList(ctx, "") }},
		{"非公開記事一覧", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.PrivatePostList(ctx) }},
		{"画像一覧", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.ImageList(ctx) }},
		{"SwitchBot デバイス一覧", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.SwitchBotList(ctx) }},
		{"Discord 招待コード", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.DiscordInvite(ctx) }},
		{"操作ログ", true, func(ctx context.Context, c *api.Client) (json.RawMessage, error) { return c.LogsList(ctx, api.LogsQuery{}) }},
	}
}

// prettyJSON は生 JSON を整形する。
func prettyJSON(raw json.RawMessage) string {
	buf := &bytes.Buffer{}
	if err := json.Indent(buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}
