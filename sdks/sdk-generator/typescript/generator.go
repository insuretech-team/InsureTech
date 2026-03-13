package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

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

	return ioutil.WriteFile(pkgPath, []byte(content), 0644)
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
  return createClient(createConfig({
    baseUrl: config.baseUrl || 'https://api.insuretech.com',
    headers: {
      'Authorization': ` + "`Bearer ${config.apiKey}`" + `,
      ...config.headers,
    },
  }));
}

// Re-export for convenience
export { createClient, createConfig } from './client';
`

	return ioutil.WriteFile(wrapperPath, []byte(wrapper), 0644)
}

func fixServiceExports(config *GeneratorConfig) error {
	// hey-api v0.73+ generates sdk.gen.ts and types.gen.ts
	// We just need to re-export them along with our custom client helper
	indexPath := filepath.Join(config.OutputPath, "src", "index.ts")
	
	index := `// Main SDK Entry Point
// Auto-generated with custom enhancements

// Export all generated services and types
export * from './sdk.gen';
export * from './types.gen';

// Export custom client helper
export { createInsureTechClient } from './client-wrapper';
export type { InsureTechClientConfig } from './client-wrapper';
`

	return ioutil.WriteFile(indexPath, []byte(index), 0644)
}

func generateAdditionalFiles(config *GeneratorConfig) error {
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

func generateReadme(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "README.md.tmpl")
	outputPath := filepath.Join(config.OutputPath, "README.md")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[string]interface{}{
		"PackageName": config.PackageName,
		"Version":     config.Version,
		"License":     config.License,
	}

	return tmpl.Execute(f, data)
}

func generateVitestConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "vitest.config.ts.tmpl")
	outputPath := filepath.Join(config.OutputPath, "vitest.config.ts")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}

func generatePrettierConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, ".prettierrc.tmpl")
	outputPath := filepath.Join(config.OutputPath, ".prettierrc")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}

func generateTsConfig(config *GeneratorConfig) error {
	tmplPath := filepath.Join(config.TemplatesPath, "tsconfig.json.tmpl")
	outputPath := filepath.Join(config.OutputPath, "tsconfig.json")

	data, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}
