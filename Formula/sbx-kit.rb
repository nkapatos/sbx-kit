# frozen_string_literal: true

class SbxKit < Formula
  desc "CLI for Docker AI Sandboxes templates, kits, and project init"
  homepage "https://github.com/nkapatos/sbx-kit"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.2/sbx-kit_darwin_arm64.tar.gz"
      sha256 "c18e0fa5b3ba6edf2e6d7d70081eeb22cc991fc69b6f2eeab6faab7cd2b52496"
    end
    on_intel do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.1.2/sbx-kit_darwin_amd64.tar.gz"
      sha256 "030855847a3cb8bbb90f60164a7e949569fea652b505862085d82b6fee46c088"
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
      Point sbx-kit at a catalog:
        sbx-kit setup
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/sbx-kit version")
    assert_match "#compdef sbx-kit", (share/"zsh/site-functions/_sbx-kit").read
  end
end
