package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateContacts(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keys.Quit):
            m.Quitting = true
            return m, tea.Quit
        case key.Matches(msg, m.keys.Back):
            m.Chosen = 0
            return m, nil
        case key.Matches(msg, m.keys.findOption):
            return m, nil
        case key.Matches(msg, m.keys.Enter):
            c, _ := m.contacts.SelectedItem().(contact)
            m.conversation = c.name

            // Wrap content before setting it
            if len(m.messages[m.conversation]) > 0 {
                m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
                    Render(strings.Join(m.messages[m.conversation], "\n")))
            } else {
                m.viewport.SetContent(lipgloss.NewStyle().Bold(true).Render(
                    "Welcome to the chat room!\nType a message and press Enter to send."))
            }
            m.viewport.GotoBottom()
        }
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    }

    var cmd tea.Cmd
    m.contacts, cmd = m.contacts.Update(msg)

    return m, cmd
}

func contactsView(m Model) string {
    var conversations string
    if len(m.contacts.Items()) == 0 {
        conversations = "No conversations (yet)"
    } else {
        conversations = lipgloss.NewStyle().Margin(contactMargin.height, contactMargin.width).
            Render(m.contacts.View())
    }
    return fmt.Sprintf(contactWrapping, conversations)
}

