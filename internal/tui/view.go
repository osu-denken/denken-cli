package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var errNotLoggedIn = errors.New("ログインが必要です (denken-cli login)")

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func (m model) View() string {
	switch {
	case m.loading:
		return "\n  実行中...\n"
	case m.screen == resultScreen:
		return m.resultView()
	default:
		return m.menuView()
	}
}

func (m model) menuView() string {
	b := &strings.Builder{}
	fmt.Fprintln(b, titleStyle.Render("denken-cli"))
	fmt.Fprintln(b, dimStyle.Render(m.statusLine()))
	fmt.Fprintln(b)
	for i, it := range m.items {
		cursor := "  "
		label := it.title
		if it.auth {
			label += dimStyle.Render(" (要ログイン)")
		}
		if i == m.cursor {
			cursor = selectedStyle.Render("> ")
			label = selectedStyle.Render(it.title)
		}
		fmt.Fprintf(b, "%s%s\n", cursor, label)
	}
	fmt.Fprintln(b, dimStyle.Render("\n↑/↓ 移動  Enter 実行  q 終了"))
	return b.String()
}

func (m model) statusLine() string {
	if m.cfg.LoggedIn() {
		return "ログイン中: " + m.cfg.Session.Email
	}
	return "未ログイン"
}

func (m model) resultView() string {
	b := &strings.Builder{}
	fmt.Fprintln(b, titleStyle.Render(m.items[m.cursor].title))
	fmt.Fprintln(b)
	if m.errMsg != "" {
		fmt.Fprintln(b, errorStyle.Render("エラー: "+m.errMsg))
	} else {
		fmt.Fprintln(b, m.result)
	}
	fmt.Fprintln(b, dimStyle.Render("\nEsc で戻る  q 終了"))
	return b.String()
}
