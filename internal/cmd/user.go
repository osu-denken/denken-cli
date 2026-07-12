package cmd

import (
	"context"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

func newRegisterCmd(app *appContext) *cobra.Command {
	var email, password, passphrase string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "新規ユーザーを登録する (招待コードが必要)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegister(app, email, password, passphrase)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "学籍番号")
	cmd.Flags().StringVar(&password, "password", "", "パスワード")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "招待コード")
	return cmd
}

func runRegister(app *appContext, email, password, passphrase string) error {
	email, err := orPrompt(email, "学籍番号: ")
	if err != nil {
		return err
	}
	email = resolveEmail(email)
	if password == "" {
		if password, err = promptSecret("パスワード: "); err != nil {
			return err
		}
	}
	if passphrase, err = orPrompt(passphrase, "招待コード: "); err != nil {
		return err
	}
	ctx, cancel := newContext()
	defer cancel()
	res, err := app.client().Register(ctx, email, password, passphrase)
	if err != nil {
		return err
	}
	return saveSession(app, res)
}

func newUserCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "user", Short: "ユーザー情報の操作"}
	cmd.AddCommand(
		newUserInfoCmd(app), newUserUpdateCmd(app), newUserExistsCmd(app),
		newUserResetPasswordCmd(app), newUserVerifyEmailCmd(app), newUserProvidersCmd(app),
	)
	return cmd
}

func newUserInfoCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "info", "認証済みユーザーの詳細を取得する", (*api.Client).Info)
}

func newUserExistsCmd(app *appContext) *cobra.Command {
	return strFlagCmd(app, strFlag{
		use: "exists", short: "指定した学籍番号のユーザーが存在するか確認する",
		name: "email", help: "学籍番号", required: true,
		xform: resolveEmail, call: (*api.Client).Exists,
	})
}

func newUserUpdateCmd(app *appContext) *cobra.Command {
	var displayName, photoURL, password string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "表示名・プロフィール写真・パスワードを更新する",
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]string{}
			addIfSet(body, "displayName", displayName)
			addIfSet(body, "photoUrl", photoURL)
			addIfSet(body, "password", password)
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return app.client().UpdateUser(ctx, body)
			})
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "表示名")
	cmd.Flags().StringVar(&photoURL, "photo-url", "", "プロフィール写真 URL")
	cmd.Flags().StringVar(&password, "password", "", "新しいパスワード")
	return cmd
}

func newUserResetPasswordCmd(app *appContext) *cobra.Command {
	return strFlagCmd(app, strFlag{
		use: "reset-password", short: "パスワードリセットメールを送信する",
		name: "email", help: "学籍番号", required: true,
		xform: resolveEmail, call: (*api.Client).ResetPassword,
	})
}

func newUserVerifyEmailCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "verify-email", "確認メールを再送する", (*api.Client).VerifyEmail)
}

func newUserProvidersCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "providers", "紐づくログイン手段を表示する", (*api.Client).Providers)
}

func addIfSet(m map[string]string, key, value string) {
	if value != "" {
		m[key] = value
	}
}
