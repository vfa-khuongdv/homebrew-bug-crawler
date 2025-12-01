class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.5"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.5/bug-crawler-darwin-arm64"
      sha256 "71836ad540d94fd7102c2b555b1ab94544f04c130d09abd81ce1228425779b45"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.5/bug-crawler-darwin-amd64"
      sha256 "2422a11eaebcccd03477551b34f10bf48a4292d22a59fd849a597357eedfbf02"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.5/bug-crawler-linux-amd64"
    sha256 "2b94185d6fd37d89f379c2d4ca982303bb68fbf7b702d079ead645884931a73d"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
