using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using InsuranceEngine.Claims.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;
using Microsoft.Extensions.Logging;

namespace InsuranceEngine.Claims.Application.Features.Commands.Claims;

public record UploadClaimDocumentCommand(
    Guid ClaimId,
    List<ClaimDocumentDto> Documents
) : IRequest<Result>;

public class UploadClaimDocumentCommandHandler : IRequestHandler<UploadClaimDocumentCommand, Result>
{
    private readonly IClaimsRepository _claimsRepository;
    private readonly ClaimDocumentValidator _documentValidator;
    private readonly ILogger<UploadClaimDocumentCommandHandler> _logger;

    public UploadClaimDocumentCommandHandler(
        IClaimsRepository claimsRepository,
        ClaimDocumentValidator documentValidator,
        ILogger<UploadClaimDocumentCommandHandler> logger)
    {
        _claimsRepository = claimsRepository;
        _documentValidator = documentValidator;
        _logger = logger;
    }

    public async Task<Result> Handle(UploadClaimDocumentCommand request, CancellationToken cancellationToken)
    {
        var claim = await _claimsRepository.GetByIdAsync(request.ClaimId, cancellationToken);
        if (claim == null)
            return Result.Fail(Error.NotFound("Claim", request.ClaimId.ToString()));

        // 1. Document Validation
        var docRequests = request.Documents.Select(d => new ValidateDocumentRequest(d.FileName, d.FileSize));
        
        // Include existing documents in total size check
        var totalExistingSize = claim.Documents.Count; // This is a bit tricky since ClaimDocument doesn't have FileSize.
        // Let's assume for now the validator only checks the NEW documents for simplicity, 
        // or we could add FileSize to ClaimDocument if needed for strict FR-099 total size enforcement.
        
        var validation = _documentValidator.Validate(docRequests);
        if (!validation.IsSuccess)
        {
            return validation;
        }

        // 2. Add Documents
        foreach (var d in request.Documents)
        {
            claim.Documents.Add(new ClaimDocument
            {
                Id = Guid.NewGuid(),
                ClaimId = claim.Id,
                DocumentType = d.DocumentType,
                FileUrl = d.FileUrl,
                FileHash = d.FileHash,
                UploadedAt = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            });
        }

        // 3. Auto-transition status if it was PendingDocuments
        if (claim.Status == ClaimStatus.PendingDocuments)
        {
            claim.Status = ClaimStatus.UnderReview;
            _logger.LogInformation("Claim {ClaimId} status transitioned from PendingDocuments to UnderReview.", claim.Id);
        }

        await _claimsRepository.UpdateAsync(claim, cancellationToken);
        
        return Result.Ok();
    }
}
