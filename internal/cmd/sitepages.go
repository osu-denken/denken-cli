package cmd

import (
	"context"
	"os"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/osu-denken/denken-cli/internal/editor"
	"github.com/spf13/cobra"
)

func newSitePagesCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "site", Short: "公式サイト本体の固定ページ (content/ 配下)"}
	cmd.AddCommand(
		authRawCmd(app, "list", "編集できるファイルの一覧を返す", (*api.Client).SitePagesList),
		newSitePagesGetCmd(app), newSitePagesUpdateCmd(app), newSitePagesEditCmd(app),
	)
	return cmd
}

func newSitePagesEditCmd(app *appContext) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "edit --path <p>",
		Short: "固定ページを $EDITOR で開いて編集、保存する",
		RunE: func(cmd *cobra.Command, args []string) error {
			fetch := func(ctx context.Context) (string, error) {
				raw, err := app.client().SitePagesGet(ctx, path)
				return jsonField(raw, "content"), err
			}
			save := func(ctx context.Context, content string) error {
				_, err := app.client().SitePagesUpdate(ctx, path, content)
				return err
			}
			return app.runEdit(editor.Ext(path), path, fetch, save)
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "ファイルパス")
	cmd.MarkFlagRequired("path")
	return cmd
}

func newSitePagesGetCmd(app *appContext) *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "固定ページファイルの中身を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().SitePagesGet(ctx, path)
			})
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
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().SitePagesUpdate(ctx, path, string(content))
			})
		},
	}
	cmd.Flags().StringVar(&path, "path", "", "ファイルパス")
	cmd.Flags().StringVar(&file, "file", "", "内容ファイルのパス")
	cmd.MarkFlagRequired("path")
	cmd.MarkFlagRequired("file")
	return cmd
}
