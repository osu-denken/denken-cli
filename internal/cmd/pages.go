package cmd

import (
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
	cmd.AddCommand(newTerminalGetCmd(app), newTerminalUpdateCmd(app))
	return cmd
}

func newTerminalGetCmd(app *appContext) *cobra.Command {
	var page string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "welcome.md の内容を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			raw, err := app.client().TerminalGet(ctx, page)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
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
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().TerminalUpdate(ctx, string(content))
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "内容ファイルのパス")
	cmd.MarkFlagRequired("file")
	return cmd
}
