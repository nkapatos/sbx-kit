# frozen_string_literal: true

class SbxKit < Formula
  desc "CLI for Docker AI Sandboxes templates, kits, and project init"
  homepage "https://github.com/nkapatos/sbx-kit"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.0/sbx-kit_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64"
    end
    on_intel do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.0/sbx-kit_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64"
    end
  end

  head "https://github.com/nkapatos/sbx-kit.git", branch: "master"

  depends_on :macos
  depends_on "go" => :build if build.head?

  def install
    if build.head?
      ldflags = %W[
        -s -w
        -X github.com/nkapatos/sbx-kit/cli/internal/version.Version=#{version}
      ]
      cd "cli" do
        system "go", "build", *std_go_args(ldflags: ldflags.join(" "), output: bin/"sbx-kit"), "./cmd/sbx-kit"
      end
    else
      bin.install "sbx-kit"
    end

    generate_completions_from_executable(bin/"sbx-kit", "completion")
  end

  def caveats
    <<~EOS
      Point sbx-kit at a recipes tree (parent of catalogs):
        sbx-kit setup
    EOS
  end

  test do
    assert_match "sbx-kit", shell_output("#{bin}/sbx-kit version")
    assert_match "#compdef sbx-kit", (share/"zsh/site-functions/_sbx-kit").read
  end
end
