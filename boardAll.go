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

const margin = 4

var board *Board

type Stage struct {
	Name   string
	Status int
}

type Job struct {
	Title       string
	Description string
	Stage       int
}

var stages []Stage = []Stage{
	{Name: "Stage 0", Status: 0},
	{Name: "Stage 1", Status: 1},
	{Name: "Stage 2", Status: 2},
	{Name: "Stage 3", Status: 3},
	{Name: "Stage 4", Status: 4},
	{Name: "Stage 5", Status: 5},
}

var jobs []Job = []Job{
	{Title: "Job 1", Description: "Description 1", Stage: 0},
	{Title: "Job 2", Description: "Description 2", Stage: 5},
}

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

type column struct {
	focus  bool
	stage  Stage
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

func newColumn(stage Stage) column {
	var focus bool
	if stage.Status == stages[0].Status {
		focus = true
	}

	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(false)

	defaultList.Title = stage.Name

	return column{focus: focus, stage: stage, list: defaultList}
}

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
		case key.Matches(msg, keys.Enter):
			println("enter")
		}
	}
	c.list, cmd = c.list.Update(msg)
	return c, cmd
}

func (c column) View() string {
	return c.getStyle().Render(c.list.View())
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

func (b *Board) initColumns() {
	b.cols = make([]column, len(stages))
	for i, stage := range stages {
		b.cols[i] = newColumn(stage)
	}
}

type JobListItem struct {
	Job
}

func (j JobListItem) Title() string       { return j.Job.Title }
func (j JobListItem) Description() string { return j.Job.Description }
func (j JobListItem) FilterValue() string { return j.Job.Title }

func NewJobListItem(job Job) JobListItem {
	return JobListItem{Job: job}
}

func (b *Board) initLists() {
	b.initColumns()

	for _, job := range jobs {
		for i, col := range b.cols {
			if col.stage.Status == job.Stage {
				jobItem := NewJobListItem(job)
				b.cols[i].list.InsertItem(len(b.cols[i].list.Items()), jobItem)
				break
			}
		}
	}
	b.loaded = true
}

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         column
	index       int
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

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Right key.Binding
	Left  key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
	Back  key.Binding
}

var keys = keyMap{
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

	initialFocus := status(stages[0].Status)

	return &Board{
		help:    help,
		focused: initialFocus,
	}
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

	var columnViews []string
	for _, col := range m.cols {
		columnViews = append(columnViews, col.View())
	}

	board := lipgloss.JoinHorizontal(lipgloss.Left, columnViews...)

	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}

type Task struct {
	status      status
	title       string
	description string
}

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
	for i, stage := range stages {
		if stage.Status == int(s) {
			if i == len(stages)-1 {
				return status(stages[0].Status)
			}
			return status(stages[i+1].Status)
		}
	}
	return s
}

func (s status) getPrev() status {
	for i, stage := range stages {
		if stage.Status == int(s) {
			if i == 0 {
				return status(stages[len(stages)-1].Status)
			}
			return status(stages[i-1].Status)
		}
	}
	return s
}
