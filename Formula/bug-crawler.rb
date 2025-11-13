class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/homebrew-bug-crawler"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-amd64"
      sha256 "6faf1adfb17fffd8e0a0e65e29eb4889f8e9d0b4cf9c36966dc8986f859c14c8"
    end
    on_arm do
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-arm64"
      sha256 "4254b2340cb718b657493ca35544f152414108f80de40b4c45917a7d1d5b7b4c"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-linux-amd64"
    sha256 "0a44f8f811ef4a403deef9b1973d32dbb74dca50ef8b10ffa7673e5ab8a015be"
  end

  def install
    # Homebrew downloads the binary directly, just need to make it executable
    bin.install Dir.glob("bug-crawler-*").first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
end
