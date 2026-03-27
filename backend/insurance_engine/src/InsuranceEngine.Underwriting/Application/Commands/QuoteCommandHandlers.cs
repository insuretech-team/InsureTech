using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Underwriting.Application.Commands;

public sealed class CreateQuoteCommandHandler : IRequestHandler<CreateQuoteCommand, Result<string>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CreateQuoteCommandHandler> _logger;

    public CreateQuoteCommandHandler(DbContext dbContext, ILogger<CreateQuoteCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(CreateQuoteCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var quoteId = Guid.NewGuid();
            var quoteNumber = $"QT-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
            var now = DateTime.UtcNow;

            var basePremium = CalculateBasePremium(request.SumAssured, request.TermYears, request.ApplicantAge);
            var totalPremium = basePremium;

            var sql = @"
                INSERT INTO insurance_schema.quotes (
                    quote_id, quote_number, beneficiary_id, insurer_product_id, status,
                    sum_assured, term_years, premium_payment_mode, base_premium, 
                    total_premium, applicant_age, smoker, valid_until, created_at
                ) VALUES (
                    @QuoteId, @QuoteNumber, @BeneficiaryId, @ProductId, @Status,
                    @SumAssured, @TermYears, @PremiumPaymentMode, @BasePremium, 
                    @TotalPremium, @ApplicantAge, @Smoker, @ValidUntil, @CreatedAt
                )";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            await connection.ExecuteAsync(sql, new
            {
                QuoteId = quoteId,
                QuoteNumber = quoteNumber,
                BeneficiaryId = request.BeneficiaryId,
                ProductId = request.ProductId,
                Status = "DRAFT",
                SumAssured = request.SumAssured,
                TermYears = request.TermYears,
                PremiumPaymentMode = request.PremiumPaymentMode,
                BasePremium = basePremium,
                TotalPremium = totalPremium,
                ApplicantAge = request.ApplicantAge,
                Smoker = request.Smoker,
                ValidUntil = now.AddDays(30),
                CreatedAt = now
            });

            _logger.LogInformation("Quote created: {QuoteId} ({QuoteNumber})", quoteId, quoteNumber);
            return Result<string>.Ok(quoteId.ToString());
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create quote");
            return Result<string>.Fail("QUOTE_CREATION_FAILED", ex.Message);
        }
    }

    private decimal CalculateBasePremium(decimal sumAssured, int termYears, int age)
    {
        var rate = 0.005m;
        var ageFactor = 1 + (age - 30) * 0.05m;
        if (ageFactor < 0.5m) ageFactor = 0.5m;
        return sumAssured * rate * termYears * ageFactor / 1000;
    }
}

public sealed class SubmitQuoteForUnderwritingCommandHandler : IRequestHandler<SubmitQuoteForUnderwritingCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<SubmitQuoteForUnderwritingCommandHandler> _logger;

    public SubmitQuoteForUnderwritingCommandHandler(DbContext dbContext, ILogger<SubmitQuoteForUnderwritingCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(SubmitQuoteForUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.quotes
                SET status = 'PENDING_UNDERWRITING', updated_at = @UpdatedAt
                WHERE quote_id = @QuoteId::uuid AND status = 'DRAFT' AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                QuoteId = request.QuoteId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.Fail("QUOTE_NOT_FOUND_OR_INVALID_STATE", "Quote not found or cannot be submitted");

            _logger.LogInformation("Quote submitted for underwriting: {QuoteId}", request.QuoteId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit quote for underwriting {QuoteId}", request.QuoteId);
            return Result<bool>.Fail("QUOTE_SUBMIT_FAILED", ex.Message);
        }
    }
}

public sealed class ApproveQuoteCommandHandler : IRequestHandler<ApproveQuoteCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ApproveQuoteCommandHandler> _logger;

    public ApproveQuoteCommandHandler(DbContext dbContext, ILogger<ApproveQuoteCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(ApproveQuoteCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.quotes
                SET status = 'APPROVED', updated_at = @UpdatedAt
                WHERE quote_id = @QuoteId::uuid AND status = 'PENDING_UNDERWRITING' AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                QuoteId = request.QuoteId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.Fail("QUOTE_NOT_FOUND_OR_INVALID_STATE", "Quote not found or cannot be approved");

            _logger.LogInformation("Quote approved: {QuoteId}", request.QuoteId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve quote {QuoteId}", request.QuoteId);
            return Result<bool>.Fail("QUOTE_APPROVE_FAILED", ex.Message);
        }
    }
}

public sealed class RejectQuoteCommandHandler : IRequestHandler<RejectQuoteCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<RejectQuoteCommandHandler> _logger;

    public RejectQuoteCommandHandler(DbContext dbContext, ILogger<RejectQuoteCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(RejectQuoteCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.quotes
                SET status = 'REJECTED', updated_at = @UpdatedAt
                WHERE quote_id = @QuoteId::uuid AND status = 'PENDING_UNDERWRITING' AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                QuoteId = request.QuoteId,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.Fail("QUOTE_NOT_FOUND_OR_INVALID_STATE", "Quote not found or cannot be rejected");

            _logger.LogInformation("Quote rejected: {QuoteId}, Reason: {Reason}", request.QuoteId, request.Reason);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject quote {QuoteId}", request.QuoteId);
            return Result<bool>.Fail("QUOTE_REJECT_FAILED", ex.Message);
        }
    }
}
