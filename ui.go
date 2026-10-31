package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// BUBBLE TEA: VIEW
func (m model) View() string {
	var content string

	switch m.currentScreen {

	case menuScreen:
		content = m.menuView()

	case inputScreen:
		content = m.inputView()

	case listScreen:
		content = m.listView()

	case resultScreen:
		content = m.resultView()

	case targetSelectScreen:
		content = m.targetSelectView()
	}

	frame := borderStyle
	if m.width > 0 {
		frame = frame.Width(max(1, m.width-6))
	}

	return "\n" + frame.Render(content) + "\n"
}

// MENU VIEW
func (m model) menuView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("MULTI-PING MONITOR"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Service Health Monitoring CLI"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("TARGETS"))
	b.WriteString("\n")

	menu := []string{
		"List Targets",
		"Add Target",
		"Edit Target",
		"Delete Target",
		"Run Health Check",
		"Exit",
	}

	for i, item := range menu {
		line := fmt.Sprintf("%d. %s", i+1, item)

		if m.cursor == i {
			b.WriteString(selectedStyle.Render("❯ " + line))
		} else {
			b.WriteString(normalStyle.Render("  " + line))
		}

		b.WriteString("\n")
	}

	b.WriteString("\n")

	if m.message != "" {
		b.WriteString(messageStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	b.WriteString(subtitleStyle.Render(
		"↑/↓ Navigate • Enter Select • 1-6 Quick Select • Q Quit",
	))

	return b.String()
}

// INPUT VIEW
func (m model) inputView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("TARGET CONFIGURATION"))
	b.WriteString("\n\n")

	prompt := ""

	switch m.action {

	case addName:
		prompt = "Service Name"

	case addURL:
		prompt = "Host or IP"

	case editName:
		prompt = "New Service Name (blank to keep current)"

	case editURL:
		prompt = "New Host or IP (blank to keep current)"

	}

	b.WriteString(labelStyle.Render(prompt))
	b.WriteString("\n\n")

	b.WriteString(inputStyle.Render("> " + m.input + "▌"))
	b.WriteString("\n\n")

	if m.message != "" {
		b.WriteString(messageStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	b.WriteString(subtitleStyle.Render(
		"Enter Submit  •  Esc Cancel",
	))

	return b.String()
}

func (m model) targetSelectView() string {
	var b strings.Builder
	title := "SELECT TARGET TO EDIT"
	if m.action == deleteID {
		title = "SELECT TARGET TO DELETE"
	}

	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	for i, target := range m.targets {
		line := fmt.Sprintf("[%d] %s — %s", target.ID, target.Name, target.URL)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("❯ " + line))
		} else {
			b.WriteString(normalStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("↑/↓ Navigate • Enter Select • Esc Back"))
	return b.String()
}

// LIST VIEW
func (m model) listView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("TARGET LIST"))
	b.WriteString("\n\n")

	if len(m.targets) == 0 {
		b.WriteString(subtitleStyle.Render("No targets configured."))
	} else {
		for _, target := range m.targets {
			b.WriteString(labelStyle.Render(fmt.Sprintf("[%d] ", target.ID)))
			b.WriteString(lipgloss.NewStyle().Bold(true).Render(target.Name))
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(wrapText(subtitleStyle.Render(target.URL), m.contentWidth()))
			b.WriteString("\n\n")
		}
	}

	b.WriteString(subtitleStyle.Render(
		"Esc / Enter Back  •  Q Quit",
	))

	return b.String()
}

// RESULT VIEW
func (m model) resultView() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("HEALTH CHECK"))
	b.WriteString("\n\n")

	if m.checking {
		b.WriteString(warningStyle.Render("● Monitoring targets..."))
		b.WriteString("\n\n")

		b.WriteString(subtitleStyle.Render(
			"Checks run every 5 seconds. Press Enter or Esc to stop.",
		))
		b.WriteString("\n\n")

		if len(m.results) == 0 {
			return b.String()
		}
	}

	up, down, warn := 0, 0, 0

	for _, result := range m.results {
		var status string

		switch result.Status {

		case "UP":
			up++
			status = successStyle.Render("[UP]")

		case "DOWN":
			down++
			status = errorStyle.Render("[DOWN]")

		case "WARN":
			warn++
			status = warningStyle.Render("[WARN]")
		}

		details := fmt.Sprintf("%s", result.Latency.Round(time.Millisecond))

		if result.Error != "" {
			details = result.Error
		}

		b.WriteString(fmt.Sprintf("%s %s\n", status, result.Target.Name))
		b.WriteString("  ")
		b.WriteString(wrapText(subtitleStyle.Render(details), m.contentWidth()))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render(
		"────────────────────────────────────────",
	))
	b.WriteString("\n\n")

	b.WriteString(
		successStyle.Render(fmt.Sprintf("UP: %d  ", up)),
	)

	b.WriteString(
		errorStyle.Render(fmt.Sprintf("DOWN: %d  ", down)),
	)

	b.WriteString(
		warningStyle.Render(fmt.Sprintf("WARN: %d", warn)),
	)

	b.WriteString("\n\n")

	b.WriteString(subtitleStyle.Render(
		fmt.Sprintf(
			"Last check completed in %s",
			m.checkTime.Round(time.Millisecond),
		),
	))

	b.WriteString("\n\n")

	b.WriteString(subtitleStyle.Render(
		"Enter / Esc Stop  •  Q Quit",
	))

	return b.String()
}

func (m model) contentWidth() int {
	if m.width <= 0 {
		return 60
	}
	return max(20, m.width-10)
}

func wrapText(value string, width int) string {
	return lipgloss.NewStyle().Width(width).Render(value)
}
