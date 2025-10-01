package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/corpus/internal/blog"
	"github.com/muesli/reflow/wordwrap"
)

type ArticleVM struct {
	article  *blog.Article
	viewport viewport.Model
	ready    bool
}

func (m *ArticleVM) Init() tea.Cmd {
	return nil
}

func (m *ArticleVM) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" {
			return m, tea.Quit
		}
		if k := msg.String(); k == "enter" {
			content := wordwrap.String(m.article.Content, m.viewport.Width)
			m.viewport.SetContent(content)
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

		content := wordwrap.String(m.article.Content, m.viewport.Width)
		m.viewport.SetContent(content)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ArticleVM) View() string {

	return fmt.Sprintf("%s\n", m.viewport.View())
}
