---
name: Sibling-repo probing audit
description: Confirms bootstrap.sh and bootstrap.ps1 are the ONLY files implementing -v<N+k> sibling-repo probing. Re-audit before adding any cross-repo URL discovery.
type: audit
---

# Sibling-Repo Probing — Authoritative Locations

**Audit date**: 2026-04-20
**Result**: ✅ PASS — sibling-repo probing is isolated to two files.

## Authoritative implementations

| File | Role |
|------|------|
| `bootstrap.sh` | Bash sibling-repo probing entry point |
| `bootstrap.ps1` | PowerShell sibling-repo probing entry point |

Both probe `https://raw.githubusercontent.com/<owner>/<base>-v<N+k>/main/install.{sh,ps1}` for `k = MAX_LOOKAHEAD..0`, pick the highest HTTP 200, and delegate.

## Search performed

```bash
grep -rln --include='*.sh' --include='*.ps1' --include='*.go' --include='*.yml' \
  -E 'MAX_LOOKAHEAD|find_latest_sibling|Find-LatestSibling|probe_install_script|Test-InstallScript|sibling'
```

Hits outside bootstrap files (all confirmed unrelated — false positives):

| File | Line | Why it's not sibling-probing |
|------|------|-----------------------------|
| `.github/workflows/release.yml` | enforcement step | The `Enforce version-pinning contract` step lists `bootstrap.{sh,ps1}` as **forbidden strings** in generated install scripts. This is the guard, not a violation. |
| `cmd/update.go` | doc comment | "sibling clone" refers to a local sibling **directory** (`movie-cli-v5/` next to the binary), not sibling **repos**. |
| `updater/cleanup.go` | comment | "sibling worker on a different PATH" — local process concern. |
| `updater/repo.go` | candidate resolver | Looks for a sibling repo **directory** on disk for local builds. |

## Forbidden locations (must NEVER probe siblings)

- `.github/workflows/release.yml` — install.{ps1,sh} generators must stay version-pinned. Enforced by the `Enforce version-pinning contract on install scripts` step. Spec: `spec/12-ci-cd-pipeline/06-version-pinned-install-scripts.md`.
- `cmd/update.go` / `updater/*.go` — `movie update` self-updates within the currently-installed repo. Cross-repo upgrades are bootstrap's job, not the updater's.
- `README.md` install one-liners — use `/releases/latest/` redirect (single-repo, latest tag).

## Re-audit trigger

Re-run the grep above before:
- Adding any `-v<N>` URL pattern anywhere
- Editing the install-script generators in `release.yml`
- Adding cross-repo discovery to `movie update`

If the audit grows beyond 2 files, that's a regression — push the new logic back into `bootstrap.{sh,ps1}` or document why it must live elsewhere.
