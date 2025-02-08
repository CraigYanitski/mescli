package main

import (
	"fmt"
	"io"
	// "log"
	"strings"

	// "github.com/CraigYanitski/mescli/client"
	// "github.com/CraigYanitski/mescli/typeset"
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
    convMargin = margin{2, 1}
)

// styles
var (
    // option styles
    optionStyle          = lipgloss.NewStyle()
    selectedOptionStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("164"))
    titleStyle           = lipgloss.NewStyle()
    paginationStyle      = list.DefaultStyles().PaginationStyle
    helpStyle            = list.DefaultStyles().HelpStyle
    subtleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    dotStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(" • ")

    // contact styles
    contactStyleName          = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
    contactStyleDesc          = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("245"))
    selectedContactStyleName  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("164"))
    selectedContactStyleDesc  = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("164"))

    // chat styles
    senderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("164"))
    promptStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
)

// additional output
var (
    conversationGap = "\n\n"
    optionWrapping = optionStyle.Margin(optionMargin.height, optionMargin.width).
        Render("\nPlease choose an option\n%s\n")
        //+
        // subtleStyle.Render("j/k, up/down: select") + dotStyle +
        // subtleStyle.Render("enter: choose") + dotStyle +
        // subtleStyle.Render("q, esc: quit")
    contactWrapping = contactStyleName.Margin(contactMargin.height, contactMargin.width).
        Render("\nConversations\n%s\n")
        //+
        // subtleStyle.Render("j/k, up/down: select") + dotStyle +
        // subtleStyle.Render("enter: choose") + dotStyle +
        // subtleStyle.Render("q, esc: quit")
)
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

	option := fmt.Sprintf("%s", i.str)

    var str string
	if index == m.Index() {
		str = selectedOptionStyle.Render("> " + option)
	} else {
	    str = optionStyle.Render("  " + option)
    }

	fmt.Fprint(w, str)
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
	if index == m.Index() {
        str = selectedContactStyleName.Render(i.name) + "\n" + 
            selectedContactStyleDesc.Render(i.desc)
	} else {
        str = contactStyleName.Render(i.name) + "\n" +
            contactStyleDesc.Render(i.desc)
    }

	fmt.Fprint(w, str + "\n")
}

// model parameters
type Model struct {
    // geometry
    height  int
    width   int
    // options
    options  list.Model
    Chosen   int
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
func InitialModel() Model {
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
    ta.KeyMap.InsertNewline.SetEnabled(false)

    vp := viewport.New(30, 5)
    welcomeMsg := lipgloss.NewStyle().Bold(true).Render(
        "Welcome to the chat room!\nType a message and press Enter to send.",
    )
    vp.SetContent(welcomeMsg)
    vp.Style.Margin(convMargin.height, convMargin.width)

    return Model {
        options:      o,
        contacts:     c,
        textarea:     ta,
        messages:     messages,
        viewport:     vp,
        senderStyle:  senderStyle,
        err:          nil,
    }
}

func (m Model) Init() tea.Cmd {
    return textarea.Blink
}

////////////
// UPDATE //
////////////

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Use the appropriate update function
    if m.Chosen == 0 {
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

func updateChoices(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEsc, tea.KeyCtrlC:
            m.Quitting = true
            return m, tea.Quit
        case tea.KeyEnter:
            o, _ := m.options.SelectedItem().(option)
            m.Chosen = o.o
            if m.Chosen == 2 {
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

func updateContacts(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            m.Quitting = true
            return m, tea.Quit
        case tea.KeyEsc, tea.KeyBackspace:
            m.Chosen = 0
            return m, nil
        case tea.KeyEnter:
            c, _ := m.contacts.SelectedItem().(contact)
            m.conversation = c.name

            if len(m.messages[m.conversation]) > -1 {
                // Wrap content before setting it
                m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).
                    Render(strings.Join(m.messages[m.conversation], "\n")))
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

func updateConversation(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
    var s string
    if m.Chosen == 0 {
        s = optionsView(m)
    } else if m.Chosen == 2 {
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

func optionsView(m Model) string {
    options := lipgloss.NewStyle().Margin(optionMargin.height, optionMargin.width).
        Render(m.options.View())
    return fmt.Sprintf(optionWrapping, options)
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

func conversationView(m Model) string {
    return fmt.Sprintf(
        "\n%s\n\n%s%s%s",
        m.conversation,
        m.viewport.View(),
        conversationGap,
        m.textarea.View(),
    )
}

///////////////////
// Miscellaneous //
///////////////////

func (m Model) resize(width, height int) (Model) {
    m.height = height
    m.width = width
    m.options.SetSize(width - 2*optionMargin.width, 
        height - lipgloss.Height(optionWrapping) - 2*optionMargin.height)
    m.contacts.SetSize(width -2*contactMargin.width, 
        height - lipgloss.Height(contactWrapping) - 2*contactMargin.height)
    m.viewport.Width = width
    m.textarea.SetWidth(width)
    m.viewport.Height = height - m.textarea.Height() - 2*lipgloss.Height(conversationGap)
    return m
}

