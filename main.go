package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := initGitlabClient()
	if err != nil {
		fmt.Println("Error initializing GitLab client:", err)
		os.Exit(1)
	}

	mainList := configureRootList()
	m := model{listStack: []list.Model{mainList}}
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
