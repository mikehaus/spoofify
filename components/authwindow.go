package components 

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

  "mikehaus/spoofify/helpers"
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

func handleSelection(m model) (tea.Model, tea.Cmd) {
	if m.list.Index() == 0 {
		return m, authWithSpotify()
	}

	return handleQuit(m)
}

func handleQuit(m model) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// Create a get request with spotify
func authWithSpotify() tea.Cmd {
	// TODO: integrate this with Oauth2 client 
  var CLIENT_ID = os.Getenv("SPOTIFY_CLIENT_ID")
  var state = helpers.GetUriAuthState()

  return authenticateSpotifyInBrowser()
}

// Opens default browser to spotify to log in to spotify 
func authenticateSpotifyInBrowser() tea.Cmd {
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
	args = append(args, helpers.GetAuthUrl())
	exec.Command(cmd, args...).Start()
	return nil
}
