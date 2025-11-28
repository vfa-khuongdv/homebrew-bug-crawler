class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.4"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-arm64"
      sha256 "adee96adaa7ac770162596917bfa7ccc0e3a20b1f76341fd9bb060f197cc56c5"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-amd64"
      sha256 "114377b068df899fa0f31c801ff65740c838ef0a1a775745118d81ecd2c66258"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-linux-amd64"
    sha256 "122c73cc0eeeed975c088203b23458d67e75f2cdc833a6b8e38d8c6aa7011b15"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
