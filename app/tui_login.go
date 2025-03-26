package main

import (
	"errors"
	"fmt"
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
    if viper.GetString("access_token") != "" {
        m.loggedIn = true
        return m, nil
    }

    cmds := make([]tea.Cmd, len(m.loginInputs))
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlC, tea.KeyEscape:
            m.Quitting = true
            return m, tea.Quit 
        case tea.KeyTab, tea.KeyDown:
            m.loginFocus = (m.loginFocus + 1) % len(m.loginInputs)
        case tea.KeyShiftTab, tea.KeyUp:
            m.loginFocus = ((m.loginFocus - 1) % len(m.loginInputs) + len(m.loginInputs)) % len(m.loginInputs)
        case tea.KeyEnter:
            err := loginWithPassword(
                m.loginInputs[loginEmail].Value(), 
                m.loginInputs[loginPassword].Value(),
            )
            if err != nil {
                m.loginMsg = fmt.Sprintf(loginMsgWrapping, errorStyle.Render("Invalid login"))
                return m, nil
            }
            m.loggedIn = true
            m.loginFocus = 0
            for i, _ := range m.loginInputs {
                m.loginInputs[i].SetValue("")
            }
        case tea.KeyCtrlN:
            m.created = false
            return m, nil
        }
        for i := range m.loginInputs {
            m.loginInputs[i].Blur()
        }
        m.loginInputs[m.loginFocus].Focus()
    }
    for i := range m.loginInputs {
        m.loginInputs[i], cmds[i] = m.loginInputs[i].Update(msg)
    }
    return m, tea.Batch(cmds...)
}

func loginView(m Model) string{
    // obscure password
    pw := m.loginInputs[loginPassword].Value()
    san := strings.Repeat("*", len(pw))
    m.loginInputs[loginPassword].SetValue(san)
    // set output string
    s := fmt.Sprintf(
        loginWrapping, 
        m.logo,
        m.loginInputs[loginEmail].View(), 
        m.loginInputs[loginPassword].View(),
        m.loginMsg,
    )
    // restore password
    m.loginInputs[loginPassword].SetValue(pw)
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

