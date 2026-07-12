package editor

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// fakeEditor は一時ファイルに固定文字列を追記するスクリプトを作り、そのパスを返す。
func fakeEditor(t *testing.T, appended string) string {
	t.Helper()
	dir := t.TempDir()
	if runtime.GOOS == "windows" {
		path := filepath.Join(dir, "ed.bat")
		writeScript(t, path, "@echo "+appended+">>%1\r\n")
		return path
	}
	path := filepath.Join(dir, "ed.sh")
	writeScript(t, path, "#!/bin/sh\necho '"+appended+"' >> \"$1\"\n")
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeScript(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestEditModifies(t *testing.T) {
	t.Setenv("DENKEN_EDITOR", fakeEditor(t, "added"))
	out, err := Edit("original\n", ".md")
	if err != nil {
		t.Fatalf("Edit: %v", err)
	}
	if out == "original\n" || len(out) <= len("original\n") {
		t.Fatalf("編集結果が反映されていない: %q", out)
	}
}

func TestEditUnchanged(t *testing.T) {
	// 何もしないエディタ (存在するが追記しない) を使う。
	noop := filepath.Join(t.TempDir(), scriptName())
	if runtime.GOOS == "windows" {
		writeScript(t, noop, "@rem noop\r\n")
	} else {
		writeScript(t, noop, "#!/bin/sh\nexit 0\n")
		os.Chmod(noop, 0o755)
	}
	t.Setenv("DENKEN_EDITOR", noop)
	_, err := Edit("same", ".md")
	if !errors.Is(err, ErrUnchanged) {
		t.Fatalf("ErrUnchanged を期待したが: %v", err)
	}
}

func scriptName() string {
	if runtime.GOOS == "windows" {
		return "ed.bat"
	}
	return "ed.sh"
}
