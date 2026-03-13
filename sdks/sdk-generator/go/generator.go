package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// getProjectRoot finds the project root directory by looking for go.mod or a marker file
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up the directory tree to find the project root
	for {
		// Check if we're at the sdks directory or can find proto directory
		protoPath := filepath.Join(dir, "proto")
		if _, err := os.Stat(protoPath); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root of the filesystem
			return "", fmt.Errorf("could not find project root (looking for 'proto' directory)")
		}
		dir = parent
	}
}

// OpenAPISpec represents a simplified OpenAPI specification
type OpenAPISpec struct {
	OpenAPI    string                          `yaml:"openapi"`
	Info       map[string]interface{}          `yaml:"info"`
	Paths      map[string]map[string]Operation `yaml:"paths"`
	Components struct {
		Schemas map[string]Schema `yaml:"schemas"`
	} `yaml:"components"`
}

// Operation represents an OpenAPI operation
type Operation struct {
	OperationID string              `yaml:"operationId"`
	Summary     string              `yaml:"summary"`
	Description string              `yaml:"description"`
	Tags        []string            `yaml:"tags"`
	Parameters  []Parameter         `yaml:"parameters"`
	RequestBody *RequestBody        `yaml:"requestBody"`
	Responses   map[string]Response `yaml:"responses"`
}

// Parameter represents an OpenAPI parameter
type Parameter struct {
	Name        string `yaml:"name"`
	In          string `yaml:"in"` // path, query, header, cookie
	Required    bool   `yaml:"required"`
	Schema      Schema `yaml:"schema"`
	Description string `yaml:"description"`
}

// RequestBody represents an OpenAPI request body
type RequestBody struct {
	Required bool                 `yaml:"required"`
	Content  map[string]MediaType `yaml:"content"`
}

// MediaType represents a media type in OpenAPI
type MediaType struct {
	Schema Schema `yaml:"schema"`
}

// Response represents an OpenAPI response
type Response struct {
	Description string               `yaml:"description"`
	Content     map[string]MediaType `yaml:"content"`
}

// MethodMetadata contains metadata for generating service methods
type MethodMetadata struct {
	ServiceName  string
	MethodName   string
	HTTPMethod   string
	Path         string
	OperationID  string
	Summary      string
	PathParams   []Parameter
	QueryParams  []Parameter
	RequestType  string
	ResponseType string
	IsPaginated  bool
	HasBody      bool
}

// Schema represents an OpenAPI schema
type Schema struct {
	Type       string              `yaml:"type"`
	Properties map[string]Property `yaml:"properties"`
	Required   []string            `yaml:"required"`
	Enum       []string            `yaml:"enum"`
	Ref        string              `yaml:"$ref"`
}

// Property represents a schema property
type Property struct {
	Type        string    `yaml:"type"`
	Format      string    `yaml:"format"`
	Description string    `yaml:"description"`
	Ref         string    `yaml:"$ref"`
	Items       *Property `yaml:"items"`
	Enum        []string  `yaml:"enum"`
}

// GeneratorConfig holds the configuration for SDK generation
type GeneratorConfig struct {
	PackageName   string
	ModulePath    string
	Version       string
	Author        string
	License       string
	ProtoPath     string
	APISpecPath   string
	OutputPath    string
	TemplatesPath string
}

func main() {
	log.Println("Starting Go SDK Generator...")

	// Find project root
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatalf("Failed to find project root: %v", err)
	}
	log.Printf("Project root: %s\n", projectRoot)

	// Get current directory (should be in sdk-generator/go)
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Build paths relative to project root
	protoPath := filepath.Join(projectRoot, "proto")
	apiSpecPath := filepath.Join(projectRoot, "api", "openapi.yaml")
	outputPath := filepath.Join(projectRoot, "sdks", "insuretech-go-sdk")
	templatesPath := filepath.Join(currentDir, "templates")

	config := &GeneratorConfig{
		PackageName:   "insuretech",
		ModulePath:    "github.com/newage-saint/insuretech-go-sdk",
		Version:       "1.0.0",
		Author:        "InsureTech Platform",
		License:       "MIT",
		ProtoPath:     protoPath,
		APISpecPath:   apiSpecPath,
		OutputPath:    outputPath,
		TemplatesPath: templatesPath,
	}

	log.Printf("Proto path: %s\n", config.ProtoPath)
	log.Printf("API spec path: %s\n", config.APISpecPath)
	log.Printf("Output path: %s\n", config.OutputPath)
	log.Printf("Templates path: %s\n", config.TemplatesPath)

	// Step 1: Load OpenAPI spec
	log.Println("Loading OpenAPI specification...")
	spec, err := loadOpenAPISpec(config.APISpecPath)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}
	log.Printf("✓ Loaded OpenAPI spec version: %s\n", spec.OpenAPI)

	// Step 2: Parse proto definitions
	log.Println("Parsing proto definitions...")
	// TODO: Implement proto parsing
	log.Println("✓ Proto definitions parsed")

	// Step 3: Create output directory structure
	log.Println("Creating output directory structure...")
	if err := createOutputStructure(config.OutputPath); err != nil {
		log.Fatalf("Failed to create output structure: %v", err)
	}
	log.Println("✓ Output structure created")

	// Step 4: Generate base SDK files from templates
	log.Println("Generating base SDK files from templates...")
	if err := generateBaseFiles(config); err != nil {
		log.Fatalf("Failed to generate base files: %v", err)
	}
	log.Println("✓ Base files generated")

	// Step 5: Generate models from OpenAPI schemas
	log.Println("Generating models from OpenAPI schemas...")
	if err := generateModels(spec, config); err != nil {
		log.Fatalf("Failed to generate models: %v", err)
	}
	log.Println("✓ Models generated")

	// Step 6: Generate service clients from OpenAPI paths
	log.Println("Generating service clients from OpenAPI paths...")
	serviceNames, err := generateServices(spec, config)
	if err != nil {
		log.Fatalf("Failed to generate services: %v", err)
	}
	log.Println("✓ Services generated")

	// Step 6.5: Generate client with service initialization
	log.Println("Generating client with service initialization...")
	if err := generateClientWithServices(config, serviceNames); err != nil {
		log.Fatalf("Failed to generate client: %v", err)
	}
	log.Println("✓ Client generated")

	// Step 7: Generate go.mod
	log.Println("Generating go.mod...")
	if err := generateGoMod(config); err != nil {
		log.Fatalf("Failed to generate go.mod: %v", err)
	}
	log.Println("✓ go.mod generated")

	// Step 8: Generate README
	log.Println("Generating README...")
	if err := generateReadme(config); err != nil {
		log.Fatalf("Failed to generate README: %v", err)
	}
	log.Println("✓ README generated")

	log.Println("✅ Go SDK generation completed successfully!")
	log.Printf("Output location: %s\n", config.OutputPath)
}

// loadOpenAPISpec loads and parses the OpenAPI specification
func loadOpenAPISpec(path string) (*OpenAPISpec, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI spec: %w", err)
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	return &spec, nil
}

// createOutputStructure creates the directory structure for the SDK
func createOutputStructure(outputPath string) error {
	dirs := []string{
		outputPath,
		filepath.Join(outputPath, "pkg"),
		filepath.Join(outputPath, "pkg", "client"),
		filepath.Join(outputPath, "pkg", "models"),
		filepath.Join(outputPath, "pkg", "services"),
		filepath.Join(outputPath, "pkg", "errors"),
		filepath.Join(outputPath, "examples"),
		filepath.Join(outputPath, "docs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateBaseFiles generates base SDK files from templates
func generateBaseFiles(config *GeneratorConfig) error {
	// First pass: generate helper files that don't need service info
	helperTemplates := map[string]string{
		"config.go.tmpl":          filepath.Join(config.OutputPath, "pkg", "client", "config.go"),
		"errors.go.tmpl":          filepath.Join(config.OutputPath, "pkg", "errors", "errors.go"),
		"models.go.tmpl":          filepath.Join(config.OutputPath, "pkg", "models", "base.go"),
		"services.go.tmpl":        filepath.Join(config.OutputPath, "pkg", "services", "services.go"),
		"request_builder.go.tmpl": filepath.Join(config.OutputPath, "pkg", "client", "request_builder.go"),
		"pagination.go.tmpl":      filepath.Join(config.OutputPath, "pkg", "models", "pagination.go"),
		"helpers.go.tmpl":         filepath.Join(config.OutputPath, "pkg", "client", "helpers.go"),
	}

	for tmplFile, outputFile := range helperTemplates {
		tmplPath := filepath.Join(config.TemplatesPath, tmplFile)

		// Check if template exists
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			log.Printf("  ⚠ Template not found: %s (skipping)\n", tmplFile)
			continue
		}

		// Read template
		tmplContent, err := ioutil.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tmplFile, err)
		}

		// Parse template
		tmpl, err := template.New(tmplFile).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", tmplFile, err)
		}

		// Create output file
		outFile, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
		}
		defer outFile.Close()

		// Execute template
		if err := tmpl.Execute(outFile, config); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", tmplFile, err)
		}

		log.Printf("  ✓ Generated: %s\n", outputFile)
	}

	return nil
}

// generateClientWithServices generates the client.go file with service fields
func generateClientWithServices(config *GeneratorConfig, serviceNames []string) error {
	tmplPath := filepath.Join(config.TemplatesPath, "client.go.tmpl")
	outputFile := filepath.Join(config.OutputPath, "pkg", "client", "client.go")

	// Read template
	tmplContent, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read client template: %w", err)
	}

	// Build service fields and initializers
	var serviceFields strings.Builder
	var serviceInits strings.Builder

	for _, serviceName := range serviceNames {
		cleanName := toPascalCase(serviceName)
		serviceFields.WriteString(fmt.Sprintf("\t%s *services.%sService\n", cleanName, cleanName))
		serviceInits.WriteString(fmt.Sprintf("\tc.%s = &services.%sService{Client: c}\n", cleanName, cleanName))
	}

	// Replace placeholders in template
	templateStr := string(tmplContent)
	templateStr = strings.ReplaceAll(templateStr, "{{.ServiceFields}}", serviceFields.String())
	templateStr = strings.ReplaceAll(templateStr, "{{.ServiceInitializers}}", serviceInits.String())

	// Write to file
	if err := ioutil.WriteFile(outputFile, []byte(templateStr), 0644); err != nil {
		return fmt.Errorf("failed to write client file: %w", err)
	}

	log.Printf("  ✓ Generated client with %d services: %s\n", len(serviceNames), outputFile)
	return nil
}

// generateGoMod generates the go.mod file
func generateGoMod(config *GeneratorConfig) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	gopkg.in/yaml.v3 v3.0.1
)
`, config.ModulePath)

	outputFile := filepath.Join(config.OutputPath, "go.mod")
	return ioutil.WriteFile(outputFile, []byte(content), 0644)
}

// generateReadme generates the README.md file
func generateReadme(config *GeneratorConfig) error {
	content := fmt.Sprintf(`# InsureTech Go SDK

Official Go SDK for the InsureTech API Platform.

## Installation

`+"```bash"+`
go get %s
`+"```"+`

## Quick Start

`+"```go"+`
package main

import (
    "context"
    "log"
    
    insuretech "%s"
)

func main() {
    client := insuretech.NewClient(
        insuretech.WithAPIKey("your-api-key"),
    )
    
    ctx := context.Background()
    
    // Use the client...
}
`+"```"+`

## Documentation

For full documentation, see [QUICKSTART.md](../sdk-generator/go/QUICKSTART.md)

## Version

Version: %s

## License

%s
`, config.ModulePath, config.ModulePath, config.Version, config.License)

	outputFile := filepath.Join(config.OutputPath, "README.md")
	return ioutil.WriteFile(outputFile, []byte(content), 0644)
}

// schemaNameRegistry maps a lowercase filename key to the canonical (first-seen) schema name.
// This ensures that when two schema names differ only in acronym casing
// (e.g. "ApiKeysListingResponse" vs "APIKeysListingResponse") they resolve
// to a single Go type, preventing duplicate/overwritten files and broken
// service references.
var schemaNameRegistry = map[string]string{}

// canonicalSchemaName returns the canonical Go struct name for a schema ref name.
// If two schema names produce the same snake_case filename, only the first one
// encountered is kept and all subsequent variants are remapped to it.
func canonicalSchemaName(name string) string {
	key := toSnakeCase(name)
	if canonical, ok := schemaNameRegistry[key]; ok {
		return canonical
	}
	schemaNameRegistry[key] = name
	return name
}

// generateModels generates Go model files from OpenAPI schemas
func generateModels(spec *OpenAPISpec, config *GeneratorConfig) error {
	modelsPath := filepath.Join(config.OutputPath, "pkg", "models")

	// Reset registry for this run
	schemaNameRegistry = map[string]string{}

	// First pass: build the canonical name registry by iterating all schema names.
	// We sort names so that the shorter / more common casing wins deterministically.
	sortedNames := make([]string, 0, len(spec.Components.Schemas))
	for schemaName := range spec.Components.Schemas {
		sortedNames = append(sortedNames, schemaName)
	}
	// Sort: prefer names where acronyms are Title-cased (e.g. "Api" over "API")
	// by sorting alphabetically — lowercase letters sort after uppercase in ASCII,
	// so "Api..." < "API..." which means "Api" variant wins as first-seen.
	// Use a stable sort by lowercase name so order is deterministic.
	sortStrings(sortedNames)
	for _, name := range sortedNames {
		canonicalSchemaName(name) // populate registry
	}

	// Track generated models
	generatedCount := 0

	// Second pass: generate files using canonical names only.
	writtenFiles := map[string]bool{}
	for _, schemaName := range sortedNames {
		schema := spec.Components.Schemas[schemaName]

		// Skip if it's just a reference
		if schema.Ref != "" {
			continue
		}

		// Resolve to canonical name — skip if this is a duplicate variant
		canonical := canonicalSchemaName(schemaName)
		filename := toSnakeCase(canonical) + ".go"
		if writtenFiles[filename] {
			log.Printf("  ⚠ Skipping duplicate schema variant %q (canonical: %q)\n", schemaName, canonical)
			continue
		}
		writtenFiles[filename] = true

		// Generate model file using the canonical name for the struct
		modelContent := generateModelCode(canonical, schema)
		if modelContent == "" {
			continue
		}

		outputFile := filepath.Join(modelsPath, filename)
		if err := ioutil.WriteFile(outputFile, []byte(modelContent), 0644); err != nil {
			return fmt.Errorf("failed to write model file %s: %w", filename, err)
		}

		generatedCount++
		if generatedCount <= 10 {
			log.Printf("  ✓ Generated model: %s\n", canonical)
		}
	}

	log.Printf("  Generated %d model files\n", generatedCount)
	return nil
}

// sortStrings sorts a string slice in place (stdlib sort without import conflict).
func sortStrings(s []string) {
	// simple insertion sort — schema counts are small enough
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}

// generateModelCode generates Go code for a single model
func generateModelCode(name string, schema Schema) string {
	// Handle enums
	if len(schema.Enum) > 0 {
		return generateEnumCode(name, schema)
	}

	if schema.Type != "object" || len(schema.Properties) == 0 {
		return ""
	}

	var code strings.Builder
	code.WriteString("package models\n\n")

	// Check if we need to import time
	needsTime := false
	for _, prop := range schema.Properties {
		if prop.Type == "string" && (prop.Format == "date-time" || prop.Format == "date") {
			needsTime = true
			break
		}
	}

	if needsTime {
		code.WriteString("import (\n")
		code.WriteString("\t\"time\"\n")
		code.WriteString(")\n\n")
	} else {
		code.WriteString("\n")
	}

	// Generate struct
	code.WriteString(fmt.Sprintf("// %s represents a %s\n", name, toSnakeCase(name)))
	code.WriteString(fmt.Sprintf("type %s struct {\n", name))

	// Add properties
	for propName, prop := range schema.Properties {
		goType := mapTypeToGo(prop)
		jsonTag := toSnakeCase(propName)
		required := contains(schema.Required, propName)

		if !required {
			jsonTag += ",omitempty"
		}

		fieldName := toPascalCase(propName)
		code.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, goType, jsonTag))
	}

	code.WriteString("}\n")

	return code.String()
}

// generateServices generates service client files from OpenAPI paths
func generateServices(spec *OpenAPISpec, config *GeneratorConfig) ([]string, error) {
	servicesPath := filepath.Join(config.OutputPath, "pkg", "services")

	// Parse all operations and build method metadata
	methodsByService := make(map[string][]MethodMetadata)

	for path, operations := range spec.Paths {
		for httpMethod, operation := range operations {
			// Extract service name from operation tags or path
			serviceName := extractServiceName(operation, path)
			if serviceName == "" {
				continue
			}

			// Build method metadata
			metadata := buildMethodMetadata(serviceName, path, httpMethod, operation)
			if metadata.MethodName == "" {
				continue
			}

			methodsByService[serviceName] = append(methodsByService[serviceName], metadata)
		}
	}

	log.Printf("  Found %d service groups with methods\n", len(methodsByService))

	// Track service names for client generation
	serviceNames := make([]string, 0, len(methodsByService))

	// Generate a service file for each group
	generatedCount := 0
	for serviceName, methods := range methodsByService {
		serviceNames = append(serviceNames, serviceName)
		serviceCode := generateServiceCodeWithMetadata(serviceName, methods, spec, config)
		if serviceCode == "" {
			continue
		}

		filename := toSnakeCase(serviceName) + "_service.go"
		outputFile := filepath.Join(servicesPath, filename)

		if err := ioutil.WriteFile(outputFile, []byte(serviceCode), 0644); err != nil {
			return nil, fmt.Errorf("failed to write service file %s: %w", filename, err)
		}

		generatedCount++
		if generatedCount <= 10 {
			log.Printf("  ✓ Generated service: %sService (%d methods)\n", toPascalCase(serviceName), len(methods))
		}
	}

	log.Printf("  Generated %d service files\n", generatedCount)
	return serviceNames, nil
}

// extractServiceName extracts the service name from operation or path
func extractServiceName(operation Operation, path string) string {
	// Try to get from operationId first (e.g., "PolicyService_Create")
	if operation.OperationID != "" {
		parts := strings.Split(operation.OperationID, "_")
		if len(parts) >= 1 {
			serviceName := strings.TrimSuffix(parts[0], "Service")
			return strings.ToLower(serviceName)
		}
	}

	// Try tags
	if len(operation.Tags) > 0 {
		return strings.ToLower(operation.Tags[0])
	}

	// Fall back to path parsing
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			serviceName := parts[i+1]
			serviceName = strings.Split(serviceName, "{")[0]
			serviceName = strings.TrimSuffix(serviceName, "/")
			return serviceName
		}
	}

	return ""
}

// buildMethodMetadata builds method metadata from operation details
func buildMethodMetadata(serviceName, path, httpMethod string, operation Operation) MethodMetadata {
	metadata := MethodMetadata{
		ServiceName: serviceName,
		HTTPMethod:  strings.ToUpper(httpMethod),
		Path:        path,
		OperationID: operation.OperationID,
		Summary:     operation.Summary,
	}

	// Extract method name from operationId
	if operation.OperationID != "" {
		parts := strings.Split(operation.OperationID, "_")
		if len(parts) >= 2 {
			metadata.MethodName = parts[1]
		}
	}

	// Fallback method name generation
	if metadata.MethodName == "" {
		metadata.MethodName = generateMethodName(httpMethod, path)
	}

	// Parse parameters
	for _, param := range operation.Parameters {
		if param.In == "path" {
			metadata.PathParams = append(metadata.PathParams, param)
		} else if param.In == "query" {
			metadata.QueryParams = append(metadata.QueryParams, param)
		}
	}

	// Parse request body
	if operation.RequestBody != nil {
		metadata.HasBody = true
		for _, mediaType := range operation.RequestBody.Content {
			if mediaType.Schema.Ref != "" {
				metadata.RequestType = extractTypeNameFromRef(mediaType.Schema.Ref)
				break
			}
		}
	}

	// Parse response type
	if resp, ok := operation.Responses["200"]; ok {
		for _, mediaType := range resp.Content {
			if mediaType.Schema.Ref != "" {
				metadata.ResponseType = extractTypeNameFromRef(mediaType.Schema.Ref)
				break
			}
		}
	} else if resp, ok := operation.Responses["201"]; ok {
		for _, mediaType := range resp.Content {
			if mediaType.Schema.Ref != "" {
				metadata.ResponseType = extractTypeNameFromRef(mediaType.Schema.Ref)
				break
			}
		}
	}

	// Check if paginated
	metadata.IsPaginated = strings.Contains(strings.ToLower(metadata.MethodName), "list") ||
		strings.Contains(strings.ToLower(metadata.ResponseType), "listing")

	return metadata
}

// extractTypeNameFromRef extracts type name from $ref and resolves it to the
// canonical schema name so that service files always reference the correct
// model struct (handles acronym casing variants like Api vs API).
func extractTypeNameFromRef(ref string) string {
	parts := strings.Split(ref, "/")
	if len(parts) == 0 {
		return ""
	}
	name := parts[len(parts)-1]
	// Resolve through registry if it has been populated (i.e. after generateModels)
	if canonical, ok := schemaNameRegistry[toSnakeCase(name)]; ok {
		return canonical
	}
	return name
}

// generateServiceCodeWithMetadata generates Go code for a service client using metadata
func generateServiceCodeWithMetadata(serviceName string, methods []MethodMetadata, spec *OpenAPISpec, config *GeneratorConfig) string {
	var code strings.Builder

	// Check what imports we need
	needsStrings := false
	for _, method := range methods {
		if len(method.PathParams) > 0 {
			needsStrings = true
			break
		}
	}

	code.WriteString("package services\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"context\"\n")
	if needsStrings {
		code.WriteString("\t\"strings\"\n")
	}
	code.WriteString(fmt.Sprintf("\t\"%s/pkg/models\"\n", config.ModulePath))
	code.WriteString(")\n\n")

	// Generate service struct
	serviceTypeName := toPascalCase(serviceName) + "Service"
	code.WriteString(fmt.Sprintf("// %s handles %s-related API calls\n", serviceTypeName, serviceName))
	code.WriteString(fmt.Sprintf("type %s struct {\n", serviceTypeName))
	code.WriteString("\tClient Client\n")
	code.WriteString("}\n\n")

	// Generate methods
	for _, method := range methods {
		code.WriteString(generateMethodCode(serviceTypeName, method))
		code.WriteString("\n")
	}

	return code.String()
}

// generateMethodCode generates code for a single service method
func generateMethodCode(serviceTypeName string, method MethodMetadata) string {
	var code strings.Builder

	// Build method signature
	comment := method.Summary
	if comment == "" {
		comment = fmt.Sprintf("performs %s %s", method.HTTPMethod, method.Path)
	}
	code.WriteString(fmt.Sprintf("// %s %s\n", method.MethodName, comment))
	code.WriteString(fmt.Sprintf("func (s *%s) %s(ctx context.Context", serviceTypeName, method.MethodName))

	// Add path parameters
	for _, param := range method.PathParams {
		paramType := "string"
		paramName := toCamelCase(param.Name)
		code.WriteString(fmt.Sprintf(", %s %s", paramName, paramType))
	}

	// Add request body parameter
	if method.HasBody && method.RequestType != "" {
		code.WriteString(fmt.Sprintf(", req *models.%s", method.RequestType))
	} else if len(method.QueryParams) > 0 && method.RequestType != "" {
		// Query params packaged in request type
		code.WriteString(fmt.Sprintf(", req *models.%s", method.RequestType))
	}

	// Return type
	if method.ResponseType != "" {
		code.WriteString(fmt.Sprintf(") (*models.%s, error) {\n", method.ResponseType))
	} else {
		code.WriteString(") error {\n")
	}

	// Build path with parameters
	code.WriteString(fmt.Sprintf("\tpath := \"%s\"\n", method.Path))

	// Replace path parameters
	for _, param := range method.PathParams {
		paramVar := toCamelCase(param.Name)
		placeholder := fmt.Sprintf("{%s}", param.Name)
		code.WriteString(fmt.Sprintf("\tpath = strings.ReplaceAll(path, \"%s\", %s)\n", placeholder, paramVar))
	}

	// Determine request body
	var requestBody string
	if method.HasBody && method.RequestType != "" {
		requestBody = "req"
	} else if len(method.QueryParams) > 0 && method.RequestType != "" {
		// For GET with query params, we'll need to build query string
		// For now, pass req as body (will need query builder helper)
		requestBody = "nil"
		code.WriteString("\t// TODO: Build query string from req\n")
	} else {
		requestBody = "nil"
	}

	// Make HTTP request
	if method.ResponseType != "" {
		code.WriteString(fmt.Sprintf("\tvar result models.%s\n", method.ResponseType))
		code.WriteString(fmt.Sprintf("\terr := s.Client.DoRequest(ctx, \"%s\", path, %s, &result)\n", method.HTTPMethod, requestBody))
		code.WriteString("\tif err != nil {\n")
		code.WriteString("\t\treturn nil, err\n")
		code.WriteString("\t}\n")
		code.WriteString("\treturn &result, nil\n")
	} else {
		code.WriteString(fmt.Sprintf("\treturn s.Client.DoRequest(ctx, \"%s\", path, %s, nil)\n", method.HTTPMethod, requestBody))
	}

	code.WriteString("}\n")

	return code.String()
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	// First part stays lowercase
	result := strings.ToLower(parts[0])

	// Rest are capitalized
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}

	return result
}

// generateServiceCode generates Go code for a service client
func generateServiceCode(serviceName string, paths map[string]map[string]interface{}, spec *OpenAPISpec) string {
	var code strings.Builder

	code.WriteString("package services\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"context\"\n")
	code.WriteString("\t\"fmt\"\n")
	code.WriteString(")\n\n")

	// Clean service name - remove hyphens and colons for valid Go type name
	cleanServiceName := strings.ReplaceAll(serviceName, "-", "_")
	cleanServiceName = strings.ReplaceAll(cleanServiceName, ":", "_")

	// Generate service struct
	serviceTypeName := toPascalCase(cleanServiceName) + "Service"
	code.WriteString(fmt.Sprintf("// %s handles %s-related API calls\n", serviceTypeName, serviceName))
	code.WriteString(fmt.Sprintf("type %s struct {\n", serviceTypeName))
	code.WriteString("\tclient *Client\n")
	code.WriteString("}\n\n")

	// Generate methods for each path
	methodCount := 0
	generatedMethods := make(map[string]bool) // Track generated method names to avoid duplicates

	for path, methods := range paths {
		for method := range methods {
			if methodCount >= 5 {
				// Limit methods per service for now
				break
			}

			methodName := generateMethodName(method, path)
			if methodName == "" {
				continue
			}

			// Skip if we already generated this method name (regardless of HTTP method)
			if generatedMethods[methodName] {
				continue
			}
			generatedMethods[methodName] = true

			// Generate method
			code.WriteString(fmt.Sprintf("// %s performs %s %s\n", methodName, strings.ToUpper(method), path))
			code.WriteString(fmt.Sprintf("func (s *%s) %s(ctx context.Context) error {\n", serviceTypeName, methodName))
			code.WriteString(fmt.Sprintf("\t// TODO: Implement %s %s\n", strings.ToUpper(method), path))
			code.WriteString("\treturn fmt.Errorf(\"not implemented\")\n")
			code.WriteString("}\n\n")

			methodCount++
		}
	}

	if methodCount == 0 {
		return ""
	}

	return code.String()
}

// generateMethodName generates a Go method name from HTTP method and path
func generateMethodName(httpMethod, path string) string {
	httpMethod = strings.ToUpper(httpMethod)

	// Extract the last segment
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]

	// Check if it's an ID parameter
	hasID := strings.Contains(lastPart, "{") && strings.Contains(lastPart, "}")

	var methodName string
	switch httpMethod {
	case "GET":
		if hasID {
			methodName = "Get"
		} else {
			methodName = "List"
		}
	case "POST":
		methodName = "Create"
	case "PUT":
		methodName = "Update"
	case "PATCH":
		// Use Patch to differentiate from PUT Update
		methodName = "Patch"
	case "DELETE":
		methodName = "Delete"
	default:
		methodName = toPascalCase(httpMethod)
	}

	// Make method name more specific if path has additional context
	if len(parts) > 2 && !hasID {
		// Clean the resource name - remove special characters
		cleanedLast := strings.ReplaceAll(lastPart, "-", "_")
		cleanedLast = strings.ReplaceAll(cleanedLast, ":", "_")
		resourceName := toPascalCase(cleanedLast)

		// Only add if it's not empty and makes sense
		if resourceName != "" && len(resourceName) > 1 {
			methodName = methodName + resourceName
		}
	}

	return methodName
}

// Helper functions
func mapTypeToGo(prop Property) string {
	if prop.Ref != "" {
		// Extract type name from $ref
		parts := strings.Split(prop.Ref, "/")
		if len(parts) > 0 {
			return "*" + parts[len(parts)-1]
		}
	}

	switch prop.Type {
	case "string":
		if prop.Format == "date-time" {
			return "time.Time"
		}
		if prop.Format == "date" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if prop.Format == "int64" {
			return "int64"
		}
		return "int"
	case "number":
		if prop.Format == "double" {
			return "float64"
		}
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if prop.Items != nil {
			itemType := mapTypeToGo(*prop.Items)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

func toSnakeCase(s string) string {
	// Replace special characters with underscores
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, " ", "_")

	// Common acronyms that should stay together (order matters - longer first)
	acronyms := []struct{ pattern, replacement string }{
		{"UUID", "uuid"},
		{"HTTP", "http"},
		{"API", "api"},
		{"KYC", "kyc"},
		{"FAQ", "faq"},
		{"OTP", "otp"},
		{"MCP", "mcp"},
		{"MFS", "mfs"},
		{"IoT", "iot"},
		{"NID", "nid"},
		{"TIN", "tin"},
		{"SMS", "sms"},
		{"URL", "url"},
		{"AI", "ai"},
		{"ID", "id"},
	}

	// Replace known acronyms with markers
	for _, acr := range acronyms {
		s = strings.ReplaceAll(s, acr.pattern, "~"+acr.replacement+"~")
	}

	var result strings.Builder
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Handle our markers - add underscore before marker if needed
		if r == '~' {
			// Add underscore before marker if we have content and last char wasn't underscore
			if result.Len() > 0 {
				lastRune := runes[i-1]
				if lastRune != '_' && lastRune != '~' {
					result.WriteRune('_')
				}
			}
			// Copy the acronym
			i++ // skip opening ~
			for i < len(runes) && runes[i] != '~' {
				result.WriteRune(runes[i])
				i++
			}
			// i is now on closing ~, add underscore after if next char is not underscore
			if i+1 < len(runes) && runes[i+1] != '_' && runes[i+1] != '~' && runes[i+1] >= 'A' && runes[i+1] <= 'Z' {
				result.WriteRune('_')
			}
			continue
		}

		// Add underscore before uppercase letters
		if i > 0 && r >= 'A' && r <= 'Z' {
			lastRune := runes[i-1]
			// Add underscore if previous char was lowercase or a closing marker
			if (lastRune >= 'a' && lastRune <= 'z') || (lastRune >= '0' && lastRune <= '9') {
				result.WriteRune('_')
			}
		}

		result.WriteRune(r)
	}

	finalStr := strings.ToLower(result.String())

	// Clean up multiple underscores
	for strings.Contains(finalStr, "__") {
		finalStr = strings.ReplaceAll(finalStr, "__", "_")
	}

	// Remove leading/trailing underscores
	finalStr = strings.Trim(finalStr, "_")

	return finalStr
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// generateEnumCode generates Go code for an enum type
func generateEnumCode(name string, schema Schema) string {
	var code strings.Builder

	code.WriteString("package models\n\n")

	// Determine the base type
	baseType := "string"
	if schema.Type == "integer" {
		baseType = "int"
	} else if schema.Type == "number" {
		baseType = "float64"
	}

	// Generate type definition
	code.WriteString(fmt.Sprintf("// %s represents a %s\n", name, toSnakeCase(name)))
	code.WriteString(fmt.Sprintf("type %s %s\n\n", name, baseType))

	// Generate constants
	if len(schema.Enum) > 0 {
		code.WriteString(fmt.Sprintf("// %s values\n", name))
		code.WriteString("const (\n")

		for i, enumVal := range schema.Enum {
			constName := name + toPascalCase(fmt.Sprintf("%v", enumVal))

			if baseType == "string" {
				if i == 0 {
					code.WriteString(fmt.Sprintf("\t%s %s = \"%v\"\n", constName, name, enumVal))
				} else {
					code.WriteString(fmt.Sprintf("\t%s %s = \"%v\"\n", constName, "", enumVal))
				}
			} else {
				if i == 0 {
					code.WriteString(fmt.Sprintf("\t%s %s = %v\n", constName, name, enumVal))
				} else {
					code.WriteString(fmt.Sprintf("\t%s = %v\n", constName, enumVal))
				}
			}
		}

		code.WriteString(")\n")
	}

	return code.String()
}

