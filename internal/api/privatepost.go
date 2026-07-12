package api

import (
	"context"
	"encoding/json"
)

// PrivatePostList は非公開記事を新しい順に一覧する (本文なし)。PrivatePostView 権限が必要。
func (c *Client) PrivatePostList(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/private-posts/list", nil)
}

// PrivatePostGet は非公開記事の本文を取得する。PrivatePostView 権限が必要。
func (c *Client) PrivatePostGet(ctx context.Context, slug string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/private-posts/get", map[string]string{"slug": slug})
}

// PrivatePostUpdate は非公開記事を新規作成または上書きする。PrivatePostEdit 権限が必要。
func (c *Client) PrivatePostUpdate(ctx context.Context, slug, title, content string) (json.RawMessage, error) {
	body := map[string]string{"slug": slug, "title": title, "content": content}
	return c.postRaw(ctx, "/private-posts/update", body)
}

// PrivatePostDelete は非公開記事を削除する。PrivatePostEdit 権限が必要。
func (c *Client) PrivatePostDelete(ctx context.Context, slug string) (json.RawMessage, error) {
	return c.postRaw(ctx, "/private-posts/delete", map[string]string{"slug": slug})
}
