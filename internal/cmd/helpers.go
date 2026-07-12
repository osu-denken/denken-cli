package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
)

// rawJSON は API が返す生 JSON。
type rawJSON = json.RawMessage

// apiCall は生 JSON を返す API 呼び出し。
type apiCall func(context.Context) (rawJSON, error)

// runRaw は「(auth なら認証) → API 呼び出し → JSON 表示」の定型処理をまとめる。
func (a *appContext) runRaw(auth bool, call apiCall) error {
	ctx, cancel := newContext()
	defer cancel()
	if auth {
		if err := a.requireAuth(ctx); err != nil {
			return err
		}
	}
	raw, err := call(ctx)
	if err != nil {
		return err
	}
	return a.printJSON(raw)
}

// authRawCmd は「認証必須・引数なし・JSON を表示する」コマンドを生成する。
func authRawCmd(app *appContext, use, short string, call apiCall) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(true, call)
		},
	}
}

// argCmd は文字列引数を1つ取り JSON を表示するコマンドを作る。auth で認証要否を指定する。
func argCmd(app *appContext, use, short string, auth bool, call func(context.Context, string) (rawJSON, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(auth, func(ctx context.Context) (rawJSON, error) {
				return call(ctx, args[0])
			})
		},
	}
}
