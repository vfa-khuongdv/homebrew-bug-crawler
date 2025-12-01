package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bug-crawler/pkg/analyzer"
	"github.com/bug-crawler/pkg/auth"
	"github.com/bug-crawler/pkg/backlog"
	"github.com/bug-crawler/pkg/bitbucket"
	"github.com/bug-crawler/pkg/cli"
	"github.com/bug-crawler/pkg/github"
	"github.com/bug-crawler/pkg/platform"
	"github.com/bug-crawler/pkg/report"
)

func main() {
	fmt.Println("🐛 Bug Crawler - Multi-Platform PR Bug Analysis Tool")
	fmt.Println("==========================================")

	// Initialize managers
	tokenMgr := auth.NewTokenManager()
	cliTool := cli.NewCLI()
	ctx := context.Background()

	// Step 0: Select Platform
	fmt.Println("\nStep 0: Chọn Platform")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	selectedPlatform, err := cliTool.PromptSelectPlatform()
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn platform:", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Đã chọn: %s\n", strings.ToUpper(selectedPlatform))

	// Step 1: Get Token and Platform-Specific Credentials
	fmt.Println("\nStep 1: Xác Thực")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	var token, email, spaceID, domain string

	// Try to get saved token
	savedToken, err := tokenMgr.GetTokenForPlatform(selectedPlatform)
	if err == nil {
		fmt.Printf("✓ Token đã được tìm thấy từ file config cho %s\n", selectedPlatform)
		token = savedToken

		// Get additional credentials if needed
		if selectedPlatform == "bitbucket" {
			email, _ = tokenMgr.GetBitbucketEmail()
		} else if selectedPlatform == "backlog" {
			spaceID, _ = tokenMgr.GetBacklogSpaceID()
			domain, _ = tokenMgr.GetBacklogDomain()
		}
	}

	// Prompt for missing credentials
	if token == "" {
		var promptLabel string
		switch selectedPlatform {
		case "github":
			promptLabel = "GitHub Personal Access Token"
		case "bitbucket":
			promptLabel = "Bitbucket API Token"
			fmt.Println("\n📝 Tạo API Token tại: https://bitbucket.org/account/settings/personal/api-tokens/")
			fmt.Println("   Chọn scopes: User (Read), Workspace (Read), Repository (Read), Pull Request (Read)")
		case "backlog":
			promptLabel = "Backlog API Key"
		}

		fmt.Printf("\nNhập %s:\n", promptLabel)

		var inputToken string
		var err error

		if selectedPlatform == "backlog" {
			inputToken, err = cliTool.PromptBacklogApiKey()
		} else {
			inputToken, err = cliTool.PromptToken(promptLabel)
		}

		if err != nil {
			fmt.Println("❌ Lỗi khi nhập token:", err)
			os.Exit(1)
		}
		token = inputToken

		// Ask to save token
		if saveToken, err := cliTool.PromptSaveToken(); err == nil && saveToken {
			if err := tokenMgr.SaveTokenForPlatform(selectedPlatform, token); err != nil {
				fmt.Println("⚠️  Lỗi khi lưu token:", err)
			} else {
				fmt.Println("✓ Token đã được lưu")
			}
		}
	}

	// Get platform-specific additional credentials
	if selectedPlatform == "bitbucket" && email == "" {
		email, err = cliTool.PromptBitbucketEmail()
		if err != nil {
			fmt.Println("❌ Lỗi khi nhập email:", err)
			os.Exit(1)
		}
		_ = tokenMgr.SaveBitbucketEmail(email)
	} else if selectedPlatform == "backlog" && spaceID == "" {
		spaceID, err = cliTool.PromptBacklogSpaceID()
		if err != nil {
			fmt.Println("❌ Lỗi khi nhập space ID:", err)
			os.Exit(1)
		}
		_ = tokenMgr.SaveBacklogSpaceID(spaceID)
	}

	if selectedPlatform == "backlog" && domain == "" {
		domain, err = cliTool.PromptBacklogDomain()
		if err != nil {
			fmt.Println("❌ Lỗi khi chọn domain:", err)
			os.Exit(1)
		}
		_ = tokenMgr.SaveBacklogDomain(domain)
	}

	// Step 2: Initialize Platform Client
	fmt.Println("\nStep 2: Khởi Tạo Client")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	var platformClient platform.Platform
	switch selectedPlatform {
	case "github":
		platformClient, err = github.NewClient(token)
	case "bitbucket":
		platformClient, err = bitbucket.NewClient(email, token)
	case "backlog":
		platformClient, err = backlog.NewClient(spaceID, token, domain)
	default:
		fmt.Printf("❌ Platform không được hỗ trợ: %s\n", selectedPlatform)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("❌ Lỗi khi khởi tạo client:", err)
		os.Exit(1)
	}

	// Verify token
	if err := platformClient.VerifyToken(ctx); err != nil {
		fmt.Println("❌ Token không hợp lệ hoặc đã hết hạn:", err)
		os.Exit(1)
	}
	fmt.Println("✓ Token xác thực thành công")

	// Step 3: Select Scan Mode
	fmt.Println("\nStep 3: Chọn Chế Độ Scan")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	scanMode, err := cliTool.PromptSelectScanMode()
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn chế độ scan:", err)
		os.Exit(1)
	}

	// Step 4: Select Repositories
	fmt.Println("\nStep 4: Chọn Repositories")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	scanSource, err := cliTool.PromptSelectScanSource()
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn loại scan:", err)
		os.Exit(1)
	}

	var allRepos []*platform.RepositoryInfo

	if scanSource == "user" {
		fmt.Println("📦 Đang quét repositories của bạn...")
		userRepos, err := platformClient.GetCurrentUserRepositories(ctx)
		if err != nil {
			fmt.Println("❌ Lỗi khi quét repositories:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Tìm được %d repositories\n", len(userRepos))
		allRepos = userRepos
	} else {
		fmt.Println("🏢 Lấy danh sách organizations...")
		orgs, err := platformClient.GetCurrentUserOrganizations(ctx)
		if err != nil {
			fmt.Println("❌ Lỗi khi lấy organizations:", err)
			os.Exit(1)
		}

		if len(orgs) == 0 {
			fmt.Println("❌ Không tìm thấy organizations nào")
			os.Exit(1)
		}

		fmt.Printf("✓ Tìm được %d organizations\n", len(orgs))

		fmt.Println("\nChọn organizations để scan:")
		selectedOrgs, err := cliTool.PromptSelectOrganizations(orgs)
		if err != nil {
			fmt.Println("❌ Lỗi khi chọn organizations:", err)
			os.Exit(1)
		}

		fmt.Println("\n📦 Đang quét repositories từ organizations...")
		repoMap := make(map[string]bool)
		for _, org := range selectedOrgs {
			fmt.Printf("🔄 %s...\n", org)
			orgRepos, err := platformClient.GetOrganizationRepositories(ctx, org)
			if err != nil {
				fmt.Printf("⚠️  Lỗi: %v\n", err)
				continue
			}
			fmt.Printf("   ✓ %d repositories\n", len(orgRepos))

			for _, repo := range orgRepos {
				if !repoMap[repo.FullName] {
					allRepos = append(allRepos, repo)
					repoMap[repo.FullName] = true
				}
			}
		}
	}

	if len(allRepos) == 0 {
		fmt.Println("❌ Không tìm thấy repositories nào")
		os.Exit(1)
	}

	fmt.Println("\n" + strings.Repeat("-", 43))
	fmt.Printf("✓ Tổng cộng: %d repositories\n", len(allRepos))
	fmt.Println(strings.Repeat("-", 43))

	var repoNames []string
	for _, repo := range allRepos {
		repoNames = append(repoNames, repo.FullName)
	}

	selectedRepos, err := cliTool.PromptSelectMultipleRepositories(repoNames)
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn repositories:", err)
		os.Exit(1)
	}

	repos := selectedRepos
	if len(repos) == 0 {
		fmt.Println("❌ Vui lòng chọn ít nhất 1 repository")
		os.Exit(1)
	}

	fmt.Println("\n" + strings.Repeat("=", 43))
	fmt.Printf("📋 Repositories đã chọn (%d):\n", len(repos))
	fmt.Println(strings.Repeat("=", 43))
	for i, repo := range repos {
		fmt.Printf("%2d. ✓ %s\n", i+1, repo)
	}
	fmt.Println(strings.Repeat("=", 43))

	// Step 5: Select Date Range
	fmt.Println("\nStep 5: Chọn Khoảng Thời Gian")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	startDate, endDate, err := cliTool.PromptDateRange()
	if err != nil {
		fmt.Println("❌ Lỗi khi nhập ngày:", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Sẽ phân tích PR từ %s đến %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Step 6: Select Bug Type (if in bug detection mode)
	var bugType string
	if scanMode == "bug" {
		fmt.Println("\nStep 6: Chọn Loại Bug")
		fmt.Println("-" + strings.Repeat("-", 40) + "-")

		bugType, err = cliTool.PromptSelectBugType()
		if err != nil {
			fmt.Println("❌ Lỗi khi chọn loại bug:", err)
			os.Exit(1)
		}

		if bugType == "bug" {
			fmt.Println("✓ Sẽ scan bug từ labels")
		} else {
			fmt.Println("✓ Sẽ scan bug_review")
		}
	} else {
		fmt.Println("\nStep 6: Code Review Compliance Scan")
		fmt.Println("-" + strings.Repeat("-", 40) + "-")
		fmt.Println("✓ Sẽ scan PR theo quy tắc code review")
	}

	// Step 7: Crawler PR
	fmt.Println("\nStep 7: Crawler PR từ " + strings.ToUpper(selectedPlatform))
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	startTime := time.Now()
	bugAnalyzer := analyzer.NewBugAnalyzer()
	prRuleAnalyzer := analyzer.NewPRRuleAnalyzer()
	allResults := make([]*analyzer.BugResult, 0)
	allPRRuleResults := make([]*analyzer.PRRuleResult, 0)
	totalPRsCrawled := 0

	maxWorkers := 3
	if len(repos) < 3 {
		maxWorkers = len(repos)
	}
	if len(repos) > 10 {
		maxWorkers = 5
	}

	fmt.Printf("🚀 Quét %d repositories với %d workers (song song)...\n", len(repos), maxWorkers)

	scanJobs, err := platformClient.GetPullRequestsFromRepositoriesConcurrent(ctx, repos, startDate, endDate, maxWorkers)
	if err != nil {
		fmt.Printf("❌ Lỗi khi quét repositories: %v\n", err)
	}

	for _, job := range scanJobs {
		if job.Error != nil {
			fmt.Printf("❌ Lỗi khi lấy PR từ %s/%s: %v\n", job.Owner, job.RepoName, job.Error)
			continue
		}

		fmt.Printf("✓ %s/%s: %d PR\n", job.Owner, job.RepoName, len(job.PRData))
		totalPRsCrawled += len(job.PRData)

		if scanMode == "pr_rules" && len(job.PRData) > 0 {
			prNumbers := make([]int, len(job.PRData))
			for i, pr := range job.PRData {
				prNumbers[i] = pr.Number
			}

			reviewsMap, err := platformClient.GetPullRequestReviewsConcurrent(ctx, job.Owner, job.RepoName, prNumbers, 5)
			if err != nil {
				// Silently continue on error
			} else {
				for _, pr := range job.PRData {
					if reviews, exists := reviewsMap[pr.Number]; exists {
						pr.Reviews = reviews
					}
				}
			}
		}

		if scanMode == "pr_rules" {
			results := prRuleAnalyzer.AnalyzePRRules(job.PRData)
			allPRRuleResults = append(allPRRuleResults, results...)
		} else {
			results := bugAnalyzer.AnalyzePRs(job.PRData, bugType)
			allResults = append(allResults, results...)
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("✓ Hoàn thành crawl trong: %.2f giây\n", elapsedTime.Seconds())

	// Step 8: Report Results
	fmt.Println("\nStep 8: Thống Kê Kết Quả")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	reporter := report.NewReporter()

	if scanMode == "pr_rules" {
		reporter.PrintPRRulesSummary(allPRRuleResults)
		reporter.PrintPRRulesDetails(allPRRuleResults)

		csvFile := "pr_rules_report.csv"
		if err := reporter.ExportPRRulesCSV(csvFile, allPRRuleResults); err != nil {
			fmt.Printf("❌ Lỗi khi export CSV: %v\n", err)
		}
	} else {
		var filteredResults []*analyzer.BugResult
		switch bugType {
		case "bug_review":
			for _, result := range allResults {
				if result.DetectionType == "bug_review" {
					filteredResults = append(filteredResults, result)
				}
			}
		case "bug":
			for _, result := range allResults {
				if result.DetectionType == "label" {
					filteredResults = append(filteredResults, result)
				}
			}
		}

		stats := reporter.GenerateStatistics(filteredResults)
		stats.TotalPRsCrawled = totalPRsCrawled

		if stats.TotalPRsCrawled > 0 {
			stats.BugPercentage = float64(stats.BugRelatedPRs) * 100 / float64(stats.TotalPRsCrawled)
		}

		reporter.PrintSummary(stats)
		reporter.PrintDetails(stats)

		if stats.BugRelatedPRs > 0 {
			csvFile := "bug_report.csv"
			if err := reporter.ExportCSV(csvFile, stats); err != nil {
				fmt.Printf("❌ Lỗi khi export CSV: %v\n", err)
			}
		}
	}

	fmt.Println("\n✓ Hoàn thành!")
}
