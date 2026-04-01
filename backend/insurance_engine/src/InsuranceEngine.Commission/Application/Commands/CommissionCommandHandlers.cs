using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Commission.Application.Commands;

// ===== CalculateCommission =====
public sealed class CalculateCommissionCommandHandler : IRequestHandler<CalculateCommissionCommand, CalculateCommissionResponse>
{
    private readonly IRepository<CommissionEntity> _commissionRepository;
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly ILogger<CalculateCommissionCommandHandler> _logger;

    public CalculateCommissionCommandHandler(
        IRepository<CommissionEntity> commissionRepository,
        IRepository<PolicyEntity> policyRepository,
        ILogger<CalculateCommissionCommandHandler> logger)
    {
        _commissionRepository = commissionRepository;
        _policyRepository = policyRepository;
        _logger = logger;
    }

    public async Task<CalculateCommissionResponse> Handle(CalculateCommissionCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _policyRepository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
                return new CalculateCommissionResponse { Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" } };

            // Commission rate based on type and recipient
            var rate = DetermineRate(request.CommissionType, request.RecipientType);
            var commissionAmount = (long)(policy.PremiumAmount * rate);

            var breakdown = System.Text.Json.JsonSerializer.Serialize(new
            {
                premiumAmount = policy.PremiumAmount,
                rate,
                commissionType = request.CommissionType,
                recipientType = request.RecipientType,
                calculatedAmount = commissionAmount
            });

            var entity = new CommissionEntity
            {
                CommissionId = Guid.NewGuid(),
                CommissionNumber = $"COM-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}",
                PolicyId = Guid.Parse(request.PolicyId),
                CommissionType = request.CommissionType,
                PartnerId = request.RecipientType == "PARTNER" ? Guid.Parse(request.RecipientId) : null,
                AgentId = request.RecipientType == "AGENT" ? Guid.Parse(request.RecipientId) : null,
                CommissionRate = rate,
                CommissionAmount = commissionAmount,
                CommissionCurrency = "BDT",
                CalculationBreakdown = breakdown,
                Status = "PENDING",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _commissionRepository.AddAsync(entity, cancellationToken);

            _logger.LogInformation("Commission calculated: {CommissionNumber}, Amount: {Amount}", entity.CommissionNumber, commissionAmount);

            return new CalculateCommissionResponse
            {
                CommissionId = entity.CommissionId.ToString(),
                CommissionNumber = entity.CommissionNumber,
                Amount = new Money { Amount = commissionAmount, Currency = "BDT" },
                CalculationBreakdown = breakdown
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to calculate commission");
            return new CalculateCommissionResponse { Error = new Error { Code = "CALC_FAILED", Message = ex.Message } };
        }
    }

    private static decimal DetermineRate(string commissionType, string recipientType) => (commissionType, recipientType) switch
    {
        ("ACQUISITION", "AGENT") => 0.15m,
        ("ACQUISITION", "PARTNER") => 0.10m,
        ("RENEWAL", "AGENT") => 0.07m,
        ("RENEWAL", "PARTNER") => 0.05m,
        _ => 0.10m
    };
}

// ===== CreatePayout =====
public sealed class CreatePayoutCommandHandler : IRequestHandler<CreatePayoutCommand, CreatePayoutResponse>
{
    private readonly IRepository<CommissionEntity> _commissionRepository;
    private readonly IRepository<CommissionPayoutEntity> _payoutRepository;
    private readonly ILogger<CreatePayoutCommandHandler> _logger;

    public CreatePayoutCommandHandler(
        IRepository<CommissionEntity> commissionRepository,
        IRepository<CommissionPayoutEntity> payoutRepository,
        ILogger<CreatePayoutCommandHandler> logger)
    {
        _commissionRepository = commissionRepository;
        _payoutRepository = payoutRepository;
        _logger = logger;
    }

    public async Task<CreatePayoutResponse> Handle(CreatePayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var recipientId = Guid.Parse(request.RecipientId);
            List<CommissionEntity> commissions;

            if (request.CommissionIds != null && request.CommissionIds.Count > 0)
            {
                var ids = request.CommissionIds.Select(Guid.Parse).ToList();
                commissions = await _commissionRepository.FindAsync(c => ids.Contains(c.CommissionId) && c.Status == "PENDING", cancellationToken);
            }
            else
            {
                commissions = await _commissionRepository.FindAsync(
                    c => (c.PartnerId == recipientId || c.AgentId == recipientId) && c.Status == "PENDING" && c.DeletedAt == null, cancellationToken);
            }

            if (commissions.Count == 0)
                return new CreatePayoutResponse { Error = new Error { Code = "NO_COMMISSIONS", Message = "No pending commissions found" } };

            var totalAmount = commissions.Sum(c => c.CommissionAmount);
            var payout = new CommissionPayoutEntity
            {
                PayoutId = Guid.NewGuid(),
                PayoutNumber = $"PAY-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}",
                RecipientType = request.RecipientType,
                RecipientId = recipientId,
                TotalAmount = totalAmount,
                CommissionCount = commissions.Count,
                PeriodStart = DateTime.Parse(request.PeriodStart),
                PeriodEnd = DateTime.Parse(request.PeriodEnd),
                Status = "PENDING",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _payoutRepository.AddAsync(payout, cancellationToken);

            // Link commissions to payout
            foreach (var c in commissions)
            {
                c.PayoutId = payout.PayoutId;
                c.UpdatedAt = DateTime.UtcNow;
                await _commissionRepository.UpdateAsync(c, cancellationToken);
            }

            return new CreatePayoutResponse
            {
                PayoutId = payout.PayoutId.ToString(),
                PayoutNumber = payout.PayoutNumber,
                TotalAmount = new Money { Amount = totalAmount, Currency = "BDT" },
                CommissionCount = commissions.Count
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create payout");
            return new CreatePayoutResponse { Error = new Error { Code = "PAYOUT_FAILED", Message = ex.Message } };
        }
    }
}

// ===== ProcessPayout =====
public sealed class ProcessPayoutCommandHandler : IRequestHandler<ProcessPayoutCommand, ProcessPayoutResponse>
{
    private readonly IRepository<CommissionPayoutEntity> _payoutRepository;
    private readonly IRepository<CommissionEntity> _commissionRepository;
    private readonly ILogger<ProcessPayoutCommandHandler> _logger;

    public ProcessPayoutCommandHandler(
        IRepository<CommissionPayoutEntity> payoutRepository,
        IRepository<CommissionEntity> commissionRepository,
        ILogger<ProcessPayoutCommandHandler> logger)
    {
        _payoutRepository = payoutRepository;
        _commissionRepository = commissionRepository;
        _logger = logger;
    }

    public async Task<ProcessPayoutResponse> Handle(ProcessPayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var payout = await _payoutRepository.GetByIdAsync(Guid.Parse(request.PayoutId), cancellationToken);
            if (payout == null)
                return new ProcessPayoutResponse { Error = new Error { Code = "PAYOUT_NOT_FOUND", Message = "Payout not found" } };

            if (payout.Status != "PENDING")
                return new ProcessPayoutResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Payout cannot be processed from status '{payout.Status}'" } };

            payout.Status = "PROCESSED";
            payout.PaymentMethod = request.PaymentMethod;
            payout.PaymentReference = request.PaymentReference;
            payout.PaidAt = DateTime.UtcNow;
            payout.UpdatedAt = DateTime.UtcNow;
            await _payoutRepository.UpdateAsync(payout, cancellationToken);

            // Mark linked commissions as PAID
            var commissions = await _commissionRepository.FindAsync(c => c.PayoutId == payout.PayoutId, cancellationToken);
            foreach (var c in commissions)
            {
                c.Status = "PAID";
                c.PaidAt = DateTime.UtcNow;
                c.UpdatedAt = DateTime.UtcNow;
                await _commissionRepository.UpdateAsync(c, cancellationToken);
            }

            return new ProcessPayoutResponse
            {
                Message = "Payout processed successfully",
                PaidAt = DateTime.UtcNow.ToString("yyyy-MM-ddTHH:mm:ssZ")
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process payout");
            return new ProcessPayoutResponse { Error = new Error { Code = "PROCESS_FAILED", Message = ex.Message } };
        }
    }
}
