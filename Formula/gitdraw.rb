class Gitdraw < Formula
  desc "Draw text and art on your GitHub contribution graph"
  homepage "https://github.com/1etu/gitdraw"
  version "1.0.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/1etu/gitdraw/releases/download/v#{version}/gitdraw-cli-macos-arm64.zip"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    end
    on_intel do
      url "https://github.com/1etu/gitdraw/releases/download/v#{version}/gitdraw-cli-macos-amd64.zip"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end
  end

  def install
    if Hardware::CPU.arm?
      bin.install "gitdraw-cli-macos-arm64" => "gitdraw"
    else
      bin.install "gitdraw-cli-macos-amd64" => "gitdraw"
    end
  end

  test do
    assert_match "gitdraw", shell_output("#{bin}/gitdraw --version")
  end
end
