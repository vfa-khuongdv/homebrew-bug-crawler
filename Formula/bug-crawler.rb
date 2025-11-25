class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-arm64"
      sha256 "3113ac18abb0251cabffae60f157142f7b169471a50e5aa0493e95666559b192"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-darwin-amd64"
      sha256 "e527627ef63d3841f9d3cdd5f7e312e37b47827cf83ec4d18e03b4f6aaa890a9"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.2/bug-crawler-linux-amd64"
    sha256 "6bc49863ff05e3f6476a88bb6dfef0581172abc1a31c1e3cc6f25b4abe81a402"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
