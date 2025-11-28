class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.4"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-arm64"
      sha256 "030dabfc4331572f9236b38e73968e58de3285668c2f6113145af34cfaf55184"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-amd64"
      sha256 "d713b952b8a6cbb280350a592ae72b52829f90b161ce740faf2231b12ad63f95"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-linux-amd64"
    sha256 "20f3ff71d7cfbe472dd8a4ee2064a5075a58ba403896964807f923e0d35f4860"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
