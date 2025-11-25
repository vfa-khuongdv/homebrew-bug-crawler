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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr)

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

func TestAnalyzePR_TypeBug(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *github.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
	}{
		{
			name: "type: bug in description",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "type: bug - fix login issue",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "type: bug",
		},
		{
			name: "type:bug without space",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "type:bug",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "type: bug",
		},
		{
			name: "TYPE: BUG case insensitive",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "TYPE: BUG in the system",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "type: bug",
		},
		{
			name: "type: bug in title should be ignored",
			pr: &github.PullRequestData{
				Title:       "type: bug - fix",
				Description: "Some description",
				Labels:      []string{},
			},
			wantBug:  false,
			wantType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr)

			if result.IsBugRelated != tt.wantBug {
				t.Errorf("IsBugRelated = %v, want %v", result.IsBugRelated, tt.wantBug)
			}
			if result.DetectionType != tt.wantType {
				t.Errorf("DetectionType = %v, want %v", result.DetectionType, tt.wantType)
			}
			if result.MatchedKeyword != tt.wantKeyword {
				t.Errorf("MatchedKeyword = %v, want %v", result.MatchedKeyword, tt.wantKeyword)
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
			name: "no matching label",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Some description",
				Labels:      []string{"enhancement", "documentation"},
			},
			wantBug:  false,
			wantType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr)

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

func TestAnalyzePR_GeneralKeywords(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *github.PullRequestData
		wantBug     bool
		wantType    string
		wantKeyword string
	}{
		{
			name: "fix keyword in description",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "This PR fixes the login issue",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "fix",
		},
		{
			name: "bug keyword in description",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Found a bug in the system",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "bug",
		},
		{
			name: "crash keyword",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Prevent crash on startup",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "crash",
		},
		{
			name: "error keyword",
			pr: &github.PullRequestData{
				Title:       "Some PR",
				Description: "Handle error properly",
				Labels:      []string{},
			},
			wantBug:     true,
			wantType:    "keyword",
			wantKeyword: "error",
		},
		{
			name: "keyword in title should be ignored",
			pr: &github.PullRequestData{
				Title:       "Fix bug in login",
				Description: "Some description",
				Labels:      []string{},
			},
			wantBug:  false,
			wantType: "",
		},
		{
			name: "no bug keywords",
			pr: &github.PullRequestData{
				Title:       "Add new feature",
				Description: "Implement user profile page",
				Labels:      []string{},
			},
			wantBug:  false,
			wantType: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr)

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

func TestAnalyzePR_Priority(t *testing.T) {
	analyzer := NewBugAnalyzer()

	tests := []struct {
		name        string
		pr          *github.PullRequestData
		wantType    string
		wantKeyword string
		wantCount   int
	}{
		{
			name: "bug_review takes priority over type: bug",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "bug_review: 2 and type: bug",
				Labels:      []string{},
			},
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   2,
		},
		{
			name: "bug_review takes priority over labels",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "bug_review: 3",
				Labels:      []string{"bug", "critical"},
			},
			wantType:    "bug_review",
			wantKeyword: "bug_review",
			wantCount:   3,
		},
		{
			name: "type: bug takes priority over labels",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "type: bug in the code",
				Labels:      []string{"bug"},
			},
			wantType:    "keyword",
			wantKeyword: "type: bug",
			wantCount:   0,
		},
		{
			name: "type: bug takes priority over general keywords",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "type: bug - fix crash issue",
				Labels:      []string{},
			},
			wantType:    "keyword",
			wantKeyword: "type: bug",
			wantCount:   0,
		},
		{
			name: "labels take priority over general keywords",
			pr: &github.PullRequestData{
				Title:       "Fix",
				Description: "fix the issue",
				Labels:      []string{"hotfix"},
			},
			wantType:    "label",
			wantKeyword: "hotfix",
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.AnalyzePR(tt.pr)

			if !result.IsBugRelated {
				t.Error("Expected IsBugRelated to be true")
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

func TestAnalyzePRs(t *testing.T) {
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

	results := analyzer.AnalyzePRs(prs)

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	if !results[0].IsBugRelated {
		t.Error("First PR should be bug-related")
	}
	if results[1].IsBugRelated {
		t.Error("Second PR should not be bug-related")
	}
	if !results[2].IsBugRelated {
		t.Error("Third PR should be bug-related")
	}
}

func TestGetBugCount(t *testing.T) {
	analyzer := NewBugAnalyzer()

	results := []*BugResult{
		{IsBugRelated: true},
		{IsBugRelated: false},
		{IsBugRelated: true},
		{IsBugRelated: true},
		{IsBugRelated: false},
	}

	count := analyzer.GetBugCount(results)
	if count != 3 {
		t.Errorf("Expected bug count 3, got %d", count)
	}
}
