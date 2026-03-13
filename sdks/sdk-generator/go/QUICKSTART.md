# Go SDK Quick Start Guide

## Installation

```bash
go get github.com/insuretech/go-sdk
```

## Basic Usage

### 1. Initialize the Client

```go
package main

import (
    "context"
    "log"
    
    insuretech "github.com/insuretech/go-sdk"
)

func main() {
    // Create a new client with API key authentication
    client := insuretech.NewClient(
        insuretech.WithAPIKey("your-api-key"),
        insuretech.WithBaseURL("https://api.insuretech.com"),
    )
    
    ctx := context.Background()
    
    // Use the client...
}
```

### 2. Authentication

#### API Key Authentication
```go
client := insuretech.NewClient(
    insuretech.WithAPIKey("your-api-key"),
)
```

#### OAuth2 Authentication
```go
client := insuretech.NewClient(
    insuretech.WithOAuth2("your-access-token"),
)
```

### 3. Making API Calls

#### Create a Policy
```go
policy, err := client.Policies.Create(ctx, &insuretech.CreatePolicyRequest{
    ProductID: "prod_123",
    CustomerID: "cust_456",
    StartDate: time.Now(),
    EndDate: time.Now().AddDate(1, 0, 0),
    Premium: 1000.00,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created policy: %s\n", policy.ID)
```

#### List Policies
```go
policies, err := client.Policies.List(ctx, &insuretech.ListPoliciesRequest{
    Page: 1,
    PageSize: 20,
    Status: "active",
})
if err != nil {
    log.Fatal(err)
}
for _, policy := range policies.Items {
    fmt.Printf("Policy: %s - %s\n", policy.ID, policy.Status)
}
```

#### Get a Policy
```go
policy, err := client.Policies.Get(ctx, "policy_123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Policy: %+v\n", policy)
```

### 4. Error Handling

```go
policy, err := client.Policies.Get(ctx, "policy_123")
if err != nil {
    if apiErr, ok := err.(*insuretech.APIError); ok {
        fmt.Printf("API Error: %s (Code: %d)\n", apiErr.Message, apiErr.StatusCode)
        switch apiErr.StatusCode {
        case 404:
            fmt.Println("Policy not found")
        case 401:
            fmt.Println("Unauthorized")
        default:
            fmt.Printf("Unexpected error: %v\n", apiErr)
        }
    } else {
        log.Fatal(err)
    }
}
```

### 5. Working with Claims

```go
// Create a claim
claim, err := client.Claims.Create(ctx, &insuretech.CreateClaimRequest{
    PolicyID: "policy_123",
    ClaimType: "accident",
    ClaimDate: time.Now(),
    Amount: 5000.00,
    Description: "Vehicle accident claim",
})
if err != nil {
    log.Fatal(err)
}

// Update claim status
claim, err = client.Claims.Update(ctx, claim.ID, &insuretech.UpdateClaimRequest{
    Status: "approved",
    ApprovedAmount: 4500.00,
})
```

### 6. Pagination

```go
page := 1
pageSize := 50

for {
    response, err := client.Policies.List(ctx, &insuretech.ListPoliciesRequest{
        Page: page,
        PageSize: pageSize,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Process policies
    for _, policy := range response.Items {
        fmt.Printf("Policy: %s\n", policy.ID)
    }
    
    // Check if there are more pages
    if !response.HasMore {
        break
    }
    page++
}
```

### 7. Context and Timeouts

```go
// Create a context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

policy, err := client.Policies.Get(ctx, "policy_123")
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    } else {
        log.Fatal(err)
    }
}
```

### 8. Advanced Configuration

```go
client := insuretech.NewClient(
    insuretech.WithAPIKey("your-api-key"),
    insuretech.WithBaseURL("https://api.insuretech.com"),
    insuretech.WithTimeout(30*time.Second),
    insuretech.WithRetry(3, time.Second),
    insuretech.WithHTTPClient(&http.Client{
        Transport: &http.Transport{
            MaxIdleConns: 10,
            IdleConnTimeout: 90 * time.Second,
        },
    }),
)
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    insuretech "github.com/insuretech/go-sdk"
)

func main() {
    // Initialize client
    client := insuretech.NewClient(
        insuretech.WithAPIKey("your-api-key"),
        insuretech.WithBaseURL("https://api.insuretech.com"),
    )
    
    ctx := context.Background()
    
    // Create a policy
    policy, err := client.Policies.Create(ctx, &insuretech.CreatePolicyRequest{
        ProductID: "prod_health_001",
        CustomerID: "cust_12345",
        StartDate: time.Now(),
        EndDate: time.Now().AddDate(1, 0, 0),
        Premium: 1200.00,
        Currency: "USD",
    })
    if err != nil {
        log.Fatalf("Failed to create policy: %v", err)
    }
    
    fmt.Printf("✓ Created policy: %s\n", policy.ID)
    
    // Create a claim
    claim, err := client.Claims.Create(ctx, &insuretech.CreateClaimRequest{
        PolicyID: policy.ID,
        ClaimType: "medical",
        ClaimDate: time.Now(),
        Amount: 500.00,
        Description: "Medical consultation",
    })
    if err != nil {
        log.Fatalf("Failed to create claim: %v", err)
    }
    
    fmt.Printf("✓ Created claim: %s\n", claim.ID)
    
    // List all claims for the policy
    claims, err := client.Claims.List(ctx, &insuretech.ListClaimsRequest{
        PolicyID: policy.ID,
        Page: 1,
        PageSize: 10,
    })
    if err != nil {
        log.Fatalf("Failed to list claims: %v", err)
    }
    
    fmt.Printf("✓ Found %d claims\n", len(claims.Items))
}
```

## Resources

- [Full API Documentation](https://docs.insuretech.com)
- [GitHub Repository](https://github.com/insuretech/go-sdk)
- [Examples](https://github.com/insuretech/go-sdk/tree/main/examples)
- [Support](https://support.insuretech.com)
