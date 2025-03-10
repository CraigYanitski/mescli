package main

import (
	"github.com/charmbracelet/bubbles/list"
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
    loginWrapping = "\n%s\n\n\n\n\n\n%s\n\n%s\n\n\n%s\n"
    loginMsgWrapping = "enter to submit credentials\nctrl+n to create a new account\n\n%s"
    updateWrapping = "\n%s\n\n\n\n\n\n%s\n\n%s\n\n%s\n\n%s\n\n\n%s\n"
    updateMsgWrapping = "enter to submit credentials\nctrl+n to update your account\n\n%s"
    conversationWrapping = "\n%s\n\n%s\n\n%s"
    optionWrapping = optionStyle.Margin(optionMargin.height, optionMargin.width).
        Render("\nPlease choose an option\n%s\n")
    contactWrapping = contactStyleName.Margin(contactMargin.height, contactMargin.width).
        Render("\nConversations\n%s\n")
    helpWrapping = helpStyle.Margin(contactMargin.height, contactMargin.width).
        Render("\nKey Bindings\n%s\n")
)
type (
    errMsg error
)

