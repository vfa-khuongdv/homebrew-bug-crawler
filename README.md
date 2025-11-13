# Bug Crawler - GitHub PR Bug Analysis Tool

Ứng dụng terminal Go để phân tích và thống kê số lượng bug từ Pull Request trên GitHub dựa vào description.

## Tính năng

- ✅ Quản lý GitHub token an toàn (lưu vào file config)
- ✅ **4 chế độ quét repositories**:
  - Nhập thủ công
  - Quét repositories của user
  - Quét repositories của organization
  - Quét repositories của tài khoản hiện tại
- ✅ Tự động sử dụng tất cả repositories tìm được (không cần chọn lại)
- ✅ Phân tích PR trong khoảng thời gian tùy chọn
- ✅ Detect bug dựa trên keywords và labels
- ✅ Thống kê chi tiết và tóm tắt
- ✅ Export kết quả dạng CSV

## Cài đặt

### Yêu cầu
- Go 1.21 hoặc cao hơn
- GitHub Token (có thể tạo tại https://github.com/settings/tokens)

### Build

```bash
cd bug_crawler
go mod tidy
go build -o bug-crawler ./cmd/main.go
```

## Sử dụng

### Chạy ứng dụng

```bash
./bug-crawler
```

### Luồng sử dụng

1. **Nhập GitHub Token**: Nhập token của bạn hoặc sử dụng token đã lưu
   - Lần đầu, bạn sẽ được yêu cầu nhập token
   - Token sẽ được lưu vào `~/.config/bug-crawler/token` nếu bạn chọn
   - Lần tiếp theo, token sẽ được tải tự động

2. **Chọn Repositories**: Bạn có 4 cách để quét repositories:
   
   **2a. Nhập thủ công**
   - Nhập danh sách repositories theo format: `owner/repo` (ví dụ: `golang/go`)
   - Nhập từng repo trên một dòng
   - Nhấn Enter 2 lần để kết thúc
   
   **2b. Quét repositories của user**
   - Nhập username GitHub
   - Ứng dụng sẽ tự động quét tất cả repositories của user đó
   - Sử dụng tất cả repositories tìm được
   
   **2c. Quét repositories của organization**
   - Nhập tên organization
   - Ứng dụng sẽ tự động quét tất cả repositories của organization
   - Sử dụng tất cả repositories tìm được
   
   **2d. Quét repositories của bạn**
   - Tự động quét tất cả repositories thuộc tài khoản GitHub của bạn
   - Sử dụng tất cả repositories tìm được

3. **Chọn Khoảng Thời Gian**: Nhập ngày bắt đầu và kết thúc
   - Format: `YYYY-MM-DD` (ví dụ: `2024-01-01`)

4. **Phân Tích**: Ứng dụng sẽ crawler PR và phân tích tự động

5. **Kết Quả**:
   - In tóm tắt thống kê
   - In chi tiết từng PR liên quan bug
   - Export file CSV nếu có PR liên quan bug

## Cấu trúc Thư Mục

```
bug_crawler/
├── cmd/
│   └── main.go              # Entry point
├── pkg/
│   ├── auth/                # Quản lý GitHub token
│   ├── github/              # GitHub API client
│   ├── analyzer/            # Phân tích bug
│   ├── cli/                 # Interactive CLI
│   └── report/              # Thống kê & reporting
├── go.mod                   # Go module
├── go.sum                   # Checksums
└── README.md               # Tài liệu
```

## Chế Độ Quét Repositories

Ứng dụng hỗ trợ 4 chế độ quét repositories:

### 1. Nhập Thủ công (Manual)
- Tự do nhập các repositories theo format `owner/repo`
- Phù hợp khi bạn biết chính xác repositories muốn phân tích
- Ví dụ: `golang/go`, `kubernetes/kubernetes`

### 2. Quét User
- Quét tất cả repositories của một GitHub user
- Sau đó chọn repositories muốn phân tích
- Ví dụ: Quét user `torvalds` để xem repositories của Linus Torvalds

### 3. Quét Organization
- Quét tất cả repositories của một GitHub organization
- Sau đó chọn repositories muốn phân tích
- Ví dụ: Quét organization `golang` để xem tất cả repositories của Go project

### 4. Quét User Hiện Tại
- Tự động quét tất cả repositories của tài khoản GitHub bạn
- Rất hữu ích để phân tích tất cả projects của bạn

### Cách Chọn Repositories
Khi ứng dụng liệt kê danh sách repositories:
- Nhập index repositories (ví dụ: `1,3,5`)
- Hoặc nhập `all` để chọn tất cả

## Phương Pháp Detect Bug

Ứng dụng phát hiện bug dựa trên:

1. **Keywords** trong title/description:
   - `bug`, `fix`, `hotfix`, `patch`
   - `crash`, `error`, `issue`, `problem`
   - `failed`, `exception`, `broken`

2. **Labels** của PR:
   - Regex: `(?i:bug|fix|hotfix|critical|error|issue)`

3. **Cơ chế Phát Hiện**:
   - `keyword`: Phát hiện qua keywords
   - `label`: Phát hiện qua labels
   - `both`: Phát hiện qua cả keywords và labels

## GitHub Token

### Cách tạo Personal Access Token

1. Đăng nhập vào GitHub
2. Vào Settings → Developer settings → Personal access tokens → Tokens (classic)
3. Click "Generate new token (classic)"
4. Chọn scope: `repo` (full control of private repositories)
5. Click "Generate token"
6. Copy token và lưu nơi an toàn

### Env Variable

Bạn cũng có thể sử dụng environment variable:

```bash
export GITHUB_TOKEN=your_token_here
./bug-crawler
```

## Ví dụ

```bash
$ ./bug-crawler

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
✓ Token xác thực thành công

Step 3: Chọn Repositories
----------------------------------------
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
  ▸ 2. Quét repositories của user
    3. Quét repositories của organization
    4. Quét repositories của tôi

GitHub Username: golang
Đang quét repositories của golang...
✓ Tìm được 80 repositories
(Tự động sử dụng tất cả 80 repositories)

Step 4: Chọn Khoảng Thời Gian
----------------------------------------
Ngày bắt đầu (YYYY-MM-DD): 2024-01-01
Ngày kết thúc (YYYY-MM-DD): 2024-12-31

Step 5: Crawler PR từ GitHub
----------------------------------------
Đang lấy PR từ golang/go...
✓ Tìm được 125 PR
Đang lấy PR từ golang/mock...
✓ Tìm được 35 PR
...

Step 6: Thống Kê Kết Quả
--------------------------------------------
============================================================
THỐNG KÊ BUG REVIEW CODE
============================================================
Tổng số PR: 1250
PR liên quan bug: 156
Phát hiện qua keyword: 128
Phát hiện qua label: 45
Tỷ lệ bug: 12.48%
============================================================

CHI TIẾT CÁC PR LIÊN QUAN BUG:
...

Kết quả đã được export vào: bug_report.csv

✓ Hoàn thành!
```

## Dependencies

- `github.com/google/go-github/v56` - GitHub API client
- `github.com/manifoldco/promptui` - Interactive CLI prompts

## Tương Lai

- [ ] Support GraphQL query để fetch dữ liệu nhanh hơn
- [ ] Support định nghĩa custom keywords
- [ ] Support export JSON, HTML format
- [ ] Caching PR data để tăng tốc độ
- [ ] Support filtering by author, status
- [ ] Web UI dashboard

## License

MIT

## Tác giả

Made with ❤️ for Go developers
