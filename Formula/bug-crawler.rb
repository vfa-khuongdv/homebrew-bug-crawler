class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.3"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.3/bug-crawler-darwin-arm64"
      sha256 "8eab210bc4f4ad77b33bb0dce4ccba73d174b6eed7983f20aab3ea528a88f877"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.3/bug-crawler-darwin-amd64"
      sha256 "7246674019bd2b925c0c7c3c5926b87ffd0287189cc35241285a7955d6d48ef5"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.3/bug-crawler-linux-amd64"
    sha256 "d8506d0eec3fd7bef78f1eb5ba4089b620dd5a9b368c82980aa2d699c49f00e3"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
