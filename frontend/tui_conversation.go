package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func updateConversation(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    var (
        tiCmd tea.Cmd
        vpCmd tea.Cmd
    )

    m.textarea, tiCmd = m.textarea.Update(msg)
    m.viewport, vpCmd = m.viewport.Update(msg)

    // if len(m.messages[m.conversation]) > -1 {
    //     // Wrap content before setting it
    //     m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
    //         Render(strings.Join(m.messages[m.conversation], "\n")))
    // }
    // m.viewport.GotoBottom()

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            //fmt.Println(m.textarea.Value())
            m.Quitting = true
            return m, tea.Quit
        case tea.KeyEsc:
            m.conversation = ""
            return m, nil
        case tea.KeyCtrlH:
            m.viewHelp = true
        case tea.KeyEnter:
            if strings.TrimSpace(m.textarea.Value()) != "" {
                renderer, err := glamour.NewTermRenderer(
                    glamour.WithStylePath("tokyo-night"), 
                    glamour.WithWordWrap(m.viewport.Width - len(m.senderPrompt)),
                )
                if err != nil {
                    renderer, _ = glamour.NewTermRenderer()
                }
                messageMD, err := renderer.Render(m.textarea.Value())
                if err != nil {
                    messageMD = m.textarea.Value()
                }
                messageMD = strings.TrimSpace(messageMD)
                message := m.Prompt + strings.Replace(messageMD, "m  ", "m", 1)
                m.messages[m.conversation] = append(m.messages[m.conversation], 
                    message)
                m.viewport.SetContent(strings.Join(m.messages[m.conversation], "\n"))
                m.textarea.Reset()
                m.viewport.GotoBottom()
            } else {
                m.textarea.Reset()
            }
        }
    case errMsg:
        m.err = msg
        return m, nil
    }

    return m, tea.Batch(tiCmd, vpCmd)
}

func conversationView(m Model) string {
    return fmt.Sprintf(
        conversationWrapping,
        converstionStyle.Render(m.conversation),
        m.viewport.View(),
        m.textarea.View(),
    )
}

