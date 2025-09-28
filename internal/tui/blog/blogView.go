package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/corpus/internal/blog"
)

type Model struct {
	Username        string
	Articles        []blog.Article
	SelectedArticle int
}

func InitialModel() Model {
	return Model{
		Username:        "coleflen",
		Articles:        []blog.Article{*blog.New("1", "title1", "content1"), *blog.New("1", "title1", "content1"), *blog.New("1", "title1", "content1"), *blog.New("1", "title1", "content1")},
		SelectedArticle: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.SelectedArticle > 0 {
				m.SelectedArticle--
			}

		case "down", "j":
			if m.SelectedArticle < len(m.Articles)-1 {
				m.SelectedArticle++
			}
		}

	}
	return m, nil
}

func (m Model) View() string {
	s := "List of Articles\n\n"

	// Iterate over our choices
	for i, choice := range m.Articles {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.SelectedArticle == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice.Title)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
