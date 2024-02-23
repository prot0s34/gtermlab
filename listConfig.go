package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	title, desc string
	itemType    string
	handler     func(*model) tea.Cmd
	projectID   int
}

func configureList(items []list.Item, title string) list.Model {
	l := list.New(items, customDelegate{}, 0, 0)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	return l
}

func configureRootList() list.Model {
	items := []list.Item{
		item{title: "Starred", desc: "List of starred projects", itemType: "root", handler: handleStarredItem},
		item{title: "Search", desc: "Filter by name", itemType: "root", handler: nil},
	}

	rootTitle := fmt.Sprintf("󰮠 %s", gitlabInstanceName)
	return configureList(items, rootTitle)
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
