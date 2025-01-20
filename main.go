package main

import (
	"fmt"
	"log"
	"time"

  "mikehaus/spoofify/auth"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// First we need to check if user is authorized
// If they are direct them to the main window
// If not do signup flow

type model int

type tickMsg time.Time

var (
	tickStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D56F4")).PaddingTop(2).PaddingBottom(4)
)

func main() {
	p := tea.NewProgram(model(5), tea.WithAltScreen())
  renderAuthList()
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func renderAuthList() {
  auth.AuthWindow()
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	return tickStyle.Render(fmt.Sprintf("\n\n     Hi. This program will end in %d seconds...", m))
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
