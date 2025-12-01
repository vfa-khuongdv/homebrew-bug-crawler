package analyzer

import (
	"testing"

	"github.com/bug-crawler/pkg/platform"
)

func TestAnalyzePR_BugReview(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *platform.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
		wantCount   int
	}{
		{
			name: "bug_review with count",
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			result := analyzer.AnalyzePR(tt.pr, "bug_review", "github")

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
		pr          *platform.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
	}{
		{
			name: "bug label",
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
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
			pr: &platform.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"enhancement", "documentation"},
			},
			wantBug:  false,
			wantType: "",
		},
		{
			name: "empty labels",
			pr: &platform.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{},
			},
			wantBug:  false,
			wantType: "",
		},
		{
			name: "nil labels",
			pr: &platform.PullRequestData{
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
			result := analyzer.AnalyzePR(tt.pr, "label", "github") // Use "label" or default

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

	pr := &platform.PullRequestData{
		Title:       "Fix",
		Description: "bug_review: 3",
		Labels:      []string{"bug", "critical"},
	}

	t.Run("bug_review mode ignores labels", func(t *testing.T) {
		result := analyzer.AnalyzePR(pr, "bug_review", "github")
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
		result := analyzer.AnalyzePR(pr, "label", "github")
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

func TestAnalyzePR_BitbucketBacklog(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *platform.PullRequestData
		platform    string
		wantBug     bool
		wantType    string
		wantKeyword string
	}{
		{
			name: "Bitbucket: type: bug in description",
			pr: &platform.PullRequestData{
				Title:       "Fix something",
				Description: "This is a fix. type: bug",
				Labels:      []string{},
			},
			platform:    "bitbucket",
			wantBug:     true,
			wantType:    "description_regex",
			wantKeyword: "type: bug",
		},
		{
			name: "Backlog: type: bug in description",
			pr: &platform.PullRequestData{
				Title:       "Fix something",
				Description: "This is a fix. type: bug",
				Labels:      []string{},
			},
			platform:    "backlog",
			wantBug:     true,
			wantType:    "description_regex",
			wantKeyword: "type: bug",
		},
		{
			name: "Bitbucket: type:bug (no space)",
			pr: &platform.PullRequestData{
				Title:       "Fix something",
				Description: "type:bug",
				Labels:      []string{},
			},
			platform:    "bitbucket",
			wantBug:     true,
			wantType:    "description_regex",
			wantKeyword: "type: bug",
		},
		{
			name: "Bitbucket: TYPE: BUG (case insensitive)",
			pr: &platform.PullRequestData{
				Title:       "Fix something",
				Description: "TYPE: BUG",
				Labels:      []string{},
			},
			platform:    "bitbucket",
			wantBug:     true,
			wantType:    "description_regex",
			wantKeyword: "type: bug",
		},
		{
			name: "GitHub: type: bug in description (should be ignored)",
			pr: &platform.PullRequestData{
				Title:       "Fix something",
				Description: "type: bug",
				Labels:      []string{},
			},
			platform:    "github",
			wantBug:     false,
			wantType:    "",
			wantKeyword: "",
		},
		{
			name: "Bitbucket: no type: bug",
			pr: &platform.PullRequestData{
				Title:       "Feature",
				Description: "New feature",
				Labels:      []string{},
			},
			platform:    "bitbucket",
			wantBug:     false,
			wantType:    "",
			wantKeyword: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr, "bug", tt.platform)

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

func TestAnalyzePR_NotBugRelated(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name string
		pr   *platform.PullRequestData
	}{
		{
			name: "feature PR",
			pr: &platform.PullRequestData{
				Title:       "Add new feature",
				Description: "Implement user profile page",
				Labels:      []string{"enhancement"},
			},
		},
		{
			name: "documentation PR",
			pr: &platform.PullRequestData{
				Title:       "Update README",
				Description: "Add installation instructions",
				Labels:      []string{"documentation"},
			},
		},
		{
			name: "refactor PR",
			pr: &platform.PullRequestData{
				Title:       "Refactor code",
				Description: "Improve code structure",
				Labels:      []string{"refactor"},
			},
		},
		{
			name: "PR with bug keyword in description but no label (in label mode)",
			pr: &platform.PullRequestData{
				Title:       "Some PR",
				Description: "This fixes a bug in the system",
				Labels:      []string{},
			},
		},
		{
			name: "PR with fix keyword in title only",
			pr: &platform.PullRequestData{
				Title:       "Fix typo in documentation",
				Description: "Some description",
				Labels:      []string{},
			},
		},
		{
			name: "empty PR",
			pr: &platform.PullRequestData{
				Title:       "",
				Description: "",
				Labels:      []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr, "label", "github")

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

	prs := []*platform.PullRequestData{
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

	results := analyzer.AnalyzePRs(prs, "bug_review", "github")

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

	prs := []*platform.PullRequestData{
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

	results := analyzer.AnalyzePRs(prs, "label", "github")

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

	results := analyzer.AnalyzePRs([]*platform.PullRequestData{}, "label", "github")

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestAnalyzePRs_NilList(t *testing.T) {
	analyzer := NewBugAnalyzer()

	results := analyzer.AnalyzePRs(nil, "label", "github")

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
		name      string
		text      string
		keywords  []string
		wantValid bool
	}{
		{
			name:      "All keywords present (7/7)",
			text:      "Description of the work. Changes Made to the code. Self-Review checklist completed. Functionality is working. Security is handled. Error Handling is in place. Code Style follows conventions. Dependencies are managed.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "6 keywords present (>2 required)",
			text:      "Description of the work. Changes Made to the code. Self-Review checklist. Functionality works. Security is good. Error Handling included. Dependencies managed.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Exactly 3 keywords (minimum required, >2)",
			text:      "Description provided. Changes Made documented. Functionality tested.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Only 2 keywords (should fail, need >2)",
			text:      "Description provided. Changes Made documented.",
			keywords:  DescriptionKeywords,
			wantValid: false,
		},
		{
			name:      "Only 1 keyword (should fail)",
			text:      "Description of the work.",
			keywords:  DescriptionKeywords,
			wantValid: false,
		},
		{
			name:      "Case insensitive matching (7/7)",
			text:      "DESCRIPTION of work. CHANGES MADE to files. SELF-REVIEW done. FUNCTIONALITY tested. SECURITY checked. ERROR HANDLING implemented. CODE STYLE followed. DEPENDENCIES updated.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Mixed case keywords (7/7)",
			text:      "description here. Changes Made included. self-review completed. Functionality present. Security evaluated. error handling done. Code Style checked. Dependencies reviewed.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Review comment all 5 keywords (>2 required)",
			text:      "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Code Style is excellent. Code Readability is clear.",
			keywords:  ReviewCommentKeywords,
			wantValid: true,
		},
		{
			name:      "Review comment 4 keywords (>2 required)",
			text:      "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Code Style is excellent.",
			keywords:  ReviewCommentKeywords,
			wantValid: true,
		},
		{
			name:      "Review comment exactly 3 keywords (minimum)",
			text:      "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive.",
			keywords:  ReviewCommentKeywords,
			wantValid: true,
		},
		{
			name:      "Review comment only 2 keywords (should fail)",
			text:      "Reviewed the functionality carefully. Security looks good.",
			keywords:  ReviewCommentKeywords,
			wantValid: false,
		},
		{
			name:      "Review comment only 1 keyword (should fail)",
			text:      "Reviewed the functionality carefully.",
			keywords:  ReviewCommentKeywords,
			wantValid: false,
		},
		{
			name:      "Empty text (0 keywords)",
			text:      "",
			keywords:  ReviewCommentKeywords,
			wantValid: false,
		},
		{
			name:      "Keywords with special characters (7/7)",
			text:      "This addresses: Description, Changes Made, Self-Review, Functionality, Security, Error Handling, Code Style, Dependencies.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Abbreviated tags (D1, CM1, SR1, F1, S1, EH1, C1)",
			text:      "D1: Feature overview. CM1: Updated code. SR1: Self checked. F1: Works well. S1: Secure. EH1: Error handled. C1: Code formatted. DEP1: Deps updated.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Mixed full and abbreviated keywords (7/7)",
			text:      "Description provided. CM1: Updated. Self-Review done. F1: Working. Security verified. Error Handling solid. Code Style good. Dependencies managed.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Partial keyword match only (should fail - need real keywords)",
			text:      "Describe the change. Change made. Self. Function. Secure. Error. Code.",
			keywords:  DescriptionKeywords,
			wantValid: false,
		},
		{
			name:      "Real-world PR description with 3+ keywords",
			text:      "Description: Fixed login bug. Changes Made: Updated auth logic. Self-Review: Tested locally. Functionality: Works correctly.",
			keywords:  DescriptionKeywords,
			wantValid: true,
		},
		{
			name:      "Real-world minimal PR (only 2 keywords)",
			text:      "Description: Update version. Changes Made: Bumped to 1.0.7.",
			keywords:  DescriptionKeywords,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := analyzer.CheckKeywordsInText(tt.text, tt.keywords)

			if valid != tt.wantValid {
				t.Errorf("CheckKeywordsInText() valid = %v, want %v", valid, tt.wantValid)
				t.Logf("Text: %q", tt.text)
			}
		})
	}
}

func TestCheckReviewComments(t *testing.T) {
	analyzer := NewPRRuleAnalyzer()

	tests := []struct {
		name      string
		reviews   []*platform.ReviewData
		wantValid bool
	}{
		{
			name: "All 5 keywords present",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Reviewed the functionality carefully. Security looks good. Error Handling is comprehensive. Code Style is excellent. Readability is clear.",
				},
			},
			wantValid: true,
		},
		{
			name: "Exactly 3 keywords (minimum required)",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality works well. Security is good. Error Handling is solid.",
				},
			},
			wantValid: true,
		},
		{
			name: "4 out of 5 keywords (more than minimum)",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is working well. Security is good. Error Handling looks solid. Code Style needs improvement.",
				},
			},
			wantValid: true,
		},
		{
			name: "Multiple reviewers with combined 3+ keywords",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is working well.",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "COMMENTED",
					CommentBody:   "Security is good. Error Handling is solid.",
				},
			},
			wantValid: true,
		},
		{
			name: "Only 2 keywords (below minimum)",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is good. Security is solid. Performance looks fine.",
				},
			},
			wantValid: false,
		},
		{
			name: "Only 1 keyword (well below minimum)",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality is good. Performance and efficiency are great.",
				},
			},
			wantValid: false,
		},
		{
			name: "No keywords",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Looks good to me. Great work!",
				},
			},
			wantValid: false,
		},
		{
			name:      "No reviews",
			reviews:   []*platform.ReviewData{},
			wantValid: false,
		},
		{
			name: "Review with empty comment",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "",
				},
			},
			wantValid: false,
		},
		{
			name: "Case insensitive - 3 keywords",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "FUNCTIONALITY is good. SECURITY is solid. ERROR HANDLING is complete.",
				},
			},
			wantValid: true,
		},
		{
			name: "Mixed case - all 5 keywords",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "functionality is reviewed. Security checked. error handling implemented. code style verified. Readability is excellent.",
				},
			},
			wantValid: true,
		},
		{
			name: "Real example from WawaTalk PR #80",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "vfa-khuongdv",
					State:         "COMMENTED",
					CommentBody:   "## Review Summary\nAll checklist items **PASS** ✅\n### Checklist Results\n- **F1: Functionality** ✅ – The code correctly implements layout improvements\n- **S1: Security** ✅ – No user input handling changes\n- **EH1: Error Handling** ✅ – No critical external calls added\n- **C1: Code Style/Readability** ✅ – Code follows clean formatting",
				},
			},
			wantValid: true,
		},
		{
			name: "Multiple reviews - combined to reach minimum",
			reviews: []*platform.ReviewData{
				{
					ReviewerLogin: "reviewer1",
					State:         "APPROVED",
					CommentBody:   "Functionality works.",
				},
				{
					ReviewerLogin: "reviewer2",
					State:         "APPROVED",
					CommentBody:   "Security looks good.",
				},
				{
					ReviewerLogin: "reviewer3",
					State:         "APPROVED",
					CommentBody:   "Error Handling is solid.",
				},
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := analyzer.CheckReviewComments(tt.reviews)

			if valid != tt.wantValid {
				t.Errorf("CheckReviewComments() valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestAnalyzePRRule(t *testing.T) {
	analyzer := NewPRRuleAnalyzer()

	tests := []struct {
		name                   string
		pr                     *platform.PullRequestData
		wantDescriptionValid   bool
		wantReviewCommentValid bool
		wantCompliant          bool
	}{
		{
			name: "Fully compliant PR",
			pr: &platform.PullRequestData{
				Number:      1,
				Title:       "Add feature",
				Description: "Description: Feature overview. Changes Made: Added new component. Self-Review: Tested manually. Functionality: Works as expected. Security: No issues. Error Handling: Try-catch implemented. Code Style: Follows conventions. Dependencies: Updated.",
				Reviews: []*platform.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "Functionality looks great. Security is solid. Error Handling is comprehensive. Code Style is excellent. Code Readability is clear.",
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
			pr: &platform.PullRequestData{
				Number:      2,
				Title:       "Fix bug",
				Description: "This is a simple PR without proper keywords",
				Reviews: []*platform.ReviewData{
					{
						ReviewerLogin: "reviewer1",
						State:         "APPROVED",
						CommentBody:   "Functionality works. Security is good. Error Handling is fine. Code Style is okay. Code Readability is excellent.",
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
			pr: &platform.PullRequestData{
				Number:      3,
				Title:       "Refactor code",
				Description: "Description: Refactoring. Changes Made: Updated logic. Self-Review: Checked. Functionality: Tested. Security: Verified. Error Handling: Handled. Code Style: Formatted. Dependencies: Updated.",
				Reviews: []*platform.ReviewData{
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
			pr: &platform.PullRequestData{
				Number:      4,
				Title:       "Update docs",
				Description: "Just updated the documentation file",
				Reviews: []*platform.ReviewData{
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
			pr: &platform.PullRequestData{
				Number:      5,
				Title:       "Add feature",
				Description: "Description: New feature. Changes Made: Files added. Self-Review: Done. Functionality: Tested. Security: Checked. Error Handling: Implemented. Code Style: Verified. Dependencies: Updated.",
				Reviews:     []*platform.ReviewData{},
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

		})
	}
}
