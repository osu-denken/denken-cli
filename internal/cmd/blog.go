package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newBlogCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "blog", Short: "ブログ記事の操作"}
	cmd.AddCommand(newBlogListCmd(app), newBlogGetCmd(app), newBlogUpdateCmd(app))
	return cmd
}

func newBlogListCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "記事の一覧を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			items, err := app.client().BlogList(ctx)
			if err != nil {
				return err
			}
			for _, it := range items {
				title, _ := it.Meta["title"].(string)
				fmt.Fprintf(app.out, "%-40s %s\n", it.Name, title)
			}
			return nil
		},
	}
}

func newBlogGetCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "get <slug>",
		Short: "記事の本文とメタデータを取得する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			entry, err := app.client().BlogGet(ctx, args[0])
			if err != nil {
				return err
			}
			data, _ := json.MarshalIndent(entry, "", "  ")
			return app.printJSON(data)
		},
	}
}

func newBlogUpdateCmd(app *appContext) *cobra.Command {
	var slug, title, contentFile string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "記事を新規作成または更新する (要 BlogEdit 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBlogUpdate(app, slug, title, contentFile)
		},
	}
	cmd.Flags().StringVar(&slug, "slug", "", "記事のスラッグ")
	cmd.Flags().StringVar(&title, "title", "", "記事タイトル")
	cmd.Flags().StringVar(&contentFile, "file", "", "本文 (Markdown) のファイルパス")
	cmd.MarkFlagRequired("slug")
	cmd.MarkFlagRequired("file")
	return cmd
}

func runBlogUpdate(app *appContext, slug, title, contentFile string) error {
	content, err := os.ReadFile(contentFile)
	if err != nil {
		return err
	}
	meta := map[string]any{}
	if title != "" {
		meta["title"] = title
	}
	ctx, cancel := newContext()
	defer cancel()
	if err := app.requireAuth(ctx); err != nil {
		return err
	}
	if err := app.client().BlogUpdate(ctx, slug, meta, string(content)); err != nil {
		return err
	}
	fmt.Fprintln(app.out, "更新しました:", slug)
	return nil
}
