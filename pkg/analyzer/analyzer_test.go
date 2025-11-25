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
