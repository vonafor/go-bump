package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitLabHosting struct {
	client    *gitlab.Client
	projectID string
}

func NewGitLabHosting(library string) *GitLabHosting {
	parts := strings.Split(library, "/")
	if len(parts) != 3 {
		panic(fmt.Sprintf("invalid library name: %s", library))
	}
	projectID := parts[len(parts)-2] + "/" + parts[len(parts)-1]
	client := gitlab.NewClient(nil, os.Getenv("GOBUMP_API_KEY"))
	err := client.SetBaseURL("https://" + parts[0])
	if err != nil {
		panic(err)
	}
	return &GitLabHosting{client: client, projectID: projectID}
}

func (h *GitLabHosting) CreateMR(title, sourceBranch, targetBranch string) (string, error) {
	opt := &gitlab.CreateMergeRequestOptions{
		Title:        gitlab.String(title),
		SourceBranch: gitlab.String(sourceBranch),
		TargetBranch: gitlab.String(targetBranch),
	}

	mergeRequest, _, err := h.client.MergeRequests.CreateMergeRequest(h.projectID, opt)
	if err != nil {
		return "", err
	}
	return mergeRequest.WebURL, err
}
