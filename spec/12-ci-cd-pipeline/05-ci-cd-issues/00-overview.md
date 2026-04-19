# CI/CD Lint Issues Log

This folder catalogues **every individual lint/build failure** the CI has flagged, one file per issue. The aggregate playbook lives at `../04-ci-cd-build-fixes.md`; this folder is the granular per-incident record so the AI can pattern-match a new error against an exact past case.

## Naming

`NN-short-slug.md`, e.g. `01-misspell-british-american.md`. Bump `NN` for each new issue.

## File template

Every file follows this structure:

1. **Symptom** — verbatim CI output line(s)
2. **Trigger** — the linter/version/rule that fired
3. **Root cause** — why the AI/human wrote the offending code
4. **Fix pattern** — the exact transformation to apply
5. **Prevention rule** — a one-liner the AI must internalise so the same class never recurs (mirrored to `mem://constraints/...` or `mem://ci-cd/01-build-fixes-playbook`)
6. **History** — versions where this issue was hit + fixed

## Index

| # | Issue | Linter | First hit | Fixed in |
|---|-------|--------|-----------|----------|
| 01 | misspell — British → American | `misspell` (US locale) | v2.128.1 | v2.128.1, v2.128.2 |
| 02 | gofmt — struct tag double-padding | `gofmt` | v2.128.2 | v2.128.2 |
| 03 | gofmt — doc-comment list indent (3-space) | `gofmt` | v2.128.1 | v2.128.1 |
| 04 | govet fieldalignment | `govet` | v2.128.1 | v2.128.1, v2.128.2 |
| 05 | acronym MixedCaps (project-specific) | custom guard (see `spec/01-coding-guidelines/.../09-acronym-naming.md`) | v2.128.3 | v2.128.3 |
| 06 | release published with missing archives (HTTP 404 in installer) | release.yml (no asset-presence guard) | v2.97.0 | v2.128.4 |
| 07 | release dual-trigger race — branch + tag both fire, second run clobbers release | release.yml (`on.push.tags` + `on.push.branches` for self-tagging publish) | v2.97.0 | v2.128.5 |
