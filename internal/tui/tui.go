package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/osu-denken/denken-cli/internal/api"
	"github.com/osu-denken/denken-cli/internal/config"
)

type screen int

const (
	menuScreen screen = iota
	resultScreen
)

type model struct {
	cfg     *config.Config
	items   []action
	cursor  int
	screen  screen
	loading bool
	result  string
	errMsg  string
}

type resultMsg struct {
	text string
	err  error
}

// Run は対話 TUI を起動する。
func Run(cfg *config.Config) error {
	m := model{cfg: cfg, items: actions()}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case resultMsg:
		m.loading = false
		m.screen = resultScreen
		if msg.err != nil {
			m.errMsg = msg.err.Error()
			m.result = ""
		} else {
			m.errMsg = ""
			m.result = msg.text
		}
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc":
		m.screen = menuScreen
		return m, nil
	case "up", "k":
		m.moveCursor(-1)
	case "down", "j":
		m.moveCursor(1)
	case "enter":
		if m.screen == menuScreen && !m.loading {
			m.loading = true
			return m, m.runSelected()
		}
	}
	return m, nil
}

func (m *model) moveCursor(delta int) {
	if m.screen != menuScreen {
		return
	}
	m.cursor = (m.cursor + delta + len(m.items)) % len(m.items)
}

func (m model) runSelected() tea.Cmd {
	act := m.items[m.cursor]
	token := m.cfg.Session.IDToken
	baseURL := m.cfg.BaseURL
	loggedIn := m.cfg.LoggedIn()
	return func() tea.Msg {
		if act.auth && !loggedIn {
			return resultMsg{err: errNotLoggedIn}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		raw, err := act.run(ctx, api.New(baseURL, token))
		if err != nil {
			return resultMsg{err: err}
		}
		return resultMsg{text: prettyJSON(raw)}
	}
}
