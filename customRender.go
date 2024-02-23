package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type customDelegate struct{}

func (d customDelegate) Height() int {
	return 3
}

func (d customDelegate) Spacing() int {
	return 0
}

func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title := fmt.Sprintf("%s", i.Title())
	desc := i.Description()

	var titleStyled, descStyled string

	if index == m.Index() {
		titleStyled = selectedItemStyle.Copy().Inherit(baseTitleStyle).Render(title)
		descStyled = selectedItemStyle.Copy().Inherit(baseDescriptionStyle).Render(desc)
	} else {
		titleStyled = baseTitleStyle.Render(title)
		descStyled = baseDescriptionStyle.Render(desc)
	}

	fmt.Fprint(w, titleStyled+"\n"+descStyled)
}
