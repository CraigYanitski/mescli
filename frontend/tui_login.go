package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

// login struct
const (
    loginEmail = iota
    loginPassword
)

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

func loginView(m Model) string{
    // obscure password
    pw := m.inputs[loginPassword].Value()
    san := strings.Repeat("*", len(pw))
    m.inputs[loginPassword].SetValue(san)
    // set output string
    s := fmt.Sprintf(
        loginWrapping, 
        m.logo,
        m.inputs[loginEmail].View(), 
        m.inputs[loginPassword].View(),
    )
    // restore password
    m.inputs[loginPassword].SetValue(pw)
    return s
}

