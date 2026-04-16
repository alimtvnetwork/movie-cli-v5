package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v4/apperror"
)

// Cleanup removes leftover temp binaries and backup files from previous updates.
// It cleans BOTH the active binary's directory AND the deploy directory (if different),
// covering:
//   - <name>-update-*[.exe]   handoff copies from any prior PID
//   - movie-update-*[.exe]    legacy handoff copies (pre-rebrand)
//   - *.old                   rename-first deploy backups
//   - *.bak                   legacy backup files
func Cleanup() (int, error) {
	cleaned := 0

	selfPath, err := os.Executable()
	if err != nil {
		return 0, apperror.Wrap("cannot determine executable path", err)
	}
	if resolved, evalErr := filepath.EvalSymlinks(selfPath); evalErr == nil {
		selfPath = resolved
	}

	baseName := binaryBaseName(selfPath)
	dirs := candidateDirs(selfPath)

	for _, dir := range dirs {
		cleaned += cleanDir(dir, baseName, selfPath)
	}

	return cleaned, nil
}

// binaryBaseName returns the running binary name without extension,
// e.g. "gitmap.exe" -> "gitmap", "movie" -> "movie".
func binaryBaseName(selfPath string) string {
	name := filepath.Base(selfPath)
	ext := filepath.Ext(name)
	if ext != "" {
		name = strings.TrimSuffix(name, ext)
	}
	return name
}

// candidateDirs returns the unique directories to scan for leftovers:
// the active binary's directory and the OS temp directory.
func candidateDirs(selfPath string) []string {
	seen := map[string]struct{}{}
	var dirs []string
	add := func(d string) {
		abs, err := filepath.Abs(d)
		if err != nil {
			return
		}
		key := strings.ToLower(abs)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		dirs = append(dirs, abs)
	}
	add(filepath.Dir(selfPath))
	add(os.TempDir())
	return dirs
}

// cleanDir removes all known leftover patterns in a single directory.
func cleanDir(dir, baseName, selfPath string) int {
	patterns := []string{
		filepath.Join(dir, baseName+"-update-*"),
		filepath.Join(dir, "movie-update-*"), // legacy
		filepath.Join(dir, "*.old"),
		filepath.Join(dir, "*.bak"),
	}
	cleaned := 0
	for _, pattern := range patterns {
		cleaned += cleanGlob(pattern, selfPath)
	}
	return cleaned
}

// cleanGlob removes files matching a glob pattern, skipping the running binary.
func cleanGlob(pattern, selfPath string) int {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}

	cleaned := 0
	for _, match := range matches {
		if shouldSkip(match, selfPath) {
			continue
		}
		if err := os.Remove(match); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not remove %s: %v\n", filepath.Base(match), err)
			continue
		}
		fmt.Printf("  Removed: %s\n", filepath.Base(match))
		cleaned++
	}
	return cleaned
}

// shouldSkip returns true when the match is the running binary or a directory.
func shouldSkip(match, selfPath string) bool {
	abs, _ := filepath.Abs(match)
	if pathsEqual(abs, selfPath) {
		return true
	}
	info, err := os.Stat(match)
	if err != nil {
		return true
	}
	return info.IsDir()
}

// pathsEqual compares two paths case-insensitively on Windows.
func pathsEqual(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}
