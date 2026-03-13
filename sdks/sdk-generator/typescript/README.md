# TypeScript SDK Generator

This directory contains the TypeScript SDK generator that uses `@hey-api/openapi-ts` with custom Go-based post-processing.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    OpenAPI Spec (YAML)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              @hey-api/openapi-ts Generator                  │
│  - Generates types, services, and client                    │
│  - Modern, actively maintained                              │
│  - Supports latest OpenAPI 3.1 features                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│           Custom Go Post-Processor (generator.go)           │
│  - Customizes package.json                                  │
│  - Adds client wrapper for better DX                        │
│  - Fixes exports                                            │
│  - Generates additional files (README, configs)             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              Production-Ready TypeScript SDK                │
│  - Fully typed                                              │
│  - Tree-shakeable                                           │
│  - ESM + CJS support                                        │
│  - Comprehensive documentation                              │
└─────────────────────────────────────────────────────────────┘
```

## Why @hey-api/openapi-ts?

1. **Modern & Maintained**: Actively developed, supports latest OpenAPI 3.1
2. **Better Output**: Generates cleaner, more idiomatic TypeScript
3. **Flexible**: Highly configurable with plugins
4. **Type-Safe**: Excellent TypeScript support with proper generics
5. **Tree-Shakeable**: Generates code that can be tree-shaken
6. **Multiple Clients**: Supports fetch, axios, and custom clients

## Files

- `openapi-ts.config.ts` - Configuration for @hey-api/openapi-ts
- `package.json` - Dependencies for the generator
- `generator.go` - Custom Go post-processor
- `generate.ps1` - Main generation script
- `templates/` - Templates for additional files (README, configs)

## Usage

### Quick Start

```powershell
# From this directory
.\generate.ps1
```

This will:
1. Install @hey-api/openapi-ts dependencies
2. Run the hey-api generator
3. Build and run custom Go post-processor
4. Install SDK dependencies
5. Build the SDK
6. Run tests

### Manual Steps

```powershell
# 1. Install generator dependencies
npm install

# 2. Run hey-api generator
npm run generate

# 3. Build custom post-processor
$env:GOWORK="off"
go build -o generator.exe generator.go

# 4. Run post-processor
.\generator.exe

# 5. Build SDK
cd ../../insuretech-typescript-sdk
npm install
npm run build
```

## Configuration

### hey-api Configuration (`openapi-ts.config.ts`)

```typescript
export default defineConfig({
  client: '@hey-api/client-fetch',  // Use fetch client
  input: '../../../api/openapi.yaml',  // OpenAPI spec
  output: {
    path: '../../insuretech-typescript-sdk/src',
    format: 'prettier',  // Auto-format output
    lint: 'eslint',  // Auto-lint output
  },
  types: {
    enums: 'javascript',  // Generate JS enums
    dates: 'types+transform',  // Transform dates
  },
  services: {
    asClass: true,  // Generate services as classes
  },
});
```

### Custom Post-Processing

The Go post-processor (`generator_v2.go`) performs:

1. **Package.json Customization**
   - Sets correct package name (`@lifeplus/insuretech-sdk`)
   - Sets version (`0.1.0`)
   - Adds repository information

2. **Client Wrapper**
   - Creates `InsureTechClient` class
   - Simplifies authentication
   - Provides better developer experience

3. **Export Fixes**
   - Ensures proper module exports
   - Re-exports all generated types

4. **Additional Files**
   - Generates README.md
   - Generates vitest.config.ts
   - Generates .prettierrc

## Customization

### Adding Custom Modifications

Edit `generator.go` to add custom post-processing:

```go
func applyCustomizations(config *GeneratorConfig) error {
    // Add your custom modifications here
    
    // Example: Add custom utility functions
    if err := addUtilityFunctions(config); err != nil {
        return err
    }
    
    return nil
}

func addUtilityFunctions(config *GeneratorConfig) error {
    utilsPath := filepath.Join(config.OutputPath, "src", "utils.ts")
    
    utils := `// Custom utility functions
export function formatMoney(amount: number, currency: string): string {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency,
    }).format(amount / 100);
}
`
    
    return ioutil.WriteFile(utilsPath, []byte(utils), 0644)
}
```

### Modifying hey-api Configuration

Edit `openapi-ts.config.ts`:

```typescript
export default defineConfig({
  // Change client library
  client: '@hey-api/client-axios',  // Use axios instead of fetch
  
  // Add plugins
  plugins: [
    '@hey-api/typescript',
    '@hey-api/schemas',
    {
      name: '@hey-api/sdk',
      asClass: true,
    },
  ],
  
  // Customize output
  output: {
    path: '../../insuretech-typescript-sdk/src',
    format: 'prettier',
    lint: 'eslint',
  },
});
```

## Generated SDK Structure

```
insuretech-typescript-sdk/
├── src/
│   ├── index.ts              # Main entry (custom)
│   ├── client.ts             # Client wrapper (custom)
│   ├── sdk/                  # Generated by hey-api
│   │   ├── types.ts          # All TypeScript types
│   │   ├── services/         # Service classes
│   │   │   ├── AuthService.ts
│   │   │   ├── PoliciesService.ts
│   │   │   └── ...
│   │   └── index.ts          # SDK exports
│   └── utils.ts              # Custom utilities (optional)
├── package.json              # Customized by post-processor
├── tsconfig.json             # TypeScript config
├── vitest.config.ts          # Test config (generated)
├── .prettierrc               # Prettier config (generated)
└── README.md                 # Documentation (generated)
```

## Advantages Over Custom Generator

### Before (Custom Generator)
- ❌ Manual template maintenance
- ❌ Complex parsing logic
- ❌ Incomplete type generation
- ❌ Missing edge cases
- ❌ Hard to update for new OpenAPI features

### After (hey-api + Custom)
- ✅ Maintained by community
- ✅ Handles all OpenAPI 3.1 features
- ✅ Complete type generation
- ✅ Handles edge cases
- ✅ Easy to update
- ✅ Custom post-processing for specific needs

## Troubleshooting

### hey-api fails to generate

```powershell
# Check OpenAPI spec is valid
npm run generate -- --dry-run

# Check for syntax errors
npx @hey-api/openapi-ts --help
```

### Post-processor fails

```powershell
# Build with verbose output
go build -v -o generator.exe generator.go

# Run with error details
.\generator.exe
```

### SDK build fails

```powershell
cd ../../insuretech-typescript-sdk

# Check for type errors
npm run typecheck

# Check for lint errors
npm run lint

# Build with verbose output
npm run build -- --verbose
```

## Migration from Old Generator

If migrating from the old custom generator:

1. **Backup existing SDK**
   ```powershell
   Copy-Item -Recurse ../../insuretech-typescript-sdk ../../insuretech-typescript-sdk.backup
   ```

2. **Run new generator**
   ```powershell
   .\generate.ps1
   ```

3. **Compare outputs**
   ```powershell
   # Use a diff tool to compare
   code --diff ../../insuretech-typescript-sdk.backup ../../insuretech-typescript-sdk
   ```

4. **Test thoroughly**
   ```powershell
   cd ../../insuretech-typescript-sdk
   npm test
   npm run build
   ```

## Resources

- [@hey-api/openapi-ts Documentation](https://heyapi.vercel.app/)
- [OpenAPI 3.1 Specification](https://spec.openapis.org/oas/v3.1.0)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

## Support

For issues with:
- **hey-api generator**: Check [hey-api GitHub](https://github.com/hey-api/openapi-ts)
- **Custom post-processor**: Check `generator.go` code
- **Generated SDK**: Check SDK's README.md
