class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.6"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-darwin-arm64"
      sha256 "89cd5b01e97f2baec71576ed05e2877df3197a00cc3090066a3af998593a9f3a"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-darwin-amd64"
      sha256 "f7cb97856d33e5d571aa87616865fe85365200152ac16d1f32db417b2f1b5f81"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.6/bug-crawler-linux-amd64"
    sha256 "501e502f97ff295640a683d60e9b5da373b3d913591b6646526fece7879347ce"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
