package auth

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var txtStyle = lipgloss.NewStyle().Margin(1, 2)

var ENV_CLIENT_ID = os.Getenv("SPOTIFY_CLIENT_ID")
var redirect_uri = "http://localhost:8080/callback"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
		if msg.String() == "enter" {
			return m, openBrowserAuth()
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

func openBrowserAuth(url string) tea.Cmd {
	// TODO: this is deprecated, should refactor in near term
	rand.Seed(time.Now().UnixNano())

	var state = randSeq(16)
	var scope = "FILL IN SCOPE STR HERE"

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	exec.Command(cmd, args...).Start()
	return nil
}

// Helper used to generate random state value for Spotify auth req
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
