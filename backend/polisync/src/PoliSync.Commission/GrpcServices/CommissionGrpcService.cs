using System.Collections.Concurrent;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Commission.Entity.V1;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using Insuretech.Partner.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PartnerCommission = Insuretech.Partner.Entity.V1.Commission;

namespace PoliSync.Commission.GrpcServices;

public sealed class CommissionGrpcService : CommissionService.CommissionServiceBase
{
    private static readonly ConcurrentDictionary<string, PartnerCommission> CommissionsStore = new();
    private static readonly ConcurrentDictionary<string, CommissionPayout> PayoutStore = new();
    private static readonly ConcurrentDictionary<string, List<string>> PayoutCommissionMap = new();

    private readonly IMediator _mediator;
    private readonly ILogger<CommissionGrpcService> _logger;

    public CommissionGrpcService(IMediator mediator, ILogger<CommissionGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override Task<CalculateCommissionResponse> CalculateCommission(CalculateCommissionRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId) || string.IsNullOrWhiteSpace(request.RecipientId))
        {
            return Task.FromResult(new CalculateCommissionResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId and RecipientId are required")
            });
        }

        var type = ParseCommissionType(request.CommissionType);
        var rate = type switch
        {
            CommissionType.Acquisition => 0.15,
            CommissionType.Renewal => 0.05,
            CommissionType.ClaimsAssistance => 0.02,
            _ => 0.03
        };

        const long basePremiumPaisa = 1_000_000; // 10,000 BDT reference baseline
        var amount = (long)Math.Round(basePremiumPaisa * rate, MidpointRounding.AwayFromZero);
        var now = Timestamp.FromDateTime(DateTime.UtcNow);

        var commission = new PartnerCommission
        {
            CommissionId = Guid.NewGuid().ToString("N"),
            PolicyId = request.PolicyId,
            PartnerId = request.RecipientType.Equals("partner", StringComparison.OrdinalIgnoreCase) ? request.RecipientId : string.Empty,
            AgentId = request.RecipientType.Equals("agent", StringComparison.OrdinalIgnoreCase) ? request.RecipientId : string.Empty,
            Type = type,
            CommissionAmount = NewMoney(amount),
            CommissionRate = rate,
            Status = CommissionStatus.Pending,
            CreatedAt = now,
            UpdatedAt = now
        };

        CommissionsStore[commission.CommissionId] = commission;
        _logger.LogInformation("Commission calculated: {CommissionId}", commission.CommissionId);

        return Task.FromResult(new CalculateCommissionResponse
        {
            CommissionId = commission.CommissionId,
            CommissionNumber = BuildCommissionNumber(commission.CommissionId),
            Amount = commission.CommissionAmount,
            CalculationBreakdown = $"base={basePremiumPaisa}; rate={rate:P2}; amount={amount}"
        });
    }

    public override Task<GetCommissionResponse> GetCommission(GetCommissionRequest request, ServerCallContext context)
    {
        if (!CommissionsStore.TryGetValue(request.CommissionId, out var commission))
        {
            return Task.FromResult(new GetCommissionResponse { Error = BuildError("NOT_FOUND", "Commission not found") });
        }

        return Task.FromResult(new GetCommissionResponse { Commission = commission });
    }

    public override Task<ListCommissionsResponse> ListCommissions(ListCommissionsRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

        var query = CommissionsStore.Values.AsEnumerable();
        if (!string.IsNullOrWhiteSpace(request.RecipientId))
        {
            query = query.Where(c => c.PartnerId == request.RecipientId || c.AgentId == request.RecipientId);
        }

        var status = ParseCommissionStatus(request.Status);
        if (status != CommissionStatus.Unspecified)
        {
            query = query.Where(c => c.Status == status);
        }

        var ordered = query.OrderByDescending(c => c.CreatedAt?.Seconds ?? 0).ToList();
        var paged = ordered.Skip((page - 1) * pageSize).Take(pageSize).ToList();

        var totalAmount = ordered.Sum(c => c.CommissionAmount?.Amount ?? 0);
        var response = new ListCommissionsResponse
        {
            TotalCount = ordered.Count,
            TotalAmount = NewMoney(totalAmount)
        };
        response.Commissions.AddRange(paged);
        return Task.FromResult(response);
    }

    public override Task<CreatePayoutResponse> CreatePayout(CreatePayoutRequest request, ServerCallContext context)
    {
        var selectedIds = request.CommissionIds.Count > 0
            ? request.CommissionIds.ToHashSet(StringComparer.Ordinal)
            : null;

        var commissions = CommissionsStore.Values
            .Where(c => c.Status == CommissionStatus.Pending)
            .Where(c => string.IsNullOrWhiteSpace(request.RecipientId) || c.PartnerId == request.RecipientId || c.AgentId == request.RecipientId)
            .Where(c => selectedIds == null || selectedIds.Contains(c.CommissionId))
            .ToList();

        if (commissions.Count == 0)
        {
            return Task.FromResult(new CreatePayoutResponse { Error = BuildError("NOT_FOUND", "No commissions eligible for payout") });
        }

        var payoutId = Guid.NewGuid().ToString("N");
        var payout = new CommissionPayout
        {
            Id = payoutId,
            PayoutNumber = BuildPayoutNumber(payoutId),
            RecipientType = request.RecipientType,
            RecipientId = request.RecipientId,
            PeriodStart = ParseDateToTimestamp(request.PeriodStart),
            PeriodEnd = ParseDateToTimestamp(request.PeriodEnd),
            TotalAmount = NewMoney(commissions.Sum(c => c.CommissionAmount.Amount)),
            CommissionCount = commissions.Count,
            Status = PayoutStatus.Pending
        };

        PayoutStore[payoutId] = payout;
        PayoutCommissionMap[payoutId] = commissions.Select(c => c.CommissionId).ToList();

        foreach (var commission in commissions)
        {
            commission.Status = CommissionStatus.Processing;
            commission.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        }

        return Task.FromResult(new CreatePayoutResponse
        {
            PayoutId = payout.Id,
            PayoutNumber = payout.PayoutNumber,
            TotalAmount = payout.TotalAmount,
            CommissionCount = payout.CommissionCount
        });
    }

    public override Task<ProcessPayoutResponse> ProcessPayout(ProcessPayoutRequest request, ServerCallContext context)
    {
        if (!PayoutStore.TryGetValue(request.PayoutId, out var payout))
        {
            return Task.FromResult(new ProcessPayoutResponse { Error = BuildError("NOT_FOUND", "Payout not found") });
        }

        var paidAt = Timestamp.FromDateTime(DateTime.UtcNow);
        payout.Status = PayoutStatus.Paid;
        payout.PaymentMethod = request.PaymentMethod;
        payout.PaymentReference = request.PaymentReference;
        payout.PaidAt = paidAt;

        if (PayoutCommissionMap.TryGetValue(request.PayoutId, out var commissionIds))
        {
            foreach (var id in commissionIds)
            {
                if (CommissionsStore.TryGetValue(id, out var commission))
                {
                    commission.Status = CommissionStatus.Paid;
                    commission.PaymentId = request.PayoutId;
                    commission.PaidAt = paidAt;
                    commission.UpdatedAt = paidAt;
                }
            }
        }

        return Task.FromResult(new ProcessPayoutResponse
        {
            Message = "Payout processed",
            PaidAt = paidAt.ToDateTime().ToString("O")
        });
    }

    public override Task<GetCommissionStatementResponse> GetCommissionStatement(GetCommissionStatementRequest request, ServerCallContext context)
    {
        var commissions = CommissionsStore.Values
            .Where(c => c.PartnerId == request.RecipientId || c.AgentId == request.RecipientId)
            .ToList();

        var totalEarned = commissions.Sum(c => c.CommissionAmount?.Amount ?? 0);
        var totalPaid = commissions.Where(c => c.Status == CommissionStatus.Paid).Sum(c => c.CommissionAmount?.Amount ?? 0);

        var response = new GetCommissionStatementResponse
        {
            RecipientId = request.RecipientId,
            PeriodStart = request.PeriodStart,
            PeriodEnd = request.PeriodEnd,
            TotalEarned = NewMoney(totalEarned),
            TotalPaid = NewMoney(totalPaid),
            PendingAmount = NewMoney(totalEarned - totalPaid)
        };

        var byType = commissions
            .GroupBy(c => c.Type)
            .Select(g => new CommissionSummary
            {
                Type = g.Key.ToString(),
                Count = g.Count(),
                TotalAmount = NewMoney(g.Sum(c => c.CommissionAmount?.Amount ?? 0))
            });

        response.ByType.AddRange(byType);
        return Task.FromResult(response);
    }

    public override Task<GetRevenueShareReportResponse> GetRevenueShareReport(GetRevenueShareReportRequest request, ServerCallContext context)
    {
        var commissions = CommissionsStore.Values.Where(c => !string.IsNullOrWhiteSpace(c.PolicyId)).ToList();
        var gross = commissions.Sum(c => c.CommissionAmount?.Amount ?? 0) * 5; // proxy gross estimate
        var platformShare = (long)(gross * 0.20);
        var insurerShare = gross - platformShare;

        var response = new GetRevenueShareReportResponse
        {
            InsurerId = request.InsurerId,
            TotalGrossPremium = NewMoney(gross),
            TotalPlatformShare = NewMoney(platformShare),
            TotalInsurerShare = NewMoney(insurerShare),
            PolicyCount = commissions.Select(c => c.PolicyId).Distinct().Count()
        };

        response.ByRevenueModel["DEFAULT"] = new RevenueModelAmount
        {
            Amount = NewMoney(platformShare)
        };

        return Task.FromResult(response);
    }

    private static CommissionType ParseCommissionType(string value)
    {
        return ParseEnum(value, CommissionType.Unspecified);
    }

    private static CommissionStatus ParseCommissionStatus(string value)
    {
        return ParseEnum(value, CommissionStatus.Unspecified);
    }

    private static TEnum ParseEnum<TEnum>(string value, TEnum fallback) where TEnum : struct, System.Enum
    {
        if (string.IsNullOrWhiteSpace(value))
        {
            return fallback;
        }

        var token = value.Trim();
        if (System.Enum.TryParse<TEnum>(token, true, out var direct))
        {
            return direct;
        }

        var parts = token.Split('_', StringSplitOptions.RemoveEmptyEntries);
        for (var i = 0; i < parts.Length; i++)
        {
            var candidate = string.Concat(parts.Skip(i).Select(ToPascalCase));
            if (System.Enum.TryParse<TEnum>(candidate, true, out var parsed))
            {
                return parsed;
            }
        }

        return fallback;
    }

    private static string ToPascalCase(string segment)
    {
        if (string.IsNullOrWhiteSpace(segment)) return string.Empty;
        var lower = segment.ToLowerInvariant();
        return char.ToUpperInvariant(lower[0]) + lower[1..];
    }

    private static string BuildCommissionNumber(string commissionId) => $"COM-{DateTime.UtcNow:yyyyMMdd}-{commissionId[..8].ToUpperInvariant()}";
    private static string BuildPayoutNumber(string payoutId) => $"PYO-{DateTime.UtcNow:yyyyMMdd}-{payoutId[..8].ToUpperInvariant()}";

    private static Timestamp ParseDateToTimestamp(string value)
    {
        return DateTime.TryParse(value, out var parsed)
            ? Timestamp.FromDateTime(DateTime.SpecifyKind(parsed, DateTimeKind.Utc))
            : Timestamp.FromDateTime(DateTime.UtcNow);
    }

    private static Money NewMoney(long amount, string currency = "BDT") => new() { Amount = amount, Currency = currency };
    private static Error BuildError(string code, string message) => new() { Code = code, Message = message };
}
