package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListView struct {
	list list.Model
}

func (l ListView) Update(msg tea.Msg) (View, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l ListView) View() string {
	return docStyle.Render(l.list.View())
}
