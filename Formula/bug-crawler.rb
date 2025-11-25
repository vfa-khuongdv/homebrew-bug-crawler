class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-arm64"
      sha256 "6b804472884ab9392069d420856c2f6094d9314cc7e93757d74ddddaba47e975"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-amd64"
      sha256 "fc7f3539f42b3941cb598fcd8111de8e216f61cdaca4ce4016f0c7622d40a90a"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-linux-amd64"
    sha256 "56f0cbcee0807c4f9cccca228e4fc4d1e9c7b28902e57e548a20f2926fb4460a"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
