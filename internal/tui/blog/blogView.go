package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/corpus/internal/blog"
)

type BlogVM struct {
	Username        string
	Articles        []*blog.Article
	SelectedArticle int
	articleVM       *ArticleVM
	renderArticle   bool
}

func InitialModel() BlogVM {
	return BlogVM{
		Username:        "coleflen",
		Articles:        []*blog.Article{blog.New("1", "title1", "content1"), blog.New("1", "title2", "content2"), blog.New("1", "title3", "content3"), blog.New("1", "title4", "content4")},
		SelectedArticle: 0,
		renderArticle:   false,
	}
}

func (m *BlogVM) Init() tea.Cmd {
	return nil
}

// ok so I want to add many screens and allow a hierarchical representation. One way is to implement
// multiple models that follow the pattern then have a master model with a list of models. It would only
// have to swap the model then I could call the associated methods as normal
func (m *BlogVM) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			if m.articleVM == nil {
				m.articleVM = &ArticleVM{
					article:  m.Articles[m.SelectedArticle],
					ready:    false,
					viewport: viewport.New(80, 20),
				}
			}

			m.renderArticle = true
			m.articleVM.Update(msg)

		case "esc", "escape":
			m.renderArticle = false
		}

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.articleVM != nil {
			m.articleVM.Update(msg)
		}
	}
	return m, nil
}

func (m BlogVM) View() string {
	// header
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}).Background(lipgloss.AdaptiveColor{Light: "#7c3aed", Dark: "#7c3aed"}).Padding(0, 1)
	headerLeft := " Corpus "
	headerRight := fmt.Sprintf("user: %s  •  %d articles", m.Username, len(m.Articles))
	header := headerStyle.Render(headerLeft + " " + headerRight)

	// footer
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#9a9a9a", Dark: "#6b6b6b"}).Padding(0, 1)
	footer := footerStyle.Render("q: quit  •  enter: open  •  esc: back")

	if m.renderArticle {
		// show header, article view, footer
		return header + "\n\n" + m.articleVM.View() + "\n" + footer
	}

	s := header + "\n\nList of Articles\n\n"

	var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}).Bold(true)
	var unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#9a9a9a", Dark: "#6b6b6b"})

	// Iterate over our choices
	// totalWidth is the target width for the title+author line. Adjust as needed
	const totalWidth = 80

	for i, choice := range m.Articles {
		// Prepare visible pieces
		title := choice.Title
		author := choice.AuthorID
		excerpt := func(s string, words int) string {
			fs := strings.Fields(s)
			if len(fs) <= words {
				return strings.Join(fs, " ")
			}
			return strings.Join(fs[:words], " ") + "..."
		}(choice.Content, 12)

		// Build prefix and author text
		var prefix string
		if m.SelectedArticle == i {
			prefix = "> "
		} else {
			prefix = "  "
		}
		authorText := author

		// Compute available space for title
		prefixLen := lipgloss.Width(prefix)
		authorLen := lipgloss.Width(authorText)
		titleMax := totalWidth - prefixLen - authorLen - 1
		if titleMax < 1 {
			titleMax = 1
		}

		// Truncate title to fit
		if lipgloss.Width(title) > titleMax {
			title = truncateString(title, titleMax)
		}

		// Compute padding spaces so author is right-aligned to totalWidth
		spaces := titleMax - lipgloss.Width(title) + 1
		if spaces < 1 {
			spaces = 1
		}

		if m.SelectedArticle == i {
			s += selectedStyle.Render(prefix+title) + strings.Repeat(" ", spaces) + selectedStyle.Render(authorText) + "\n"
			s += "  " + selectedStyle.Render(excerpt) + "\n"
		} else {
			s += unselectedStyle.Render(prefix+title) + strings.Repeat(" ", spaces) + unselectedStyle.Render(authorText) + "\n"
			s += "  " + unselectedStyle.Render(excerpt) + "\n"
		}
	}

	// The footer
	s += footer

	// Send the UI for rendering
	return s
}

// truncateString shortens s to fit within max visible columns, adding "..." if truncated.
func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= max {
		return s
	}
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		candidate := string(runes[:i+1]) + "..."
		if lipgloss.Width(candidate) > max {
			if i == 0 {
				return "..."
			}
			return string(runes[:i]) + "..."
		}
	}
	return s
}
