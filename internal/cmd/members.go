package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newMembersCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "members", Short: "部員名簿の操作 (要 MemberManage 権限)"}
	cmd.AddCommand(
		newMembersListCmd(app), newMembersDetailCmd(app),
		newMembersApproveCmd(app), newMembersRejectCmd(app), newMembersUpdateCmd(app),
	)
	return cmd
}

func newMembersListCmd(app *appContext) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "部員一覧を取得する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().MembersList(ctx, status)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "pre-active/active/withdrawn/graduated/rejected")
	return cmd
}

func newMembersDetailCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "detail <id>", "部員一人の詳細を取得する", func(c cmdCtx, id int) (any, error) {
		return app.client().MembersDetail(c.ctx, id)
	})
}

func newMembersApproveCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "approve <id>", "仮部員を承認する (要 MemberApprove 権限)", func(c cmdCtx, id int) (any, error) {
		return app.client().MembersApprove(c.ctx, id)
	})
}

func newMembersRejectCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "reject <id>", "仮部員の登録を却下する (要 MemberApprove 権限)", func(c cmdCtx, id int) (any, error) {
		return app.client().MembersReject(c.ctx, id)
	})
}
