package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newPortalCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "portal", Short: "ポータル・外部サービス連携"}
	cmd.AddCommand(newPortalInfoCmd(app), newGithubCmd(app), newDiscordCmd(app))
	return cmd
}

func newPortalInfoCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "info", "ポータル用の情報をまとめて取得する", app.client().Portal)
}

func newDiscordCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "discord", Short: "Discord 連携"}
	cmd.AddCommand(authRawCmd(app, "invite", "Discord サーバーの招待コードを取得する", app.client().DiscordInvite))
	return cmd
}

func newGithubCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "github", Short: "GitHub 連携"}
	cmd.AddCommand(
		newGithubInviteCmd(app), newGithubJoinCmd(app),
		authRawCmd(app, "username", "連携済み GitHub ログイン名を取得する", app.client().GithubUsername),
		authRawCmd(app, "oauth-start", "GitHub OAuth の認可 URL を取得する (要 BlogEdit 権限)", app.client().GithubOAuthStart),
		newGithubTokenCmd(app),
	)
	return cmd
}

func newGithubInviteCmd(app *appContext) *cobra.Command {
	var email string
	cmd := &cobra.Command{
		Use:   "invite",
		Short: "指定メールを GitHub Organization に招待する (要 MemberManage 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().GithubInvite(ctx, email)
			})
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "招待するメールアドレス")
	cmd.MarkFlagRequired("email")
	return cmd
}

func newGithubJoinCmd(app *appContext) *cobra.Command {
	var username string
	cmd := &cobra.Command{
		Use:   "join",
		Short: "自分が Organization への招待を受け取る (連携済みなら username 省略可)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().GithubJoin(ctx, username)
			})
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "GitHub ユーザー名")
	return cmd
}
