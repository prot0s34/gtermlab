package main

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
)

var gitClient *gitlab.Client

func initGitlabClient() error {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITLAB_TOKEN environment variable is not set")
	}

	var err error
	gitClient, err = gitlab.NewClient(token)
	if err != nil {
		return fmt.Errorf("failed to create GitLab client: %v", err)
	}

	return nil
}
