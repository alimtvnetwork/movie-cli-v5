# 02 — gofmt: struct tag double-padding

## Symptom

```
Error: doctor/json.go:13:1: File is not properly formatted (gofmt)

  	Schema    string         `json:"schema"`
^
```

## Trigger

`gofmt` rejects manually-padded gaps between field type and struct tag. gofmt's rule: exactly ONE space between the widest field type in the block and the tag column. Any extra space fails.

## Root cause

AI hand-aligned the tag column for visual symmetry while editing struct fields, not realising gofmt would re-collapse the spacing.

## Fix pattern

Never hand-pad struct tags. Either:
- Write minimal spacing (one space) and let gofmt fix the column on save.
- Run `gofmt -w <file>` after any struct edit.

```go
// BAD (manual extra spaces)
type T struct {
    Schema    string         `json:"schema"`
    HasErr    bool           `json:"has_errors"`
}

// GOOD (gofmt-aligned, single space gap)
type T struct {
    Schema string        `json:"schema"`
    HasErr bool          `json:"has_errors"`
}
```

Same rule applies to `const` blocks and `var` blocks.

## Prevention rule

> **Never hand-pad struct tags or const/var blocks. Run `gofmt -w` after every struct edit — accept its output without question.**

## History

- **v2.128.2** — first hit: `doctor/json.go` JSONReport tag column. Fixed.
