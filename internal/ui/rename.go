package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// renameMsg is sent when the user confirms a rename.
type renameMsg struct {
	oldName string
	newName string
}

// renameCancelMsg is sent when the user cancels a rename.
type renameCancelMsg struct{}

// RenameModel is an inline text input for renaming a session.
type RenameModel struct {
	input   textinput.Model
	oldName string
}

// NewRenameModel creates a rename input pre-filled with the current name.
func NewRenameModel(currentName string) RenameModel {
	ti := textinput.New()
	ti.SetValue(currentName)
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40
	ti.Prompt = "Rename: "
	ti.PromptStyle = inputPromptStyle

	return RenameModel{
		input:   ti,
		oldName: currentName,
	}
}

func (m RenameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RenameModel) Update(msg tea.Msg) (RenameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			newName := m.input.Value()
			if newName != "" && newName != m.oldName {
				return m, func() tea.Msg {
					return renameMsg{oldName: m.oldName, newName: newName}
				}
			}
			return m, func() tea.Msg { return renameCancelMsg{} }
		case "esc":
			return m, func() tea.Msg { return renameCancelMsg{} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m RenameModel) View() string {
	return m.input.View()
}
