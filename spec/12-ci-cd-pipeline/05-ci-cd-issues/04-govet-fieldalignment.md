# 04 — govet fieldalignment

## Symptom

```
Error: doctor/diagnose.go:39:13: fieldalignment: struct with 64 pointer bytes could be 56 (govet)
type Report struct {
            ^

Error: doctor/json.go:12:17: fieldalignment: struct with 80 pointer bytes could be 72 (govet)
type JSONReport struct {
                ^
```

## Trigger

`govet`'s `fieldalignment` analyzer (enabled in `golangci-lint`) flags struct field orderings that waste padding bytes on 64-bit platforms. CI fails on any wasteful layout.

## Root cause

AI ordered struct fields by **logical grouping** (related fields together) instead of by **size descending** (which minimises padding). Common offender: bools placed in the middle of a struct.

## Fix pattern

Reorder fields by size, descending. Rule of thumb:

1. Pointers / interfaces / strings (16 bytes)
2. Slices / maps (24 bytes — actually largest, but conventionally placed second)
3. Other 8-byte values (int64, uintptr)
4. 4-byte values (int32, float32)
5. **Bools, bytes, small enums LAST** — they pack together at the tail

```go
// BAD (bools mid-struct cause padding holes)
type Report struct {
    Findings []Finding
    Source   string  // padded after slice
    Target   string
    HasErr   bool
}

// GOOD (strings → slice → bools last)
type Report struct {
    Source    string
    Target    string
    DeployDir string
    Findings  []Finding
    HasErr    bool
    HasFix    bool
}
```

When in doubt, run:

```bash
go vet -vettool=$(which fieldalignment) ./...
```

It prints the suggested reordering directly.

## Prevention rule

> **Struct field order: strings/pointers → slices/maps → 8-byte → 4-byte → bools/bytes LAST. After adding any field, mentally count padding or run `fieldalignment`.**

## History

- **v2.128.1** — first hit: `doctor.Report` (64 → 56 bytes). Fixed by moving `Findings` slice to the end.
- **v2.128.2** — recurrence: `doctor.JSONReport` (80 → 72 bytes). Fixed by moving `HasErr`/`HasFix` bools to the end.
