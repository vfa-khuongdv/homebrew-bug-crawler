# 🔐 GitHub Token Setup Guide

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
