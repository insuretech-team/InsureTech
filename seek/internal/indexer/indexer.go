package indexer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/bmatcuk/doublestar/v4"
)

var (
	defaultExcludedDirs = map[string]struct{}{
		".cache":       {},
		".git":         {},
		".hg":          {},
		".idea":        {},
		".next":        {},
		".nuxt":        {},
		".projindex":   {},
		".seek":        {},
		".svn":         {},
		".turbo":       {},
		".vscode":      {},
		"bin":          {},
		"build":        {},
		"coverage":     {},
		"dist":         {},
		"logs":         {},
		"node_modules": {},
		"obj":          {},
		"out":          {},
	}
	defaultExcludedExts = map[string]struct{}{
		".7z":      {},
		".a":       {},
		".ai":      {},
		".apk":     {},
		".avif":    {},
		".bmp":     {},
		".class":   {},
		".dll":     {},
		".doc":     {},
		".docx":    {},
		".dylib":   {},
		".eot":     {},
		".eps":     {},
		".exe":     {},
		".gif":     {},
		".gz":      {},
		".ico":     {},
		".jar":     {},
		".jpeg":    {},
		".jpg":     {},
		".lock":    {},
		".mov":     {},
		".mp3":     {},
		".mp4":     {},
		".o":       {},
		".otf":     {},
		".pdf":     {},
		".png":     {},
		".psd":     {},
		".so":      {},
		".sqlite":  {},
		".sqlite3": {},
		".svgz":    {},
		".tar":     {},
		".tif":     {},
		".tiff":    {},
		".ttf":     {},
		".wasm":    {},
		".webm":    {},
		".webp":    {},
		".woff":    {},
		".woff2":   {},
		".xls":     {},
		".xlsx":    {},
		".zip":     {},
	}
)

type BuildOptions struct {
	Root           string
	DBPath         string
	IncludeExt     []string
	ExcludeDirs    []string
	ExcludeGlobs   []string
	MaxFileSize    int64
	IncludeHidden  bool
	FollowSymlinks bool
	Reset          bool
	Verbose        bool
	Progress       func(BuildProgress)
}

type BuildStats struct {
	Root           string
	DBPath         string
	IndexedFiles   int
	UpdatedFiles   int
	UnchangedFiles int
	RemovedFiles   int
	SkippedFiles   int
	Errors         int
	IndexedLines   int
	IndexedTerms   int
}

type BuildProgress struct {
	Stage          string
	Path           string
	ScannedFiles   int
	IndexedFiles   int
	UpdatedFiles   int
	UnchangedFiles int
	SkippedFiles   int
	Errors         int
}

type GrepOptions struct {
	Pattern    string
	PathGlobs  []string
	Limit      int
	Fixed      bool
	IgnoreCase bool
}

func Build(opts BuildOptions) (BuildStats, error) {
	var stats BuildStats

	if opts.MaxFileSize <= 0 {
		opts.MaxFileSize = 4 * 1024 * 1024
	}

	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return stats, err
	}
	root = filepath.Clean(root)

	dbPath, err := ResolveDBPath(root, opts.DBPath)
	if err != nil {
		return stats, err
	}

	store, err := OpenStore(dbPath)
	if err != nil {
		return stats, err
	}
	defer store.Close()

	if err := store.Init(); err != nil {
		return stats, err
	}
	if opts.Reset {
		if err := store.Reset(); err != nil {
			return stats, err
		}
	}
	if err := store.AssertRoot(root, opts.Reset); err != nil {
		return stats, err
	}

	states, err := store.FileStates()
	if err != nil {
		return stats, err
	}

	stats.Root = root
	stats.DBPath = dbPath
	seen := make(map[string]struct{})
	excludedDirs := normalizeDirSet(opts.ExcludeDirs)
	includeExts := normalizeExtSet(opts.IncludeExt)
	dbRootDir := filepath.Dir(dbPath)
	scannedFiles := 0
	lastProgress := time.Now()

	reportProgress := func(stage, path string, force bool) {
		if opts.Progress == nil {
			return
		}
		if !force && scannedFiles > 0 && scannedFiles%250 != 0 && time.Since(lastProgress) < 2*time.Second {
			return
		}
		lastProgress = time.Now()
		opts.Progress(BuildProgress{
			Stage:          stage,
			Path:           path,
			ScannedFiles:   scannedFiles,
			IndexedFiles:   stats.IndexedFiles,
			UpdatedFiles:   stats.UpdatedFiles,
			UnchangedFiles: stats.UnchangedFiles,
			SkippedFiles:   stats.SkippedFiles,
			Errors:         stats.Errors,
		})
	}

	reportProgress("start", root, true)

	walkErr := filepath.WalkDir(root, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			stats.Errors++
			if opts.Verbose {
				fmt.Printf("walk error: %v\n", walkErr)
			}
			return nil
		}
		if currentPath == root {
			return nil
		}
		if !opts.FollowSymlinks && entry.Type()&fs.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			stats.SkippedFiles++
			return nil
		}

		rel, err := filepath.Rel(root, currentPath)
		if err != nil {
			return err
		}
		rel = normalizePath(rel)
		name := entry.Name()

		if entry.IsDir() {
			if shouldSkipDir(rel, name, opts.IncludeHidden, excludedDirs, opts.ExcludeGlobs, dbRootDir, currentPath) {
				return filepath.SkipDir
			}
			return nil
		}
		scannedFiles++

		if shouldSkipFile(rel, name, includeExts, opts.IncludeHidden, opts.ExcludeGlobs, currentPath, dbPath) {
			stats.SkippedFiles++
			reportProgress("skip", rel, false)
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			stats.Errors++
			return nil
		}
		if info.Size() > opts.MaxFileSize {
			stats.SkippedFiles++
			reportProgress("skip", rel, false)
			return nil
		}

		seen[rel] = struct{}{}
		if current, ok := states[rel]; ok && current.Size == info.Size() && current.ModTimeUnix == info.ModTime().Unix() {
			stats.UnchangedFiles++
			reportProgress("unchanged", rel, false)
			return nil
		}

		content, isText, err := readTextFile(currentPath, opts.MaxFileSize)
		if err != nil {
			stats.Errors++
			if opts.Verbose {
				fmt.Printf("read error: %s (%v)\n", rel, err)
			}
			return nil
		}
		if !isText {
			stats.SkippedFiles++
			reportProgress("skip", rel, false)
			return nil
		}

		lines := splitLines(content)
		termHits, termTotal := tokenizeLines(lines)
		digest := sha256.Sum256([]byte(content))

		if err := store.UpsertFileIndex(rel, FileRecord{
			Size:        info.Size(),
			ModTimeUnix: info.ModTime().Unix(),
			IndexedAt:   time.Now().Unix(),
			SHA256:      hex.EncodeToString(digest[:]),
			LineCount:   len(lines),
		}, lines, termHits); err != nil {
			return err
		}

		stats.IndexedFiles++
		stats.IndexedLines += len(lines)
		stats.IndexedTerms += termTotal
		if _, exists := states[rel]; exists {
			stats.UpdatedFiles++
		}
		reportProgress("index", rel, false)

		return nil
	})
	if walkErr != nil {
		return stats, walkErr
	}

	removed, err := store.DeleteMissing(seen)
	if err != nil {
		return stats, err
	}
	stats.RemovedFiles = removed

	if err := store.SetBuildMeta(root, time.Now().UTC()); err != nil {
		return stats, err
	}
	reportProgress("done", root, true)

	return stats, nil
}

func ResolveDBPath(root, dbPath string) (string, error) {
	if strings.TrimSpace(root) == "" {
		return "", errors.New("root path is required")
	}
	if strings.TrimSpace(dbPath) == "" {
		return filepath.Join(root, ".seek", "index.sqlite"), nil
	}
	if filepath.IsAbs(dbPath) {
		return filepath.Clean(dbPath), nil
	}
	return filepath.Abs(dbPath)
}

func shouldSkipDir(rel, name string, includeHidden bool, extraDirs map[string]struct{}, excludeGlobs []string, dbRootDir, currentPath string) bool {
	if samePath(currentPath, dbRootDir) {
		return true
	}
	if !includeHidden && strings.HasPrefix(name, ".") {
		return true
	}
	if _, found := defaultExcludedDirs[strings.ToLower(name)]; found {
		return true
	}
	if _, found := extraDirs[strings.ToLower(name)]; found {
		return true
	}
	return matchesAnyGlob(rel, excludeGlobs)
}

func shouldSkipFile(rel, name string, includeExts map[string]struct{}, includeHidden bool, excludeGlobs []string, absPath, dbPath string) bool {
	if samePath(absPath, dbPath) {
		return true
	}
	if !includeHidden && strings.HasPrefix(name, ".") {
		return true
	}
	if matchesAnyGlob(rel, excludeGlobs) {
		return true
	}
	ext := strings.ToLower(filepath.Ext(name))
	if len(includeExts) > 0 {
		_, ok := includeExts[ext]
		return !ok
	}
	if _, found := defaultExcludedExts[ext]; found {
		return true
	}
	return false
}

func normalizeDirSet(values []string) map[string]struct{} {
	out := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(strings.ToLower(value))
		if value == "" {
			continue
		}
		out[value] = struct{}{}
	}
	return out
}

func normalizeExtSet(values []string) map[string]struct{} {
	if len(values) == 0 {
		return nil
	}

	out := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.TrimSpace(strings.ToLower(value))
		if value == "" {
			continue
		}
		if !strings.HasPrefix(value, ".") {
			value = "." + value
		}
		out[value] = struct{}{}
	}
	return out
}

func matchesAnyGlob(rel string, globs []string) bool {
	if len(globs) == 0 {
		return false
	}
	for _, glob := range globs {
		match, err := doublestar.PathMatch(glob, rel)
		if err == nil && match {
			return true
		}
	}
	return false
}

func splitLines(content string) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func tokenizeLines(lines []string) ([]LineTerms, int) {
	lineTerms := make([]LineTerms, 0, len(lines))
	total := 0

	for idx, line := range lines {
		terms := tokenize(line)
		if len(terms) == 0 {
			continue
		}

		keys := make([]string, 0, len(terms))
		for term := range terms {
			keys = append(keys, term)
		}
		sort.Strings(keys)
		total += len(keys)

		hits := make([]TermHit, 0, len(keys))
		for _, key := range keys {
			hits = append(hits, TermHit{Term: key, Frequency: terms[key]})
		}

		lineTerms = append(lineTerms, LineTerms{
			LineNo: idx + 1,
			Hits:   hits,
		})
	}

	return lineTerms, total
}

func tokenize(line string) map[string]int {
	terms := make(map[string]int)
	var current []rune

	flush := func() {
		if len(current) == 0 {
			return
		}
		token := strings.ToLower(string(current))
		terms[token]++
		current = current[:0]
	}

	for _, r := range line {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			current = append(current, r)
			continue
		}
		flush()
	}
	flush()

	return terms
}

func readTextFile(path string, maxSize int64) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	if int64(len(data)) > maxSize {
		return "", false, nil
	}

	if len(data) >= 2 {
		if bytes.Equal(data[:2], []byte{0xFF, 0xFE}) {
			return decodeUTF16(data[2:], true), true, nil
		}
		if bytes.Equal(data[:2], []byte{0xFE, 0xFF}) {
			return decodeUTF16(data[2:], false), true, nil
		}
	}
	if len(data) >= 3 && bytes.Equal(data[:3], []byte{0xEF, 0xBB, 0xBF}) {
		data = data[3:]
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return "", false, nil
	}
	if !utf8.Valid(data) {
		return "", false, nil
	}

	return string(data), true, nil
}

func decodeUTF16(data []byte, littleEndian bool) string {
	u16 := make([]uint16, 0, len(data)/2)
	for i := 0; i+1 < len(data); i += 2 {
		var value uint16
		if littleEndian {
			value = uint16(data[i]) | uint16(data[i+1])<<8
		} else {
			value = uint16(data[i])<<8 | uint16(data[i+1])
		}
		u16 = append(u16, value)
	}
	return string(utf16.Decode(u16))
}

func normalizePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

func samePath(left, right string) bool {
	return filepath.Clean(left) == filepath.Clean(right)
}
