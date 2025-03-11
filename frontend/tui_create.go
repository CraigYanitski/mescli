package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func updateCreate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    cmds := make([]tea.Cmd, len(m.updateInputs))
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEscape:
            m.Quitting = true
            return m, tea.Quit 
        case tea.KeyTab, tea.KeyDown:
            m.updateFocus = (m.updateFocus + 1) % len(m.updateInputs)
        case tea.KeyShiftTab, tea.KeyUp:
            m.updateFocus = ((m.updateFocus - 1) % len(m.updateInputs) + len(m.updateInputs)) % len(m.updateInputs)
        case tea.KeyEnter:
            if err := m.updateInputs[updateEmail].Err; err != nil {
                m.createMsg = fmt.Sprintf(createMsgWrapping, errorStyle.Render(err.Error()))
                return m, nil
            } else if err = m.updateInputs[updatePassword].Err; err != nil {
                m.createMsg = fmt.Sprintf(createMsgWrapping, errorStyle.Render(err.Error()))
                return m, nil
            } else if m.updateInputs[updateRetypePassword].Value() != m.updateInputs[updatePassword].Value() {
                m.createMsg = fmt.Sprintf(createMsgWrapping, errorStyle.Render("Passwords do not match"))
                return m, nil
            }
            err := createAccount(
                m.updateInputs[updateName].Value(),
                m.updateInputs[updateEmail].Value(),
                m.updateInputs[updatePassword].Value(),
            )
            if err != nil {
                m.createMsg = fmt.Sprintf(createMsgWrapping, errorStyle.Render("Update failed"))
                return m, nil
            }
            m.created = true
            m.loggedIn = true
            m.updateFocus = 0
            for i, _ := range m.updateInputs {
                m.updateInputs[i].SetValue("")
            }
        }
        for i := range m.updateInputs {
            m.updateInputs[i].Blur()
        }
        m.updateInputs[m.updateFocus].Focus()
    }
    for i := range m.updateInputs {
        m.updateInputs[i], cmds[i] = m.updateInputs[i].Update(msg)
    }
    return m, tea.Batch(cmds...)
}

func createView(m Model) string{
    // obscure password
    pw := m.updateInputs[updatePassword].Value()
    pwTemp := strings.Repeat("*", len(pw))
    m.updateInputs[updatePassword].SetValue(pwTemp)
    rpw := m.updateInputs[updatePassword].Value()
    rpwTemp := strings.Repeat("*", len(rpw))
    m.updateInputs[updateRetypePassword].SetValue(rpwTemp)
    // set output string
    s := fmt.Sprintf(
        updateWrapping, 
        m.logo,
        m.updateInputs[updateName].View(), 
        m.updateInputs[updateEmail].View(), 
        m.updateInputs[updatePassword].View(),
        m.updateInputs[updateRetypePassword].View(),
        m.createMsg,
    )
    // restore password
    m.updateInputs[updatePassword].SetValue(pw)
    m.updateInputs[updateRetypePassword].SetValue(rpw)
    return s
}

