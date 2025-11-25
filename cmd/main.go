package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bug-crawler/pkg/analyzer"
	"github.com/bug-crawler/pkg/auth"
	"github.com/bug-crawler/pkg/cli"
	githubclient "github.com/bug-crawler/pkg/github"
	"github.com/bug-crawler/pkg/report"
)

func main() {
	fmt.Println("🐛 Bug Crawler - GitHub PR Bug Analysis Tool")
	fmt.Println("==========================================")

	// 1. Quản lý token
	tokenMgr := auth.NewTokenManager()
	cliTool := cli.NewCLI()

	var token string
	fmt.Println("Step 1: GitHub Token")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	// Cố gắng lấy token từ environment hoặc file config
	if savedToken, err := tokenMgr.GetToken(""); err == nil {
		fmt.Println("✓ Token đã được tìm thấy từ file config")
		token = savedToken
	} else {
		// Yêu cầu user nhập token
		inputToken, err := cliTool.PromptToken()
		if err != nil {
			fmt.Println("❌ Lỗi khi nhập token:", err)
			os.Exit(1)
		}

		token = inputToken

		// Hỏi user có muốn lưu token không
		if saveToken, err := cliTool.PromptSaveToken(); err == nil && saveToken {
			if err := tokenMgr.SaveToken(token); err != nil {
				fmt.Println("⚠️  Lỗi khi lưu token:", err)
			} else {
				fmt.Println("✓ Token đã được lưu")
			}
		}
	}

	// 2. Khởi tạo GitHub client
	fmt.Println("\nStep 2: Xác thực GitHub")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	ctx := context.Background()
	ghClient, err := githubclient.NewClient(token)
	if err != nil {
		fmt.Println("❌ Lỗi khi khởi tạo GitHub client:", err)
		os.Exit(1)
	}

	// Kiểm tra token hợp lệ
	if err := ghClient.VerifyToken(ctx); err != nil {
		fmt.Println("❌ Token không hợp lệ hoặc đã hết hạn:", err)
		os.Exit(1)
	}
	fmt.Println("✓ Token xác thực thành công")

	// 3. Chọn loại scan và lấy repositories
	fmt.Println("\nStep 3: Chọn Repositories")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	// Chọn loại scan
	scanSource, err := cliTool.PromptSelectScanSource()
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn loại scan:", err)
		os.Exit(1)
	}

	var allRepos []*githubclient.RepositoryInfo

	if scanSource == "user" {
		// Scan repositories của user hiện tại
		fmt.Println("📦 Đang quét repositories của bạn...")
		userRepos, err := ghClient.GetCurrentUserRepositories(ctx)
		if err != nil {
			fmt.Println("❌ Lỗi khi quét repositories:", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Tìm được %d repositories\n", len(userRepos))
		allRepos = userRepos
	} else {
		// Scan repositories của organizations
		fmt.Println("🏢 Lấy danh sách organizations...")
		orgs, err := ghClient.GetCurrentUserOrganizations(ctx)
		if err != nil {
			fmt.Println("❌ Lỗi khi lấy organizations:", err)
			os.Exit(1)
		}

		if len(orgs) == 0 {
			fmt.Println("❌ Không tìm thấy organizations nào")
			os.Exit(1)
		}

		fmt.Printf("✓ Tìm được %d organizations\n", len(orgs))

		// Cho phép user chọn organizations
		fmt.Println("\nChọn organizations để scan:")
		selectedOrgs, err := cliTool.PromptSelectOrganizations(orgs)
		if err != nil {
			fmt.Println("❌ Lỗi khi chọn organizations:", err)
			os.Exit(1)
		}

		// Quét repositories từ các organizations đã chọn
		fmt.Println("\n📦 Đang quét repositories từ organizations...")
		repoMap := make(map[string]bool)
		for _, org := range selectedOrgs {
			fmt.Printf("🔄 %s...\n", org)
			orgRepos, err := ghClient.GetOrganizationRepositories(ctx, org)
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

	// Chuyển đổi repository objects thành chuỗi
	var repoNames []string
	for _, repo := range allRepos {
		repoNames = append(repoNames, repo.FullName)
	}

	// Cho phép user chọn repositories từ danh sách quét được
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

	// Hiển thị danh sách repositories đã chọn
	fmt.Println("\n" + strings.Repeat("=", 43))
	fmt.Printf("📋 Repositories đã chọn (%d):\n", len(repos))
	fmt.Println(strings.Repeat("=", 43))
	for i, repo := range repos {
		fmt.Printf("%2d. ✓ %s\n", i+1, repo)
	}
	fmt.Println(strings.Repeat("=", 43))

	// 4. Chọn khoảng thời gian
	fmt.Println("\nStep 4: Chọn Khoảng Thời Gian")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	startDate, endDate, err := cliTool.PromptDateRange()
	if err != nil {
		fmt.Println("❌ Lỗi khi nhập ngày:", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Sẽ phân tích PR từ %s đến %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// 5. Chọn loại bug để scan
	fmt.Println("\nStep 5: Chọn Loại Bug")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	bugType, err := cliTool.PromptSelectBugType()
	if err != nil {
		fmt.Println("❌ Lỗi khi chọn loại bug:", err)
		os.Exit(1)
	}

	if bugType == "bug" {
		fmt.Println("✓ Sẽ scan bug từ labels")
	} else {
		fmt.Println("✓ Sẽ scan bug_review")
	}

	// 6. Crawler PR
	fmt.Println("\nStep 6: Crawler PR từ GitHub")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	bugAnalyzer := analyzer.NewBugAnalyzer()
	allResults := make([]*analyzer.BugResult, 0)
	totalPRsCrawled := 0 // Đếm tổng số PR thực tế

	for _, repoStr := range repos {
		parts := strings.Split(repoStr, "/")
		if len(parts) != 2 {
			fmt.Printf("❌ Format repository không hợp lệ: %s\n", repoStr)
			continue
		}

		owner := parts[0]
		repoName := parts[1]

		fmt.Printf("Đang lấy PR từ %s/%s...\n", owner, repoName)

		prs, err := ghClient.GetPullRequests(ctx, owner, repoName, startDate, endDate)
		if err != nil {
			fmt.Printf("❌ Lỗi khi lấy PR từ %s/%s: %v\n", owner, repoName, err)
			continue
		}

		fmt.Printf("✓ Tìm được %d PR\n", len(prs))
		totalPRsCrawled += len(prs) // Cộng vào tổng

		// Phân tích PR
		results := bugAnalyzer.AnalyzePRs(prs, bugType)
		allResults = append(allResults, results...)
	}

	// Lọc kết quả theo loại bug đã chọn
	var filteredResults []*analyzer.BugResult
	switch bugType {
	case "bug_review":
		// Chỉ lấy PR có DetectionType là "bug_review"
		for _, result := range allResults {
			if result.DetectionType == "bug_review" {
				filteredResults = append(filteredResults, result)
			}
		}
	case "bug":
		// Chỉ lấy PR có DetectionType là "label" (bug từ labels)
		for _, result := range allResults {
			if result.DetectionType == "label" {
				filteredResults = append(filteredResults, result)
			}
		}
	}

	// 7. Thống kê và in báo cáo
	fmt.Println("\nStep 7: Thống Kê Kết Quả")
	fmt.Println("-" + strings.Repeat("-", 40) + "-")

	reporter := report.NewReporter()
	stats := reporter.GenerateStatistics(filteredResults)
	stats.TotalPRsCrawled = totalPRsCrawled // Ghi nhận tổng số PR thực tế

	// Tính lại BugPercentage dựa trên tổng PR thực tế được crawl
	if stats.TotalPRsCrawled > 0 {
		stats.BugPercentage = float64(stats.BugRelatedPRs) * 100 / float64(stats.TotalPRsCrawled)
	}

	reporter.PrintSummary(stats)
	reporter.PrintDetails(stats)

	// 8. Export CSV (optional)
	if stats.BugRelatedPRs > 0 {
		csvFile := "bug_report.csv"
		if err := reporter.ExportCSV(csvFile, stats); err != nil {
			fmt.Printf("❌ Lỗi khi export CSV: %v\n", err)
		}
	}

	fmt.Println("\n✓ Hoàn thành!")
}
