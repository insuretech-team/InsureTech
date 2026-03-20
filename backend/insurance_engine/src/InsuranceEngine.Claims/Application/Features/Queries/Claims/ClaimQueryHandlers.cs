using MediatR;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.SharedKernel.CQRS;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Claims.Application.Features.Queries.Claims;

public class ClaimQueryHandlers : 
    IRequestHandler<GetClaimByIdQuery, Result<ClaimResponseDto>>,
    IRequestHandler<ListClaimsByCustomerQuery, Result<PaginatedResult<ClaimResponseDto>>>
{
    private readonly IClaimsRepository _claimsRepository;

    public ClaimQueryHandlers(IClaimsRepository claimsRepository)
    {
        _claimsRepository = claimsRepository;
    }

    public async Task<Result<ClaimResponseDto>> Handle(GetClaimByIdQuery request, CancellationToken cancellationToken)
    {
        var claim = await _claimsRepository.GetByIdAsync(request.Id, cancellationToken);
        if (claim == null)
            return Result<ClaimResponseDto>.Fail(Error.NotFound("Claim", request.Id.ToString()));

        return Result<ClaimResponseDto>.Success(MapToDto(claim));
    }

    public async Task<Result<PaginatedResult<ClaimResponseDto>>> Handle(ListClaimsByCustomerQuery request, CancellationToken cancellationToken)
    {
        var claims = await _claimsRepository.ListByCustomerAsync(request.CustomerId, request.Page, request.PageSize, cancellationToken);
        var total = await _claimsRepository.GetTotalCountByCustomerAsync(request.CustomerId, cancellationToken);
        
        return Result<PaginatedResult<ClaimResponseDto>>.Success(new PaginatedResult<ClaimResponseDto>(
            claims.Select(MapToDto).ToList(),
            total,
            request.Page,
            request.PageSize
        ));
    }

    private static ClaimResponseDto MapToDto(Claim claim)
    {
        return new ClaimResponseDto
        {
            Id = claim.Id,
            ClaimNumber = claim.ClaimNumber,
            PolicyId = claim.PolicyId,
            CustomerId = claim.CustomerId,
            Status = claim.Status,
            Type = claim.Type,
            ClaimedAmount = (decimal)claim.ClaimedAmount / 100,
            ApprovedAmount = (decimal)claim.ApprovedAmount / 100,
            Currency = claim.ClaimedCurrency,
            IncidentDate = claim.IncidentDate,
            IncidentDescription = claim.IncidentDescription,
            PlaceOfIncident = claim.PlaceOfIncident ?? "",
            SubmittedAt = claim.CreatedAt,
            RejectionReason = claim.RejectionReason,
            Approvals = claim.Approvals.Select(a => new ClaimApprovalDto
            {
                Id = a.Id,
                Decision = a.Decision.ToString(),
                Level = a.ApprovalLevel,
                Notes = a.Notes,
                DecidedAt = a.CreatedAt
            }).ToList(),
            Documents = claim.Documents.Select(d => new ClaimDocumentDto
            {
                Id = d.Id,
                DocumentType = d.DocumentType,
                FileUrl = d.FileUrl,
                IsVerified = d.Verified,
                UploadedAt = d.CreatedAt
            }).ToList()
        };
    }
}
