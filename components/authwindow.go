package components

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"mikehaus/spoofify/helpers"
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
	loggedIn    bool
	list        list.Model
	spotifyAuth *helpers.SpotifyAuth
}

// I need to expose a single auth
type AuthWindow struct {
	SpotifyAuth *helpers.SpotifyAuth
}

// MARK: External exports
func NewAuthWindow(auth *helpers.SpotifyAuth) *AuthWindow {
	return &AuthWindow{
		SpotifyAuth: auth,
	}
}

func (w *AuthWindow) Render() {
	m := initialModel(w.SpotifyAuth)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// MARK: View
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			handleQuit(m)
		}
		if msg.String() == "enter" {
			handleSelection(m)
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

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel(s *helpers.SpotifyAuth) model {
	items := []list.Item{
		item{title: "Login", desc: "Login with your Spotify login info"},
		item{title: "Quit", desc: "Quit application"},
	}

	m := model{
		loggedIn:    false,
		list:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		spotifyAuth: s,
	}

	m.list.Title = "Please select an option"

	return m
}

func handleSelection(m model) (tea.Model, tea.Cmd) {
	if m.list.Index() == 0 {
		return m, authenticateSpotifyInBrowser(m.spotifyAuth)
	}

	return handleQuit(m)
}

func handleQuit(m model) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// Opens default browser to spotify to log in to spotify
func authenticateSpotifyInBrowser(auth *helpers.SpotifyAuth) tea.Cmd {
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

	// TODO: This isn't opening a url so need to figure that out
	// client, url := helpers.GenerateSpotifyOAuthClient()
	url := auth.SpotifyOAuthUrl()

	args = append(args, url)
	exec.Command(cmd, args...).Start()
	return nil
}
