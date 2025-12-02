class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.6"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-darwin-arm64"
      sha256 "c30b5d1eedfd8b433e720cbe18fe99c57cacc5f14898689e0a48ac72b3c0331a"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-darwin-amd64"
      sha256 "6c5ddf4e5dee0c70eba7334bb9f9d5603bd14fbf6ec39dcc2de08d10478e63af"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-linux-amd64"
    sha256 "8c00bee2475952b7829321580307509d60633232d28b060cda8c77358c4472f3"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
