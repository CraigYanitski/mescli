package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CraigYanitski/mescli/client"
	"github.com/CraigYanitski/mescli/typeset"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

func main() {
    // BubbleTea interface
    p := tea.NewProgram(initialModel())

    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("X3DH test")
    fmt.Printf("---------------\n")

    // Initialise clients in conversation
    alice := &client.Client{Name: "Alice"}
    _ = alice.Initialise()
    bob := &client.Client{Name: "Bob"}
    _ = bob.Initialise()

    // Perform extended triple Diffie-Hellman exchange
    _ = alice.EstablishX3DH(bob)
    fmt.Printf("\nX3DH initialised\n")
    _ = bob.CompleteX3DH(alice)
    fmt.Printf("\nX3DH established\n")

    // Check if exchange was successful
    // Confirm whether or not they are equal, and thus the exchange is complete
    var result string
    if alice.CheckSecretEqual(bob) {
        result, _ = typeset.FormatString("\nDiffie-Hellman secrets match - extended triple Diffie-Hellman exchange complete", 
            []string{"italics", "green"})
    } else {
        result, _ = typeset.FormatString("\nDiffie-Hellman secrets don't match - error in establishing X3DH exchange! Secrets are not equal!!", 
            []string{"italics", "red"})
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Encryption test")
    fmt.Printf("---------------\n")

    // Try to send a message from Alice to Bob
    alicePub, _ := alice.Identity()
    bobPub, _ := bob.Identity()
    message := "Hi Bob!!"
    ciphertext, err := alice.SendMessage(message, []string{"blue"}, bobPub)
    if err != nil {
        panic(err)
    }
    plaintext, err := bob.ReceiveMessage(ciphertext, alicePub)
    if err != nil {
        panic(err)
    }

    // Define progress strings
    initMessage, _ := typeset.FormatString("\ninitial message (%d): ", []string{"yellow", "bold"})
    initMessage += "%s\n"
    encrMessage, _ := typeset.FormatString("\nencrypted message (%d): ", []string{"yellow", "bold"})
    encrMessage += "0x%x\n"
    decrMessage, _ := typeset.FormatString("\ndecrypted message (%d): ", []string{"yellow", "bold"})
    decrMessage += "%s\n"

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result, _ = typeset.FormatString("\nMessage Encryption successful!!", []string{"green"})
    } else {
        result, _ = typeset.FormatString("\nError in message encryption!", []string{"red"})
    }
    fmt.Println(result)

    // Test prelude
    fmt.Printf("\n---------------\n")
    fmt.Println("Message length test")
    fmt.Printf("---------------\n")

    // Try to send a message from Alice to Bob
    message = "I am wondering about how much text I can put in a message before it encryption truncates. " +
              "There is obviously some entropy limit that cannot be surpassed given the SHA256 hashing function. " +
              "Perhaps this sentence will not make it through the transmission? " +
              "I should start splitting the message into chunks before finishing the encryption. " +
              "This message is clearly a good way to test this functionality."
    ciphertext, err = alice.SendMessage(message, []string{}, bobPub)
    if err != nil {
        panic(err)
    }
    plaintext, err = bob.ReceiveMessage(ciphertext, alicePub)
    if err != nil {
        panic(err)
    }

    // Print progress
    fmt.Printf(initMessage, len(message), message)
    fmt.Printf(encrMessage, len(ciphertext), ciphertext)
    fmt.Printf(decrMessage, len(plaintext), plaintext)

    // Compare result
    if strings.Contains(plaintext, message) {
        result, _ = typeset.FormatString("\nMessage Encryption successful!!", []string{"italics", "green"})
    } else {
        result, _ = typeset.FormatString("\nError in message encryption!", []string{"italics", "red"})
    }
    fmt.Println(result)
}

type (
    errMsg error
)

type model struct {
    viewport     viewport.Model
    messages     []string
    textarea     textarea.Model
    senderStyle  lipgloss.Style
    err          error
}

func initialModel() model {
    ta := textarea.New()
    ta.Placeholder = "Message"
    ta.Focus()

    ta.Prompt = "| "
    ta.CharLimit = 256

    ta.FocusedStyle.Prompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))

    ta.SetWidth(30)
    ta.SetHeight(3)

    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

    ta.ShowLineNumbers = false

    vp := viewport.New(30, 5)
    vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")
    vp.Style.Bold(true)
    
    ta.KeyMap.InsertNewline.SetEnabled(false)

    return model {
        textarea:     ta,
        messages:     []string{},
        viewport:     vp,
        senderStyle:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")),
        err:          nil,
    }
}

func (m model) Init() tea.Cmd {
    return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        tiCmd tea.Cmd
        vpCmd tea.Cmd
    )

    m.textarea, tiCmd = m.textarea.Update(msg)
    m.viewport, vpCmd = m.viewport.Update(msg)

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.textarea.SetWidth(msg.Width)
        m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

        if len(m.messages) > 0 {
            // Wrap content before setting it
            m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
        }
        m.viewport.GotoBottom()
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
            fmt.Println(m.textarea.Value())
            return m, tea.Quit
        case tea.KeyEnter:
            if strings.TrimSpace(m.textarea.Value()) != "" {
                m.viewport.Style.Bold(false)
                m.messages = append(m.messages, m.senderStyle.Render("You: ") + m.textarea.Value())
                m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
                    Render(strings.Join(m.messages, "\n")))
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

func (m model) View() string {
    return fmt.Sprintf(
        "%s%s%s",
        m.viewport.View(),
        gap,
        m.textarea.View(),
    )
}

