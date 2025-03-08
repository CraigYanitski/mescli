package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
    logo      string
    err       error
    Quitting  bool
}

// model initialiser
func InitialModel() Model {
    // open logo
    var logo string
    file, err := os.ReadFile("assets/logo.txt")
    if err != nil {
        log.Printf("error: %s", err)
        logo = "mescli"
    }
    logo = string(file)
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
        logo:          logo,
        err:           nil,
    }
}

func (m Model) Init() tea.Cmd {
    return textarea.Blink
}

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
        s = choicesView(m)
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

