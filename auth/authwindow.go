package auth 

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var txtStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string {
	return i.title
}

func (i item) FilterValue() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

type model struct {
	loggedIn bool
	list     list.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := txtStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
  return txtStyle.Render(m.list.View())
}

func AuthWindow() {
	m := initialModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

  if _, err := p.Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	items := []list.Item{
		item{title: "Login", desc: "Login with your Spotify login info"},
		item{title: "Quit", desc: "Quit application"},
	}

	m := model{
		loggedIn: false,
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
	}

	m.list.Title = "Please select an option"

	return m
}
