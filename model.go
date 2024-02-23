package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Update(msg tea.Msg) (View, tea.Cmd)
	View() string
}

type ListView struct {
	list list.Model
}

type model struct {
	views             []View
	lastWindowSize    tea.WindowSizeMsg
	SelectedProjectID int
}

func (l ListView) Update(msg tea.Msg) (View, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l ListView) View() string {
	return docStyle.Render(l.list.View())
}
func (m *model) pushView(v View) {
	if listView, ok := v.(ListView); ok {
		horizontalPadding, verticalPadding := 4, 4
		width := m.lastWindowSize.Width - horizontalPadding*2
		height := m.lastWindowSize.Height - verticalPadding*2
		listView.list.SetSize(width, height)

		v = listView
	}

	m.views = append(m.views, v)
}

func (m *model) popView() {
	if len(m.views) > 1 {
		m.views = m.views[:len(m.views)-1]
	}
}

func (m *model) currentView() View {
	if len(m.views) > 0 {
		return m.views[len(m.views)-1]
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if currentView, ok := m.currentView().(ListView); ok {
				selectedItem, ok := currentView.list.SelectedItem().(item)
				if ok && selectedItem.handler != nil {
					cmd := selectedItem.handler(m)
					return m, cmd
				}
			}

		case "esc":
			m.popView()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.lastWindowSize = msg

		if currentView, ok := m.currentView().(ListView); ok {
			horizontalPadding, verticalPadding := 4, 4
			width := msg.Width - horizontalPadding*2
			height := msg.Height - verticalPadding*2
			currentView.list.SetSize(width, height)

			m.views[len(m.views)-1] = currentView
		}
	}

	if currentView := m.currentView(); currentView != nil {
		updatedView, cmd := currentView.Update(msg)
		m.views[len(m.views)-1] = updatedView
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	currentView := m.currentView()
	if currentView != nil {
		return currentView.View()
	}
	return ""
}
