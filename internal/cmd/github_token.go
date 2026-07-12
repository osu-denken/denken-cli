package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newGithubTokenCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "token", Short: "GitHub PAT の確認、保存、削除 (要 BlogEdit 権限)"}
	cmd.AddCommand(
		authRawCmd(app, "get", "保存済み PAT の状態を確認する", app.client().GithubTokenGet),
		newGithubTokenSaveCmd(app),
		authRawCmd(app, "delete", "保存済み PAT を削除する", app.client().GithubTokenDelete),
	)
	return cmd
}

func newGithubTokenSaveCmd(app *appContext) *cobra.Command {
	var token string
	cmd := &cobra.Command{
		Use:   "save",
		Short: "GitHub PAT を暗号化して保存する",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if token == "" {
				if token, err = promptSecret("GitHub PAT: "); err != nil {
					return err
				}
			}
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().GithubTokenSave(ctx, token)
			})
		},
	}
	cmd.Flags().StringVar(&token, "token", "", "GitHub PAT")
	return cmd
}
