package nagios

import (
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

var (
	// https://docs.gitlab.com/ee/api/merge_requests.html#list-project-merge-requests
	// opened, closed, locked, or merged.
	gitlabMergeRequestsOpenedState = "opened"
)

type gitlabProjectMRChecker struct {
	client *gitlab.Client
}

func newGitlabProjectMRChecker(endpoint, apiToken string) (*gitlabProjectMRChecker, error) {
	c, err := gitlab.NewClient(apiToken, gitlab.WithBaseURL(endpoint))
	if err != nil {
		return nil, err
	}
	return &gitlabProjectMRChecker{
		client: c,
	}, nil
}

func (g gitlabProjectMRChecker) CheckMergeRequests(project string, targetBranch string) ([]MergeRequest, error) {
	var gmr []MergeRequest
	mr, _, err := g.client.MergeRequests.ListProjectMergeRequests(project, &gitlab.ListProjectMergeRequestsOptions{
		State:        &gitlabMergeRequestsOpenedState,
		TargetBranch: &targetBranch,
	})
	if err != nil {
		return gmr, errors.Wrap(err, "listing project merge-requests")
	}

	for _, cmr := range mr {
		gmr = append(gmr, MergeRequest{
			CreatedAt: *cmr.CreatedAt,
			UpdatedAt: *cmr.UpdatedAt,
			ID:        *&cmr.ID,
			Title:     cmr.Title,
		})
	}

	return gmr, nil
}
