package analyzer

import (
	"testing"

	"github.com/bug-crawler/pkg/github"
)

func TestAnalyzePR_BugReview(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *github.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
		wantCount   int
	}{
		{
			name: "bug_review with count",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "This PR fixes bug_review: 3 issues",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   3,
		},
		{
			name: "bug_review with single bug",
			pr: &github.PullRequestData{
				Title:       "Fix issue",
				Description: "bug_review: 1",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   1,
		},
		{
			name: "bug_review case insensitive",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "BUG_REVIEW: 5",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   5,
		},
		{
			name: "bug_review with large count",
			pr: &github.PullRequestData{
				Title:       "Major fix",
				Description: "bug_review: 123",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   123,
		},
		{
			name: "bug_review with extra spaces",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "bug_review:     10",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   10,
		},
		{
			name: "bug_review with count 0 - not detected",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "bug_review: 0",
				Labels:      []string{},
			},
			wantBug:   false,
			wantType:  "",
			wantCount: 0,
		},
		{
			name: "bug_review without count - not detected",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "bug_review:",
				Labels:      []string{},
			},
			wantBug:   false,
			wantType:  "",
			wantCount: 0,
		},
		{
			name: "bug_review in title - ignored",
			pr: &github.PullRequestData{
				Title:       "bug_review: 5",
				Description: "Some description",
				Labels:      []string{},
			},
			wantBug:   false,
			wantType:  "",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr, "bug_review")

			if result.IsBugRelated != tt.wantBug {
				t.Errorf("IsBugRelated = %v, want %v", result.IsBugRelated, tt.wantBug)
			}
			if result.DetectionType != tt.wantType {
				t.Errorf("DetectionType = %v, want %v", result.DetectionType, tt.wantType)
			}
			if result.MatchedKeyword != tt.wantKeyword {
				t.Errorf("MatchedKeyword = %v, want %v", result.MatchedKeyword, tt.wantKeyword)
			}
			if result.BugCount != tt.wantCount {
				t.Errorf("BugCount = %v, want %v", result.BugCount, tt.wantCount)
			}
		})
	}
}

func TestAnalyzePR_Labels(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *github.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
	}{
		{
			name: "bug label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"bug"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "bug",
		},
		{
			name: "fix label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"enhancement", "fix"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "fix",
		},
		{
			name: "hotfix label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"hotfix"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "hotfix",
		},
		{
			name: "critical label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"critical"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "critical",
		},
		{
			name: "error label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"error"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "error",
		},
		{
			name: "issue label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"issue"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "issue",
		},
		{
			name: "BUG label case insensitive",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"BUG"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "BUG",
		},
		{
			name: "multiple labels - first match wins",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"enhancement", "bug", "fix"},
			},
			wantBug:     true,
			wantType:    "label",
			wantKeyword: "bug",
		},
		{
			name: "no matching label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"enhancement", "documentation"},
			},
			wantBug:  false,
			wantType: "",
		},
		{
			name: "empty labels",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{},
			},
			wantBug:  false,
			wantType: "",
		},
		{
			name: "nil labels",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      nil,
			},
			wantBug:  false,
			wantType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr, "label") // Use "label" or default

			if result.IsBugRelated != tt.wantBug {
				t.Errorf("IsBugRelated = %v, want %v", result.IsBugRelated, tt.wantBug)
			}
			if result.DetectionType != tt.wantType {
				t.Errorf("DetectionType = %v, want %v", result.DetectionType, tt.wantType)
			}
			if tt.wantBug && result.MatchedKeyword != tt.wantKeyword {
				t.Errorf("MatchedKeyword = %v, want %v", result.MatchedKeyword, tt.wantKeyword)
			}
		})
	}
}

func TestAnalyzePR_ExclusiveModes(t *testing.T) {
	analyzer := NewBugAnalyzer()

	pr := &github.PullRequestData{
		Title:       "Fix",
		Description: "bug_review: 3",
		Labels:      []string{"bug", "critical"},
	}

	t.Run("bug_review mode ignores labels", func(t *testing.T) {
		result := analyzer.AnalyzePR(pr, "bug_review")
		if !result.IsBugRelated {
			t.Error("Should be detected as bug related")
		}
		if result.DetectionType != "bug_review" {
			t.Errorf("Expected detection type 'bug_review', got '%s'", result.DetectionType)
		}
		if result.BugCount != 3 {
			t.Errorf("Expected bug count 3, got %d", result.BugCount)
		}
	})

	t.Run("label mode ignores bug_review", func(t *testing.T) {
		result := analyzer.AnalyzePR(pr, "label")
		if !result.IsBugRelated {
			t.Error("Should be detected as bug related")
		}
		if result.DetectionType != "label" {
			t.Errorf("Expected detection type 'label', got '%s'", result.DetectionType)
		}
		if result.BugCount != 0 {
			t.Errorf("Expected bug count 0 in label mode, got %d", result.BugCount)
		}
	})
}

func TestAnalyzePR_NotBugRelated(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name string
		pr   *github.PullRequestData
	}{
		{
			name: "feature PR",
			pr: &github.PullRequestData{
				Title:       "Add new feature",
				Description: "Implement user profile page",
				Labels:      []string{"enhancement"},
			},
		},
		{
			name: "documentation PR",
			pr: &github.PullRequestData{
				Title:       "Update README",
				Description: "Add installation instructions",
				Labels:      []string{"documentation"},
			},
		},
		{
			name: "refactor PR",
			pr: &github.PullRequestData{
				Title:       "Refactor code",
				Description: "Improve code structure",
				Labels:      []string{"refactor"},
			},
		},
		{
			name: "PR with bug keyword in description but no label (in label mode)",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "This fixes a bug in the system",
				Labels:      []string{},
			},
		},
		{
			name: "PR with fix keyword in title only",
			pr: &github.PullRequestData{
				Title:       "Fix typo in documentation",
				Description: "Some description",
				Labels:      []string{},
			},
		},
		{
			name: "empty PR",
			pr: &github.PullRequestData{
				Title:       "",
				Description: "",
				Labels:      []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr, "label")

			if result.IsBugRelated {
				t.Errorf("Expected IsBugRelated to be false, got true with type=%s, keyword=%s",
					result.DetectionType, result.MatchedKeyword)
			}
			if result.DetectionType != "" {
				t.Errorf("Expected DetectionType to be empty, got %v", result.DetectionType)
			}
			if result.BugCount != 0 {
				t.Errorf("Expected BugCount to be 0, got %v", result.BugCount)
			}
		})
	}
}

func TestAnalyzePRs_BugReviewMode(t *testing.T) {
	analyzer := NewBugAnalyzer()

	prs := []*github.PullRequestData{
		{
			Title:       "Fix bug",
			Description: "bug_review: 2",
			Labels:      []string{},
		},
		{
			Title:       "Add feature",
			Description: "New feature implementation",
			Labels:      []string{},
		},
		{
			Title:       "Hotfix",
			Description: "Critical issue",
			Labels:      []string{"hotfix"},
		},
	}

	results := analyzer.AnalyzePRs(prs, "bug_review")

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First PR should be detected
	if !results[0].IsBugRelated {
		t.Error("First PR should be bug-related")
	}
	if results[0].DetectionType != "bug_review" {
		t.Errorf("First PR should be detected by bug_review, got %s", results[0].DetectionType)
	}

	// Third PR (hotfix label) should NOT be detected in bug_review mode
	if results[2].IsBugRelated {
		t.Error("Third PR should NOT be bug-related in bug_review mode")
	}
}

func TestAnalyzePRs_LabelMode(t *testing.T) {
	analyzer := NewBugAnalyzer()

	prs := []*github.PullRequestData{
		{
			Title:       "Fix bug",
			Description: "bug_review: 2",
			Labels:      []string{},
		},
		{
			Title:       "Add feature",
			Description: "New feature implementation",
			Labels:      []string{},
		},
		{
			Title:       "Hotfix",
			Description: "Critical issue",
			Labels:      []string{"hotfix"},
		},
	}

	results := analyzer.AnalyzePRs(prs, "label")

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// First PR (bug_review) should NOT be detected in label mode
	if results[0].IsBugRelated {
		t.Error("First PR should NOT be bug-related in label mode")
	}

	// Third PR (hotfix label) should be detected
	if !results[2].IsBugRelated {
		t.Error("Third PR should be bug-related")
	}
	if results[2].DetectionType != "label" {
		t.Errorf("Third PR should be detected by label, got %s", results[2].DetectionType)
	}
}

func TestAnalyzePRs_EmptyList(t *testing.T) {
	analyzer := NewBugAnalyzer()

	results := analyzer.AnalyzePRs([]*github.PullRequestData{}, "label")

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestAnalyzePRs_NilList(t *testing.T) {
	analyzer := NewBugAnalyzer()

	results := analyzer.AnalyzePRs(nil, "label")

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestGetBugCount(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name      string
		results   []*BugResult
		wantCount int
	}{
		{
			name: "mixed results",
			results: []*BugResult{
				{IsBugRelated: true},
				{IsBugRelated: false},
				{IsBugRelated: true},
				{IsBugRelated: true},
				{IsBugRelated: false},
			},
			wantCount: 3,
		},
		{
			name: "all bugs",
			results: []*BugResult{
				{IsBugRelated: true},
				{IsBugRelated: true},
				{IsBugRelated: true},
			},
			wantCount: 3,
		},
		{
			name: "no bugs",
			results: []*BugResult{
				{IsBugRelated: false},
				{IsBugRelated: false},
			},
			wantCount: 0,
		},
		{
			name:      "empty results",
			results:   []*BugResult{},
			wantCount: 0,
		},
		{
			name:      "nil results",
			results:   nil,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := analyzer.GetBugCount(tt.results)
			if count != tt.wantCount {
				t.Errorf("Expected bug count %d, got %d", tt.wantCount, count)
			}
		})
	}
}

func TestNewBugAnalyzer(t *testing.T) {
	analyzer := NewBugAnalyzer()

	if analyzer == nil {
		t.Fatal("NewBugAnalyzer returned nil")
	}

	if analyzer.bugLabelRegex == nil {
		t.Error("bugLabelRegex should not be nil")
	}
}

// Tests for PRRuleAnalyzer

func TestCheckKeywordsInText(t *testing.T) {
	analyzer := NewPRRuleAnalyzer()

	tests := []struct {
		name        string
		text        string
		keywords    []string
		wantValid   bool
		wantMissing int
	}{
		{
			name:        "All keywords present",
			text:        "Description of the work. Changes Made to the code. Self-Review checklist completed. Functionality is working. Security is handled. Error Handling is in place. Code Style follows conventions.",
			keywords:    DescriptionKeywords,
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name:        "Missing one keyword",
			text:        "Description of the work. Changes Made to the code. Self-Review checklist. Functionality works. Security is good. Error Handling included.",
			keywords:    DescriptionKeywords,
			wantValid:   false,
			wantMissing: 1,
		},
		{
			name:        "Missing multiple keywords",
			text:        "Description provided. Changes Made documented.",
			keywords:    DescriptionKeywords,
			wantValid:   false,
			wantMissing: 5,
		},
		{
			name:        "Case insensitive matching",
			text:        "DESCRIPTION of work. CHANGES MADE to files. SELF-REVIEW done. FUNCTIONALITY tested. SECURITY checked. ERROR HANDLING implemented. CODE STYLE followed.",
			keywords:    DescriptionKeywords,
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name:        "Mixed case keywords",
			text:        "description here. Changes Made included. self-review completed. Functionality present. Security evaluated. error handling done. Code Style checked.",
			keywords:    DescriptionKeywords,
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name:        "Review comment keywords present",
			text:        "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Code Style is excellent.",
			keywords:    ReviewCommentKeywords,
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name:        "Review comment missing keyword",
			text:        "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Performance is excellent.",
			keywords:    ReviewCommentKeywords,
			wantValid:   false,
			wantMissing: 1,
		},
		{
			name:        "Empty text",
			text:        "",
			keywords:    ReviewCommentKeywords,
			wantValid:   false,
			wantMissing: 4,
		},
		{
			name:        "Keywords with special characters",
			text:        "This addresses: Description, Changes Made, Self-Review, Functionality, Security, Error Handling, Code Style.",
			keywords:    DescriptionKeywords,
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name:        "Partial keyword match (should fail)",
			text:        "Describe the change. Change made. Self. Function. Secure. Error. Code.",
			keywords:    DescriptionKeywords,
			wantValid:   false,
			wantMissing: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, missing := analyzer.CheckKeywordsInText(tt.text, tt.keywords)

			if valid != tt.wantValid {
				t.Errorf("CheckKeywordsInText() valid = %v, want %v", valid, tt.wantValid)
			}

			if len(missing) != tt.wantMissing {
				t.Errorf("CheckKeywordsInText() missing count = %d, want %d", len(missing), tt.wantMissing)
				t.Logf("Missing keywords: %v", missing)
			}
		})
	}
}

func TestCheckReviewComments(t *testing.T) {
	analyzer := NewPRRuleAnalyzer()

	tests := []struct {
		name        string
		reviews     []*github.ReviewData
		wantValid   bool
		wantMissing int
	}{
		{
			name: "All review comment keywords present",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Code Style is excellent.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name: "Multiple reviewers with combined coverage",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is working well. Security is good.",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "COMMENTED",
					CommentBody:   "Error Handling looks solid. Code Style needs improvement.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name: "Missing one keyword in single review",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is good. Security is solid. Error Handling is implemented. Performance looks fine.",
				},
			},
			wantValid:   false,
			wantMissing: 1,
		},
		{
			name: "Multiple reviews with missing keywords",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is good.",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "COMMENTED",
					CommentBody:   "Security is fine.",
				},
			},
			wantValid:   false,
			wantMissing: 2,
		},
		{
			name:        "No reviews",
			reviews:     []*github.ReviewData{},
			wantValid:   false,
			wantMissing: 4,
		},
		{
			name: "Review with empty comment",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "",
				},
			},
			wantValid:   false,
			wantMissing: 4,
		},
		{
			name: "Multiple reviews with empty comments",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "COMMENTED",
					CommentBody:   "",
				},
			},
			wantValid:   false,
			wantMissing: 4,
		},
		{
			name: "Case insensitive review comments",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "FUNCTIONALITY is good. SECURITY is solid. ERROR HANDLING is complete. CODE STYLE is correct.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name: "Mixed case review keywords",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "functionality is reviewed. Security checked. error handling implemented. code style verified.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name: "Review with all keywords spread across reviews",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality works.",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "COMMENTED",
					CommentBody:   "Security looks good.",
				},
				{
					ReviewerLogin: "reviewer3",
					State:         "APPROVED",
					CommentBody:   "Error Handling is comprehensive.",
				},
				{
					ReviewerLogin: "reviewer4",
					State:         "COMMENTED",
					CommentBody:   "Code Style follows conventions.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
		{
			name: "Single review with all keywords",
			reviews: []*github.ReviewData{
				{
					ReviewerLogin: "senior_reviewer",
					State:         "APPROVED",
					CommentBody:   "Excellent work! Functionality is robust, Security is well-handled, Error Handling covers edge cases, and Code Style is consistent with project standards.",
				},
			},
			wantValid:   true,
			wantMissing: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, missing := analyzer.CheckReviewComments(tt.reviews)

			if valid != tt.wantValid {
				t.Errorf("CheckReviewComments() valid = %v, want %v", valid, tt.wantValid)
			}

			if len(missing) != tt.wantMissing {
				t.Errorf("CheckReviewComments() missing count = %d, want %d", len(missing), tt.wantMissing)
				t.Logf("Missing keywords: %v", missing)
			}
		})
	}
}

func TestAnalyzePRRule(t *testing.T) {
	analyzer := NewPRRuleAnalyzer()

	tests := []struct {
		name                   string
		pr                     *github.PullRequestData
		wantDescriptionValid   bool
		wantReviewCommentValid bool
		wantCompliant          bool
	}{
		{
			name: "Fully compliant PR",
			pr: &github.PullRequestData{
				Number:      1,
				Title:       "Add feature",
				Description: "Description: Feature overview. Changes Made: Added new component. Self-Review: Tested manually. Functionality: Works as expected. Security: No issues. Error Handling: Try-catch implemented. Code Style: Follows conventions.",
				Reviews: []*github.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "Functionality looks great. Security is solid. Error Handling is comprehensive. Code Style is excellent.",
					},
				},
				HTMLURL: "https://github.com/test/pull/1",
			},
			wantDescriptionValid:   true,
			wantReviewCommentValid: true,
			wantCompliant:          true,
		},
		{
			name: "Missing description keywords",
			pr: &github.PullRequestData{
				Number:      2,
				Title:       "Fix bug",
				Description: "This is a simple PR without proper keywords",
				Reviews: []*github.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "Functionality works. Security is good. Error Handling is fine. Code Style is okay.",
					},
				},
				HTMLURL: "https://github.com/test/pull/2",
			},
			wantDescriptionValid:   false,
			wantReviewCommentValid: true,
			wantCompliant:          false,
		},
		{
			name: "Missing review comment keywords",
			pr: &github.PullRequestData{
				Number:      3,
				Title:       "Refactor code",
				Description: "Description: Refactoring. Changes Made: Updated logic. Self-Review: Checked. Functionality: Tested. Security: Verified. Error Handling: Handled. Code Style: Formatted.",
				Reviews: []*github.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "LGTM!",
					},
				},
				HTMLURL: "https://github.com/test/pull/3",
			},
			wantDescriptionValid:   true,
			wantReviewCommentValid: false,
			wantCompliant:          false,
		},
		{
			name: "Missing both description and review keywords",
			pr: &github.PullRequestData{
				Number:      4,
				Title:       "Update docs",
				Description: "Just updated the documentation file",
				Reviews: []*github.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "Looks good",
					},
				},
				HTMLURL: "https://github.com/test/pull/4",
			},
			wantDescriptionValid:   false,
			wantReviewCommentValid: false,
			wantCompliant:          false,
		},
		{
			name: "No reviews",
			pr: &github.PullRequestData{
				Number:      5,
				Title:       "Add feature",
				Description: "Description: New feature. Changes Made: Files added. Self-Review: Done. Functionality: Tested. Security: Checked. Error Handling: Implemented. Code Style: Verified.",
				Reviews:     []*github.ReviewData{},
				HTMLURL:     "https://github.com/test/pull/5",
			},
			wantDescriptionValid:   true,
			wantReviewCommentValid: false,
			wantCompliant:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePRRule(tt.pr)

			if result.PRDescriptionValid != tt.wantDescriptionValid {
				t.Errorf("AnalyzePRRule() PRDescriptionValid = %v, want %v", result.PRDescriptionValid, tt.wantDescriptionValid)
			}

			if result.ReviewCommentValid != tt.wantReviewCommentValid {
				t.Errorf("AnalyzePRRule() ReviewCommentValid = %v, want %v", result.ReviewCommentValid, tt.wantReviewCommentValid)
			}

			if result.PRCompliant != tt.wantCompliant {
				t.Errorf("AnalyzePRRule() PRCompliant = %v, want %v", result.PRCompliant, tt.wantCompliant)
			}

			if len(result.MissingDescKeywords) > 0 && tt.wantDescriptionValid {
				t.Errorf("AnalyzePRRule() should not have missing keywords when valid, got: %v", result.MissingDescKeywords)
			}

			if len(result.MissingReviewKeywords) > 0 && tt.wantReviewCommentValid {
				t.Errorf("AnalyzePRRule() should not have missing review keywords when valid, got: %v", result.MissingReviewKeywords)
			}
		})
	}
}
