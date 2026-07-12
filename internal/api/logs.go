package api

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// LogsQuery は操作ログ一覧の絞り込み条件。
type LogsQuery struct {
	Type   string
	Cursor string
	Limit  int
}

// LogsList は操作ログを新しい順に一覧する。LogView 権限が必要。
func (c *Client) LogsList(ctx context.Context, q LogsQuery) (json.RawMessage, error) {
	values := url.Values{}
	if q.Type != "" {
		values.Set("type", q.Type)
	}
	if q.Cursor != "" {
		values.Set("cursor", q.Cursor)
	}
	if q.Limit > 0 {
		values.Set("limit", strconv.Itoa(q.Limit))
	}
	var out json.RawMessage
	if err := c.GetJSON(ctx, "/logs/list", values, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Ping はサーバーの稼働確認を行う。pong という文字列を返す。
func (c *Client) Ping(ctx context.Context) (string, error) {
	return c.GetText(ctx, "/ping")
}
