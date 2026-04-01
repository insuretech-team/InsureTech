package indexer

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	_ "modernc.org/sqlite"
)

type Store struct {
	db   *sql.DB
	path string
}

type FileState struct {
	ID          int64
	Size        int64
	ModTimeUnix int64
}

type FileRecord struct {
	Size        int64
	ModTimeUnix int64
	IndexedAt   int64
	SHA256      string
	LineCount   int
}

type TermHit struct {
	Term      string
	Frequency int
}

type LineTerms struct {
	LineNo int
	Hits   []TermHit
}

type SearchResult struct {
	Path      string
	LineNo    int
	Line      string
	Frequency int
}

type ContextLine struct {
	LineNo int
	Line   string
}

type IndexStats struct {
	Root          string
	Files         int
	Lines         int
	DistinctTerms int
	TermHits      int
	LastBuild     string
}

func OpenStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &Store{db: db, path: path}
	if err := store.configure(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) configure() error {
	pragmas := []string{
		"PRAGMA busy_timeout=5000;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA temp_store=MEMORY;",
	}
	for _, pragma := range pragmas {
		if _, err := s.db.Exec(pragma); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Init() error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS meta (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL UNIQUE,
			size INTEGER NOT NULL,
			mod_time_unix INTEGER NOT NULL,
			indexed_at_unix INTEGER NOT NULL,
			sha256 TEXT NOT NULL,
			line_count INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS lines (
			file_id INTEGER NOT NULL,
			line_no INTEGER NOT NULL,
			content TEXT NOT NULL,
			PRIMARY KEY (file_id, line_no),
			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS terms (
			term TEXT NOT NULL,
			file_id INTEGER NOT NULL,
			line_no INTEGER NOT NULL,
			frequency INTEGER NOT NULL,
			PRIMARY KEY (term, file_id, line_no),
			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_files_path ON files(path);`,
		`CREATE INDEX IF NOT EXISTS idx_terms_term ON terms(term);`,
		`CREATE INDEX IF NOT EXISTS idx_terms_file ON terms(file_id);`,
		`CREATE INDEX IF NOT EXISTS idx_lines_file ON lines(file_id);`,
	}

	for _, stmt := range schema {
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Reset() error {
	statements := []string{
		`DROP TABLE IF EXISTS terms;`,
		`DROP TABLE IF EXISTS lines;`,
		`DROP TABLE IF EXISTS files;`,
		`DROP TABLE IF EXISTS meta;`,
	}
	for _, stmt := range statements {
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}
	return s.Init()
}

func (s *Store) AssertRoot(root string, allowOverride bool) error {
	currentRoot, err := s.meta("root")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if currentRoot != "" && currentRoot != root && !allowOverride {
		return fmt.Errorf("index belongs to %q, not %q; use --reset or choose a different --db", currentRoot, root)
	}
	return nil
}

func (s *Store) SetBuildMeta(root string, builtAt time.Time) error {
	if err := s.setMeta("root", root); err != nil {
		return err
	}
	if err := s.setMeta("last_build_utc", builtAt.Format(time.RFC3339)); err != nil {
		return err
	}
	if err := s.setMeta("schema_version", "1"); err != nil {
		return err
	}
	return nil
}

func (s *Store) FileStates() (map[string]FileState, error) {
	rows, err := s.db.Query(`SELECT id, path, size, mod_time_unix FROM files`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states := make(map[string]FileState)
	for rows.Next() {
		var state FileState
		var path string
		if err := rows.Scan(&state.ID, &path, &state.Size, &state.ModTimeUnix); err != nil {
			return nil, err
		}
		states[path] = state
	}

	return states, rows.Err()
}

func (s *Store) UpsertFileIndex(path string, record FileRecord, lines []string, lineTerms []LineTerms) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(`
		INSERT INTO files (path, size, mod_time_unix, indexed_at_unix, sha256, line_count)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			size = excluded.size,
			mod_time_unix = excluded.mod_time_unix,
			indexed_at_unix = excluded.indexed_at_unix,
			sha256 = excluded.sha256,
			line_count = excluded.line_count
	`, path, record.Size, record.ModTimeUnix, record.IndexedAt, record.SHA256, record.LineCount); err != nil {
		return err
	}

	var fileID int64
	if err = tx.QueryRow(`SELECT id FROM files WHERE path = ?`, path).Scan(&fileID); err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM lines WHERE file_id = ?`, fileID); err != nil {
		return err
	}
	if _, err = tx.Exec(`DELETE FROM terms WHERE file_id = ?`, fileID); err != nil {
		return err
	}

	lineStmt, err := tx.Prepare(`INSERT INTO lines (file_id, line_no, content) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer lineStmt.Close()

	termStmt, err := tx.Prepare(`INSERT INTO terms (term, file_id, line_no, frequency) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer termStmt.Close()

	for idx, line := range lines {
		if _, err = lineStmt.Exec(fileID, idx+1, line); err != nil {
			return err
		}
	}

	for _, lineTerms := range lineTerms {
		for _, hit := range lineTerms.Hits {
			if _, err = termStmt.Exec(hit.Term, fileID, lineTerms.LineNo, hit.Frequency); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteMissing(seen map[string]struct{}) (int, error) {
	rows, err := s.db.Query(`SELECT path FROM files`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var toDelete []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return 0, err
		}
		if _, ok := seen[path]; !ok {
			toDelete = append(toDelete, path)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	if len(toDelete) == 0 {
		return 0, nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`DELETE FROM files WHERE path = ?`)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	for _, path := range toDelete {
		if _, err := stmt.Exec(path); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return len(toDelete), nil
}

func (s *Store) SearchTerm(term string, pathGlobs []string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.Query(`
		SELECT f.path, t.line_no, l.content, t.frequency
		FROM terms t
		JOIN files f ON f.id = t.file_id
		JOIN lines l ON l.file_id = t.file_id AND l.line_no = t.line_no
		WHERE t.term = ?
		ORDER BY f.path, t.line_no
	`, strings.ToLower(strings.TrimSpace(term)))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]SearchResult, 0, limit)
	for rows.Next() {
		var result SearchResult
		if err := rows.Scan(&result.Path, &result.LineNo, &result.Line, &result.Frequency); err != nil {
			return nil, err
		}
		if !matchesPathGlobs(result.Path, pathGlobs) {
			continue
		}
		results = append(results, result)
		if len(results) >= limit {
			break
		}
	}

	return results, rows.Err()
}

func (s *Store) Grep(opts GrepOptions) ([]SearchResult, error) {
	if opts.Limit <= 0 {
		opts.Limit = 50
	}

	var (
		rows *sql.Rows
		err  error
	)

	if opts.Fixed {
		query := `
			SELECT f.path, l.line_no, l.content
			FROM lines l
			JOIN files f ON f.id = l.file_id
			WHERE instr(` + maybeLowerExpr(opts.IgnoreCase, "l.content") + `, ?) > 0
			ORDER BY f.path, l.line_no
		`
		rows, err = s.db.Query(query, maybeLowerValue(opts.IgnoreCase, opts.Pattern))
	} else {
		rows, err = s.db.Query(`
			SELECT f.path, l.line_no, l.content
			FROM lines l
			JOIN files f ON f.id = l.file_id
			ORDER BY f.path, l.line_no
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regex *regexp.Regexp
	if !opts.Fixed {
		pattern := opts.Pattern
		if opts.IgnoreCase {
			pattern = "(?i)" + pattern
		}
		regex, err = regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
	}

	results := make([]SearchResult, 0, opts.Limit)
	for rows.Next() {
		var result SearchResult
		if err := rows.Scan(&result.Path, &result.LineNo, &result.Line); err != nil {
			return nil, err
		}
		if !matchesPathGlobs(result.Path, opts.PathGlobs) {
			continue
		}
		if !opts.Fixed && !regex.MatchString(result.Line) {
			continue
		}
		results = append(results, result)
		if len(results) >= opts.Limit {
			break
		}
	}

	return results, rows.Err()
}

func (s *Store) Stats() (IndexStats, error) {
	stats := IndexStats{}
	_ = s.db.QueryRow(`SELECT value FROM meta WHERE key = 'root'`).Scan(&stats.Root)
	_ = s.db.QueryRow(`SELECT value FROM meta WHERE key = 'last_build_utc'`).Scan(&stats.LastBuild)
	if stats.LastBuild == "" {
		stats.LastBuild = "unknown"
	}

	if err := s.db.QueryRow(`SELECT COUNT(*) FROM files`).Scan(&stats.Files); err != nil {
		return stats, err
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM lines`).Scan(&stats.Lines); err != nil {
		return stats, err
	}
	if err := s.db.QueryRow(`SELECT COUNT(DISTINCT term) FROM terms`).Scan(&stats.DistinctTerms); err != nil {
		return stats, err
	}
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM terms`).Scan(&stats.TermHits); err != nil {
		return stats, err
	}

	return stats, nil
}

func (s *Store) GetContext(path string, lineNo, before, after int) ([]ContextLine, error) {
	start := lineNo - before
	if start < 1 {
		start = 1
	}
	end := lineNo + after

	rows, err := s.db.Query(`
		SELECT l.line_no, l.content
		FROM lines l
		JOIN files f ON f.id = l.file_id
		WHERE f.path = ? AND l.line_no BETWEEN ? AND ?
		ORDER BY l.line_no
	`, path, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	context := make([]ContextLine, 0, before+after+1)
	for rows.Next() {
		var line ContextLine
		if err := rows.Scan(&line.LineNo, &line.Line); err != nil {
			return nil, err
		}
		context = append(context, line)
	}
	return context, rows.Err()
}

func (s *Store) setMeta(key, value string) error {
	_, err := s.db.Exec(`
		INSERT INTO meta (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}

func (s *Store) meta(key string) (string, error) {
	var value string
	err := s.db.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&value)
	return value, err
}

func matchesPathGlobs(path string, globs []string) bool {
	if len(globs) == 0 {
		return true
	}
	for _, glob := range globs {
		match, err := doublestar.PathMatch(glob, path)
		if err == nil && match {
			return true
		}
	}
	return false
}

func maybeLowerExpr(ignoreCase bool, expr string) string {
	if ignoreCase {
		return "lower(" + expr + ")"
	}
	return expr
}

func maybeLowerValue(ignoreCase bool, value string) string {
	if ignoreCase {
		return strings.ToLower(value)
	}
	return value
}
