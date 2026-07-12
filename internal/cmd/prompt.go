package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// promptLine は1行の入力を求めて返す (改行は除去する)。
func promptLine(label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// promptSecret はエコーを伏せて秘密情報の入力を求める。
func promptSecret(label string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	data, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// orPrompt は値が空ならプロンプトで補う。
func orPrompt(value, label string) (string, error) {
	if value != "" {
		return value, nil
	}
	return promptLine(label)
}
