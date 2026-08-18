#!/usr/bin/env bash
# Cut or publish an sbx-kit release: tag, changelog, GitHub Release, Formula.
set -euo pipefail

CLIFF_CONFIG="${CLIFF_CONFIG:-.github/cliff.toml}"
FORMULA="Formula/sbx-kit.rb"

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

wait_tarball() {
  local url="$1" out="$2"
  local i
  for i in 1 2 3 4 5 6 7 8; do
    if curl -fsSL "$url" -o "$out"; then
      return 0
    fi
    sleep 3
  done
  echo "failed to download $url" >&2
  exit 1
}

update_formula() {
  local tag="$1" sha="$2" repo="$3"
  local url="https://github.com/${repo}/archive/refs/tags/${tag}.tar.gz"
  sed -i -E \
    -e "s|^([[:space:]]*url \")[^\"]+(\")[[:space:]]*$|\\1${url}\\2|" \
    -e "s|^([[:space:]]*sha256 \")[^\"]+(\")[[:space:]]*$|\\1${sha}\\2|" \
    "$FORMULA"
  grep -Fq "$url" "$FORMULA" || { echo "failed to set formula url" >&2; exit 1; }
  grep -Fq "$sha" "$FORMULA" || { echo "failed to set formula sha256" >&2; exit 1; }
}

main() {
  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

  # Tag this checkout (already tested). Do not pull — that would release untested HEAD.
  local tag
  tag="$(resolve_dispatch_tag)"
  create_tag "$tag"
  echo "Publishing $tag"

  git cliff --config "$CLIFF_CONFIG" --latest --strip header > /tmp/release-notes.md
  if [[ ! -s /tmp/release-notes.md ]]; then
    echo "$tag" > /tmp/release-notes.md
  fi

  if gh release view "$tag" >/dev/null 2>&1; then
    echo "GitHub release $tag already exists; updating notes"
    gh release edit "$tag" --notes-file /tmp/release-notes.md
  else
    gh release create "$tag" --title "$tag" --notes-file /tmp/release-notes.md --verify-tag
  fi

  local tar_url tar_path sha
  tar_url="https://github.com/${GITHUB_REPOSITORY}/archive/refs/tags/${tag}.tar.gz"
  tar_path="/tmp/${tag}.tar.gz"
  wait_tarball "$tar_url" "$tar_path"
  sha="$(sha256sum "$tar_path" | awk '{print $1}')"
  echo "tarball sha256=$sha"

  git cliff --config "$CLIFF_CONFIG" -o CHANGELOG.md
  update_formula "$tag" "$sha" "$GITHUB_REPOSITORY"

  git add "$FORMULA" CHANGELOG.md
  if git diff --cached --quiet; then
    echo "Formula and CHANGELOG already match $tag"
    return 0
  fi
  git commit -m "chore(formula): point tap at ${tag}

[skip ci]"
  git push origin "HEAD:${DEFAULT_BRANCH:?}"
}

main "$@"
