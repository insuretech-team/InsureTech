package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// writeFileIfChanged writes content to path only when the file doesn't exist or
// its content differs — preventing unnecessary git churn on repeated pipeline runs.
func writeFileIfChanged(path string, content []byte, perm os.FileMode) error {
	if existing, err := ioutil.ReadFile(path); err == nil {
		if bytes.Equal(existing, content) {
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, content, perm)
}

func writeFileIfChangedStr(path string, content string, perm os.FileMode) error {
	return writeFileIfChanged(path, []byte(content), perm)
}

// GeneratorConfig holds configuration
type GeneratorConfig struct {
	OpenAPIPath   string
	OutputPath    string
	TemplatesPath string
	PackageName   string
	Version       string
	License       string
}

func main() {
	fmt.Println("🚀 TypeScript SDK Generator (hey-api + custom)")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println()

	// Set up paths
	workspaceRoot := filepath.Join("..", "..", "..")
	openapiPath := filepath.Join(workspaceRoot, "api", "openapi.yaml")
	outputPath := filepath.Join(workspaceRoot, "sdks", "insuretech-typescript-sdk")
	templatesPath := "./templates"

	config := &GeneratorConfig{
		OpenAPIPath:   openapiPath,
		OutputPath:    outputPath,
		TemplatesPath: templatesPath,
		PackageName:   "@lifeplus/insuretech-sdk",
		Version:       "0.1.0",
		License:       "MIT",
	}

	// Step 1: Check if OpenAPI spec exists
	fmt.Println("📖 Checking OpenAPI specification...")
	if _, err := os.Stat(config.OpenAPIPath); os.IsNotExist(err) {
		fmt.Printf("❌ OpenAPI spec not found at: %s\n", config.OpenAPIPath)
		os.Exit(1)
	}
	fmt.Println("✓ OpenAPI spec found")
	fmt.Println()

	// Step 2: Install hey-api dependencies
	fmt.Println("📦 Installing @hey-api/openapi-ts...")
	if err := installDependencies(); err != nil {
		fmt.Printf("❌ Failed to install dependencies: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Dependencies installed")
	fmt.Println()

	// Step 3: Run hey-api generator
	fmt.Println("⚙️  Running @hey-api/openapi-ts generator...")
	if err := runHeyApiGenerator(); err != nil {
		fmt.Printf("❌ Failed to run hey-api generator: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Base SDK generated")
	fmt.Println()

	// Step 4: Apply custom modifications
	fmt.Println("🔧 Applying custom modifications...")
	if err := applyCustomizations(config); err != nil {
		fmt.Printf("❌ Failed to apply customizations: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Customizations applied")
	fmt.Println()

	// Step 5: Generate additional files
	fmt.Println("📝 Generating additional files...")
	if err := generateAdditionalFiles(config); err != nil {
		fmt.Printf("❌ Failed to generate additional files: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Additional files generated")
	fmt.Println()

	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println("✅ TypeScript SDK generation completed successfully!")
	fmt.Println()
	fmt.Printf("📍 Output location: %s\n", config.OutputPath)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. cd", config.OutputPath)
	fmt.Println("  2. npm install")
	fmt.Println("  3. npm run build")
	fmt.Println("  4. npm test")
	fmt.Println()
}

func installDependencies() error {
	cmd := exec.Command("npm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runHeyApiGenerator() error {
	cmd := exec.Command("npm", "run", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func applyCustomizations(config *GeneratorConfig) error {
	// Custom modification 1: Fix package.json
	if err := customizePackageJson(config); err != nil {
		return fmt.Errorf("failed to customize package.json: %w", err)
	}

	// Custom modification 2: Add custom client wrapper
	if err := addClientWrapper(config); err != nil {
		return fmt.Errorf("failed to add client wrapper: %w", err)
	}

	// Custom modification 3: Fix service exports
	if err := fixServiceExports(config); err != nil {
		return fmt.Errorf("failed to fix service exports: %w", err)
	}

	// Custom modification 4: Unwrap ApiResponse envelope from generated response types
	// Transforms "200: ApiResponse & { data?: T }" → "200: T" so that
	// result.data is T directly (matches the response interceptor in client-wrapper.ts).
	if err := unwrapResponseTypes(config); err != nil {
		return fmt.Errorf("failed to unwrap response types: %w", err)
	}

	return nil
}

func customizePackageJson(config *GeneratorConfig) error {
	pkgPath := filepath.Join(config.OutputPath, "package.json")

	// Read existing package.json
	data, err := ioutil.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Replace package name
	content = strings.ReplaceAll(content, `"name": "insuretech-typescript-sdk"`,
		fmt.Sprintf(`"name": "%s"`, config.PackageName))

	// Replace version
	content = strings.ReplaceAll(content, `"version": "1.0.0"`,
		fmt.Sprintf(`"version": "%s"`, config.Version))

	// Add repository info if not present
	if !strings.Contains(content, `"repository"`) {
		// Insert before devDependencies
		repoInfo := `,
  "repository": {
    "type": "git",
    "url": "https://github.com/lifeplus/InsureTech"
  },
  "bugs": {
    "url": "https://github.com/lifeplus/InsureTech/issues"
  },
  "homepage": "https://github.com/lifeplus/InsureTech#readme"`

		content = strings.ReplaceAll(content, `"devDependencies"`, repoInfo+`,
  "devDependencies"`)
	}

	// Add runtime dependencies if not present
	if !strings.Contains(content, `"dependencies"`) {
		// Insert before devDependencies
		deps := `,
  "dependencies": {
    "@hey-api/client-fetch": "^0.1.0"
  }`

		content = strings.ReplaceAll(content, `"devDependencies"`, deps+`,
  "devDependencies"`)
	}

	return writeFileIfChangedStr(pkgPath, content, 0644)
}

func addClientWrapper(config *GeneratorConfig) error {
	// Create a custom client wrapper that provides better DX
	wrapperPath := filepath.Join(config.OutputPath, "src", "client-wrapper.ts")

	wrapper := `// Custom Client Wrapper for InsureTech SDK
// Provides a configured client instance for use with generated services

import { createClient, createConfig } from './client';

export interface InsureTechClientConfig {
  /** API key for authentication */
  apiKey: string;
  /** Base URL for the API (optional, defaults to production) */
  baseUrl?: string;
  /** Additional headers to include in all requests */
  headers?: Record<string, string>;
}

/**
 * Create a configured client for the InsureTech API
 * 
 * @example
 * ` + "```typescript" + `
 * import { createInsureTechClient, AiService } from '@lifeplus/insuretech-sdk';
 * 
 * const client = createInsureTechClient({
 *   apiKey: 'your-api-key',
 *   baseUrl: 'https://api.insuretech.com'
 * });
 * 
 * // Use with any service method
 * const response = await AiService.aiServiceChat({
 *   client,
 *   body: { message: 'Hello' }
 * });
 * ` + "```" + `
 */
export function createInsureTechClient(config: InsureTechClientConfig) {
  const c = createClient(createConfig({
    baseUrl: config.baseUrl || 'https://api.insuretech.com',
    headers: {
      'Authorization': ` + "`Bearer ${config.apiKey}`" + `,
      ...config.headers,
    },
  }));

  // ── Unwrap ApiResponse envelope ─────────────────────────────────────────
  // The gateway wraps every response as { success, data, error, meta }.
  // hey-api puts the parsed JSON into result.data, so without this
  // interceptor consumers would need result.data.data to reach the payload.
  // By replacing the Response body with just the inner "data" field we make
  // result.data === T directly — no double-wrap.
  c.interceptors.response.use(async (response) => {
    const ct = response.headers.get('content-type') ?? '';
    if (!ct.includes('application/json')) return response;
    // Clone so we can read the body without consuming the original.
    const text = await response.clone().text();
    if (!text) return response;
    try {
      const envelope = JSON.parse(text);
      // Only unwrap if it looks like our standard ApiResponse envelope.
      if (
        typeof envelope === 'object' &&
        envelope !== null &&
        'success' in envelope &&
        'data' in envelope
      ) {
        // Success: unwrap envelope.data so result.data === T
        // Error: unwrap envelope.error so result.error has gateway error details
        const inner = envelope.success ? envelope.data : envelope.error;

        // Preserve Set-Cookie and X-CSRF-Token across the body rewrite.
        // Set-Cookie is a forbidden header in the Fetch API — constructing a
        // new Response(..., { headers }) silently drops it in both browser and
        // Node.js (undici). We copy it to the readable header x-set-cookie so
        // that server-side Next.js API route handlers (e.g. the login route)
        // can still forward the session cookie to the browser.
        const newHeaders = new Headers(response.headers);
        const setCookie = response.headers.get('set-cookie');
        if (setCookie) newHeaders.set('x-set-cookie', setCookie);
        const csrfToken = response.headers.get('x-csrf-token');
        if (csrfToken) newHeaders.set('x-csrf-token', csrfToken);

        return new Response(JSON.stringify(inner ?? {}), {
          status: response.status,
          statusText: response.statusText,
          headers: newHeaders,
        });
      }
    } catch { /* not JSON — pass through */ }
    return response;
  });

  return c;
}

// Re-export for convenience
export { createClient, createConfig } from './client';
`

	return writeFileIfChangedStr(wrapperPath, wrapper, 0644)
}

func fixServiceExports(config *GeneratorConfig) error {
	// hey-api v0.73+ generates sdk.gen.ts and types.gen.ts
	// index.ts re-exports everything including the unified ApiResponse<T> envelope.
	indexPath := filepath.Join(config.OutputPath, "src", "index.ts")

	index := `// Auto-generated SDK Entry Point — DO NOT EDIT
// Generated by InsureTech API Pipeline

// ─── Core ApiResponse<T> envelope + type guards ───────────────────────────────
export type {
  ApiResponse,
  ResponseMeta,
  PaginationMeta,
  PaginationRequest,
  Money,
  Address,
  Timestamp,
  DateString,
  UUID,
} from './types';
export { unwrapData, isApiSuccess, isApiError } from './types';

// ─── Structured error classes ─────────────────────────────────────────────────
export { InsureTechApiError, ApiError } from './errors';
export type { ApiErrorDetail, FieldViolation } from './errors';

// ─── Generated services and types (hey-api) ───────────────────────────────────
export * from './sdk.gen';
export * from './types.gen';

// ─── Custom client helper ─────────────────────────────────────────────────────
export { createInsureTechClient } from './client-wrapper';
export type { InsureTechClientConfig } from './client-wrapper';
`

	return writeFileIfChangedStr(indexPath, index, 0644)
}

func generateAdditionalFiles(config *GeneratorConfig) error {
	// Generate ApiResponse envelope types (errors.ts, types.ts)
	if err := generateErrorsFile(config); err != nil {
		return fmt.Errorf("failed to generate errors.ts: %w", err)
	}

	if err := generateTypesFile(config); err != nil {
		return fmt.Errorf("failed to generate types.ts: %w", err)
	}

	// Generate README
	if err := generateReadme(config); err != nil {
		return err
	}

	// Generate vitest config
	if err := generateVitestConfig(config); err != nil {
		return err
	}

	// Generate prettier config
	if err := generatePrettierConfig(config); err != nil {
		return err
	}

	// Generate tsconfig
	if err := generateTsConfig(config); err != nil {
		return err
	}

	return nil
}

// generateErrorsFile writes src/errors.ts from errors.ts.tmpl.
// Contains InsureTechApiError class and ApiErrorDetail/FieldViolation interfaces
// aligned with the Go gateway's respond package.
func generateErrorsFile(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "errors.ts.tmpl")
	outputPath := filepath.Join(config.OutputPath, "src", "errors.ts")
	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}
	return writeFileIfChanged(outputPath, data, 0644)
}

// generateTypesFile writes src/types.ts from types.ts.tmpl.
// Contains ApiResponse<T>, ResponseMeta, PaginationMeta and helper utilities
// (unwrapData, isApiSuccess, isApiError) aligned with the openapi.yaml schema.
func generateTypesFile(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "types.ts.tmpl")
	outputPath := filepath.Join(config.OutputPath, "src", "types.ts")
	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}
	return writeFileIfChanged(outputPath, data, 0644)
}

func generateReadme(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "README.md.tmpl")
	outputPath := filepath.Join(config.OutputPath, "README.md")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"PackageName": config.PackageName,
		"Version":     config.Version,
		"License":     config.License,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	return writeFileIfChanged(outputPath, buf.Bytes(), 0644)
}

func generateVitestConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "vitest.config.ts.tmpl")
	outputPath := filepath.Join(config.OutputPath, "vitest.config.ts")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return writeFileIfChanged(outputPath, data, 0644)
}

func generatePrettierConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, ".prettierrc.tmpl")
	outputPath := filepath.Join(config.OutputPath, ".prettierrc")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return writeFileIfChanged(outputPath, data, 0644)
}

func generateTsConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "tsconfig.json.tmpl")
	outputPath := filepath.Join(config.OutputPath, "tsconfig.json")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return writeFileIfChanged(outputPath, data, 0644)
}

// unwrapResponseTypes post-processes the generated types.gen.ts to remove the
// ApiResponse envelope wrapper from response type definitions.
//
// hey-api generates response types like:
//
//	200: ApiResponse & {
//	    data?: LoginResponse;
//	};
//
// Because the client-wrapper response interceptor already strips the envelope
// at the HTTP layer, the TypeScript types should reflect the unwrapped payload:
//
//	200: LoginResponse;
//
// This function performs a regex transformation on the generated types file
// to align the types with the runtime unwrapping behaviour.
func unwrapResponseTypes(config *GeneratorConfig) error {
	typesPath := filepath.Join(config.OutputPath, "src", "types.gen.ts")
	data, err := ioutil.ReadFile(typesPath)
	if err != nil {
		return err
	}

	content := string(data)

	// Pattern: "    200: ApiResponse & {\n        data?: SomeType;\n    };"
	// Replace with: "    200: SomeType;"
	// Works for 200 and 201 success responses.
	re := regexp.MustCompile(`(\s+)(200|201): ApiResponse & \{\s*\n\s+data\?: ([^;]+);\s*\n\s+\};`)
	content = re.ReplaceAllString(content, "${1}${2}: ${3};")

	return writeFileIfChanged(typesPath, []byte(content), 0644)
}
