package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/CraigYanitski/mescli/client"
	"github.com/CraigYanitski/mescli/typeset"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// margins
type margin struct{
    width   int
    height  int
}
var (
    optionMargin = margin{2, 1}
    contactMargin = margin{2, 1}
)

// styles
var (
    // option styles
    optionStyle        = lipgloss.NewStyle()
    selectedItemStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))
    titleStyle         = lipgloss.NewStyle()
    paginationStyle    = list.DefaultStyles().PaginationStyle
    helpStyle          = list.DefaultStyles().HelpStyle
    subtleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    dotStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(" • ")

    // contact styles
    contactStyle              = lipgloss.NewStyle()
    selectedContactStyleName  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))
    selectedContactStyleDesc  = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("170"))

    // chat styles
    senderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
    promptStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
)

// additional output
var (
    conversationGap = "\n\n"
    optionWrapping = "\nPlease choose an option\n%s\n" +
        subtleStyle.Render("j/k, up/down: select") + dotStyle +
        subtleStyle.Render("enter: choose") + dotStyle +
        subtleStyle.Render("q, esc: quit")
    contactWrapping = "\nConversations\n%s\n" +
        subtleStyle.Render("j/k, up/down: select") + dotStyle +
        subtleStyle.Render("enter: choose") + dotStyle +
        subtleStyle.Render("q, esc: quit")
)

func main() {
    // bubble tea interface
    p := tea.NewProgram(initialModel())

    // run
    m, err := p.Run()
    if err != nil {
        log.Fatal(err)
    }

    // run tests if selected
    if m, ok := m.(model); ok && m.chosen == 2 {
        runTests()
    }

    // output additional padding
    fmt.Print("\n\n")
}

func runTests() {
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

// generic item list struct
type itemDelegate struct{}
func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// list information for options
type option struct {
    str  string
    o    int
}
func (o option) FilterValue() string { return "" }
type optionDelegate struct{
    itemDelegate
}
func (d optionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(option)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.str)

	fn := optionStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// list information for contacts
type contact struct {
    name, desc  string
}
func (c contact) FilterValue() string { return "" }
type contactDelegate struct{
    itemDelegate
}
func (d contactDelegate) Height() int { return 2 }
func (d contactDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(contact)
	if !ok {
		return
	}

    var str string
	// name := fmt.Sprintf("%s\n  %s", i.name, i.desc)

	// fn := contactStyle.Render
	if index == m.Index() {
        str = selectedContactStyleName.Render(i.name) + "\n" + 
            selectedContactStyleDesc.Render(i.desc)
		// fn = func(s ...string) string {
		// 	return selectedItemStyle.Render(strings.Join(s, " "))
		// }
	} else {
        str = contactStyle.Render(strings.Join([]string{i.name, i.desc}, "\n"))
    }

	fmt.Fprint(w, str)
}

// model parameters
type model struct {
    // geometry
    height  int
    width   int
    // options
    options  list.Model
    chosen   int
    // contacts
    contacts      list.Model
    conversation  string
    // conversation
    viewport     viewport.Model
    messages     map[string][]string
    textarea     textarea.Model
    senderStyle  lipgloss.Style
    // misc
    err       error
    Quitting  bool
}

// model initialiser
func initialModel() model {
    options := []list.Item{
        option{str: "View conversations", o: 1},
        option{str: "Run custom tests", o: 2},
    }
    o := list.New(options, optionDelegate{}, 20, 10)
    o.SetShowTitle(false)
    o.SetShowStatusBar(false)
    o.SetFilteringEnabled(false)
    o.Styles.PaginationStyle = paginationStyle
    o.Styles.HelpStyle = helpStyle

    contacts := []list.Item{
        contact{name: "Test contact 1", desc: "encrypted"},
        contact{name: "Test contact 2", desc: "encrypted"},
        contact{name: "Test contact 3", desc: "encrypted"},
    }
    c := list.New(contacts, contactDelegate{}, 20, 10)
    c.SetShowTitle(false)
    c.SetShowStatusBar(false)
    c.SetFilteringEnabled(false)
    c.Styles.Title = titleStyle
    c.Styles.PaginationStyle = paginationStyle
    c.Styles.HelpStyle = helpStyle

    messages := make(map[string][]string)

    ta := textarea.New()
    ta.Placeholder = "Message"
    ta.Focus()

    ta.Prompt = "| "
    ta.CharLimit = 256

    ta.FocusedStyle.Prompt = promptStyle

    ta.SetWidth(30)
    ta.SetHeight(3)

    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

    ta.ShowLineNumbers = false

    vp := viewport.New(30, 5)
    welcomeMsg := lipgloss.NewStyle().Bold(true).Render(
        "Welcome to the chat room!\nType a message and press Enter to send.",
    )
    vp.SetContent(welcomeMsg)
    
    ta.KeyMap.InsertNewline.SetEnabled(false)

    return model {
        options:      o,
        contacts:     c,
        textarea:     ta,
        messages:     messages,
        viewport:     vp,
        senderStyle:  senderStyle,
        err:          nil,
    }
}

func (m model) Init() tea.Cmd {
    return textarea.Blink
}

////////////
// UPDATE //
////////////

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Use the appropriate update function
    if m.chosen == 0 {
        return updateChoices(msg, m)
    } else if m.conversation != "" {
        return updateConversation(msg, m)
    } else {
        return updateContacts(msg, m)
    }
}

////////////////
// SUB-UPDATE //
////////////////

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEsc, tea.KeyCtrlC:
            m.Quitting = true
            return m, tea.Quit
        case tea.KeyEnter:
            o, _ := m.options.SelectedItem().(option)
            m.chosen = o.o
            if m.chosen == 2 {
                return m, tea.Quit
            }
        }
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    }

    var cmd tea.Cmd
    m.options, cmd = m.options.Update(msg)

    return m, cmd
}

func updateContacts(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            m.Quitting = true
            return m, tea.Quit
        case tea.KeyEsc, tea.KeyBackspace:
            m.chosen = 0
            return m, nil
        case tea.KeyEnter:
            c, _ := m.contacts.SelectedItem().(contact)
            m.conversation = c.name
        }
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    }

    var cmd tea.Cmd
    m.contacts, cmd = m.contacts.Update(msg)

    return m, cmd
}

func updateConversation(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
    var (
        tiCmd tea.Cmd
        vpCmd tea.Cmd
    )

    m.textarea, tiCmd = m.textarea.Update(msg)
    m.viewport, vpCmd = m.viewport.Update(msg)

    if len(m.messages[m.conversation]) > -1 {
        // Wrap content before setting it
        m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
            Render(strings.Join(m.messages[m.conversation], "\n")))
    }
    m.viewport.GotoBottom()

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m = m.resize(msg.Width, msg.Height)
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            fmt.Println(m.textarea.Value())
            return m, tea.Quit
        case tea.KeyEsc:
            m.conversation = ""
            return m, nil
        case tea.KeyEnter:
            if strings.TrimSpace(m.textarea.Value()) != "" {
                m.viewport.Style.Bold(false)
                m.messages[m.conversation] = append(m.messages[m.conversation], 
                    m.senderStyle.Render("You: ") + m.textarea.Value())
                m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
                    Render(strings.Join(m.messages[m.conversation], "\n")))
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

//////////
// VIEW //
//////////

func (m model) View() string {
    var s string
    if m.chosen == 0 {
        s = optionsView(m)
    } else if m.chosen == 2 {
        s = ""
    } else if m.conversation != "" {
        s = conversationView(m)
    } else {
        s = contactsView(m)
    }
    return s
}

//////////////
// SUB-VIEW //
//////////////

func optionsView(m model) string {
    options := lipgloss.NewStyle().Margin(optionMargin.height, optionMargin.width).
        Render(m.options.View())
    return fmt.Sprintf(optionWrapping, options)
}

func contactsView(m model) string {
    var conversations string
    if len(m.contacts.Items()) == 0 {
        conversations = "No conversations (yet)"
    } else {
        conversations = lipgloss.NewStyle().Margin(contactMargin.height, contactMargin.width).
            Render(m.contacts.View())
    }
    return fmt.Sprintf(contactWrapping, conversations)
}

func conversationView(m model) string {
    return fmt.Sprintf(
        "%s%s%s%s",
        m.conversation,
        m.viewport.View(),
        conversationGap,
        m.textarea.View(),
    )
}

///////////////////
// Miscellaneous //
///////////////////

func (m model) resize(width, height int) (model) {
    m.height = height
    m.width = width
    m.options.SetSize(width - 2*optionMargin.width, 
        height - lipgloss.Height(optionWrapping) - 2*optionMargin.height)
    m.contacts.SetSize(width -2*contactMargin.width, 
        height - lipgloss.Height(contactWrapping) - 2*contactMargin.height)
    m.viewport.Width = width
    m.textarea.SetWidth(width)
    m.viewport.Height = height - m.textarea.Height() - lipgloss.Height(conversationGap)
    return m
}

