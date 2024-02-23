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

func handleStarredItem(m *model) tea.Cmd {
	projects, err := getStarredProjects()
	if err != nil {
		fmt.Println("Error fetching starred projects:", err)
		return nil
	}

	detailList := configureList(projects, "Starred Projects")
	detailView := ListView{list: detailList}
	m.pushView(detailView)

	return nil
}

func handleProjectDetailItem(m *model, projectID int) tea.Cmd {
	detailItems := []list.Item{
		item{
			title:    "Pipelines",
			desc:     "List of pipelines",
			itemType: "pipeline",
			handler: func(m *model) tea.Cmd {
				return handlePipelines(m, projectID)
			},
		},
		item{
			title:    "MRs",
			desc:     "List of merge requests",
			itemType: "mr",
			handler: func(m *model) tea.Cmd {
				return handleMRs(m, projectID)
			},
		},
		item{
			title:    "Branches",
			desc:     "List of branches",
			itemType: "branch",
			handler:  handleBranches,
		},
	}

	detailsList := configureList(detailItems, "Project Details")
	detailsView := ListView{list: detailsList}
	m.pushView(detailsView)

	return nil
}

func handlePipelines(m *model, projectID int) tea.Cmd {
	pipelines, err := getPipelines(projectID)
	if err != nil {
		fmt.Println("Error fetching pipelines:", err)
		return nil
	}

	pipelinesList := configureList(pipelines, "Pipelines")
	pipelinesView := ListView{list: pipelinesList}
	m.pushView(pipelinesView)

	return nil
}

func handleMRs(m *model, projectID int) tea.Cmd {
	mrs, err := getMRs(projectID)
	if err != nil {
		fmt.Println("Error fetching MRs:", err)
		return nil
	}

	mrsList := configureList(mrs, "Merge Requests")
	mrsView := ListView{list: mrsList}
	m.pushView(mrsView)

	return nil
}

func handleBranches(m *model) tea.Cmd {
	// placeholder
	fmt.Println("Branches selected")
	return nil
}

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
