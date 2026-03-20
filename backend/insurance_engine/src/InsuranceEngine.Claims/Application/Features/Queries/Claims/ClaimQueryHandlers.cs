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
            ProcessingType = claim.ProcessingType,
            ClaimedAmount = new MoneyDto(claim.ClaimedAmount, claim.ClaimedCurrency),
            ApprovedAmount = new MoneyDto(claim.ApprovedAmount, claim.ApprovedCurrency),
            SettledAmount = new MoneyDto(claim.SettledAmount, claim.SettledCurrency),
            DeductibleAmount = new MoneyDto(claim.DeductibleAmount, claim.DeductibleCurrency),
            CoPayAmount = new MoneyDto(claim.CoPayAmount, claim.CoPayCurrency),
            IncidentDate = claim.IncidentDate,
            IncidentDescription = claim.IncidentDescription,
            PlaceOfIncident = claim.PlaceOfIncident,
            SubmittedAt = claim.SubmittedAt,
            ApprovedAt = claim.ApprovedAt,
            SettledAt = claim.SettledAt,
            RejectionReason = claim.RejectionReason,
            AppealOptionAvailable = claim.AppealOptionAvailable,
            FraudCheck = claim.FraudCheck != null ? new FraudCheckResultDto
            {
                Id = claim.FraudCheck.Id,
                FraudScore = claim.FraudCheck.FraudScore,
                RiskFactors = claim.FraudCheck.RiskFactors,
                Flagged = claim.FraudCheck.Flagged,
                ReviewedBy = claim.FraudCheck.ReviewedBy,
                ReviewedAt = claim.FraudCheck.ReviewedAt
            } : null,
            Approvals = claim.Approvals.Select(a => new ClaimApprovalDto
            {
                Id = a.Id,
                ApproverId = a.ApproverId,
                ApproverRole = a.ApproverRole,
                ApprovalLevel = a.ApprovalLevel,
                Decision = a.Decision,
                ApprovedAmount = new MoneyDto(a.ApprovedAmount, a.ApprovedCurrency),
                Notes = a.Notes,
                ApprovedAt = a.ApprovedAt
            }).ToList(),
            Documents = claim.Documents.Select(d => new ClaimDocumentDto
            {
                Id = d.Id,
                DocumentType = d.DocumentType,
                FileUrl = d.FileUrl,
                FileHash = d.FileHash,
                Verified = d.Verified,
                UploadedAt = d.UploadedAt
            }).ToList(),
            CreatedAt = claim.CreatedAt,
            UpdatedAt = claim.UpdatedAt
        };
    }
}

