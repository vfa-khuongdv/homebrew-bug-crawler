package analyzer

import (
	"regexp"
	"strings"

	"github.com/bug-crawler/pkg/github"
)

// BugAnalyzer phân tích PR để detect bug
type BugAnalyzer struct {
	bugKeywords   []string
	fixKeywords   []string
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
		bugKeywords: []string{
			"bug", "fix", "hotfix", "patch",
			"crash", "error", "issue", "problem",
			"failed", "exception", "broken",
		},
		fixKeywords: []string{
			"fix", "resolve", "close", "closes",
			"fixed", "resolved", "patch", "hotfix",
		},
		bugLabelRegex: regexp.MustCompile(`(?i:bug|fix|hotfix|critical|error|issue)`),
	}
}

// AnalyzePR phân tích một PR để detect bug
func (ba *BugAnalyzer) AnalyzePR(pr *github.PullRequestData) *BugResult {
	result := &BugResult{
		PR:            pr,
		IsBugRelated:  false,
		DetectionType: "",
		BugCount:      0,
	}

	descLower := strings.ToLower(pr.Description)

	// Priority 1: Kiểm tra bug_review pattern
	bugCount, found := ba.extractBugReviewCount(descLower)
	if found {
		result.IsBugRelated = true
		result.BugCount = bugCount
		result.DetectionType = "bug_review"
		result.MatchedKeyword = "bug_review"
		return result
	}

	// Priority 2: Kiểm tra "type: bug" keyword
	if ba.checkTypeBugKeyword(descLower) {
		result.IsBugRelated = true
		result.DetectionType = "keyword"
		result.MatchedKeyword = "type: bug"
		return result
	}

	// Priority 3: Kiểm tra bug labels
	for _, label := range pr.Labels {
		if ba.bugLabelRegex.MatchString(label) {
			result.IsBugRelated = true
			result.DetectionType = "label"
			result.MatchedKeyword = label
			return result
		}
	}

	// Priority 4: Fallback - Kiểm tra general bug keywords
	for _, keyword := range ba.bugKeywords {
		if strings.Contains(descLower, strings.ToLower(keyword)) {
			result.IsBugRelated = true
			result.MatchedKeyword = keyword
			result.DetectionType = "keyword"
			return result
		}
	}

	return result
}

// checkTypeBugKeyword kiểm tra pattern "type: bug" trong text
func (ba *BugAnalyzer) checkTypeBugKeyword(text string) bool {
	// Tìm pattern: type: bug (case insensitive)
	re := regexp.MustCompile(`(?i)type:\s*bug`)
	return re.MatchString(text)
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
func (ba *BugAnalyzer) AnalyzePRs(prs []*github.PullRequestData) []*BugResult {
	results := make([]*BugResult, 0)
	for _, pr := range prs {
		result := ba.AnalyzePR(pr)
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
