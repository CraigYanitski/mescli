package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
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
            ok, err := loginWithPassword(
                m.inputs[loginEmail].Value(), 
                m.inputs[loginPassword].Value(),
            )
            if err != nil {
                log.Println(err)
                m.loginMsg += "\n\nInvalid login"
            }
            if ok {
                m.loggedIn = true
            }
        case tea.KeyCtrlN:
            if err := m.inputs[loginEmail].Err; err != nil {
                m.loginMsg += fmt.Sprintf("\n\n%s", err)
                return m, nil
            } else if err = m.inputs[loginPassword].Err; err != nil {
                m.loginMsg += fmt.Sprintf("\n\n%s", err)
                return m, nil
            }
            ok, err := createAccount(
                m.inputs[loginEmail].Value(),
                m.inputs[loginPassword].Value(),
            )
            if err != nil {
                log.Println(err)
                m.loginMsg += "\n\nInvalid login"
                return m, nil
            }
            if ok {
                m.loggedIn = true
            }
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
        m.loginMsg,
    )
    // restore password
    m.inputs[loginPassword].SetValue(pw)
    return s
}

func emailValidator(s string) error {
    ok, err := regexp.MatchString(`([a-zA-Z0-9\._-]+)\@`, s)
    if err == nil && !ok {
        err = errors.New("not a valid email")
    }
    if err != nil {
        return err
    }
    ok, err = regexp.MatchString(`(\@[a-zA-Z0-9\._-]+)\.`, s)
    if err == nil && !ok {
        err = errors.New("not a valid email")
    }
    if err != nil {
        return err
    }
    ok, err = regexp.MatchString(`(\.[a-zA-Z0-9]+)`, s)
    if err == nil && !ok {
        err = errors.New("not a valid email")
    }
    return err
}

func passwordValidator(s string) error {
    ok, err := regexp.MatchString(`[a-z]+`, s)
    if err == nil && !ok {
        err = errors.New("password does not contain lowercase characters")
    }
    if err != nil {
        return err
    }
    ok, err = regexp.MatchString(`[A-Z]+`, s)
    if err == nil && !ok {
        err = errors.New("password does not contain uppercase characters")
    }
    if err != nil {
        return err
    }
    ok, err = regexp.MatchString(`[^a-zA-Z0-9]+`, s)
    if err == nil && !ok {
        err = errors.New("password does not contain special characters")
    }
    if err != nil {
        return err
    }
    ok, err = regexp.MatchString(`.{8,}`, s)
    if err == nil && !ok {
        err = errors.New("password not long enough")
    }
    return err
}

