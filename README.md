# 🐛 Bug Crawler - Multi-Platform PR Bug Analysis Tool

> Công cụ tự động phân tích và thống kê bug từ Pull Request trên GitHub, Bitbucket và Backlog

Ứng dụng terminal Go để phân tích tự động các Pull Request trên các nền tảng Git, phát hiện bug dựa trên keywords và labels, rồi xuất kết quả dạng CSV cho báo cáo.

## ✨ Tính Năng Chính

- 🔐 **Quản lý token an toàn** - Lưu token vào file config được mã hóa
- 📦 **Hỗ trợ đa nền tảng**:
  - GitHub
  - Bitbucket
  - Backlog
- 🎯 **Tự động xử lý** - Sử dụng tất cả repositories tìm được
- 📅 **Lọc theo thời gian** - Phân tích PR trong khoảng thời gian tùy chọn
- 🔍 **2 phương pháp phát hiện bug thông minh**:
  - Label-based: Phát hiện từ PR labels (`bug`, `fix`, `hotfix`, `critical`, `error`, `issue`)
  - Tag-based: Phát hiện từ pattern `bug_review: <number>` trong PR description
- 📊 **Thống kê chi tiết** - Tóm tắt và chi tiết từng PR liên quan bug
- 📁 **Export CSV** - Xuất kết quả dạng CSV cho báo cáo

## 🚀 Giới Thiệu Nhanh

**Bug Crawler** giúp team **giảm 80% thời gian scan bug sau mỗi sprint**. Chỉ cần nhập token, chọn repositories, tool sẽ tự động:
- Phát hiện tất cả PR liên quan bug
- Thống kê chi tiết theo keywords/labels
- Xuất báo cáo CSV sẵn sàng gửi BPM

**Sử dụng:**
```bash
brew install vfa-khuongdv/homebrew-bug-crawler/bug-crawler
bug-crawler
```

## 📥 Cài Đặt

### Yêu Cầu
- **Go 1.23+** (nếu build từ source)
- **Personal Access Token** cho platform tương ứng

### Cách 1: Cài đặt qua Homebrew ⭐ (Khuyên dùng)

```bash
# Cài đặt
brew tap vfa-khuongdv/homebrew-bug-crawler
brew install bug-crawler

# Chạy ứng dụng
bug-crawler
```

```bash
# Gỡ cài đặt
brew untap vfa-khuongdv/homebrew-bug-crawler
brew uninstall bug-crawler
```

```bash
# Cập nhật phiên bản Homebrew Tap
brew update
brew upgrade bug-crawler
```

### Cách 2: Build từ Source

```bash
git clone https://github.com/vfa-khuongdv/homebrew-bug-crawler.git
cd homebrew-bug-crawler

# Tải dependencies
go mod tidy

# Build
go build -o bug-crawler ./cmd/main.go

# Chạy
./bug-crawler
```

### Update Package

```bash
brew upgrade bug-crawler
```

## 📖 Sử Dụng

### Chạy Ứng Dụng

**Cài qua Homebrew:**
```bash
bug-crawler
```

**Build từ source:**
```bash
./bug-crawler
```

### 🔄 Luồng Sử Dụng Chi Tiết (7 Bước)

#### **Bước 1: Chọn Platform**
- Chọn platform bạn muốn scan: GitHub, Bitbucket, hoặc Backlog

#### **Bước 2: Xác Thực**
- Nhập token/API key tương ứng
- Ứng dụng tự động xác thực với API của platform
- Hiển thị thông tin tài khoản đã đăng nhập

#### **Bước 3: Chọn Scan Source**
Bạn có 2 lựa chọn:

**Option 1: Repositories của bạn (User)**
- Tự động quét tất cả repositories của tài khoản GitHub của bạn
- Nhanh, phù hợp phân tích toàn bộ projects cá nhân

**Option 2: Repositories của Organizations**
- Hiển thị danh sách organizations bạn là thành viên
- Chọn một hoặc nhiều organizations
- Ứng dụng sẽ quét tất cả repositories từ organizations đã chọn

#### **Bước 4: Chọn Repositories**
- Ứng dụng hiển thị danh sách repositories từ scan source
- Bạn có thể:
  - Chọn từng repository bằng arrow keys + Space
  - Nhập `all` để chọn tất cả
  - Hoặc nhập số index cách nhau bằng dấu phẩy (ví dụ: `1,3,5`)

#### **Bước 5: Chọn Khoảng Thời Gian**
- Nhập ngày bắt đầu: `YYYY-MM-DD` (ví dụ: `2024-01-01`)
- Nhập ngày kết thúc: `YYYY-MM-DD` (ví dụ: `2024-12-31`)
- Ứng dụng chỉ phân tích PR tạo trong khoảng thời gian này

#### **Bước 6: Chọn Loại Bug**
Bạn có 2 lựa chọn:

**Option 1: Scan bug (từ labels)**
- Phát hiện PR có labels liên quan bug
- Labels được tìm kiếm: `bug`, `fix`, `hotfix`, `critical`, `error`, `issue`

**Option 2: Scan bug_review**
- Phát hiện PR có pattern `bug_review: <number>` trong description
- Extract số lượng bugs từ tag này

#### **Bước 7: Crawler, Phân Tích & Báo Cáo**
- Ứng dụng lấy tất cả PR từ repositories được chọn
- Phân tích từng PR dựa trên loại bug đã chọn
- In tóm tắt thống kê
- In chi tiết PR liên quan bug
- Export kết quả vào `bug_report.csv`

## 📁 Cấu Trúc Dự Án

```
homebrew-bug-crawler/
├── cmd/
│   └── main.go                      # Entry point chính
├── pkg/
│   ├── auth/
│   │   └── auth.go                  # Quản lý GitHub token
│   ├── cli/
│   │   └── cli.go                   # Interactive CLI interface
│   ├── github/
│   │   └── client.go                # GitHub API client
│   ├── analyzer/
│   │   ├── analyzer.go              # Phân tích bug logic
│   │   └── analyzer_test.go         # Unit tests
│   └── report/
│       ├── report.go                # Thống kê & reporting
│       └── report_test.go           # Unit tests
├── Formula/
│   └── bug-crawler.rb               # Homebrew formula
├── docs/
│   └── bug-detection-guide.md       # Guide chi tiết phát hiện bug
├── go.mod                           # Go module definitions
├── go.sum                           # Dependency checksums
├── README.md                        # Documentation
├── Makefile                         # Build script
├── USAGE.md                         # Hướng dẫn sử dụng chi tiết
├── TOKEN_SETUP.md                   # Hướng dẫn tạo GitHub token
```

## 🎯 Các Chế Độ Quét Repositories

### 1. Repositories của Bạn (User)
- **Mục đích**: Quét tất cả repositories của tài khoản GitHub của bạn
- **Cách sử dụng**: Chọn option này, ứng dụng sẽ tự động quét
- **Ưu điểm**: Nhanh, không cần nhập gì, phân tích toàn bộ projects cá nhân

### 2. Repositories của Organizations
- **Mục đích**: Quét repositories từ một hoặc nhiều organizations
- **Cách sử dụng**: 
  - Hiển thị danh sách organizations bạn là thành viên
  - Bạn chọn organizations bằng arrow keys + Space
  - Ứng dụng sẽ quét tất cả repositories từ organizations đã chọn
- **Ưu điểm**: Linh hoạt, phân tích team/organization projects

### Cách Chọn Repositories
Khi ứng dụng liệt kê danh sách repositories:
- **Arrow keys (↑↓)**: Di chuyển
- **Space**: Toggle chọn/bỏ chọn
- **Enter**: Xác nhận lựa chọn
- **Hoặc nhập số index**: `1,3,5` để chọn repositories 1, 3, 5
- **Hoặc nhập `all`**: Để chọn tất cả repositories

## 🔍 Phương Pháp Phát Hiện Bug

Ứng dụng hỗ trợ **2 phương pháp** phát hiện bug:

### 1. **Phương Pháp 1: Scan bug (Label-based)**
Phát hiện PR có labels liên quan bug

**Labels được tìm kiếm** (case-insensitive regex):
- Bug-related: `bug`, `fix`, `hotfix`, `critical`
- Error-related: `error`, `issue`

**Cách hoạt động:**
- Kiểm tra tất cả labels của PR
- Nếu có label khớp với pattern → Detect bug → `DetectionType: "label"`

**Ví dụ:**
- PR với labels `["bug", "p0"]` → ✅ Phát hiện
- PR với labels `["documentation"]` → ❌ Không phát hiện

### 2. **Phương Pháp 2: Scan bug_review (Tag-based)**
Phát hiện PR có pattern `bug_review: <number>` trong description

**Pattern tìm kiếm:** `bug_review:\s*(\d+)`

**Cách hoạt động:**
- Tìm pattern `bug_review: <number>` trong PR description
- Extract số lượng bugs từ tag này
- Nếu tìm thấy → Detect bug → `DetectionType: "bug_review"`
- Lưu số lượng bugs trong `BugCount` field

**Ví dụ:**
- Description: "bug_review: 5" → ✅ Phát hiện, BugCount = 5
- Description: "bug_review: 12" → ✅ Phát hiện, BugCount = 12
- Description: "No bugs found" → ❌ Không phát hiện

### Kết Quả Phân Tích

Mỗi PR được lưu với các thông tin:
```go
type BugResult struct {
    PR             *PullRequestData    // Thông tin PR
    IsBugRelated   bool                 // Có liên quan bug?
    DetectionType  string               // "label" hoặc "bug_review"
    MatchedKeyword string               // Label hoặc keyword tìm được
    BugCount       int                  // Số bugs từ bug_review tag
}
```

### Thống Kê

Báo cáo sẽ phân tách:
- **PR phát hiện qua label**: Số lượng
- **PR phát hiện qua bug_review**: Số lượng + Tổng bugs
- **Tỷ lệ bug**: (PR bug-related / Tổng PR) * 100%

## 🔑 Cách Tạo GitHub Personal Access Token

### Bước 1: Đăng Nhập GitHub
- Truy cập https://github.com và đăng nhập tài khoản của bạn

### Bước 2: Mở Settings
- Click vào avatar góc phải → **Settings**

### Bước 3: Developer Settings
- Scroll xuống, click **Developer settings** (phía trái)
- Click **Personal access tokens** → **Tokens (classic)**

### Bước 4: Tạo Token Mới
- Click **Generate new token (classic)**
- Nhập tên token (ví dụ: "bug-crawler")
- **Chọn Scope**: Tích vào `repo` (full control of private repositories)
- Click **Generate token** (dưới cùng)

### Bước 5: Copy Token
- ⚠️ **Quan trọng**: Copy token ngay lập tức (chỉ hiển thị một lần)
- Lưu nơi an toàn

### Sử Dụng Token

**Cách 1: Nhập trong ứng dụng**
- Chạy `bug-crawler` → nhập token khi được yêu cầu

**Cách 2: Environment Variable**
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
bug-crawler
```

**Cách 3: Lưu vào File Config**
- Khi chạy ứng dụng, chọn "Yes" khi được hỏi lưu token
- Token sẽ được lưu vào `~/.config/bug-crawler/token`

## 🎬 Ví Dụ Thực Tế

```bash
$ bug-crawler

🐛 Bug Crawler - GitHub PR Bug Analysis Tool
==========================================

Step 1: GitHub Token
----------------------------------------
GitHub Token: ••••••••••••••••••••••••••••••••••
Lưu token vào file config?
  ▸ Có
    Không
✓ Token đã được lưu

Step 2: Xác thực GitHub
----------------------------------------
👤 Đăng nhập thành công với: khuongdv
✓ Token xác thực thành công

Step 3: Chọn Scan Source
----------------------------------------
Chọn loại để scan
  1. Repositories của tôi (User)
  ▸ 2. Repositories của Organizations
✓ Chọn Organizations

Step 3: Chọn Organizations
----------------------------------------
Chọn Organizations (↑↓=navigate, Space=select, Enter=confirm)
[✓] Golang
[✓] Kubernetes
[ ] Docker
    ...

Selected: 2/10

Step 4: Chọn Repositories
----------------------------------------
📦 Đang quét repositories từ organizations...
🔄 golang...
   ✓ 80 repositories
🔄 kubernetes...
   ✓ 45 repositories

=============================================
📋 Repositories đã chọn (125):
=============================================
 1. ✓ golang/go
 2. ✓ golang/tools
 3. ✓ kubernetes/kubernetes
 ...
=============================================

Step 5: Chọn Khoảng Thời Gian
----------------------------------------
Ngày bắt đầu (YYYY-MM-DD): 2024-01-01
Ngày kết thúc (YYYY-MM-DD): 2024-12-31
✓ Sẽ phân tích PR từ 2024-01-01 đến 2024-12-31

Step 6: Chọn Loại Bug
----------------------------------------
Chọn loại bug để scan
  1. Scan bug (từ labels)
  ▸ 2. Scan bug_review
✓ Sẽ scan bug_review

Step 7: Crawler PR từ GitHub
----------------------------------------
Đang lấy PR từ golang/go...
✓ Tìm được 125 PR
Đang lấy PR từ golang/tools...
✓ Tìm được 35 PR
...

Step 8: Thống Kê Kết Quả
--------------------------------------------
============================================================
THỐNG KÊ BUG
============================================================
Tổng số PR: 1250
PR liên quan bug: 156
  ├─ Phát hiện qua bug_review tag: 120 (Tổng bugs: 245)
  └─ Phát hiện qua label: 36
Tỷ lệ bug: 12.48%
============================================================

CHI TIẾT CÁC PR LIÊN QUAN BUG:
========================================================================================================================
PR#     TITLE                                    AUTHOR      PHÁT HIỆN   BUGS/KEYWORD/LABEL
2345    [Bug] Fix critical memory leak           john-doe    bug_review  5
5678    Fix panic on invalid input               jane-smith  label       bug
8901    Hotfix: Database connection timeout      bob-wilson  label       fix
...

Kết quả đã được export vào: bug_report.csv

✓ Hoàn thành!
```

## 📊 Kết Quả & Hiểu Dữ Liệu

### Định Dạng Kết Quả

**Tóm Tắt Thống Kê:**
```
============================================================
THỐNG KÊ BUG
============================================================
Tổng số PR: 1250
PR liên quan bug: 156
  ├─ Phát hiện qua bug_review tag: 120 (Tổng bugs: 245)
  └─ Phát hiện qua label: 36
Tỷ lệ bug: 12.48%
============================================================
```

**Giải thích:**
- **Tổng số PR**: Tất cả PR trong khoảng thời gian được chọn
- **PR liên quan bug**: PR được phát hiện có liên quan bug
- **bug_review tag**: Số PR có pattern `bug_review: <number>`
  - **Tổng bugs**: Tổng cộng số bugs từ tất cả `bug_review` tags
- **label**: Số PR có labels liên quan bug
- **Tỷ lệ bug**: (PR bug-related / Tổng PR) * 100%

**Chi Tiết PR:**
```
PR#     TITLE                                    AUTHOR      PHÁT HIỆN   BUGS/KEYWORD/LABEL
2345    [Bug] Fix critical memory leak           john-doe    bug_review  5
5678    Fix panic on invalid input               jane-smith  label       bug
```

### File CSV Export

File `bug_report.csv` chứa:
- **PR Number**: Số PR
- **Title**: Tiêu đề PR
- **Author**: Tác giả PR
- **Created Date**: Ngày tạo PR
- **Detection Method**: Cách phát hiện (label/bug_review)
- **Repository**: Repository name
- **Bugs/Keyword/Label**: Số bugs hoặc tên label
- **Opened Date**: Ngày mở PR
- **PR Link**: Link đến PR

## 📚 Dependencies

| Package | Mục Đích | Version |
|---------|---------|---------|
| `github.com/google/go-github` | GitHub API client | v56.0.0 |
| `github.com/manifoldco/promptui` | Interactive CLI prompts | v0.9.0 |
| `github.com/gdamore/tcell` | Terminal UI support | v2.9.0 |

## 🔧 Phát Triển

### Chạy Tests
```bash
go test ./...
```

### Build Binary
```bash
go build -o bug-crawler ./cmd/main.go
```

### Build cho Nhiều OS
```bash
# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o bug-crawler-darwin-amd64 ./cmd/main.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o bug-crawler-darwin-arm64 ./cmd/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o bug-crawler-linux-amd64 ./cmd/main.go
```

## 🚀 Roadmap

- [ ] Support GraphQL queries để fetch dữ liệu nhanh hơn
- [ ] Định nghĩa custom keywords & patterns
- [ ] Export JSON, HTML format
- [ ] Caching PR data để tăng tốc độ
- [ ] Filtering advanced (by author, status, assignee)
- [ ] Web UI dashboard
- [ ] GitHub Actions integration
- [ ] Support batch processing

## 📝 Tài Liệu Khác

- **[USAGE.md](./USAGE.md)** - Hướng dẫn sử dụng chi tiết
- **[TOKEN_SETUP.md](./TOKEN_SETUP.md)** - Cách tạo GitHub token
- **[docs/bug-detection-guide.md](./docs/bug-detection-guide.md)** - Giải thích chi tiết về phát hiện bug

## ❓ FAQ

**Q: Tôi có thể phân tích repositories private không?**  
A: Có, cần GitHub token có scope `repo` để phân tích repositories private.

**Q: Token của tôi có được lưu an toàn không?**  
A: Token được lưu tại `~/.config/bug-crawler/token` trên máy của bạn. Đây là file cục bộ.

**Q: Làm sao để đổi token?**  
A: Xóa file `~/.config/bug-crawler/token` hoặc chạy `bug-crawler` và chọn nhập token mới.

**Q: Phân tích bao lâu?**  
A: Tùy thuộc vào số lượng repositories và PR. Thường từ vài giây đến vài phút.

**Q: Có cách nào để tăng tốc độ không?**  
A: Chọn repositories cụ thể hoặc khoảng thời gian hẹp để giảm số lượng PR cần phân tích.

## 📞 Support

- 🐛 Báo bug: [GitHub Issues](https://github.com/vfa-khuongdv/homebrew-bug-crawler/issues)
- 💬 Thảo luận: [GitHub Discussions](https://github.com/vfa-khuongdv/homebrew-bug-crawler/discussions)

## 📄 License

MIT License

## 👨‍💻 Contributors

- **khuongdv** - Creator & Maintainer

Cảm ơn đã sử dụng Bug Crawler! ⭐ Star repository nếu bạn thấy hữu ích!
