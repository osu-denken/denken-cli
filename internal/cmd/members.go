package cmd

import (
	"context"

	"github.com/osu-denken/denken-cli/internal/api"
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
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().MembersList(ctx, status)
			})
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "pre-active/active/withdrawn/graduated/rejected")
	return cmd
}

func newMembersDetailCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "detail <id>", "部員一人の詳細を取得する", (*api.Client).MembersDetail)
}

func newMembersApproveCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "approve <id>", "仮部員を承認する (要 MemberApprove 権限)", (*api.Client).MembersApprove)
}

func newMembersRejectCmd(app *appContext) *cobra.Command {
	return memberIDCmd(app, "reject <id>", "仮部員の登録を却下する (要 MemberApprove 権限)", (*api.Client).MembersReject)
}
