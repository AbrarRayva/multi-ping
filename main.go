package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DATA STRUCTURES
type Target struct {
	ID   int
	Name string
	URL  string
}

type CheckResult struct {
	Target  Target
	Status  string
	Code    int
	Latency time.Duration
	Error   string
}

type screen int

const (
	menuScreen screen = iota
	inputScreen
	listScreen
	resultScreen
	targetSelectScreen
)

type formAction int

const (
	addName formAction = iota
	addURL
	editID
	editName
	editURL
	deleteID
)

type healthCheckMsg struct {
	results  []CheckResult
	duration time.Duration
}

type healthCheckTickMsg struct{}

// STYLES
var (
	colorPrimary = lipgloss.Color("#6d2fff")
	colorGreen   = lipgloss.Color("#22C55E")
	colorRed     = lipgloss.Color("#EF4444")
	colorYellow  = lipgloss.Color("#EAB308")
	colorMuted   = lipgloss.Color("#888888")
	colorWhite   = lipgloss.Color("#FFFFFF")
	colorGray    = lipgloss.Color("#AAAAAA")
	colorBorder  = lipgloss.Color("#444444")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Background(lipgloss.Color("#444444")).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)

	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)

	labelStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	inputStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	messageStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)
)

// MODEL
type model struct {
	targets []Target
	nextID  int

	currentScreen screen
	cursor        int

	action      formAction
	input       string
	pendingName string
	editingID   int

	message string

	results   []CheckResult
	checking  bool
	checkTime time.Duration
	width     int
}

func initialModel() model {
	return model{
		targets:       []Target{},
		nextID:        1,
		currentScreen: menuScreen,
		cursor:        0,
	}
}

// Bubble Tea: Init
func (m model) Init() tea.Cmd {
	return nil
}

// BUBBLE TEA: UPDATE
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case healthCheckMsg:
		m.results = msg.results
		m.checkTime = msg.duration
		if m.checking {
			return m, tea.Tick(healthCheckInterval, func(time.Time) tea.Msg {
				return healthCheckTickMsg{}
			})
		}
		return m, nil

	case healthCheckTickMsg:
		if m.checking {
			return m, runHealthCheck(m.targets)
		}
		return m, nil

	case tea.KeyMsg:

		// MENU
		if m.currentScreen == menuScreen {
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "k":
				m.cursor = (m.cursor - 1 + 6) % 6

			case "down", "j":
				m.cursor = (m.cursor + 1) % 6

			case "enter":
				return m.selectMenu()

			default:
				// Number keys immediately open their menu.
				if len(msg.Runes) == 1 {
					key := msg.Runes[0]

					if key >= '1' && key <= '6' {
						m.cursor = int(key - '1')
						return m.selectMenu()
					}
				}
			}

			return m, nil
		}

		// INPUT
		if m.currentScreen == inputScreen {
			switch msg.Type {

			case tea.KeyEsc:
				m.currentScreen = menuScreen
				m.input = ""
				m.message = ""

			case tea.KeyEnter:
				return m.submitInput()

			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.input) > 0 {
					_, size := utf8.DecodeLastRuneInString(m.input)
					m.input = m.input[:len(m.input)-size]
				}

			case tea.KeyRunes:
				m.input += string(msg.Runes)
			}

			return m, nil
		}

		if m.currentScreen == targetSelectScreen {
			switch msg.String() {
			case "up", "k":
				m.cursor = (m.cursor - 1 + len(m.targets)) % len(m.targets)
			case "down", "j":
				m.cursor = (m.cursor + 1) % len(m.targets)
			case "enter":
				return m.selectTarget()
			case "esc":
				m.currentScreen = menuScreen
				m.message = ""
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		// LIST / RESULTS
		switch msg.String() {

		case "esc", "enter":
			m.checking = false
			m.currentScreen = menuScreen
			m.message = ""

		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

// MENU HANDLER
func (m model) selectMenu() (tea.Model, tea.Cmd) {
	switch m.cursor {

	case 0:
		m.currentScreen = listScreen
		m.message = ""

	case 1:
		m.currentScreen = inputScreen
		m.action = addName
		m.input = ""
		m.message = ""

	case 2:
		m.action = editID
		if len(m.targets) == 0 {
			m.message = "No targets configured."
			return m, nil
		}
		m.currentScreen = targetSelectScreen
		m.cursor = 0
		m.message = ""

	case 3:
		m.action = deleteID
		if len(m.targets) == 0 {
			m.message = "No targets configured."
			return m, nil
		}
		m.currentScreen = targetSelectScreen
		m.cursor = 0
		m.message = ""

	case 4:
		if len(m.targets) == 0 {
			m.message = "No targets configured."
			return m, nil
		}

		m.currentScreen = resultScreen
		m.checking = true
		m.message = ""

		return m, runHealthCheck(m.targets)

	case 5:
		return m, tea.Quit
	}

	return m, nil
}

// INPUT HANDLER
func (m model) submitInput() (tea.Model, tea.Cmd) {
	value := strings.TrimSpace(m.input)

	switch m.action {

	// ADD TARGET
	case addName:
		if value == "" {
			m.message = "Service name cannot be empty."
			return m, nil
		}

		m.pendingName = value
		m.action = addURL
		m.input = ""
		m.message = ""

	case addURL:
		normalizedURL, valid := normalizeTarget(value)
		if !valid {
			m.message = "Invalid host or IP address."
			return m, nil
		}

		m.targets = append(m.targets, Target{
			ID:   m.nextID,
			Name: m.pendingName,
			URL:  normalizedURL,
		})

		m.nextID++

		m.currentScreen = menuScreen
		m.input = ""
		m.message = "Target added successfully."

	// EDIT TARGET
	case editID:
	case editName:
		if value != "" {
			m.pendingName = value
		}

		m.action = editURL
		m.input = ""
		m.message = ""

	case editURL:
		for i := range m.targets {
			if m.targets[i].ID != m.editingID {
				continue
			}

			if value != "" {
				normalizedURL, valid := normalizeTarget(value)
				if !valid {
					m.message = "Invalid host or IP address. Changes were not saved."
					return m, nil
				}

				m.targets[i].URL = normalizedURL
			}

			m.targets[i].Name = m.pendingName
			break
		}

		m.currentScreen = menuScreen
		m.input = ""
		m.message = "Target updated successfully."

	// DELETE TARGET
	case deleteID:
	}

	return m, nil
}

func (m model) selectTarget() (tea.Model, tea.Cmd) {
	target := m.targets[m.cursor]

	switch m.action {
	case editID:
		m.editingID = target.ID
		m.pendingName = target.Name
		m.action = editName
		m.currentScreen = inputScreen
		m.input = ""
		m.message = ""

	case deleteID:
		m.targets = append(m.targets[:m.cursor], m.targets[m.cursor+1:]...)
		m.currentScreen = menuScreen
		m.cursor = 0
		m.message = "Target deleted successfully."
	}

	return m, nil
}

// MAIN
func main() {
	program := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := program.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
