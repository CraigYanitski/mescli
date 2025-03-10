package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keys.Quit):
            m.Quitting = true
            return m, tea.Quit
        case key.Matches(msg, m.keys.Back):
            m.loggedIn = false
            return m, nil
        case key.Matches(msg, m.keys.Enter):
            o, _ := m.options.SelectedItem().(option)
            m.Chosen = o.o
            if m.Chosen == 2 {
                m.updated = false
                m.Chosen = 0
                m.updateMsg = fmt.Sprintf(updateMsgWrapping, "")
            }
        }
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    }

    var cmd tea.Cmd
    m.options, cmd = m.options.Update(msg)

    return m, cmd
}

func choicesView(m Model) string {
    options := lipgloss.NewStyle().Margin(optionMargin.height, optionMargin.width).
        Render(m.options.View())
    return fmt.Sprintf(optionWrapping, options)
}

