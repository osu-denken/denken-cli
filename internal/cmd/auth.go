package cmd

import (
	"context"
	"fmt"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/osu-denken/denken-cli/internal/config"
	"github.com/spf13/cobra"
)

func newLoginCmd(app *appContext) *cobra.Command {
	var email, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "ログインして ID トークンを保存する",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(app, email, password)
		},
	}
	cmd.Flags().StringVar(&email, "email", "", "学籍番号")
	cmd.Flags().StringVar(&password, "password", "", "パスワード")
	return cmd
}

func runLogin(app *appContext, email, password string) error {
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
	ctx, cancel := newContext()
	defer cancel()

	res, err := app.client().Login(ctx, email, password)
	if err != nil {
		return err
	}
	if res.MFARequired {
		if res, err = completeMFA(ctx, app, res.MFAPendingToken); err != nil {
			return err
		}
	}
	return saveSession(app, res)
}

func completeMFA(ctx context.Context, app *appContext, pendingToken string) (*api.AuthResult, error) {
	code, err := promptLine("2段階認証コード (6桁): ")
	if err != nil {
		return nil, err
	}
	return app.client().LoginTotp(ctx, pendingToken, code)
}

func saveSession(app *appContext, res *api.AuthResult) error {
	app.applyAuth(res)
	if err := config.Save(app.cfg); err != nil {
		return err
	}
	fmt.Fprintln(app.out, "ログインしました:", app.cfg.Session.Email)
	return nil
}

func newLogoutCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "保存済みトークンを削除する",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.cfg.Session = config.Session{}
			if err := config.Save(app.cfg); err != nil {
				return err
			}
			fmt.Fprintln(app.out, "ログアウトしました")
			return nil
		},
	}
}

func newWhoamiCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "whoami", "認証済みユーザーの情報を表示する", (*api.Client).Info)
}

func newRefreshCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "リフレッシュトークンで ID トークンを更新する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.refresh(ctx); err != nil {
				return err
			}
			fmt.Fprintln(app.out, "トークンを更新しました")
			return nil
		},
	}
}
