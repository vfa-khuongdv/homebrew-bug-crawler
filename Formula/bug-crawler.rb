class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/vfa-khuongdv/bug-crawler"
  version "1.0.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.1/bug-crawler-darwin-arm64"
      sha256 "9a1cc7e375f5a52c9a964902d46668a46ac632204ceee2d9ec29570afbf86036"
    end
    if Hardware::CPU.intel?
      url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.1/bug-crawler-darwin-amd64"
      sha256 "56d637e3f6c254e5497effe35272cc6ea6fe2f60a3bd076077971a7c35cd2937"
    end
  end

  on_linux do
    url "https://github.com/vfa-khuongdv/homebrew-bug-crawler/releases/download/v1.0.1/bug-crawler-linux-amd64"
    sha256 "21e49e9637e650fe4884defca13b694c99ceb2b4c7d632bf1b7532071fc5018d"
  end

  def install
    bin.install Dir["bug-crawler-*"].first => "bug-crawler"
  end

  test do
    assert_predicate bin/"bug-crawler", :exist?
  end
end
