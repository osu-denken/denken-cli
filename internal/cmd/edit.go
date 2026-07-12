package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/osu-denken/denken-cli/internal/editor"
)

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
