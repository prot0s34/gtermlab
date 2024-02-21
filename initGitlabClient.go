package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

var gitClient *gitlab.Client
var gitlabInstanceName string
var gitlabURL string

func initGitlabClient() error {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITLAB_TOKEN environment variable is not set")
	}

	gitlabInstanceName = os.Getenv("GITLAB_URL")
	if gitlabInstanceName == "" || gitlabInstanceName == "gitlab.com" {
		gitlabURL = "https://gitlab.com"
		gitlabInstanceName = "gitlab.com"
	} else if !strings.HasPrefix(gitlabInstanceName, "http://") && !strings.HasPrefix(gitlabInstanceName, "https://") {
		gitlabURL = "https://" + gitlabInstanceName
	} else {
		gitlabURL = gitlabInstanceName
	}

	var err error
	gitClient, err = gitlab.NewClient(token, gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		return fmt.Errorf("failed to create GitLab client: %v", err)
	}

	return nil
}
