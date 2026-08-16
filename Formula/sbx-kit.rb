# frozen_string_literal: true

class SbxKit < Formula
  desc "CLI for Docker AI Sandboxes templates, kits, and project init"
  homepage "https://github.com/nkapatos/sbx-kit"
  license "MIT"

  url "https://github.com/nkapatos/sbx-kit/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_TARBALL_SHA256"
  head "https://github.com/nkapatos/sbx-kit.git", branch: "master"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X github.com/nkapatos/sbx-kit/cli/internal/version.Version=#{version}
    ]

    cd "cli" do
      system "go", "build", *std_go_args(ldflags: ldflags.join(" "), output: bin/"sbx-kit"), "./cmd/sbx-kit"
    end

    share_root = share/"sbx-kit"
    share_root.mkpath
    share_root.install "config"
    share_root.install "kits"
    share_root.install "templates" if File.directory?("templates")
    share_root.install "docs" if File.directory?("docs")
    share_root.install "skills" if File.directory?("skills")
  end

  test do
    assert_match "sbx-kit", shell_output("#{bin}/sbx-kit version")
    assert_match "shell", shell_output("#{bin}/sbx-kit recipes")
    assert_match "pi", shell_output("#{bin}/sbx-kit recipes")
    assert_match "recipe", shell_output("#{bin}/sbx-kit concepts")
    assert_match "check", shell_output("#{bin}/sbx-kit check --help")
    assert_match "load", shell_output("#{bin}/sbx-kit image load --help")
  end
end
