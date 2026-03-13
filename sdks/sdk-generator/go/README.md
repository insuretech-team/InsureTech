# Go SDK Generator

This directory contains the Go SDK generator for the InsureTech API.

## Overview

The Go SDK generator creates a client library from the OpenAPI specification located in `C:\_DEV\GO\InsureTech\api\openapi.yaml`, using the Protocol Buffer definitions in `C:\_DEV\GO\InsureTech\proto` as the source of truth.

## Structure

```
go/
├── templates/           # Go SDK templates
│   ├── client.go.tmpl  # Client wrapper template
│   ├── config.go.tmpl  # Configuration template
│   ├── models.go.tmpl  # Model modifiers template
│   └── errors.go.tmpl  # Error handling template
├── generator.go        # Main generator logic
└── README.md          # This file
```

## Usage

```bash
# Generate the Go SDK
go run generator.go

# Output will be in: C:\_DEV\GO\InsureTech\sdks\go-sdk
```

## Features

- Auto-generated client from OpenAPI spec
- Type-safe API methods
- Built-in retry logic
- Error handling
- Authentication support
- Context-aware requests

## Requirements

- Go 1.21 or higher
- OpenAPI Generator or custom generator
- Access to proto definitions

## Configuration

The generator uses:
- **Source**: `C:\_DEV\GO\InsureTech\api\openapi.yaml`
- **Proto Source**: `C:\_DEV\GO\InsureTech\proto`
- **Output**: `C:\_DEV\GO\InsureTech\sdks\go-sdk`
