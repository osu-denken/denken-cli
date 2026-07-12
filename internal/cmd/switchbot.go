package cmd

import (
	"github.com/spf13/cobra"
)

func newSwitchBotCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "switchbot", Short: "部室の SwitchBot 操作 (要 SwitchBotControl 権限)", Aliases: []string{"lock"}}
	cmd.AddCommand(
		authRawCmd(app, "validate", "トークンが有効かどうかを確認する", app.client().SwitchBotValidate),
		authRawCmd(app, "list", "デバイス一覧を返す", app.client().SwitchBotList),
		authRawCmd(app, "lock", "施錠する", app.client().SwitchBotLock),
		authRawCmd(app, "unlock", "解錠する", app.client().SwitchBotUnlock),
	)
	return cmd
}
