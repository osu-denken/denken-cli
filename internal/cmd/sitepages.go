package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newSitePagesCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "site", Short: "公式サイト本体の固定ページ (content/ 配下)"}
	cmd.AddCommand(
		authRawCmd(app, "list", "編集できるファイルの一覧を返す", func(c cmdCtx) (any, error) {
			return app.client().SitePagesList(c.ctx)
		}),
		newSitePagesGetCmd(app), newSitePagesUpdateCmd(app),
	)
	return cmd
}

func newSitePagesGetCmd(app *appContext) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "固定ページファイルの中身を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().SitePagesGet(ctx, path)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "ファイルパス")
	cmd.MarkFlagRequired("path")
	return cmd
}

func newSitePagesUpdateCmd(app *appContext) *cobra.Command {
	var path, file string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "固定ページファイルを更新し再ビルドを起動する",
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
			raw, err := app.client().SitePagesUpdate(ctx, path, string(content))
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "ファイルパス")
	cmd.Flags().StringVar(&file, "file", "", "内容ファイルのパス")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("file")
	return cmd
}
