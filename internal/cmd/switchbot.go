package cmd

import (
	"github.com/spf13/cobra"
)

func newSwitchBotCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "switchbot", Short: "部室の SwitchBot 操作 (要 SwitchBotControl 権限)", Aliases: []string{"lock"}}
	cmd.AddCommand(
		authRawCmd(app, "validate", "トークンが有効かどうかを確認する", func(c cmdCtx) (any, error) {
			return app.client().SwitchBotValidate(c.ctx)
		}),
		authRawCmd(app, "list", "デバイス一覧を返す", func(c cmdCtx) (any, error) {
			return app.client().SwitchBotList(c.ctx)
		}),
		authRawCmd(app, "lock", "施錠する", func(c cmdCtx) (any, error) {
			return app.client().SwitchBotLock(c.ctx)
		}),
		authRawCmd(app, "unlock", "解錠する", func(c cmdCtx) (any, error) {
			return app.client().SwitchBotUnlock(c.ctx)
		}),
	)
	return cmd
}
