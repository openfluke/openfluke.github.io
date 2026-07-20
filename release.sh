#!/usr/bin/env bash
# Auto-release openfluke.github.io for the current Welvet scorecard version.
#
# What it does:
#   1. Read version from ../README.md (welvet scorecard)
#   2. Build HTML book + PDF  (python3 _gen_welvet_book.py --run --pdf)
#   3. Commit site source if dirty
#   4. Push main → origin (GitHub Pages)
#   5. Create/update GitHub Release tagged with that version, attach PDF
#
# Usage:
#   ./release.sh                 # full release
#   ./release.sh --dry-run       # build only, no commit/push/release
#   ./release.sh --skip-run      # skip go-run examples (use existing manifest)
#   ./release.sh --no-push       # commit + local release notes, don't push
#
# Needs: python3, google-chrome/chromium, git, and either:
#   - gh (GitHub CLI) authenticated, OR
#   - GITHUB_TOKEN / GH_TOKEN with repo scope
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

WELVET_README="$(cd "$ROOT/.." && pwd)/README.md"
REPO_SLUG="openfluke/openfluke.github.io"
API="https://api.github.com/repos/${REPO_SLUG}"

DRY_RUN=0
SKIP_RUN=0
NO_PUSH=0

for arg in "$@"; do
  case "$arg" in
    --dry-run) DRY_RUN=1 ;;
    --skip-run) SKIP_RUN=1 ;;
    --no-push) NO_PUSH=1 ;;
    -h|--help)
      sed -n '2,20p' "$0"
      exit 0
      ;;
    *)
      echo "unknown flag: $arg" >&2
      exit 2
      ;;
  esac
done

# ---------------------------------------------------------------------------
# Version from welvet/README.md scorecard
# ---------------------------------------------------------------------------
read_version() {
  python3 - <<'PY'
from pathlib import Path
import re
text = Path("../README.md").read_text(encoding="utf-8")
earned = None
m = re.search(r"\*\*(\d+(?:\.\d+)?)\s*/\s*100\*\*\s*pts", text)
if m:
    earned = float(m.group(1))
if earned is None:
    m = re.search(r"\|\s*\*\*Version\*\*\s*\|\s*\*\*(v[\d.]+)\*\*", text)
    if m:
        v = m.group(1)
        earned = 100.0 if v == "v1.0" else float(v[3:]) if v.startswith("v0.") else None
if earned is None:
    raise SystemExit("could not parse Welvet version from ../README.md")
ver = "v1.0" if earned >= 100 else f"v0.{int(round(earned)):02d}"
print(f"{ver} {earned}")
PY
}

# ---------------------------------------------------------------------------
# GitHub helpers — prefer gh, else curl + token
# ---------------------------------------------------------------------------
have_gh() { command -v gh >/dev/null 2>&1; }

token() {
  if [[ -n "${GITHUB_TOKEN:-}" ]]; then
    echo "$GITHUB_TOKEN"
  elif [[ -n "${GH_TOKEN:-}" ]]; then
    echo "$GH_TOKEN"
  else
    echo ""
  fi
}

need_publish_tools() {
  if have_gh; then
    return 0
  fi
  if [[ -n "$(token)" ]]; then
    return 0
  fi
  echo "ERROR: need GitHub CLI (gh) or GITHUB_TOKEN to publish a release." >&2
  echo "  install: https://cli.github.com/  then  gh auth login" >&2
  echo "  or:      export GITHUB_TOKEN=ghp_…" >&2
  exit 1
}

release_exists() {
  local tag="$1"
  if have_gh; then
    gh release view "$tag" --repo "$REPO_SLUG" >/dev/null 2>&1
  else
    local code
    code=$(curl -sS -o /dev/null -w "%{http_code}" \
      -H "Authorization: Bearer $(token)" \
      -H "Accept: application/vnd.github+json" \
      "${API}/releases/tags/${tag}")
    [[ "$code" == "200" ]]
  fi
}

create_or_update_release() {
  local tag="$1"
  local pdf="$2"
  local earned="$3"
  local notes
  notes="$(cat <<EOF
## Welvet feature book ${tag}

Scorecard: **${earned}/100** → version **${tag}** (from \`welvet/README.md\`).

### Assets
- \`$(basename "$pdf")\` — printable feature book (HTML site is on GitHub Pages)

### Site
https://openfluke.github.io/welvet/

### Regenerate
\`\`\`bash
cd openfluke.github.io && ./release.sh
\`\`\`
EOF
)"

  if have_gh; then
    if release_exists "$tag"; then
      echo "  updating existing release ${tag}…"
      gh release upload "$tag" "$pdf" --repo "$REPO_SLUG" --clobber
      gh release edit "$tag" --repo "$REPO_SLUG" \
        --title "Welvet feature book ${tag}" \
        --notes "$notes"
    else
      echo "  creating release ${tag}…"
      gh release create "$tag" "$pdf" --repo "$REPO_SLUG" \
        --title "Welvet feature book ${tag}" \
        --notes "$notes"
    fi
    return
  fi

  # curl fallback
  local auth="Authorization: Bearer $(token)"
  if release_exists "$tag"; then
    echo "  uploading PDF to existing release ${tag}…"
    local id upload
    id=$(curl -sS -H "$auth" -H "Accept: application/vnd.github+json" \
      "${API}/releases/tags/${tag}" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
    # delete prior asset with same name if present
    curl -sS -H "$auth" -H "Accept: application/vnd.github+json" \
      "${API}/releases/${id}/assets" \
      | python3 -c "
import sys, json, os, urllib.request
assets=json.load(sys.stdin)
name=os.path.basename('${pdf}')
tok=os.environ.get('GITHUB_TOKEN') or os.environ.get('GH_TOKEN')
for a in assets:
    if a.get('name')==name:
        req=urllib.request.Request('${API}/releases/assets/'+str(a['id']), method='DELETE',
            headers={'Authorization':'Bearer '+tok,'Accept':'application/vnd.github+json'})
        urllib.request.urlopen(req)
" || true
    upload=$(curl -sS -H "$auth" -H "Accept: application/vnd.github+json" \
      "${API}/releases/${id}" | python3 -c "import sys,json; print(json.load(sys.stdin)['upload_url'].split('{')[0])")
    curl -sS -H "$auth" -H "Content-Type: application/pdf" \
      --data-binary @"$pdf" \
      "${upload}?name=$(basename "$pdf")" >/dev/null
  else
    echo "  creating release ${tag} via API…"
    local body
    body=$(python3 - <<PY
import json
print(json.dumps({
  "tag_name": "${tag}",
  "name": "Welvet feature book ${tag}",
  "body": """${notes}""",
  "draft": False,
  "prerelease": False,
}))
PY
)
    local resp upload
    resp=$(curl -sS -H "$auth" -H "Accept: application/vnd.github+json" \
      -H "Content-Type: application/json" \
      -d "$body" "${API}/releases")
    upload=$(echo "$resp" | python3 -c "import sys,json; print(json.load(sys.stdin)['upload_url'].split('{')[0])")
    curl -sS -H "$auth" -H "Content-Type: application/pdf" \
      --data-binary @"$pdf" \
      "${upload}?name=$(basename "$pdf")" >/dev/null
  fi
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------
if [[ ! -f "$WELVET_README" ]]; then
  echo "ERROR: welvet README not found at $WELVET_README" >&2
  exit 1
fi

read -r VERSION EARNED <<<"$(read_version)"
PDF_NAME="welvet-feature-book-${VERSION}.pdf"
PDF_PATH="dist/${PDF_NAME}"

echo "════════════════════════════════════════"
echo " Welvet book release"
echo " version:  ${VERSION}  (${EARNED}/100)"
echo " repo:     ${REPO_SLUG}"
echo " pdf:      ${PDF_PATH}"
echo "════════════════════════════════════════"

# 1) Build
GEN_FLAGS=(--pdf)
if [[ "$SKIP_RUN" -eq 0 ]]; then
  GEN_FLAGS=(--run --pdf)
fi
echo ""
echo "→ building book ${GEN_FLAGS[*]} …"
# clear Go-example TMPDIR so Chrome PDF step isn't poisoned
env -u TMPDIR -u GOTMPDIR python3 _gen_welvet_book.py "${GEN_FLAGS[@]}"

if [[ ! -f "$PDF_PATH" ]]; then
  echo "ERROR: expected PDF missing: $PDF_PATH" >&2
  exit 1
fi
echo "→ PDF ready: $PDF_PATH ($(du -h "$PDF_PATH" | cut -f1))"

if [[ "$DRY_RUN" -eq 1 ]]; then
  echo ""
  echo "dry-run: skipping commit / push / release"
  echo "would release: ${VERSION} + ${PDF_NAME}"
  exit 0
fi

# 2) Commit site source (never commit dist/ PDF — gitignored)
echo ""
echo "→ git status"
git status --short || true

if [[ -n "$(git status --porcelain)" ]]; then
  echo "→ committing site source…"
  git add -A
  # unstage anything that slipped in from dist/ (belt + suspenders)
  git reset -q -- dist/ 2>/dev/null || true
  git commit -m "$(cat <<EOF
Release Welvet feature book ${VERSION}

Scorecard ${EARNED}/100. HTML site + examples regenerated for GitHub Pages.
PDF published separately as GitHub Release asset ${PDF_NAME}.
EOF
)"
else
  echo "→ working tree clean — nothing to commit"
fi

# 3) Push
if [[ "$NO_PUSH" -eq 1 ]]; then
  echo "→ --no-push: skipping push + GitHub release"
  echo "  local PDF: $PDF_PATH"
  echo "  tag when ready: gh release create ${VERSION} ${PDF_PATH}"
  exit 0
fi

need_publish_tools

echo "→ pushing main…"
git push origin HEAD

# 4) Tag + release with PDF
echo "→ publishing GitHub Release ${VERSION}…"
# ensure annotated tag exists on remote tip
if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "  tag ${VERSION} already exists locally"
else
  git tag -a "$VERSION" -m "Welvet feature book ${VERSION} (${EARNED}/100)"
fi
git push origin "$VERSION" 2>/dev/null || git push origin "refs/tags/${VERSION}"

create_or_update_release "$VERSION" "$PDF_PATH" "$EARNED"

echo ""
echo "════════════════════════════════════════"
echo " Done · ${VERSION}"
echo " Pages:   https://openfluke.github.io/welvet/"
echo " Release: https://github.com/${REPO_SLUG}/releases/tag/${VERSION}"
echo " PDF:     https://github.com/${REPO_SLUG}/releases/latest/download/${PDF_NAME}"
echo "════════════════════════════════════════"
