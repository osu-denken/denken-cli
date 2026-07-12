package cmd

import (
	"github.com/spf13/cobra"
)

func newImageCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{Use: "image", Short: "ブログ用画像の操作"}
	cmd.AddCommand(newImageListCmd(app), newImageUploadCmd(app), newImageDeleteCmd(app))
	return cmd
}

func newImageListCmd(app *appContext) *cobra.Command {
	return authRawCmd(app, "list", "アップロード済み画像を一覧する (要 BlogEdit 権限)", func(c cmdCtx) (any, error) {
		return app.client().ImageList(c.ctx)
	})
}

func newImageUploadCmd(app *appContext) *cobra.Command {
	var file, name string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "画像をアップロードする (要 ImageUpload 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().ImageUpload(ctx, file, name)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "画像ファイルのパス (jpg/png/webp/gif, 最大20MB)")
	cmd.Flags().StringVar(&name, "name", "", "ファイル名 (省略時は UUID)")
	cmd.MarkFlagRequired("file")
	return cmd
}

func newImageDeleteCmd(app *appContext) *cobra.Command {
	var filename, sha string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "画像を削除する (要 ImageDelete 権限)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			if err := app.requireAuth(ctx); err != nil {
				return err
			}
			raw, err := app.client().ImageDelete(ctx, filename, sha)
			if err != nil {
				return err
			}
			return app.printJSON(raw)
		},
	}
	cmd.Flags().StringVar(&filename, "filename", "", "削除する画像のファイル名")
	cmd.Flags().StringVar(&sha, "sha", "", "画像の SHA")
	cmd.MarkFlagRequired("filename")
	cmd.MarkFlagRequired("sha")
	return cmd
}
