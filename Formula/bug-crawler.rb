class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-arm64"
      sha256 "ac16fa9e504aa1de7173e1cae6c129103a7f52e38c3aca3ee156a79303976ce1"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-amd64"
      sha256 "a98d263c3cfe793f87b59ae863faf6c8e974f9c156bdddd6537fc4287cce5aff"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-linux-amd64"
    sha256 "757e2b54655686be370f6c76cc52fdc755a0c07412ca20a0c32952dbcdca9e90"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
