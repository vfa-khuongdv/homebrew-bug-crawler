package report

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/bug-crawler/pkg/analyzer"
)

// Statistics chứa thống kê bug
type Statistics struct {
	TotalPRs        int
	BugRelatedPRs   int
	ByKeyword       int
	ByLabel         int
	ByBugReview     int
	TotalBugCount   int // Tổng số bugs từ bug_review tags
	BugPercentage   float64
	DetailedResults []*analyzer.BugResult
}

// Reporter tạo báo cáo
type Reporter struct{}

// NewReporter khởi tạo Reporter
func NewReporter() *Reporter {
	return &Reporter{}
}

// GenerateStatistics tạo thống kê từ kết quả phân tích
func (r *Reporter) GenerateStatistics(results []*analyzer.BugResult) *Statistics {
	stats := &Statistics{
		TotalPRs:        len(results),
		DetailedResults: results,
	}

	byLabel := 0
	byBugReview := 0
	totalBugCount := 0
	bugCount := 0

	for _, result := range results {
		if result.IsBugRelated {
			bugCount++
			switch result.DetectionType {
			case "label":
				byLabel++
			case "bug_review":
				byBugReview++
				totalBugCount += result.BugCount
			}
		}
	}

	stats.BugRelatedPRs = bugCount
	stats.ByLabel = byLabel
	stats.ByBugReview = byBugReview
	stats.TotalBugCount = totalBugCount

	if stats.TotalPRs > 0 {
		stats.BugPercentage = float64(bugCount) * 100 / float64(stats.TotalPRs)
	}

	return stats
}

// PrintSummary in tóm tắt thống kê
func (r *Reporter) PrintSummary(stats *Statistics) {
	separator := "============================================================"
	fmt.Println("\n" + separator)
	fmt.Println("THỐNG KÊ BUG")
	fmt.Println(separator)
	fmt.Printf("Tổng số PR: %d\n", stats.TotalPRs)
	fmt.Printf("PR liên quan bug: %d\n", stats.BugRelatedPRs)
	if stats.ByBugReview > 0 {
		fmt.Printf("  ├─ Phát hiện qua bug_review tag: %d (Tổng bugs: %d)\n", stats.ByBugReview, stats.TotalBugCount)
	}
	if stats.ByLabel > 0 {
		fmt.Printf("  └─ Phát hiện qua label: %d\n", stats.ByLabel)
	}
	fmt.Printf("Tỷ lệ bug: %.2f%%\n", stats.BugPercentage)
	fmt.Println(separator)
}

// PrintDetails in chi tiết từng PR
func (r *Reporter) PrintDetails(stats *Statistics) {
	if len(stats.DetailedResults) == 0 {
		return
	}

	separator := "=========================================================================================================================="
	fmt.Println("\nCHI TIẾT CÁC PR LIÊN QUAN BUG:")
	fmt.Println(separator)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PR#\tTITLE\tAUTHOR\tPHẦT HIỆN\tBUGS/KEYWORD/LABEL")

	for _, result := range stats.DetailedResults {
		if result.IsBugRelated {
			title := result.PR.Title
			if len(title) > 40 {
				title = title[:37] + "..."
			}

			detailInfo := ""
			if result.DetectionType == "bug_review" {
				detailInfo = fmt.Sprintf("%d bugs", result.BugCount)
			} else {
				detailInfo = result.MatchedKeyword
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				result.PR.Number,
				title,
				result.PR.Author,
				result.DetectionType,
				detailInfo)
		}
	}

	w.Flush()
	fmt.Println(separator)
}

// ExportJSON export kết quả dưới dạng JSON (có thể mở rộng sau)
func (r *Reporter) ExportCSV(filename string, stats *Statistics) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "PR#,Title,Author,Detection Type,Matched Keyword,Number Bug,URL")

	for _, result := range stats.DetailedResults {
		if result.IsBugRelated {
			numberBug := 1
			if result.DetectionType == "bug_review" {
				numberBug = result.BugCount
			}

			fmt.Fprintf(file, "%d,\"%s\",%s,%s,%s,%d,%s\n",
				result.PR.Number,
				result.PR.Title,
				result.PR.Author,
				result.DetectionType,
				result.MatchedKeyword,
				numberBug,
				result.PR.HTMLURL,
			)
		}
	}

	fmt.Printf("\nKết quả đã được export vào: %s\n", filename)
	return nil
}
