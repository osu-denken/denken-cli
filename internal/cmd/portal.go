package cmd

import (
	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

func newPortalCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "portal", Short: "ポータル・外部サービス連携"}
	cmd.AddCommand(newPortalInfoCmd(app), newGithubCmd(app), newDiscordCmd(app))
	return cmd
}

func newPortalInfoCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "info", "ポータル用の情報をまとめて取得する", (*api.Client).Portal)
}

func newDiscordCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "discord", Short: "Discord 連携"}
	cmd.AddCommand(authRawCmd(app, "invite", "Discord サーバーの招待コードを取得する", (*api.Client).DiscordInvite))
	return cmd
}

func newGithubCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "github", Short: "GitHub 連携"}
	cmd.AddCommand(
		newGithubInviteCmd(app), newGithubJoinCmd(app),
		authRawCmd(app, "username", "連携済み GitHub ログイン名を取得する", (*api.Client).GithubUsername),
		authRawCmd(app, "oauth-start", "GitHub OAuth の認可 URL を取得する (要 BlogEdit 権限)", (*api.Client).GithubOAuthStart),
		newGithubTokenCmd(app),
	)
	return cmd
}

func newGithubInviteCmd(app *appContext) *cobra.Command {
	return strFlagCmd(app, strFlag{
		use: "invite", short: "指定メールを GitHub Organization に招待する (要 MemberManage 権限)",
		name: "email", help: "招待するメールアドレス", auth: true, required: true,
		call: (*api.Client).GithubInvite,
	})
}

func newGithubJoinCmd(app *appContext) *cobra.Command {
	return strFlagCmd(app, strFlag{
		use: "join", short: "自分が Organization への招待を受け取る (連携済みなら username 省略可)",
		name: "username", help: "GitHub ユーザー名", auth: true,
		call: (*api.Client).GithubJoin,
	})
}
