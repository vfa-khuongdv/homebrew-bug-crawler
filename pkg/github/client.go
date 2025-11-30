package github

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bug-crawler/pkg/platform"
	"github.com/google/go-github/v56/github"
)

// Client wraps GitHub API client
type Client struct {
	client *github.Client
}

// Type aliases for backward compatibility
type ReviewData = platform.ReviewData
type PullRequestData = platform.PullRequestData
type RepositoryInfo = platform.RepositoryInfo
type RepositoryScanJob = platform.RepositoryScanJob

// NewClient initializes GitHub client
func NewClient(token string) (*Client, error) {
	client := github.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}

	return &Client{
		client: client,
	}, nil
}

// GetPullRequests retrieves pull requests within a time range
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*PullRequestData, error) {
	var prs []*PullRequestData
	opts := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		githubPRs, resp, err := c.client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy PR từ %s/%s: %w", owner, repo, err)
		}

		for _, pr := range githubPRs {
			// Bỏ qua PR ngoài khoảng thời gian
			if pr.CreatedAt.Before(startDate) || pr.CreatedAt.After(endDate) {
				continue
			}

			labels := make([]string, 0)
			for _, label := range pr.Labels {
				labels = append(labels, label.GetName())
			}

			var mergedAt *time.Time
			if pr.MergedAt != nil {
				mergedAt = &pr.MergedAt.Time
			}

			status := "open"
			if !pr.GetMergedAt().IsZero() {
				status = "merged"
			}

			prData := &platform.PullRequestData{
				Number:      pr.GetNumber(),
				Title:       pr.GetTitle(),
				Description: pr.GetBody(),
				Author:      pr.GetUser().GetLogin(),
				CreatedAt:   pr.GetCreatedAt().Time,
				MergedAt:    mergedAt,
				Labels:      labels,
				HTMLURL:     pr.GetHTMLURL(),
				Status:      status,
			}

			prs = append(prs, prData)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return prs, nil
}

// GetPullRequestReviews retrieves reviews for a pull request (including issue comments)
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*ReviewData, error) {
	var reviews []*ReviewData
	opts := &github.ListOptions{PerPage: 100}

	// Get reviews from PR review API
	for {
		githubReviews, resp, err := c.client.PullRequests.ListReviews(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy reviews từ PR %d: %w", prNumber, err)
		}

		for _, review := range githubReviews {
			reviewData := &platform.ReviewData{
				ReviewerLogin: review.GetUser().GetLogin(),
				State:         review.GetState(),
				SubmittedAt:   &review.SubmittedAt.Time,
				CommentBody:   review.GetBody(),
			}
			reviews = append(reviews, reviewData)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Get additional comments from issue comments API (including review comments)
	issueOpts := &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		comments, resp, err := c.client.Issues.ListComments(ctx, owner, repo, prNumber, issueOpts)
		if err != nil {
			fmt.Printf("⚠️  Lỗi khi lấy issue comments từ PR %d: %v\n", prNumber, err)
			break
		}

		for _, comment := range comments {
			// If comment has content, add to reviews
			if comment.GetBody() != "" {
				reviewData := &platform.ReviewData{
					ReviewerLogin: comment.GetUser().GetLogin(),
					State:         "COMMENTED",
					SubmittedAt:   &comment.CreatedAt.Time,
					CommentBody:   comment.GetBody(),
				}
				reviews = append(reviews, reviewData)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		issueOpts.Page = resp.NextPage
	}

	return reviews, nil
}

// GetPullRequestsWithReviews retrieves pull requests with review data
func (c *Client) GetPullRequestsWithReviews(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*PullRequestData, error) {
	prs, err := c.GetPullRequests(ctx, owner, repo, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get reviews for each PR
	for _, pr := range prs {
		reviews, err := c.GetPullRequestReviews(ctx, owner, repo, pr.Number)
		if err != nil {
			fmt.Printf("⚠️  Lỗi khi lấy reviews cho PR #%d: %v\n", pr.Number, err)
			continue
		}
		pr.Reviews = reviews
	}

	return prs, nil
}

// VerifyToken verifies token validity and displays scope information
func (c *Client) VerifyToken(ctx context.Context) error {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return err
	}

	fmt.Printf("👤 Đăng nhập thành công với: %s\n", user.GetLogin())

	// Get rate limit information
	rateLimits, _, err := c.client.RateLimits(ctx)
	if err == nil {
		fmt.Printf("📊 Rate limit: %d/%d requests\n",
			rateLimits.Core.Remaining,
			rateLimits.Core.Limit)
	}

	return nil
}

// GetUserRepositories retrieves user repositories
func (c *Client) GetUserRepositories(ctx context.Context, username string) ([]*RepositoryInfo, error) {
	var repos []*RepositoryInfo
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		githubRepos, resp, err := c.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy repositories của %s: %w", username, err)
		}

		for _, repo := range githubRepos {
			repoInfo := &platform.RepositoryInfo{
				FullName: repo.GetFullName(),
				Owner:    repo.GetOwner().GetLogin(),
				Name:     repo.GetName(),
				URL:      repo.GetHTMLURL(),
			}
			repos = append(repos, repoInfo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return repos, nil
}

// GetOrganizationRepositories retrieves organization repositories
func (c *Client) GetOrganizationRepositories(ctx context.Context, orgName string) ([]*RepositoryInfo, error) {
	var repos []*RepositoryInfo
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        "all", // Get all public, private, internal
	}

	for {
		githubRepos, resp, err := c.client.Repositories.ListByOrg(ctx, orgName, opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy repositories của organization %s: %w", orgName, err)
		}

		for _, repo := range githubRepos {
			repoInfo := &RepositoryInfo{
				FullName: repo.GetFullName(),
				Owner:    repo.GetOwner().GetLogin(),
				Name:     repo.GetName(),
				URL:      repo.GetHTMLURL(),
			}
			repos = append(repos, repoInfo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return repos, nil
}

// GetCurrentUserRepositories retrieves current user repositories
func (c *Client) GetCurrentUserRepositories(ctx context.Context) ([]*RepositoryInfo, error) {
	var repos []*RepositoryInfo
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		githubRepos, resp, err := c.client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy repositories của user hiện tại: %w", err)
		}

		for _, repo := range githubRepos {
			repoInfo := &RepositoryInfo{
				FullName: repo.GetFullName(),
				Owner:    repo.GetOwner().GetLogin(),
				Name:     repo.GetName(),
				URL:      repo.GetHTMLURL(),
			}
			repos = append(repos, repoInfo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return repos, nil
}

// GetCurrentUserOrganizations retrieves current user organizations
func (c *Client) GetCurrentUserOrganizations(ctx context.Context) ([]string, error) {
	var orgs []string
	opts := &github.ListOptions{PerPage: 100}

	for {
		githubOrgs, resp, err := c.client.Organizations.List(ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy organizations: %w", err)
		}

		for _, org := range githubOrgs {
			orgs = append(orgs, org.GetLogin())
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return orgs, nil
}

// GetAllUserAndOrgRepositories retrieves all repositories of user and organizations
func (c *Client) GetAllUserAndOrgRepositories(ctx context.Context) ([]*RepositoryInfo, error) {
	var allRepos []*RepositoryInfo
	repoMap := make(map[string]bool) // To avoid duplicates

	// Get user repositories
	fmt.Println("📦 Quét repositories của user...")
	userRepos, err := c.GetCurrentUserRepositories(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy user repositories: %w", err)
	}
	fmt.Printf("   ✓ Tìm được %d repositories của user\n", len(userRepos))

	for _, repo := range userRepos {
		if !repoMap[repo.FullName] {
			allRepos = append(allRepos, repo)
			repoMap[repo.FullName] = true
		}
	}

	// Get organizations
	fmt.Println("🏢 Lấy danh sách organizations...")
	orgs, err := c.GetCurrentUserOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy organizations: %w", err)
	}
	fmt.Printf("   ✓ Tìm được %d organizations\n", len(orgs))

	// Get repositories from each organization
	if len(orgs) > 0 {
		fmt.Println("📦 Quét repositories từ organizations...")
	}
	for _, org := range orgs {
		fmt.Printf("   🔄 Đang quét %s...\n", org)
		orgRepos, err := c.GetOrganizationRepositories(ctx, org)
		if err != nil {
			fmt.Printf("   ⚠️  Lỗi khi lấy repositories từ %s: %v\n", org, err)
			continue
		}
		if len(orgRepos) == 0 {
			fmt.Printf("   ⚠️  %s: Không tìm thấy repositories (kiểm tra quyền truy cập)\n", org)
		} else {
			fmt.Printf("   ✓ %s: %d repositories\n", org, len(orgRepos))
		}

		for _, repo := range orgRepos {
			if !repoMap[repo.FullName] {
				allRepos = append(allRepos, repo)
				repoMap[repo.FullName] = true
			}
		}
	}

	return allRepos, nil
}

// GetPullRequestReviewsConcurrent retrieves reviews for multiple PRs concurrently
func (c *Client) GetPullRequestReviewsConcurrent(ctx context.Context, owner, repo string, prNumbers []int, maxWorkers int) (map[int][]*ReviewData, error) {
	if maxWorkers <= 0 {
		maxWorkers = 5 // Default worker pool size
	}

	results := make(map[int][]*ReviewData)
	resultsMutex := &sync.Mutex{}

	// Create semaphore to limit concurrent requests
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, prNumber := range prNumbers {
		wg.Add(1)
		go func(prNum int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			reviews, err := c.GetPullRequestReviews(ctx, owner, repo, prNum)
			if err != nil {
				fmt.Printf("⚠️  Error fetching reviews for PR #%d: %v\n", prNum, err)
				return
			}

			resultsMutex.Lock()
			results[prNum] = reviews
			resultsMutex.Unlock()
		}(prNumber)
	}

	wg.Wait()
	return results, nil
}

// GetPullRequestsFromRepositoriesConcurrent fetches PRs from multiple repositories concurrently
func (c *Client) GetPullRequestsFromRepositoriesConcurrent(ctx context.Context, repos []string, startDate, endDate time.Time, maxWorkers int) ([]RepositoryScanJob, error) {
	if maxWorkers <= 0 {
		maxWorkers = 3 // Default worker pool size for repo scanning
	}

	results := make([]RepositoryScanJob, 0)
	resultsMutex := &sync.Mutex{}

	// Create semaphore to limit concurrent requests
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, repoStr := range repos {
		// Parse repo string
		var owner, repoName string
		parts := len(repoStr)
		for i := 0; i < parts; i++ {
			if repoStr[i] == '/' {
				owner = repoStr[:i]
				repoName = repoStr[i+1:]
				break
			}
		}

		if owner == "" || repoName == "" {
			results = append(results, platform.RepositoryScanJob{
				Owner:    owner,
				RepoName: repoName,
				Error:    fmt.Errorf("invalid repository format: %s", repoStr),
			})
			continue
		}

		wg.Add(1)
		go func(o, r string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			prs, err := c.GetPullRequests(ctx, o, r, startDate, endDate)
			job := platform.RepositoryScanJob{
				Owner:    o,
				RepoName: r,
				PRData:   prs,
				Error:    err,
			}

			resultsMutex.Lock()
			results = append(results, job)
			resultsMutex.Unlock()
		}(owner, repoName)
	}

	wg.Wait()
	return results, nil
}

// GetPullRequestsWithReviewsConcurrent retrieves PRs with reviews concurrently
func (c *Client) GetPullRequestsWithReviewsConcurrent(ctx context.Context, owner, repo string, startDate, endDate time.Time, maxWorkers int) ([]*PullRequestData, error) {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	prs, err := c.GetPullRequests(ctx, owner, repo, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Extract PR numbers
	prNumbers := make([]int, len(prs))
	for i, pr := range prs {
		prNumbers[i] = pr.Number
	}

	// Fetch reviews concurrently
	reviewsMap, err := c.GetPullRequestReviewsConcurrent(ctx, owner, repo, prNumbers, maxWorkers)
	if err != nil {
		return nil, err
	}

	// Attach reviews to PRs
	for _, pr := range prs {
		if reviews, exists := reviewsMap[pr.Number]; exists {
			pr.Reviews = reviews
		}
	}

	return prs, nil
}
