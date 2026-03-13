"""
Phase 2: Update proto structure to proto/Insuretech/v1/
and add detailed CQRS/MediatR Insurance Engine structure
"""

import re

# Read the current file
with open("SRS_V3_FINAL_DRAFT.md", "r", encoding="utf-8") as f:
    content = f.read()

# 4. Update all proto syntax declarations to use proper structure
old_proto_syntax = 'syntax = "proto3";\n\npackage labaid.entities;'
new_proto_syntax = '''syntax = "proto3";

package insuretech.v1.entities;

option csharp_namespace = "Insuretech.V1.Entities";
option go_package = "insuretech/v1/entities";
option java_package = "com.insuretech.v1.entities";
option java_multiple_files = true;'''

content = content.replace(old_proto_syntax, new_proto_syntax)

# Also update other proto packages
content = content.replace('package labaid.entities', 'package insuretech.v1.entities')
content = content.replace('package labaid.insurance', 'package insuretech.v1.services')
content = content.replace('option csharp_namespace = "LabAid.Insurance.Grpc";', 'option csharp_namespace = "Insuretech.V1.Services";')
content = content.replace('option csharp_namespace = "LabAid.Entities";', 'option csharp_namespace = "Insuretech.V1.Entities";')
content = content.replace('option go_package = "labaid/entities";', 'option go_package = "insuretech/v1/entities";')

# 5. Update proto file organization structure
old_proto_structure = """### 6.4 Proto File Organization

```
proto/
├── entities/
│   ├── user.proto
│   ├── policy.proto
│   ├── claim.proto
│   ├── payment.proto
│   ├── partner.proto
│   ├── product.proto
│   ├── quote.proto
│   ├── kyc_document.proto
│   ├── audit_log.proto
│   └── notification_log.proto
├── services/
│   ├── insurance_engine.proto
│   ├── partner_management.proto
│   ├── ai_engine.proto
│   ├── payment_service.proto
│   ├── notification_service.proto
│   └── analytics_service.proto
└── common/
    ├── localized_string.proto
    ├── pagination.proto
    └── error.proto
```"""

new_proto_structure = """### 6.4 Proto File Organization

```
proto/
├── insuretech/
│   └── v1/
│       ├── entities/
│       │   ├── user.proto
│       │   ├── policy.proto
│       │   ├── claim.proto
│       │   ├── payment.proto
│       │   ├── partner.proto
│       │   ├── product.proto
│       │   ├── quote.proto
│       │   ├── kyc_document.proto
│       │   ├── audit_log.proto
│       │   └── notification_log.proto
│       ├── events/
│       │   ├── policy_events.proto
│       │   ├── claim_events.proto
│       │   ├── payment_events.proto
│       │   └── user_events.proto
│       ├── services/
│       │   ├── insurance_engine.proto
│       │   ├── partner_management.proto
│       │   ├── ai_engine.proto
│       │   ├── payment_service.proto
│       │   ├── notification_service.proto
│       │   └── analytics_service.proto
│       └── common/
│           ├── localized_string.proto
│           ├── pagination.proto
│           ├── error.proto
│           └── metadata.proto
```

**Versioning Strategy:**
- `v1` = Current production API
- `v2` = Breaking changes (future)
- Maintain backward compatibility within same version"""

content = content.replace(old_proto_structure, new_proto_structure)

# 6. Add detailed CQRS/MediatR Insurance Engine structure
detailed_insurance_engine = """
[[[PAGEBREAK]]]

### 4.5 Insurance Engine - Detailed CQRS Architecture

The Insurance Engine is the core domain service built with **C# .NET 8**, implementing **CQRS (Command Query Responsibility Segregation)** with **MediatR** for clean separation of concerns.

**Technology Stack:**
- **.NET 8** - LTS framework
- **MediatR 12.x** - In-process messaging and CQRS
- **FluentValidation 11.x** - Request validation
- **AutoMapper 12.x** - Object-to-object mapping
- **Entity Framework Core 8.x** - ORM for PostgreSQL
- **Grpc.AspNetCore 2.x** - gRPC server
- **Polly 8.x** - Resilience and transient fault handling
- **Serilog** - Structured logging

#### 4.5.1 CQRS Pattern Implementation

```
InsuranceEngine/
├── API/
│   ├── Grpc/
│   │   ├── Services/
│   │   │   ├── InsuranceEngineGrpcService.cs
│   │   │   ├── PolicyGrpcService.cs
│   │   │   └── PremiumGrpcService.cs
│   │   └── Interceptors/
│   │       ├── LoggingInterceptor.cs
│   │       └── ExceptionInterceptor.cs
│   └── Program.cs
├── Application/
│   ├── Commands/               # Write operations (CQRS)
│   │   ├── Policies/
│   │   │   ├── IssuePolicy/
│   │   │   │   ├── IssuePolicyCommand.cs
│   │   │   │   ├── IssuePolicyCommandHandler.cs
│   │   │   │   ├── IssuePolicyCommandValidator.cs
│   │   │   │   └── IssuePolicyCommandTests.cs
│   │   │   ├── CancelPolicy/
│   │   │   │   ├── CancelPolicyCommand.cs
│   │   │   │   └── CancelPolicyCommandHandler.cs
│   │   │   └── RenewPolicy/
│   │   │       ├── RenewPolicyCommand.cs
│   │   │       └── RenewPolicyCommandHandler.cs
│   │   ├── Quotes/
│   │   │   └── CreateQuote/
│   │   │       ├── CreateQuoteCommand.cs
│   │   │       └── CreateQuoteCommandHandler.cs
│   │   └── Claims/
│   │       ├── SubmitClaim/
│   │       │   ├── SubmitClaimCommand.cs
│   │       │   └── SubmitClaimCommandHandler.cs
│   │       └── ApproveClaim/
│   │           ├── ApproveClaimCommand.cs
│   │           └── ApproveClaimCommandHandler.cs
│   ├── Queries/                # Read operations (CQRS)
│   │   ├── Policies/
│   │   │   ├── GetPolicyById/
│   │   │   │   ├── GetPolicyByIdQuery.cs
│   │   │   │   └── GetPolicyByIdQueryHandler.cs
│   │   │   ├── GetPoliciesByUser/
│   │   │   │   ├── GetPoliciesByUserQuery.cs
│   │   │   │   └── GetPoliciesByUserQueryHandler.cs
│   │   │   └── SearchPolicies/
│   │   │       ├── SearchPoliciesQuery.cs
│   │   │       └── SearchPoliciesQueryHandler.cs
│   │   ├── Premium/
│   │   │   └── CalculatePremium/
│   │   │       ├── CalculatePremiumQuery.cs
│   │   │       └── CalculatePremiumQueryHandler.cs
│   │   └── Claims/
│   │       └── GetClaimStatus/
│   │           ├── GetClaimStatusQuery.cs
│   │           └── GetClaimStatusQueryHandler.cs
│   ├── Behaviors/              # MediatR Pipeline Behaviors
│   │   ├── ValidationBehavior.cs
│   │   ├── LoggingBehavior.cs
│   │   ├── PerformanceBehavior.cs
│   │   └── TransactionBehavior.cs
│   ├── DTOs/
│   │   ├── PolicyDto.cs
│   │   ├── QuoteDto.cs
│   │   └── ClaimDto.cs
│   └── Mappings/
│       └── AutoMapperProfile.cs
├── Domain/
│   ├── Entities/
│   │   ├── Policy.cs
│   │   ├── Quote.cs
│   │   ├── Claim.cs
│   │   ├── Premium.cs
│   │   └── Nominee.cs
│   ├── ValueObjects/
│   │   ├── Money.cs
│   │   ├── PolicyNumber.cs
│   │   └── DateRange.cs
│   ├── Events/                 # Domain Events
│   │   ├── PolicyIssuedEvent.cs
│   │   ├── PolicyCancelledEvent.cs
│   │   ├── ClaimSubmittedEvent.cs
│   │   └── ClaimApprovedEvent.cs
│   ├── Interfaces/
│   │   ├── IPolicyRepository.cs
│   │   ├── IQuoteRepository.cs
│   │   └── IClaimRepository.cs
│   └── Services/               # Domain Services
│       ├── PremiumCalculationService.cs
│       ├── UnderwritingService.cs
│       └── RiskAssessmentService.cs
├── Infrastructure/
│   ├── Persistence/
│   │   ├── InsuranceDbContext.cs
│   │   ├── Configurations/
│   │   │   ├── PolicyConfiguration.cs
│   │   │   └── ClaimConfiguration.cs
│   │   └── Repositories/
│   │       ├── PolicyRepository.cs
│   │       ├── QuoteRepository.cs
│   │       └── ClaimRepository.cs
│   ├── GrpcClients/           # Adapters to other services
│   │   ├── PartnerManagementClient.cs
│   │   ├── AIEngineClient.cs
│   │   ├── PaymentServiceClient.cs
│   │   └── NotificationServiceClient.cs
│   └── EventBus/
│       └── KafkaEventPublisher.cs
└── Proto/
    └── insuretech/
        └── v1/
            ├── services/
            │   └── insurance_engine.proto
            └── entities/
                ├── policy.proto
                ├── quote.proto
                └── claim.proto
```

#### 4.5.2 MediatR Request/Response Flow

```csharp
// Command Example: Issue Policy
public class IssuePolicyCommand : IRequest<IssuePolicyResponse>
{
    public string UserId { get; set; }
    public string ProductId { get; set; }
    public decimal SumAssured { get; set; }
    public int TenureYears { get; set; }
    public List<NomineeDto> Nominees { get; set; }
}

public class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, IssuePolicyResponse>
{
    private readonly IPolicyRepository _policyRepository;
    private readonly ILogger<IssuePolicyCommandHandler> _logger;
    private readonly IPartnerManagementClient _partnerClient;
    private readonly IEventPublisher _eventPublisher;

    public IssuePolicyCommandHandler(
        IPolicyRepository policyRepository,
        ILogger<IssuePolicyCommandHandler> logger,
        IPartnerManagementClient partnerClient,
        IEventPublisher eventPublisher)
    {
        _policyRepository = policyRepository;
        _logger = logger;
        _partnerClient = partnerClient;
        _eventPublisher = eventPublisher;
    }

    public async Task<IssuePolicyResponse> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        // 1. Validate business rules
        // 2. Create domain entity
        var policy = Policy.Create(
            userId: request.UserId,
            productId: request.ProductId,
            sumAssured: new Money(request.SumAssured, "BDT"),
            tenure: request.TenureYears
        );

        // 3. Add nominees
        foreach (var nomineeDto in request.Nominees)
        {
            policy.AddNominee(nomineeDto.Name, nomineeDto.Relationship, nomineeDto.SharePercentage);
        }

        // 4. Persist
        await _policyRepository.AddAsync(policy, cancellationToken);
        await _policyRepository.UnitOfWork.SaveChangesAsync(cancellationToken);

        // 5. Publish domain event
        await _eventPublisher.PublishAsync(new PolicyIssuedEvent(policy.Id, policy.PolicyNumber), cancellationToken);

        _logger.LogInformation("Policy {PolicyNumber} issued for user {UserId}", policy.PolicyNumber, request.UserId);

        return new IssuePolicyResponse
        {
            PolicyId = policy.Id,
            PolicyNumber = policy.PolicyNumber,
            Status = policy.Status.ToString()
        };
    }
}
```

```csharp
// Query Example: Get Policy By ID
public class GetPolicyByIdQuery : IRequest<PolicyDto>
{
    public string PolicyId { get; set; }
}

public class GetPolicyByIdQueryHandler : IRequestHandler<GetPolicyByIdQuery, PolicyDto>
{
    private readonly IPolicyRepository _policyRepository;
    private readonly IMapper _mapper;

    public async Task<PolicyDto> Handle(GetPolicyByIdQuery request, CancellationToken cancellationToken)
    {
        var policy = await _policyRepository.GetByIdAsync(request.PolicyId, cancellationToken);
        
        if (policy == null)
            throw new NotFoundException($"Policy {request.PolicyId} not found");

        return _mapper.Map<PolicyDto>(policy);
    }
}
```

#### 4.5.3 FluentValidation Example

```csharp
public class IssuePolicyCommandValidator : AbstractValidator<IssuePolicyCommand>
{
    public IssuePolicyCommandValidator()
    {
        RuleFor(x => x.UserId)
            .NotEmpty().WithMessage("User ID is required");

        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.SumAssured)
            .GreaterThan(0).WithMessage("Sum assured must be greater than 0")
            .LessThanOrEqualTo(10000000).WithMessage("Sum assured cannot exceed 1 Crore BDT");

        RuleFor(x => x.TenureYears)
            .InclusiveBetween(1, 30).WithMessage("Tenure must be between 1 and 30 years");

        RuleFor(x => x.Nominees)
            .NotEmpty().WithMessage("At least one nominee is required")
            .Must(HaveValidSharePercentages).WithMessage("Nominee share percentages must sum to 100");
    }

    private bool HaveValidSharePercentages(List<NomineeDto> nominees)
    {
        return Math.Abs(nominees.Sum(n => n.SharePercentage) - 100) < 0.01;
    }
}
```

#### 4.5.4 MediatR Pipeline Behaviors

```csharp
// Validation Behavior - Automatically validates all commands/queries
public class ValidationBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    where TRequest : IRequest<TResponse>
{
    private readonly IEnumerable<IValidator<TRequest>> _validators;

    public ValidationBehavior(IEnumerable<IValidator<TRequest>> validators)
    {
        _validators = validators;
    }

    public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
    {
        if (!_validators.Any()) return await next();

        var context = new ValidationContext<TRequest>(request);
        
        var validationResults = await Task.WhenAll(
            _validators.Select(v => v.ValidateAsync(context, cancellationToken)));

        var failures = validationResults
            .SelectMany(r => r.Errors)
            .Where(f => f != null)
            .ToList();

        if (failures.Any())
            throw new ValidationException(failures);

        return await next();
    }
}

// Logging Behavior
public class LoggingBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    where TRequest : IRequest<TResponse>
{
    private readonly ILogger<LoggingBehavior<TRequest, TResponse>> _logger;

    public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
    {
        var requestName = typeof(TRequest).Name;
        _logger.LogInformation("Handling {RequestName}", requestName);

        var response = await next();

        _logger.LogInformation("Handled {RequestName}", requestName);

        return response;
    }
}

// Performance Behavior
public class PerformanceBehavior<TRequest, TResponse> : IPipelineBehavior<TRequest, TResponse>
    where TRequest : IRequest<TResponse>
{
    private readonly ILogger<PerformanceBehavior<TRequest, TResponse>> _logger;
    private readonly Stopwatch _timer;

    public PerformanceBehavior(ILogger<PerformanceBehavior<TRequest, TResponse>> logger)
    {
        _timer = new Stopwatch();
        _logger = logger;
    }

    public async Task<TResponse> Handle(TRequest request, RequestHandlerDelegate<TResponse> next, CancellationToken cancellationToken)
    {
        _timer.Start();

        var response = await next();

        _timer.Stop();

        var elapsedMilliseconds = _timer.ElapsedMilliseconds;

        if (elapsedMilliseconds > 500) // Log slow requests
        {
            var requestName = typeof(TRequest).Name;
            _logger.LogWarning("Long Running Request: {Name} ({ElapsedMilliseconds} milliseconds) {@Request}",
                requestName, elapsedMilliseconds, request);
        }

        return response;
    }
}
```

#### 4.5.5 Dependency Injection Setup

```csharp
// Program.cs
var builder = WebApplication.CreateBuilder(args);

// Add MediatR
builder.Services.AddMediatR(cfg => {
    cfg.RegisterServicesFromAssembly(typeof(Program).Assembly);
    
    // Add pipeline behaviors (order matters!)
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(LoggingBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(ValidationBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(PerformanceBehavior<,>));
    cfg.AddBehavior(typeof(IPipelineBehavior<,>), typeof(TransactionBehavior<,>));
});

// Add FluentValidation
builder.Services.AddValidatorsFromAssembly(typeof(Program).Assembly);

// Add AutoMapper
builder.Services.AddAutoMapper(typeof(Program).Assembly);

// Add gRPC
builder.Services.AddGrpc(options =>
{
    options.Interceptors.Add<LoggingInterceptor>();
    options.Interceptors.Add<ExceptionInterceptor>();
});

// Add EF Core
builder.Services.AddDbContext<InsuranceDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("InsuranceDb")));

// Add Repositories
builder.Services.AddScoped<IPolicyRepository, PolicyRepository>();
builder.Services.AddScoped<IQuoteRepository, QuoteRepository>();
builder.Services.AddScoped<IClaimRepository, ClaimRepository>();

// Add Domain Services
builder.Services.AddScoped<PremiumCalculationService>();
builder.Services.AddScoped<UnderwritingService>();
builder.Services.AddScoped<RiskAssessmentService>();

// Add gRPC Clients (to other services)
builder.Services.AddGrpcClient<PartnerManagement.PartnerManagementClient>(o =>
{
    o.Address = new Uri(builder.Configuration["Services:PartnerManagement"]);
});

var app = builder.Build();

// Map gRPC services
app.MapGrpcService<InsuranceEngineGrpcService>();
app.MapGrpcService<PolicyGrpcService>();
app.MapGrpcService<PremiumGrpcService>();

app.Run();
```

**Key Benefits of CQRS + MediatR:**

1. **Separation of Concerns:** Commands (write) and Queries (read) are separate
2. **Single Responsibility:** Each handler does one thing
3. **Testability:** Easy to unit test handlers in isolation
4. **Cross-Cutting Concerns:** Validation, logging, performance tracking via pipeline behaviors
5. **Scalability:** Can scale read and write models independently
6. **Maintainability:** Easy to add new features without modifying existing code

"""

# Insert before the current 4.4 section
content = content.replace("### 4.4 VSA Internal Structure Example", detailed_insurance_engine + "\n\n### 4.4 VSA Internal Structure Example")

print("✅ Phase 2 complete!")
print("- Updated proto structure to proto/insuretech/v1/")
print("- Added events/ folder to proto structure")
print("- Added detailed CQRS/MediatR Insurance Engine documentation with:")
print("  * Complete folder structure")
print("  * MediatR command/query examples with syntax highlighting")
print("  * FluentValidation examples")
print("  * Pipeline behaviors (Validation, Logging, Performance)")
print("  * Dependency injection setup")

# Write back
with open("SRS_V3_FINAL_DRAFT.md", "w", encoding="utf-8") as f:
    f.write(content)

print("\nFile updated successfully!")
