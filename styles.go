package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	baseTitleStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingLeft(2)

	baseDescriptionStyle = lipgloss.NewStyle().
				Italic(true).
				PaddingLeft(2)
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#626262"))

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#FFA500")).
				Bold(true).Border(lipgloss.NormalBorder(), false, false, false, true)

	descriptionStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#949494")).
				Render
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)
