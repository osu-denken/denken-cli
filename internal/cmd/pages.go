package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

func newPagesCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "pages", Short: "ターミナル welcome.md と固定ページの編集 (要 PageEdit 権限)"}
	cmd.AddCommand(newTerminalCmd(app), newSitePagesCmd(app))
	return cmd
}

func newTerminalCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "terminal", Short: "トップページのターミナル welcome.md"}
	cmd.AddCommand(newTerminalGetCmd(app), newTerminalUpdateCmd(app), newTerminalEditCmd(app))
	return cmd
}

func newTerminalEditCmd(app *appContext) *cobra.Command {
	var page string
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "welcome.md を $EDITOR で開いて編集、保存する",
		RunE: func(cmd *cobra.Command, args []string) error {
			fetch := func(ctx context.Context) (string, error) {
				raw, err := app.client().TerminalGet(ctx, page)
				return jsonField(raw, "content"), err
			}
			save := func(ctx context.Context, content string) error {
				_, err := app.client().TerminalUpdate(ctx, content)
				return err
			}
			return app.runEdit(".md", page, fetch, save)
		},
	}
	cmd.Flags().StringVar(&page, "page", "welcome", "ページ名")
	return cmd
}

func newTerminalGetCmd(app *appContext) *cobra.Command {
	var page string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "welcome.md の内容を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(false, func(ctx context.Context) (rawJSON, error) {
				return app.client().TerminalGet(ctx, page)
			})
		},
	}
	cmd.Flags().StringVar(&page, "page", "welcome", "ページ名")
	return cmd
}

func newTerminalUpdateCmd(app *appContext) *cobra.Command {
	var file string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "welcome.md を更新し再ビルドを起動する",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().TerminalUpdate(ctx, string(content))
			})
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "内容ファイルのパス")
	cmd.MarkFlagRequired("file")
	return cmd
}
