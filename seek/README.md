# seek

`seek` is a standalone Go CLI that builds a portable SQLite search index for any project folder. It is designed for large, mixed-language repositories where you want fast exact-word lookup and line-based `grep`-style pattern search without rescanning the filesystem every time.

## Features

- Standalone module with its own `go.mod`
- Embedded SQLite database stored under the indexed project by default
- Exact-word search backed by an inverted index
- Regex or fixed-string line search over indexed content
- Incremental rebuilds using file size and modification time
- Works across mixed source trees, docs, configs, and other text-heavy projects
- Skips common binary, archive, media, and build-output files by default

## Default index location

If `--db` is omitted, the index is stored at:

```text
<project-root>/.seek/index.sqlite
```

## Commands

Set a default root once:

```powershell
seek config set-root E:\Projects\InsureTech
```

Show the saved default root:

```powershell
seek config show
```

Build an index:

```powershell
go run ./cmd/seek build --root E:\Projects\InsureTech
```

Build only selected file types:

```powershell
go run ./cmd/seek build --root E:\Projects\InsureTech --type go,ts,proto
```

Search for an exact word:

```powershell
go run ./cmd/seek search --root E:\Projects\InsureTech --term policy
```

Run a regex search:

```powershell
go run ./cmd/seek grep --root E:\Projects\InsureTech --pattern "Claim.*Status"
```

Run a fast fixed-string search on selected paths:

```powershell
go run ./cmd/seek grep --root E:\Projects\InsureTech --pattern payment --fixed --ignore-case --path "backend/**,api/**"
```

Inspect index stats:

```powershell
go run ./cmd/seek stats --root E:\Projects\InsureTech
```

After you save a default root, `--root` becomes optional:

```powershell
seek build --reset
seek search --term session_token
seek grep --pattern "Claim.*Status"
seek stats
```

Watch and reindex continuously:

```powershell
seek watch --root E:\Projects\InsureTech --interval 2
```

Use the interactive shell:

```powershell
seek interactive --root E:\Projects\InsureTech --type go,ts --context 2
```

## Useful build flags

- `--include-ext ".go,.ts,.js,.md,.yaml,.yml,.proto"` limits indexing to selected extensions
- `--type "go,ts,proto"` uses named extension groups
- `--exclude-dir "vendor,tmp,.output"` adds more skipped directories
- `--exclude-glob "documentation/**,web_shared/**"` excludes specific path patterns
- `--max-file-size-mb 8` raises the default per-file size cap
- `--reset` discards the previous index before rebuilding
- `--save-root` saves the resolved root as the new default root while running the command
- `--context 2` shows surrounding lines for `search` and `grep`

## .seekignore

You can create a `.seekignore` file in the indexed project root to exclude extra paths from indexing.

Example:

```text
docs/
documentation/
*.tmp
gen/**
```

## Notes

- Exact-word search is case-insensitive because terms are normalized to lowercase.
- `grep --fixed` is the fastest way to search for literal strings.
- Regex search is line-based and works best with `--path` filters on very large indexes.
- `watch` is filesystem-event based and uses `--interval` only for a lightweight root-level fallback check.
- `interactive` is a lightweight REPL for repeated search and grep queries.
