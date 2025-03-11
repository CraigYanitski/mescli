package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// login struct
const (
    updateName = iota
    updateEmail
    updatePassword
    updateRetypePassword
)

func updateUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
    cmds := make([]tea.Cmd, len(m.updateInputs))
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC:
            m.Quitting = true
            return m, tea.Quit 
        case tea.KeyEscape:
            m.updated = true
            return m, nil
        case tea.KeyTab, tea.KeyDown:
            m.updateFocus = (m.updateFocus + 1) % len(m.updateInputs)
        case tea.KeyShiftTab, tea.KeyUp:
            m.updateFocus = ((m.updateFocus - 1) % len(m.updateInputs) + len(m.updateInputs)) % len(m.updateInputs)
        case tea.KeyEnter:
            if err := m.updateInputs[updateEmail].Err; err != nil {
                m.updateMsg = fmt.Sprintf(updateMsgWrapping, errorStyle.Render(err.Error()))
                return m, nil
            } else if err = m.updateInputs[updatePassword].Err; err != nil {
                m.updateMsg = fmt.Sprintf(updateMsgWrapping, errorStyle.Render(err.Error()))
                return m, nil
            } else if m.updateInputs[updateRetypePassword].Value() != m.updateInputs[updatePassword].Value() {
                m.updateMsg = fmt.Sprintf(updateMsgWrapping, errorStyle.Render("Passwords do not match"))
                return m, nil
            }
            err := updateAccount(
                m.updateInputs[updateName].Value(),
                m.updateInputs[updateEmail].Value(),
                m.updateInputs[updatePassword].Value(),
            )
            if err != nil {
                m.updateMsg = fmt.Sprintf(updateMsgWrapping, errorStyle.Render("Update failed"))
                return m, nil
            }
            m.updated = true
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

func updateView(m Model) string{
    // obscure password
    pw := m.updateInputs[updatePassword].Value()
    pwTemp := strings.Repeat("*", len(pw))
    m.updateInputs[updatePassword].SetValue(pwTemp)
    rpw := m.updateInputs[updateRetypePassword].Value()
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
        m.updateMsg,
    )
    // restore password
    m.updateInputs[updatePassword].SetValue(pw)
    m.updateInputs[updateRetypePassword].SetValue(rpw)
    return s
}

