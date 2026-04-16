# Changelog

All notable changes to this project will be documented in this file.

## v2.95.0

### Fixed
- **install.sh now exports PATH for the running script and prints a copy-paste refresh hint** — when invoked via `curl … | bash`, the script runs in a subshell and cannot mutate the parent interactive shell's environment (Unix process isolation). Previously the rc-file write happened silently and users were left wondering why `movie` was "not found" until they opened a new terminal.
- The installer now ends with a clear one-liner (e.g. `export PATH="$HOME/.local/bin:$PATH"` or `fish_add_path …`) the user can paste to refresh their current shell immediately.
- Messaging updated: "Added to ~/.zshrc (new shells will pick it up)" instead of just "Added to ~/.zshrc", removing ambiguity about why the current shell still doesn't see the binary.

## v2.94.0

### Fixed
- **Installer now updates PATH for the current PowerShell session** — previously `[Environment]::SetEnvironmentVariable("PATH", ..., "User")` only affected *future* sessions, so users got "The term 'movie' is not recognized" immediately after install and had to open a new terminal. The installer now also refreshes `$env:PATH` in the running session so `movie` works right away.

## v2.93.0

### Changed
- **Module path renamed** — `github.com/alimtvnetwork/movie-cli-v4` → `github.com/alimtvnetwork/movie-cli-v5` across the entire project (104 files: Go imports, `go.mod`, README, CI workflows, install scripts, docs).
- All GitHub URLs (`github.com/alimtvnetwork/movie-cli-v4` → `…/movie-cli-v5`) updated in README badges, install one-liners, and release-asset URLs.

### Migration
- Run `go mod tidy` after pulling to refresh the module cache.
- Local clones tracking the old remote should update their `origin` URL: `git remote set-url origin https://github.com/alimtvnetwork/movie-cli-v5.git`.

## v2.92.0

### Changed
- **Updater scope locked down** — Go `updater/` package no longer runs `git checkout`/`pull`/`fetch` or build commands. All git mutations and build/deploy steps are delegated to `run.ps1 -Update`.
- Replaced `prepareRepoBranch` with `preflightRepo` in `updater/run.go` — only validates clean working tree before handoff.
- Eliminates "pathspec did not match" errors caused by attempting to check out gitmap release labels (e.g. `release/v1.31.0`) that aren't real branches.

### Added
- New memory constraint `mem://constraints/updater-scope.md` — documents the Go-vs-PowerShell scope split so future contributors don't reintroduce git logic in Go.

## v2.91.0

### Fixed
- Removed branch-switching logic from Go updater that caused `pathspec 'release/v1.31.0' did not match` failures on `movie update`.

## v2.16.0

### Changed
- **Extracted helpers from 12 oversized functions** — all command functions now comply with the ≤50-line guideline
  - `runMovieLsTable` (81→38) — extracted `printLsTableHeader`, `printLsTableRow`, `printLsTableDivider`, `formatRating`
  - `printMediaDetailTable` (75→30) — extracted `buildDetailTableRows` with declarative optional field list
  - `printMediaDetail` (69→33) — extracted `printDetailHeader`, `printDetailIdentifiers`, `printDetailRatings`, `printDetailCredits`, `printDetailFinancials`, `printDetailDescription`, `printDetailFiles`
  - `runMovieScan` (87→37) — extracted `createScanContext`, `executeScan`, `finalizeScan`
  - `runMovieRescan` (71→38) — extracted `fetchRescanEntries`, `processRescanEntries`, `printRescanResult`
  - `runMovieCd` (67→25) — extracted `listScanFolders`, `matchScanFolder`
  - `runWatchLoop` (64→30) — extracted `seedWatchSeen`, `processWatchCycle`, `logWatchScanHistory`
  - `writeHTMLReport` (64→33) — extracted `buildHTMLReportItems`, `splitGenreList`
  - `writeScanSummary` (62→27) — extracted `buildSummaryItems`, `categorizeByGenre`
  - `promptDestination` (62→38) — extracted `loadDestinationDirs`, `loadConfigDir`
  - `runMovieDuplicates` (60→24) — extracted `findDuplicateGroups`, `printDuplicateGroups`, `resolveDuplicatePath`
  - `runMoviePopout` (64) — already well-structured, no change needed

## v2.15.0

### Fixed
- **Update handoff now blocks (foreground)** — changed `cmd.Start()` + `Process.Release()` to `cmd.Run()` so the terminal stays stable and the user sees all worker output
- **Reads `.gitmap/release/latest.json`** — the update command now reads gitmap to determine the correct release branch, and checks out that branch before pulling
- New `updater/gitmap.go` with `GitMapRelease` struct and `readGitMapLatest()` reader

## v2.13.0

### Changed
- **Added 11 option structs for 39 functions with >3 params** — all functions now comply with the ≤3 parameter guideline
  - `StatsCounts` — groups totalMovies/totalTV/total for stats rendering
  - `MoveContext` — groups database/scanner/sourceDir/files/home for move flows
  - `CleanupContext` — groups scanner/database/batchID for popout cleanup
  - `ScanServiceConfig` — groups scanDir/outputDir/database/creds for post-scan services
  - `ScanLoopConfig` — groups client/scanDir/batchID/format flags for main scan loop
  - `ScanOutputOpts` — groups useTable/useJSON flags for scan output
  - `SuggestCollector` — groups client/existingIDs/count for suggestion helpers
  - `LsPage` — groups offset/pageSize/total for list pagination
  - `RecursiveWalkOpts` — groups baseParts/maxDepth for recursive directory walks
  - `ThumbnailInput` — groups client/database/media/posterPath/outputDir for thumbnail downloads
  - `HistoryLogInput` — groups basePath/title/year/fromPath/toPath for history logging
- New `cmd/types.go` with all option struct definitions
- Updated `Test-SourceFiles` in `run.ps1` (88 total)

## v2.12.0

### Changed
- **Split 3 oversized files to under 300 lines** — zero files now exceed the 300-line limit
  - `tmdb/client.go` (323→212) + new `tmdb/http.go` (125) — HTTP retry/response logic extracted
  - `db/schema.go` (314→48) + new `db/schema_tables.go` (263) + `db/schema_indexes.go` (40) — DDL split by table group
  - `cmd/movie_suggest.go` (314→154) + new `cmd/movie_suggest_helpers.go` (169) — genre analysis, discovery, and printing extracted
- Updated `Test-SourceFiles` in `run.ps1` with 4 new files (87 total)

## v2.11.0

### Fixed
- **Self-update now targets the exact original executable path** — the worker passes both deploy directory and binary filename into `run.ps1`, so rebuild/redeploy lands on the same binary that launched `movie update`
- **`run.ps1` deploy override completed** — added `-BinaryNameOverride` alongside `-DeployPath` so update mode no longer depends on `powershell.json` filename defaults

## v2.10.0

### Fixed
- **Self-update now redeploys the exact binary that launched `movie update`** — the original executable path is passed into the handoff worker and forwarded into `run.ps1` as a deploy-path override
- **True handoff flow** — the parent process now starts the copied worker and exits immediately so the original binary can release its file lock before rebuild/deploy
- **Repo-driven rebuild path** — the worker still runs `run.ps1` from the cloned/local GitHub repo, but now targets the original executable directory instead of only the default `powershell.json` deploy path
- **Hidden worker contract tightened** — `update-runner` now requires both `--repo-path` and `--target-binary`

## v2.9.0

### Changed
- **Eliminated all 41 `} else {` violations in Go code** — converted to early-return guard clauses, `continue` in loops, and extracted helpers across 30 files
- Remaining 2 `} else {` are in PowerShell template strings (`updater/script.go`) — not Go logic
- Extracted `runDryRunPlainOutput` and `incrementTypeCountPtr` helpers from oversized `runDryRunScan`
- Simplified `CountMedia` and `ListWatchlist` with early-return patterns
- Files changed: `movie_config.go`, `movie_discover.go`, `movie_history.go`, `movie_info.go`, `movie_logs.go`, `movie_move_batch.go`, `movie_play.go`, `movie_popout.go`, `movie_popout_cleanup.go`, `movie_popout_discover.go`, `movie_redo_exec.go`, `movie_redo_handlers.go`, `movie_rescan.go`, `movie_rest.go`, `movie_rest_report.go`, `movie_scan.go`, `movie_scan_json_output.go`, `movie_scan_loop.go`, `movie_scan_process.go`, `movie_scan_process_helpers.go`, `movie_scan_table.go`, `movie_search.go`, `movie_stats.go`, `movie_tmdb.go`, `movie_undo_handlers.go`, `update.go`, `db/cleanup.go`, `db/media_query.go`, `db/watchlist.go`

## v2.8.0

### Added
- **Pre-build source file validation** — `Test-SourceFiles` function in `run.ps1` validates 83 critical source files exist before compilation, catching missing files early (ported from gitmap-v2 pattern)

## v2.7.0

### Fixed
- **Updater: wrong GitHub repo URL** — `repoURL` used `movie-cli-v5.git` but actual GitHub repo is `movie-cli-v3`; sibling dir search also looked for wrong name
- **run.ps1: stale version file path** — referenced `version/version.go` (renamed to `version/info.go`), causing version detection to fail
- **run.ps1: wrong ldflags module path** — used `movie-cli-v3` instead of `movie-cli-v5` Go module path in build ldflags

### Added
- **run.ps1: `-Deploy` and `-Update` flags** — matches gitmap-v2 pattern; `-Deploy` forces deploy, `-Update` enables rename-first PATH sync
- **run.ps1: PATH binary sync** — when deployed binary differs from PATH binary, auto-syncs with retry and rename-first fallback (ported from gitmap-v2)
- **Updater: passes `-Update` flag to run.ps1** — enables PATH sync during `movie update` flow

## v2.6.0

### Changed
- **P4: Option structs for >3 params** — introduced 6 new input structs (`ErrorLogEntry`, `MoveInput`, `ScanHistoryInput`, `ActionInput`, `WatchlistInput`, `ScanStats`) to replace functions with 4–9 positional parameters; reduced violations from 58 → 47 across 18 files

## v2.5.0

### Changed
- **P3: Replaced all `fmt.Errorf` with `apperror.Wrap()`** — eliminated all 106 `fmt.Errorf` calls across the codebase; all errors now use `apperror.Wrap`, `Wrapf`, or `New` for consistent structured error handling

## v2.4.0

### Changed
- **P2: Eliminated nested ifs** — refactored top 10 worst files using early returns and guard clauses; flattened deeply nested conditionals across scan, move, rename, popout, suggest, rest, and undo commands

## v2.3.0

### Changed
- **Schema fix** — `db/schema.go` multi-value `d.Exec()` error fixed (single-value context)

## v2.2.0

### Changed
- **File splits** — extracted `movie_popout_discover.go`, `movie_popout_cleanup.go`, `movie_scan_loop.go` to keep files under 200 lines; removed duplicate function declarations

## v1.31.0

### Added
- **Version in CLI header box** — scan output now shows `🎬  Movie CLI v1.31.0` centered in the banner (matches gitmap style)

### Changed
- **Spec v1.1** (`spec/10-cli-output-spec.md`) — added flag reference table, JSON item schema, table column definitions, exit codes, flag interaction edge cases, metadata line priority order

## v1.30.0

### Added
- **`--rest` flag for `movie scan`** — starts REST server and opens HTML report in browser after scan completes
- **`--port` flag for `movie scan`** — customize REST server port when using `--rest`
- **REST API request logging** — every HTTP request logged via `errlog.Info` with method, path, status, duration
- **Thumbnails in output folder** — saved to `.movie-output/thumbnails/{slug}-{id}.jpg` with relative paths
- **Thumbnails served via REST** — `/thumbnails/` route serves poster images for the HTML report
- **Gitmap-style CLI output** — box header, numbered items with type icons (🎬/📺), ratings, tree-style output files
- **CLI output spec** — `spec/10-cli-output-spec.md` documents the full output format

### Changed
- Thumbnail naming: `{slug}-{tmdbID}.jpg` flat in `thumbnails/` dir (was nested subdirectories)
- Thumbnail path stored as relative (`thumbnails/xxx.jpg`) for portability
- REST HTML report uses `/thumbnails/` route for images instead of absolute file paths
- Scan output modernized: numbered items, category icons, structured sections

## v1.28.0

### Added
- **Centralized error logging system** (`errlog/logger.go`) — all errors are now logged to:
  - `.movie-output/logs/error.txt` (file-based, append-only, with timestamp/source/stack trace)
  - `error_logs` DB table (queryable, with level/source/function/command/workdir/stack trace)
- **`error_logs` table** (`db/errorlog.go`) — new table with columns: timestamp, level (ERROR/WARN/INFO), source, function, command, work_dir, message, stack_trace; includes `RecentErrorLogs()` query
- **`errlog` package** — `Error()`, `Warn()`, `Info()` functions with automatic caller detection, stack trace capture (errors only), and dual output (file + DB)
- **DB writer injection** — `errlog.SetDBWriter()` allows wiring DB logging without circular imports

### Changed
- **`movie scan` errors** — DB search, stat, insert, update, JSON write, TMDb, and thumbnail errors now use `errlog` instead of raw `fmt.Fprintf(os.Stderr)`
- **`movie rest` errors** — JSON encode, template render, watchlist update, tag add, config read errors now use `errlog`
- **Error entries include**: timestamp, severity, source file:line, function name, CLI command, working directory, message, and full Go stack trace

## v1.27.0

### Changed
- **Modernized HTML report** — complete UI overhaul: sticky toolbar with inline search, genre/rating/sort dropdowns, type filter pills, dark zinc theme, result count, empty state, keyboard shortcut (`/` to search, `Esc` to close modal), responsive layout
- **Search now searches titles, directors, and cast** — not just titles
- **Genre filter dropdown** — auto-populated from scan data
- **Rating filter dropdown** — filter by minimum rating (5+ through 9+)
- **Sort options** — sort by title, rating, or year (ascending/descending)
- **Connected REST indicator** — banner shows green dot when REST server is detected

### Fixed
- **`writeJSON` error swallowed** — `json.Encoder.Encode` error now logged to stderr
- **`tmpl.Execute` error swallowed** — template render error now logged to stderr
- **`GetConfig` errors swallowed** — `tmdb_api_key` and `tmdb_token` config read errors now logged
- **`database.Exec` watchlist update error swallowed** — now logged to stderr
- **`database.AddTag` watched tag error swallowed** — now logged to stderr
- **JS error handling** — all `catch(e)` blocks now show specific error messages; `fetch` non-ok responses show HTTP status/body

## v1.26.0

### Added
- **`GET /` on REST server** — serves a live HTML library report rendered from the database; always up-to-date, no need to open a static file

## v1.25.0

### Added
- **HTML report: tag management** — add/remove tags per card with inline input; tags shown as purple pills with ✕ to remove
- **HTML report: mark watched** — 👁 button marks a movie as watched via REST API; card gets green border and "watched" tag
- **HTML report: similar movies** — 🔍 button opens a modal with TMDb recommendations (poster, title, year, rating, description)
- **HTML report: watched filter** — new "✅ Watched" filter button in the toolbar
- **HTML report: tags auto-load** — when REST server is detected, all tags load automatically on page open

## v1.24.0

### Added
- **`GET/POST/DELETE /api/tags`** — full tag management via REST: list all tags with counts, list tags per media, add tag, remove tag
- **`GET /api/media/{id}/similar`** — fetches TMDb recommendations for a media item
- **`PATCH /api/media/{id}/watched`** — marks a media item as watched (updates watchlist + adds "watched" tag)
- **Refactored REST handlers** — new endpoints in `cmd/movie_rest_handlers.go` to keep files under 200 lines

## v1.23.0

### Added
- **`movie rest --open`** — automatically opens the HTML report in the default browser when the REST server starts; supports macOS (`open`), Windows (`rundll32`), and Linux (`xdg-open`)

## v1.22.0

### Added
- **`movie rest`** — starts a local REST API server (default port 8086, `--port` to override) exposing library endpoints: `GET /api/media`, `GET/DELETE/PATCH /api/media/{id}`, `GET /api/stats`; enables interactive features in the HTML report
- **HTML report** — `movie scan` now generates `report.html` in `.movie-output/` with responsive card layout showing thumbnail, title, year, rating, genre, director, cast, description, and tagline; includes search, filter, and delete via REST API
- **`templates/report.html`** — external HTML template file (not embedded in Go code); bundled via Go `embed` at compile time through `templates/embed.go`

## v1.21.0

### Added
- **`movie db`** — prints the resolved database path, data directory, and record counts for debugging

## v1.20.0

### Changed
- **Renamed `<package>/<package>.go` files** — `db/db.go` → `db/open.go`, `cleaner/cleaner.go` → `cleaner/parse.go`, `updater/updater.go` → `updater/run.go`, `version/version.go` → `version/info.go`; enforced as a permanent naming convention

## v1.19.0

### Added
- **`movie history --format table`** — output move history as a formatted table with columns: #, Title, From, To, Date, Status

## v1.18.0

### Added
- **Binary-relative data storage** — all data (database, thumbnails, JSON metadata) is now stored in `data/` next to the CLI binary, not the working directory
- **`run.ps1` deploys data folder** — build script copies data directory alongside the deployed binary

## v1.17.0

### Added
- **`movie ls --format table`** — output library listing as a formatted table with columns: #, Title, Year, Type, Rating, Genre, Director (no interactive pager)

### Changed
- **Refactored `movie_ls.go`** — split 313-line file into `movie_ls.go` (196), `movie_ls_table.go` (99), and `movie_ls_detail.go` (120)

## v1.16.0

### Changed
- **Refactored `movie_search.go`** — extracted save-and-print logic into `cmd/movie_search_save.go` (135 lines); `movie_search.go` reduced from 240 to 135 lines

## v1.15.0

### Added
- **`movie stats --format table`** — output library statistics as a formatted key-value table with sections for counts, storage, genres, and ratings

## v1.14.0

### Changed
- **Refactored `movie_info.go`** — extracted `fetchMovieDetails` and `fetchTVDetails` into `cmd/movie_fetch_details.go`

## v1.13.0

### Fixed
- **`movie update` fresh-clone flow** — when no local repo exists, a new clone is now reported as bootstrap success instead of incorrectly saying "Already up to date"
- **Self-update specs** — documented repo bootstrap vs existing-repo pull behavior using the GitMap-aligned update flow

## v1.12.0

### Added
- **`movie search --format table`** — output TMDb search results as a formatted table (no interactive prompt); columns: #, Title, Year, Type, Rating, TMDb ID
- **`movie info --format table`** — output media detail as a key-value formatted table; shows all metadata fields dynamically

## v1.11.0

### Added
- **`movie search --format json`** — output TMDb search results as a JSON array to stdout (no interactive prompt); pipeable to `jq` and scripts
- **`movie info --format json`** — output media detail as a JSON object to stdout; includes source field ("local" or "tmdb")

## v1.10.0

### Added
- **`movie ls --format json`** — output entire library as a JSON array to stdout; includes id, title, year, type, ratings, genre, file path, and file size per item
- **`movie stats --format json`** — output library statistics as a JSON object to stdout; includes counts, storage, top genres, and average ratings

## v1.9.0

### Added
- **`movie scan --format json`** — output scan results as structured JSON to stdout for piping to `jq`, scripts, or other tools; includes metadata, counts, and per-item details; works with `--dry-run` too

## v1.8.0

### Fixed
- **`movie scan` no longer fails when TMDb is unset** — media with no TMDb match/key now store `tmdb_id` as `NULL` instead of `0`, so bulk scans no longer hit `UNIQUE constraint failed: media.tmdb_id`
- **Interactive TMDb setup before scan** — when TMDb is not configured, `movie scan` now prompts for a TMDb API key and TMDb access token before scanning starts; leaving both blank continues without metadata
- **TMDb bearer token support** — scan can now authenticate with either `tmdb_api_key` or `tmdb_token`

## v1.7.1

### Changed
- **Refactored `movie_scan.go`** — split from ~500 lines into 4 focused files:
  - `movie_scan.go` (~120 lines) — command definition, orchestrator, helpers
  - `movie_scan_collect.go` (~110 lines) — video file discovery and path utilities
  - `movie_scan_process.go` (~170 lines) — per-file processing and TMDb enrichment
  - `movie_scan_table.go`, `movie_scan_json.go`, `movie_scan_summary.go` — unchanged

## v1.7.0

### Added
- **`movie scan --format table`** — display scan results as a formatted table with columns for #, filename, clean title, year, type, rating, and status; works with `--dry-run` too

## v1.6.0

### Added
- **`movie scan --dry-run`** — preview what would be scanned (files found, cleaned titles, types) without writing to DB or creating `.movie-output/`

## v1.5.0

### Added
- **`movie scan --depth N` (`-d`)** — limit recursive scan to N subdirectory levels (0 = unlimited); e.g. `movie scan -r -d 2`

## v1.4.0

### Fixed
- **`movie update` works from anywhere** — no longer requires CWD to be inside the git repo; finds the repo next to the binary, clones fresh if needed

## v1.3.0

### Added
- **`movie scan --recursive` (`-r`)** — scan all subdirectories recursively instead of just top-level entries; skips `.movie-output` and hidden directories automatically

### Changed
- **Refactored scan internals** — extracted `collectVideoFiles`, `processVideoFile`, and `enrichFromTMDb` helpers for cleaner architecture and reuse

## v1.2.0
### Changed
- **`movie scan` defaults to current directory** — running `movie scan` without arguments now scans the CWD instead of a config-stored `scan_dir` path
- **Scan output to `.movie-output/`** — all scan results (per-item JSON, summary.json with categories/descriptions/metadata) are now written to `.movie-output/` inside the scanned folder

### Added
- **`summary.json`** — comprehensive scan report with total counts, genre-based categories, and full TMDb metadata per item

## v1.1.0

### Fixed
- **`run.ps1` version stamping** — now reads the version from `version/version.go` and injects commit/build date into the correct `version` package variables
- **`run.ps1` version summary** — now reports the binary that was just built/deployed instead of accidentally showing an older `movie` found earlier in `PATH`
- **Deployed changelog visibility** — `run.ps1` now copies `CHANGELOG.md` beside the deployed binary and verifies `movie changelog --latest`

## v0.2.4

### Fixed
- **`GetConfig` false warnings** — `movie_info.go` and `movie_scan.go` now explicitly ignore `sql: no rows in result set` from `GetConfig`, preventing false-positive error messages when config keys are unset
- **Indentation fix** — corrected misleading indentation in `movie_scan.go` error block

### Changed
- **JSON export completeness** — `movie_export.go` now includes all 6 previously missing metadata fields: `Runtime`, `Language`, `Budget`, `Revenue`, `TrailerURL`, `Tagline`

## v0.2.3

### Fixed
- **`db/media.go` silent scan error** — `TopGenres` now returns a wrapped error on `rows.Scan` failure instead of silently using `continue`
- **`movie_info.go` poster error swallowed** — `DownloadPoster` failures now logged to stderr
- **`movie_scan.go` poster error swallowed** — `DownloadPoster` failures now logged to stderr
- **`movie_scan.go` subdirectory read error** — `os.ReadDir` failures in subdirectory scanning now logged instead of silently skipped
- **`movie_undo.go` permission error masked** — `os.Stat` now distinguishes permission errors from file-not-found and logs them separately

## v0.2.2

### Fixed
- **`movie_search.go` unchecked `GetConfig`** — API key lookup now checks for errors before proceeding
- **`movie_suggest.go` unchecked `GetConfig`** — API key lookup now checks for errors and handles `sql: no rows` correctly
- **`movie_resolve.go` unbounded query** — `resolveByTitle` now uses `LIMIT 1` to prevent scanning full table
- **`db/media.go` missing `rows.Err()` check** — `TopGenres` now checks `rows.Err()` after iteration loop

### Changed
- **`movie_search.go` duplicate detail fetch removed** — eliminated redundant `GetMovieDetails`/`GetTVDetails` calls that were already handled by shared `fetchMovieDetails`/`fetchTVDetails` helpers

## v0.2.1

### Fixed
- **`movie_move.go` unchecked error** — `database.GetConfig("movies_dir")` error now handled instead of silently ignored
- **`movie_move.go` unchecked error** — `database.GetConfig("tv_dir")` error now handled instead of silently ignored
- **`movie_move_helpers.go` cross-drive cleanup** — copy+delete fallback now removes the source file after successful copy
- **`movie_rename.go` unchecked `InsertMoveHistory`** — rename history logging error now reported to stderr
- **`movie_play.go` unchecked `exec.Command` error** — player launch error now reported to stderr
- **`movie_stats.go` unchecked `CountMedia`** — movie/TV count errors now handled instead of silently returning zero
- **`movie_watch.go` unchecked `GetConfig`** — API key lookup now checks for errors before proceeding
- **`tmdb/client.go` unchecked `json.NewDecoder` error** — HTTP response body decoding errors now properly returned
- **`updater/updater.go` unchecked exec errors** — `git pull` and `go build` errors now returned instead of silently ignored

## v1.0.0

### Added
- **Batch move** (`movie move --all`) — move all video files at once with auto-routing to movies/TV directories, preview table, and `[y/N]` confirmation
- **JSON metadata export** — `movie scan` now writes per-file JSON metadata to `./data/json/movie/` and `json/tv/`
- **Genre-based discovery** — `movie suggest` uses `DiscoverByGenre` for TMDb genre-based recommendations (3-phase: genre discovery → recommendations → trending fallback)
- **`GenreNameToID()` helper** — reverse genre map in tmdb package for name-to-ID lookups
- **CI pipeline** (`.github/workflows/ci.yml`) — lint (`go vet` + `golangci-lint`), vulnerability scanning (`govulncheck`), parallel test matrix, cross-compiled builds (6 targets), SHA deduplication
- **Release pipeline** (`.github/workflows/release.yml`) — triggers on `release/**` branches and `v*` tags, cross-compiled binaries, SHA256 checksums, version-pinned install scripts, changelog extraction
- **Cross-platform install scripts** — `install.sh` (Linux/macOS) and `install.ps1` (Windows) with checksum verification and PATH setup
- **`.golangci.yml`** — sensible linter defaults (errcheck, govet, staticcheck, gocritic, misspell, errorlint, etc.)
- **Undo confirmation prompt** — `movie undo` shows from/to paths and asks `[y/N]` before reverting
- **Tag command** (`movie tag`) — add, remove, and list tags on media entries
- **Comprehensive CLI help** — root command shows version + categorized help with examples; `movie --version` flag; `movie version` shows Go/OS/arch

### Changed
- **`movie ls`** now only shows scan-indexed items (filters by non-empty `original_file_path`)
- **`movie suggest`** upgraded from recommendations-only to 3-phase strategy (DiscoverByGenre → Recommendations → Trending)
- **Repository migrated** from `movie-cli-v1` to `movie-cli-v2` to `movie-cli-v5` across all imports, workflows, and docs

### Fixed
- Timestamp bug — `saveHistoryLog` now uses `time.Now().Format(time.RFC3339)` instead of hardcoded "now"
- Deduplicated TMDb fetch logic — shared `fetchMovieDetails()`/`fetchTVDetails()` helpers
- Cross-drive move fallback — copy+delete when `os.Rename` fails with `EXDEV`

## v0.1.0

### Added
- Core CLI with Cobra: `hello`, `version`, `self-update` commands
- Movie management: `scan`, `ls`, `search`, `info`, `suggest`, `move`, `rename`, `undo`, `play`, `stats`, `config`
- SQLite database with WAL mode, 5 tables, 7 indexes
- TMDb API client (search, details, credits, recommendations, trending, posters)
- Filename cleaner (junk removal, year extraction, TV detection)
- PowerShell build & deploy pipeline (`run.ps1`)
- Full project specification in `spec/`
