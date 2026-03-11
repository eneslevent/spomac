package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/eneslevent/spomac/internal/spotify"
)

// Colors
var (
	green   = lipgloss.Color("#1DB954")
	white   = lipgloss.Color("#FFFFFF")
	gray    = lipgloss.Color("#535353")
	dimGray = lipgloss.Color("#404040")
)

// Styles
var (
	trackStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white)

	artistStyle = lipgloss.NewStyle().
			Foreground(green)

	albumStyle = lipgloss.NewStyle().
			Foreground(gray)

	timeStyle = lipgloss.NewStyle().
			Foreground(gray)

	barFull = lipgloss.NewStyle().
		Foreground(green)

	barEmpty = lipgloss.NewStyle().
			Foreground(dimGray)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(green).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)
)

type tickMsg time.Time

type Model struct {
	state   spotify.PlayerState
	err     error
	width   int
	volume  int
	quitting bool
}

func NewModel() Model {
	return Model{width: 50}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return doTick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		if m.width < 30 {
			m.width = 30
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case " ", "space", "2":
			_ = spotify.PlayPause()
			return m, doTick()

		case "1", "left":
			_ = spotify.Previous()
			return m, doTick()

		case "3", "right":
			_ = spotify.Next()
			return m, doTick()

		case "up":
			newVol := m.volume + 5
			if newVol > 100 {
				newVol = 100
			}
			_ = spotify.SetVolume(newVol)
			m.volume = newVol
			return m, nil

		case "down":
			newVol := m.volume - 5
			if newVol < 0 {
				newVol = 0
			}
			_ = spotify.SetVolume(newVol)
			m.volume = newVol
			return m, nil
		}

	case tickMsg:
		if !spotify.IsRunning() {
			_ = spotify.Launch()
			m.err = fmt.Errorf("launching Spotify, please wait")
			return m, doTick()
		}
		state, err := spotify.GetState()
		m.state = state
		m.err = err
		vol, verr := spotify.GetVolume()
		if verr == nil {
			m.volume = vol
		}
		return m, doTick()
	}

	return m, nil
}

func (m Model) View() tea.View {
	if m.quitting {
		v := tea.NewView("")
		v.AltScreen = true
		return v
	}

	if m.err != nil {
		content := errorStyle.Render("Spotify is not running or no track is playing")
		v := tea.NewView("\n" + boxStyle.Width(m.width).Render(content) + "\n")
		v.AltScreen = true
		v.WindowTitle = "spomac"
		return v
	}

	// Status icon
	status := artistStyle.Render("▶")
	if !m.state.IsPlaying {
		status = timeStyle.Render("⏸")
	}

	// Artist - Album - Track on single line
	info := artistStyle.Render(m.state.Artist) + albumStyle.Render(" - "+m.state.Album+" - ") + trackStyle.Render(m.state.Track)

	// Progress bar
	barWidth := m.width - 6
	if barWidth < 10 {
		barWidth = 10
	}

	var percent float64
	if m.state.Duration > 0 {
		percent = m.state.Position / m.state.Duration
	}
	if percent > 1 {
		percent = 1
	}

	filled := int(float64(barWidth) * percent)
	empty := barWidth - filled

	bar := barFull.Render(strings.Repeat("█", filled)) +
		barEmpty.Render(strings.Repeat("░", empty))

	posMin := int(m.state.Position) / 60
	posSec := int(m.state.Position) % 60
	durMin := int(m.state.Duration) / 60
	durSec := int(m.state.Duration) % 60

	// Volume indicator
	volText := timeStyle.Render(fmt.Sprintf("♪ %d%%", m.volume))

	// Right-align volume on the time line
	timeStr := fmt.Sprintf("%d:%02d / %d:%02d", posMin, posSec, durMin, durSec)
	volStr := fmt.Sprintf("♪ %d%%", m.volume)
	innerWidth := m.width - 2
	gap := innerWidth - len(timeStr) - len(volStr)
	if gap < 1 {
		gap = 1
	}
	bottomLine := timeStyle.Render(timeStr) + strings.Repeat(" ", gap) + volText

	// Compose
	content := fmt.Sprintf(
		"%s %s\n%s\n%s",
		status, info,
		bar,
		bottomLine,
	)

	v := tea.NewView("\n" + boxStyle.Width(m.width).Render(content) + "\n")
	v.AltScreen = true
	v.WindowTitle = "spomac"
	return v
}

