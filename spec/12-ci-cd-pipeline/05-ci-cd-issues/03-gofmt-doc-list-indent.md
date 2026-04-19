# 03 — gofmt: doc-comment numbered list indent

## Symptom

```
Error: doctor/diagnose.go:6:1: File is not properly formatted (gofmt)

//   1. Active PATH binary differs from powershell.json deployPath
^
```

## Trigger

`gofmt`'s doc-comment formatter (Go 1.19+) normalises numbered lists in `//`-comments to **2-space indent**: `//  1. foo`. A 3-space indent (`//   1. foo`) is rewritten on save and the diff fails CI.

## Root cause

AI used the visually pleasing 3-space indent for numbered lists in package doc comments, matching how Markdown renders. gofmt has its own opinion.

## Fix pattern

Use 2-space indent for numbered lists in doc comments:

```go
// GOOD
//
//  1. First item
//  2. Second item

// BAD (3 spaces — gofmt rewrites)
//
//   1. First item
//   2. Second item
```

Bonus: if the comment block is on the file's `package` line and is the only / first such block in the package, start with `// Package <name>` (godoc convention). Never use `// filename.go — desc`.

## Prevention rule

> **Doc-comment numbered lists use 2-space indent. The first/only file in a package starts with `// Package <name>`, not `// filename.go — desc`.**

## History

- **v2.128.1** — first hit: `doctor/diagnose.go` numbered list of 4 failure modes. Fixed; converted to `// Package doctor` doc comment.
