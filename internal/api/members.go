package api

import (
	"context"
	"encoding/json"
	"net/url"
)

// MembersList は部員一覧を取得する。status は任意の絞り込み。
func (c *Client) MembersList(ctx context.Context, status string) (json.RawMessage, error) {
	q := url.Values{}
	if status != "" {
		q.Set("status", status)
	}
	return c.postRawQuery(ctx, "/members/list", q, nil)
}

// MembersDetail は部員一人の詳細を取得する。
func (c *Client) MembersDetail(ctx context.Context, id int) (json.RawMessage, error) {
	return c.postRaw(ctx, "/members/detail", map[string]int{"id": id})
}

// MembersApprove は仮部員を承認して在籍にする。MemberApprove 権限が必要。
func (c *Client) MembersApprove(ctx context.Context, id int) (json.RawMessage, error) {
	return c.postRaw(ctx, "/members/approve", map[string]int{"id": id})
}

// MembersReject は仮部員の登録を却下する。MemberApprove 権限が必要。
func (c *Client) MembersReject(ctx context.Context, id int) (json.RawMessage, error) {
	return c.postRaw(ctx, "/members/reject", map[string]int{"id": id})
}

// MembersUpdate は部員情報を更新する。id 以外は変更したい項目のみ含める。
func (c *Client) MembersUpdate(ctx context.Context, body map[string]any) (json.RawMessage, error) {
	return c.postRaw(ctx, "/members/update", body)
}

func (c *Client) postRawQuery(ctx context.Context, path string, query url.Values, body any) (json.RawMessage, error) {
	full := path
	if len(query) > 0 {
		full += "?" + query.Encode()
	}
	return c.postRaw(ctx, full, body)
}
