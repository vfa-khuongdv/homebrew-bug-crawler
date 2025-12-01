package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bug-crawler/pkg/platform"
)

const (
	bitbucketAPIURL = "https://api.bitbucket.org/2.0"
)

// Client wraps Bitbucket API client
type Client struct {
	httpClient *http.Client
	email      string
	apiToken   string
}

// NewClient initializes Bitbucket client
func NewClient(email, apiToken string) (*Client, error) {
	if email == "" || apiToken == "" {
		return nil, fmt.Errorf("atlassian account email and API token are required")
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		email:      email,
		apiToken:   apiToken,
	}, nil
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, urlPath string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlPath, nil)
	if err != nil {
		return nil, err
	}

	// Basic authentication with Atlassian account email
	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, fmt.Errorf("bitbucket API error: 401 unauthorized - kiểm tra lại atlassian account email (%s) và API token. API token cần được tạo tại https://id.atlassian.com/manage-profile/security/api-tokens", c.email)
		}
		return nil, fmt.Errorf("bitbucket API error: %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// VerifyToken verifies token validity
func (c *Client) VerifyToken(ctx context.Context) error {
	urlPath := fmt.Sprintf("%s/user", bitbucketAPIURL)
	body, err := c.doRequest(ctx, "GET", urlPath)
	if err != nil {
		// Provide helpful message if verification fails
		if strings.Contains(err.Error(), "401") {
			return fmt.Errorf("authentication failed - verify your API token has these scopes: User (Read), Workspace (Read), Repository (Read), Pull Request (Read). Token may also be expired or revoked. See: https://id.atlassian.com/manage-profile/security/api-tokens")
		}
		return err
	}

	var user struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
	}

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	fmt.Printf("👤 Đăng nhập thành công với: %s (%s)\n", user.Username, user.DisplayName)
	return nil
}

// GetCurrentUserRepositories retrieves current user repositories
func (c *Client) GetCurrentUserRepositories(ctx context.Context) ([]*platform.RepositoryInfo, error) {
	var repos []*platform.RepositoryInfo

	// Get user's workspaces first
	workspaces, err := c.GetCurrentUserOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("⚠️  No workspaces found for current user")
		return repos, nil
	}

	fmt.Printf("📦 Found %d workspace(s), scanning for repositories...\n", len(workspaces))

	// For each workspace, get the repositories
	for _, workspace := range workspaces {
		workspaceRepos, err := c.GetOrganizationRepositories(ctx, workspace)
		if err != nil {
			fmt.Printf("⚠️  Failed to get repositories from workspace '%s': %v\n", workspace, err)
			continue
		}
		repos = append(repos, workspaceRepos...)
	}

	return repos, nil
}

// GetOrganizationRepositories retrieves organization (workspace) repositories
func (c *Client) GetOrganizationRepositories(ctx context.Context, workspace string) ([]*platform.RepositoryInfo, error) {
	var repos []*platform.RepositoryInfo
	urlPath := fmt.Sprintf("%s/repositories/%s", bitbucketAPIURL, workspace)

	for urlPath != "" {
		body, err := c.doRequest(ctx, "GET", urlPath)
		if err != nil {
			return nil, err
		}

		var response struct {
			Values []struct {
				FullName string `json:"full_name"`
				Name     string `json:"name"`
				Owner    struct {
					Username string `json:"username"`
				} `json:"owner"`
				Links struct {
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
			} `json:"values"`
			Next string `json:"next"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		for _, repo := range response.Values {
			repoInfo := &platform.RepositoryInfo{
				FullName: repo.FullName,
				Owner:    repo.Owner.Username,
				Name:     repo.Name,
				URL:      repo.Links.HTML.Href,
			}
			repos = append(repos, repoInfo)
		}

		urlPath = response.Next
	}

	return repos, nil
}

// GetCurrentUserOrganizations retrieves current user workspaces
func (c *Client) GetCurrentUserOrganizations(ctx context.Context) ([]string, error) {
	var workspaces []string
	urlPath := fmt.Sprintf("%s/workspaces", bitbucketAPIURL)

	for urlPath != "" {
		body, err := c.doRequest(ctx, "GET", urlPath)
		if err != nil {
			return nil, err
		}

		var response struct {
			Values []struct {
				Slug string `json:"slug"`
			} `json:"values"`
			Next string `json:"next"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		for _, workspace := range response.Values {
			workspaces = append(workspaces, workspace.Slug)
		}

		urlPath = response.Next
	}

	return workspaces, nil
}

// GetPullRequests retrieves pull requests within a time range
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string, startDate, endDate time.Time) ([]*platform.PullRequestData, error) {
	var prs []*platform.PullRequestData
	urlPath := fmt.Sprintf("%s/repositories/%s/%s/pullrequests?state=MERGED&state=OPEN", bitbucketAPIURL, owner, repo)

	for urlPath != "" {
		body, err := c.doRequest(ctx, "GET", urlPath)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy PR từ %s/%s: %w", owner, repo, err)
		}

		var response struct {
			Values []struct {
				ID          int    `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				State       string `json:"state"`
				Author      struct {
					DisplayName string `json:"display_name"`
				} `json:"author"`
				CreatedOn time.Time  `json:"created_on"`
				UpdatedOn time.Time  `json:"updated_on"`
				MergedOn  *time.Time `json:"closed_on"`
				Links     struct {
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
			} `json:"values"`
			Next string `json:"next"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		for _, pr := range response.Values {
			// Filter by date range
			if pr.CreatedOn.Before(startDate) || pr.CreatedOn.After(endDate) {
				continue
			}

			status := "open"
			if pr.State == "MERGED" {
				status = "merged"
			}

			prData := &platform.PullRequestData{
				Number:      pr.ID,
				Title:       pr.Title,
				Description: pr.Description,
				Author:      pr.Author.DisplayName,
				CreatedAt:   pr.CreatedOn,
				MergedAt:    pr.MergedOn,
				Labels:      []string{}, // Bitbucket doesn't have labels on PRs by default
				HTMLURL:     pr.Links.HTML.Href,
				Status:      status,
			}

			prs = append(prs, prData)
		}

		urlPath = response.Next
	}

	return prs, nil
}

// GetPullRequestReviews retrieves reviews for a pull request
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*platform.ReviewData, error) {
	var reviews []*platform.ReviewData
	urlPath := fmt.Sprintf("%s/repositories/%s/%s/pullrequests/%d/comments", bitbucketAPIURL, owner, repo, prNumber)

	for urlPath != "" {
		body, err := c.doRequest(ctx, "GET", urlPath)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi lấy reviews từ PR %d: %w", prNumber, err)
		}

		var response struct {
			Values []struct {
				Content struct {
					Raw string `json:"raw"`
				} `json:"content"`
				User struct {
					DisplayName string `json:"display_name"`
				} `json:"user"`
				CreatedOn time.Time `json:"created_on"`
			} `json:"values"`
			Next string `json:"next"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		for _, comment := range response.Values {
			reviewData := &platform.ReviewData{
				ReviewerLogin: comment.User.DisplayName,
				State:         "COMMENTED",
				SubmittedAt:   &comment.CreatedOn,
				CommentBody:   comment.Content.Raw,
			}
			reviews = append(reviews, reviewData)
		}

		urlPath = response.Next
	}

	return reviews, nil
}

// GetPullRequestReviewsConcurrent retrieves reviews for multiple PRs concurrently
func (c *Client) GetPullRequestReviewsConcurrent(ctx context.Context, owner, repo string, prNumbers []int, maxWorkers int) (map[int][]*platform.ReviewData, error) {
	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	results := make(map[int][]*platform.ReviewData)
	resultsMutex := &sync.Mutex{}

	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, prNumber := range prNumbers {
		wg.Add(1)
		go func(prNum int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

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
func (c *Client) GetPullRequestsFromRepositoriesConcurrent(ctx context.Context, repos []string, startDate, endDate time.Time, maxWorkers int) ([]platform.RepositoryScanJob, error) {
	if maxWorkers <= 0 {
		maxWorkers = 3
	}

	results := make([]platform.RepositoryScanJob, 0)
	resultsMutex := &sync.Mutex{}

	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, repoStr := range repos {
		// Parse repo string (format: owner/repo)
		var owner, repoName string
		for i := 0; i < len(repoStr); i++ {
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
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

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
