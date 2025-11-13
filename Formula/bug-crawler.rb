class BugCrawler < Formula
  desc "GitHub PR Bug Analysis Tool - analyze and count bugs from pull requests"
  homepage "https://github.com/yourusername/homebrew-bug-crawler"
  url "https://github.com/yourusername/homebrew-bug-crawler/releases/download/v{version}/bug-crawler-darwin-amd64"
  sha256 "{sha256_hash}"
  version "{version}"

  def install
    bin.install "bug-crawler-darwin-amd64" => "bug-crawler"
  end

  test do
    system "#{bin}/bug-crawler", "--version"
  end
end
