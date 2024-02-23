package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := initGitlabClient()
	if err != nil {
		fmt.Println("Error initializing GitLab client:", err)
		os.Exit(1)
	}

	rootList := configureRootList()
	rootView := ListView{list: rootList}

	m := model{views: []View{rootView}}
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
