package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"seek/internal/indexer"
)

var typeExtensions = map[string][]string{
	"csharp": {".cs"},
	"config": {".conf", ".env", ".ini", ".json", ".toml", ".xml", ".yaml", ".yml"},
	"css":    {".css", ".less", ".sass", ".scss"},
	"docs":   {".md", ".mdx", ".rst", ".txt"},
	"go":     {".go"},
	"html":   {".htm", ".html", ".svelte"},
	"js":     {".cjs", ".js", ".jsx", ".mjs"},
	"json":   {".json"},
	"md":     {".md", ".mdx"},
	"proto":  {".proto"},
	"ps":     {".ps1", ".psd1", ".psm1"},
	"py":     {".py"},
	"sh":     {".bash", ".sh", ".zsh"},
	"sql":    {".sql"},
	"ts":     {".ts", ".tsx"},
	"web":    {".css", ".htm", ".html", ".js", ".jsx", ".svelte", ".ts", ".tsx"},
	"yaml":   {".yaml", ".yml"},
}

type renderOptions struct {
	ContextLines int
}

func buildExcludeGlobs(root string, cliGlobs []string) ([]string, error) {
	filePatterns, err := loadSeekIgnore(root)
	if err != nil {
		return nil, err
	}
	globs := make([]string, 0, len(cliGlobs)+len(filePatterns))
	globs = append(globs, cliGlobs...)
	globs = append(globs, filePatterns...)
	return globs, nil
}

func loadSeekIgnore(root string) ([]string, error) {
	path := filepath.Join(root, ".seekignore")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	var patterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			// Re-include rules are not supported yet; skip quietly.
			continue
		}
		line = filepath.ToSlash(line)
		line = strings.TrimPrefix(line, "./")
		line = strings.TrimPrefix(line, "/")
		if line == "" {
			continue
		}

		switch {
		case strings.HasSuffix(line, "/"):
			line = strings.TrimSuffix(line, "/")
			patterns = append(patterns, line+"/**", "**/"+line+"/**")
		case strings.Contains(line, "/"):
			patterns = append(patterns, line)
		default:
			patterns = append(patterns, "**/"+line, "**/"+line+"/**")
		}
	}

	return patterns, nil
}

func resolveTypeFilters(typeArg string) (map[string]struct{}, []string, error) {
	typeNames := splitCSV(typeArg)
	if len(typeNames) == 0 {
		return nil, nil, nil
	}

	extSet := make(map[string]struct{})
	extList := make([]string, 0, len(typeNames))
	for _, name := range typeNames {
		key := strings.ToLower(strings.TrimSpace(name))
		extensions, ok := typeExtensions[key]
		if !ok {
			return nil, nil, fmt.Errorf("unknown type %q", name)
		}
		for _, ext := range extensions {
			if _, exists := extSet[ext]; exists {
				continue
			}
			extSet[ext] = struct{}{}
			extList = append(extList, ext)
		}
	}

	sort.Strings(extList)
	return extSet, extList, nil
}

func mergeIncludeExtensions(rawIncludeExt string, typeArg string) ([]string, error) {
	fromFlags := splitCSV(rawIncludeExt)
	typeSet, typeExts, err := resolveTypeFilters(typeArg)
	if err != nil {
		return nil, err
	}
	if len(fromFlags) == 0 && len(typeExts) == 0 {
		return nil, nil
	}

	merged := make(map[string]struct{})
	for _, ext := range fromFlags {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		merged[ext] = struct{}{}
	}
	for ext := range typeSet {
		merged[ext] = struct{}{}
	}

	out := make([]string, 0, len(merged))
	for ext := range merged {
		out = append(out, ext)
	}
	sort.Strings(out)
	return out, nil
}

func filterAndRankResults(results []indexer.SearchResult, query string, typeSet map[string]struct{}, limit int) []indexer.SearchResult {
	if len(results) == 0 {
		return results
	}

	candidates := make([]scoredResult, 0, len(results))
	for _, result := range results {
		if !matchesTypeFilter(result.Path, typeSet) {
			continue
		}
		candidates = append(candidates, scoredResult{
			result: result,
			score:  scoreResult(result, query),
		})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].score == candidates[j].score {
			if candidates[i].result.Path == candidates[j].result.Path {
				return candidates[i].result.LineNo < candidates[j].result.LineNo
			}
			return candidates[i].result.Path < candidates[j].result.Path
		}
		return candidates[i].score > candidates[j].score
	})

	if limit > 0 && len(candidates) > limit {
		candidates = candidates[:limit]
	}

	out := make([]indexer.SearchResult, 0, len(candidates))
	for _, candidate := range candidates {
		out = append(out, candidate.result)
	}
	return out
}

func candidateLimit(limit int) int {
	if limit <= 0 {
		return 300
	}
	if limit < 50 {
		return 300
	}
	return limit * 8
}

func parseContext(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func parseIntervalSeconds(value int) int {
	if value < 1 {
		return 2
	}
	return value
}

func matchesTypeFilter(path string, typeSet map[string]struct{}) bool {
	if len(typeSet) == 0 {
		return true
	}
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := typeSet[ext]
	return ok
}

type scoredResult struct {
	result indexer.SearchResult
	score  float64
}

func scoreResult(result indexer.SearchResult, query string) float64 {
	path := strings.ToLower(filepath.ToSlash(result.Path))
	line := strings.ToLower(result.Line)
	query = strings.ToLower(strings.TrimSpace(query))
	base := strings.ToLower(filepath.Base(path))
	ext := strings.ToLower(filepath.Ext(path))

	score := extensionWeight(ext)
	score += pathWeight(path)

	if query != "" {
		if strings.Contains(base, query) {
			score += 18
		}
		if strings.Contains(path, query) {
			score += 8
		}
		if strings.Contains(line, query) {
			score += 6
		}
	}
	if result.Frequency > 0 {
		score += float64(result.Frequency) * 1.5
	}
	if result.LineNo > 0 {
		score += 4 / float64(result.LineNo+4)
	}
	score -= float64(strings.Count(path, "/")) * 0.15

	return score
}

func extensionWeight(ext string) float64 {
	switch ext {
	case ".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".proto", ".sql", ".ps1", ".sh", ".yaml", ".yml":
		return 30
	case ".json", ".toml", ".ini", ".conf":
		return 20
	case ".md", ".txt", ".rst":
		return 10
	case ".html", ".htm", ".svelte":
		return 5
	default:
		return 0
	}
}

func pathWeight(path string) float64 {
	score := 0.0
	if strings.Contains(path, "/cmd/") || strings.Contains(path, "/src/") || strings.Contains(path, "/internal/") || strings.Contains(path, "/backend/") {
		score += 18
	}
	if strings.Contains(path, "/gen/") || strings.Contains(path, "proto-generated") || strings.Contains(path, ".pb.") {
		score -= 12
	}
	if strings.Contains(path, "/docs/") || strings.Contains(path, "/documentation/") || strings.Contains(path, "/api/docs/") {
		score -= 8
	}
	return score
}

func parseIntArg(value string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return n
}
