package cmd

import (
	"github.com/spf13/cobra"
)

func newRegisterCmd(app *appContext) *cobra.Command {
	var email, password, passphrase string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "新規ユーザーを登録する (合言葉または招待コードが必要)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegister(app, email, password, passphrase)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "メールアドレス")
	cmd.Flags().StringVar(&password, "password", "", "パスワード (省略時はプロンプト)")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "合言葉または招待コード")
	return cmd
}

func runRegister(app *appContext, email, password, passphrase string) error {
	email, err := orPrompt(email, "メールアドレス: ")
	if err != nil {
		return err
	}
	if password == "" {
		if password, err = promptSecret("パスワード: "); err != nil {
			return err
		}
	}
	if passphrase, err = orPrompt(passphrase, "合言葉/招待コード: "); err != nil {
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
	return authRawCmd(app, "info", "認証済みユーザーの詳細を取得する", func(ctx cmdCtx) (any, error) {
		return app.client().Info(ctx.ctx)
	})
}

func newUserExistsCmd(app *appContext) *cobra.Command {
	var email string
	cmd := &cobra.Command{
		Use:   "exists",
		Short: "指定メールのユーザーが存在するか確認する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			raw, err := app.client().Exists(ctx, email)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "確認するメールアドレス")
	cmd.MarkFlagRequired("email")
	return cmd
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
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().UpdateUser(ctx, body)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "表示名")
	cmd.Flags().StringVar(&photoURL, "photo-url", "", "プロフィール写真 URL")
	cmd.Flags().StringVar(&password, "password", "", "新しいパスワード")
	return cmd
}

func newUserResetPasswordCmd(app *appContext) *cobra.Command {
	var email string
	cmd := &cobra.Command{
		Use:   "reset-password",
		Short: "パスワードリセットメールを送信する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			raw, err := app.client().ResetPassword(ctx, email)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "メールアドレス")
	cmd.MarkFlagRequired("email")
	return cmd
}

func newUserVerifyEmailCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "verify-email", "確認メールを再送する", func(ctx cmdCtx) (any, error) {
		return app.client().VerifyEmail(ctx.ctx)
	})
}

func newUserProvidersCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "providers", "紐づくログイン手段を表示する", func(ctx cmdCtx) (any, error) {
		return app.client().Providers(ctx.ctx)
	})
}

func addIfSet(m map[string]string, key, value string) {
	if value != "" {
		m[key] = value
	}
}
