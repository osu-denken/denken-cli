package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func newTotpCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "totp", Short: "2段階認証 (TOTP) の設定"}
	cmd.AddCommand(newTotpSetupCmd(app), newTotpEnableCmd(app), newTotpDisableCmd(app))
	return cmd
}

func newTotpSetupCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "setup", "シークレットと QR を発行する (まだ有効化しない)", func(c cmdCtx) (any, error) {
		return app.client().TotpSetup(c.ctx)
	})
}

func newTotpEnableCmd(app *appContext) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "コードを検証して2段階認証を有効化する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authCodeAction(app, code, func(c cmdCtx, code string) (any, error) {
				return app.client().TotpEnable(c.ctx, code)
			})
		},
	}
	cmd.Flags().StringVar(&code, "code", "", "6桁の認証コード")
	cmd.MarkFlagRequired("code")
	return cmd
}

func newTotpDisableCmd(app *appContext) *cobra.Command {
	var code string
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "コード (またはリカバリコード) を検証して解除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authCodeAction(app, code, func(c cmdCtx, code string) (any, error) {
				return app.client().TotpDisable(c.ctx, code)
			})
		},
	}
	cmd.Flags().StringVar(&code, "code", "", "認証コードまたはリカバリコード")
	cmd.MarkFlagRequired("code")
	return cmd
}

func authCodeAction(app *appContext, code string, fn func(cmdCtx, string) (any, error)) error {
	ctx, cancel := newContext()
	defer cancel()
	if err := app.requireAuth(ctx); err != nil {
		return err
	}
	res, err := fn(cmdCtx{ctx}, code)
	if err != nil {
		return err
	}
	return app.printJSON(res.(json.RawMessage))
}
