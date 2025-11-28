package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v56/github"
)

// Client wraps GitHub API client
type Client struct {
	client *github.Client
}

// ReviewData chứa thông tin review của một reviewer
type ReviewData struct {
	ReviewerLogin string // Người review
	State         string // "APPROVED", "COMMENTED", "CHANGES_REQUESTED", "PENDING"
	SubmittedAt   *time.Time
	CommentBody   string // Nội dung comment của reviewer
}

// PullRequestData chứa thông tin PR cần thiết
type PullRequestData struct {
	Number      int
	Title       string
	Description string
	Author      string
	CreatedAt   time.Time
	MergedAt    *time.Time
	Labels      []string
	HTMLURL     string
	Status      string        // "open" or "merged"
	Reviews     []*ReviewData // Danh sách reviews
}

// NewClient khởi tạo GitHub client
func NewClient(token string) (*Client, error) {
	client := github.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}

	return &Client{
		client: client,
	}, nil
}

// GetPullRequests lấy danh sách PR trong khoảng thời gian
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

			prData := &PullRequestData{
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

// GetPullRequestReviews lấy danh sách reviews của một PR (bao gồm cả issue comments)
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*ReviewData, error) {
	var reviews []*ReviewData
	opts := &github.ListOptions{PerPage: 100}

	// Lấy reviews từ PR review API
	for {
		githubReviews, resp, err := c.client.PullRequests.ListReviews(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy reviews từ PR %d: %w", prNumber, err)
		}

		for _, review := range githubReviews {
			reviewData := &ReviewData{
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

	// Lấy thêm comments từ issue comments API (bao gồm review comments)
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
			// Nếu comment có nội dung, thêm vào reviews
			if comment.GetBody() != "" {
				reviewData := &ReviewData{
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

// GetPullRequestsWithReviews lấy danh sách PR cùng với review data của từng PR
func (c *Client) GetPullRequestsWithReviews(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*PullRequestData, error) {
	prs, err := c.GetPullRequests(ctx, owner, repo, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Lấy reviews cho mỗi PR
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

// VerifyToken kiểm tra token hợp lệ và hiển thị thông tin scopes
func (c *Client) VerifyToken(ctx context.Context) error {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return err
	}

	fmt.Printf("👤 Đăng nhập thành công với: %s\n", user.GetLogin())

	// Lấy thông tin rate limit
	rateLimits, _, err := c.client.RateLimits(ctx)
	if err == nil {
		fmt.Printf("📊 Rate limit: %d/%d requests\n",
			rateLimits.Core.Remaining,
			rateLimits.Core.Limit)
	}

	return nil
}

// RepositoryInfo chứa thông tin repository
type RepositoryInfo struct {
	FullName string
	Owner    string
	Name     string
	URL      string
}

// GetUserRepositories lấy danh sách repositories của user
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

// GetOrganizationRepositories lấy danh sách repositories của organization
func (c *Client) GetOrganizationRepositories(ctx context.Context, orgName string) ([]*RepositoryInfo, error) {
	var repos []*RepositoryInfo
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Type:        "all", // Lấy cả public, private, internal
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

// GetCurrentUserRepositories lấy repositories của user hiện tại
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

// GetCurrentUserOrganizations lấy danh sách organizations của user hiện tại
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

// GetAllUserAndOrgRepositories lấy tất cả repositories của user và các organizations
func (c *Client) GetAllUserAndOrgRepositories(ctx context.Context) ([]*RepositoryInfo, error) {
	var allRepos []*RepositoryInfo
	repoMap := make(map[string]bool) // Để tránh trùng lặp

	// Lấy repositories của user hiện tại
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

	// Lấy danh sách organizations
	fmt.Println("🏢 Lấy danh sách organizations...")
	orgs, err := c.GetCurrentUserOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy organizations: %w", err)
	}
	fmt.Printf("   ✓ Tìm được %d organizations\n", len(orgs))

	// Lấy repositories từ mỗi organization
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
