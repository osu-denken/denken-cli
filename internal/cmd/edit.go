package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/osu-denken/denken-cli/internal/editor"
)

// runEdit は「認証 → 現在の内容を取得 → エディタで編集 → 保存」の一連を行う。
// fetch/save が呼ばれる context はそれぞれ独立し、編集中は時間制限を受けない。
func (a *appContext) runEdit(ext, doneLabel string, fetch func(context.Context) (string, error), save func(context.Context, string) error) error {
	current, err := a.authFetch(fetch)
	if err != nil {
		return err
	}
	edited, ok, err := a.openEditor(current, ext)
	if err != nil || !ok {
		return err
	}
	ctx, cancel := newContext()
	defer cancel()
	if err := save(ctx, edited); err != nil {
		return err
	}
	fmt.Fprintln(a.out, "更新しました:", doneLabel)
	return nil
}

// authFetch は認証を保証したうえで fetch を実行する。
func (a *appContext) authFetch(fetch func(context.Context) (string, error)) (string, error) {
	ctx, cancel := newContext()
	defer cancel()
	if err := a.requireAuth(ctx); err != nil {
		return "", err
	}
	return fetch(ctx)
}

// openEditor は current をエディタで開き、編集後の内容と変更有無を返す。
// 変更が無ければ ok=false を返し、メッセージを表示する。
func (a *appContext) openEditor(current, ext string) (edited string, ok bool, err error) {
	edited, err = editor.Edit(current, ext)
	if errors.Is(err, editor.ErrUnchanged) {
		fmt.Fprintln(a.out, "変更がなかったため中止しました。")
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return edited, true, nil
}

// jsonField は生 JSON から文字列フィールドを取り出す (無ければ空文字)。
func jsonField(raw json.RawMessage, key string) string {
	m := map[string]any{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return ""
	}
	s, _ := m[key].(string)
	return s
}
