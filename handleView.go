package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

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
