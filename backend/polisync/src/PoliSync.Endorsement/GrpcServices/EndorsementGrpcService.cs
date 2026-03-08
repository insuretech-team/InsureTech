using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Endorsement.Entity.V1;
using Insuretech.Endorsement.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.GrpcServices;

public sealed class EndorsementGrpcService : EndorsementService.EndorsementServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<EndorsementGrpcService> _logger;
    private readonly IEndorsementDataGateway _dataGateway;

    public EndorsementGrpcService(
        IMediator mediator,
        ILogger<EndorsementGrpcService> logger,
        IEndorsementDataGateway dataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _dataGateway = dataGateway;
    }

    public override async Task<RequestEndorsementResponse> RequestEndorsement(RequestEndorsementRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new RequestEndorsementResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        var endorsementType = ParseEndorsementType(request.Type);
        if (endorsementType == EndorsementType.Unspecified)
        {
            return new RequestEndorsementResponse
            {
                Error = BuildError("VALIDATION_ERROR", "Type is required and must be valid")
            };
        }

        var endorsementId = Guid.NewGuid().ToString("N");
        var now = DateTime.UtcNow;
        var endorsement = new EndorsementEntity
        {
            Id = endorsementId,
            EndorsementNumber = $"END-{now:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
            PolicyId = request.PolicyId,
            Type = endorsementType,
            Reason = request.Reason,
            Changes = request.Changes,
            PremiumAdjustment = CalculatePremiumAdjustment(endorsementType),
            PremiumRefundRequired = endorsementType == EndorsementType.RiderRemoval,
            Status = EndorsementStatus.Pending,
            RequestedBy = "SYSTEM",
            EffectiveDate = ParseDateOrDefault(request.EffectiveDate, now.Date.AddDays(1))
        };

        try
        {
            var created = await _dataGateway.CreateEndorsementAsync(endorsement, GetCancellationToken(context));
            _logger.LogInformation("Endorsement requested: {EndorsementId}", created.Id);

            return new RequestEndorsementResponse
            {
                EndorsementId = created.Id,
                EndorsementNumber = created.EndorsementNumber,
                Message = "Endorsement request submitted"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to create endorsement for policy {PolicyId}", request.PolicyId);
            return new RequestEndorsementResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetEndorsementResponse> GetEndorsement(GetEndorsementRequest request, ServerCallContext context)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, GetCancellationToken(context));
            if (endorsement is null)
            {
                return new GetEndorsementResponse
                {
                    Error = BuildError("NOT_FOUND", "Endorsement not found")
                };
            }

            return new GetEndorsementResponse
            {
                Endorsement = endorsement
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get endorsement {EndorsementId}", request.EndorsementId);
            return new GetEndorsementResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ListEndorsementsResponse> ListEndorsements(ListEndorsementsRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId))
        {
            return new ListEndorsementsResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId is required")
            };
        }

        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;
        var status = ParseEndorsementStatus(request.Status);

        try
        {
            var endorsements = await _dataGateway.ListEndorsementsByPolicyAsync(request.PolicyId, GetCancellationToken(context));

            var filtered = status == EndorsementStatus.Unspecified
                ? endorsements
                : endorsements.Where(x => x.Status == status).ToList();

            var ordered = filtered.OrderByDescending(x => x.EffectiveDate?.Seconds ?? 0).ToList();
            var pageItems = ordered.Skip((page - 1) * pageSize).Take(pageSize).ToList();

            var response = new ListEndorsementsResponse { TotalCount = ordered.Count };
            response.Endorsements.AddRange(pageItems);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list endorsements for policy {PolicyId}", request.PolicyId);
            return new ListEndorsementsResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ApproveEndorsementResponse> ApproveEndorsement(ApproveEndorsementRequest request, ServerCallContext context)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, GetCancellationToken(context));
            if (endorsement is null)
            {
                return new ApproveEndorsementResponse
                {
                    Error = BuildError("NOT_FOUND", "Endorsement not found")
                };
            }

            endorsement.Status = EndorsementStatus.Applied;
            endorsement.ApprovedBy = request.ApprovedBy;
            endorsement.ApprovedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            if (!string.IsNullOrWhiteSpace(request.Comments))
            {
                endorsement.Changes = $"{endorsement.Changes};comments={request.Comments}";
            }

            await _dataGateway.UpdateEndorsementAsync(endorsement, GetCancellationToken(context));

            return new ApproveEndorsementResponse
            {
                Message = "Endorsement approved and applied"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to approve endorsement {EndorsementId}", request.EndorsementId);
            return new ApproveEndorsementResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RejectEndorsementResponse> RejectEndorsement(RejectEndorsementRequest request, ServerCallContext context)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, GetCancellationToken(context));
            if (endorsement is null)
            {
                return new RejectEndorsementResponse
                {
                    Error = BuildError("NOT_FOUND", "Endorsement not found")
                };
            }

            endorsement.Status = EndorsementStatus.Rejected;
            endorsement.Reason = request.Reason;
            await _dataGateway.UpdateEndorsementAsync(endorsement, GetCancellationToken(context));

            return new RejectEndorsementResponse
            {
                Message = "Endorsement rejected"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to reject endorsement {EndorsementId}", request.EndorsementId);
            return new RejectEndorsementResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    private static EndorsementType ParseEndorsementType(string value)
    {
        return ParseEnum(value, EndorsementType.Unspecified);
    }

    private static EndorsementStatus ParseEndorsementStatus(string value)
    {
        return ParseEnum(value, EndorsementStatus.Unspecified);
    }

    private static Timestamp ParseDateOrDefault(string input, DateTime fallbackUtc)
    {
        if (DateTime.TryParse(input, out var parsed))
        {
            return Timestamp.FromDateTime(DateTime.SpecifyKind(parsed, DateTimeKind.Utc));
        }

        return Timestamp.FromDateTime(DateTime.SpecifyKind(fallbackUtc, DateTimeKind.Utc));
    }

    private static Money CalculatePremiumAdjustment(EndorsementType type)
    {
        var delta = type switch
        {
            EndorsementType.SumAssuredChange => 15_000L,
            EndorsementType.PremiumAdjustment => 10_000L,
            EndorsementType.RiderAddition => 5_000L,
            EndorsementType.RiderRemoval => -5_000L,
            _ => 0L
        };

        return new Money { Amount = delta, Currency = "BDT" };
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

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };
}
