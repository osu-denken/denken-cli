package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newBlogCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "blog", Short: "ブログ記事の操作"}
	cmd.AddCommand(newBlogListCmd(app), newBlogGetCmd(app), newBlogUpdateCmd(app), newBlogEditCmd(app))
	return cmd
}

func newBlogEditCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <slug>",
		Short: "記事の本文を $EDITOR で開いて編集、保存する (要 BlogEdit 権限)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBlogEdit(app, args[0])
		},
	}
}

func runBlogEdit(app *appContext, slug string) error {
	var meta map[string]any
	fetch := func(ctx context.Context) (string, error) {
		entry, err := app.client().BlogGet(ctx, slug)
		if err != nil {
			return "", err
		}
		meta = entry.Meta
		return entry.Content, nil
	}
	save := func(ctx context.Context, content string) error {
		return app.client().BlogUpdate(ctx, slug, meta, content)
	}
	return app.runEdit(".md", slug, fetch, save)
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
	meta := map[string]any{}
	if title != "" {
		meta["title"] = title
	}
	return app.runFileUpdate(contentFile, slug, func(ctx context.Context, content string) error {
		return app.client().BlogUpdate(ctx, slug, meta, content)
	})
}
