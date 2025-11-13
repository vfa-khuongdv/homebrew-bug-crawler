# Bug Crawler - Hướng Dẫn Sử Dụng Chi Tiết

## 📋 Mục Lục
1. [Cài Đặt](#cài-đặt)
2. [Chạy Ứng Dụng](#chạy-ứng-dụng)
3. [Luồng Sử Dụng Chi Tiết](#luồng-sử-dụng-chi-tiết)
4. [Chế Độ Quét Repositories](#chế-độ-quét-repositories)
5. [Cách Phát Hiện Bug](#cách-phát-hiện-bug)
6. [Hiểu Kết Quả](#hiểu-kết-quả)
7. [Xử Lý Lỗi](#xử-lý-lỗi)
8. [Mẹo & Best Practices](#mẹo--best-practices)
9. [Ví Dụ Thực Tế](#ví-dụ-thực-tế)

---

## Cài Đặt

### Cách 1: Cài đặt qua Homebrew (Khuyên dùng)

```bash
# Thêm Homebrew Tap
brew tap vfa-khuongdv/homebrew-bug-crawler

# Cài đặt
brew install bug-crawler

# Chạy ứng dụng
bug-crawler
```

### Cách 2: Build từ Source

**Yêu cầu:**
- Go 1.23 trở lên

```bash
git clone https://github.com/vfa-khuongdv/homebrew-bug-crawler.git
cd homebrew-bug-crawler
go mod tidy
go build -o bug-crawler ./cmd/main.go

# Chạy
./bug-crawler
```

### Yêu Cầu GitHub Token

Ứng dụng cần GitHub Personal Access Token để truy cập API GitHub.

**Cách tạo token:**
1. Truy cập https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Điền tên: `bug-crawler`
4. Chọn scopes:
   - `public_repo` - Để truy cập public repositories
   - `repo` - Để truy cập cả private repositories (nếu cần)
5. Click "Generate token"
6. Copy token (⚠️ chỉ hiển thị một lần)

---

## Chạy Ứng Dụng

**Nếu cài qua Homebrew:**
```bash
bug-crawler
```

**Nếu build từ source:**
```bash
./bug-crawler
```

---

## Luồng Sử Dụng Chi Tiết

### Step 1: GitHub Token

Ứng dụng sẽ tự động:
- **Lần đầu**: Yêu cầu bạn nhập GitHub token
- **Các lần sau**: Sử dụng token đã lưu từ `~/.config/bug-crawler/token`

Bạn có thể chọn:
- ✓ Lưu token: Không phải nhập lại mỗi lần chạy
- ✗ Không lưu: Token chỉ dùng lần này (an toàn hơn)

### Step 2: Chọn Chế Độ Quét Repositories

Ứng dụng sẽ yêu cầu bạn chọn một trong 4 chế độ:

```
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
  2. Quét repositories của user
  3. Quét repositories của organization
  4. Quét repositories của tôi
```

### Step 3: Khoảng Thời Gian (Tuỳ chọn)

```
Từ ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-01-01
Đến ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-12-31
```

Để bỏ qua, chỉ cần nhấn Enter. Khi bỏ qua, ứng dụng sẽ phân tích tất cả PR.

### Step 4: Xem Kết Quả

Ứng dụng sẽ:
- 📊 Phân tích tất cả PR trong khoảng thời gian
- 🐛 Detect bug dựa trên keywords, labels, và bug_review tag
- 📁 Xuất kết quả vào `bug_report.csv`
- 📈 Hiển thị thống kê tóm tắt

---

## Chế Độ Quét Repositories

### Mode 1: Nhập Thủ Công (Manual)

**Khi nào dùng?**
- Bạn biết chính xác repositories muốn phân tích
- Số lượng repositories ít (< 10)
- Repositories nằm rải rác ở nhiều owner khác nhau

**Ví dụ:**
```
Nhập repositories (format: owner/repo, mỗi repo trên một dòng):
Repo 1: golang/go
Repo 2: kubernetes/kubernetes
Repo 3: docker/cli
Repo 4: (nhấn Enter để xong)
```

### Mode 2: Quét Repositories của User

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của một developer
- Bạn muốn so sánh chất lượng code giữa các project của cùng một người
- Kiểm tra repositories công khai của một lập trình viên nổi tiếng

**Ví dụ:**
```
GitHub Username: linus
Đang quét repositories của linus...
✓ Tìm được 8 repositories
(Tự động sử dụng tất cả)
```

### Mode 3: Quét Repositories của Organization

**Khi nào dùng?**
- Bạn là thành viên của một organization
- Bạn muốn phân tích toàn bộ codebase của organization
- Quản lý chất lượng code cho toàn công ty/team

**Ví dụ:**
```
Organization Name: golang
Đang quét repositories của organization golang...
✓ Tìm được 45 repositories
(Tự động sử dụng tất cả)
```

### Mode 4: Quét Repositories của Tôi

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của chính mình
- Quản lý chất lượng code trên tất cả projects của bạn
- Cách nhanh nhất nếu bạn có nhiều repositories

**Ví dụ:**
```
Đang quét repositories của bạn...
✓ Tìm được 15 repositories
(Tự động sử dụng tất cả)
```

---

## Cách Phát Hiện Bug

Tool phát hiện bug dựa trên **3 phương pháp** (theo thứ tự ưu tiên):

### 1. Bug Review Tag (ƯTIÊN NHẤT) ⭐

Nếu PR description chứa `bug_review: <number>`, tool sẽ:
- Nhận diện ngay đó là bug-related PR
- Ghi nhận số lượng bug = `<number>`

**Ví dụ:**
```
PR Title: Fix authentication issues

PR Description:
Fixed several authentication bugs reported by QA team.
bug_review: 3

✓ Detected: Bug-related (3 bugs từ bug_review tag)
```

### 2. Keywords trong Title/Description

Nếu PR title hoặc description chứa một trong các từ khóa sau:

```
bug, fix, hotfix, patch, crash, error, issue, problem, failed, exception, broken
```

**Ví dụ:**
```
PR Title: Fix critical bug in payment processing
✓ Detected: Bug-related (keyword: "bug")

PR Title: Hotfix for session timeout
✓ Detected: Bug-related (keyword: "hotfix")

PR Title: Patch memory leak in cache module
✓ Detected: Bug-related (keyword: "patch")
```

### 3. Labels

Nếu PR có label khớp với pattern (case-insensitive):

```
bug, fix, hotfix, critical, error, issue
```

**Ví dụ:**
```
PR Labels: ["bug", "critical", "security"]
✓ Detected: Bug-related (labels match)

PR Labels: ["feature", "documentation"]
✗ Not detected: No bug-related labels
```

---

## Hiểu Kết Quả

### Tệp CSV Output

Ứng dụng tạo file `bug_report.csv` với các cột:

| Cột | Ý Nghĩa |
|-----|---------|
| Repository | Tên repository (owner/repo) |
| PR Number | Số PR (#123) |
| PR Title | Tiêu đề PR |
| Author | Người tạo PR |
| Created At | Ngày tạo PR |
| Merged At | Ngày merge PR |
| Is Bug Related | True/False - có phải bug-related PR |
| Detection Type | bug_review/keyword/label/both |
| Matched Keyword | Keyword nào được detect (nếu có) |
| Bug Count | Số lượng bug (từ bug_review tag) |
| Description | PR description |

### Thống Kê Tóm Tắt

```
Bug Crawler Report
==================
Repository: golang/go
Khoảng thời gian: 2024-01-01 đến 2024-12-31

Tổng số PR: 350              # Tất cả PR trong khoảng thời gian
PR liên quan bug: 45         # PR được detect là bug-related
Tỷ lệ bug: 12.86%            # (45/350) * 100
Tổng số bug (bug_review): 78 # Tổng bugs từ bug_review tags
```

---

## Xử Lý Lỗi

### 1. Token không hợp lệ

```
❌ Token không hợp lệ hoặc đã hết hạn
```

**Giải pháp:**
- Tạo token mới tại https://github.com/settings/tokens
- Xóa file config: `rm -rf ~/.config/bug-crawler/`
- Chạy lại ứng dụng và nhập token mới

### 2. Rate limit exceeded

```
❌ API rate limit exceeded. Please wait 1 hour before trying again.
```

**Giải pháp:**
- Chạy lại sau 1 giờ (GitHub reset rate limit hàng giờ)
- Hoặc sử dụng token khác nếu có
- Giảm số repositories nếu có thể
- Sử dụng khoảng thời gian nhỏ hơn để giảm số lượng API calls

### 3. Organization/User không tồn tại

```
❌ 404 Not Found: User or organization not found
```

**Giải pháp:**
- Kiểm tra username hoặc organization name
- Đảm bảo tên được nhập chính xác (case-sensitive)
- Truy cập trang GitHub để xác nhận tên đúng

### 4. Repositories không tìm thấy

```
⚠️  Không tìm được repositories nào
```

**Giải pháp:**
- Đối với user: Kiểm tra xem user có repository công khai không
- Đối với organization: Kiểm tra xem bạn có quyền truy cập không
- Token cần có scope `public_repo` hoặc `repo` phù hợp

### 5. Connection timeout

```
❌ Request timeout: Connection failed
```

**Giải pháp:**
- Kiểm tra kết nối internet
- Chạy lại sau vài phút
- Giảm khoảng thời gian để giảm thời gian xử lý

---

## Mẹo & Best Practices

1. **Test trước**: Lần đầu dùng, hãy test với 1-2 repositories nhỏ để làm quen
2. **Lưu token**: Lưu token để không phải nhập mỗi lần chạy
3. **Khoảng thời gian**: 
   - Sử dụng toàn năm (vd: `2024-01-01` đến `2024-12-31`) để báo cáo BPM
   - Hoặc bỏ qua để phân tích tất cả PR
4. **Organization lớn**: Có thể mất vài phút nếu > 100 repositories
5. **CSV Export**: File được lưu tại `bug_report.csv` trong thư mục hiện tại
6. **Tái chạy**: Chạy lại sẽ ghi đè file cũ → lưu với tên khác nếu muốn giữ kết quả cũ:
   ```bash
   mv bug_report.csv bug_report_2024.csv
   bug-crawler  # Chạy lại, tạo file mới
   ```
7. **Performance**: Token có rate limit ~5000 requests/giờ, cứ mỗi repository cần ~2-3 requests

---

## ví Dụ Thực Tế

### Ví dụ 1: Phân tích repositories của một user

```bash
$ bug-crawler

🐛 Bug Crawler - GitHub PR Bug Analysis Tool
==========================================

Step 1: GitHub Token
-----------------------------------------
✓ Token đã được tìm thấy từ file config

Step 2: Chọn Repositories
-----------------------------------------
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
▸ 2. Quét repositories của user

GitHub Username: golang
Đang quét repositories của golang...
✓ Tìm được 12 repositories

Step 3: Khoảng Thời Gian (Tuỳ chọn)
-----------------------------------------
Từ ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-01-01
Đến ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-12-31

Đang phân tích...
✓ go: 125 PR (18 bug-related)
✓ net: 87 PR (9 bug-related)
✓ crypto: 56 PR (7 bug-related)
✓ time: 42 PR (5 bug-related)
...

✅ Hoàn thành!
📊 Thống kê tổng quan:
   - Tổng PR: 350
   - PR liên quan bug: 45
   - Tỷ lệ: 12.86%
   
📄 Kết quả được lưu tại: bug_report.csv
```

### Ví dụ 2: Phân tích repositories của organization

```bash
$ bug-crawler

🐛 Bug Crawler - GitHub PR Bug Analysis Tool
==========================================

Step 1: GitHub Token
-----------------------------------------
Nhập GitHub Token: [token pasted...]
Bạn có muốn lưu token không? (y/n): y
✓ Token đã được lưu

Step 2: Chọn Repositories
-----------------------------------------
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
  2. Quét repositories của user
▸ 3. Quét repositories của organization

Organization Name: kubernetes
Đang quét repositories của organization kubernetes...
✓ Tìm được 142 repositories

(Ứng dụng sẽ phân tích tất cả 142 repositories...)
```

### Ví dụ 3: Nhập thủ công repositories

```bash
$ bug-crawler

Step 2: Chọn Repositories
-----------------------------------------
▸ 1. Nhập thủ công (owner/repo)

Nhập repositories (format: owner/repo, mỗi repo trên một dòng):
Repo 1: golang/go
Repo 2: kubernetes/kubernetes
Repo 3: docker/cli
Repo 4: 
(Ứng dụng sẽ phân tích 3 repositories...)
```

---

## 🔒 An Toàn & Bảo Mật

- **Token lưu cục bộ**: Được lưu tại `~/.config/bug-crawler/token` với quyền `0600`
- **Không gửi token lên**: Chỉ dùng token để gọi GitHub API
- **HTTPS**: Tất cả request đều sử dụng HTTPS
- **Tuỳ chọn lưu**: Bạn quyết định có lưu token hay không
- **Không lưu trữ dữ liệu**: Dữ liệu chỉ được lưu vào CSV local

# Bug Crawler - Hướng Dẫn Sử Dụng Chi Tiết

## 📋 Mục Lục
1. [Cài Đặt](#cài-đặt)
2. [Chạy Ứng Dụng](#chạy-ứng-dụng)
3. [Luồng Sử Dụng Chi Tiết](#luồng-sử-dụng-chi-tiết)
4. [Chế Độ Quét Repositories](#chế-độ-quét-repositories)
5. [Cách Phát Hiện Bug](#cách-phát-hiện-bug)
6. [Hiểu Kết Quả](#hiểu-kết-quả)
7. [Xử Lý Lỗi](#xử-lý-lỗi)

## Cài Đặt

### Cách 1: Cài đặt qua Homebrew (Khuyên dùng)

```bash
# Thêm Homebrew Tap
brew tap vfa-khuongdv/homebrew-bug-crawler

# Cài đặt
brew install bug-crawler

# Chạy ứng dụng
bug-crawler
```

### Cách 2: Build từ Source

**Yêu cầu:**
- Go 1.23 trở lên

```bash
git clone https://github.com/vfa-khuongdv/homebrew-bug-crawler.git
cd homebrew-bug-crawler
go mod tidy
go build -o bug-crawler ./cmd/main.go

# Chạy
./bug-crawler
```

### Yêu cầu GitHub Token

1. Truy cập https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Điền tên: `bug-crawler`
4. Chọn scopes:
   - `public_repo` - Để truy cập public repositories
   - `repo` - Để truy cập cả private repositories (nếu cần)
5. Click "Generate token"
6. Copy token (⚠️ chỉ hiển thị một lần)

---

## Chạy Ứng Dụng

**Nếu cài qua Homebrew:**
```bash
bug-crawler
```

**Nếu build từ source:**
```bash
./bug-crawler
```

---

## Luồng Sử Dụng Chi Tiết

### Step 1: GitHub Token

```
Step 1: GitHub Token
-----------------------------------------
```

Ứng dụng sẽ tự động:
- **Lần đầu**: Yêu cầu bạn nhập GitHub token
- **Các lần sau**: Sử dụng token đã lưu từ `~/.config/bug-crawler/token`

Lựa chọn:
- ✓ Lưu token: Không phải nhập lại mỗi lần chạy
- ✗ Không lưu: Token chỉ dùng lần này (an toàn hơn)

### Step 2: Chọn Chế Độ Quét Repositories

Ứng dụng sẽ yêu cầu bạn chọn một trong 4 chế độ:
```
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
  2. Quét repositories của user
  3. Quét repositories của organization
  4. Quét repositories của tôi
```

### Step 3: Nhập Khoảng Thời Gian (Tuỳ chọn)

```
Từ ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-01-01
Đến ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-12-31
```

Để bỏ qua, chỉ cần nhấn Enter.

### Step 4: Xem Kết Quả

Ứng dụng sẽ:
- 📊 Phân tích tất cả PR trong khoảng thời gian
- 🐛 Detect bug dựa trên keywords, labels, và bug_review tag
- 📁 Xuất kết quả vào `bug_report.csv`

---

## Chế Độ Quét Repositories

### Mode 1: Nhập Thủ Công (Manual)

**Khi nào dùng?**
- Bạn biết chính xác repositories muốn phân tích
- Số lượng repositories ít (< 10)
- Repositories nằm rải rác ở nhiều owner khác nhau

**Ví dụ:**
```
Nhập repositories (format: owner/repo, mỗi repo trên một dòng):
Repo 1: golang/go
Repo 2: kubernetes/kubernetes
Repo 3: docker/cli
Repo 4: (nhấn Enter để xong)
```

### Mode 2: Quét Repositories của User

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của một developer
- Bạn muốn so sánh chất lượng code giữa các project của cùng một người

**Ví dụ:**
```
GitHub Username: linus
Đang quét repositories của linus...
✓ Tìm được 8 repositories
(Tự động sử dụng tất cả)
```

### Mode 3: Quét Repositories của Organization

**Khi nào dùng?**
- Bạn là thành viên của một organization
- Bạn muốn phân tích toàn bộ codebase của organization

**Ví dụ:**
```
Organization Name: golang
Đang quét repositories của organization golang...
✓ Tìm được 45 repositories
(Tự động sử dụng tất cả)
```

### Mode 4: Quét Repositories của Tôi

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của chính mình
- Quản lý chất lượng code trên tất cả projects của bạn

**Ví dụ:**
```
Đang quét repositories của bạn...
✓ Tìm được 15 repositories
(Tự động sử dụng tất cả)
```

---

## Cách Phát Hiện Bug

Tool phát hiện bug dựa trên **3 phương pháp** (theo thứ tự ưu tiên):

### 1. Bug Review Tag (ƯTIÊN NHẤT) ⭐

Nếu PR description chứa `bug_review: <number>`, tool sẽ:
- Nhận diện ngay đó là bug-related PR
- Ghi nhận số lượng bug = `<number>`

**Ví dụ:**
```
PR Title: Fix authentication issues

PR Description:
Fixed several authentication bugs reported by QA team.
bug_review: 3

✓ Detected: Bug-related (3 bugs từ bug_review tag)
```

### 2. Keywords trong Title/Description

Nếu PR title hoặc description chứa một trong các từ khóa:
```
bug, fix, hotfix, patch, crash, error, issue, problem,
failed, exception, broken
```

**Ví dụ:**
```
PR Title: Fix critical bug in payment processing
✓ Detected: Bug-related (keyword: "bug")

PR Title: Hotfix for session timeout
✓ Detected: Bug-related (keyword: "hotfix")
```

### 3. Labels (Issues Label)

Nếu PR có label khớp với pattern:
```
bug, fix, hotfix, critical, error, issue (case-insensitive)
```

**Ví dụ:**
```
PR Labels: ["bug", "critical", "security"]
✓ Detected: Bug-related (labels match)
```

---

## Hiểu Kết Quả

### Tệp CSV Output

Ứng dụng tạo file `bug_report.csv` với các cột:

| Cột | Ý Nghĩa |
|-----|---------|
| Repository | Tên repository (owner/repo) |
| PR Number | Số PR (#123) |
| PR Title | Tiêu đề PR |
| Author | Người tạo PR |
| Created At | Ngày tạo PR |
| Merged At | Ngày merge PR |
| Is Bug Related | True/False - có phải bug-related PR |
| Detection Type | bug_review/keyword/label/both |
| Matched Keyword | Keyword nào được detect (nếu có) |
| Bug Count | Số lượng bug (từ bug_review tag) |
| Description | PR description |

### Thống Kê Tóm Tắt

```
Bug Crawler Report
==================
Repository: golang/go
Khoảng thời gian: 2024-01-01 đến 2024-12-31

Tổng số PR: 350              # Tất cả PR trong khoảng thời gian
PR liên quan bug: 45         # PR được detect là bug-related
Tỷ lệ bug: 12.86%            # (45/350) * 100
Tổng số bug (bug_review): 78 # Tổng bugs từ bug_review tags
```

---

## Xử Lý Lỗi

### 1. Token không hợp lệ
```
❌ Token không hợp lệ hoặc đã hết hạn
```

**Giải pháp:**
- Tạo token mới tại https://github.com/settings/tokens
- Xóa file config: `rm -rf ~/.config/bug-crawler/`
- Chạy lại ứng dụng và nhập token mới

### 2. Rate limit exceeded
```
❌ API rate limit exceeded
```

**Giải pháp:**
- Chạy lại sau 1 giờ (GitHub reset rate limit)
- Hoặc sử dụng token khác nếu có
- Giảm số repositories nếu có thể

### 3. Organization/User không tồn tại
```
❌ 404 Not Found
```

**Giải pháp:**
- Kiểm tra username hoặc organization name
- Đảm bảo tên được nhập chính xác (case-sensitive)

### 4. Repositories không tìm thấy
```
⚠️  Không tìm được repositories nào
```

**Giải pháp:**
- Đối với user: Kiểm tra xem user có repository công khai không
- Đối với organization: Kiểm tra xem bạn có quyền truy cập không
- Token cần có scope `public_repo` hoặc `repo` phù hợp

---

## 💡 Mẹo & Best Practices

1. **Test trước**: Lần đầu dùng, hãy test với 1-2 repositories nhỏ
2. **Lưu token**: Lưu token để không phải nhập mỗi lần
3. **Khoảng thời gian**: Dùng toàn năm (vd: `2024-01-01` đến `2024-12-31`) để báo cáo BPM
4. **Organization lớn**: Có thể mất vài phút nếu > 100 repositories
5. **CSV Export**: File được lưu tại `bug_report.csv` trong thư mục hiện tại
6. **Tái chạy**: Chạy lại với khoảng thời gian khác sẽ ghi đè file cũ (recommended: lưu với tên khác)

---

## 📊 Ví Dụ Thực Tế Toàn Bộ Workflow

```bash
$ bug-crawler

🐛 Bug Crawler - GitHub PR Bug Analysis Tool
==========================================

Step 1: GitHub Token
-----------------------------------------
✓ Token đã được tìm thấy từ file config

Step 2: Chọn Repositories
-----------------------------------------
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
▸ 2. Quét repositories của user

GitHub Username: golang
Đang quét repositories của golang...
✓ Tìm được 12 repositories

Step 3: Khoảng Thời Gian (Tuỳ chọn)
-----------------------------------------
Từ ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-01-01
Đến ngày (YYYY-MM-DD, nhấn Enter để bỏ qua): 2024-12-31

Đang phân tích...
✓ go: 125 PR (18 bug-related)
✓ net: 87 PR (9 bug-related)
✓ crypto: 56 PR (7 bug-related)
...

✅ Hoàn thành!
📄 Kết quả được lưu tại: bug_report.csv
```

---

## 🔒 An Toàn & Bảo Mật

- **Token lưu cục bộ**: Được lưu tại `~/.config/bug-crawler/token`
- **Không gửi token lên**: Chỉ dùng token để gọi GitHub API
- **HTTPS**: Tất cả request đều sử dụng HTTPS
- **Tuỳ chọn lưu**: Bạn quyết định có lưu token hay không

---

## ⚙️ Xử Lý Lỗi

### 2. Rate limit exceeded
```
❌ API rate limit exceeded
```
→ Chạy lại sau 1 giờ

### 3. Organization không tồn tại
```
❌ 404 Not Found
```
→ Kiểm tra username/organization name

---

## 💡 Mẹo

1. Lần đầu: Nhập thủ công 1-2 repositories để test
2. Token: Lưu token để không phải nhập lại mỗi lần
3. Khoảng thời gian: `2024-01-01` đến `2024-12-31` để phân tích cả năm
4. Organization lớn: Có thể mất vài phút nếu > 100 repositories
5. File CSV: Được lưu tại `bug_report.csv`

---

## 📊 Hiểu Kết Quả

```
Tổng số PR: 350              # Tất cả PR trong khoảng thời gian
PR liên quan bug: 45         # PR có từ khóa hoặc label bug-related
Tỷ lệ bug: 12.86%            # (45/350) * 100
```

Keywords: `bug`, `fix`, `hotfix`, `patch`, `crash`, `error`, `issue`, `problem`
