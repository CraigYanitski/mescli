package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// add struct functions to implement help.KeyMap

// ShortHelp returns bindings to show in the abbreviated help view. It's part
// of the help.KeyMap interface.
func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.viewport.KeyMap.Up,
		m.viewport.KeyMap.Down,
	}

	kb = append(kb,
		m.viewport.KeyMap.HalfPageUp,
		m.viewport.KeyMap.HalfPageDown,
	)

	return append(kb,
		m.viewport.KeyMap.PageUp,
		m.viewport.KeyMap.PageDown,
	)
}

// FullHelp returns bindings to show the full help view. It's part of the
// help.KeyMap interface.
func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.viewport.KeyMap.Up,
		m.viewport.KeyMap.Down,
        m.viewport.KeyMap.HalfPageUp,
		m.viewport.KeyMap.HalfPageDown,
        m.viewport.KeyMap.PageUp,
        m.viewport.KeyMap.PageDown,
	}}

    taKB := [][]key.Binding {{
        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "return to previous screen")),
        key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "show this help screen")),
        m.textarea.KeyMap.Paste,
        m.textarea.KeyMap.InsertNewline,
        m.textarea.KeyMap.CharacterForward,
        m.textarea.KeyMap.CharacterBackward,
    }}

	return append(kb, taKB...)
}

func updateHelp(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            m.Quitting = true
            m.Chosen = 3
            m.viewHelp = false
            return m, tea.Quit
        case tea.KeyEsc, tea.KeyCtrlQ:
            m.viewHelp = false
        }
    }
    return m, nil
}

func helpView(m Model) string {
    m.help.ShowAll = true
    help := helpStyle.Width(m.viewport.Width - 6).Margin(optionMargin.height, optionMargin.width).
        Render(m.help.View(m))
    return fmt.Sprintf(helpWrapping, help)
}

