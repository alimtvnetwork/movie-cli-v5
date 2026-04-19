# 01 — misspell: British → American

## Symptom

```
Error: cmd/doctor.go:6:4: `catalogued` is a misspelling of `cataloged` (misspell)
Error: updater/self_replace.go:31:4: `Behaviour` is a misspelling of `Behavior` (misspell)
Error: doctor/diagnose.go:39:16: `optimised` is a misspelling of `optimized` (misspell)
```

## Trigger

`golangci-lint` enables the `misspell` linter with `locale: US`. Any British spelling in code, comments, doc strings, error messages, CHANGELOG, or spec docs fails the build.

## Root cause

AI defaults to British spellings in prose ("behaviour", "optimised", "catalogued", "normalised") because both forms are grammatically valid and the model wasn't constrained to one locale.

## Fix pattern

Mechanical swap. Full table:

| British (BAD)  | American (GOOD) |
|----------------|-----------------|
| behaviour      | behavior        |
| optimised      | optimized       |
| catalogued     | cataloged       |
| normalised     | normalized      |
| analyse(d)     | analyze(d)      |
| organise(d)    | organize(d)     |
| recognise(d)   | recognize(d)    |
| serialise(d)   | serialize(d)    |
| initialise(d)  | initialize(d)   |
| utilise(d)     | utilize(d)      |
| minimise(d)    | minimize(d)     |
| maximise(d)    | maximize(d)     |
| colour         | color           |
| favour(ite)    | favor(ite)      |
| centre(d)      | center(ed)      |
| labour         | labor           |
| licence (n.)   | license         |
| cancelled      | canceled        |
| travelling     | traveling       |
| modelling      | modeling        |
| labelled       | labeled         |
| fulfil(led)    | fulfill(ed)     |

## Prevention rule

> **Before typing ANY word ending in `-our`, `-ise`, `-ised`, `-isation`, `-yse`, `-tre`, `-lled`, `-lling`, swap to American form.**

Mirrored to:
- `mem://ci-cd/01-build-fixes-playbook` (full table)
- `mem://index.md` Core (one-liner reminder)

## History

- **v2.128.1** — first hit: `catalogued` (cmd/doctor.go), `behaviour` (updater/self_replace.go). Fixed same release.
- **v2.128.2** — recurrence: `optimised` (doctor/diagnose.go). Fixed. Spec entry created.
