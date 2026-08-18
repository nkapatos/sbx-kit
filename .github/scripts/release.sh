#!/usr/bin/env bash
# Cut or publish an sbx-kit release: tag, changelog, GitHub Release, Formula.
set -euo pipefail

CLIFF_CONFIG="${CLIFF_CONFIG:-.github/cliff.toml}"
FORMULA="Formula/sbx-kit.rb"
TAP_REPO="${HOMEBREW_TAP_REPO:-${GITHUB_REPOSITORY%%/*}/homebrew-sbx-kit}"

normalize_tag() {
  local v="$1"
  v="${v#v}"
  if [[ ! "$v" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "invalid version: $1 (want X.Y.Z)" >&2
    exit 1
  fi
  echo "v$v"
}

latest_tag() {
  git describe --tags --abbrev=0 --match 'v*.*.*' 2>/dev/null || true
}

same_version() {
  local a b
  a="$(normalize_tag "$1")"
  b="$(normalize_tag "$2")"
  [[ "$a" == "$b" ]]
}

resolve_dispatch_tag() {
  local override="${INPUT_VERSION:-}"
  override="$(echo "$override" | tr -d '[:space:]')"
  if [[ -n "$override" ]]; then
    normalize_tag "$override"
    return
  fi

  local latest
  latest="$(latest_tag)"
  if [[ -z "$latest" ]]; then
    echo "v0.1.0"
    return
  fi

  local bumped
  bumped="$(git cliff --config "$CLIFF_CONFIG" --bumped-version)"
  bumped="$(echo "$bumped" | tr -d '[:space:]')"
  if [[ -z "$bumped" ]]; then
    echo "git-cliff --bumped-version produced no output" >&2
    exit 1
  fi
  bumped="$(normalize_tag "$bumped")"
  if same_version "$bumped" "$latest"; then
    echo "nothing to release: no feat/fix/breaking commits since $latest (pass version to override)" >&2
    exit 1
  fi
  echo "$bumped"
}

create_tag() {
  local tag="$1"
  if git rev-parse "refs/tags/$tag" >/dev/null 2>&1; then
    echo "tag $tag already exists" >&2
    exit 1
  fi
  git tag -a "$tag" -m "$tag"
  git push origin "$tag"
}

build_darwin_asset() {
  local arch="$1" tag="$2"
  local ver="${tag#v}"
  local dir="dist/sbx-kit_darwin_${arch}"
  mkdir -p "$dir"
  (
    cd cli
    CGO_ENABLED=0 GOOS=darwin GOARCH="$arch" go build \
      -ldflags "-s -w -X github.com/nkapatos/sbx-kit/cli/internal/version.Version=${ver}" \
      -o "../${dir}/sbx-kit" ./cmd/sbx-kit
  )
  tar -C "$dir" -czf "dist/sbx-kit_darwin_${arch}.tar.gz" sbx-kit
}

set_formula_asset() {
  local arch="$1" tag="$2" sha="$3" repo="$4"
  local url="https://github.com/${repo}/releases/download/${tag}/sbx-kit_darwin_${arch}.tar.gz"
  awk -v arch="$arch" -v newurl="$url" -v sha="$sha" '
    $0 ~ "sbx-kit_darwin_" arch ".tar.gz" {
      sub(/url "[^"]+"/, "url \"" newurl "\"")
      print
      if (getline > 0) {
        sub(/sha256 "[^"]+"/, "sha256 \"" sha "\"")
        print
      }
      next
    }
    { print }
  ' "$FORMULA" > "${FORMULA}.tmp" && mv "${FORMULA}.tmp" "$FORMULA"
  grep -Fq "$url" "$FORMULA" || { echo "failed to set formula url for ${arch}" >&2; exit 1; }
  grep -Fq "$sha" "$FORMULA" || { echo "failed to set formula sha256 for ${arch}" >&2; exit 1; }
}

tap_git() {
  git -c "http.extraheader=AUTHORIZATION: bearer ${HOMEBREW_TAP_TOKEN}" "$@"
}

# Copy Formula (and a stub README on first run) to nkapatos/homebrew-sbx-kit.
# Requires HOMEBREW_TAP_TOKEN with contents:write on that repo.
push_formula_to_tap() {
  local tag="$1"
  if [[ -z "${HOMEBREW_TAP_TOKEN:-}" ]]; then
    echo "HOMEBREW_TAP_TOKEN is required to push Formula to ${TAP_REPO}" >&2
    echo "Add a PAT (contents:write on ${TAP_REPO}) as repo secret HOMEBREW_TAP_TOKEN." >&2
    exit 1
  fi

  local tap_dir
  tap_dir="$(mktemp -d)/tap"
  if ! tap_git clone "https://github.com/${TAP_REPO}.git" "$tap_dir"; then
    echo "failed to clone https://github.com/${TAP_REPO} — create that public repo first" >&2
    exit 1
  fi

  mkdir -p "$tap_dir/Formula"
  cp "$FORMULA" "$tap_dir/Formula/sbx-kit.rb"
  if [[ ! -f "$tap_dir/README.md" ]]; then
    cat > "$tap_dir/README.md" <<EOF
# sbx-kit

\`\`\`
brew tap ${GITHUB_REPOSITORY%%/*}/sbx-kit
brew install sbx-kit
\`\`\`
EOF
  fi

  (
    cd "$tap_dir"
    git config user.name "github-actions[bot]"
    git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
    branch="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || true)"
    if [[ -z "$branch" || "$branch" == "HEAD" ]]; then
      git checkout -B main
      branch=main
    fi
    git add Formula/sbx-kit.rb README.md
    if git diff --cached --quiet; then
      echo "tap Formula already matches $tag"
      exit 0
    fi
    git commit -m "chore(formula): point tap at ${tag}"
    tap_git push origin "HEAD:${branch}"
  )
}

main() {
  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

  # Tag this checkout (already tested). Do not pull — that would release untested HEAD.
  local tag
  tag="$(resolve_dispatch_tag)"
  echo "Publishing $tag"

  mkdir -p dist
  build_darwin_asset arm64 "$tag"
  build_darwin_asset amd64 "$tag"

  create_tag "$tag"

  git cliff --config "$CLIFF_CONFIG" --latest --strip header > /tmp/release-notes.md
  if [[ ! -s /tmp/release-notes.md ]]; then
    echo "$tag" > /tmp/release-notes.md
  fi

  local assets=(dist/sbx-kit_darwin_arm64.tar.gz dist/sbx-kit_darwin_amd64.tar.gz)
  if gh release view "$tag" >/dev/null 2>&1; then
    echo "GitHub release $tag already exists; updating notes and assets"
    gh release edit "$tag" --notes-file /tmp/release-notes.md
    gh release upload "$tag" "${assets[@]}" --clobber
  else
    gh release create "$tag" --title "$tag" --notes-file /tmp/release-notes.md --verify-tag "${assets[@]}"
  fi

  local sha_arm sha_amd
  sha_arm="$(sha256sum dist/sbx-kit_darwin_arm64.tar.gz | awk '{print $1}')"
  sha_amd="$(sha256sum dist/sbx-kit_darwin_amd64.tar.gz | awk '{print $1}')"
  echo "arm64 sha256=$sha_arm"
  echo "amd64 sha256=$sha_amd"

  git cliff --config "$CLIFF_CONFIG" -o CHANGELOG.md
  set_formula_asset arm64 "$tag" "$sha_arm" "$GITHUB_REPOSITORY"
  set_formula_asset amd64 "$tag" "$sha_amd" "$GITHUB_REPOSITORY"

  git add "$FORMULA" CHANGELOG.md
  if git diff --cached --quiet; then
    echo "Formula and CHANGELOG already match $tag"
  else
    git commit -m "chore(formula): point tap at ${tag}

[skip ci]"
    git push origin "HEAD:${DEFAULT_BRANCH:?}"
  fi

  push_formula_to_tap "$tag"
}

main "$@"
