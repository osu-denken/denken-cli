package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
)

// cmdCtx はコマンド実行時のコンテキストをまとめる。
type cmdCtx struct {
	ctx context.Context
}

// rawFunc は生 JSON を返す API 呼び出し。
type rawFunc func(cmdCtx) (any, error)

// authRawCmd は「認証必須・引数なし・JSON を表示する」コマンドを生成する。
func authRawCmd(app *appContext, use, short string, fn rawFunc) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			return runAndPrint(app, fn, ctx)
		},
	}
}

// authSlugCmd は文字列引数を1つ取り認証必須で JSON を表示するコマンドを作る。
func authSlugCmd(app *appContext, use, short string, fn func(cmdCtx, string) (any, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			res, err := fn(cmdCtx{ctx}, args[0])
			if err != nil {
				return err
			}
			return app.printJSON(res.(json.RawMessage))
		},
	}
}

func runAndPrint(app *appContext, fn rawFunc, ctx context.Context) error {
	res, err := fn(cmdCtx{ctx})
	if err != nil {
		return err
	}
	raw, ok := res.(json.RawMessage)
	if !ok {
		return nil
	}
	return app.printJSON(raw)
}
