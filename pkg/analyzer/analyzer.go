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

const (
	MinKeywordsDescription   = 2 // Number of keywords required in PR description
	MinKeywordsReviewComment = 2 // Number of keywords required in PR review comment
)

// BugAnalyzer phân tích PR để detect bug
type BugAnalyzer struct {
	bugLabelRegex *regexp.Regexp
}

// BugResult kết quả phân tích PR
type BugResult struct {
	PR             *github.PullRequestData
	IsBugRelated   bool
	DetectionType  string // "bug_review", "bug"
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
	PRDescriptionValid    bool // Có đủ <MaxKeywordsDescription> keywords trong description
	ReviewCommentValid    bool // Review comment có đủ <MaxKeywordsReviewComment> keywords
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

// CheckKeywordsInText checks if the given text contains a sufficient number of specified keywords.
// It uses regular expressions for specific keywords and falls back to substring matching for others.
//
// Parameters:
//
//	text: The input string to be analyzed (e.g., PR description).
//	keywords: A slice of strings representing the keywords to look for.
//
// Returns:
//
//	bool: True if the number of matched keywords is greater than MinKeywordsDescription, false otherwise.
//	[]string: A slice of keywords that were not found in the text.
func (pra *PRRuleAnalyzer) CheckKeywordsInText(text string, keywords []string) (bool, []string) {
	textLower := strings.ToLower(text)
	var missing []string

	// patterns defines regular expressions for specific keywords and their abbreviated forms.
	// This allows for flexible matching beyond exact substring comparison.
	patterns := map[string]*regexp.Regexp{
		"Description":    regexp.MustCompile(`(?i:description|desc|d\d)`),
		"Changes Made":   regexp.MustCompile(`(?i:changes made|changes|cm\d|change\d)`),
		"Self-Review":    regexp.MustCompile(`(?i:self-review|self review|sr\d)`),
		"Functionality":  regexp.MustCompile(`(?i:functionality|f\d)`),
		"Security":       regexp.MustCompile(`(?i:security|s\d)`),
		"Error Handling": regexp.MustCompile(`(?i:error handling|eh\d)`),
		"Code Style":     regexp.MustCompile(`(?i:code style|readability|c\d)`),
		"Dependencies":   regexp.MustCompile(`(?i:dependencies|dep\d)`),
	}

	for _, keyword := range keywords {
		if pattern, exists := patterns[keyword]; exists {
			if !pattern.MatchString(textLower) {
				missing = append(missing, keyword)
			}
		} else {
			// Fallback to substring matching if a specific regex pattern is not defined for the keyword.
			keywordLower := strings.ToLower(keyword)
			if !strings.Contains(textLower, keywordLower) {
				missing = append(missing, keyword)
			}
		}
	}

	// Determine if the number of found keywords meets the minimum requirement.
	// MinKeywordsDescription is typically 3, meaning at greater than 2 keywords are required.
	matchedCount := len(keywords) - len(missing)
	return matchedCount > MinKeywordsDescription, missing
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

// CheckReviewComments analyzes the provided review comments to ensure they contain a sufficient number of predefined keywords.
// It aggregates all comments from reviewers and checks for the presence of keywords,
// supporting both direct substring matching and regular expression patterns for keywords and their abbreviated tags.
//
// Parameters:
//
//	reviews []*github.ReviewData: A slice of review data, each potentially containing a comment body.
//
// Returns:
//
//	bool: True if the aggregated review comments contain more than `MinKeywordsReviewComment` (typically 3, meaning at more than 3)
//	      of the `ReviewCommentKeywords`, indicating compliance. False otherwise.
//	[]string: A slice of keywords from `ReviewCommentKeywords` that were not found in the aggregated review comments.
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

	// Allow as long as > MinKeywordsReviewComment keywords are found
	matchedCount := len(ReviewCommentKeywords) - len(missing)
	return matchedCount > MinKeywordsReviewComment, missing
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
