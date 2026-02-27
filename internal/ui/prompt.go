package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// SelectItem represents an item in a selection list.
type SelectItem struct {
	ID    string
	Label string
}

// ─── Single Select ────────────────────────────────────────────────────────────

type singleSelectModel struct {
	title      string
	allItems   []SelectItem
	filtered   []SelectItem
	cursor     int
	chosen     string
	done       bool
	cancelled  bool
	textInput  textinput.Model
}

var (
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	titleStyle     = lipgloss.NewStyle().Bold(true)
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func initialSingleSelect(title string, items []SelectItem) singleSelectModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40

	return singleSelectModel{
		title:     title,
		allItems:  items,
		filtered:  items,
		textInput: ti,
	}
}

func (m singleSelectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m singleSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filtered) > 0 {
				m.chosen = m.filtered[m.cursor].ID
				m.done = true
				return m, tea.Quit
			}
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			m.filterItems()
			m.cursor = 0
			return m, cmd
		}
	}
	return m, nil
}

func (m *singleSelectModel) filterItems() {
	query := strings.TrimSpace(m.textInput.Value())
	if query == "" {
		m.filtered = m.allItems
		return
	}
	labels := make([]string, len(m.allItems))
	for i, item := range m.allItems {
		labels[i] = item.Label
	}
	matches := fuzzy.FindNormalizedFold(query, labels)
	matchSet := make(map[string]bool, len(matches))
	for _, m := range matches {
		matchSet[m] = true
	}
	var filtered []SelectItem
	for _, item := range m.allItems {
		if matchSet[item.Label] {
			filtered = append(filtered, item)
		}
	}
	m.filtered = filtered
}

func (m singleSelectModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render(m.title) + "\n")
	sb.WriteString(m.textInput.View() + "\n\n")

	maxDisplay := 15
	start := 0
	if m.cursor >= maxDisplay {
		start = m.cursor - maxDisplay + 1
	}
	end := start + maxDisplay
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		item := m.filtered[i]
		if i == m.cursor {
			sb.WriteString(cursorStyle.Render("▶ ") + selectedStyle.Render(item.Label) + "\n")
		} else {
			sb.WriteString("  " + item.Label + "\n")
		}
	}

	if len(m.filtered) == 0 {
		sb.WriteString(dimStyle.Render("  No matches found\n"))
	}

	sb.WriteString(dimStyle.Render("\n↑/↓ navigate • enter select • esc cancel"))
	return sb.String()
}

// SelectOne shows an interactive single-select list. Returns the chosen item's ID.
func SelectOne(title string, items []SelectItem) (string, error) {
	m := initialSingleSelect(title, items)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	final := result.(singleSelectModel)
	if final.cancelled {
		return "", fmt.Errorf("selection cancelled")
	}
	return final.chosen, nil
}

// ─── Multi Select ─────────────────────────────────────────────────────────────

type multiSelectModel struct {
	title     string
	allItems  []SelectItem
	filtered  []SelectItem
	cursor    int
	selected  map[string]bool
	done      bool
	cancelled bool
	textInput textinput.Model
}

func initialMultiSelect(title string, items []SelectItem) multiSelectModel {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40

	return multiSelectModel{
		title:     title,
		allItems:  items,
		filtered:  items,
		selected:  make(map[string]bool),
		textInput: ti,
	}
}

func (m multiSelectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "up", "ctrl+p":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case " ":
			if len(m.filtered) > 0 {
				id := m.filtered[m.cursor].ID
				m.selected[id] = !m.selected[id]
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			m.filterItems()
			m.cursor = 0
			return m, cmd
		}
	}
	return m, nil
}

func (m *multiSelectModel) filterItems() {
	query := strings.TrimSpace(m.textInput.Value())
	if query == "" {
		m.filtered = m.allItems
		return
	}
	labels := make([]string, len(m.allItems))
	for i, item := range m.allItems {
		labels[i] = item.Label
	}
	matches := fuzzy.FindNormalizedFold(query, labels)
	matchSet := make(map[string]bool)
	for _, m := range matches {
		matchSet[m] = true
	}
	var filtered []SelectItem
	for _, item := range m.allItems {
		if matchSet[item.Label] {
			filtered = append(filtered, item)
		}
	}
	m.filtered = filtered
}

func (m multiSelectModel) View() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render(m.title) + "\n")
	sb.WriteString(m.textInput.View() + "\n\n")

	maxDisplay := 15
	start := 0
	if m.cursor >= maxDisplay {
		start = m.cursor - maxDisplay + 1
	}
	end := start + maxDisplay
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		item := m.filtered[i]
		check := "[ ]"
		if m.selected[item.ID] {
			check = selectedStyle.Render("[✔]")
		}

		if i == m.cursor {
			sb.WriteString(cursorStyle.Render("▶ ") + check + " " + selectedStyle.Render(item.Label) + "\n")
		} else {
			sb.WriteString("  " + check + " " + item.Label + "\n")
		}
	}

	if len(m.filtered) == 0 {
		sb.WriteString(dimStyle.Render("  No matches found\n"))
	}

	// Show selected count
	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	sb.WriteString(dimStyle.Render(fmt.Sprintf("\n↑/↓ navigate • space toggle • enter confirm (%d selected) • esc cancel", count)))
	return sb.String()
}

// MultiSelect shows an interactive multi-select list. Returns selected item IDs.
func MultiSelect(title string, items []SelectItem) ([]string, error) {
	m := initialMultiSelect(title, items)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}
	final := result.(multiSelectModel)
	if final.cancelled {
		return nil, fmt.Errorf("selection cancelled")
	}

	var chosen []string
	for _, item := range final.allItems {
		if final.selected[item.ID] {
			chosen = append(chosen, item.ID)
		}
	}
	return chosen, nil
}
