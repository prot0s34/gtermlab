package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
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
		items = append(items, item{
			title:    project.NameWithNamespace,
			desc:     fmt.Sprintf("Last Activity At: %s", project.LastActivityAt.String()),
			itemType: "starredProject",
			handler:  handleProjectDetailItem,
		})
	}

	return items, nil
}
