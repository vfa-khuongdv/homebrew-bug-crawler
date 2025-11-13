# Cách đưa Bug Crawler lên Homebrew

## 1. Chuẩn bị Repository

### Bước 1.1: Thêm Git Tags
```bash
git tag v1.0.0
git push origin v1.0.0
```

### Bước 1.2: Tạo GitHub Releases
1. Vào https://github.com/yourusername/homebrew-bug-crawler/releases
2. Tạo release mới từ tag v1.0.0
3. Tên release: `Bug Crawler v1.0.0`
4. Mô tả: Thêm changelog
5. Upload binaries (xem Bước 2)

## 2. Build Binaries

### Bước 2.1: Build cho macOS (Intel)
```bash
GOOS=darwin GOARCH=amd64 go build -o bug-crawler-darwin-amd64 ./cmd/main.go
```

### Bước 2.2: Build cho macOS (Apple Silicon)
```bash
GOOS=darwin GOARCH=arm64 go build -o bug-crawler-darwin-arm64 ./cmd/main.go
```

### Bước 2.3: Build cho Linux
```bash
GOOS=linux GOARCH=amd64 go build -o bug-crawler-linux-amd64 ./cmd/main.go
```

### Bước 2.4: Tính SHA256
```bash
shasum -a 256 bug-crawler-darwin-amd64
shasum -a 256 bug-crawler-darwin-arm64
shasum -a 256 bug-crawler-linux-amd64
```

## 3. Tạo Homebrew Tap (Repository)

### Bước 3.1: Tạo Repository mới
1. Tạo repository mới trên GitHub: `homebrew-bug-crawler`
2. Clone repository:
```bash
git clone https://github.com/yourusername/homebrew-bug-crawler.git
cd homebrew-bug-crawler
```

### Bước 3.2: Cấu trúc Tap
```
homebrew-bug-crawler/
├── Formula/
│   └── bug-crawler.rb
├── Casks/
├── README.md
└── LICENSE
```

## 4. Viết Homebrew Formula

File: `Formula/bug-crawler.rb`

```ruby
class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/yourusername/bug-crawler"
  
  # Intel Mac
  url "https://github.com/yourusername/bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-amd64"
  sha256 "SHA256_HASH_HERE"
  
  # Apple Silicon
  on_macos do
    on_arm64 do
      url "https://github.com/yourusername/bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-arm64"
      sha256 "SHA256_HASH_HERE"
    end
  end

  def install
    bin.install "bug-crawler" if OS.mac?
  end

  test do
    # Test the binary exists
    assert_predicate bin/"bug-crawler", :exist?
  end
end
```

## 5. Đưa lên Homebrew Core (Tùy chọn)

### Bước 5.1: Điều kiện yêu cầu
- Repository có > 30 stars
- Có documentation rõ ràng
- Được maintain tốt
- Không duplicate functionality

### Bước 5.2: Submit PR
1. Fork: https://github.com/Homebrew/homebrew-core
2. Thêm formula vào `Formula/bug-crawler.rb`
3. Chạy tests: `brew install --verbose bug-crawler.rb`
4. Submit PR

## 6. Sử dụng Tap (Cách dễ nhất)

### Bước 6.1: Thêm Tap
```bash
brew tap yourusername/homebrew-bug-crawler
```

### Bước 6.2: Cài đặt
```bash
brew install bug-crawler
```

### Bước 6.3: Update
```bash
brew upgrade bug-crawler
```

## 7. Cấu hình GitHub Actions (Tự động build)

File: `.github/workflows/release.yml` đã tạo sẵn

Workflow này sẽ tự động:
- Build khi có tag push
- Tạo binaries cho macOS (Intel & ARM64) và Linux
- Upload lên GitHub Releases

## Các bước cụ thể để bắt đầu

1. **Update Repository Main**
   ```bash
   # Thêm version vào code
   # Commit changes
   git add .
   git commit -m "Release v1.0.0"
   ```

2. **Tạo Tag**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Tạo Homebrew Tap**
   ```bash
   # Tạo repo mới: homebrew-bug-crawler
   git clone https://github.com/YOUR_USERNAME/homebrew-bug-crawler.git
   cd homebrew-bug-crawler
   mkdir -p Formula
   # Copy Formula/bug-crawler.rb
   git add .
   git commit -m "Initial commit"
   git push origin main
   ```

4. **Cài đặt từ Tap của riêng mình**
   ```bash
   brew tap your-username/homebrew-bug-crawler https://github.com/YOUR_USERNAME/homebrew-bug-crawler.git
   brew install bug-crawler
   ```

## Lưu ý quan trọng

- Đặt `yourusername` bằng GitHub username thực của bạn
- Thay SHA256 bằng hash thực tế từ binaries
- Thêm LICENSE file
- Viết README tốt
- Test formula trước khi push

## Kiểm tra Formula

```bash
# Validate syntax
brew formula-audit Formula/bug-crawler.rb

# Test install
brew install --verbose Formula/bug-crawler.rb

# Check installation
which bug-crawler
bug-crawler --help
```
