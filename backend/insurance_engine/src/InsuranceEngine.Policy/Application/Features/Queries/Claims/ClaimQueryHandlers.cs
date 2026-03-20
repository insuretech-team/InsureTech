using MediatR;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Policy.Application.Features.Queries.Claims;

public class ClaimQueryHandlers : 
    IRequestHandler<GetClaimByIdQuery, Result<ClaimResponseDto>>,
    IRequestHandler<ListClaimsByCustomerQuery, Result<List<ClaimResponseDto>>>
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

    public async Task<Result<List<ClaimResponseDto>>> Handle(ListClaimsByCustomerQuery request, CancellationToken cancellationToken)
    {
        var claims = await _claimsRepository.ListByCustomerAsync(request.CustomerId, request.Page, request.PageSize, cancellationToken);
        return Result<List<ClaimResponseDto>>.Success(claims.Select(MapToDto).ToList());
    }

    private static ClaimResponseDto MapToDto(InsuranceEngine.Policy.Domain.Entities.Claim claim)
    {
        return new ClaimResponseDto
        {
            Id = claim.Id,
            ClaimNumber = claim.ClaimNumber,
            PolicyId = claim.PolicyId,
            CustomerId = claim.CustomerId,
            Status = claim.Status.ToString(),
            ClaimType = claim.Type.ToString(),
            ClaimedAmount = claim.ClaimedAmount / 100.0m,
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
