package tui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CraigYanitski/mescli/internal/requests"
	"github.com/CraigYanitski/mescli/internal/utils"
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
                rawMsg := utils.RawMessage{
                    Sender: utils.SelfType,
                    Message: strings.TrimSpace(m.textarea.Value()),
                    Time: time.Now(),
                }
                m.cfg.Messages[m.conversation] = append(
                    m.cfg.Messages[m.conversation], 
                    rawMsg,
                )
                message := renderMessage(m, rawMsg)
                m.messages[m.conversation] = append(
                    m.messages[m.conversation], 
                    message,
                )
				err := requests.SendMessage(m.conversation, message)
				if err != nil {
					m.messages[m.conversation] = append(
						m.messages[m.conversation], 
						"Send failed...",
					)
				}
                ok := requests.WriteMessages(m.cfg.Messages)
                if !ok {
                    log.Fatal("error writing messages")
                }
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
        conversationStyle.Render(m.conversation),
        m.viewport.View(),
        m.textarea.View(),
    )
}

func renderMessage(m Model, rawMsg utils.RawMessage) string {
    var prompt string
    switch rawMsg.Sender {
    case utils.SelfType:
        prompt = m.Prompt
    case utils.ContactType:
        if m.receivePrompt != "" {
            prompt = m.receiveStyle.Render(m.receivePrompt)
        } else {
            prompt = m.receiveStyle.Render(m.conversation)
        }
    }
    renderer, err := glamour.NewTermRenderer(
        glamour.WithStylePath("tokyo-night"), 
        glamour.WithWordWrap(m.viewport.Width - len(m.senderPrompt)),
    )
    if err != nil {
        renderer, _ = glamour.NewTermRenderer()
    }
    messageMD, err := renderer.Render(rawMsg.Message)
    if err != nil {
        // fallback to unformatted text if there is an issue rendering the markdown
        messageMD = rawMsg.Message
    }
    messageMD = strings.TrimSpace(messageMD)
    message := prompt + strings.Replace(messageMD, "m  ", "m", 1)
    return message
}

