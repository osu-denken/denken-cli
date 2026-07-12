package cmd

import (
	"context"
	"strconv"

	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/spf13/cobra"
)

// memberIDCmd は <id> を1つ取り認証必須で JSON を表示するコマンドを作る。
func memberIDCmd(app *appContext, use, short string, call func(*api.Client, context.Context, int) (rawJSON, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
				return call(app.client(), ctx, id)
			})
		},
	}
}

func newMembersUpdateCmd(app *appContext) *cobra.Command {
	f := &memberUpdateFlags{}
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "部員情報を更新する (項目ごとに必要権限が異なる)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return runMembersUpdate(app, id, f)
		},
	}
	f.bind(cmd)
	return cmd
}

func runMembersUpdate(app *appContext, id int, f *memberUpdateFlags) error {
	body := f.body()
	body["id"] = id
	return app.runRaw(true, func(ctx context.Context) (rawJSON, error) {
		return app.client().MembersUpdate(ctx, body)
	})
}

// memberUpdateFlags は members update の可変項目。
type memberUpdateFlags struct {
	name, furigana, email, tel, status, joinDate string
	roleBits, permBits                           int
}

func (f *memberUpdateFlags) bind(cmd *cobra.Command) {
	fl := cmd.Flags()
	fl.StringVar(&f.name, "name", "", "氏名")
	fl.StringVar(&f.furigana, "furigana", "", "ふりがな")
	fl.StringVar(&f.email, "email", "", "メールアドレス")
	fl.StringVar(&f.tel, "tel", "", "電話番号 (幹部のみ)")
	fl.StringVar(&f.status, "status", "", "在籍状態 (要 MemberDelete)")
	fl.StringVar(&f.joinDate, "join-date", "", "入部日 (YYYY-MM-DD)")
	fl.IntVar(&f.roleBits, "role-bits", -1, "役職ビット (要 MemberRoleEdit)")
	fl.IntVar(&f.permBits, "perm-bits", -1, "権限ビット (要 MemberPermissionEdit)")
}

func (f *memberUpdateFlags) body() map[string]any {
	body := map[string]any{}
	setStr(body, "name", f.name)
	setStr(body, "furigana", f.furigana)
	setStr(body, "email", f.email)
	setStr(body, "tel", f.tel)
	setStr(body, "status", f.status)
	setStr(body, "joinDate", f.joinDate)
	if f.roleBits >= 0 {
		body["roleBits"] = f.roleBits
	}
	if f.permBits >= 0 {
		body["permBits"] = f.permBits
	}
	return body
}

func setStr(m map[string]any, key, value string) {
	if value != "" {
		m[key] = value
	}
}
