# frozen_string_literal: true

class SbxKit < Formula
  desc "CLI for Docker AI Sandboxes templates, kits, and project init"
  homepage "https://github.com/nkapatos/sbx-kit"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.1/sbx-kit_darwin_arm64.tar.gz"
      sha256 "9eaebce3dd6b7df303b5301bf9765b7d9afb3a948e6faa1f610acbce4c4565ea"
    end
    on_intel do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.1/sbx-kit_darwin_amd64.tar.gz"
      sha256 "7f3c151d43165a44d44aba7a947dfdbdd98b4ff6eae2b152d6f5d91549d4f079"
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
