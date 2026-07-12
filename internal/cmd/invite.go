package cmd

import (
	"github.com/spf13/cobra"
)

func newInviteCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "invite", Short: "招待コードの操作"}
	cmd.AddCommand(newInviteValidateCmd(app), newInviteCreateCmd(app), newInviteDeleteCmd(app))
	return cmd
}

func newInviteValidateCmd(app *appContext) *cobra.Command {
	return argCmd(app, "validate <code>", "招待コードが有効かどうかを検証する", false, app.client().InviteValidate)
}

func newInviteCreateCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "create", "新しい招待コードを生成する (要 InviteCodeCreate 権限)", app.client().InviteCreate)
}

func newInviteDeleteCmd(app *appContext) *cobra.Command {
	return argCmd(app, "delete <code>", "指定した招待コードを無効化する", true, app.client().InviteDelete)
}
