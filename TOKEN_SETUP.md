# 🔐 Token Setup Guide

Hướng dẫn thiết lập token cho các nền tảng: GitHub, Bitbucket, và Backlog.

---

# GitHub Token Setup

## Cần Token Scopes Nào?

Để ứng dụng Bug Crawler có thể truy cập đầy đủ các repositories và pull requests, bạn cần các **scopes** sau:

### 🔑 Required Scopes:

| Scope | Mục đích |
|-------|---------|
| `repo` | ✅ **QUAN TRỌNG** - Truy cập vào repositories công khai và riêng tư |
| `read:org` | ✅ **QUAN TRỌNG** - Đọc thông tin organization |
| `read:user` | ✅ Đọc thông tin user profile |

### 📋 Full Recommended Scopes:
```
repo
read:org
read:user
```

---

## 📝 Hướng dẫn tạo Personal Access Token (PAT)

### Bước 1: Truy cập GitHub Settings
1. Đăng nhập vào GitHub: https://github.com
2. Nhấp vào avatar góc phải → **Settings**
3. Sidebar trái → **Developer settings**
4. **Personal access tokens** → **Tokens (classic)**

### Bước 2: Tạo Token Mới
1. Nhấp **Generate new token** → **Generate new token (classic)**

### Bước 3: Cấu hình Token

**Token name:**
```
bug-crawler-token
```

**Expiration:**
```
90 days (hoặc tùy chọn của bạn)
```

**Select scopes:** ✅ Chọn những scope này:
- ✅ `repo` - Full control of private repositories
- ✅ `read:org` - Read org and team membership
- ✅ `read:user` - Read user profile data

### Bước 4: Tạo & Sao chép Token
1. Nhấp **Generate token**
2. **Sao chép token ngay lập tức** - Bạn sẽ không thể xem lại!
3. Lưu token ở nơi an toàn

---

## 🚀 Sử dụng Token

### Option 1: Nhập khi chạy ứng dụng
```bash
./bug-crawler
# Chương trình sẽ yêu cầu nhập token
# Bạn có thể chọn lưu token vào file config
```

### Option 2: Sử dụng Environment Variable
```bash
export GITHUB_TOKEN="your_token_here"
./bug-crawler
```

### Option 3: Token được lưu tự động
```bash
# Lần đầu tiên
./bug-crawler
# → Nhập token
# → Chọn "Có" để lưu token

# Lần tiếp theo, token sẽ tự động tải từ:
# ~/.config/bug-crawler/token
```

---

## ✅ Kiểm Tra Token Có Đủ Quyền

Khi chạy ứng dụng, nó sẽ hiển thị:

```
Step 2: Xác thực GitHub
✓ Token xác thực thành công
👤 Đăng nhập thành công với: your-username
📊 Rate limit: 4990/5000 requests
```

### Nếu thấy lỗi:

**❌ "Token không hợp lệ hoặc đã hết hạn"**
- Token đã hết hạn → Tạo token mới
- Token bị xóa → Tạo token mới

**❌ "Không tìm thấy repositories"**
- Token không có `repo` scope → Tạo lại token với đủ scopes
- Tài khoản GitHub không có repositories → Tạo hoặc fork repositories

**⚠️ "Không tìm thấy repositories từ organization"**
- Organization không có quyền truy cập → Kiểm tra lại membership
- Token không có `read:org` scope → Tạo lại token với `read:org`

---

## 🔒 Security Tips

✅ **Làm tốt:**
- Sử dụng Personal Access Token (PAT) thay vì password
- Hạn chế scopes (chỉ chọn những scopes cần thiết)
- Đặt expiration date cho token
- Xóa token khi không sử dụng

❌ **KHÔNG làm:**
- Không commit token vào Git repository
- Không chia sẻ token công khai
- Không sử dụng token trong production URLs
- Không lưu token trong plain text files (ngoài ~/.config)

---

## 🆘 Troubleshooting

### Vấn đề 1: "Không tìm thấy repositories của organization"
**Nguyên nhân:** 
- Token không có `read:org` scope
- Không phải member của organization

**Giải pháp:**
1. Tạo token mới với `read:org` scope
2. Đảm bảo bạn là active member của organization

### Vấn đề 2: "Rate limit exceeded"
**Nguyên nhân:**
- Quá nhiều requests trong 1 giờ
- Token hết rate limit

**Giải pháp:**
- Chờ 1 giờ để reset rate limit
- Sử dụng token khác

### Vấn đề 3: "Permission denied"
**Nguyên nhân:**
- Repository hoặc organization là private
- Token không có quyền truy cập

**Giải pháp:**
1. Kiểm tra bạn có quyền truy cập repository/org
2. Tạo token mới với `repo` scope
3. Liên hệ admin organization

---

## 📞 Cần Giúp?

Nếu vấn đề vẫn không giải quyết:

1. **Kiểm tra GitHub Status:** https://www.githubstatus.com/
2. **GitHub Documentation:** https://docs.github.com/en/rest
3. **Personal Access Token Docs:** https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens

---

# Bitbucket Token Setup

## Cần Token Scopes Nào?

### 🔑 Required Scopes:

API Tokens cho Bitbucket Cloud cần các scopes sau để đọc repositories và pull requests:

| Scope | Mục đích |
|-------|---------|
| `read:repository:bitbucket` | ✅ **QUAN TRỌNG** - Truy cập đầy đủ vào repositories |
| `read:pullrequest:bitbucket` | ✅ **QUAN TRỌNG** - Truy cập đầy đủ vào pull requests |
| `read:user:bitbucket` | ✅ Đọc thông tin user profile |
| `read:workspace:bitbucket` | ✅ Đọc thông tin workspace |

---

## 📝 Hướng dẫn tạo API Token (Bitbucket)

### Bước 1: Truy cập Atlassian Account Settings
1. Đăng nhập vào Bitbucket: https://bitbucket.org
2. Nhấp vào **Settings** (cog icon) góc phải trên → **Atlassian account settings**

### Bước 2: Truy cập Security Tab
1. Trong Atlassian Account page, nhấp vào **Security** tab
2. Nhấp vào **Create and manage API tokens**

### Bước 3: Tạo API Token Mới
1. Nhấp **Create API token with scopes**

### Bước 4: Cấu hình Token

**Token name (Tên token):**
```
bug-crawler
```

**Expiry (Hạn sử dụng):**
```
90 days (hoặc tùy chọn của bạn)
```

Nhấp **Next**

### Bước 5: Chọn App
1. Chọn **Bitbucket** làm app
2. Nhấp **Next**

### Bước 6: Chọn Scopes
Chọn các scopes sau:
- ✅ `repository` - Full access to repositories
- ✅ `read:repository:bitbucket` - Read access to repositories
- ✅ `pullrequest` - Full access to pull requests
- ✅ `read:pullrequest:bitbucket` - Read access to pull requests
- ✅ `read:user:bitbucket` - Read user profile data
- ✅ `read:workspace:bitbucket` - Read workspace information

Nhấp **Next**

### Bước 7: Review & Tạo Token
1. Review thông tin token
2. Nhấp **Create token**
3. **Sao chép API token ngay lập tức** - Bạn sẽ không thể xem lại!
4. Lưu token ở nơi an toàn

**⚠️ Lưu ý:** API token chỉ hiển thị một lần. Nếu bạn mất token, bạn phải tạo token mới.

---

## 🚀 Sử dụng Bitbucket Token

### Option 1: Nhập khi chạy ứng dụng
```bash
./bug-crawler --platform bitbucket
# Chương trình sẽ yêu cầu nhập username và API token
# Bạn có thể chọn lưu token vào file config
```

### Option 2: Sử dụng Environment Variables
```bash
export BITBUCKET_USERNAME="your_username"
export BITBUCKET_TOKEN="your_api_token"
./bug-crawler --platform bitbucket
```

### Option 3: Token được lưu tự động
```bash
# Lần đầu tiên
./bug-crawler --platform bitbucket
# → Nhập username
# → Nhập API token
# → Chọn "Có" để lưu token

# Lần tiếp theo, token sẽ tự động tải từ:
# ~/.config/bug-crawler/bitbucket
```

---

## ✅ Kiểm Tra Token Có Đủ Quyền

Khi chạy ứng dụng với Bitbucket, nó sẽ hiển thị:

```
Step 2: Xác thực Bitbucket
✓ Token xác thực thành công
👤 Đăng nhập thành công với: your-username
📊 Rate limit: 60/60 requests per hour
```

### Nếu thấy lỗi:

**❌ "Token không hợp lệ"**
- Username hoặc token không chính xác → Tạo API token mới
- Token đã bị xóa → Tạo API token mới

**❌ "Không tìm thấy repositories"**
- Token không có `repository` scope → Tạo lại API token với đủ scopes
- Tài khoản Bitbucket không có repositories → Tạo repositories mới

---

## 🔒 Bitbucket Security Tips

✅ **Làm tốt:**
- Sử dụng API Token thay vì personal password hoặc app password
- Hạn chế scopes (chỉ chọn những scopes cần thiết)
- Đặt expiration date cho token
- Xóa API token khi không sử dụng hoặc bị compromise
- Lưu token ở nơi an toàn

❌ **KHÔNG làm:**
- Không commit token vào Git repository
- Không chia sẻ token công khai
- Không sử dụng personal password
- Không lưu token trong plain text files (ngoài ~/.config)
- Không sử dụng cùng một token cho nhiều ứng dụng

**Tham khảo thêm:**
- [API Tokens - Bitbucket Support](https://support.atlassian.com/bitbucket-cloud/docs/api-tokens/)
- [Create an API Token - Bitbucket Support](https://support.atlassian.com/bitbucket-cloud/docs/create-an-api-token/)
- [API Token Permissions - Bitbucket Support](https://support.atlassian.com/bitbucket-cloud/docs/api-token-permissions/)

---

# Backlog Token Setup

---

## 📝 Hướng dẫn tạo API Token (Backlog)

### Bước 1: Truy cập Backlog Settings
1. Đăng nhập vào Backlog: https://[your-space].backlog.jp (hoặc .com)
2. Nhấp vào avatar góc phải → **個人設定** (Personal Settings)
3. **API** → **API Tokens**

### Bước 2: Tạo API Token Mới
1. Nhấp **新規作成** (Create New)

### Bước 3: Cấu hình API Token

**説明 (Description):**
```
bug-crawler
```

**有効期間 (Expiration):**
```
1 năm (hoặc tùy chọn của bạn)
```

### Bước 4: Tạo & Sao chép Token
1. Nhấp **作成** (Create)
2. **Sao chép token ngay lập tức** - Bạn sẽ không thể xem lại!
3. Lưu token ở nơi an toàn

---

## 🚀 Sử dụng Backlog Token

### Option 1: Nhập khi chạy ứng dụng
```bash
./bug-crawler --platform backlog
```
# Chương trình sẽ yêu cầu nhập space key, API token
# Bạn có thể chọn lưu token vào file config

### Option 2: Sử dụng Environment Variables
```bash
export BACKLOG_SPACE_KEY="your_space_key"
export BACKLOG_API_TOKEN="your_api_token"
./bug-crawler --platform backlog
```

### Option 3: Token được lưu tự động
```bash
# Lần đầu tiên
./bug-crawler --platform backlog
# → Nhập space key (ví dụ: mycompany)
# → Nhập API token
# → Chọn "Có" để lưu token

# Lần tiếp theo, token sẽ tự động tải từ:
# ~/.config/bug-crawler/backlog
```

---

## ✅ Kiểm Tra Token Có Đủ Quyền

Khi chạy ứng dụng với Backlog, nó sẽ hiển thị:

```
Step 2: Xác thực Backlog
✓ Token xác thực thành công
👤 Space: your-space-key
📊 Projects found: 5
📊 Rate limit: 300/300 requests per hour
```

### Nếu thấy lỗi:

**❌ "Token không hợp lệ"**
- Space key hoặc token không chính xác → Kiểm tra lại setting
- Token đã hết hạn → Tạo API token mới

**❌ "Không có quyền truy cập projects"**
- Token không có quyền truy cập → Kiểm tra user role
- Projects không tồn tại → Tạo projects mới

**❌ "Không tìm thấy issues"**
- Projects không có issues → Tạo issues mới
- Token không có quyền đọc issues → Kiểm tra user role

---

## 🔒 Backlog Security Tips

✅ **Làm tốt:**
- Sử dụng API Token thay vì password
- Đặt expiration date cho token
- Xóa API token khi không sử dụng
- Lưu token ở nơi an toàn

❌ **KHÔNG làm:**
- Không commit token vào Git repository
- Không chia sẻ token công khai
- Không lưu token trong plain text files (ngoài ~/.config)
- Không sử dụng API token trong URLs công khai

---

## 📞 Cần Giúp?

### GitHub
1. **Kiểm tra GitHub Status:** https://www.githubstatus.com/
2. **GitHub Documentation:** https://docs.github.com/en/rest
3. **Personal Access Token Docs:** https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens

### Bitbucket
1. **API Tokens Guide:** https://support.atlassian.com/bitbucket-cloud/docs/api-tokens/
2. **Create an API Token:** https://support.atlassian.com/bitbucket-cloud/docs/create-an-api-token/
3. **API Token Permissions:** https://support.atlassian.com/bitbucket-cloud/docs/api-token-permissions/
4. **API Documentation:** https://developer.atlassian.com/cloud/bitbucket/rest/intro/

### Backlog
1. **Backlog Documentation:** https://backlog.com/ja/
2. **API Documentation:** https://developer.backlog.jp/api/2/
3. **Support:** https://support.backlog.jp/