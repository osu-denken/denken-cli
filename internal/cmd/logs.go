package cmd

import (
	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

func newLogsCmd(app *appContext) *cobra.Command {
	q := api.LogsQuery{}
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "操作ログを一覧する (要 LogView 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().LogsList(ctx, q)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&q.Type, "type", "", "種別で絞り込む")
	cmd.Flags().StringVar(&q.Cursor, "cursor", "", "続きを取得するカーソル")
	cmd.Flags().IntVar(&q.Limit, "limit", 0, "1〜100 (既定 50)")
	return cmd
}

func newPingCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "ping",
		Short: "サーバーの稼働確認を行う",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			res, err := app.client().Ping(ctx)
			if err != nil {
				return err
			}
			cmd.Println(res)
			return nil
		},
	}
}
