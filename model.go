package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	listStack         []list.Model
	lastWindowSize    tea.WindowSizeMsg
	SelectedProjectID int
}

func (m *model) pushList(l list.Model) {
	if m.lastWindowSize != (tea.WindowSizeMsg{}) {
		horizontalPadding := 4
		verticalPadding := 4
		l.SetSize(m.lastWindowSize.Width-horizontalPadding, m.lastWindowSize.Height-verticalPadding)
	}

	m.listStack = append(m.listStack, l)
}
func (m *model) popList() {
	if len(m.listStack) > 1 {
		m.listStack = m.listStack[:len(m.listStack)-1]
	}
}

func (m *model) currentList() *list.Model {
	return &m.listStack[len(m.listStack)-1]
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			currentList := m.currentList()
			if selectedItem, ok := currentList.SelectedItem().(item); ok && selectedItem.handler != nil {
				m.SelectedProjectID = selectedItem.projectID
				return m, selectedItem.handler(m)
			}
			return m, nil
		case "esc":
			m.popList()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.lastWindowSize = msg
		currentList := m.currentList()
		horizontalPadding, verticalPadding := 4, 4
		currentList.SetSize(msg.Width-horizontalPadding, msg.Height-verticalPadding)
	}

	currentList := m.currentList()
	var cmd tea.Cmd
	*currentList, cmd = currentList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	currentList := m.currentList()
	return docStyle.Render(currentList.View())
}
