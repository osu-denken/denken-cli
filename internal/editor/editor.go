package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// ErrUnchanged は編集後も内容が変わらなかったことを表す。
var ErrUnchanged = fmt.Errorf("内容が変更されませんでした")

// Edit は initial を一時ファイルに書き出してエディタで開き、保存後の内容を返す。
// ext は一時ファイルの拡張子 (例: ".md")。内容が変わらなければ ErrUnchanged を返す。
func Edit(initial, ext string) (string, error) {
	file, err := writeTemp(initial, ext)
	if err != nil {
		return "", err
	}
	defer os.Remove(file)

	if err := launch(file); err != nil {
		return "", err
	}
	edited, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	if string(edited) == initial {
		return "", ErrUnchanged
	}
	return string(edited), nil
}

func writeTemp(initial, ext string) (string, error) {
	f, err := os.CreateTemp("", "denken-*"+ext)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.WriteString(initial); err != nil {
		return "", err
	}
	return f.Name(), nil
}

func launch(file string) error {
	name, args := editorCommand()
	args = append(args, file)
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// editorCommand は使用するエディタと引数を決める。
// 優先順位: DENKEN_EDITOR > VISUAL > EDITOR > プラットフォーム既定。
func editorCommand() (string, []string) {
	for _, key := range []string{"DENKEN_EDITOR", "VISUAL", "EDITOR"} {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			fields := strings.Fields(v)
			return fields[0], fields[1:]
		}
	}
	if runtime.GOOS == "windows" {
		return "notepad", nil
	}
	return "vi", nil
}

// Ext はパスから拡張子を返す (無ければ .txt)。
func Ext(path string) string {
	if e := filepath.Ext(path); e != "" {
		return e
	}
	return ".txt"
}
