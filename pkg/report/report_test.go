package report

import (
	"os"
	"strings"
	"testing"

	"github.com/bug-crawler/pkg/analyzer"
	"github.com/bug-crawler/pkg/platform"
)

func TestExportCSV(t *testing.T) {
	// Setup test data
	results := []*analyzer.BugResult{
		{
			PR: &platform.PullRequestData{
				Number:  1,
				Title:   "Fix bug 1",
				Author:  "user1",
				HTMLURL: "http://github.com/org/repo/pull/1",
			},
			IsBugRelated:   true,
			DetectionType:  "label",
			MatchedKeyword: "bug",
			BugCount:       0,
		},
		{
			PR: &platform.PullRequestData{
				Number:  2,
				Title:   "Review bug fix",
				Author:  "user2",
				HTMLURL: "http://github.com/org/repo/pull/2",
			},
			IsBugRelated:   true,
			DetectionType:  "bug_review",
			MatchedKeyword: "bug_review",
			BugCount:       3,
		},
		{
			PR: &platform.PullRequestData{
				Number:  3,
				Title:   "Feature 1",
				Author:  "user3",
				HTMLURL: "http://github.com/org/repo/pull/3",
			},
			IsBugRelated:  false,
			DetectionType: "",
			BugCount:      0,
		},
	}

	stats := &Statistics{
		DetailedResults: results,
	}

	reporter := NewReporter()
	filename := "test_report.csv"
	defer func() { _ = os.Remove(filename) }()

	// Execute
	err := reporter.ExportCSV(filename, stats)
	if err != nil {
		t.Fatalf("ExportCSV failed: %v", err)
	}

	// Verify
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read generated CSV: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 { // Header + 2 bug related PRs
		t.Errorf("Expected 3 lines in CSV, got %d", len(lines))
	}

	// Check header
	expectedHeader := "PR#,Title,Author,Detection Type,Matched Keyword,Number Bug,Date Opened,URL"
	if lines[0] != expectedHeader {
		t.Errorf("Header mismatch.\nExpected: %s\nGot:      %s", expectedHeader, lines[0])
	}

	// Check row 1 (Label detection)
	// PR#,Title,Author,Detection Type,Matched Keyword,Number Bug,URL
	// 1,"Fix bug 1",user1,label,bug,1,http://github.com/org/repo/pull/1
	if !strings.Contains(lines[1], ",1,") {
		t.Errorf("Row 1 should contain number_bug=1. Got: %s", lines[1])
	}

	// Check row 2 (Bug Review detection)
	// 2,"Review bug fix",user2,bug_review,bug_review,3,http://github.com/org/repo/pull/2
	if !strings.Contains(lines[2], ",3,") {
		t.Errorf("Row 2 should contain number_bug=3. Got: %s", lines[2])
	}
}
