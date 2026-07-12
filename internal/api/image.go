package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// ImageList はアップロード済み画像を一覧する。BlogEdit 権限が必要。
func (c *Client) ImageList(ctx context.Context) (json.RawMessage, error) {
	return c.postRaw(ctx, "/image/list", nil)
}

// ImageDelete は画像を削除する。ImageDelete 権限が必要。
func (c *Client) ImageDelete(ctx context.Context, filename, sha string) (json.RawMessage, error) {
	body := map[string]string{"filename": filename, "sha": sha}
	return c.postRaw(ctx, "/image/delete", body)
}

// ImageUpload は画像をアップロードする。name は空なら UUID になる。ImageUpload 権限が必要。
func (c *Client) ImageUpload(ctx context.Context, filePath, name string) (json.RawMessage, error) {
	body, contentType, err := buildImageForm(filePath, name)
	if err != nil {
		return nil, err
	}
	var out json.RawMessage
	if err := c.SendRaw(ctx, "POST", "/image/upload", body, contentType, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func buildImageForm(filePath, name string) (io.Reader, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", err
	}
	if name != "" {
		if err := writer.WriteField("name", name); err != nil {
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return buf, writer.FormDataContentType(), nil
}
