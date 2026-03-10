package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"my_skill_manager/internal/repo"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func truncate(s string, limit int) string {
	words := strings.Fields(s)
	if len(words) <= limit {
		return s
	}
	return strings.Join(words[:limit], " ") + "..."
}

func RunSkillsTUI(repoName string, skills []repo.Skill) error {
	items := make([]list.Item, len(skills))
	for i, s := range skills {
		desc := s.Description
		if desc == "" {
			desc = "No description provided."
		}
		items[i] = item{
			title: s.Name,
			desc:  "  " + truncate(desc, 50),
		}
	}

	delegate := list.NewDefaultDelegate()
	// Each item takes 3 lines in the view + spacing.
	// Default delegate uses 2 lines by default, but we can't easily force "exactly 3 items visible"
	// without knowing window height. We'll stick to default delegate which is standard.

	m := model{list: list.New(items, delegate, 0, 0)}
	m.list.Title = fmt.Sprintf("Skills in repo: %s", repoName)
	m.list.SetFilteringEnabled(true)

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
