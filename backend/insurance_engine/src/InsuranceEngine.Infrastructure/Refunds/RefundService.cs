using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.Infrastructure.Notifications;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Infrastructure.Refunds;

public enum RefundReason
{
    Unspecified = 0,
    FreeLookCancellation = 1,
    CustomerRequest = 2,
    UnderwritingRejection = 3,
    DuplicatePolicy = 4,
    DeathOfInsured = 5,
    PolicyLapsed = 6,
    Fraud = 7,
    ProposalRejection = 8
}

public enum RefundStatus
{
    Unspecified = 0,
    Pending = 1,
    Calculating = 2,
    Approved = 3,
    Processing = 4,
    Completed = 5,
    Failed = 6,
    Rejected = 7
}

public class RefundCalculationResult
{
    public decimal TotalPremiumPaid { get; set; }
    public decimal PremiumUsed { get; set; }
    public decimal CancellationCharge { get; set; }
    public decimal RefundableAmount { get; set; }
    public string CalculationDetails { get; set; } = string.Empty;
    public int DaysUnused { get; set; }
    public int TotalPolicyDays { get; set; }
    public decimal UnusedDaysPercentage { get; set; }
}

public class RefundRequest
{
    public string PolicyId { get; set; } = string.Empty;
    public string Reason { get; set; } = string.Empty;
    public string ReasonDetails { get; set; } = string.Empty;
}

public interface IRefundService
{
    Task<RefundCalculationResult> CalculateProRataRefundAsync(
        string policyId,
        decimal totalPremiumPaid,
        DateTime policyStartDate,
        DateTime policyEndDate,
        CancellationToken ct = default);

    Task<string> RequestRefundAsync(RefundRequest request, CancellationToken ct = default);

    Task<RefundCalculationResult> GetRefundCalculationAsync(string refundId, CancellationToken ct = default);

    Task ApproveRefundAsync(string refundId, string approvedBy, string? comments = null, CancellationToken ct = default);

    Task ProcessRefundAsync(string refundId, string paymentMethod, string? paymentReference = null, CancellationToken ct = default);

    Task NotifyRefundStatusAsync(string userId, string refundId, RefundStatus status, decimal amount, CancellationToken ct = default);
}

public class RefundService : IRefundService
{
    private readonly InsuranceServiceClient _client;
    private readonly INotificationService _notificationService;
    private readonly ILogger<RefundService> _logger;

    private const decimal DefaultCancellationChargePercent = 10m;
    private const decimal FreeLookCancellationChargePercent = 0m;

    public RefundService(
        InsuranceServiceClient client,
        INotificationService notificationService,
        ILogger<RefundService> logger)
    {
        _client = client;
        _notificationService = notificationService;
        _logger = logger;
    }

    public async Task<RefundCalculationResult> CalculateProRataRefundAsync(
        string policyId,
        decimal totalPremiumPaid,
        DateTime policyStartDate,
        DateTime policyEndDate,
        CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Refund.Services.V1.CalculateRefundRequest
            {
                PolicyId = policyId,
                Reason = "CUSTOMER_REQUEST"
            };

            var response = await _client.Refunds.CalculateRefundAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogWarning("Go backend refund calculation failed, calculating locally: {Error}", 
                    response.Error.Message);
                
                return CalculateLocally(totalPremiumPaid, policyStartDate, policyEndDate, RefundReason.CustomerRequest);
            }

            var result = new RefundCalculationResult
            {
                TotalPremiumPaid = decimal.Parse(response.TotalPremiumPaid),
                PremiumUsed = decimal.Parse(response.PremiumUsed),
                CancellationCharge = decimal.Parse(response.CancellationCharge),
                RefundableAmount = response.RefundableAmount?.Amount ?? 0,
                CalculationDetails = response.CalculationDetails
            };

            _logger.LogInformation(
                "Refund calculated for policy {PolicyId}: Total={Total}, Used={Used}, Charge={Charge}, Refund={Refund}",
                policyId, result.TotalPremiumPaid, result.PremiumUsed, result.CancellationCharge, result.RefundableAmount);

            return result;
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Error calling Go refund service, calculating locally");
            return CalculateLocally(totalPremiumPaid, policyStartDate, policyEndDate, RefundReason.CustomerRequest);
        }
    }

    public async Task<string> RequestRefundAsync(RefundRequest request, CancellationToken ct = default)
    {
        try
        {
            var protoRequest = new Insuretech.Refund.Services.V1.RequestRefundRequest
            {
                PolicyId = request.PolicyId,
                Reason = MapReasonToProto(request.Reason),
                ReasonDetails = request.ReasonDetails
            };

            var response = await _client.Refunds.RequestRefundAsync(protoRequest, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                _logger.LogError("Failed to request refund: {Error}", response.Error.Message);
                throw new InvalidOperationException($"Refund request failed: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Refund requested: {RefundId} ({RefundNumber}) for policy {PolicyId}",
                response.RefundId, response.RefundNumber, request.PolicyId);

            return response.RefundId;
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error requesting refund for policy {PolicyId}", request.PolicyId);
            throw;
        }
    }

    public async Task<RefundCalculationResult> GetRefundCalculationAsync(string refundId, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Refund.Services.V1.GetRefundRequest
            {
                RefundId = refundId
            };

            var response = await _client.Refunds.GetRefundAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to get refund: {response.Error.Message}");
            }

            var refund = response.Refund;

            return new RefundCalculationResult
            {
                TotalPremiumPaid = refund?.TotalPremiumPaid?.Amount ?? 0,
                PremiumUsed = refund?.PremiumUsed?.Amount ?? 0,
                CancellationCharge = refund?.CancellationCharge?.Amount ?? 0,
                RefundableAmount = refund?.RefundableAmount?.Amount ?? 0,
                CalculationDetails = refund?.CalculationDetails ?? string.Empty
            };
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error getting refund {RefundId}", refundId);
            throw;
        }
    }

    public async Task ApproveRefundAsync(string refundId, string approvedBy, string? comments = null, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Refund.Services.V1.ApproveRefundRequest
            {
                RefundId = refundId,
                ApprovedBy = approvedBy,
                Comments = comments ?? string.Empty
            };

            var response = await _client.Refunds.ApproveRefundAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to approve refund: {response.Error.Message}");
            }

            _logger.LogInformation("Refund approved: {RefundId} by {ApprovedBy}", refundId, approvedBy);
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error approving refund {RefundId}", refundId);
            throw;
        }
    }

    public async Task ProcessRefundAsync(string refundId, string paymentMethod, string? paymentReference = null, CancellationToken ct = default)
    {
        try
        {
            var request = new Insuretech.Refund.Services.V1.ProcessRefundRequest
            {
                RefundId = refundId,
                PaymentMethod = paymentMethod,
                PaymentReference = paymentReference ?? string.Empty
            };

            var response = await _client.Refunds.ProcessRefundAsync(request, _client.BuildCallOptions(ct));

            if (response.Error != null && !string.IsNullOrEmpty(response.Error.Message))
            {
                throw new InvalidOperationException($"Failed to process refund: {response.Error.Message}");
            }

            _logger.LogInformation(
                "Refund processed: {RefundId}, Method={PaymentMethod}, Ref={PaymentReference}",
                refundId, paymentMethod, paymentReference);
        }
        catch (Exception ex) when (ex is not InvalidOperationException)
        {
            _logger.LogError(ex, "Error processing refund {RefundId}", refundId);
            throw;
        }
    }

    public async Task NotifyRefundStatusAsync(
        string userId,
        string refundId,
        RefundStatus status,
        decimal amount,
        CancellationToken ct = default)
    {
        var statusMessage = status switch
        {
            RefundStatus.Pending => $"Your refund request (ID: {refundId}) has been submitted and is pending approval.",
            RefundStatus.Approved => $"Your refund (ID: {refundId}) of BDT {amount:N2} has been approved and will be processed soon.",
            RefundStatus.Processing => $"Your refund (ID: {refundId}) of BDT {amount:N2} is being processed.",
            RefundStatus.Completed => $"Your refund (ID: {refundId}) of BDT {amount:N2} has been completed successfully.",
            RefundStatus.Rejected => $"Your refund request (ID: {refundId}) has been rejected. Please contact support for details.",
            _ => $"Refund update for {refundId}: Status changed to {status}."
        };

        await _notificationService.SendEmailAsync(
            userId,
            $"Refund Status Update - {refundId}",
            statusMessage,
            new Dictionary<string, string>
            {
                ["refund_id"] = refundId,
                ["status"] = status.ToString(),
                ["amount"] = amount.ToString("N2")
            },
            ct);

        _logger.LogInformation("Refund status notification sent to user {UserId}: {RefundId} -> {Status}",
            userId, refundId, status);
    }

    private RefundCalculationResult CalculateLocally(
        decimal totalPremiumPaid,
        DateTime policyStartDate,
        DateTime policyEndDate,
        RefundReason reason)
    {
        var totalDays = (policyEndDate - policyStartDate).Days;
        var daysElapsed = (DateTime.UtcNow - policyStartDate).Days;
        var daysUnused = Math.Max(0, totalDays - daysElapsed);
        var unusedDaysPercentage = totalDays > 0 ? (decimal)daysUnused / totalDays * 100 : 0;

        decimal cancellationChargePercent;
        if (reason == RefundReason.FreeLookCancellation)
        {
            cancellationChargePercent = FreeLookCancellationChargePercent;
        }
        else
        {
            cancellationChargePercent = DefaultCancellationChargePercent;
        }

        var premiumUsed = totalPremiumPaid - (totalPremiumPaid * unusedDaysPercentage / 100);
        var cancellationCharge = (totalPremiumPaid - premiumUsed) * cancellationChargePercent / 100;
        var refundableAmount = totalPremiumPaid - premiumUsed - cancellationCharge;

        if (refundableAmount < 0)
            refundableAmount = 0;

        var details = new
        {
            calculation_date = DateTime.UtcNow.ToString("O"),
            policy_start_date = policyStartDate.ToString("O"),
            policy_end_date = policyEndDate.ToString("O"),
            total_policy_days = totalDays,
            days_elapsed = daysElapsed,
            days_unused = daysUnused,
            unused_days_percentage = Math.Round(unusedDaysPercentage, 2),
            cancellation_charge_percent = cancellationChargePercent,
            formula = new
            {
                premium_used = $"total_premium × (days_unused / total_days) = {totalPremiumPaid} × ({daysUnused} / {totalDays})",
                cancellation_charge = $"(total_premium - premium_used) × charge_% = ({totalPremiumPaid} - {Math.Round(premiumUsed, 2)}) × {cancellationChargePercent}%",
                refundable = "total_premium - premium_used - cancellation_charge"
            }
        };

        var result = new RefundCalculationResult
        {
            TotalPremiumPaid = totalPremiumPaid,
            PremiumUsed = Math.Round(premiumUsed, 2),
            CancellationCharge = Math.Round(cancellationCharge, 2),
            RefundableAmount = Math.Round(refundableAmount, 2),
            CalculationDetails = System.Text.Json.JsonSerializer.Serialize(details),
            DaysUnused = daysUnused,
            TotalPolicyDays = totalDays,
            UnusedDaysPercentage = Math.Round(unusedDaysPercentage, 2)
        };

        _logger.LogInformation(
            "Local pro-rata calculation: Total={Total}, Used={Used}, Charge={Charge}, Refund={Refund} ({UnusedDays}% unused)",
            result.TotalPremiumPaid, result.PremiumUsed, result.CancellationCharge, 
            result.RefundableAmount, result.UnusedDaysPercentage);

        return result;
    }

    private static string MapReasonToProto(string reason)
    {
        return reason.ToUpper() switch
        {
            "FREELOOK" or "FREE_LOOK" or "FREELOOK_CANCELLATION" => "FREE_LOOK_CANCELLATION",
            "CUSTOMER_REQUEST" or "CUSTOMER" => "CUSTOMER_REQUEST",
            "UNDERWRITING_REJECTION" or "UNDERWRITING" => "UNDERWRITING_REJECTION",
            "DUPLICATE_POLICY" or "DUPLICATE" => "DUPLICATE_POLICY",
            "DEATH_OF_INSURED" or "DEATH" => "DEATH_OF_INSURED",
            "POLICY_LAPSED" or "LAPSED" => "POLICY_LAPSED",
            "FRAUD" => "FRAUD",
            "PROPOSAL_REJECTION" or "PROPOSAL" => "PROPOSAL_REJECTION",
            _ => "CUSTOMER_REQUEST"
        };
    }
}

public class MockRefundService : IRefundService
{
    private readonly INotificationService _notificationService;
    private readonly ILogger<MockRefundService> _logger;

    public MockRefundService(INotificationService notificationService, ILogger<MockRefundService> logger)
    {
        _notificationService = notificationService;
        _logger = logger;
    }

    public Task<RefundCalculationResult> CalculateProRataRefundAsync(
        string policyId,
        decimal totalPremiumPaid,
        DateTime policyStartDate,
        DateTime policyEndDate,
        CancellationToken ct = default)
    {
        var totalDays = (policyEndDate - policyStartDate).Days;
        var daysElapsed = Math.Max(1, (DateTime.UtcNow - policyStartDate).Days);
        var daysUnused = Math.Max(0, totalDays - daysElapsed);
        var unusedPercent = totalDays > 0 ? (decimal)daysUnused / totalDays * 100 : 100;

        var premiumUsed = totalPremiumPaid * (100 - unusedPercent) / 100;
        var cancellationCharge = (totalPremiumPaid - premiumUsed) * 0.10m;
        var refundable = totalPremiumPaid - premiumUsed - cancellationCharge;

        _logger.LogInformation(
            "[MOCK] Pro-rata refund calculated for {PolicyId}: Premium={Premium}, Used={Used}, Unused%={UnusedPct}%, Refund={Refund}",
            policyId, totalPremiumPaid, premiumUsed, unusedPercent, refundable);

        return Task.FromResult(new RefundCalculationResult
        {
            TotalPremiumPaid = totalPremiumPaid,
            PremiumUsed = Math.Round(premiumUsed, 2),
            CancellationCharge = Math.Round(cancellationCharge, 2),
            RefundableAmount = Math.Round(refundable, 2),
            CalculationDetails = $"{{\"method\": \"pro-rata\", \"days_unused\": {daysUnused}, \"unused_percent\": {unusedPercent}}}",
            DaysUnused = daysUnused,
            TotalPolicyDays = totalDays,
            UnusedDaysPercentage = Math.Round(unusedPercent, 2)
        });
    }

    public Task<string> RequestRefundAsync(RefundRequest request, CancellationToken ct = default)
    {
        var refundId = Guid.NewGuid().ToString();
        _logger.LogInformation(
            "[MOCK] Refund requested: {RefundId} for policy {PolicyId}, Reason: {Reason}",
            refundId, request.PolicyId, request.Reason);
        return Task.FromResult(refundId);
    }

    public Task<RefundCalculationResult> GetRefundCalculationAsync(string refundId, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Get refund calculation: {RefundId}", refundId);
        return Task.FromResult(new RefundCalculationResult
        {
            TotalPremiumPaid = 10000,
            PremiumUsed = 3000,
            CancellationCharge = 300,
            RefundableAmount = 6700,
            CalculationDetails = "{\"mock\": true}"
        });
    }

    public Task ApproveRefundAsync(string refundId, string approvedBy, string? comments = null, CancellationToken ct = default)
    {
        _logger.LogInformation("[MOCK] Refund approved: {RefundId} by {ApprovedBy}", refundId, approvedBy);
        return Task.CompletedTask;
    }

    public Task ProcessRefundAsync(string refundId, string paymentMethod, string? paymentReference = null, CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Refund processed: {RefundId}, Method={PaymentMethod}, Ref={PaymentReference}",
            refundId, paymentMethod, paymentReference);
        return Task.CompletedTask;
    }

    public Task NotifyRefundStatusAsync(string userId, string refundId, RefundStatus status, decimal amount, CancellationToken ct = default)
    {
        _logger.LogInformation(
            "[MOCK] Refund notification sent: User={UserId}, Refund={RefundId}, Status={Status}, Amount={Amount}",
            userId, refundId, status, amount);
        return Task.CompletedTask;
    }
}
