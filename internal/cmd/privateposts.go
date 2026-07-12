package cmd

import (
	"fmt"
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
	current, title, err := fetchPrivatePost(app, slug)
	if err != nil {
		return err
	}
	edited, ok, err := app.openEditor(current, ".md")
	if err != nil || !ok {
		return err
	}
	ctx, cancel := newContext()
	defer cancel()
	if _, err := app.client().PrivatePostUpdate(ctx, slug, title, edited); err != nil {
		return err
	}
	fmt.Fprintln(app.out, "更新しました:", slug)
	return nil
}

func fetchPrivatePost(app *appContext, slug string) (content, title string, err error) {
	ctx, cancel := newContext()
	defer cancel()
	if err = app.requireAuth(ctx); err != nil {
		return "", "", err
	}
	raw, err := app.client().PrivatePostGet(ctx, slug)
	if err != nil {
		return "", "", err
	}
	return jsonField(raw, "content"), jsonField(raw, "title"), nil
}

func newPPListCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "list", "非公開記事を一覧する (要 PrivatePostView 権限)", func(c cmdCtx) (any, error) {
		return app.client().PrivatePostList(c.ctx)
	})
}

func newPPGetCmd(app *appContext) *cobra.Command {
	return authSlugCmd(app, "get <slug>", "非公開記事の本文を取得する", func(c cmdCtx, slug string) (any, error) {
		return app.client().PrivatePostGet(c.ctx, slug)
	})
}

func newPPDeleteCmd(app *appContext) *cobra.Command {
	return authSlugCmd(app, "delete <slug>", "非公開記事を削除する (要 PrivatePostEdit 権限)", func(c cmdCtx, slug string) (any, error) {
		return app.client().PrivatePostDelete(c.ctx, slug)
	})
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
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().PrivatePostUpdate(ctx, slug, title, string(content))
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&slug, "slug", "", "スラッグ (英小文字・数字・ハイフン)")
	cmd.Flags().StringVar(&title, "title", "", "タイトル")
	cmd.Flags().StringVar(&file, "file", "", "本文ファイルのパス")
	cmd.MarkFlagRequired("slug")
	cmd.MarkFlagRequired("file")
	return cmd
}
