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
	return &cobra.Command{
		Use:   "validate <code>",
		Short: "招待コードが有効かどうかを検証する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			raw, err := app.client().InviteValidate(ctx, args[0])
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
}

func newInviteCreateCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "create", "新しい招待コードを生成する (要 InviteCodeCreate 権限)", func(c cmdCtx) (any, error) {
		return app.client().InviteCreate(c.ctx)
	})
}

func newInviteDeleteCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <code>",
		Short: "指定した招待コードを無効化する",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().InviteDelete(ctx, args[0])
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
}
