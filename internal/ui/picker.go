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

// pickerItem represents one row in the picker — a project, session, or window.
type pickerItem struct {
	isSession    bool
	isWindow     bool
	name         string
	path         string // project path (projects only)
	windows      int    // window count (sessions only)
	attached     bool   // whether client is attached (sessions only)
	expanded     bool   // whether session is expanded (sessions only)
	sessionName  string // parent session name (windows only)
	windowIndex  int    // tmux window index (windows only)
	windowActive bool   // active window indicator (windows only)
}

// PickerModel is the Bubbletea model for the two-section picker.
type PickerModel struct {
	cfg          *config.Config
	projects     []pickerItem
	sessions     []pickerItem
	displayItems []pickerItem // flattened list the cursor navigates
	expanded     map[string][]tmux.WindowInfo
	cursor       int // index into displayItems
	mode         mode
	rename       RenameModel
	err          error
	quitting     bool
	width        int
	height       int
}

// switchMsg tells the program to switch to a session and quit.
type switchMsg struct{ name string }

// NewPicker creates a new picker model.
func NewPicker(cfg *config.Config) PickerModel {
	m := PickerModel{
		cfg:      cfg,
		expanded: make(map[string][]tmux.WindowInfo),
	}
	m.refreshItems()

	// Auto-expand the current tmux session.
	if current, err := tmux.CurrentSession(); err == nil && current != "" {
		for i, item := range m.displayItems {
			if item.isSession && item.name == current {
				m.cursor = i
				_ = m.expandSession(&item)
				break
			}
		}
	}

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
		_, isExpanded := m.expanded[s.Name]
		m.sessions = append(m.sessions, pickerItem{
			isSession: true,
			name:      s.Name,
			windows:   s.Windows,
			attached:  s.Attached,
			expanded:  isExpanded,
		})
	}

	// Prune stale expanded entries.
	activeNames := make(map[string]bool)
	for _, s := range m.sessions {
		activeNames[s.name] = true
	}
	for name := range m.expanded {
		if !activeNames[name] {
			delete(m.expanded, name)
		}
	}

	m.rebuildDisplayItems()
}

func (m *PickerModel) rebuildDisplayItems() {
	m.displayItems = nil

	for _, item := range m.projects {
		m.displayItems = append(m.displayItems, item)
	}

	for _, item := range m.sessions {
		m.displayItems = append(m.displayItems, item)
		if item.expanded {
			if wins, ok := m.expanded[item.name]; ok {
				for _, w := range wins {
					m.displayItems = append(m.displayItems, pickerItem{
						isWindow:     true,
						name:         w.Name,
						sessionName:  item.name,
						windowIndex:  w.Index,
						windowActive: w.Active,
					})
				}
			}
		}
	}
}

func (m PickerModel) totalItems() int {
	return len(m.displayItems)
}

func (m PickerModel) selectedItem() *pickerItem {
	if len(m.displayItems) == 0 || m.cursor < 0 || m.cursor >= len(m.displayItems) {
		return nil
	}
	return &m.displayItems[m.cursor]
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
		// Transfer expanded state from old name to new name.
		if wins, ok := m.expanded[msg.oldName]; ok {
			delete(m.expanded, msg.oldName)
			m.expanded[msg.newName] = wins
		}
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

			if item.isWindow {
				// Switch to the specific window.
				target := fmt.Sprintf("%s:%d", item.sessionName, item.windowIndex)
				return m, func() tea.Msg { return switchMsg{name: target} }
			}

			if item.isSession {
				// Toggle expand/collapse.
				if item.expanded {
					m.collapseSession(item)
				} else {
					if err := m.expandSession(item); err != nil {
						m.err = err
					}
				}
				break
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

		case key.Matches(msg, keys.Right):
			item := m.selectedItem()
			if item == nil {
				break
			}

			if item.isSession {
				if item.expanded {
					// Move cursor to first child window.
					if m.cursor+1 < m.totalItems() && m.displayItems[m.cursor+1].isWindow {
						m.cursor++
					}
				} else {
					// Expand the session.
					if err := m.expandSession(item); err != nil {
						m.err = err
					}
				}
				break
			}

			if item.isWindow {
				// Switch to the window (same as Enter).
				target := fmt.Sprintf("%s:%d", item.sessionName, item.windowIndex)
				return m, func() tea.Msg { return switchMsg{name: target} }
			}

			// Project — open/switch (same as Enter).
			proj := m.cfg.FindProject(item.name)
			if proj == nil {
				break
			}
			if tmux.SessionExists(proj.Name) {
				return m, func() tea.Msg { return switchMsg{name: proj.Name} }
			}
			if err := m.createSession(proj); err != nil {
				m.err = err
				break
			}
			return m, func() tea.Msg { return switchMsg{name: proj.Name} }

		case key.Matches(msg, keys.Left):
			item := m.selectedItem()
			if item == nil {
				break
			}

			if item.isSession && item.expanded {
				m.collapseSession(item)
				break
			}

			if item.isWindow {
				// Jump to parent session.
				m.cursor = m.findParentSessionIndex()
			}

		case key.Matches(msg, keys.Kill):
			item := m.selectedItem()
			if item != nil && (item.isSession || item.isWindow) {
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
			} else if item != nil && item.isWindow {
				target := fmt.Sprintf("%s:%d", item.sessionName, item.windowIndex)
				if err := tmux.KillWindow(target); err != nil {
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

// expandSession expands a session to show its windows.
func (m *PickerModel) expandSession(item *pickerItem) error {
	wins, err := tmux.ListWindows(item.name)
	if err != nil {
		return err
	}
	m.expanded[item.name] = wins
	for i := range m.sessions {
		if m.sessions[i].name == item.name {
			m.sessions[i].expanded = true
			break
		}
	}
	m.rebuildDisplayItems()
	return nil
}

// collapseSession collapses a session, hiding its windows.
func (m *PickerModel) collapseSession(item *pickerItem) {
	sessionIdx := m.cursor
	delete(m.expanded, item.name)
	for i := range m.sessions {
		if m.sessions[i].name == item.name {
			m.sessions[i].expanded = false
			break
		}
	}
	m.rebuildDisplayItems()
	if m.cursor > sessionIdx {
		m.cursor = sessionIdx
	}
	if m.cursor >= m.totalItems() && m.cursor > 0 {
		m.cursor = m.totalItems() - 1
	}
}

// findParentSessionIndex scans backwards from the current cursor to find the parent session.
func (m *PickerModel) findParentSessionIndex() int {
	for i := m.cursor - 1; i >= 0; i-- {
		if m.displayItems[i].isSession {
			return i
		}
	}
	return m.cursor
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

	// Render using displayItems with section headers.
	inSessions := false
	projectsRendered := false

	for i, item := range m.displayItems {
		// Insert Projects header before first project.
		if !projectsRendered && !item.isSession && !item.isWindow {
			b.WriteString(headerStyle.Render("Projects") + "\n")
			b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
			projectsRendered = true
		}

		// Insert Sessions header at the boundary.
		if !inSessions && (item.isSession || item.isWindow) {
			if !projectsRendered {
				b.WriteString(headerStyle.Render("Projects") + "\n")
				b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
				b.WriteString(normalStyle.Render("(no projects configured)") + "\n")
				projectsRendered = true
			}
			b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
			b.WriteString(headerStyle.Render("Active Sessions") + "\n")
			b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
			inSessions = true
		}

		b.WriteString(m.renderItem(i, item, w))
	}

	// Handle empty states.
	if !projectsRendered {
		b.WriteString(headerStyle.Render("Projects") + "\n")
		b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
		b.WriteString(normalStyle.Render("(no projects configured)") + "\n")
	}

	if !inSessions {
		b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
		b.WriteString(headerStyle.Render("Active Sessions") + "\n")
		b.WriteString(separatorStyle.Render(strings.Repeat("─", w)) + "\n")
		b.WriteString(normalStyle.Render("(no active sessions)") + "\n")
	}

	// Mode-specific footer.
	switch m.mode {
	case modeConfirmKill:
		item := m.selectedItem()
		if item != nil {
			b.WriteString("\n")
			kind := "session"
			if item.isWindow {
				kind = "window"
			}
			b.WriteString(confirmStyle.Render(fmt.Sprintf("  Kill %s %q? (y/n)", kind, item.name)) + "\n")
		}
	case modeRename:
		b.WriteString("\n")
		b.WriteString(m.rename.View() + "\n")
	default:
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("hjkl navigate  ⏎ select  d kill  r rename  q quit") + "\n")
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

	if item.isWindow {
		indicator := "  "
		if item.windowActive {
			indicator = windowActiveIndicator.Render() + " "
		}
		name := style.Render(item.name)
		return fmt.Sprintf("     %s%s%s\n", cursor, indicator, name)
	}

	if item.isSession {
		chevron := "▶"
		if item.expanded {
			chevron = "▼"
		}
		indicator := activeIndicator.Render()
		name := style.Render(item.name)
		info := ""
		if !item.expanded {
			info = "  " + pathStyle.Render(fmt.Sprintf("%d windows", item.windows))
		}
		return fmt.Sprintf(" %s%s %s %s%s\n", cursor, indicator, chevron, name, info)
	}

	name := style.Render(item.name)
	path := pathStyle.Render(item.path)
	return fmt.Sprintf(" %s%s  %s\n", cursor, name, path)
}
