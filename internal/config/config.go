package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// DefaultBaseURL は OSU-Denken Web API の本番エンドポイント。
const DefaultBaseURL = "https://api.osudenken4dev.workers.dev"

// Session は Firebase 認証で得たトークン一式を保持する。
type Session struct {
	Email        string    `json:"email,omitempty"`
	IDToken      string    `json:"idToken,omitempty"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt,omitempty"`
}

// Config は永続化される CLI 設定。
type Config struct {
	BaseURL string  `json:"baseUrl"`
	Session Session `json:"session"`
}

// LoggedIn は有効な ID トークンを保持しているかを返す。
func (c *Config) LoggedIn() bool {
	return c.Session.IDToken != ""
}

// Expired はトークンの有効期限が切れている(もしくは間近)かを返す。
func (s *Session) Expired() bool {
	if s.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(s.ExpiresAt.Add(-30 * time.Second))
}

// Path は設定ファイルの絶対パスを返す。
// 環境変数 DENKEN_CONFIG が優先される。
func Path() (string, error) {
	if p := os.Getenv("DENKEN_CONFIG"); p != "" {
		return p, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "denken-cli", "config.json"), nil
}

// Load は設定を読み込む。ファイルが無ければ既定値を返す。
func Load() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{BaseURL: DefaultBaseURL}, nil
	}
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	return cfg, nil
}

// Save は設定をディスクへ書き出す。ディレクトリは必要に応じて作成する。
func Save(cfg *Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
