package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/osu-denken/denken-cli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "CLI 設定の確認・変更"}
	cmd.AddCommand(newConfigShowCmd(app), newConfigSetURLCmd(app))
	return cmd
}

func newConfigShowCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "現在の設定を表示する (トークンは伏せる)",
		RunE: func(cmd *cobra.Command, args []string) error {
			view := map[string]any{
				"baseUrl":  app.cfg.BaseURL,
				"email":    app.cfg.Session.Email,
				"loggedIn": app.cfg.LoggedIn(),
			}
			data, _ := json.MarshalIndent(view, "", "  ")
			fmt.Fprintln(app.out, string(data))
			return nil
		},
	}
}

func newConfigSetURLCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "set-url <url>",
		Short: "API のベース URL を保存する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app.cfg.BaseURL = args[0]
			if err := config.Save(app.cfg); err != nil {
				return err
			}
			fmt.Fprintln(app.out, "ベース URL を保存しました:", args[0])
			return nil
		},
	}
}
