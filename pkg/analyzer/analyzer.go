package analyzer

import (
	"regexp"
	"strings"

	"github.com/bug-crawler/pkg/github"
)

// PRRuleKeywords defines required keywords for PR rules
var (
	DescriptionKeywords = []string{
		"Description",
		"Changes Made",
		"Self-Review",
		"Functionality",
		"Security",
		"Error Handling",
		"Code Style",
		"Dependencies",
	}

	ReviewCommentKeywords = []string{
		"Functionality",
		"Security",
		"Error Handling",
		"Code Style",
		"Code Readability",
	}
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

// PRRuleResult chứa kết quả phân tích PR theo quy tắc code review
type PRRuleResult struct {
	PR                    *github.PullRequestData
	PRDescriptionValid    bool // Có đủ 7 keywords trong description
	ReviewCommentValid    bool // Review comment có đủ 4 keywords
	PRCompliant           bool // Tuân thủ đầy đủ tất cả quy tắc
	MissingDescKeywords   []string
	MissingReviewKeywords []string
}

// PRRuleAnalyzer phân tích PR theo quy tắc code review
type PRRuleAnalyzer struct{}

// NewPRRuleAnalyzer khởi tạo PRRuleAnalyzer
func NewPRRuleAnalyzer() *PRRuleAnalyzer {
	return &PRRuleAnalyzer{}
}

// CheckKeywordsInText kiểm tra xem text có chứa keywords (case-insensitive)
// Hỗ trợ cả full keywords và abbreviated tags (D1, DK1, etc)
func (pra *PRRuleAnalyzer) CheckKeywordsInText(text string, keywords []string) (bool, []string) {
	textLower := strings.ToLower(text)
	var missing []string

	// Regex patterns to match keywords and abbreviated tags for description
	patterns := map[string]*regexp.Regexp{
		"Description":    regexp.MustCompile(`(?i:description|desc|d\d)`),
		"Changes Made":   regexp.MustCompile(`(?i:changes made|changes|cm\d|change\d)`),
		"Self-Review":    regexp.MustCompile(`(?i:self-review|self review|sr\d)`),
		"Functionality":  regexp.MustCompile(`(?i:functionality|f\d)`),
		"Security":       regexp.MustCompile(`(?i:security|s\d)`),
		"Error Handling": regexp.MustCompile(`(?i:error handling|eh\d)`),
		"Code Style":     regexp.MustCompile(`(?i:code style|readability|c\d)`),
	}

	for _, keyword := range keywords {
		if pattern, exists := patterns[keyword]; exists {
			if !pattern.MatchString(textLower) {
				missing = append(missing, keyword)
			}
		} else {
			// Fallback to substring matching if pattern not defined
			keywordLower := strings.ToLower(keyword)
			if !strings.Contains(textLower, keywordLower) {
				missing = append(missing, keyword)
			}
		}
	}

	// Allow as long as text is not empty
	return true, missing
}

// AnalyzePRRule phân tích một PR theo quy tắc code review
func (pra *PRRuleAnalyzer) AnalyzePRRule(pr *github.PullRequestData) *PRRuleResult {
	result := &PRRuleResult{
		PR:                    pr,
		PRDescriptionValid:    false,
		ReviewCommentValid:    false,
		PRCompliant:           false,
		MissingDescKeywords:   []string{},
		MissingReviewKeywords: []string{},
	}

	// Bước 1: Kiểm tra PR Description có đủ keywords
	valid, missing := pra.CheckKeywordsInText(pr.Description, DescriptionKeywords)
	result.PRDescriptionValid = valid
	result.MissingDescKeywords = missing

	// Bước 2: Kiểm tra Review Comment
	valid, missing = pra.CheckReviewComments(pr.Reviews)
	result.ReviewCommentValid = valid
	result.MissingReviewKeywords = missing

	// Bước 3: Xác định PR compliant
	result.PRCompliant = result.PRDescriptionValid && result.ReviewCommentValid

	return result
}

// CheckReviewComments kiểm tra review comments có đủ keywords (case-insensitive)
// Yêu cầu: >2 keywords (tối thiểu 3 keywords), hỗ trợ cả abbreviated tags (F1, S1, etc)
func (pra *PRRuleAnalyzer) CheckReviewComments(reviews []*github.ReviewData) (bool, []string) {
	if len(reviews) == 0 {
		return false, ReviewCommentKeywords
	}

	// Gộp tất cả comments từ các reviewers
	allComments := ""
	for _, review := range reviews {
		if review.CommentBody != "" {
			allComments += " " + review.CommentBody
		}
	}

	if allComments == "" {
		return false, ReviewCommentKeywords
	}

	textLower := strings.ToLower(allComments)
	var missing []string

	// Regex patterns to match keywords and abbreviated tags for review comments
	patterns := map[string]*regexp.Regexp{
		"Functionality":    regexp.MustCompile(`(?i:functionality|f\d)`),
		"Security":         regexp.MustCompile(`(?i:security|s\d)`),
		"Error Handling":   regexp.MustCompile(`(?i:error handling|eh\d)`),
		"Code Style":       regexp.MustCompile(`(?i:code style|c\d)`),
		"Code Readability": regexp.MustCompile(`(?i:readability|code readability|cr\d)`),
	}

	for _, keyword := range ReviewCommentKeywords {
		if pattern, exists := patterns[keyword]; exists {
			if !pattern.MatchString(textLower) {
				missing = append(missing, keyword)
			}
		} else {
			// Fallback to substring matching if pattern not defined
			keywordLower := strings.ToLower(keyword)
			if !strings.Contains(textLower, keywordLower) {
				missing = append(missing, keyword)
			}
		}
	}

	// Allow as long as >2 keywords are found (at least 3 keywords)
	matchedCount := len(ReviewCommentKeywords) - len(missing)
	return matchedCount > 2, missing
}

// AnalyzePRRules phân tích danh sách PR theo quy tắc code review
func (pra *PRRuleAnalyzer) AnalyzePRRules(prs []*github.PullRequestData) []*PRRuleResult {
	results := make([]*PRRuleResult, 0)
	for _, pr := range prs {
		result := pra.AnalyzePRRule(pr)
		results = append(results, result)
	}
	return results
}
