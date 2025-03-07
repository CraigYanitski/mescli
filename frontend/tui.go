package main

import (
	"fmt"
	"io"

	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
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
    // login styles
    inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(164))

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
    converstionStyle  = lipgloss.NewStyle().Bold(true)
    outputStyle       = lipgloss.NewStyle()
    senderStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("164"))
    promptStyle       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))
)

// additional output
var (
    loginWrapping = "\n%s\n\n%s\n"
    conversationWrapping = "\n%s\n\n%s\n\n%s"
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
    helpWrapping = helpStyle.Margin(contactMargin.height, contactMargin.width).
        Render("\nKey Bindings\n%s\n")
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

// list key map
type listKeyMap struct {
    optionUp        key.Binding
    optionDown      key.Binding
    toggleHelpMenu  key.Binding
    Enter           key.Binding
    Back            key.Binding
    Quit            key.Binding
}
func newListKeyMap() *listKeyMap {
    return &listKeyMap{
        optionUp: key.NewBinding(
            key.WithKeys("up", "k"),
            key.WithHelp("up", "previous option"),
        ),
        optionDown: key.NewBinding(
            key.WithKeys("down", "j"),
            key.WithHelp("down", "next option"),
        ),
        toggleHelpMenu: key.NewBinding(
            key.WithKeys("ctrl+h", "h"),
            key.WithHelp("ctrl+h | h", "display help menu"),
        ),
        Enter: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp("enter", "select option"),
        ),
        Back: key.NewBinding(
            key.WithKeys("esc", "backspace"),
            key.WithHelp("esc | backspace", "previous menu"),
        ),
        Quit: key.NewBinding(
            key.WithKeys("ctrl+c", "q"),
            key.WithHelp("ctrl+c | q", "quit mescli"),
        ),
    }
}

// login struct
const (
    loginEmail = iota
    loginPassword
)


// model parameters
type Model struct {
    // geometry
    height  int
    width   int
    // list key map
    keys  *listKeyMap
    // login
    loggedIn  bool
    inputs    []textinput.Model
    focused   int
    // options
    options  list.Model
    Chosen   int
    // contacts
    contacts      list.Model
    conversation  string
    // conversation
    viewport      viewport.Model
    messages      map[string][]string
    textarea      textarea.Model
    senderPrompt  string
    senderStyle   lipgloss.Style
    Prompt        string
    help          help.Model
    // help
    viewHelp  bool
    // misc
    err         error
    Quitting    bool
}

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

/////////////////
// INITIALIZER //
/////////////////

// model initialiser
func InitialModel() Model {
    // login textinput
    inputs := make([]textinput.Model, 2)
    inputs[loginEmail] = textinput.New()
    inputs[loginEmail].Placeholder = "email"
    inputs[loginEmail].Focus()
    inputs[loginEmail].CharLimit = 256
    inputs[loginEmail].Width = 50
    inputs[loginEmail].Prompt = ""
    inputs[loginPassword] = textinput.New()
    inputs[loginPassword].Placeholder = "password"
    inputs[loginPassword].CharLimit = 256
    inputs[loginPassword].Width = 50
    inputs[loginPassword].Prompt = ""

    // option list
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
    o.SetShowHelp(false)
    
    // contact list
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
    c.SetShowHelp(false)

    // contact messages
    messages := make(map[string][]string)

    // conversation textarea
    ta := textarea.New()
    // DefaultKeyMap is the default set of key bindings for navigating and acting
    // upon the textarea.
    var textareaKeyMap = textarea.KeyMap{
    	CharacterForward: key.NewBinding(
            key.WithKeys("right", "ctrl+f"), 
            key.WithHelp("right", "character forward"),
        ),
    	CharacterBackward: key.NewBinding(
            key.WithKeys("left", "ctrl+b"), 
            key.WithHelp("left", "character backward"),
        ),
    	WordForward: key.NewBinding(
            key.WithKeys("ctrl+right", "alt+f"), 
            key.WithHelp("ctrl+right", "word forward"),
        ),
    	WordBackward: key.NewBinding(
            key.WithKeys("ctrl+left", "alt+b"), 
            key.WithHelp("ctrl+left", "word backward"),
        ),
    	LineNext: key.NewBinding(
            key.WithKeys("down"), 
            key.WithHelp("down", "next line"),
        ),
    	LinePrevious: key.NewBinding(
            key.WithKeys("up"), 
            key.WithHelp("up", "previous line"),
        ),
    	DeleteWordBackward: key.NewBinding(
            key.WithKeys("alt+backspace", "ctrl+w"), 
            key.WithHelp("alt+backspace", "delete word backward"),
        ),
    	DeleteWordForward: key.NewBinding(
            key.WithKeys("alt+delete", "alt+d"), 
            key.WithHelp("alt+delete", "delete word forward"),
        ),
    	DeleteAfterCursor: key.NewBinding(
            key.WithKeys("ctrl+k"), 
            key.WithHelp("ctrl+k", "delete after cursor"),
        ),
    	DeleteBeforeCursor: key.NewBinding(
            key.WithKeys("ctrl+u"), 
            key.WithHelp("ctrl+u", "delete before cursor"),
        ),
    	InsertNewline: key.NewBinding(
            key.WithKeys("ctrl+s"), 
            key.WithHelp("ctrl+s", "insert newline"),
        ),
    	DeleteCharacterBackward: key.NewBinding(
            key.WithKeys("backspace", "ctrl+h"), 
            key.WithHelp("backspace", "delete character backward"),
        ),
    	DeleteCharacterForward: key.NewBinding(
            key.WithKeys("delete", "ctrl+d"), 
            key.WithHelp("delete", "delete character forward"),
        ),
    	LineStart: key.NewBinding(
            key.WithKeys("home", "ctrl+a"), 
            key.WithHelp("home", "line start"),
        ),
    	LineEnd: key.NewBinding(
            key.WithKeys("end", "ctrl+e"), 
            key.WithHelp("end", "line end"),
        ),
    	Paste: key.NewBinding(
            key.WithKeys("ctrl+v"), 
            key.WithHelp("ctrl+v", "paste"),
        ),
    	InputBegin: key.NewBinding(
            key.WithKeys("alt+<", "ctrl+home"), 
            key.WithHelp("alt+<", "input begin"),
        ),
    	InputEnd: key.NewBinding(
            key.WithKeys("alt+>", "ctrl+end"), 
            key.WithHelp("alt+>", "input end"),
        ),
    
    	CapitalizeWordForward: key.NewBinding(
            key.WithKeys("alt+c"), 
            key.WithHelp("alt+c", "capitalize word forward"),
        ),
    	LowercaseWordForward: key.NewBinding(
            key.WithKeys("alt+l"), 
            key.WithHelp("alt+l", "lowercase word forward"),
        ),
    	UppercaseWordForward: key.NewBinding(
            key.WithKeys("alt+u"), 
            key.WithHelp("alt+u", "uppercase word forward"),
        ),
    
    	TransposeCharacterBackward: key.NewBinding(
            key.WithKeys("ctrl+t"), 
            key.WithHelp("ctrl+t", "transpose character backward"),
        ),
    }
    ta.KeyMap = textareaKeyMap
    ta.Placeholder = "Enter message to send"
    ta.Focus()
    ta.Prompt = "| "
    ta.CharLimit = 1024
    ta.FocusedStyle.Prompt = promptStyle
    ta.SetWidth(30)
    ta.SetHeight(3)
    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
    ta.ShowLineNumbers = false
    ta.KeyMap.InsertNewline.SetEnabled(true)
    // conversation viewport
    vp := viewport.New(30, 5)
    viewportKeyMap := viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "conversation page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "conversation page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("ctrl+u", "conversation ½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "conversation ½ page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("ctrl+up"),
			key.WithHelp("ctrl+↑", "conversation up"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+down"),
			key.WithHelp("ctrl+↓", "conversation down"),
		),
		// Left: key.NewBinding(
		// 	key.WithKeys("left", "h"),
		// 	key.WithHelp("←/h", "move left"),
		// ),
		// Right: key.NewBinding(
		// 	key.WithKeys("right", "l"),
		// 	key.WithHelp("→/l", "move right"),
		// ),
	}
    vp.KeyMap = viewportKeyMap
    welcomeMsg := lipgloss.NewStyle().Bold(true).Render(
        "Welcome to the chat room!\nType a message and press Enter to send.",
    )
    vp.SetContent(welcomeMsg)
    vp.Style.Margin(convMargin.height, convMargin.width)
    senderPrompt := "You: "

    return Model {
        loggedIn:      false,
        inputs:        inputs,
        focused:       0,
        keys:          newListKeyMap(),
        options:       o,
        contacts:      c,
        textarea:      ta,
        messages:      messages,
        viewport:      vp,
        senderStyle:   senderStyle,
        senderPrompt:  senderPrompt,
        Prompt:        senderStyle.Render(senderPrompt),
        help:          help.New(),
        err:           nil,
    }
}

func (m Model) Init() tea.Cmd {
    return textarea.Blink
}

////////////
// UPDATE //
////////////

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Check if program exits
    if m.Quitting {
        return m, tea.Quit
    }
    // Use the appropriate update function
    if !m.loggedIn {
        return updateLogin(msg, m)
    } else if m.conversation != "" {
        return updateConversation(msg, m)
    } else if m.Chosen == 0 {
        return updateChoices(msg, m)
    } else if m.Chosen == 1 {
        return updateContacts(msg, m)
    } else if m.viewHelp {
        return updateHelp(msg, m)
    } else {
        // m.View()
        return m, tea.Quit
    }
}

////////////////
// SUB-UPDATE //
////////////////

func updateLogin(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    if viper.GetString("refresh_token") != "" {
        m.loggedIn = true
        return m, nil
    }

    cmds := make([]tea.Cmd, len(m.inputs))
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEsc:
            m.Quitting = true
            return m, tea.Quit 
        case tea.KeyTab:
            m.focused = (m.focused + 1) % len(m.inputs)
        case tea.KeyShiftTab:
            m.focused = (m.focused % len(m.inputs) + len(m.inputs)) % len(m.inputs)
        case tea.KeyEnter:
            // loginUsingPassword()
            m.loggedIn = true
        }
        for i := range m.inputs {
            m.inputs[i].Blur()
        }
        m.inputs[m.focused].Focus()
    }
    for i := range m.inputs {
        m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
    }
    return m, tea.Batch(cmds...)
}

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
        switch {
        case key.Matches(msg, m.keys.Quit):
            m.Quitting = true
            return m, tea.Quit
        case key.Matches(msg, m.keys.Back):
            m.Chosen = 0
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

//////////
// VIEW //
//////////

func (m Model) View() string {
    if m.Quitting {
        return "Bye!"
    }
    var s string
    if !m.loggedIn {
        s = loginView(m)
    } else if m.conversation != "" {
        s = conversationView(m)
    } else if m.Chosen == 0 {
        s = optionsView(m)
    } else if m.Chosen == 1 {
        s = contactsView(m)
    } else if m.Chosen == 2 {
        s = ""
    } else if m.viewHelp {
        s = helpView(m)
    } else {
        s = ""
    }
    return s
}

//////////////
// SUB-VIEW //
//////////////

func loginView(m Model) string{
    pw := m.inputs[loginPassword].Value()
    san := strings.Repeat("*", len(pw))
    m.inputs[loginPassword].SetValue(san)
    s := fmt.Sprintf(
        loginWrapping, 
        m.inputs[loginEmail].View(), 
        m.inputs[loginPassword].View(),
    )
    m.inputs[loginPassword].SetValue(pw)
    return s
}

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
        conversationWrapping,
        converstionStyle.Render(m.conversation),
        m.viewport.View(),
        m.textarea.View(),
    )
}

func helpView(m Model) string {
    m.help.ShowAll = true
    help := helpStyle.Width(m.viewport.Width - 6).Margin(optionMargin.height, optionMargin.width).
        Render(m.help.View(m))
    return fmt.Sprintf(helpWrapping, help)
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
    m.viewport.Height = height - m.textarea.Height() - lipgloss.Height(conversationWrapping)
    return m
}

