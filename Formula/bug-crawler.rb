class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.4"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-arm64"
      sha256 "dae3115962657fdf75ed141925cb1b5bee77bf969ff665cd90a08a0c8247199e"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-darwin-amd64"
      sha256 "730c6cfe72946f1b633b115fd2956d583408743b4b829813ed082ae0cacd3ae7"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.4/bug-crawler-linux-amd64"
    sha256 "77f4f0bda9c95793d4e4521cb80e6fd57903abb0012fd193c584d55bf1103d4d"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
