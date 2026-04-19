# 05 — Acronym MixedCaps (project-specific)

## Symptom

```
Run # Enforce spec/01-coding-guidelines/.../09-acronym-naming.md:
Error: Acronym MixedCaps violations found
  ./doctor/json.go:13:type JSONReport struct {
  ./doctor/json.go:18: Findings  []JSONFinding `json:"findings"`
  ./doctor/json.go:24:type JSONFinding struct {
  ./doctor/json.go:44:func (r *Report) toJSON() JSONReport {
  ./doctor/json.go:56:func toJSONFindings(findings []Finding) []JSONFinding {
Fix: rename IMDb→Imdb, TMDb→Tmdb, API→Api, HTTP→Http, URL→Url,
     JSON→Json, SQL→Sql, HTML→Html, XML→Xml in identifiers.
Error: Process completed with exit code 1.
```

## Trigger

Custom CI guard step in `.github/workflows/ci.yml` that enforces `spec/01-coding-guidelines/03-coding-guidelines-spec/03-golang/09-acronym-naming.md`. The grep pattern is:

```bash
grep -rn -E '\b(IMDb|TMDb|API|HTTP|URL|JSON|SQL|HTML|XML)[A-Z]' --include='*.go' .
```

Any match outside the trailing-initialism allowlist (section 2.1 of the spec) fails the build.

## Root cause

This project **inverts** Go's standard "initialisms stay ALL-CAPS" convention (Effective Go) in favour of strict MixedCaps (`JsonReport`, `ImdbCache`, `ApiKey`). Reasons (from the spec):

- Removes ambiguity at rename time (`IMDbID` vs `IMDbId` vs `IMDBID`)
- Satisfies golangci-lint consistency checks
- Matches the rest of the codebase that grew this way (`db.SetImdbLookup`, `cmd.attachImdbCacheUnless`)

AI defaults to Go-standard `JSONReport`/`HTTPClient`/`APIKey` because that's what `revive`/`stylecheck` enforce by default. **This repo's rule is the opposite.**

## Fix pattern

Rename every identifier — exported and unexported — that contains one of the listed acronyms followed by another uppercase letter:

| Acronym | BAD            | GOOD           |
|---------|----------------|----------------|
| IMDb    | `IMDbCache`    | `ImdbCache`    |
| TMDb    | `TMDbClient`   | `TmdbClient`   |
| API     | `APIKey`       | `ApiKey`       |
| HTTP    | `HTTPClient`   | `HttpClient`   |
| URL     | `URLPattern`   | `UrlPattern`   |
| JSON    | `JSONReport`   | `JsonReport`   |
| SQL     | `SQLBuilder`   | `SqlBuilder`   |
| HTML    | `HTMLEscape`   | `HtmlEscape`   |
| XML     | `XMLParser`    | `XmlParser`    |

Rule applies whether the identifier is exported (`ApiKey`) or unexported (`apiKey`).

### Allowed exceptions (do NOT rename)

1. **Trailing-initialism short locals** — `imdbID`, `tmdbID`, `imgURL`, `reqURL`, `posterURL`. Short-lived, unexported, acronym is the LAST token.
2. **Bare 2-letter `ID`** — `Media.ID`, `User.ID`. Single 2-letter `ID` not preceded by another acronym (e.g. `IMDbID` IS a violation; bare `ID` alone is fine).
3. **SQL string literals** — `SELECT ImdbId, TmdbId FROM Media`. Already MixedCaps in DB schema.
4. **External contracts** — env vars (`TMDB_API_KEY`), HTTP headers (`X-API-Key`), JSON tags (` `json:"id"` `).
5. **Prose comments** — "the TMDb API", "the IMDb id".

### Sed pattern (use word boundaries!)

```bash
sed -i 's/\bJSONReport\b/JsonReport/g' **/*.go
sed -i 's/\bJSONFinding\b/JsonFinding/g' **/*.go
sed -i 's/\bPrintJSON\b/PrintJson/g' **/*.go
sed -i 's/\btoJSON\b/toJson/g' **/*.go
sed -i 's/\btoJSONFindings\b/toJsonFindings/g' **/*.go
sed -i 's/\bemitJSON\b/emitJson/g' **/*.go
sed -i 's/\bdoctorJSON\b/doctorJson/g' **/*.go
```

Without `\b` word boundaries you will corrupt comments and string literals.

## Prevention rule

> **In Go identifiers (exported or unexported), write multi-letter acronyms as MixedCaps: `Json` not `JSON`, `Imdb` not `IMDb`, `Tmdb` not `TMDb`, `Api` not `API`, `Http` not `HTTP`, `Url` not `URL`, `Sql` not `SQL`, `Html` not `HTML`, `Xml` not `XML`. Project rule overrides Effective Go.**

Mirrored to:
- `mem://constraints/acronym-mixedcaps`
- `mem://index.md` Core
- `mem://ci-cd/01-build-fixes-playbook` (full table)

### Granting exemptions for Go stdlib overrides

When a method MUST keep the standard form because it overrides/implements a stdlib interface (e.g. `MarshalJSON`, `UnmarshalJSON`, `ServeHTTP`, `URL()` on `net/url.URL`), add the identifier to the CI guard's allowlist regex in `.github/workflows/ci.yml`:

```yaml
| grep -vE '\b(MarshalJSON|UnmarshalJSON|ServeHTTP|imdbID|tmdbID|...)\b'
```

Document each exemption inline with `// stdlib interface override — keep as-is per spec section 2.x` so future readers know it's intentional.

## History

- **v2.115.0** — wholesale sweep across 17 Go files; spec authored.
- **v2.128.3** — first recurrence after spec: `doctor/json.go` (`JSONReport`/`JSONFinding`/`PrintJSON`/`toJSON`/`toJSONFindings`) + `cmd/doctor.go` (`doctorJSON`/`emitJSON`). Fixed; this issue file authored.
