package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rmvaldesd/tplm/internal/config"
	"github.com/rmvaldesd/tplm/internal/tmux"
)

type mode int

const (
	modeNormal mode = iota
	modeConfirmKill
	modeRename
)

// pickerItem represents one row in the picker — either a project or a session.
type pickerItem struct {
	isSession bool
	name      string
	path      string // project path (projects only)
	windows   int    // window count (sessions only)
	attached  bool   // whether client is attached (sessions only)
}

// PickerModel is the Bubbletea model for the two-section picker.
type PickerModel struct {
	cfg       *config.Config
	projects  []pickerItem
	sessions  []pickerItem
	cursor    int // index into the combined list (projects then sessions)
	mode      mode
	rename    RenameModel
	err       error
	quitting  bool
	width     int
	height    int
}

// switchMsg tells the program to switch to a session and quit.
type switchMsg struct{ name string }

// NewPicker creates a new picker model.
func NewPicker(cfg *config.Config) PickerModel {
	m := PickerModel{cfg: cfg}
	m.refreshItems()
	return m
}

func (m *PickerModel) refreshItems() {
	m.projects = nil
	for _, p := range m.cfg.Projects {
		m.projects = append(m.projects, pickerItem{
			name: p.Name,
			path: p.Path,
		})
	}

	m.sessions = nil
	sessions, _ := tmux.ListSessions()
	for _, s := range sessions {
		m.sessions = append(m.sessions, pickerItem{
			isSession: true,
			name:      s.Name,
			windows:   s.Windows,
			attached:  s.Attached,
		})
	}
}

func (m PickerModel) totalItems() int {
	return len(m.projects) + len(m.sessions)
}

func (m PickerModel) selectedItem() *pickerItem {
	total := m.totalItems()
	if total == 0 || m.cursor < 0 || m.cursor >= total {
		return nil
	}
	if m.cursor < len(m.projects) {
		return &m.projects[m.cursor]
	}
	return &m.sessions[m.cursor-len(m.projects)]
}

func (m PickerModel) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case switchMsg:
		// Perform the switch and exit.
		m.quitting = true
		if err := tmux.SwitchClient(msg.name); err != nil {
			m.err = err
		}
		return m, tea.Quit

	case renameMsg:
		if err := tmux.RenameSession(msg.oldName, msg.newName); err != nil {
			m.err = err
		}
		m.mode = modeNormal
		m.refreshItems()
		return m, nil

	case renameCancelMsg:
		m.mode = modeNormal
		return m, nil
	}

	// Delegate to sub-modes.
	switch m.mode {
	case modeConfirmKill:
		return m.updateConfirmKill(msg)
	case modeRename:
		return m.updateRename(msg)
	default:
		return m.updateNormal(msg)
	}
}

func (m PickerModel) updateNormal(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}

		case key.Matches(msg, keys.Down):
			if m.cursor < m.totalItems()-1 {
				m.cursor++
			}

		case key.Matches(msg, keys.Select):
			item := m.selectedItem()
			if item == nil {
				break
			}
			if item.isSession {
				return m, func() tea.Msg { return switchMsg{name: item.name} }
			}
			// It's a project — create session if needed, then switch.
			proj := m.cfg.FindProject(item.name)
			if proj == nil {
				break
			}
			if tmux.SessionExists(proj.Name) {
				return m, func() tea.Msg { return switchMsg{name: proj.Name} }
			}
			// Create session from layout.
			if err := m.createSession(proj); err != nil {
				m.err = err
				break
			}
			return m, func() tea.Msg { return switchMsg{name: proj.Name} }

		case key.Matches(msg, keys.Kill):
			item := m.selectedItem()
			if item != nil && item.isSession {
				m.mode = modeConfirmKill
			}

		case key.Matches(msg, keys.Rename):
			item := m.selectedItem()
			if item != nil && item.isSession {
				m.rename = NewRenameModel(item.name)
				m.mode = modeRename
			}
		}
	}

	return m, nil
}

func (m PickerModel) updateConfirmKill(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Confirm):
			item := m.selectedItem()
			if item != nil && item.isSession {
				if err := tmux.KillSession(item.name); err != nil {
					m.err = err
				}
				m.refreshItems()
				if m.cursor >= m.totalItems() && m.cursor > 0 {
					m.cursor--
				}
			}
			m.mode = modeNormal
		case key.Matches(msg, keys.Cancel):
			m.mode = modeNormal
		}
	}
	return m, nil
}

func (m PickerModel) updateRename(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.rename, cmd = m.rename.Update(msg)
	return m, cmd
}

func (m *PickerModel) createSession(proj *config.Project) error {
	if err := tmux.NewSession(proj.Name, proj.Path); err != nil {
		return err
	}

	layout := m.cfg.GetLayout(proj)
	if err := tmux.ApplyLayout(proj.Name, layout, proj.Path); err != nil {
		return err
	}

	if len(proj.OnStart) > 0 {
		if err := tmux.RunOnStart(proj.Name, layout, proj.OnStart); err != nil {
			return err
		}
	}
	return nil
}

func (m PickerModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder
	w := m.width
	if w == 0 {
		w = 60
	}

	// Title bar.
	title := titleStyle.Render("tplm")
	hint := pathStyle.Render("d:kill  r:rename")
	gap := w - lipgloss.Width(title) - lipgloss.Width(hint)
	if gap < 1 {
		gap = 1
	}
	b.WriteString(title + strings.Repeat(" ", gap) + hint + "\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")

	// Projects section.
	b.WriteString(headerStyle.Render("Projects") + "\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")

	for i, item := range m.projects {
		b.WriteString(m.renderItem(i, item, w))
	}

	if len(m.projects) == 0 {
		b.WriteString(normalStyle.Render("(no projects configured)") + "\n")
	}

	b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")

	// Sessions section.
	b.WriteString(headerStyle.Render("Active Sessions") + "\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")

	for i, item := range m.sessions {
		idx := len(m.projects) + i
		b.WriteString(m.renderItem(idx, item, w))
	}

	if len(m.sessions) == 0 {
		b.WriteString(normalStyle.Render("(no active sessions)") + "\n")
	}

	// Mode-specific footer.
	switch m.mode {
	case modeConfirmKill:
		item := m.selectedItem()
		if item != nil {
			b.WriteString("\n")
			b.WriteString(confirmStyle.Render(fmt.Sprintf("  Kill session %q? (y/n)", item.name)) + "\n")
		}
	case modeRename:
		b.WriteString("\n")
		b.WriteString(m.rename.View() + "\n")
	default:
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑↓ navigate  ⏎ open/switch  q quit") + "\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(confirmStyle.Render(fmt.Sprintf("  Error: %v", m.err)) + "\n")
	}

	return b.String()
}

func (m PickerModel) renderItem(idx int, item pickerItem, width int) string {
	cursor := "  "
	style := normalStyle
	if idx == m.cursor {
		cursor = "> "
		style = selectedStyle
	}

	if item.isSession {
		indicator := activeIndicator.Render()
		name := style.Render(item.name)
		info := pathStyle.Render(fmt.Sprintf("%d windows", item.windows))
		return fmt.Sprintf(" %s%s %s  %s\n", cursor, indicator, name, info)
	}

	name := style.Render(item.name)
	path := pathStyle.Render(item.path)
	return fmt.Sprintf(" %s%s  %s\n", cursor, name, path)
}
