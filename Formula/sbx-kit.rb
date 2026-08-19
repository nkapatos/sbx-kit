# frozen_string_literal: true

class SbxKit < Formula
  desc "Recipes, kits, and custom images on top of Docker sbx"
  homepage "https://github.com/nkapatos/sbx-kit"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.2.0/sbx-kit_darwin_arm64.tar.gz"
      sha256 "9983b81728003c052dfbb6d171845b90498a2a19679c39d0d679c3b9cdbc47d1"
    end
    on_intel do
      url "https://github.com/nkapatos/sbx-kit/releases/download/v0.2.0/sbx-kit_darwin_amd64.tar.gz"
      sha256 "1b1c1bc39a630e4bf521adf175f9ec7f2bce1a0d890011fdb1bcd83183228242"
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
