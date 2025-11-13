# Bug Crawler - Hướng Dẫn Sử Dụng Chi Tiết

## 📋 Mục Lục
1. [Chuẩn Bị](#chuẩn-bị)
2. [Chế Độ Quét Repositories](#chế-độ-quét-repositories)
3. [Ví Dụ Thực Tế](#ví-dụ-thực-tế)
4. [Xử Lý Lỗi](#xử-lý-lỗi)

## Chuẩn Bị

### 1. Tạo GitHub Personal Access Token

1. Truy cập https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Điền tên: `bug-crawler`
4. Chọn scopes:
   - `public_repo` - Để truy cập public repositories
   - `repo` - Để truy cập cả private repositories (nếu cần)
5. Click "Generate token"
6. Copy token (⚠️ chỉ hiển thị một lần)

### 2. Build Ứng Dụng

```bash
cd /path/to/bug_crawler
go mod tidy
go build -o bug-crawler ./cmd/main.go
```

## Chế Độ Quét Repositories

### Mode 1: Nhập Thủ công (Manual)

**Khi nào dùng?**
- Bạn biết chính xác repositories muốn phân tích
- Số lượng repositories ít (< 10)
- Repositories nằm rải rác ở nhiều owner khác nhau

**Ví dụ:**
```
Step 3: Chọn Repositories
----------------------------------------
Chọn cách quét repositories
  1. Nhập thủ công (owner/repo)
  ▸ Repo 1: golang/go
    Repo 2: kubernetes/kubernetes
    Repo 3: docker/cli
    Repo 4: (nhấn Enter để xong)
```

### Mode 2: Quét Repositories của User

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của một developer
- Bạn muốn so sánh chất lượng code giữa các project của cùng một người
- Số lượng repositories của user lớn (>20)

**Ví dụ:**
```
GitHub Username: torvalds
Đang quét repositories của torvalds...
✓ Tìm được 5 repositories
(Tự động sử dụng: torvalds/subsurface, torvalds/linux, ...)
```

### Mode 3: Quét Repositories của Organization

**Khi nào dùng?**
- Bạn là thành viên của một organization
- Bạn muốn phân tích toàn bộ codebase của organization
- Tổ chức có nhiều repositories

**Ví dụ:**
```
Organization Name: kubernetes
Đang quét repositories của organization kubernetes...
✓ Tìm được 142 repositories
(Tự động sử dụng tất cả 142 repositories)
```

### Mode 4: Quét Repositories của Tôi

**Khi nào dùng?**
- Bạn muốn phân tích tất cả repositories của chính mình
- Quản lý chất lượng code trên tất cả projects của bạn
- Cách nhanh nhất nếu bạn có nhiều repositories

**Ví dụ:**
```
Đang quét repositories của bạn...
✓ Tìm được 23 repositories
(Tự động sử dụng tất cả 23 repositories)
```

---

## Xử Lý Lỗi

### 1. Token không hợp lệ
```
❌ Token không hợp lệ hoặc đã hết hạn
```
→ Tạo token mới tại https://github.com/settings/tokens

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
