using MediatR;
using Microsoft.Extensions.Logging;
using Microsoft.EntityFrameworkCore;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.SharedKernel.Infrastructure;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Underwriting.Application.Commands;

// ===== RequestQuote =====
public sealed class RequestQuoteCommandHandler : IRequestHandler<RequestQuoteCommand, RequestQuoteResponse>
{
    private readonly IRepository<QuoteEntity> _repository;
    private readonly ILogger<RequestQuoteCommandHandler> _logger;

    public RequestQuoteCommandHandler(IRepository<QuoteEntity> repository, ILogger<RequestQuoteCommandHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<RequestQuoteResponse> Handle(RequestQuoteCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quoteId = Guid.NewGuid();
            var quoteNumber = $"QT-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            // TM-001: Age Limits (18 - 65)
            if (request.ApplicantAge < 18 || request.ApplicantAge > 65)
            {
                return new RequestQuoteResponse
                {
                    Error = new Error { Code = "INVALID_AGE", Message = "Applicant age must be between 18 and 65 years." }
                };
            }

            // TM-001: Sum Assured Bounds (Max 1,000,000 BDT for auto-quoting)
            if (request.SumAssured < 5000 || request.SumAssured > 1000000)
            {
                return new RequestQuoteResponse
                {
                    Error = new Error { Code = "INVALID_SUM_ASSURED", Message = "Sum Assured must be between 5,000 and 1,000,000 BDT for this product." }
                };
            }

            var basePremium = CalculateBasePremium(request.SumAssured, request.TermYears, request.ApplicantAge, request.Smoker);
            var riderPremium = (request.RiderCodes?.Count ?? 0) * (long)(basePremium * 0.05m);
            var taxAmount = (long)((basePremium + riderPremium) * 0.15m);
            var totalPremium = basePremium + riderPremium + taxAmount;

            var entity = new QuoteEntity
            {
                QuoteId = quoteId,
                QuoteNumber = quoteNumber,
                BeneficiaryId = Guid.Parse(request.BeneficiaryId),
                InsurerProductId = Guid.Parse(request.InsurerProductId),
                Status = "DRAFT",
                SumAssured = request.SumAssured,
                TermYears = request.TermYears,
                PremiumPaymentMode = request.PremiumPaymentMode,
                BasePremium = basePremium,
                RiderPremium = riderPremium,
                TaxAmount = taxAmount,
                TotalPremium = totalPremium,
                SelectedRiders = request.RiderCodes != null ? System.Text.Json.JsonSerializer.Serialize(request.RiderCodes) : null,
                ApplicantAge = request.ApplicantAge,
                Smoker = request.Smoker,
                ValidUntil = DateTime.UtcNow.AddDays(30),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _repository.AddAsync(entity, cancellationToken);

            _logger.LogInformation("Quote created: {QuoteNumber}", quoteNumber);

            return new RequestQuoteResponse
            {
                QuoteId = quoteId.ToString(),
                QuoteNumber = quoteNumber,
                BasePremium = new Money { Amount = basePremium, Currency = "BDT" },
                TotalPremium = new Money { Amount = totalPremium, Currency = "BDT" },
                ValidUntil = entity.ValidUntil.ToString("yyyy-MM-dd"),
                Message = "Quote generated successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create quote");
            return new RequestQuoteResponse { Error = new Error { Code = "QUOTE_FAILED", Message = ex.Message } };
        }
    }

    private static long CalculateBasePremium(long sumAssured, int termYears, int age, bool smoker)
    {
        var rate = 0.005m;
        var ageFactor = 1m + (age - 30) * 0.02m;
        if (ageFactor < 0.5m) ageFactor = 0.5m;
        if (smoker) ageFactor *= 1.25m;
        return (long)(sumAssured * rate * termYears * ageFactor / 1000m);
    }
}

// ===== SubmitHealthDeclaration =====
public sealed class SubmitHealthDeclarationCommandHandler : IRequestHandler<SubmitHealthDeclarationCommand, SubmitHealthDeclarationResponse>
{
    private readonly IRepository<QuoteEntity> _quoteRepository;
    private readonly IRepository<HealthDeclarationEntity> _declarationRepository;
    private readonly ILogger<SubmitHealthDeclarationCommandHandler> _logger;

    public SubmitHealthDeclarationCommandHandler(
        IRepository<QuoteEntity> quoteRepository,
        IRepository<HealthDeclarationEntity> declarationRepository,
        ILogger<SubmitHealthDeclarationCommandHandler> logger)
    {
        _quoteRepository = quoteRepository;
        _declarationRepository = declarationRepository;
        _logger = logger;
    }

    public async Task<SubmitHealthDeclarationResponse> Handle(SubmitHealthDeclarationCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _quoteRepository.GetByIdAsync(Guid.Parse(request.QuoteId), cancellationToken);
            if (quote == null)
                return new SubmitHealthDeclarationResponse { Error = new Error { Code = "QUOTE_NOT_FOUND", Message = "Quote not found" } };

            var medicalExamRequired = request.HasPreExistingConditions || request.Smoker;
            var autoApproval = !medicalExamRequired && !request.AlcoholConsumer && request.OccupationRiskLevel != "HIGH";

            var entity = new HealthDeclarationEntity
            {
                DeclarationId = Guid.NewGuid(),
                QuoteId = Guid.Parse(request.QuoteId),
                HeightCm = request.HeightCm,
                WeightKg = request.WeightKg,
                HasPreExistingConditions = request.HasPreExistingConditions,
                PreExistingConditions = request.PreExistingConditions,
                Smoker = request.Smoker,
                AlcoholConsumer = request.AlcoholConsumer,
                OccupationRiskLevel = request.OccupationRiskLevel,
                MedicalExamRequired = medicalExamRequired,
                AutoApprovalPossible = autoApproval,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _declarationRepository.AddAsync(entity, cancellationToken);

            // Auto-submit for underwriting
            quote.Status = "PENDING_UNDERWRITING";
            quote.UpdatedAt = DateTime.UtcNow;
            await _quoteRepository.UpdateAsync(quote, cancellationToken);

            return new SubmitHealthDeclarationResponse
            {
                Message = "Health declaration submitted",
                MedicalExamRequired = medicalExamRequired,
                AutoApprovalPossible = autoApproval
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit health declaration");
            return new SubmitHealthDeclarationResponse { Error = new Error { Code = "DECLARATION_FAILED", Message = ex.Message } };
        }
    }
}

// ===== ApproveUnderwriting =====
public sealed class ApproveUnderwritingCommandHandler : IRequestHandler<ApproveUnderwritingCommand, ApproveUnderwritingResponse>
{
    private readonly IRepository<QuoteEntity> _quoteRepository;
    private readonly IRepository<UnderwritingDecisionEntity> _decisionRepository;
    private readonly ILogger<ApproveUnderwritingCommandHandler> _logger;

    public ApproveUnderwritingCommandHandler(
        IRepository<QuoteEntity> quoteRepository,
        IRepository<UnderwritingDecisionEntity> decisionRepository,
        ILogger<ApproveUnderwritingCommandHandler> logger)
    {
        _quoteRepository = quoteRepository;
        _decisionRepository = decisionRepository;
        _logger = logger;
    }

    public async Task<ApproveUnderwritingResponse> Handle(ApproveUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _quoteRepository.GetByIdAsync(Guid.Parse(request.QuoteId), cancellationToken);
            if (quote == null)
                return new ApproveUnderwritingResponse { Error = new Error { Code = "QUOTE_NOT_FOUND", Message = "Quote not found" } };

            // FR-177: Ensure health declaration exists before approval
            if (quote.Status != "PENDING_UNDERWRITING")
            {
                return new ApproveUnderwritingResponse 
                { 
                    Error = new Error { Code = "PRECONDITION_FAILED", Message = "Health declaration must be submitted before underwriting approval." } 
                };
            }

            var decision = new UnderwritingDecisionEntity
            {
                DecisionId = Guid.NewGuid(),
                QuoteId = quote.QuoteId,
                UnderwriterId = Guid.Parse(request.UnderwriterId),
                Decision = "APPROVED",
                RiskLevel = request.RiskLevel,
                PremiumAdjusted = request.PremiumAdjusted,
                AdjustedPremium = request.AdjustedPremium,
                Comments = request.Comments,
                CreatedAt = DateTime.UtcNow
            };

            await _decisionRepository.AddAsync(decision, cancellationToken);

            quote.Status = "APPROVED";
            if (request.PremiumAdjusted && request.AdjustedPremium.HasValue)
                quote.TotalPremium = request.AdjustedPremium.Value;
            quote.UpdatedAt = DateTime.UtcNow;
            await _quoteRepository.UpdateAsync(quote, cancellationToken);

            return new ApproveUnderwritingResponse { DecisionId = decision.DecisionId.ToString(), Message = "Underwriting approved" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve underwriting");
            return new ApproveUnderwritingResponse { Error = new Error { Code = "APPROVE_FAILED", Message = ex.Message } };
        }
    }
}

// ===== RejectUnderwriting =====
public sealed class RejectUnderwritingCommandHandler : IRequestHandler<RejectUnderwritingCommand, RejectUnderwritingResponse>
{
    private readonly IRepository<QuoteEntity> _quoteRepository;
    private readonly IRepository<UnderwritingDecisionEntity> _decisionRepository;
    private readonly ILogger<RejectUnderwritingCommandHandler> _logger;

    public RejectUnderwritingCommandHandler(
        IRepository<QuoteEntity> quoteRepository,
        IRepository<UnderwritingDecisionEntity> decisionRepository,
        ILogger<RejectUnderwritingCommandHandler> logger)
    {
        _quoteRepository = quoteRepository;
        _decisionRepository = decisionRepository;
        _logger = logger;
    }

    public async Task<RejectUnderwritingResponse> Handle(RejectUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _quoteRepository.GetByIdAsync(Guid.Parse(request.QuoteId), cancellationToken);
            if (quote == null)
                return new RejectUnderwritingResponse { Error = new Error { Code = "QUOTE_NOT_FOUND", Message = "Quote not found" } };

            var decision = new UnderwritingDecisionEntity
            {
                DecisionId = Guid.NewGuid(),
                QuoteId = quote.QuoteId,
                UnderwriterId = Guid.Parse(request.UnderwriterId),
                Decision = "REJECTED",
                RiskLevel = request.RiskLevel,
                RejectionReason = request.Reason,
                Comments = request.Comments,
                CreatedAt = DateTime.UtcNow
            };

            await _decisionRepository.AddAsync(decision, cancellationToken);

            quote.Status = "REJECTED";
            quote.UpdatedAt = DateTime.UtcNow;
            await _quoteRepository.UpdateAsync(quote, cancellationToken);

            return new RejectUnderwritingResponse { DecisionId = decision.DecisionId.ToString(), Message = "Underwriting rejected" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject underwriting");
            return new RejectUnderwritingResponse { Error = new Error { Code = "REJECT_FAILED", Message = ex.Message } };
        }
    }
}

// ===== ConvertQuoteToPolicy =====
public sealed class ConvertQuoteToPolicyCommandHandler : IRequestHandler<ConvertQuoteToPolicyCommand, ConvertQuoteToPolicyResponse>
{
    private readonly IRepository<QuoteEntity> _quoteRepository;
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly InsuranceDbContext _dbContext;
    private readonly ILogger<ConvertQuoteToPolicyCommandHandler> _logger;

    public ConvertQuoteToPolicyCommandHandler(
        IRepository<QuoteEntity> quoteRepository,
        IRepository<PolicyEntity> policyRepository,
        InsuranceDbContext dbContext,
        ILogger<ConvertQuoteToPolicyCommandHandler> logger)
    {
        _quoteRepository = quoteRepository;
        _policyRepository = policyRepository;
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<ConvertQuoteToPolicyResponse> Handle(ConvertQuoteToPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quote = await _quoteRepository.GetByIdAsync(Guid.Parse(request.QuoteId), cancellationToken);
            if (quote == null)
                return new ConvertQuoteToPolicyResponse { Error = new Error { Code = "QUOTE_NOT_FOUND", Message = "Quote not found" } };

            if (quote.Status != "APPROVED")
                return new ConvertQuoteToPolicyResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Quote must be APPROVED to convert, current: '{quote.Status}'" } };

            // Get sequence number for policy
            var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open)
                await connection.OpenAsync(cancellationToken);
            using var cmd = connection.CreateCommand();
            cmd.CommandText = "SELECT nextval('insurance_schema.policy_number_seq')";
            var seqResult = await cmd.ExecuteScalarAsync(cancellationToken);
            var seq = Convert.ToInt64(seqResult);

            var policyNumber = $"LBT-{DateTime.UtcNow.Year}-0001-{seq.ToString().PadLeft(6, '0')}";

            var policy = new PolicyEntity
            {
                PolicyId = Guid.NewGuid(),
                PolicyNumber = policyNumber,
                ProductId = quote.InsurerProductId,
                CustomerId = quote.BeneficiaryId,
                QuoteId = quote.QuoteId,
                Status = "PENDING_PAYMENT",
                PremiumAmount = quote.TotalPremium,
                PremiumCurrency = "BDT",
                SumInsured = quote.SumAssured,
                SumInsuredCurrency = "BDT",
                TenureMonths = quote.TermYears * 12,
                StartDate = DateTime.UtcNow,
                EndDate = DateTime.UtcNow.AddYears(quote.TermYears),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _policyRepository.AddAsync(policy, cancellationToken);

            quote.Status = "CONVERTED";
            quote.ConvertedPolicyId = policy.PolicyId;
            quote.ConvertedAt = DateTime.UtcNow;
            quote.UpdatedAt = DateTime.UtcNow;
            await _quoteRepository.UpdateAsync(quote, cancellationToken);

            _logger.LogInformation("Quote {QuoteNumber} converted to Policy {PolicyNumber}", quote.QuoteNumber, policyNumber);

            return new ConvertQuoteToPolicyResponse
            {
                PolicyId = policy.PolicyId.ToString(),
                PolicyNumber = policyNumber,
                Message = "Quote converted to policy successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to convert quote to policy");
            return new ConvertQuoteToPolicyResponse { Error = new Error { Code = "CONVERT_FAILED", Message = ex.Message } };
        }
    }
}
