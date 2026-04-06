using InsuranceEngine.Endorsements.Infrastructure;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Endorsements;

public class EndorsementType
{
    public const string AddressChange = "ADDRESS_CHANGE";
    public const string SumIncrease = "SUM_INCREASE";
    public const string SumDecrease = "SUM_DECREASE";
    public const string NomineeChange = "NOMINEE_CHANGE";
    public const string BankAccountChange = "BANK_ACCOUNT_CHANGE";
    public const string ContactChange = "CONTACT_CHANGE";
}

public enum EndorsementStatus
{
    Pending = 0,
    UnderReview = 1,
    Approved = 2,
    Rejected = 3,
    Processed = 4,
    Cancelled = 5
}

public enum EndorsementApprovalRequirement
{
    None = 0,
    AutoApproved = 1,
    RequiresUnderwriterReview = 2,
    RequiresManagerApproval = 3
}

public class SumChangeEndorsementResult
{
    public string PolicyId { get; set; } = string.Empty;
    public string EndorsementId { get; set; } = string.Empty;
    public string EndorsementNumber { get; set; } = string.Empty;
    public string EndorsementType { get; set; } = string.Empty;
    public decimal OldSumAssured { get; set; }
    public decimal NewSumAssured { get; set; }
    public decimal PremiumDifference { get; set; }
    public decimal RefundAmount { get; set; }
    public decimal AdditionalPremium { get; set; }
    public DateTime EffectiveDate { get; set; }
    public DateTime ProcessedAt { get; set; }
    public string Status { get; set; } = string.Empty;
    public string Message { get; set; } = string.Empty;
}

public class EndorsementDocumentResult
{
    public string DocumentId { get; set; } = string.Empty;
    public string DocumentType { get; set; } = string.Empty;
    public string DocumentNumber { get; set; } = string.Empty;
    public string FileUrl { get; set; } = string.Empty;
    public byte[]? FileContent { get; set; }
    public DateTime GeneratedAt { get; set; }
    public bool IsSuccess { get; set; }
    public string? ErrorMessage { get; set; }
}

public class EndorsementRequest
{
    public string PolicyId { get; set; } = string.Empty;
    public string EndorsementType { get; set; } = string.Empty;
    public string? Reason { get; set; }
    public Dictionary<string, object>? Changes { get; set; }
}

public interface IEndorsementProcessingService
{
    Task<SumChangeEndorsementResult> ProcessSumDecreaseEndorsementAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        DateTime policyStartDate,
        CancellationToken ct = default);

    Task<SumChangeEndorsementResult> ProcessSumIncreaseEndorsementAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        CancellationToken ct = default);

    Task<EndorsementDocumentResult> GenerateEndorsementDocumentAsync(
        string policyId,
        string endorsementId,
        string endorsementNumber,
        string endorsementType,
        Dictionary<string, string>? changes = null,
        CancellationToken ct = default);

    Task<bool> ValidateEndorsementRequestAsync(EndorsementRequest request, CancellationToken ct = default);

    Task<decimal> CalculateSumDecreaseRefundAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        DateTime policyStartDate,
        CancellationToken ct = default);

    Task<decimal> CalculateSumIncreasePremiumAsync(
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        CancellationToken ct = default);
}

public class EndorsementProcessingService : IEndorsementProcessingService
{
    private readonly ISqlEndorsementDataGateway _gateway;
    private readonly ILogger<EndorsementProcessingService> _logger;

    private const decimal MinimumSumAssuredChangePercent = 0.10m;

    public EndorsementProcessingService(
        ISqlEndorsementDataGateway gateway,
        ILogger<EndorsementProcessingService> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<SumChangeEndorsementResult> ProcessSumDecreaseEndorsementAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        DateTime policyStartDate,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "Processing sum decrease endorsement for policy {PolicyId}: {OldSum} -> {NewSum}",
            policyId, currentSumAssured, newSumAssured);

        var changePercent = Math.Abs(currentSumAssured - newSumAssured) / currentSumAssured;
        
        if (newSumAssured >= currentSumAssured)
        {
            return new SumChangeEndorsementResult
            {
                PolicyId = policyId,
                Status = "REJECTED",
                Message = "New sum assured must be less than current sum assured for sum decrease endorsement"
            };
        }

        if (changePercent < MinimumSumAssuredChangePercent)
        {
            return new SumChangeEndorsementResult
            {
                PolicyId = policyId,
                Status = "REQUIRES_APPROVAL",
                Message = $"Sum decrease of {changePercent:P0} requires underwriting approval (minimum {MinimumSumAssuredChangePercent:P0})"
            };
        }

        var refundAmount = await CalculateSumDecreaseRefundAsync(
            policyId, currentSumAssured, newSumAssured, currentPremium, policyStartDate, ct);

        var endorsementId = Guid.NewGuid().ToString();
        var endorsementNumber = $"END-{DateTime.UtcNow:yyyyMMdd}-{endorsementId[..8].ToUpper()}";

        var endorsement = new InsuranceEngine.SharedKernel.Persistence.Entities.EndorsementEntity
        {
            EndorsementId = endorsementId,
            EndorsementNumber = endorsementNumber,
            PolicyId = policyId,
            Type = EndorsementType.SumDecrease,
            Status = "APPROVED",
            OldSumAssured = currentSumAssured,
            NewSumAssured = newSumAssured,
            RefundAmount = refundAmount,
            EffectiveDate = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        await _gateway.CreateAsync(endorsement, ct);

        var result = new SumChangeEndorsementResult
        {
            PolicyId = policyId,
            EndorsementId = endorsementId,
            EndorsementNumber = endorsementNumber,
            EndorsementType = EndorsementType.SumDecrease,
            OldSumAssured = currentSumAssured,
            NewSumAssured = newSumAssured,
            PremiumDifference = 0,
            RefundAmount = refundAmount,
            AdditionalPremium = 0,
            EffectiveDate = DateTime.UtcNow,
            ProcessedAt = DateTime.UtcNow,
            Status = "APPROVED",
            Message = $"Sum decreased from {currentSumAssured:N0} to {newSumAssured:N0}. Refund of {refundAmount:N2} BDT will be processed."
        };

        _logger.LogInformation(
            "Sum decrease endorsement approved for {PolicyId}: Refund={RefundAmount:N2} BDT",
            policyId, result.RefundAmount);

        return result;
    }

    public async Task<SumChangeEndorsementResult> ProcessSumIncreaseEndorsementAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "Processing sum increase endorsement for policy {PolicyId}: {OldSum} -> {NewSum}",
            policyId, currentSumAssured, newSumAssured);

        var changePercent = Math.Abs(newSumAssured - currentSumAssured) / currentSumAssured;

        if (newSumAssured <= currentSumAssured)
        {
            return new SumChangeEndorsementResult
            {
                PolicyId = policyId,
                Status = "REJECTED",
                Message = "New sum assured must be greater than current sum assured for sum increase endorsement"
            };
        }

        if (changePercent >= MinimumSumAssuredChangePercent)
        {
            return new SumChangeEndorsementResult
            {
                PolicyId = policyId,
                Status = "REQUIRES_APPROVAL",
                Message = $"Sum increase of {changePercent:P0} requires underwriting approval (minimum {MinimumSumAssuredChangePercent:P0})"
            };
        }

        var additionalPremium = await CalculateSumIncreasePremiumAsync(
            currentSumAssured, newSumAssured, currentPremium, ct);

        var endorsementId = Guid.NewGuid().ToString();
        var endorsementNumber = $"END-{DateTime.UtcNow:yyyyMMdd}-{endorsementId[..8].ToUpper()}";

        var endorsement = new InsuranceEngine.SharedKernel.Persistence.Entities.EndorsementEntity
        {
            EndorsementId = endorsementId,
            EndorsementNumber = endorsementNumber,
            PolicyId = policyId,
            Type = EndorsementType.SumIncrease,
            Status = "APPROVED",
            OldSumAssured = currentSumAssured,
            NewSumAssured = newSumAssured,
            AdditionalPremium = additionalPremium,
            EffectiveDate = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        await _gateway.CreateAsync(endorsement, ct);

        var result = new SumChangeEndorsementResult
        {
            PolicyId = policyId,
            EndorsementId = endorsementId,
            EndorsementNumber = endorsementNumber,
            EndorsementType = EndorsementType.SumIncrease,
            OldSumAssured = currentSumAssured,
            NewSumAssured = newSumAssured,
            PremiumDifference = additionalPremium,
            RefundAmount = 0,
            AdditionalPremium = additionalPremium,
            EffectiveDate = DateTime.UtcNow,
            ProcessedAt = DateTime.UtcNow,
            Status = "APPROVED",
            Message = $"Sum increased from {currentSumAssured:N0} to {newSumAssured:N0}. Additional premium of {additionalPremium:N2} BDT is required."
        };

        _logger.LogInformation(
            "Sum increase endorsement approved for {PolicyId}: Additional Premium={AdditionalPremium:N2} BDT",
            policyId, result.AdditionalPremium);

        return result;
    }

    public async Task<EndorsementDocumentResult> GenerateEndorsementDocumentAsync(
        string policyId,
        string endorsementId,
        string endorsementNumber,
        string endorsementType,
        Dictionary<string, string>? changes = null,
        CancellationToken ct = default)
    {
        _logger.LogInformation(
            "Generating endorsement document for {EndorsementNumber} ({EndorsementType})",
            endorsementNumber, endorsementType);

        try
        {
            var documentId = Guid.NewGuid().ToString();
            var fileUrl = $"/documents/endorsements/{endorsementId}/{documentId}.pdf";

            var document = new InsuranceEngine.SharedKernel.Persistence.Entities.EndorsementDocumentEntity
            {
                DocumentId = documentId,
                EndorsementId = endorsementId,
                DocumentType = "ENDORSEMENT",
                DocumentNumber = endorsementNumber,
                FileUrl = fileUrl,
                Status = "GENERATED",
                GeneratedAt = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow
            };

            await _gateway.CreateDocumentAsync(document, ct);

            _logger.LogInformation(
                "Endorsement document generated: {DocumentId} for endorsement {EndorsementId}",
                documentId, endorsementId);

            return new EndorsementDocumentResult
            {
                DocumentId = documentId,
                DocumentType = "ENDORSEMENT",
                DocumentNumber = endorsementNumber,
                FileUrl = fileUrl,
                GeneratedAt = DateTime.UtcNow,
                IsSuccess = true
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to generate endorsement document for {EndorsementNumber}", endorsementNumber);
            return new EndorsementDocumentResult
            {
                DocumentId = string.Empty,
                DocumentType = "ENDORSEMENT",
                DocumentNumber = endorsementNumber,
                GeneratedAt = DateTime.UtcNow,
                IsSuccess = false,
                ErrorMessage = ex.Message
            };
        }
    }

    public Task<bool> ValidateEndorsementRequestAsync(EndorsementRequest request, CancellationToken ct = default)
    {
        if (string.IsNullOrEmpty(request.PolicyId))
            return Task.FromResult(false);

        if (string.IsNullOrEmpty(request.EndorsementType))
            return Task.FromResult(false);

        var validTypes = new[] 
        { 
            EndorsementType.AddressChange, 
            EndorsementType.SumIncrease, 
            EndorsementType.SumDecrease,
            EndorsementType.NomineeChange,
            EndorsementType.BankAccountChange,
            EndorsementType.ContactChange
        };

        if (!validTypes.Contains(request.EndorsementType))
            return Task.FromResult(false);

        return Task.FromResult(true);
    }

    public Task<decimal> CalculateSumDecreaseRefundAsync(
        string policyId,
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        DateTime policyStartDate,
        CancellationToken ct = default)
    {
        var totalDays = 365;
        var daysElapsed = Math.Max(1, (DateTime.UtcNow - policyStartDate).Days);
        var daysRemaining = Math.Max(0, totalDays - daysElapsed);
        var remainingPercent = (decimal)daysRemaining / totalDays;

        var premiumReductionPerDay = (currentPremium / totalDays);
        var refundAmount = premiumReductionPerDay * daysRemaining;

        _logger.LogInformation(
            "Sum decrease refund calculated for {PolicyId}: {DaysRemaining} days remaining, refund = {Refund:N2}",
            policyId, daysRemaining, refundAmount);

        return Task.FromResult(Math.Round(refundAmount, 2));
    }

    public Task<decimal> CalculateSumIncreasePremiumAsync(
        decimal currentSumAssured,
        decimal newSumAssured,
        decimal currentPremium,
        CancellationToken ct = default)
    {
        var premiumRate = currentPremium / currentSumAssured;
        var additionalSum = newSumAssured - currentSumAssured;
        var additionalPremium = additionalSum * premiumRate;

        _logger.LogInformation(
            "Sum increase premium calculated: Additional sum = {AdditionalSum:N0}, Premium = {Premium:N2}",
            additionalSum, additionalPremium);

        return Task.FromResult(Math.Round(additionalPremium, 2));
    }
}
