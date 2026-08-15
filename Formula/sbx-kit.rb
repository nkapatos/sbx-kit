# frozen_string_literal: true

# Homebrew formula for sbx-kit (tap = this repo).
#
#   brew tap nkapatos/sbx-kit https://github.com/nkapatos/sbx-kit
#   brew install sbx-kit
#   brew install --HEAD nkapatos/sbx-kit/sbx-kit   # before first tag
#
# Install layout:
#   bin/sbx-kit
#   share/sbx-kit/{config,kits,templates,docs}

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

  def caveats
    <<~EOS
      sbx-kit data (example recipes) lives in:
        #{share}/sbx-kit

      Override with SBX_TREE=/path/to/checkout for local development or another
      templates/kits/catalog tree.

      Host vault (created on first run/rm/upgrade/status):
        ~/.local/share/sbx-kit/profiles/   portable state archives
        ~/.local/state/sbx-kit/            project↔recipe bindings

      You still need the Docker `sbx` CLI (>= 0.34.0; kits are schemaVersion 1). Check with:
        sbx-kit version

      Stock recipes (sbx kind + kits):
        sbx-kit concepts
        sbx-kit recipes
        sbx secret set deepseek
        sbx-kit run shell --yes
        sbx-kit check

      Custom images (optional):
        sbx-kit image ls
        sbx-kit image load --engine docker kit-core
        sbx-kit image load --engine docker kit-cursor
        sbx template ls
        sbx-kit run kit-cursor --yes
    EOS
  end

  test do
    assert_match "sbx-kit", shell_output("#{bin}/sbx-kit version")
    assert_match "shell", shell_output("#{bin}/sbx-kit recipes")
    assert_match "recipe", shell_output("#{bin}/sbx-kit concepts")
    assert_match "check", shell_output("#{bin}/sbx-kit check --help")
    assert_match "load", shell_output("#{bin}/sbx-kit image load --help")
  end
end
