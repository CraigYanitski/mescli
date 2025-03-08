package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

