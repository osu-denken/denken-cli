package cmd

import (
	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

func newSwitchBotCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "switchbot", Short: "部室の SwitchBot 操作 (要 SwitchBotControl 権限)", Aliases: []string{"lock"}}
	cmd.AddCommand(
		authRawCmd(app, "validate", "トークンが有効かどうかを確認する", (*api.Client).SwitchBotValidate),
		authRawCmd(app, "list", "デバイス一覧を返す", (*api.Client).SwitchBotList),
		authRawCmd(app, "lock", "施錠する", (*api.Client).SwitchBotLock),
		authRawCmd(app, "unlock", "解錠する", (*api.Client).SwitchBotUnlock),
	)
	return cmd
}
