package cmd

import (
	"context"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

func newTotpCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "totp", Short: "2段階認証 (TOTP) の設定"}
	cmd.AddCommand(newTotpSetupCmd(app), newTotpEnableCmd(app), newTotpDisableCmd(app))
	return cmd
}

func newTotpSetupCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "setup", "シークレットと QR を発行する (まだ有効化しない)", (*api.Client).TotpSetup)
}

func newTotpEnableCmd(app *appContext) *cobra.Command {
	return codeCmd(app, "enable", "コードを検証して2段階認証を有効化する", "6桁の認証コード", (*api.Client).TotpEnable)
}

func newTotpDisableCmd(app *appContext) *cobra.Command {
	return codeCmd(app, "disable", "コード (またはリカバリコード) を検証して解除する", "認証コードまたはリカバリコード", (*api.Client).TotpDisable)
}

// codeCmd は --code を1つ取り認証必須で JSON を表示するコマンドを作る。
func codeCmd(app *appContext, use, short, codeHelp string, call func(*api.Client, context.Context, string) (rawJSON, error)) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return call(app.client(), ctx, code)
			})
		},
	}
	cmd.Flags().StringVar(&code, "code", "", codeHelp)
	cmd.MarkFlagRequired("code")
	return cmd
}
