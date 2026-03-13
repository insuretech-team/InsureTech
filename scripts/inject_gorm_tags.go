// inject_gorm_tags.go
// Unified Go generation + GORM tag injection tool for InsureTech.
//
// Responsibilities:
//  1) (Optional) Run `buf generate` to produce Go/TS/C# code from protos
//  2) Inject GORM struct tags into all entity .pb.go files, either from
//     trailing `// @inject_tag: gorm:"..."` comments, or auto-derived from
//     the json:"..." + protobuf:"..." tags on every proto struct field
//  3) (Optional, Windows-only) Regenerate proto registry
//
// Usage:
//   go run scripts/inject_gorm_tags.go [--generate] [--registry] [--verbose] [--dry-run]
//
// Flags:
//   --generate       Run `buf generate` before GORM tag injection
//   --registry       Run proto registry generation after injection (Windows only)
//   --verbose        Enable verbose output
//   --dry-run        Show what would be done without writing files

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without writing files")
	generate := flag.Bool("generate", false, "Run buf generate before GORM tag injection")
	registry := flag.Bool("registry", false, "Run proto registry generation after injection (Windows only)")
	flag.Parse()

	fmt.Println("===========================================")
	fmt.Println("InsureTech Proto Generation & GORM Injector")
	fmt.Println("===========================================")
	fmt.Printf("OS: %s, Arch: %s\n\n", runtime.GOOS, runtime.GOARCH)

	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("ERROR: Could not find project root: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Project root: %s\n\n", projectRoot)

	// Step 1 (optional): buf generate
	if *generate {
		fmt.Println("[1/3] Running buf generate...")
		if err := runCmd(projectRoot, *verbose, "buf", "generate"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: buf generate failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("  OK: Code generated\n")
	}

	// Step 2: GORM tag injection
	if *generate {
		fmt.Println("[2/3] Injecting GORM tags...")
	} else {
		fmt.Println("Injecting GORM tags...")
	}

	genGoRoot := filepath.Join(projectRoot, "gen", "go", "insuretech")
	entityDirs, err := findEntityDirs(genGoRoot)
	if err != nil {
		fmt.Printf("ERROR: Failed to find entity directories: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d entity directories\n\n", len(entityDirs))

	var totalFiles, modifiedFiles, skippedFiles int
	var totalInjected int

	for _, dir := range entityDirs {
		files, err := filepath.Glob(filepath.Join(dir, "*.pb.go"))
		if err != nil {
			continue
		}

		for _, file := range files {
			totalFiles++

			changed, injected, err := injectFile(file, *dryRun)
			if err != nil {
				fmt.Printf("  ERROR: %s - %v\n", filepath.Base(file), err)
				continue
			}

			totalInjected += injected
			if injected == 0 {
				skippedFiles++
				if *verbose {
					fmt.Printf("  SKIP (no injectable tags): %s\n", filepath.Base(file))
				}
				continue
			}

			if changed {
				modifiedFiles++
				if *verbose {
					fmt.Printf("  OK: %s (injected %d)\n", filepath.Base(file), injected)
				}
			} else if *verbose {
				fmt.Printf("  OK (already injected): %s (found %d)\n", filepath.Base(file), injected)
			}
		}
	}

	fmt.Println()
	fmt.Printf("  Total files: %d, Modified: %d, Skipped: %d, Injected tags: %d\n", totalFiles, modifiedFiles, skippedFiles, totalInjected)
	if *dryRun {
		fmt.Println("  [DRY RUN - no files modified]")
	}
	fmt.Println("  OK: GORM tags injected")

	// Step 3 (optional): proto registry
	if *registry {
		if *generate {
			fmt.Println("\n[3/3] Proto registry...")
		} else {
			fmt.Println("\nProto registry...")
		}
		if err := generateProtoRegistry(projectRoot, entityDirs, *dryRun); err != nil {
			fmt.Fprintf(os.Stderr, "  WARN: proto registry generation failed: %v\n", err)
		} else {
			fmt.Println("  OK: proto registry updated")
		}
	}

	fmt.Println("\n===========================================")
	fmt.Println("Done!")
	fmt.Println("===========================================")
}

func runCmd(dir string, verbose bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v: %w\n%s", name, args, err, out.String())
	}
	return nil
}

func injectFile(path string, dryRun bool) (changed bool, injected int, err error) {
	in, err := os.Open(path)
	if err != nil {
		return false, 0, err
	}
	defer in.Close()

	var outLines []string
	r := bufio.NewReader(in)
	var pendingInjectTag string

	for {
		line, readErr := r.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return false, 0, readErr
		}
		if line == "" && readErr == io.EOF {
			break
		}

		// Trim only the trailing newline; we will re-add exactly one '\n' on join.
		line = strings.TrimSuffix(line, "\n")

		// Check for leading comment
		if idx := strings.Index(line, "// @inject_tag:"); idx != -1 && !strings.Contains(line, "`") {
			pendingInjectTag = strings.TrimSpace(line[idx:])
			outLines = append(outLines, line)
			if readErr == io.EOF {
				break
			}
			continue
		}

		newLine, injectedThisLine, lineChanged := injectGormOnLine(line, pendingInjectTag)
		if injectedThisLine > 0 {
			injected += injectedThisLine
		}
		if lineChanged {
			changed = true
		}
		outLines = append(outLines, newLine)
		pendingInjectTag = "" // reset after struct field is processed

		if readErr == io.EOF {
			break
		}
	}

	if dryRun || !changed {
		return changed, injected, nil
	}

	newContent := strings.Join(outLines, "\n") + "\n"
	return true, injected, os.WriteFile(path, []byte(newContent), 0o644)
}

// injectGormOnLine injects a gorm struct tag from a trailing or preceding @inject_tag comment.
// Returns: (newLine, injectedCount, changed)
func injectGormOnLine(line string, pendingInjectTag string) (string, int, bool) {
	// Two modes:
	//  1) First-time injection from trailing @inject_tag comment
	//  2) Upgrade existing gorm tags (add serializers) even if comment is already gone

	comment := pendingInjectTag
	idx := strings.Index(line, "// @inject_tag:")
	if idx != -1 {
		comment = strings.TrimSpace(line[idx:])
		line = strings.TrimRight(line[:idx], " \t")
	}

	if comment == "" {
		// If there is no injectable comment, still try:
		//  1) Upgrade existing gorm tag (ensure serializers)
		//  2) Add gorm tag for ANY proto struct field missing a gorm tag
		if upgraded, n, changed := upgradeExistingGormTag(line); changed {
			return upgraded, n, changed
		}
		return addMissingGormTag(line)
	}

	// Only handle gorm injection for now.
	gormVal, ok := parseInjectTagValue(comment, "gorm")
	if !ok {
		return line, 0, false
	}

	// Auto-append serializers for protobuf well-known types and enums.
	// This makes proto-first generated structs usable directly with GORM.
	gormVal = ensureSerializers(line, gormVal)

	// Find the existing struct tag literal (backticks).
	tickStart := strings.Index(line, "`")
	tickEnd := strings.LastIndex(line, "`")
	if tickStart == -1 || tickEnd == -1 || tickEnd <= tickStart {
		// No tag literal to inject into.
		return line, 0, false
	}

	tagLiteral := line[tickStart+1 : tickEnd]

	// If it already has a gorm tag, replace it. Otherwise append it.
	gormIdx := strings.Index(tagLiteral, "gorm:\"")
	var newTagLiteral string
	if gormIdx != -1 {
		valStart := gormIdx + len("gorm:\"")
		valEndRel := strings.Index(tagLiteral[valStart:], "\"")
		if valEndRel != -1 {
			valEnd := valStart + valEndRel
			// Check if we need to prepend the column name
			// If our gormVal doesn't have a column, maybe preserve the old one?
			// Actually just completely replace it, BUT if the user wants `column:id;primarykey` they should have it in the proto.
			// However `addMissingGormTag` adds `column:id`. So if protoc generated `json:"id"`, then `addMissingGormTag` is called? No, it's not because we have an inject comment.
			// Let's ensure the injected tag contains `column:xxx` by extracting it from json tag if not present.
			finalGormVal := gormVal
			if !strings.Contains(finalGormVal, "column:") {
				// extract json column
				jsonKey := "json:\""
				jpos := strings.Index(tagLiteral, jsonKey)
				if jpos != -1 {
					jstart := jpos + len(jsonKey)
					jendRel := strings.Index(tagLiteral[jstart:], "\"")
					if jendRel != -1 {
						jsonVal := tagLiteral[jstart : jstart+jendRel]
						if comma := strings.Index(jsonVal, ","); comma != -1 {
							jsonVal = jsonVal[:comma]
						}
						finalGormVal = "column:" + jsonVal + ";" + finalGormVal
					}
				}
			}
			newTagLiteral = tagLiteral[:gormIdx] + "gorm:\"" + finalGormVal + "\"" + tagLiteral[valEnd+1:]
		} else {
			return line, 0, false
		}
	} else {
		// New tag. Let's make sure it has column info if possible.
		finalGormVal := gormVal
		if !strings.Contains(finalGormVal, "column:") {
			jsonKey := "json:\""
			jpos := strings.Index(tagLiteral, jsonKey)
			if jpos != -1 {
				jstart := jpos + len(jsonKey)
				jendRel := strings.Index(tagLiteral[jstart:], "\"")
				if jendRel != -1 {
					jsonVal := tagLiteral[jstart : jstart+jendRel]
					if comma := strings.Index(jsonVal, ","); comma != -1 {
						jsonVal = jsonVal[:comma]
					}
					finalGormVal = "column:" + jsonVal + ";" + finalGormVal
				}
			}
		}
		newTagLiteral = tagLiteral + " gorm:\"" + finalGormVal + "\""
	}

	newLine := line[:tickStart+1] + newTagLiteral + line[tickEnd:]
	return newLine, 1, true
}

func parseInjectTagValue(comment string, key string) (string, bool) {
	// comment example:
	//   // @inject_tag: gorm:"column:otp_id;not null"
	needle := key + ":\""
	pos := strings.Index(comment, needle)
	if pos == -1 {
		return "", false
	}
	start := pos + len(needle)
	endRel := strings.Index(comment[start:], "\"")
	if endRel == -1 {
		return "", false
	}
	return comment[start : start+endRel], true
}

func addMissingGormTag(line string) (string, int, bool) {
	// Auto-derive gorm struct tag for ANY proto struct field that has
	// protobuf:"..." + json:"..." tags but no gorm:"..." tag.
	//
	// This covers ALL field types: string, int32, bool, enum, *timestamppb.Timestamp, etc.
	// Column name is derived from the json tag value (e.g. json:"user_id,omitempty" → column:user_id).
	// Appropriate serializers are appended for timestamps and enums.

	// Must be a protobuf struct field line (has protobuf:"..." tag).
	if !strings.Contains(line, "protobuf:\"") {
		return line, 0, false
	}

	// Skip protoimpl internal fields (state, sizeCache, unknownFields).
	trimmed := strings.TrimSpace(line)
	for _, skip := range []string{"state ", "sizeCache ", "unknownFields "} {
		if strings.HasPrefix(trimmed, skip) {
			return line, 0, false
		}
	}

	// Must have a struct tag literal (backtick-delimited).
	tickStart := strings.Index(line, "`")
	tickEnd := strings.LastIndex(line, "`")
	if tickStart == -1 || tickEnd == -1 || tickEnd <= tickStart {
		return line, 0, false
	}
	tagLiteral := line[tickStart+1 : tickEnd]

	// Skip if already has a gorm tag.
	if strings.Contains(tagLiteral, "gorm:\"") {
		return line, 0, false
	}

	// Extract json field name from struct tags: json:"foo,omitempty"
	jsonKey := "json:\""
	jpos := strings.Index(tagLiteral, jsonKey)
	if jpos == -1 {
		return line, 0, false
	}
	jstart := jpos + len(jsonKey)
	jendRel := strings.Index(tagLiteral[jstart:], "\"")
	if jendRel == -1 {
		return line, 0, false
	}
	jsonVal := tagLiteral[jstart : jstart+jendRel]
	// jsonVal can be foo,omitempty
	if comma := strings.Index(jsonVal, ","); comma != -1 {
		jsonVal = jsonVal[:comma]
	}
	if jsonVal == "" || jsonVal == "-" {
		return line, 0, false
	}

	// Build gorm tag value.
	gormVal := "column:" + jsonVal

	// Add serializer for timestamp fields.
	if strings.Contains(line, "*timestamppb.Timestamp") {
		gormVal += ";serializer:proto_timestamp"
	}

	// Add serializer for enum fields (protobuf tag contains enum=...).
	if strings.Contains(line, "enum=") {
		gormVal += ";serializer:proto_enum"
	}

	newTagLiteral := tagLiteral + " gorm:\"" + gormVal + "\""
	newLine := line[:tickStart+1] + newTagLiteral + line[tickEnd:]
	return newLine, 1, true
}

func upgradeExistingGormTag(line string) (string, int, bool) {
	// Look for a struct tag literal containing gorm:"..." and ensure serializers are present.
	tickStart := strings.Index(line, "`")
	tickEnd := strings.LastIndex(line, "`")
	if tickStart == -1 || tickEnd == -1 || tickEnd <= tickStart {
		return line, 0, false
	}
	tagLiteral := line[tickStart+1 : tickEnd]
	gormIdx := strings.Index(tagLiteral, "gorm:\"")
	if gormIdx == -1 {
		return line, 0, false
	}
	// Extract gorm value.
	valStart := gormIdx + len("gorm:\"")
	valEndRel := strings.Index(tagLiteral[valStart:], "\"")
	if valEndRel == -1 {
		return line, 0, false
	}
	valEnd := valStart + valEndRel
	gormVal := tagLiteral[valStart:valEnd]
	newGormVal := ensureSerializers(line, gormVal)
	if newGormVal == gormVal {
		return line, 0, false
	}
	newTagLiteral := tagLiteral[:valStart] + newGormVal + tagLiteral[valEnd:]
	newLine := line[:tickStart+1] + newTagLiteral + line[tickEnd:]
	return newLine, 1, true
}

func ensureSerializers(fullLine, gormVal string) string {
	// Timestamp fields.
	// protoc-gen-go uses *timestamppb.Timestamp for google.protobuf.Timestamp.
	if strings.Contains(fullLine, "*timestamppb.Timestamp") {
		if !strings.Contains(gormVal, "serializer:proto_timestamp") {
			gormVal = gormVal + ";serializer:proto_timestamp"
		}
	}

	// Enum fields.
	// Enum presence is encoded in the protobuf struct tag: enum=....
	if strings.Contains(fullLine, "enum=") {
		if !strings.Contains(gormVal, "serializer:proto_enum") {
			gormVal = gormVal + ";serializer:proto_enum"
		}
	}
	return gormVal
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parent
	}
}

func findEntityDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && info.Name() == "v1" {
			parent := filepath.Base(filepath.Dir(path))
			if parent == "entity" {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	return dirs, err
}

// generateProtoRegistry creates backend/inscore/db/ops/proto_registry_gen.go
// with blank imports for all entity packages so protobuf types are registered.
func generateProtoRegistry(projectRoot string, entityDirs []string, dryRun bool) error {
	genGoRoot := filepath.Join(projectRoot, "gen", "go", "insuretech")
	target := filepath.Join(projectRoot, "backend", "inscore", "db", "ops", "proto_registry_gen.go")

	// Build import paths from entity dirs.
	var imports []string
	for _, dir := range entityDirs {
		// dir is absolute: .../gen/go/insuretech/authn/entity/v1
		// We need the relative part after gen/go/insuretech/
		rel, err := filepath.Rel(genGoRoot, dir)
		if err != nil {
			continue
		}
		// Convert backslash to forward slash for Go import path.
		rel = filepath.ToSlash(rel)
		imp := "github.com/newage-saint/insuretech/gen/go/insuretech/" + rel
		imports = append(imports, imp)
	}
	sort.Strings(imports)

	// Build file content.
	var buf bytes.Buffer
	buf.WriteString("// Code generated by inject_gorm_tags.go. DO NOT EDIT.\n")
	buf.WriteString("package ops\n\n")
	buf.WriteString("// This file exists to register all proto entity files with protoregistry.\n")
	buf.WriteString("// It is generated from gen/go/insuretech/**/entity/v1.\n\n")
	buf.WriteString("import (\n")
	for _, imp := range imports {
		buf.WriteString("\t_ \"" + imp + "\"\n")
	}
	buf.WriteString(")\n")

	if dryRun {
		fmt.Printf("  [DRY RUN] Would write %d imports to %s\n", len(imports), filepath.Base(target))
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	fmt.Printf("  Generated: %s (%d imports)\n", filepath.Base(target), len(imports))
	return os.WriteFile(target, buf.Bytes(), 0o644)
}
