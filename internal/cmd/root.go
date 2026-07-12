package cmd

import (
	"fmt"
	"os"

	"github.com/osu-denken/denken-cli/internal/config"
	"github.com/osu-denken/denken-cli/internal/tui"
	"github.com/spf13/cobra"
)

// version はビルド時に -ldflags で埋め込まれる。
var version = "dev"

// Execute はエントリポイントから呼ばれ、CLI を実行する。
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	app := &appContext{out: os.Stdout}
	var baseURL, token string

	root := &cobra.Command{
		Use:           "denken-cli",
		Short:         "OSU-Denken Web API の CLI/TUI クライアント",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initApp(app, baseURL, token)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run(app.cfg)
		},
	}

	root.PersistentFlags().StringVar(&baseURL, "base-url", "", "API のベース URL を上書きする")
	root.PersistentFlags().StringVar(&token, "token", "", "ID トークンを一時的に上書きする")

	root.AddCommand(
		newLoginCmd(app), newLogoutCmd(app), newWhoamiCmd(app), newRefreshCmd(app),
		newRegisterCmd(app), newUserCmd(app), newTotpCmd(app),
		newBlogCmd(app), newInviteCmd(app), newMembersCmd(app),
		newPrivatePostsCmd(app), newImageCmd(app), newPortalCmd(app),
		newSwitchBotCmd(app), newPagesCmd(app), newLogsCmd(app),
		newConfigCmd(app), newPingCmd(app),
	)
	return root
}

func initApp(app *appContext, baseURL, token string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	app.cfg = cfg
	app.token = token
	return nil
}
