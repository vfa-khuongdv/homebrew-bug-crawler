package platform

import (
	"context"
	"time"
)

// Platform defines the interface that all Git platforms must implement
type Platform interface {
	// VerifyToken verifies token validity and displays scope information
	VerifyToken(ctx context.Context) error

	// GetCurrentUserRepositories retrieves current user repositories
	GetCurrentUserRepositories(ctx context.Context) ([]*RepositoryInfo, error)

	// GetOrganizationRepositories retrieves organization repositories
	GetOrganizationRepositories(ctx context.Context, orgName string) ([]*RepositoryInfo, error)

	// GetCurrentUserOrganizations retrieves current user organizations
	GetCurrentUserOrganizations(ctx context.Context) ([]string, error)

	// GetPullRequestsFromRepositoriesConcurrent fetches PRs from multiple repositories concurrently
	GetPullRequestsFromRepositoriesConcurrent(ctx context.Context, repos []string, startDate, endDate time.Time, maxWorkers int) ([]RepositoryScanJob, error)

	// GetPullRequestReviewsConcurrent retrieves reviews for multiple PRs concurrently
	GetPullRequestReviewsConcurrent(ctx context.Context, owner, repo string, prNumbers []int, maxWorkers int) (map[int][]*ReviewData, error)
}

// RepositoryInfo contains repository information
type RepositoryInfo struct {
	FullName string
	Owner    string
	Name     string
	URL      string
}

// ReviewData contains review information
type ReviewData struct {
	ReviewerLogin string
	State         string
	SubmittedAt   *time.Time
	CommentBody   string
}

// PullRequestData contains pull request information
type PullRequestData struct {
	Number      int
	Title       string
	Description string
	Author      string
	CreatedAt   time.Time
	MergedAt    *time.Time
	Labels      []string
	HTMLURL     string
	Status      string
	Reviews     []*ReviewData
}

// RepositoryScanJob represents a job to scan a single repository
type RepositoryScanJob struct {
	Owner    string
	RepoName string
	PRData   []*PullRequestData
	Error    error
}

// PlatformType represents the type of Git platform
type PlatformType string

const (
	PlatformGitHub    PlatformType = "github"
	PlatformBitbucket PlatformType = "bitbucket"
	PlatformBacklog   PlatformType = "backlog"
)

// String returns the string representation of the platform type
func (p PlatformType) String() string {
	return string(p)
}

// DisplayName returns the user-friendly display name
func (p PlatformType) DisplayName() string {
	switch p {
	case PlatformGitHub:
		return "GitHub"
	case PlatformBitbucket:
		return "Bitbucket"
	case PlatformBacklog:
		return "Backlog"
	default:
		return string(p)
	}
}
