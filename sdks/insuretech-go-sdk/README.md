# InsureTech Go SDK

Official Go SDK for the InsureTech API Platform.

## Installation

```bash
go get github.com/newage-saint/insuretech-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    
    insuretech "github.com/newage-saint/insuretech-go-sdk"
)

func main() {
    client := insuretech.NewClient(
        insuretech.WithAPIKey("your-api-key"),
    )
    
    ctx := context.Background()
    
    // Use the client...
}
```

## Documentation

For full documentation, see [QUICKSTART.md](../sdk-generator/go/QUICKSTART.md)

## Version

Version: 1.0.0

## License

MIT
