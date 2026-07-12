# denken-cli
OSU-Denken Web APIのCLI/TUIクライアントツール

## 使い方

```bash
denken-cli [command] [flags]
```

引数なしで実行すると対話 TUI が起動する。

### グローバルフラグ

| フラグ | 説明 |
| --- | --- |
| `--base-url <url>` | API のベース URL を上書きする (既定は本番) |
| `--token <token>` | ID トークンを一時的に上書きする |
| `-v, --version` | バージョンを表示する |
| `-h, --help` | ヘルプを表示する |

ログインすると ID トークンとリフレッシュトークンが設定ファイルに保存され、以降の認証コマンドで自動的に使われる (期限切れ時は自動更新)。

### コマンド一覧

#### 認証

| コマンド | 説明 |
| --- | --- |
| `login [--email --password]` | ログインしてIDトークンを保存する (TOTP 自動対応) |
| `logout` | 保存済みトークンを削除する |
| `whoami` | 認証済みユーザーの情報を表示する |
| `refresh` | リフレッシュトークンでIDトークンを更新する |
| `register [--email --password --passphrase]` | 新規ユーザーを登録する (招待コードが必要) |
| `ping` | サーバーの稼働確認を行う |

#### `user` — ユーザー情報

| コマンド | 説明 |
| --- | --- |
| `user info` | 認証済みユーザーの詳細を取得する |
| `user update [--display-name --photo-url --password]` | 表示名、写真、パスワードを更新する |
| `user exists --email <mail>` | 指定メールのユーザーが存在するか確認する |
| `user reset-password --email <mail>` | パスワードリセットメールを送信する |
| `user verify-email` | 確認メールを再送する |
| `user providers` | 紐づくログイン手段を表示する |

#### `totp` — 2段階認証

| コマンド | 説明 |
| --- | --- |
| `totp setup` | シークレットと QR を発行する (まだ有効化しない) |
| `totp enable --code <code>` | コードを検証して有効化する |
| `totp disable --code <code>` | コード/リカバリコードを検証して解除する |

#### `blog` — ブログ記事

| コマンド | 説明 |
| --- | --- |
| `blog list` | 記事の一覧を取得する |
| `blog get <slug>` | 記事の本文とメタデータを取得する |
| `blog update --slug <s> --file <path> [--title]` | 記事を新規作成/更新する (要 BlogEdit) |
| `blog edit <slug>` | 記事の本文を $EDITOR で開いて編集、保存する (要 BlogEdit) |

#### `invite` — 招待コード

| コマンド | 説明 |
| --- | --- |
| `invite validate <code>` | 招待コードが有効かどうかを検証する |
| `invite create` | 新しい招待コードを生成する (要 InviteCodeCreate) |
| `invite delete <code>` | 指定した招待コードを無効化する |

#### `members` — 部員名簿 (要 MemberManage)

| コマンド | 説明 |
| --- | --- |
| `members list [--status <s>]` | 部員一覧を取得する |
| `members detail <id>` | 部員一人の詳細を取得する |
| `members approve <id>` | 仮部員を承認する (要 MemberApprove) |
| `members reject <id>` | 仮部員の登録を却下する (要 MemberApprove) |
| `members update <id> [--name --furigana --email --tel --status --join-date --role-bits --perm-bits]` | 部員情報を更新する (項目ごとに必要権限が異なる) |

#### `private-posts` (別名 `pp`) — 非公開記事

| コマンド | 説明 |
| --- | --- |
| `private-posts list` | 非公開記事を一覧する (要 PrivatePostView) |
| `private-posts get <slug>` | 本文を取得する (要 PrivatePostView) |
| `private-posts update --slug <s> --file <path> [--title]` | 新規作成/上書きする (要 PrivatePostEdit) |
| `private-posts delete <slug>` | 削除する (要 PrivatePostEdit) |
| `private-posts edit <slug>` | 本文を `$EDITOR` で開いて編集、保存する (要 PrivatePostEdit) |

#### `image` — ブログ用画像

| コマンド | 説明 |
| --- | --- |
| `image list` | アップロード済み画像を一覧する (要 BlogEdit) |
| `image upload --file <path> [--name]` | 画像をアップロードする (要 ImageUpload) |
| `image delete --filename <f> --sha <sha>` | 画像を削除する (要 ImageDelete) |

#### `portal` — ポータル / 外部サービス連携

| コマンド | 説明 |
| --- | --- |
| `portal info` | ポータル用の情報をまとめて取得する |
| `portal github username` | 連携済み GitHub ログイン名を取得する |
| `portal github invite --email <mail>` | メールを GitHub Org に招待する (要 MemberManage) |
| `portal github join [--username]` | 自分が Org 招待を受け取る |
| `portal github oauth-start` | GitHub OAuth 認可 URL を取得する (要 BlogEdit) |
| `portal github token get\|save\|delete` | GitHub PAT の確認/保存/削除 (要 BlogEdit) |
| `portal discord invite` | Discord サーバーの招待コードを取得する |

#### `switchbot` — 部室スマートロック (要 SwitchBotControl)

| コマンド | 説明 |
| --- | --- |
| `switchbot validate` | トークンが有効かどうかを確認する |
| `switchbot list` | デバイス一覧を返す |
| `switchbot lock` | 施錠する |
| `switchbot unlock` | 解錠する |

#### `pages` — 固定ページ編集 (要 PageEdit)

| コマンド | 説明 |
| --- | --- |
| `pages terminal get [--page]` | ターミナルの welcome.md を取得する |
| `pages terminal update --file <path>` | welcome.md を更新し再ビルドを起動する |
| `pages terminal edit [--page]` | welcome.md を `$EDITOR` で開いて編集、保存する |
| `pages site list` | 編集できるファイルの一覧を返す |
| `pages site get --path <p>` | 固定ページの中身を取得する |
| `pages site update --path <p> --file <path>` | 固定ページを更新し再ビルドを起動する |
| `pages site edit --path <p>` | 固定ページを `$EDITOR` で開いて編集、保存する (拡張子は path から判定) |

> `edit` 系はサーバーから現在の内容を取得してエディタで開き、保存して閉じると差分をそのまま API に反映する。内容が変わらなければ送信しない。使用するエディタは `DENKEN_EDITOR` → `VISUAL` → `EDITOR` の順で決定し、未設定なら Windows は `notepad`、その他は `vi`。

#### `logs` / `config`

| コマンド | 説明 |
| --- | --- |
| `logs [--type --cursor --limit]` | 操作ログを一覧する (要 LogView) |
| `config show` | 現在の設定を表示する (トークンは伏せる) |
| `config set-url <url>` | API のベース URL を保存する |

## 開発メモ
### 技術スタック
- 開発言語: Go
- フレームワーク: Cobra (CLI), Bubble Tea (TUI)
- 動作環境: Windows / macOS / Linux (GitHub Actionsでのビルド)
- API: [OSU-Denken Web API](https://github.com/osu-denken/web-api)

### リリース
タグを作成することでGitHub Actionsがトリガーされ、ビルドとリリースが行われる。
```bash
git tag v0.0.1
git push origin v0.0.1
```
