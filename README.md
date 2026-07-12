# denken-cli
OSU-Denken Web APIのCLI/TUIクライアントツール

## 使い方

```bash
denken-cli [command] [flags]
```

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
