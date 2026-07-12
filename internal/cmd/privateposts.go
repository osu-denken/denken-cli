package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

func newPrivatePostsCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "private-posts", Short: "非公開記事の操作", Aliases: []string{"pp"}}
	cmd.AddCommand(
		newPPListCmd(app), newPPGetCmd(app), newPPUpdateCmd(app), newPPDeleteCmd(app), newPPEditCmd(app),
	)
	return cmd
}

func newPPListCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "list", "非公開記事を一覧する (要 PrivatePostView 権限)", app.client().PrivatePostList)
}

func newPPGetCmd(app *appContext) *cobra.Command {
	return argCmd(app, "get <slug>", "非公開記事の本文を取得する", true, app.client().PrivatePostGet)
}

func newPPDeleteCmd(app *appContext) *cobra.Command {
	return argCmd(app, "delete <slug>", "非公開記事を削除する (要 PrivatePostEdit 権限)", true, app.client().PrivatePostDelete)
}

func newPPEditCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <slug>",
		Short: "非公開記事の本文を $EDITOR で開いて編集、保存する (要 PrivatePostEdit 権限)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPPEdit(app, args[0])
		},
	}
}

func runPPEdit(app *appContext, slug string) error {
	var title string
	fetch := func(ctx context.Context) (string, error) {
		raw, err := app.client().PrivatePostGet(ctx, slug)
		if err != nil {
			return "", err
		}
		title = jsonField(raw, "title")
		return jsonField(raw, "content"), nil
	}
	save := func(ctx context.Context, content string) error {
		_, err := app.client().PrivatePostUpdate(ctx, slug, title, content)
		return err
	}
	return app.runEdit(".md", slug, fetch, save)
}

func newPPUpdateCmd(app *appContext) *cobra.Command {
	var slug, title, file string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "非公開記事を新規作成または上書きする (要 PrivatePostEdit 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().PrivatePostUpdate(ctx, slug, title, string(content))
			})
		},
	}
	cmd.Flags().StringVar(&slug, "slug", "", "スラッグ (英小文字・数字・ハイフン)")
	cmd.Flags().StringVar(&title, "title", "", "タイトル")
	cmd.Flags().StringVar(&file, "file", "", "本文ファイルのパス")
	cmd.MarkFlagRequired("slug")
	cmd.MarkFlagRequired("file")
	return cmd
}
