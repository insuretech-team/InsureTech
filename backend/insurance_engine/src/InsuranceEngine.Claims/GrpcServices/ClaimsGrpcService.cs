using System;
using System.Linq;
using System.Threading.Tasks;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.Claims.Application.Features.Commands.Claims;
using InsuranceEngine.Claims.Application.Features.Queries.Claims;
using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.GrpcServices;

public sealed class ClaimsGrpcService : ClaimService.ClaimServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ClaimsGrpcService> _logger;

    public ClaimsGrpcService(IMediator mediator, ILogger<ClaimsGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<SubmitClaimResponse> SubmitClaim(SubmitClaimRequest request, ServerCallContext context)
    {
        _logger.LogInformation("SubmitClaim gRPC call for Policy: {PolicyId}", request.PolicyId);

        var command = new SubmitClaimCommand(
            Guid.Parse(request.PolicyId),
            Guid.Parse(request.CustomerId),
            MapToDomainClaimType(request.Type),
            request.ClaimedAmount?.Amount ?? 0,
            request.IncidentDate?.ToDateTime() ?? DateTime.UtcNow,
            request.IncidentDescription,
            request.PlaceOfIncident
        );

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
        {
            return new SubmitClaimResponse
            {
                ClaimId = result.Value.ToString(),
                ClaimNumber = $"CLM-{DateTime.UtcNow:yyyyMMdd}-{result.Value.ToString()[..4]}", // Example
                Message = "Claim submitted successfully"
            };
        }

        return new SubmitClaimResponse
        {
            Error = new Insuretech.Common.V1.Error
            {
                Code = "SUBMIT_FAILED",
                Message = result.Error?.Message ?? "Unknown error"
            }
        };
    }

    public override async Task<GetClaimResponse> GetClaim(GetClaimRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new GetClaimByIdQuery(Guid.Parse(request.ClaimId)));

        if (result.IsSuccess)
        {
            return new GetClaimResponse
            {
                Claim = MapToProtoClaim(result.Value)
            };
        }

        return new GetClaimResponse
        {
            Error = new Insuretech.Common.V1.Error
            {
                Code = "NOT_FOUND",
                Message = result.Error?.Message ?? "Claim not found"
            }
        };
    }

    public override async Task<ListUserClaimsResponse> ListUserClaims(ListUserClaimsRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 10 : request.PageSize;

        var result = await _mediator.Send(new ListClaimsByCustomerQuery(Guid.Parse(request.CustomerId), page, pageSize));

        if (result.IsSuccess)
        {
            var response = new ListUserClaimsResponse
            {
                TotalCount = result.Value.TotalCount
            };
            response.Claims.AddRange(result.Value.Items.Select(MapToProtoClaim));
            return response;
        }

        return new ListUserClaimsResponse
        {
            Error = new Insuretech.Common.V1.Error
            {
                Code = "LIST_FAILED",
                Message = result.Error?.Message ?? "Unknown error"
            }
        };
    }

    private static Claim MapToProtoClaim(ClaimResponseDto dto)
    {
        return new Claim
        {
            ClaimId = dto.Id.ToString(),
            ClaimNumber = dto.ClaimNumber,
            PolicyId = dto.PolicyId.ToString(),
            CustomerId = dto.CustomerId.ToString(),
            Type = MapToProtoClaimType(dto.Type),
            Status = MapToProtoClaimStatus(dto.Status),
            ClaimedAmount = new Insuretech.Common.V1.Money { Amount = (long)(dto.ClaimedAmount * 100), Currency = "BDT" },
            ApprovedAmount = new Insuretech.Common.V1.Money { Amount = (long)(dto.ApprovedAmount * 100), Currency = "BDT" },
            IncidentDate = Timestamp.FromDateTime(DateTime.SpecifyKind(dto.IncidentDate, DateTimeKind.Utc)),
            IncidentDescription = dto.IncidentDescription,
            SubmittedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(dto.SubmittedAt, DateTimeKind.Utc)),
            ProcessedAt = dto.ProcessedAt.HasValue ? Timestamp.FromDateTime(DateTime.SpecifyKind(dto.ProcessedAt.Value, DateTimeKind.Utc)) : null
        };
    }

    private static InsuranceEngine.Claims.Domain.Enums.ClaimType MapToDomainClaimType(ClaimType type) => type switch
    {
        ClaimType.Death => InsuranceEngine.Claims.Domain.Enums.ClaimType.Death,
        ClaimType.Maturity => InsuranceEngine.Claims.Domain.Enums.ClaimType.Maturity,
        ClaimType.Surrender => InsuranceEngine.Claims.Domain.Enums.ClaimType.Surrender,
        ClaimType.Health => InsuranceEngine.Claims.Domain.Enums.ClaimType.Health,
        _ => InsuranceEngine.Claims.Domain.Enums.ClaimType.Unspecified
    };

    private static ClaimType MapToProtoClaimType(InsuranceEngine.Claims.Domain.Enums.ClaimType type) => type switch
    {
        InsuranceEngine.Claims.Domain.Enums.ClaimType.Death => ClaimType.Death,
        InsuranceEngine.Claims.Domain.Enums.ClaimType.Maturity => ClaimType.Maturity,
        InsuranceEngine.Claims.Domain.Enums.ClaimType.Surrender => ClaimType.Surrender,
        InsuranceEngine.Claims.Domain.Enums.ClaimType.Health => ClaimType.Health,
        _ => ClaimType.Unspecified
    };

    private static ClaimStatus MapToProtoClaimStatus(InsuranceEngine.Claims.Domain.Enums.ClaimStatus status) => status switch
    {
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.Submitted => ClaimStatus.Submitted,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.PendingDocuments => ClaimStatus.PendingDocuments,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.UnderInvestigation => ClaimStatus.UnderInvestigation,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.Approved => ClaimStatus.Approved,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.Rejected => ClaimStatus.Rejected,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.Settled => ClaimStatus.Settled,
        InsuranceEngine.Claims.Domain.Enums.ClaimStatus.Disputed => ClaimStatus.Disputed,
        _ => ClaimStatus.Unspecified
    };
}
