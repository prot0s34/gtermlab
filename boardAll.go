package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (c *column) setSize(width, height int) {
	c.width = width / margin
	c.height = height

	c.list.SetSize(c.width, c.height-2)
}
func (b *Board) SetSize(width, height int) {
	b.width = width
	b.height = height
	for i := range b.cols {
		b.cols[i].setSize(width, height)
	}
}

const APPEND = -1

type column struct {
	focus  bool
	status status
	list   list.Model
	height int
	width  int
}

func (c *column) Focus() {
	c.focus = true
}

func (c *column) Blur() {
	c.focus = false
}

func (c *column) Focused() bool {
	return c.focus
}

func newColumn(status status) column {
	var focus bool
	if status == todo {
		focus = true
	}
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)
	return column{focus: focus, status: status, list: defaultList}
}

// Init does initial setup for the column.
func (c column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c column) Update(msg tea.Msg) (View, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
		c.list.SetSize(msg.Width/margin, msg.Height/2)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Edit):
			if len(c.list.VisibleItems()) != 0 {
				task := c.list.SelectedItem().(Task)
				f := NewForm(task.title, task.description)
				f.index = c.list.Index()
				f.col = c
				return f.Update(nil)
			}
		case key.Matches(msg, keys.New):
			f := newDefaultForm()
			f.index = APPEND
			f.col = c
			return f.Update(nil)
		case key.Matches(msg, keys.Delete):
			return c, c.DeleteCurrent()
		case key.Matches(msg, keys.Enter):
			return c, c.MoveToNext()
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
}

func (c *column) DeleteCurrent() tea.Cmd {
	if len(c.list.VisibleItems()) > 0 {
		c.list.RemoveItem(c.list.Index())
	}

	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)
	return cmd
}

func (c *column) Set(i int, t Task) tea.Cmd {
	if i != APPEND {
		return c.list.SetItem(i, t)
	}
	return c.list.InsertItem(APPEND, t)
}

func (c *column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

type moveMsg struct {
	Task
}

func (c *column) MoveToNext() tea.Cmd {
	var task Task
	var ok bool
	// If nothing is selected, the SelectedItem will return Nil.
	if task, ok = c.list.SelectedItem().(Task); !ok {
		return nil
	}
	// move item
	c.list.RemoveItem(c.list.Index())
	task.status = c.status.getNext()

	// refresh list
	var cmd tea.Cmd
	c.list, cmd = c.list.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return moveMsg{task} })
}

// Provides the mock data to fill the pipeline stages board

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}
	b.cols[todo].list.Title = "Lint"
	b.cols[todo].list.SetItems([]list.Item{
		Task{status: todo, title: "Go lint", description: "linting"},
		Task{status: todo, title: "Sonarqube", description: "static analysis"},
		Task{status: todo, title: "SAST", description: "security analysis"},
	})
	b.cols[inProgress].list.Title = "Build"
	b.cols[inProgress].list.SetItems([]list.Item{
		Task{status: inProgress, title: "Go build", description: "build the thing"},
	})
	b.cols[done].list.Title = "Deploy"
	b.cols[done].list.SetItems([]list.Item{
		Task{status: done, title: "Nexus", description: "deploy artifacts"},
	})
	b.loaded = true
}

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         column
	index       int
}

func newDefaultForm() *Form {
	return NewForm("task name", "")
}

func NewForm(title, description string) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
	}
	form.title.Placeholder = title
	form.description.Placeholder = description
	form.title.Focus()
	return &form
}

func (f Form) CreateTask() Task {
	return Task{f.col.status, f.title.Value(), f.description.Value()}
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f *Form) Update(msg tea.Msg) (View, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case column:
		f.col = msg
		f.col.list.Index()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return f, tea.Quit

		case key.Matches(msg, keys.Back):
			return board.Update(nil)
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			return board.Update(f)
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"Create a new task",
		f.title.View(),
		f.description.View(),
		f.help.View(keys))
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

type keyMap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Left   key.Binding
	Enter  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}

type Board struct {
	help     help.Model
	loaded   bool
	focused  status
	cols     []column
	quitting bool
	width    int
	height   int
}

func NewBoard() *Board {
	help := help.New()
	help.ShowAll = true
	return &Board{help: help, focused: todo}
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func (b *Board) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// var cmd tea.Cmd
		var cmds []tea.Cmd
		b.help.Width = msg.Width - margin
		for i := 0; i < len(b.cols); i++ {
			updatedCol, cmd := b.cols[i].Update(msg)
			b.cols[i] = updatedCol.(column)
			cmds = append(cmds, cmd)
		}
		b.loaded = true
		return b, tea.Batch(cmds...)
	case Form:
		return b, b.cols[b.focused].Set(msg.index, msg.CreateTask())
	case moveMsg:
		return b, b.cols[b.focused.getNext()].Set(APPEND, msg.Task)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			b.quitting = true
			return b, tea.Quit
		case key.Matches(msg, keys.Left):
			b.cols[b.focused].Blur()
			b.focused = b.focused.getPrev()
			b.cols[b.focused].Focus()
		case key.Matches(msg, keys.Right):
			b.cols[b.focused].Blur()
			b.focused = b.focused.getNext()
			b.cols[b.focused].Focus()
		}
	}
	res, cmd := b.cols[b.focused].Update(msg)
	if _, ok := res.(column); ok {
		b.cols[b.focused] = res.(column)
	} else {
		return res, cmd
	}
	return b, cmd
}

func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.cols[todo].View(),
		m.cols[inProgress].View(),
		m.cols[done].View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}

type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t *Task) Next() {
	if t.status == done {
		t.status = todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type status int

func (s status) getNext() status {
	if s == done {
		return todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == todo {
		return done
	}
	return s - 1
}

const margin = 4

var board *Board

const (
	todo status = iota
	inProgress
	done
)
