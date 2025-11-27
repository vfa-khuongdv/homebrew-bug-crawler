package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/manifoldco/promptui"
)

// CLI manages user interactions
type CLI struct{}

// NewCLI khởi tạo CLI
func NewCLI() *CLI {
	return &CLI{}
}

// PromptToken yêu cầu user nhập GitHub token
func (c *CLI) PromptToken() (string, error) {
	prompt := promptui.Prompt{
		Label: "GitHub Token",
		Mask:  '*',
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// PromptSaveToken hỏi user có muốn lưu token không
func (c *CLI) PromptSaveToken() (bool, error) {
	prompt := promptui.Select{
		Label: "Lưu token vào file config?",
		Items: []string{"Có", "Không"},
	}

	_, result, err := prompt.Run()
	return result == "Có", err
}

// RepositoryScanMode chế độ quét repositories
type RepositoryScanMode int

const (
	ScanModeManual RepositoryScanMode = iota
	ScanModeUser
	ScanModeOrganization
	ScanModeCurrentUser
)

// PromptRepositoryScanMode yêu cầu user chọn chế độ quét
func (c *CLI) PromptRepositoryScanMode() (RepositoryScanMode, error) {
	prompt := promptui.Select{
		Label: "Chọn cách quét repositories",
		Items: []string{
			"1. Nhập thủ công (owner/repo)",
			"2. Quét repositories của user",
			"3. Quét repositories của organization",
			"4. Quét repositories của tôi",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return ScanModeManual, err
	}

	return RepositoryScanMode(index), nil
}

// PromptRepositories yêu cầu user nhập repositories thủ công
func (c *CLI) PromptRepositories() ([]string, error) {
	fmt.Println("\nNhập repositories (format: owner/repo, mỗi repo trên một dòng, nhấn Enter 2 lần để xong):")

	var repos []string
	for {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Repo %d", len(repos)+1),
		}

		result, err := prompt.Run()
		if err != nil {
			return nil, err
		}

		result = strings.TrimSpace(result)
		if result == "" {
			if len(repos) == 0 {
				fmt.Println("Vui lòng nhập ít nhất 1 repository")
				continue
			}
			break
		}

		if !strings.Contains(result, "/") {
			fmt.Println("Format không hợp lệ. Vui lòng nhập dạng: owner/repo")
			continue
		}

		repos = append(repos, result)
	}

	return repos, nil
}

// PromptUsername yêu cầu user nhập username GitHub
func (c *CLI) PromptUsername() (string, error) {
	prompt := promptui.Prompt{
		Label: "GitHub Username",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// PromptOrganization yêu cầu user nhập tên organization
func (c *CLI) PromptOrganization() (string, error) {
	prompt := promptui.Prompt{
		Label: "Organization Name",
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

// PromptSelectMultipleRepositories cho phép user chọn nhiều repositories từ danh sách với arrow keys và space
func (c *CLI) PromptSelectMultipleRepositories(repos []string) ([]string, error) {
	if len(repos) == 0 {
		return nil, fmt.Errorf("không tìm thấy repositories")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		// Fallback to simple input if tcell fails
		return c.simpleSelectRepositories(repos)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		return c.simpleSelectRepositories(repos)
	}

	screen.SetStyle(tcell.StyleDefault)
	selected := make(map[int]bool)
	cursor := 0
	scrollOffset := 0

	for {
		screen.Clear()

		width, height := screen.Size()
		maxDisplay := height - 4

		// Draw header
		headerText := "Chọn Repositories (↑↓=navigate, Space=select, Enter=confirm)"
		for i, r := range headerText {
			if i < width {
				screen.SetContent(i, 0, r, nil, tcell.StyleDefault.Bold(true))
			}
		}

		// Draw separator
		for i := 0; i < width; i++ {
			screen.SetContent(i, 1, '─', nil, tcell.StyleDefault)
		}

		// Calculate scroll
		if cursor < scrollOffset {
			scrollOffset = cursor
		}
		if cursor >= scrollOffset+maxDisplay {
			scrollOffset = cursor - maxDisplay + 1
		}

		// Draw repositories
		row := 2
		for i := scrollOffset; i < len(repos) && i < scrollOffset+maxDisplay; i++ {
			repo := repos[i]
			checkbox := " "
			if selected[i] {
				checkbox = "✓"
			}

			// Highlight current line
			style := tcell.StyleDefault
			if i == cursor {
				style = style.Reverse(true).Bold(true)
			}

			// Draw line
			line := fmt.Sprintf("[%s] %s", checkbox, repo)
			if len(line) > width {
				line = line[:width-1]
			}

			for j, ch := range line {
				if j < width {
					screen.SetContent(j, row, ch, nil, style)
				}
			}
			row++
		}

		// Draw footer
		footerText := fmt.Sprintf("Selected: %d/%d", len(selected), len(repos))
		for i, r := range footerText {
			if i < width {
				screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Dim(true))
			}
		}

		screen.Show()

		// Handle events
		ev := screen.PollEvent()
		if ev == nil {
			continue
		}

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if cursor > 0 {
					cursor--
				}
			case tcell.KeyDown:
				if cursor < len(repos)-1 {
					cursor++
				}
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ': // Space
					selected[cursor] = !selected[cursor]
				}
			case tcell.KeyEnter:
				screen.Fini()
				var result []string
				for i := 0; i < len(repos); i++ {
					if selected[i] {
						result = append(result, repos[i])
					}
				}
				if len(result) == 0 {
					fmt.Println("❌ Vui lòng chọn ít nhất 1 repository")
					time.Sleep(1 * time.Second)
					return c.simpleSelectRepositories(repos)
				}
				return result, nil
			case tcell.KeyEscape:
				screen.Fini()
				return nil, fmt.Errorf("đã hủy chọn")
			}
		}
	}
}

// simpleSelectRepositories fallback khi tcell không khả dụng
func (c *CLI) simpleSelectRepositories(repos []string) ([]string, error) {
	if len(repos) == 0 {
		return nil, fmt.Errorf("không tìm thấy repositories")
	}

	selected := make(map[string]bool)

	for {
		fmt.Print("\033[2J\033[H")
		fmt.Println("Chọn repositories:")
		fmt.Println()

		for i, repo := range repos {
			checkbox := " "
			if selected[repo] {
				checkbox = "✓"
			}
			fmt.Printf("[%s] %2d. %s\n", checkbox, i+1, repo)
		}

		fmt.Println()

		var input string
		fmt.Print("Nhập số (1,3,5) hoặc 'all', Enter để xác nhận: ")
		_, _ = fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		if input == "" {
			var result []string
			for _, repo := range repos {
				if selected[repo] {
					result = append(result, repo)
				}
			}
			if len(result) == 0 {
				fmt.Println("❌ Vui lòng chọn ít nhất 1 repository")
				time.Sleep(1 * time.Second)
				continue
			}
			return result, nil
		}

		if strings.ToLower(input) == "all" {
			for _, repo := range repos {
				selected[repo] = true
			}
			continue
		}

		indices := strings.Split(input, ",")
		for _, idx := range indices {
			idx = strings.TrimSpace(idx)
			var i int
			_, err := fmt.Sscanf(idx, "%d", &i)
			if err != nil || i < 1 || i > len(repos) {
				continue
			}
			repo := repos[i-1]
			selected[repo] = !selected[repo]
		}
	}
}

// PromptSelectRepositories cho phép user chọn repositories từ danh sách
func (c *CLI) PromptSelectRepositories(repos []string) ([]string, error) {
	if len(repos) == 0 {
		return nil, fmt.Errorf("không tìm thấy repositories")
	}

	fmt.Printf("\nTìm được %d repositories. Chọn repositories (nhập index cách nhau bằng dấu phẩy, hoặc 'all' để chọn tất cả):\n", len(repos))

	// In danh sách repositories
	for i, repo := range repos {
		fmt.Printf("%3d. %s\n", i+1, repo)
	}

	prompt := promptui.Prompt{
		Label: "Chọn repositories",
	}

	input, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)

	// Nếu user chọn "all"
	if strings.ToLower(input) == "all" {
		return repos, nil
	}

	// Parse indices
	var selected []string
	indices := strings.Split(input, ",")

	for _, idx := range indices {
		idx = strings.TrimSpace(idx)
		var i int
		_, err := fmt.Sscanf(idx, "%d", &i)
		if err != nil || i < 1 || i > len(repos) {
			fmt.Printf("⚠️  Index không hợp lệ: %s\n", idx)
			continue
		}
		selected = append(selected, repos[i-1])
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("không có repositories nào được chọn")
	}

	return selected, nil
}

// PromptDateRange yêu cầu user chọn khoảng thời gian
func (c *CLI) PromptDateRange() (time.Time, time.Time, error) {
	fmt.Println("\nChọn khoảng thời gian phân tích (format: YYYY-MM-DD)")

	startDateStr, err := c.promptDate("Ngày bắt đầu")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDateStr, err := c.promptDate("Ngày kết thúc")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("định dạng ngày không hợp lệ")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("định dạng ngày không hợp lệ")
	}

	if startDate.After(endDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("ngày bắt đầu không được sau ngày kết thúc")
	}

	return startDate, endDate.AddDate(0, 0, 1), nil // Thêm 1 ngày để include cả ngày kết thúc
}

// promptDate helper function để prompt date
func (c *CLI) promptDate(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label + " (YYYY-MM-DD)",
	}

	return prompt.Run()
}

// PromptSelectScanSource cho phép chọn scan user hoặc organizations
func (c *CLI) PromptSelectScanSource() (string, error) {
	prompt := promptui.Select{
		Label: "Chọn loại để scan",
		Items: []string{
			"1. Repositories của tôi (User)",
			"2. Repositories của Organizations",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if index == 0 {
		return "user", nil
	}
	return "org", nil
}

// PromptSelectScanMode cho phép chọn giữa Bug Scan và PR Rules Scan
func (c *CLI) PromptSelectScanMode() (string, error) {
	prompt := promptui.Select{
		Label: "Chọn chế độ scan",
		Items: []string{
			"1. Bug Detection (Scan bugs)",
			"2. Code Review Compliance (Scan PR rules)",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if index == 0 {
		return "bug", nil
	}
	return "pr_rules", nil
}

// PromptSelectBugType cho phép chọn loại bug để scan
func (c *CLI) PromptSelectBugType() (string, error) {
	prompt := promptui.Select{
		Label: "Chọn loại bug để scan",
		Items: []string{
			"1. Scan bug (từ labels)",
			"2. Scan bug_review",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if index == 0 {
		return "bug", nil
	}
	return "bug_review", nil
}

// PromptSelectOrganizations cho phép chọn nhiều organizations với space key
func (c *CLI) PromptSelectOrganizations(organizations []string) ([]string, error) {
	if len(organizations) == 0 {
		return nil, fmt.Errorf("không tìm thấy organizations")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return c.simpleSelectFromList(organizations, "Organizations")
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		return c.simpleSelectFromList(organizations, "Organizations")
	}

	screen.SetStyle(tcell.StyleDefault)
	selected := make(map[int]bool)
	cursor := 0
	scrollOffset := 0

	for {
		screen.Clear()

		width, height := screen.Size()
		maxDisplay := height - 4

		// Draw header
		headerText := "Chọn Organizations (↑↓=navigate, Space=select, Enter=confirm)"
		for i, r := range headerText {
			if i < width {
				screen.SetContent(i, 0, r, nil, tcell.StyleDefault.Bold(true))
			}
		}

		// Draw separator
		for i := 0; i < width; i++ {
			screen.SetContent(i, 1, '─', nil, tcell.StyleDefault)
		}

		// Calculate scroll
		if cursor < scrollOffset {
			scrollOffset = cursor
		}
		if cursor >= scrollOffset+maxDisplay {
			scrollOffset = cursor - maxDisplay + 1
		}

		// Draw organizations
		row := 2
		for i := scrollOffset; i < len(organizations) && i < scrollOffset+maxDisplay; i++ {
			org := organizations[i]
			checkbox := " "
			if selected[i] {
				checkbox = "✓"
			}

			style := tcell.StyleDefault
			if i == cursor {
				style = style.Reverse(true).Bold(true)
			}

			line := fmt.Sprintf("[%s] %s", checkbox, org)
			if len(line) > width {
				line = line[:width-1]
			}

			for j, ch := range line {
				if j < width {
					screen.SetContent(j, row, ch, nil, style)
				}
			}
			row++
		}

		// Draw footer
		footerText := fmt.Sprintf("Selected: %d/%d", len(selected), len(organizations))
		for i, r := range footerText {
			if i < width {
				screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Dim(true))
			}
		}

		screen.Show()

		// Handle events
		ev := screen.PollEvent()
		if ev == nil {
			continue
		}

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyUp:
				if cursor > 0 {
					cursor--
				}
			case tcell.KeyDown:
				if cursor < len(organizations)-1 {
					cursor++
				}
			case tcell.KeyRune:
				switch ev.Rune() {
				case ' ': // Space
					selected[cursor] = !selected[cursor]
				}
			case tcell.KeyEnter:
				screen.Fini()
				var result []string
				for i := 0; i < len(organizations); i++ {
					if selected[i] {
						result = append(result, organizations[i])
					}
				}
				if len(result) == 0 {
					fmt.Println("❌ Vui lòng chọn ít nhất 1 organization")
					time.Sleep(1 * time.Second)
					return c.PromptSelectOrganizations(organizations) // Retry
				}
				return result, nil
			case tcell.KeyEscape:
				screen.Fini()
				return nil, fmt.Errorf("đã hủy chọn")
			}
		}
	}
}

// simpleSelectFromList fallback khi tcell không khả dụng
func (c *CLI) simpleSelectFromList(items []string, itemType string) ([]string, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("không tìm thấy %s", itemType)
	}

	selected := make(map[string]bool)

	for {
		fmt.Print("\033[2J\033[H")
		fmt.Printf("Chọn %s:\n", itemType)
		fmt.Println()

		for i, item := range items {
			checkbox := " "
			if selected[item] {
				checkbox = "✓"
			}
			fmt.Printf("[%s] %2d. %s\n", checkbox, i+1, item)
		}

		fmt.Println()

		var input string
		fmt.Print("Nhập số (1,3,5) hoặc 'all', Enter để xác nhận: ")
		_, _ = fmt.Scanln(&input)

		input = strings.TrimSpace(input)

		if input == "" {
			var result []string
			for _, item := range items {
				if selected[item] {
					result = append(result, item)
				}
			}
			if len(result) == 0 {
				fmt.Println("❌ Vui lòng chọn ít nhất 1 item")
				time.Sleep(1 * time.Second)
				continue
			}
			return result, nil
		}

		if strings.ToLower(input) == "all" {
			for _, item := range items {
				selected[item] = true
			}
			continue
		}

		indices := strings.Split(input, ",")
		for _, idx := range indices {
			idx = strings.TrimSpace(idx)
			var i int
			_, err := fmt.Sscanf(idx, "%d", &i)
			if err != nil || i < 1 || i > len(items) {
				continue
			}
			item := items[i-1]
			selected[item] = !selected[item]
		}
	}
}
