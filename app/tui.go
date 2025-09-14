package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// model parameters
type Model struct {
    // config
    cfg  *apiConfig
    // geometry
    height  int
    width   int
    // list key map
    keys  *listKeyMap
    // login
    loggedIn     bool
    loginInputs  []textinput.Model
    loginFocus   int
    loginMsg     string
    // create
    created      bool
    createMsg    string
    // update
    updated       bool
    updateInputs  []textinput.Model
    updateFocus   int
    updateMsg     string
    // options
    options  list.Model
    Chosen   int
    // contacts
    contacts      list.Model
    conversation  string
    // conversation
    viewport       viewport.Model
    messages       map[string][]string
    textarea       textarea.Model
    senderPrompt   string
    senderStyle    lipgloss.Style
    Prompt         string
    receivePrompt  string
    receiveStyle   lipgloss.Style
    help           help.Model
    // help
    viewHelp  bool
    // misc
    logo      string
    err       error
    Quitting  bool
}

// model initialiser
func InitialModel(cfg *apiConfig) Model {
    // load environment
    // TODO: simplify during installation
    godotenv.Load(".env")
    godotenv.Load("../.env")
    // open logo
    var logo string
    logoDir, ok := os.LookupEnv("MESCLI_DIR")
    if !ok {
        log.Println("there is no MESCLI_DIR env variable")
    }
    file, err := os.ReadFile(path.Join(logoDir, "assets/logo.txt"))
    if err != nil {
        log.Printf("error: %s", err)
        logo = "mescli"
    }
    logo = string(file)
    // login textinput
    loginInputs := make([]textinput.Model, 2)
    loginInputs[loginEmail] = textinput.New()
    loginInputs[loginEmail].Placeholder = "email"
    loginInputs[loginEmail].Focus()
    loginInputs[loginEmail].CharLimit = 256
    loginInputs[loginEmail].Width = 50
    loginInputs[loginEmail].Prompt = ""
    loginInputs[loginEmail].Validate = emailValidator
    loginInputs[loginPassword] = textinput.New()
    loginInputs[loginPassword].Placeholder = "password"
    loginInputs[loginPassword].CharLimit = 256
    loginInputs[loginPassword].Width = 50
    loginInputs[loginPassword].Prompt = ""
    loginInputs[loginPassword].Validate = passwordValidator

    //update textinput
    updateInputs := make([]textinput.Model, 4)
    updateInputs[updateName] = textinput.New()
    updateInputs[updateName].Placeholder = "name"
    updateInputs[updateName].Focus()
    updateInputs[updateName].CharLimit = 256
    updateInputs[updateName].Width = 50
    updateInputs[updateName].Prompt = ""
    updateInputs[updateEmail] = textinput.New()
    updateInputs[updateEmail].Placeholder = "email"
    updateInputs[updateEmail].CharLimit = 256
    updateInputs[updateEmail].Width = 50
    updateInputs[updateEmail].Prompt = ""
    updateInputs[updateEmail].Validate = emailValidator
    updateInputs[updatePassword] = textinput.New()
    updateInputs[updatePassword].Placeholder = "password"
    updateInputs[updatePassword].CharLimit = 256
    updateInputs[updatePassword].Width = 50
    updateInputs[updatePassword].Prompt = ""
    updateInputs[updatePassword].Validate = passwordValidator
    updateInputs[updateRetypePassword] = textinput.New()
    updateInputs[updateRetypePassword].Placeholder = "retype password"
    updateInputs[updateRetypePassword].CharLimit = 256
    updateInputs[updateRetypePassword].Width = 50
    updateInputs[updateRetypePassword].Prompt = ""

    // option list
    options := []list.Item{
        option{str: "View conversations", o: 1},
        option{str: "Update account", o: 2},
        option{str: "Run custom tests", o: 3},
    }
    o := list.New(options, optionDelegate{}, 20, 10)
    o.SetShowTitle(false)
    o.SetShowStatusBar(false)
    o.SetFilteringEnabled(false)
    o.Styles.PaginationStyle = paginationStyle
    o.Styles.HelpStyle = helpStyle
    o.SetShowHelp(true)

    // contact messages
    fmt.Println("TUI:", cfg.messages)
    messages := make(map[string][]string)
    messages["Test contact 1"] = []string{}
    messages["Test contact 2"] = []string{}
    messages["Test contact 3"] = []string{}

    // contact list
    contacts := []list.Item{}
    for key, _ := range messages {
        contacts = append(contacts, contact{name: key, desc: "encrypted"})
    }
    c := list.New(contacts, contactDelegate{}, 20, 10)
    c.SetShowTitle(false)
    c.SetShowStatusBar(false)
    c.SetFilteringEnabled(false)
    c.Styles.Title = titleStyle
    c.Styles.PaginationStyle = paginationStyle
    c.Styles.HelpStyle = helpStyle
    c.SetShowHelp(true)

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
    ta.SetWidth(100)
    ta.SetHeight(3)
    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
    ta.ShowLineNumbers = false
    ta.KeyMap.InsertNewline.SetEnabled(true)
    // conversation viewport
    vp := viewport.New(100, 5)
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
    receivePrompt := "> "

    return Model {
        // config
        cfg: cfg,
        // Model
        loggedIn:       false,
        loginInputs:    loginInputs,
        loginFocus:     0,
        loginMsg:       fmt.Sprintf(loginMsgWrapping, ""),
        created:        true,
        createMsg:      fmt.Sprintf(createMsgWrapping, ""),
        updated:        true,
        updateInputs:   updateInputs,
        updateFocus:    0,
        updateMsg:      updateMsgWrapping,
        keys:           newListKeyMap(),
        options:        o,
        contacts:       c,
        textarea:       ta,
        messages:       messages,
        viewport:       vp,
        senderStyle:    senderStyle,
        senderPrompt:   senderPrompt,
        Prompt:         senderStyle.Render(senderPrompt),
        receiveStyle:   receiveStyle,
        receivePrompt:  receivePrompt,
        help:           help.New(),
        logo:           logo,
        err:            nil,
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
    if !m.created {
        return updateCreate(msg, m)
    } else if !m.loggedIn {
        return updateLogin(msg, m)
    } else if !m.updated {
        return updateUpdate(msg, m)
    } else if m.viewHelp {
        return updateHelp(msg, m)
    } else if m.conversation != "" {
        return updateConversation(msg, m)
    } else if m.Chosen == 0 {
        return updateChoices(msg, m)
    } else if m.Chosen == 1 {
        return updateContacts(msg, m)
    } else if m.Chosen == 2 {
        return updateUpdate(msg, m)
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
    if !m.created {
        s = createView(m)
    } else if !m.loggedIn {
        s = loginView(m)
    } else if !m.updated {
        s = updateView(m)
    } else if m.viewHelp {
        s = helpView(m)
    } else if m.conversation != "" {
        s = conversationView(m)
    } else if m.Chosen == 0 {
        s = choicesView(m)
    } else if m.Chosen == 1 {
        s = contactsView(m)
    } else if m.Chosen == 2 {
        s = updateView(m)
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

