package app

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"seek/internal/indexer"
	appLogger "seek/internal/logger"
)

func runInteractive(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("interactive", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	root := fs.String("root", "", "Root used to resolve the default database path")
	dbPath := fs.String("db", "", "Path to the SQLite index file")
	mode := fs.String("mode", "search", "Starting mode: search or grep")
	typeArg := fs.String("type", "", "Comma-separated type filters")
	limit := fs.Int("limit", 20, "Maximum matches to print")
	contextLines := fs.Int("context", 0, "Show N lines of surrounding context")
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

	currentMode := strings.ToLower(strings.TrimSpace(*mode))
	if currentMode != "search" && currentMode != "grep" {
		currentMode = "search"
	}

	currentTypes := strings.TrimSpace(*typeArg)
	currentLimit := *limit
	currentContext := parseContext(*contextLines)

	appLogger.SearchHeader(stdout, "seek interactive", resolvedRoot)
	fmt.Fprintln(stdout, "Commands:")
	fmt.Fprintln(stdout, "  :search        switch to exact-word mode")
	fmt.Fprintln(stdout, "  :grep          switch to regex mode")
	fmt.Fprintln(stdout, "  :type go,ts    set type filter")
	fmt.Fprintln(stdout, "  :context 2     set context lines")
	fmt.Fprintln(stdout, "  :limit 25      set result limit")
	fmt.Fprintln(stdout, "  :root <path>   switch root")
	fmt.Fprintln(stdout, "  :help          show commands")
	fmt.Fprintln(stdout, "  :quit          exit")
	fmt.Fprintln(stdout)

	reader := bufio.NewScanner(os.Stdin)
	for {
		appLogger.Prompt(stdout, fmt.Sprintf("seek[%s]> ", currentMode))
		if !reader.Scan() {
			fmt.Fprintln(stdout)
			return nil
		}

		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, ":") {
			if done, err := handleInteractiveCommand(line, &resolvedRoot, &currentMode, &currentTypes, &currentLimit, &currentContext, stdout); err != nil {
				fmt.Fprintf(stderr, "%v\n", err)
			} else if done {
				return nil
			}
			continue
		}

		store, closeStore, err := openStoreForRead(resolvedRoot, *dbPath)
		if err != nil {
			fmt.Fprintf(stderr, "%v\n", err)
			continue
		}

		typeSet, _, err := resolveTypeFilters(currentTypes)
		if err != nil {
			closeStore()
			fmt.Fprintf(stderr, "%v\n", err)
			continue
		}

		switch currentMode {
		case "grep":
			results, err := store.Grep(indexer.GrepOptions{
				Pattern:   line,
				Limit:     candidateLimit(currentLimit),
				PathGlobs: nil,
			})
			if err == nil {
				results = filterAndRankResults(results, line, typeSet, currentLimit)
			}
			if err != nil {
				fmt.Fprintf(stderr, "%v\n", err)
				closeStore()
				continue
			}
			re, err := regexp.Compile(line)
			if err != nil {
				fmt.Fprintf(stderr, "%v\n", err)
				closeStore()
				continue
			}
			renderSearchResults(stdout, resolvedRoot, fmt.Sprintf("Pattern: %q", line), results, false, func(text string) string {
				return appLogger.HighlightRegex(text, re)
			}, store, renderOptions{ContextLines: currentContext})
		default:
			results, err := store.SearchTerm(line, nil, candidateLimit(currentLimit))
			if err == nil {
				results = filterAndRankResults(results, line, typeSet, currentLimit)
			}
			if err != nil {
				fmt.Fprintf(stderr, "%v\n", err)
				closeStore()
				continue
			}
			renderSearchResults(stdout, resolvedRoot, fmt.Sprintf("Exact word: %q", line), results, true, func(text string) string {
				return appLogger.Highlight(text, line)
			}, store, renderOptions{ContextLines: currentContext})
		}
		closeStore()
	}
}

func handleInteractiveCommand(command string, root *string, mode *string, typeArg *string, limit *int, contextLines *int, stdout io.Writer) (bool, error) {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false, nil
	}

	switch fields[0] {
	case ":quit", ":q", ":exit":
		fmt.Fprintln(stdout, "Bye.")
		return true, nil
	case ":search":
		*mode = "search"
		fmt.Fprintln(stdout, "Mode set to search")
	case ":grep":
		*mode = "grep"
		fmt.Fprintln(stdout, "Mode set to grep")
	case ":type":
		if len(fields) < 2 {
			*typeArg = ""
			fmt.Fprintln(stdout, "Type filter cleared")
			return false, nil
		}
		value := strings.ReplaceAll(strings.Join(fields[1:], ""), " ", "")
		_, _, err := resolveTypeFilters(value)
		if err != nil {
			return false, err
		}
		*typeArg = value
		fmt.Fprintf(stdout, "Type filter set to %s\n", *typeArg)
	case ":context":
		if len(fields) < 2 {
			return false, fmt.Errorf("usage: :context <number>")
		}
		*contextLines = parseContext(parseIntArg(fields[1], *contextLines))
		fmt.Fprintf(stdout, "Context set to %d\n", *contextLines)
	case ":limit":
		if len(fields) < 2 {
			return false, fmt.Errorf("usage: :limit <number>")
		}
		*limit = parseIntArg(fields[1], *limit)
		fmt.Fprintf(stdout, "Limit set to %d\n", *limit)
	case ":root":
		if len(fields) < 2 {
			return false, fmt.Errorf("usage: :root <path>")
		}
		resolvedRoot, err := resolveExplicitRoot(strings.Join(fields[1:], " "))
		if err != nil {
			return false, err
		}
		*root = resolvedRoot
		fmt.Fprintf(stdout, "Root set to %s\n", *root)
	case ":help":
		fmt.Fprintln(stdout, "Commands: :search :grep :type <csv> :context <n> :limit <n> :root <path> :quit")
	default:
		return false, fmt.Errorf("unknown command %q", fields[0])
	}

	return false, nil
}
