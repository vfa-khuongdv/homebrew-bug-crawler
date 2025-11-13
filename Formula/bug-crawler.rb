class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/homebrew-bug-crawler"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-arm64"
      sha256 "3b0b99477393a59583b49ed179e9b508aef2c3bbef151f41b212b95a612c1b9c"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-darwin-amd64"
      sha256 "d54fb62ab1ce92d3cf6b6f632324d460beed3f813c0447e2218d6fb07cad23d1"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.0/bug-crawler-linux-amd64"
    sha256 "c41179dd6f47707980f5067fd69904a4f12c211a6f1f53509288cce1a3515eaf"
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
