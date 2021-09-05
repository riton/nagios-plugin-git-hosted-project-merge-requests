package nagios

type GitMergeRequestChecker interface {
	CheckMergeRequests(project string, targetBranch string) ([]MergeRequest, error)
}
