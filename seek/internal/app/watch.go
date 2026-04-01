package app

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"

	"seek/internal/indexer"
	appLogger "seek/internal/logger"
)

var defaultWatchExcludedDirs = map[string]struct{}{
	".cache":       {},
	".git":         {},
	".hg":          {},
	".idea":        {},
	".next":        {},
	".nuxt":        {},
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

type watchTree struct {
	root          string
	watcher       *fsnotify.Watcher
	watched       map[string]struct{}
	includeHidden bool
	excludeDirs   map[string]struct{}
	excludeGlobs  []string
}

func runWatch(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Target project root to index")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
	typeArg := fs.String("type", "", "Comma-separated type filters for indexed files")
	includeExt := fs.String("include-ext", "", "Comma-separated file extensions to include")
	excludeDirs := fs.String("exclude-dir", "", "Comma-separated directory names to skip")
	excludeGlobs := fs.String("exclude-glob", "", "Comma-separated glob patterns to skip")
	maxFileSizeMB := fs.Int64("max-file-size-mb", 4, "Maximum file size to index in MB")
	includeHidden := fs.Bool("include-hidden", false, "Include hidden files and directories")
	followSymlinks := fs.Bool("follow-symlinks", false, "Follow symlinked files and directories")
	interval := fs.Int("interval", 2, "Root fallback check interval in seconds")
	saveRoot := fs.Bool("save-root", false, "Save the resolved root as the default root for future commands")

	if err := fs.Parse(args); err != nil {
		return err
	}

	resolvedRoot, err := resolveRoot(*root)
	if err != nil {
		return err
	}
	if *saveRoot {
		if err := saveDefaultRoot(resolvedRoot); err != nil {
			return err
		}
	}

	includeExts, err := mergeIncludeExtensions(*includeExt, *typeArg)
	if err != nil {
		return err
	}
	allExcludeGlobs, err := buildExcludeGlobs(resolvedRoot, splitCSV(*excludeGlobs))
	if err != nil {
		return err
	}

	opts := indexer.BuildOptions{
		Root:           resolvedRoot,
		DBPath:         *dbPath,
		IncludeExt:     includeExts,
		ExcludeDirs:    splitCSV(*excludeDirs),
		ExcludeGlobs:   allExcludeGlobs,
		MaxFileSize:    *maxFileSizeMB * 1024 * 1024,
		IncludeHidden:  *includeHidden,
		FollowSymlinks: *followSymlinks,
		Progress:       logBuildProgress,
	}

	tree, err := newWatchTree(resolvedRoot, *includeHidden, opts.ExcludeDirs, allExcludeGlobs)
	if err != nil {
		return err
	}
	defer tree.Close()

	if err := tree.Sync(); err != nil {
		return err
	}

	fallbackInterval := parseIntervalSeconds(*interval)
	appLogger.Infof("seek watch started: root=%s mode=fs-events fallback=%ds", resolvedRoot, fallbackInterval)
	fmt.Fprintf(stdout, "Watching %s for filesystem events. Root fallback check every %ds. Press Ctrl+C to stop.\n", resolvedRoot, fallbackInterval)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	rootTicker := time.NewTicker(time.Duration(fallbackInterval) * time.Second)
	defer rootTicker.Stop()

	var (
		debounceTimer *time.Timer
		debounceC     <-chan time.Time
		lastReason    string
	)

	startDebounce := func(reason string) {
		lastReason = reason
		if debounceTimer == nil {
			debounceTimer = time.NewTimer(450 * time.Millisecond)
			debounceC = debounceTimer.C
			return
		}
		if !debounceTimer.Stop() {
			select {
			case <-debounceTimer.C:
			default:
			}
		}
		debounceTimer.Reset(450 * time.Millisecond)
		debounceC = debounceTimer.C
	}

	runBuildCycle := func(reason string) error {
		if err := tree.Sync(); err != nil {
			appLogger.Warnf("watch sync warning: %v", err)
		}

		stats, err := indexer.Build(opts)
		if err != nil {
			return err
		}
		if stats.IndexedFiles > 0 || stats.UpdatedFiles > 0 || stats.RemovedFiles > 0 {
			fmt.Fprintf(stdout, "[%s] %s -> indexed=%d updated=%d removed=%d unchanged=%d skipped=%d errors=%d\n",
				time.Now().Format("15:04:05"),
				reason,
				stats.IndexedFiles,
				stats.UpdatedFiles,
				stats.RemovedFiles,
				stats.UnchangedFiles,
				stats.SkippedFiles,
				stats.Errors,
			)
		}
		return nil
	}

	if err := runBuildCycle("initial build"); err != nil {
		return err
	}

	rootSignature, err := tree.RootSignature()
	if err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-tree.watcher.Events:
			if !ok {
				return nil
			}
			reason, interesting, err := tree.HandleEvent(event)
			if err != nil {
				fmt.Fprintf(stderr, "watch event handling failed: %v\n", err)
			}
			if interesting {
				startDebounce(reason)
			}
		case err, ok := <-tree.watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(stderr, "watcher error: %v\n", err)
		case <-debounceC:
			debounceC = nil
			if err := runBuildCycle(lastReason); err != nil {
				fmt.Fprintf(stderr, "watch rebuild failed: %v\n", err)
			}
			rootSignature, _ = tree.RootSignature()
		case <-rootTicker.C:
			currentSignature, err := tree.RootSignature()
			if err != nil {
				fmt.Fprintf(stderr, "watch root check failed: %v\n", err)
				continue
			}
			if currentSignature != rootSignature {
				rootSignature = currentSignature
				if err := tree.Sync(); err != nil {
					fmt.Fprintf(stderr, "watch sync failed: %v\n", err)
				}
				startDebounce("root activity")
			}
		case <-stop:
			fmt.Fprintln(stdout, "Watch stopped.")
			return nil
		}
	}
}

func newWatchTree(root string, includeHidden bool, excludeDirs []string, excludeGlobs []string) (*watchTree, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &watchTree{
		root:          filepath.Clean(root),
		watcher:       watcher,
		watched:       make(map[string]struct{}),
		includeHidden: includeHidden,
		excludeDirs:   normalizeWatchDirSet(excludeDirs),
		excludeGlobs:  excludeGlobs,
	}, nil
}

func (w *watchTree) Close() error {
	return w.watcher.Close()
}

func (w *watchTree) Sync() error {
	w.removeMissing()
	return w.addRecursive(w.root)
}

func (w *watchTree) HandleEvent(event fsnotify.Event) (string, bool, error) {
	path := filepath.Clean(event.Name)
	if !isInterestingWatchOp(event.Op) {
		return "", false, nil
	}
	if !w.shouldHandlePath(path) {
		return "", false, nil
	}

	if event.Op&(fsnotify.Create|fsnotify.Rename) != 0 {
		if err := w.addIfDirectory(path); err != nil {
			return path, true, err
		}
	}
	if event.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
		w.forgetPath(path)
	}

	return describeWatchEvent(event), true, nil
}

func (w *watchTree) RootSignature() (string, error) {
	entries, err := os.ReadDir(w.root)
	if err != nil {
		return "", err
	}

	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(w.root, entry.Name())
		if !w.shouldHandlePath(path) {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%t:%d", strings.ToLower(entry.Name()), entry.IsDir(), info.ModTime().UnixNano()))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|"), nil
}

func (w *watchTree) addRecursive(start string) error {
	return filepath.WalkDir(start, func(currentPath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if !entry.IsDir() {
			return nil
		}
		if !w.shouldWatchDir(currentPath) {
			if sameFilePath(currentPath, start) {
				return nil
			}
			return filepath.SkipDir
		}
		return w.addDir(currentPath)
	})
}

func (w *watchTree) addIfDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return nil
	}
	if !w.shouldWatchDir(path) {
		return nil
	}
	return w.addRecursive(path)
}

func (w *watchTree) addDir(path string) error {
	path = filepath.Clean(path)
	if _, ok := w.watched[path]; ok {
		return nil
	}
	if err := w.watcher.Add(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	w.watched[path] = struct{}{}
	return nil
}

func (w *watchTree) forgetPath(path string) {
	path = filepath.Clean(path)
	for watchedPath := range w.watched {
		if sameFilePath(watchedPath, path) || isChildPath(path, watchedPath) {
			_ = w.watcher.Remove(watchedPath)
			delete(w.watched, watchedPath)
		}
	}
}

func (w *watchTree) removeMissing() {
	for watchedPath := range w.watched {
		if _, err := os.Stat(watchedPath); err == nil {
			continue
		}
		_ = w.watcher.Remove(watchedPath)
		delete(w.watched, watchedPath)
	}
}

func (w *watchTree) shouldWatchDir(path string) bool {
	if sameFilePath(path, w.root) {
		return true
	}
	return w.shouldHandlePath(path)
}

func (w *watchTree) shouldHandlePath(path string) bool {
	rel, ok := relativeWatchPath(w.root, path)
	if !ok {
		return false
	}
	if rel == "." {
		return true
	}

	segments := strings.Split(rel, "/")
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		lower := strings.ToLower(segment)
		if !w.includeHidden && strings.HasPrefix(segment, ".") {
			return false
		}
		if _, found := defaultWatchExcludedDirs[lower]; found {
			return false
		}
		if _, found := w.excludeDirs[lower]; found {
			return false
		}
	}

	return !matchesWatchGlobs(rel, w.excludeGlobs)
}

func normalizeWatchDirSet(values []string) map[string]struct{} {
	out := make(map[string]struct{}, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		out[value] = struct{}{}
	}
	return out
}

func relativeWatchPath(root, path string) (string, bool) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", false
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", false
	}
	rel = filepath.ToSlash(filepath.Clean(rel))
	if rel == "." {
		return ".", true
	}
	return rel, true
}

func matchesWatchGlobs(rel string, globs []string) bool {
	for _, glob := range globs {
		match, err := doublestar.PathMatch(glob, rel)
		if err == nil && match {
			return true
		}
	}
	return false
}

func isInterestingWatchOp(op fsnotify.Op) bool {
	return op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) != 0
}

func describeWatchEvent(event fsnotify.Event) string {
	path := filepath.Clean(event.Name)
	switch {
	case event.Op&fsnotify.Create != 0:
		return fmt.Sprintf("created %s", path)
	case event.Op&fsnotify.Write != 0:
		return fmt.Sprintf("updated %s", path)
	case event.Op&fsnotify.Remove != 0:
		return fmt.Sprintf("removed %s", path)
	case event.Op&fsnotify.Rename != 0:
		return fmt.Sprintf("renamed %s", path)
	default:
		return fmt.Sprintf("changed %s", path)
	}
}

func sameFilePath(left, right string) bool {
	return strings.EqualFold(filepath.Clean(left), filepath.Clean(right))
}

func isChildPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
