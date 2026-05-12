package utils

import "github.com/charmbracelet/lipgloss"

var (
    // stateful output
    StatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
    SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
    ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)
