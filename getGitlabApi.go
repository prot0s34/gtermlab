package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xanzy/go-gitlab"
)

func getStarredProjects() ([]list.Item, error) {
	if gitClient == nil {
		return nil, fmt.Errorf("GitLab client is not initialized")
	}

	opt := &gitlab.ListProjectsOptions{Starred: gitlab.Bool(true)}
	starredProjects, _, err := gitClient.Projects.ListProjects(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to get starred projects: %v", err)
	}

	items := make([]list.Item, 0, len(starredProjects))
	for _, project := range starredProjects {
		projectID := project.ID
		items = append(items, item{
			title:    project.NameWithNamespace,
			desc:     fmt.Sprintf("Last Activity At: %s", project.LastActivityAt.String()),
			itemType: "starredProject",
			handler: func(m *model) tea.Cmd {
				return handleProjectDetailItem(m, projectID)
			},
			projectID: projectID,
		})
	}

	return items, nil
}

func getPipelines(projectID int) ([]list.Item, error) {
	if gitClient == nil {
		return nil, fmt.Errorf("GitLab client is not initialized")
	}

	pipelines, _, err := gitClient.Pipelines.ListProjectPipelines(projectID, &gitlab.ListProjectPipelinesOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pipelines for project %d: %v", projectID, err)
	}

	items := make([]list.Item, 0, len(pipelines))
	for _, pipeline := range pipelines {
		items = append(items, item{
			title:    fmt.Sprintf("Pipeline #%d", pipeline.ID),
			desc:     fmt.Sprintf("Status: %s", pipeline.Status),
			itemType: "pipeline",
			handler:  nil,
		})
	}

	return items, nil
}

func getMRs(projectID int) ([]list.Item, error) {
	if gitClient == nil {
		return nil, fmt.Errorf("GitLab client is not initialized")
	}

	mergeRequests, _, err := gitClient.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get merge requests for project %d: %v", projectID, err)
	}

	items := make([]list.Item, 0, len(mergeRequests))
	for _, mergeRequest := range mergeRequests {
		items = append(items, item{
			title:    mergeRequest.Title,
			desc:     fmt.Sprintf("State: %s", mergeRequest.State),
			itemType: "mergeRequest",
			handler:  nil,
		})
	}

	return items, nil
}
