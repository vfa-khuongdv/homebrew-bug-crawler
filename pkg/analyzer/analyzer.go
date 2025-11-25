package analyzer

import (
	"regexp"
	"strings"

	"github.com/bug-crawler/pkg/github"
)

// BugAnalyzer phân tích PR để detect bug
type BugAnalyzer struct {
	bugLabelRegex *regexp.Regexp
}

// BugResult kết quả phân tích PR
type BugResult struct {
	PR             *github.PullRequestData
	IsBugRelated   bool
	DetectionType  string // "bug_review", "keyword", "label", "both"
	MatchedKeyword string
	BugCount       int // Số bugs từ bug_review tag
}

// NewBugAnalyzer khởi tạo BugAnalyzer
func NewBugAnalyzer() *BugAnalyzer {
	return &BugAnalyzer{
		bugLabelRegex: regexp.MustCompile(`(?i:bug|fix|hotfix|critical|error|issue)`),
	}
}

// AnalyzePR phân tích một PR để detect bug
func (ba *BugAnalyzer) AnalyzePR(pr *github.PullRequestData, bugType string) *BugResult {
	result := &BugResult{
		PR:            pr,
		IsBugRelated:  false,
		DetectionType: "",
		BugCount:      0,
	}

	descLower := strings.ToLower(pr.Description)

	switch bugType {
	case "bug_review":
		// 1. Kiểm tra bug_review tag
		bugCount, found := ba.extractBugReviewCount(descLower)
		if found {
			result.IsBugRelated = true
			result.BugCount = bugCount
			result.DetectionType = "bug_review"
			result.MatchedKeyword = "bug_review"
		}
		return result
	default:
		// 2. Kiểm tra labels: bug, fix, hotfix, critical, error, issue
		for _, label := range pr.Labels {
			if ba.bugLabelRegex.MatchString(label) {
				result.IsBugRelated = true
				result.DetectionType = "label"
				result.MatchedKeyword = label
				break
			}
		}
		return result
	}
}

// extractBugReviewCount tìm pattern "bug_review: <number>" trong description
func (ba *BugAnalyzer) extractBugReviewCount(desc string) (int, bool) {
	// Tìm pattern: bug_review: <number>
	re := regexp.MustCompile(`bug_review:\s*(\d+)`)
	matches := re.FindStringSubmatch(desc)

	if len(matches) >= 2 {
		count := 0
		// Parse number từ string
		for _, ch := range matches[1] {
			if ch >= '0' && ch <= '9' {
				count = count*10 + int(ch-'0')
			}
		}

		if count > 0 {
			return count, true
		}
	}

	return 0, false
}

// AnalyzePRs phân tích danh sách PR
func (ba *BugAnalyzer) AnalyzePRs(prs []*github.PullRequestData, bugType string) []*BugResult {
	results := make([]*BugResult, 0)
	for _, pr := range prs {
		result := ba.AnalyzePR(pr, bugType)
		results = append(results, result)
	}
	return results
}

// GetBugCount lấy số lượng PR liên quan bug
func (ba *BugAnalyzer) GetBugCount(results []*BugResult) int {
	count := 0
	for _, result := range results {
		if result.IsBugRelated {
			count++
		}
	}
	return count
}
