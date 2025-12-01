package backlog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bug-crawler/pkg/platform"
)

// Client wraps Backlog API client
type Client struct {
	httpClient *http.Client
	spaceID    string
	apiKey     string
	baseURL    string
}

// NewClient initializes Backlog client
func NewClient(spaceID, apiKey, domain string) (*Client, error) {
	if spaceID == "" || apiKey == "" {
		return nil, fmt.Errorf("space ID and API key are required")
	}

	if domain == "" {
		domain = "backlog.com"
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		spaceID:    spaceID,
		apiKey:     apiKey,
		baseURL:    fmt.Sprintf("https://%s.%s/api/v2", spaceID, domain),
	}, nil
}

// doRequest performs an HTTP request with API key authentication
func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values) ([]byte, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apiKey", c.apiKey)

	urlPath := fmt.Sprintf("%s%s?%s", c.baseURL, path, params.Encode())

	req, err := http.NewRequestWithContext(ctx, method, urlPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("backlog API error: %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// VerifyToken verifies API key validity
func (c *Client) VerifyToken(ctx context.Context) error {
	fmt.Printf("🔗 Connecting to: %s\n", c.baseURL)

	body, err := c.doRequest(ctx, "GET", "/space", nil)
	if err != nil {
		return err
	}

	var space struct {
		SpaceKey string `json:"spaceKey"`
		Name     string `json:"name"`
	}

	if err := json.Unmarshal(body, &space); err != nil {
		return err
	}

	fmt.Printf("👤 Đăng nhập thành công với space: %s (%s)\n", space.SpaceKey, space.Name)
	return nil
}

// GetCurrentUserRepositories retrieves Git repositories from all projects
func (c *Client) GetCurrentUserRepositories(ctx context.Context) ([]*platform.RepositoryInfo, error) {
	// Get all projects first
	body, err := c.doRequest(ctx, "GET", "/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []struct {
		ID         int    `json:"id"`
		ProjectKey string `json:"projectKey"`
		Name       string `json:"name"`
	}

	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, err
	}

	var allRepos []*platform.RepositoryInfo

	// Get repositories for each project
	for _, project := range projects {
		repos, err := c.getProjectRepositories(ctx, project.ProjectKey)
		if err != nil {
			fmt.Printf("⚠️  Lỗi khi lấy repositories từ project %s: %v\n", project.ProjectKey, err)
			continue
		}
		allRepos = append(allRepos, repos...)
	}

	return allRepos, nil
}

// getProjectRepositories retrieves repositories for a specific project
func (c *Client) getProjectRepositories(ctx context.Context, projectKey string) ([]*platform.RepositoryInfo, error) {
	path := fmt.Sprintf("/projects/%s/git/repositories", projectKey)
	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var repos []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}

	var repoInfos []*platform.RepositoryInfo
	for _, repo := range repos {
		repoInfo := &platform.RepositoryInfo{
			FullName: fmt.Sprintf("%s/%s", projectKey, repo.Name),
			Owner:    projectKey,
			Name:     repo.Name,
			URL:      fmt.Sprintf("https://%s.backlog.com/git/%s/%s", c.spaceID, projectKey, repo.Name),
		}
		repoInfos = append(repoInfos, repoInfo)
	}

	return repoInfos, nil
}

// GetOrganizationRepositories retrieves repositories for a specific project (Backlog uses projects instead of orgs)
func (c *Client) GetOrganizationRepositories(ctx context.Context, projectKey string) ([]*platform.RepositoryInfo, error) {
	return c.getProjectRepositories(ctx, projectKey)
}

// GetCurrentUserOrganizations retrieves all projects (equivalent to organizations)
func (c *Client) GetCurrentUserOrganizations(ctx context.Context) ([]string, error) {
	body, err := c.doRequest(ctx, "GET", "/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []struct {
		ProjectKey string `json:"projectKey"`
	}

	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, err
	}

	var projectKeys []string
	for _, project := range projects {
		projectKeys = append(projectKeys, project.ProjectKey)
	}

	return projectKeys, nil
}

// GetPullRequests retrieves pull requests within a time range
func (c *Client) GetPullRequests(ctx context.Context, projectKey, repoName string, startDate, endDate time.Time) ([]*platform.PullRequestData, error) {
	path := fmt.Sprintf("/projects/%s/git/repositories/%s/pullRequests", projectKey, repoName)

	params := url.Values{}
	params.Set("count", "100")

	body, err := c.doRequest(ctx, "GET", path, params)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy PR từ %s/%s: %w", projectKey, repoName, err)
	}

	var pullRequests []struct {
		Number      int    `json:"number"`
		Summary     string `json:"summary"`
		Description string `json:"description"`
		Status      struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		CreatedUser struct {
			Name string `json:"name"`
		} `json:"createdUser"`
		Created *time.Time `json:"created"`
		Updated *time.Time `json:"updated"`
		Merged  *time.Time `json:"merged"`
	}

	if err := json.Unmarshal(body, &pullRequests); err != nil {
		return nil, err
	}

	var prs []*platform.PullRequestData
	for _, pr := range pullRequests {
		// Filter by date range
		if pr.Created != nil && (pr.Created.Before(startDate) || pr.Created.After(endDate)) {
			continue
		}

		status := "open"
		if pr.Status.ID == 3 { // 3 = Merged in Backlog
			status = "merged"
		}

		createdAt := time.Now()
		if pr.Created != nil {
			createdAt = *pr.Created
		}

		prData := &platform.PullRequestData{
			Number:      pr.Number,
			Title:       pr.Summary,
			Description: pr.Description,
			Author:      pr.CreatedUser.Name,
			CreatedAt:   createdAt,
			MergedAt:    pr.Merged,
			Labels:      []string{}, // Backlog doesn't have labels on PRs
			HTMLURL:     fmt.Sprintf("https://%s.backlog.com/git/%s/%s/pullRequests/%d", c.spaceID, projectKey, repoName, pr.Number),
			Status:      status,
		}

		prs = append(prs, prData)
	}

	return prs, nil
}

// GetPullRequestReviews retrieves reviews/comments for a pull request
func (c *Client) GetPullRequestReviews(ctx context.Context, projectKey, repoName string, prNumber int) ([]*platform.ReviewData, error) {
	path := fmt.Sprintf("/projects/%s/git/repositories/%s/pullRequests/%d/comments", projectKey, repoName, prNumber)

	body, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi lấy reviews từ PR %d: %w", prNumber, err)
	}

	var comments []struct {
		Content     string     `json:"content"`
		Created     *time.Time `json:"created"`
		CreatedUser struct {
			Name string `json:"name"`
		} `json:"createdUser"`
	}

	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, err
	}

	var reviews []*platform.ReviewData
	for _, comment := range comments {
		reviewData := &platform.ReviewData{
			ReviewerLogin: comment.CreatedUser.Name,
			State:         "COMMENTED",
			SubmittedAt:   comment.Created,
			CommentBody:   comment.Content,
		}
		reviews = append(reviews, reviewData)
	}

	return reviews, nil
}

// GetPullRequestReviewsConcurrent retrieves reviews for multiple PRs concurrently
func (c *Client) GetPullRequestReviewsConcurrent(ctx context.Context, projectKey, repoName string, prNumbers []int, maxWorkers int) (map[int][]*platform.ReviewData, error) {
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

			reviews, err := c.GetPullRequestReviews(ctx, projectKey, repoName, prNum)
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

	// Limit number of concurrent workers
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, repoStr := range repos {
		// Parse repo string (format: project-key/repo-name)
		var projectKey, repoName string
		for i := 0; i < len(repoStr); i++ {
			if repoStr[i] == '/' {
				projectKey = repoStr[:i]
				repoName = repoStr[i+1:]
				break
			}
		}

		if projectKey == "" || repoName == "" {
			results = append(results, platform.RepositoryScanJob{
				Owner:    projectKey,
				RepoName: repoName,
				Error:    fmt.Errorf("invalid repository format: %s", repoStr),
			})
			continue
		}

		// Add to wait group
		wg.Add(1)
		go func(pk, rn string) {
			// Wait for semaphore
			defer wg.Done()
			// Limit number of concurrent workers
			semaphore <- struct{}{}
			// Release semaphore when done
			defer func() { <-semaphore }()

			prs, err := c.GetPullRequests(ctx, pk, rn, startDate, endDate)

			job := platform.RepositoryScanJob{
				Owner:    pk,
				RepoName: rn,
				PRData:   prs,
				Error:    err,
			}

			// Lock results mutex
			resultsMutex.Lock()
			// Append job to results
			results = append(results, job)
			// Unlock results mutex
			resultsMutex.Unlock()
		}(projectKey, repoName)
	}

	wg.Wait()
	return results, nil
}
