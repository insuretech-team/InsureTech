using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using InsuranceEngine.Policy.Domain;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class CreatePolicyCommandHandler : IRequestHandler<CreatePolicyCommand, CreatePolicyResponse>
{
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly IRepository<ProductEntity> _productRepository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly IRepository<PolicyRiderEntity> _riderRepository;
    private readonly InsuranceDbContext _dbContext;
    private readonly ILogger<CreatePolicyCommandHandler> _logger;
    private readonly IPdfGenerator _pdfGenerator;
    private readonly IKafkaPublisher _kafkaPublisher;

    public CreatePolicyCommandHandler(
        IRepository<PolicyEntity> policyRepository,
        IRepository<ProductEntity> productRepository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        IRepository<PolicyRiderEntity> riderRepository,
        InsuranceDbContext dbContext,
        ILogger<CreatePolicyCommandHandler> logger,
        IPdfGenerator pdfGenerator,
        IKafkaPublisher kafkaPublisher)
    {
        _policyRepository = policyRepository;
        _productRepository = productRepository;
        _nomineeRepository = nomineeRepository;
        _riderRepository = riderRepository;
        _dbContext = dbContext;
        _logger = logger;
        _pdfGenerator = pdfGenerator;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<CreatePolicyResponse> Handle(CreatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // 0. Validate product exists and is ACTIVE
            var product = await _productRepository.GetByIdAsync(Guid.Parse(request.ProductId), cancellationToken);
            if (product == null)
            {
                return new CreatePolicyResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "PRODUCT_NOT_FOUND", Message = "Product not found" }
                };
            }
            if (product.Status != "ACTIVE")
            {
                return new CreatePolicyResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "PRODUCT_INACTIVE", Message = "Cannot create policy for inactive product" }
                };
            }

            // 1. Get sequence number from DB for collision-safe policy number (FR-034)
            long sequenceNumber;
            if (_dbContext.Database.ProviderName == "Microsoft.EntityFrameworkCore.Sqlite")
            {
                // Professional fallback for cross-provider testing (SQLite/InMemory)
                sequenceNumber = await _dbContext.Policies.IgnoreQueryFilters().CountAsync(cancellationToken) + 1;
            }
            else
            {
                var connection = _dbContext.Database.GetDbConnection();
                if (connection.State != System.Data.ConnectionState.Open)
                    await connection.OpenAsync(cancellationToken);

                using var cmd = connection.CreateCommand();
                cmd.CommandText = "SELECT nextval('insurance_schema.policy_number_seq')";
                var seqResult = await cmd.ExecuteScalarAsync(cancellationToken);
                sequenceNumber = Convert.ToInt64(seqResult);
            }

            // 2. Create Domain Aggregate (DDD)
            var policy = PolicyAggregate.Create(
                productId: Guid.Parse(request.ProductId),
                productCode: product.ProductCode.Length >= 4 ? product.ProductCode[..4] : product.ProductCode.PadRight(4, '0'),
                insuranceType: product.Category,
                customerId: Guid.Parse(request.CustomerId),
                premium: request.PremiumAmount,
                sumInsured: request.SumInsured,
                tenure: request.TenureMonths,
                startDate: request.StartDate,
                sequenceNumber: sequenceNumber
            );

            // 3. Persist Policy Entity
            var policyEntity = new PolicyEntity
            {
                PolicyId = policy.Id,
                PolicyNumber = policy.PolicyNumber,
                ProductId = policy.ProductId,
                CustomerId = policy.CustomerId,
                PartnerId = string.IsNullOrEmpty(request.PartnerId) ? null : Guid.Parse(request.PartnerId),
                AgentId = string.IsNullOrEmpty(request.AgentId) ? null : Guid.Parse(request.AgentId),
                QuoteId = string.IsNullOrEmpty(request.QuoteId) ? null : Guid.Parse(request.QuoteId),
                Status = "PENDING_PAYMENT",
                PremiumAmount = (long)(request.PremiumAmount * 100), // Store in paisa
                PremiumCurrency = "BDT",
                SumInsured = (long)(request.SumInsured * 100),
                SumInsuredCurrency = "BDT",
                TenureMonths = request.TenureMonths,
                StartDate = request.StartDate,
                EndDate = request.StartDate.AddMonths(request.TenureMonths),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _policyRepository.AddAsync(policyEntity, cancellationToken);

            // 4. Add Nominees
            if (request.Nominees != null)
            {
                foreach (var nominee in request.Nominees)
                {
                    var nomineeEntity = new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = policy.Id,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth?.ToDateTime() ?? DateTime.UtcNow,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    };
                    await _nomineeRepository.AddAsync(nomineeEntity, cancellationToken);
                }
            }

            // 5. FR-035: Generate PDF Document (Simulated)
            await _pdfGenerator.GeneratePolicyDocumentAsync(policy.PolicyNumber, "N/A", "N/A", request.PremiumAmount);

            // 6. FR-019: Kafka Event Streaming
            var policyEvent = new PolicyIssuedEvent(
                policyEntity.PolicyId, 
                policyEntity.PolicyNumber, 
                policyEntity.CustomerId, 
                policyEntity.PremiumAmount,
                policyEntity.PartnerId,
                policyEntity.AgentId);
            await _kafkaPublisher.PublishAsync("insurance.policy.created", policyEvent);

            _logger.LogInformation("Policy created: {PolicyNumber} for Customer: {CustomerId}", policy.PolicyNumber, request.CustomerId);

            return new CreatePolicyResponse
            {
                PolicyId = policy.Id.ToString(),
                PolicyNumber = policy.PolicyNumber,
                Message = "Policy created successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create policy for customer {CustomerId}", request.CustomerId);
            return new CreatePolicyResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "POLICY_CREATION_FAILED", Message = ex.Message }
            };
        }
    }
}
