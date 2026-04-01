package app

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"seek/internal/config"
	"seek/internal/indexer"
	appLogger "seek/internal/logger"
)

func Run(args []string, stdout, stderr io.Writer) error {
	cfg := appLogger.NoFileConfig()
	cfg.Level = "info"
	cfg.Format = "text"
	cfg.Output = "console"
	cfg.Verbose = true
	_ = appLogger.Initialize(cfg)

	if len(args) == 0 {
		printUsage(stdout)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		return nil
	case "config":
		return runConfig(args[1:], stdout)
	case "watch":
		return runWatch(args[1:], stdout, stderr)
	case "interactive":
		return runInteractive(args[1:], stdout, stderr)
	case "build":
		return runBuild(args[1:], stdout)
	case "search":
		return runSearch(args[1:], stdout)
	case "grep":
		return runGrep(args[1:], stdout)
	case "stats":
		return runStats(args[1:], stdout)
	default:
		fmt.Fprintf(stderr, "Unknown command %q\n\n", args[0])
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runBuild(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Target project root to index")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
	typeArg := fs.String("type", "", "Comma-separated type filters (go,ts,js,proto,md,yaml,sql,py,ps,html,docs,config,web)")
	includeExt := fs.String("include-ext", "", "Comma-separated file extensions to include")
	excludeDirs := fs.String("exclude-dir", "", "Comma-separated directory names to skip")
	excludeGlobs := fs.String("exclude-glob", "", "Comma-separated glob patterns to skip")
	maxFileSizeMB := fs.Int64("max-file-size-mb", 4, "Maximum file size to index in MB")
	includeHidden := fs.Bool("include-hidden", false, "Include hidden files and directories")
	followSymlinks := fs.Bool("follow-symlinks", false, "Follow symlinked files and directories")
	reset := fs.Bool("reset", false, "Reset the existing index before rebuilding")
	verbose := fs.Bool("verbose", false, "Print skipped files and progress details")
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
		Reset:          *reset,
		Verbose:        *verbose,
		Progress:       logBuildProgress,
	}

	startedAt := time.Now()
	appLogger.Infof("seek build started: root=%s", resolvedRoot)
	stats, err := indexer.Build(opts)
	if err != nil {
		appLogger.Errorf("seek build failed: %v", err)
		return err
	}
	appLogger.Infof("seek build completed: indexed=%d unchanged=%d skipped=%d errors=%d elapsed=%s", stats.IndexedFiles, stats.UnchangedFiles, stats.SkippedFiles, stats.Errors, time.Since(startedAt).Round(time.Millisecond))

	fmt.Fprintf(stdout, "Indexed root: %s\n", stats.Root)
	fmt.Fprintf(stdout, "Index database: %s\n", stats.DBPath)
	fmt.Fprintf(stdout, "Files indexed: %d\n", stats.IndexedFiles)
	fmt.Fprintf(stdout, "Files updated: %d\n", stats.UpdatedFiles)
	fmt.Fprintf(stdout, "Files unchanged: %d\n", stats.UnchangedFiles)
	fmt.Fprintf(stdout, "Files removed: %d\n", stats.RemovedFiles)
	fmt.Fprintf(stdout, "Files skipped: %d\n", stats.SkippedFiles)
	fmt.Fprintf(stdout, "Errors: %d\n", stats.Errors)
	fmt.Fprintf(stdout, "Lines stored: %d\n", stats.IndexedLines)
	fmt.Fprintf(stdout, "Term hits stored: %d\n", stats.IndexedTerms)
	fmt.Fprintf(stdout, "Elapsed: %s\n", time.Since(startedAt).Round(time.Millisecond))

	return nil
}

func runSearch(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Root used to resolve the default database path")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
	term := fs.String("term", "", "Word to search for")
	limit := fs.Int("limit", 50, "Maximum number of matches to print")
	contextLines := fs.Int("context", 0, "Show N lines of surrounding context")
	pathFilter := fs.String("path", "", "Comma-separated glob filters for file paths")
	typeArg := fs.String("type", "", "Comma-separated type filters (go,ts,js,proto,md,yaml,sql,py,ps,html,docs,config,web)")
	saveRoot := fs.Bool("save-root", false, "Save the resolved root as the default root for future commands")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*term) == "" {
		return fmt.Errorf("--term is required")
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

	store, closeStore, err := openStoreForRead(resolvedRoot, *dbPath)
	if err != nil {
		return err
	}
	defer closeStore()

	typeSet, _, err := resolveTypeFilters(*typeArg)
	if err != nil {
		return err
	}

	needle := strings.TrimSpace(*term)
	results, err := store.SearchTerm(needle, splitCSV(*pathFilter), candidateLimit(*limit))
	if err != nil {
		return err
	}
	results = filterAndRankResults(results, needle, typeSet, *limit)

	if len(results) == 0 {
		fmt.Fprintf(stdout, "No exact-word matches found for %q\n", *term)
		return nil
	}

	renderSearchResults(stdout, resolvedRoot, fmt.Sprintf("Exact word: %q", needle), results, true, func(line string) string {
		return appLogger.Highlight(line, needle)
	}, store, renderOptions{ContextLines: parseContext(*contextLines)})

	return nil
}

func runGrep(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Root used to resolve the default database path")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
	pattern := fs.String("pattern", "", "Regex or fixed pattern to search for")
	limit := fs.Int("limit", 50, "Maximum number of matches to print")
	contextLines := fs.Int("context", 0, "Show N lines of surrounding context")
	pathFilter := fs.String("path", "", "Comma-separated glob filters for file paths")
	typeArg := fs.String("type", "", "Comma-separated type filters (go,ts,js,proto,md,yaml,sql,py,ps,html,docs,config,web)")
	fixed := fs.Bool("fixed", false, "Treat pattern as a fixed string")
	ignoreCase := fs.Bool("ignore-case", false, "Ignore case during matching")
	saveRoot := fs.Bool("save-root", false, "Save the resolved root as the default root for future commands")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*pattern) == "" {
		return fmt.Errorf("--pattern is required")
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

	store, closeStore, err := openStoreForRead(resolvedRoot, *dbPath)
	if err != nil {
		return err
	}
	defer closeStore()

	typeSet, _, err := resolveTypeFilters(*typeArg)
	if err != nil {
		return err
	}

	results, err := store.Grep(indexer.GrepOptions{
		Pattern:    *pattern,
		PathGlobs:  splitCSV(*pathFilter),
		Limit:      candidateLimit(*limit),
		Fixed:      *fixed,
		IgnoreCase: *ignoreCase,
	})
	if err != nil {
		return err
	}
	results = filterAndRankResults(results, strings.TrimSpace(*pattern), typeSet, *limit)

	if len(results) == 0 {
		fmt.Fprintf(stdout, "No pattern matches found for %q\n", *pattern)
		return nil
	}

	highlight := ""
	var highlightFunc func(string) string
	if *fixed {
		highlight = strings.TrimSpace(*pattern)
		highlightFunc = func(line string) string {
			return appLogger.Highlight(line, highlight)
		}
	} else {
		regexPattern := strings.TrimSpace(*pattern)
		if *ignoreCase {
			regexPattern = "(?i)" + regexPattern
		}
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return err
		}
		highlightFunc = func(line string) string {
			return appLogger.HighlightRegex(line, re)
		}
	}
	renderSearchResults(stdout, resolvedRoot, fmt.Sprintf("Pattern: %q", strings.TrimSpace(*pattern)), results, false, highlightFunc, store, renderOptions{ContextLines: parseContext(*contextLines)})

	return nil
}

func runStats(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Root used to resolve the default database path")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
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

	store, closeStore, err := openStoreForRead(resolvedRoot, *dbPath)
	if err != nil {
		return err
	}
	defer closeStore()

	stats, err := store.Stats()
	if err != nil {
		return err
	}

	fileInfo, err := os.Stat(store.Path())
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Indexed root: %s\n", stats.Root)
	fmt.Fprintf(stdout, "Index database: %s\n", store.Path())
	fmt.Fprintf(stdout, "Database size: %s\n", humanBytes(fileInfo.Size()))
	fmt.Fprintf(stdout, "Indexed files: %d\n", stats.Files)
	fmt.Fprintf(stdout, "Indexed lines: %d\n", stats.Lines)
	fmt.Fprintf(stdout, "Distinct terms: %d\n", stats.DistinctTerms)
	fmt.Fprintf(stdout, "Term hits: %d\n", stats.TermHits)
	fmt.Fprintf(stdout, "Last build: %s\n", stats.LastBuild)

	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "seek indexes a project into a portable SQLite database for fast exact-word and pattern search.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  seek config [show|set-root|clear-root]")
	fmt.Fprintln(w, "  seek build  [flags]")
	fmt.Fprintln(w, "  seek search [flags]")
	fmt.Fprintln(w, "  seek grep   [flags]")
	fmt.Fprintln(w, "  seek watch  [flags]")
	fmt.Fprintln(w, "  seek interactive [flags]")
	fmt.Fprintln(w, "  seek stats  [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  seek config set-root E:\\Projects\\InsureTech")
	fmt.Fprintln(w, "  seek build --root E:\\Projects\\InsureTech")
	fmt.Fprintln(w, "  seek build --type go,ts --save-root")
	fmt.Fprintln(w, "  seek build --root E:\\Projects\\InsureTech --save-root")
	fmt.Fprintln(w, "  seek search --term policy --context 2")
	fmt.Fprintln(w, "  seek search --root E:\\Projects\\InsureTech --term policy --type go,ts")
	fmt.Fprintln(w, "  seek grep --root E:\\Projects\\InsureTech --pattern \"Claim.*Status\" --context 2")
	fmt.Fprintln(w, "  seek grep --root . --pattern payment --fixed --ignore-case --path \"backend/**,api/**\"")
	fmt.Fprintln(w, "  seek watch --root E:\\Projects\\InsureTech")
	fmt.Fprintln(w, "  seek interactive --root E:\\Projects\\InsureTech")
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		out = append(out, part)
	}
	return out
}

func openStoreForRead(root, dbPath string) (*indexer.Store, func(), error) {
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, nil, err
	}

	resolvedDB, err := indexer.ResolveDBPath(resolvedRoot, dbPath)
	if err != nil {
		return nil, nil, err
	}
	if _, err := os.Stat(resolvedDB); err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("no index found at %s; run `seek build --root %s` first", resolvedDB, resolvedRoot)
		}
		return nil, nil, err
	}

	store, err := indexer.OpenStore(resolvedDB)
	if err != nil {
		return nil, nil, err
	}
	if err := store.Init(); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	stats, err := store.Stats()
	if err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	if strings.TrimSpace(stats.Root) == "" {
		_ = store.Close()
		return nil, nil, fmt.Errorf("no index found at %s; run `seek build --root %s` first", resolvedDB, resolvedRoot)
	}
	if err := store.AssertRoot(resolvedRoot, false); err != nil {
		_ = store.Close()
		return nil, nil, err
	}

	return store, func() {
		_ = store.Close()
	}, nil
}

func humanBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func logBuildProgress(progress indexer.BuildProgress) {
	switch progress.Stage {
	case "start":
		appLogger.Infof("scanning files under %s", progress.Path)
	case "done":
		appLogger.Infof("scan finished: scanned=%d indexed=%d updated=%d unchanged=%d skipped=%d errors=%d", progress.ScannedFiles, progress.IndexedFiles, progress.UpdatedFiles, progress.UnchangedFiles, progress.SkippedFiles, progress.Errors)
	default:
		appLogger.Infof("progress: scanned=%d indexed=%d updated=%d unchanged=%d skipped=%d errors=%d current=%s", progress.ScannedFiles, progress.IndexedFiles, progress.UpdatedFiles, progress.UnchangedFiles, progress.SkippedFiles, progress.Errors, progress.Path)
	}
}

func renderSearchResults(stdout io.Writer, root, title string, results []indexer.SearchResult, includeFrequency bool, highlight func(string) string, store *indexer.Store, options renderOptions) {
	appLogger.SearchHeader(stdout, title, root)
	currentPath := ""
	for _, result := range results {
		absolutePath := filepath.Join(root, filepath.FromSlash(result.Path))
		if absolutePath != currentPath {
			if currentPath != "" {
				fmt.Fprintln(stdout)
			}
			currentPath = absolutePath
			appLogger.SearchFile(stdout, absolutePath)
		}

		if options.ContextLines > 0 && store != nil {
			contextLines, err := store.GetContext(result.Path, result.LineNo, options.ContextLines, options.ContextLines)
			if err == nil && len(contextLines) > 0 {
				for _, contextLine := range contextLines {
					lineText := renderSnippet(contextLine.Line)
					if contextLine.LineNo == result.LineNo {
						if highlight != nil {
							lineText = highlight(lineText)
						}
						appLogger.SearchHit(stdout, contextLine.LineNo, result.Frequency, lineText, includeFrequency, "")
					} else {
						appLogger.SearchContext(stdout, contextLine.LineNo, lineText)
					}
				}
				fmt.Fprintln(stdout)
				continue
			}
		}

		lineText := renderSnippet(result.Line)
		if highlight != nil {
			lineText = highlight(lineText)
		}
		appLogger.SearchHit(stdout, result.LineNo, result.Frequency, lineText, includeFrequency, "")
	}

	appLogger.SearchSummary(stdout, len(results))
}

func renderSnippet(line string) string {
	line = strings.ReplaceAll(line, "\t", "    ")
	line = strings.TrimSpace(line)
	if line == "" {
		return "<empty line>"
	}

	const maxLen = 220
	runes := []rune(line)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + " ..."
	}

	return line
}

func runConfig(args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] == "show" {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		path, err := config.Path()
		if err != nil {
			return err
		}
		if strings.TrimSpace(cfg.DefaultRoot) == "" {
			fmt.Fprintf(stdout, "Config file: %s\n", path)
			fmt.Fprintln(stdout, "Default root: not set")
			return nil
		}
		fmt.Fprintf(stdout, "Config file: %s\n", path)
		fmt.Fprintf(stdout, "Default root: %s\n", cfg.DefaultRoot)
		if cfg.UpdatedAt != "" {
			fmt.Fprintf(stdout, "Updated at: %s\n", cfg.UpdatedAt)
		}
		return nil
	}

	switch args[0] {
	case "set-root":
		if len(args) < 2 {
			return fmt.Errorf("usage: seek config set-root <path>")
		}
		resolvedRoot, err := resolveExplicitRoot(args[1])
		if err != nil {
			return err
		}
		if err := saveDefaultRoot(resolvedRoot); err != nil {
			return err
		}
		fmt.Fprintf(stdout, "Default root saved: %s\n", resolvedRoot)
		return nil
	case "clear-root":
		if err := config.Clear(); err != nil {
			return err
		}
		fmt.Fprintln(stdout, "Default root cleared")
		return nil
	default:
		return fmt.Errorf("unknown config command %q", args[0])
	}
}

func resolveRoot(root string) (string, error) {
	if strings.TrimSpace(root) != "" {
		return resolveExplicitRoot(root)
	}

	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(cfg.DefaultRoot) != "" {
		resolvedRoot, err := resolveExplicitRoot(cfg.DefaultRoot)
		if err != nil {
			return "", fmt.Errorf("default root is invalid: %w", err)
		}
		if filepath.Clean(cfg.DefaultRoot) != resolvedRoot {
			if err := config.Save(config.Config{DefaultRoot: resolvedRoot}); err != nil {
				return "", err
			}
		}
		return resolvedRoot, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Abs(cwd)
}

func resolveExplicitRoot(root string) (string, error) {
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(resolvedRoot)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("root path is not a directory: %s", resolvedRoot)
	}
	return resolvedRoot, nil
}

func saveDefaultRoot(root string) error {
	resolvedRoot, err := resolveExplicitRoot(root)
	if err != nil {
		return err
	}
	if err := config.Save(config.Config{DefaultRoot: resolvedRoot}); err != nil {
		return err
	}
	return nil
}
