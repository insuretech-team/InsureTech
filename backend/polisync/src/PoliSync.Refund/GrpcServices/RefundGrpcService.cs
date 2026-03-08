using System.Collections.Concurrent;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Refund.Entity.V1;
using Insuretech.Refund.Services.V1;
using Insuretech.Payment.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Refund.Infrastructure;
using RefundEntity = Insuretech.Refund.Entity.V1.Refund;

namespace PoliSync.Refund.GrpcServices;

public sealed class RefundGrpcService : RefundService.RefundServiceBase
{
    private static readonly ConcurrentDictionary<string, RefundEntity> RefundsStore = new();
    private static readonly object MutationLock = new();

    private readonly IMediator _mediator;
    private readonly ILogger<RefundGrpcService> _logger;
    private readonly IRefundPaymentGateway _paymentGateway;

    public RefundGrpcService(
        IMediator mediator,
        ILogger<RefundGrpcService> logger,
        IRefundPaymentGateway paymentGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _paymentGateway = paymentGateway;
    }

    public override Task<RequestRefundResponse> RequestRefund(RequestRefundRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return Task.FromResult(new RequestRefundResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            });
        }

        var refundId = Guid.NewGuid().ToString("N");
        var refund = new RefundEntity
        {
            Id = refundId,
            RefundNumber = $"RFD-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
            PolicyId = request.PolicyId,
            Reason = ParseRefundReason(request.Reason),
            ReasonDetails = request.ReasonDetails,
            TotalPremiumPaid = NewMoney(0),
            PremiumUsed = NewMoney(0),
            CancellationCharge = NewMoney(0),
            RefundableAmount = NewMoney(0),
            CalculationDetails = string.Empty,
            Status = RefundStatus.Pending,
            RequestedBy = "SYSTEM"
        };

        RefundsStore[refundId] = refund;
        _logger.LogInformation("Refund requested: {RefundId}", refundId);

        return Task.FromResult(new RequestRefundResponse
        {
            RefundId = refund.Id,
            RefundNumber = refund.RefundNumber,
            Message = "Refund request submitted"
        });
    }

    public override Task<CalculateRefundResponse> CalculateRefund(CalculateRefundRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return Task.FromResult(new CalculateRefundResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            });
        }

        var reason = ParseRefundReason(request.Reason);
        var refund = RefundsStore.Values
            .Where(x => x.PolicyId == request.PolicyId)
            .OrderByDescending(x => x.ProcessedAt?.Seconds ?? 0)
            .FirstOrDefault();

        if (refund is null)
        {
            refund = new RefundEntity
            {
                Id = Guid.NewGuid().ToString("N"),
                RefundNumber = $"RFD-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
                PolicyId = request.PolicyId,
                Reason = reason,
                Status = RefundStatus.Pending,
                RequestedBy = "SYSTEM"
            };
            RefundsStore[refund.Id] = refund;
        }

        var totalPremiumPaid = 1_200_000L;
        var premiumUsed = reason switch
        {
            RefundReason.FreeLookCancellation => 0L,
            RefundReason.PolicyLapsed => 1_050_000L,
            RefundReason.DeathOfInsured => 400_000L,
            _ => 350_000L
        };
        var cancellationCharge = reason switch
        {
            RefundReason.Fraud => 300_000L,
            RefundReason.CustomerRequest => 100_000L,
            _ => 50_000L
        };
        var refundable = Math.Max(totalPremiumPaid - premiumUsed - cancellationCharge, 0);

        lock (MutationLock)
        {
            refund.Reason = reason;
            refund.TotalPremiumPaid = NewMoney(totalPremiumPaid);
            refund.PremiumUsed = NewMoney(premiumUsed);
            refund.CancellationCharge = NewMoney(cancellationCharge);
            refund.RefundableAmount = NewMoney(refundable);
            refund.CalculationDetails = $"total={totalPremiumPaid};used={premiumUsed};charge={cancellationCharge};reason={reason}";
            refund.Status = RefundStatus.Calculating;
        }

        return Task.FromResult(new CalculateRefundResponse
        {
            TotalPremiumPaid = totalPremiumPaid.ToString(),
            PremiumUsed = premiumUsed.ToString(),
            CancellationCharge = cancellationCharge.ToString(),
            RefundableAmount = refund.RefundableAmount,
            CalculationDetails = refund.CalculationDetails
        });
    }

    public override Task<GetRefundResponse> GetRefund(GetRefundRequest request, ServerCallContext context)
    {
        if (!RefundsStore.TryGetValue(request.RefundId, out var refund))
        {
            return Task.FromResult(new GetRefundResponse
            {
                Error = BuildError("NOT_FOUND", "Refund not found")
            });
        }

        return Task.FromResult(new GetRefundResponse
        {
            Refund = refund
        });
    }

    public override Task<ApproveRefundResponse> ApproveRefund(ApproveRefundRequest request, ServerCallContext context)
    {
        if (!RefundsStore.TryGetValue(request.RefundId, out var refund))
        {
            return Task.FromResult(new ApproveRefundResponse
            {
                Error = BuildError("NOT_FOUND", "Refund not found")
            });
        }

        lock (MutationLock)
        {
            refund.Status = RefundStatus.Approved;
            refund.ApprovedBy = request.ApprovedBy;
            if (!string.IsNullOrWhiteSpace(request.Comments))
            {
                refund.ReasonDetails = $"{refund.ReasonDetails};approval_comments={request.Comments}";
            }
        }

        return Task.FromResult(new ApproveRefundResponse
        {
            Message = "Refund approved"
        });
    }

    public override async Task<ProcessRefundResponse> ProcessRefund(ProcessRefundRequest request, ServerCallContext context)
    {
        if (!RefundsStore.TryGetValue(request.RefundId, out var refund))
        {
            return new ProcessRefundResponse
            {
                Error = BuildError("NOT_FOUND", "Refund not found")
            };
        }

        if (refund.Status != RefundStatus.Approved)
        {
            return new ProcessRefundResponse
            {
                Error = BuildError("INVALID_STATE", "Refund must be approved before processing")
            };
        }

        var paymentId = string.IsNullOrWhiteSpace(request.PaymentReference)
            ? refund.PaymentReference
            : request.PaymentReference;

        if (string.IsNullOrWhiteSpace(paymentId))
        {
            return new ProcessRefundResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PaymentReference is required to initiate payment refund")
            };
        }

        try
        {
            var paymentResponse = await _paymentGateway.InitiateRefundAsync(
                paymentId,
                refund.RefundableAmount ?? NewMoney(0),
                refund.Reason.ToString(),
                string.IsNullOrWhiteSpace(refund.ApprovedBy) ? "refund-service" : refund.ApprovedBy,
                context?.CancellationToken ?? CancellationToken.None);

            if (paymentResponse.Error is not null && !string.IsNullOrWhiteSpace(paymentResponse.Error.Code))
            {
                return new ProcessRefundResponse
                {
                    Error = BuildError(paymentResponse.Error.Code, paymentResponse.Error.Message)
                };
            }

            var processedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            lock (MutationLock)
            {
                refund.Status = RefundStatus.Completed;
                refund.PaymentMethod = request.PaymentMethod;
                refund.PaymentReference = paymentId;
                refund.PaymentRefundId = string.IsNullOrWhiteSpace(paymentResponse.RefundId)
                    ? $"PAYREF-{Guid.NewGuid():N}"[..18]
                    : paymentResponse.RefundId;
                refund.ProcessedAt = processedAt;
            }

            return new ProcessRefundResponse
            {
                Message = "Refund processed",
                ProcessedAt = processedAt.ToDateTime().ToString("O")
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to initiate payment refund for refund id {RefundId}", request.RefundId);
            return new ProcessRefundResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override Task<ListRefundsResponse> ListRefunds(ListRefundsRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

        var query = RefundsStore.Values.AsEnumerable();

        if (!string.IsNullOrWhiteSpace(request.BeneficiaryId))
        {
            query = query.Where(x => x.RequestedBy == request.BeneficiaryId);
        }

        var status = ParseRefundStatus(request.Status);
        if (status != RefundStatus.Unspecified)
        {
            query = query.Where(x => x.Status == status);
        }

        var ordered = query.OrderByDescending(x => x.ProcessedAt?.Seconds ?? 0).ToList();
        var pageItems = ordered.Skip((page - 1) * pageSize).Take(pageSize).ToList();

        var response = new ListRefundsResponse { TotalCount = ordered.Count };
        response.Refunds.AddRange(pageItems);
        return Task.FromResult(response);
    }

    private static RefundReason ParseRefundReason(string value)
    {
        return ParseEnum(value, RefundReason.CustomerRequest);
    }

    private static RefundStatus ParseRefundStatus(string value)
    {
        return ParseEnum(value, RefundStatus.Unspecified);
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

    private static Money NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };
}
