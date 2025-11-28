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
	MinKeywordsDescription   = 3 // Number of keywords required in PR description
	MinKeywordsReviewComment = 3 // Number of keywords required in PR review comment
)

// BugAnalyzer analyzes a PR to detect bug
type BugAnalyzer struct {
	bugLabelRegex *regexp.Regexp
}

// BugResult contains the result of analyzing a PR to detect bug
type BugResult struct {
	PR             *github.PullRequestData
	IsBugRelated   bool
	DetectionType  string // "bug_review", "bug"
	MatchedKeyword string
	BugCount       int // Number of bugs from bug_review tag
}

// NewBugAnalyzer initializes a BugAnalyzer
func NewBugAnalyzer() *BugAnalyzer {
	return &BugAnalyzer{
		bugLabelRegex: regexp.MustCompile(`(?i:bug|fix|hotfix|critical|error|issue)`),
	}
}

// AnalyzePR analyzes a PR to detect bug
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
		// Check bug_review tag
		bugCount, found := ba.extractBugReviewCount(descLower)
		if found {
			result.IsBugRelated = true
			result.BugCount = bugCount
			result.DetectionType = "bug_review"
			result.MatchedKeyword = "bug_review"
		}
		return result
	default:
		// Check labels: bug, fix, hotfix, critical, error, issue
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

// extractBugReviewCount extracts the bug count from the description
func (ba *BugAnalyzer) extractBugReviewCount(desc string) (int, bool) {
	// Find pattern: bug_review: <number>
	re := regexp.MustCompile(`bug_review:\s*(\d+)`)
	matches := re.FindStringSubmatch(desc)

	if len(matches) >= 2 {
		count := 0
		// Parse number from string
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

// AnalyzePRs analyzes a list of PRs
func (ba *BugAnalyzer) AnalyzePRs(prs []*github.PullRequestData, bugType string) []*BugResult {
	results := make([]*BugResult, 0)
	for _, pr := range prs {
		result := ba.AnalyzePR(pr, bugType)
		results = append(results, result)
	}
	return results
}

// GetBugCount returns the number of PRs related to bug
func (ba *BugAnalyzer) GetBugCount(results []*BugResult) int {
	count := 0
	for _, result := range results {
		if result.IsBugRelated {
			count++
		}
	}
	return count
}

// PRRuleResult contains the result of analyzing a PR based on code review rules
type PRRuleResult struct {
	PR                 *github.PullRequestData
	PRDescriptionValid bool // PR description contains at least <MinKeywordsDescription> keywords
	ReviewCommentValid bool // Review comment contains at least <MinKeywordsReviewComment> keywords
	PRCompliant        bool // PR complies with all rules
}

// PRRuleAnalyzer analyzes a PR based on code review rules
type PRRuleAnalyzer struct{}

// NewPRRuleAnalyzer creates a new PRRuleAnalyzer
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
func (pra *PRRuleAnalyzer) CheckKeywordsInText(text string, keywords []string) bool {
	textLower := strings.ToLower(text)

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

	var matched []string

	for _, keyword := range keywords {
		if pattern, exists := patterns[keyword]; exists {
			if pattern.MatchString(textLower) {
				matched = append(matched, keyword)
			}
		} else {
			// Fallback to substring matching if a specific regex pattern is not defined for the keyword.
			keywordLower := strings.ToLower(keyword)
			if strings.Contains(textLower, keywordLower) {
				matched = append(matched, keyword)
			}
		}
	}

	// Determine if the number of found keywords meets the minimum requirement.
	// MinKeywordsDescription is typically 3, meaning at greater than 2 keywords are required.
	matchedCount := len(matched)
	return matchedCount >= MinKeywordsDescription
}

// AnalyzePRRule analyzes a PR based on code review rules
func (pra *PRRuleAnalyzer) AnalyzePRRule(pr *github.PullRequestData) *PRRuleResult {
	result := &PRRuleResult{
		PR:                 pr,
		PRDescriptionValid: false,
		ReviewCommentValid: false,
		PRCompliant:        false,
	}

	// Check if PR description contains at least MinKeywordsDescription keywords
	valid := pra.CheckKeywordsInText(pr.Description, DescriptionKeywords)
	result.PRDescriptionValid = valid

	// Check if review comments contain at least MinKeywordsReviewComment keywords
	result.ReviewCommentValid = valid

	// Determine if the PR complies with all rules
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
func (pra *PRRuleAnalyzer) CheckReviewComments(reviews []*github.ReviewData) bool {
	if len(reviews) == 0 {
		return false
	}

	// Aggregate all comments from reviewers
	allComments := ""
	for _, review := range reviews {
		if review.CommentBody != "" {
			allComments += " " + review.CommentBody
		}
	}

	if allComments == "" {
		return false
	}

	textLower := strings.ToLower(allComments)

	var matched []string

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
			if pattern.MatchString(textLower) {
				matched = append(matched, keyword)
			}
		} else {
			// Fallback to substring matching if pattern not defined
			keywordLower := strings.ToLower(keyword)
			if strings.Contains(textLower, keywordLower) {
				matched = append(matched, keyword)
			}
		}
	}

	// Allow as long as > MinKeywordsReviewComment keywords are found
	matchedCount := len(matched)
	return matchedCount >= MinKeywordsReviewComment
}

// AnalyzePRRules analyzes a list of PRs based on code review rules
func (pra *PRRuleAnalyzer) AnalyzePRRules(prs []*github.PullRequestData) []*PRRuleResult {
	results := make([]*PRRuleResult, 0)
	for _, pr := range prs {
		result := pra.AnalyzePRRule(pr)
		results = append(results, result)
	}
	return results
}
