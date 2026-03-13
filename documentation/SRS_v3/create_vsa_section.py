"""
Create new VSA Architecture section to insert into document
"""

vsa_section = """
[[[PAGEBREAK]]]

## 3.5 System Architecture - VSA Pattern

### 3.5.1 Vertical Slice Architecture Overview

![VSA Architecture](VSA.png)

*Figure 1: Vertical Slice Architecture - Language-Agnostic Pattern*

The LabAid InsureTech Platform adopts **Vertical Slice Architecture (VSA)** across ALL microservices, regardless of programming language:

- **Go Services:** Gateway, Auth, DBManager, Storage, IoT Broker, Kafka Orchestration
- **C# .NET Services:** Insurance Engine, Partner Management, Analytics & Reporting  
- **Node.js Services:** Payment Service, Ticketing Service
- **Python Services:** AI Engine, OCR/PDF Service

**Key VSA Principles:**
1. **High Cohesion:** Each slice contains all layers needed for one feature
2. **Low Coupling:** Slices are independent and don't share logic
3. **Feature-Focused:** Organized by business capability, not technical layer
4. **Testability:** Each slice can be tested in isolation

### 3.5.2 Protocol Buffer Data Models

All data models are defined in **Protocol Buffers (Proto3)** for type-safe, language-agnostic communication.

**Proto Structure:**
```
proto/
└── insuretech/
    └── v1/
        ├── entities/          # Core data models
        │   ├── user.proto
        │   ├── policy.proto
        │   ├── claim.proto
        │   ├── payment.proto
        │   └── ...
        ├── events/            # Domain events
        │   ├── policy_events.proto
        │   ├── claim_events.proto
        │   └── ...
        ├── services/          # gRPC service definitions
        │   ├── insurance_engine.proto
        │   ├── partner_management.proto
        │   └── ...
        └── common/            # Shared types
            ├── localized_string.proto
            ├── pagination.proto
            └── error.proto
```

**Proto Package Convention:**
```protobuf
syntax = "proto3";

package insuretech.v1.entities;

option csharp_namespace = "Insuretech.V1.Entities";
option go_package = "insuretech/v1/entities";
option java_package = "com.insuretech.v1.entities";
option java_multiple_files = true;
```

### 3.5.3 Insurance Engine - CQRS with MediatR

The Insurance Engine (C# .NET 8) implements **CQRS** (Command Query Responsibility Segregation) using **MediatR** for clean separation:

**Technology Stack:**
- .NET 8 LTS
- MediatR 12.x (CQRS pattern)
- FluentValidation 11.x (Request validation)
- Entity Framework Core 8.x (PostgreSQL ORM)
- Grpc.AspNetCore 2.x
- Serilog (Structured logging)

**Folder Structure:**
```
InsuranceEngine/
├── API/Grpc/
│   └── Services/
│       └── InsuranceEngineGrpcService.cs
├── Application/
│   ├──Commands/              # Write operations (CQRS)
│   │   ├── Policies/IssuePolicy/
│   │   │   ├── IssuePolicyCommand.cs
│   │   │   ├── IssuePolicyCommandHandler.cs
│   │   │   └── IssuePolicyCommandValidator.cs
│   │   └── Claims/SubmitClaim/
│   ├── Queries/               # Read operations (CQRS)
│   │   ├── Policies/GetPolicyById/
│   │   │   ├── GetPolicyByIdQuery.cs
│   │   │   └── GetPolicyByIdQueryHandler.cs
│   │   └── Premium/CalculatePremium/
│   └── Behaviors/             # MediatR Pipeline
│       ├── ValidationBehavior.cs
│       ├── LoggingBehavior.cs
│       └── PerformanceBehavior.cs
├── Domain/
│   ├── Entities/
│   ├── ValueObjects/
│   └── Events/
└── Infrastructure/
    ├── Persistence/
    ├── GrpcClients/           # Adapters to other services
    └── EventBus/
```

**CQRS Command Example:**
```csharp
// Command
public class IssuePolicyCommand : IRequest<IssuePolicyResponse>
{
    public string UserId { get; set; }
    public string ProductId { get; set; }
    public decimal SumAssured { get; set; }
}

// Handler
public class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, IssuePolicyResponse>
{
    private readonly IPolicyRepository _policyRepository;
    private readonly IEventPublisher _eventPublisher;

    public async Task<IssuePolicyResponse> Handle(IssuePolicyCommand request, CancellationToken ct)
    {
        // Create policy
        var policy = Policy.Create(request.UserId, request.ProductId, request.SumAssured);
        
        // Persist
        await _policyRepository.AddAsync(policy, ct);
        
        // Publish event
        await _eventPublisher.PublishAsync(new PolicyIssuedEvent(policy.Id), ct);
        
        return new IssuePolicyResponse { PolicyId = policy.Id };
    }
}

// Validator
public class IssuePolicyCommandValidator : AbstractValidator<IssuePolicyCommand>
{
    public IssuePolicyCommandValidator()
    {
        RuleFor(x => x.SumAssured)
            .GreaterThan(0).WithMessage("Sum assured must be greater than 0")
            .LessThanOrEqualTo(10000000).WithMessage("Cannot exceed 1 Crore BDT");
    }
}
```

**CQRS Query Example:**
```csharp
// Query
public class GetPolicyByIdQuery : IRequest<PolicyDto>
{
    public string PolicyId { get; set; }
}

// Handler
public class GetPolicyByIdQueryHandler : IRequestHandler<GetPolicyByIdQuery, PolicyDto>
{
    private readonly IPolicyRepository _repository;
    private readonly IMapper _mapper;

    public async Task<PolicyDto> Handle(GetPolicyByIdQuery request, CancellationToken ct)
    {
        var policy = await _repository.GetByIdAsync(request.PolicyId, ct);
        return _mapper.Map<PolicyDto>(policy);
    }
}
```

**MediatR Pipeline Behaviors:**
```csharp
// Automatic validation for all commands/queries
public class ValidationBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
{
    private readonly IEnumerable<IValidator<TRequest>> _validators;

    public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken ct)
    {
        var failures = _validators
            .SelectMany(v => v.Validate(request).Errors)
            .Where(f => f != null)
            .ToList();

        if (failures.Any())
            throw new ValidationException(failures);

        return await next();
    }
}
```

---
"""

print("VSA section created")
print(f"Length: {len(vsa_section)} characters")
with open("vsa_architecture_section.txt", "w", encoding="utf-8") as f:
    f.write(vsa_section)
print("Saved to vsa_architecture_section.txt")
