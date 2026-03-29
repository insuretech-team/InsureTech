using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;
using InsuranceEngine.Claims.Domain;
using Microsoft.EntityFrameworkCore;
using Google.Protobuf.WellKnownTypes;
using System.Security.Cryptography;

namespace InsuranceEngine.Claims.Application.Commands;

// ===== SubmitClaim =====
public sealed class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, SubmitClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly IRepository<ClaimDocumentEntity> _documentRepository;
    private readonly InsuranceDbContext _dbContext;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;
    private readonly IKafkaPublisher _kafkaPublisher;

    public SubmitClaimCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IRepository<PolicyEntity> policyRepository,
        IRepository<ClaimDocumentEntity> documentRepository,
        InsuranceDbContext dbContext,
        ILogger<SubmitClaimCommandHandler> logger,
        IKafkaPublisher kafkaPublisher)
    {
        _claimRepository = claimRepository;
        _policyRepository = policyRepository;
        _documentRepository = documentRepository;
        _dbContext = dbContext;
        _logger = logger;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<SubmitClaimResponse> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Validate policy exists and is ACTIVE
            var policy = await _policyRepository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }
            if (policy.Status != "ACTIVE")
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_ACTIVE", Message = $"Cannot submit claim for policy in status '{policy.Status}'" }
                };
            }

            // Get sequence number for claim number (FR-083)
            var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open)
                await connection.OpenAsync(cancellationToken);

            using var cmd = connection.CreateCommand();
            cmd.CommandText = "SELECT nextval('insurance_schema.claim_number_seq')";
            var seqResult = await cmd.ExecuteScalarAsync(cancellationToken);
            var sequenceNumber = Convert.ToInt64(seqResult);

            // Domain: Create Claim
            var claim = ClaimAggregate.Submit(
                policyId: Guid.Parse(request.PolicyId),
                type: request.ClaimType,
                amount: request.ClaimAmount,
                description: request.Description,
                sequenceNumber: sequenceNumber,
                documentContent: null
            );

            // Parse incident date
            DateTime incidentDate;
            if (!DateTime.TryParse(request.IncidentDate, out incidentDate))
                incidentDate = DateTime.UtcNow;

            var claimEntity = new ClaimEntity
            {
                ClaimId = claim.Id,
                ClaimNumber = claim.ClaimNumber,
                PolicyId = Guid.Parse(request.PolicyId),
                CustomerId = Guid.Parse(request.CustomerId),
                Status = "SUBMITTED",
                Type = request.ClaimType,
                ClaimedAmount = (long)(request.ClaimAmount * 100),
                ClaimedCurrency = "BDT",
                ApprovedCurrency = "BDT",
                SettledCurrency = "BDT",
                IncidentDate = incidentDate,
                IncidentDescription = request.Description,
                SubmittedAt = DateTime.UtcNow,
                ProcessingType = "MANUAL",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _claimRepository.AddAsync(claimEntity, cancellationToken);

            // Add initial document URLs (if any)
            if (request.DocumentUrls != null)
            {
                foreach (var url in request.DocumentUrls)
                {
                    await _documentRepository.AddAsync(new ClaimDocumentEntity
                    {
                        DocumentId = Guid.NewGuid(),
                        ClaimId = claimEntity.ClaimId,
                        DocumentType = "supporting_document",
                        FileUrl = url,
                        FileHash = ComputeHash(url),
                        Verified = false,
                        UploadedAt = DateTime.UtcNow,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }

            // Kafka event
            var evt = new ClaimSubmittedEvent(claim.Id, claim.ClaimNumber, claim.PolicyId, (long)(request.ClaimAmount * 100));
            await _kafkaPublisher.PublishAsync("insurance.claims.submitted", evt);

            _logger.LogInformation("Claim submitted: {ClaimNumber} for Policy: {PolicyId}", claim.ClaimNumber, request.PolicyId);

            return new SubmitClaimResponse
            {
                ClaimId = claim.Id.ToString(),
                ClaimNumber = claim.ClaimNumber,
                Message = "Claim submitted successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit claim for policy {PolicyId}", request.PolicyId);
            return new SubmitClaimResponse
            {
                Error = new Error { Code = "CLAIM_SUBMIT_FAILED", Message = ex.Message }
            };
        }
    }

    private static string ComputeHash(string input)
    {
        var bytes = SHA256.HashData(System.Text.Encoding.UTF8.GetBytes(input));
        return Convert.ToHexStringLower(bytes);
    }
}

// ===== ApproveClaim =====
public sealed class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, ApproveClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IRepository<ClaimApprovalEntity> _approvalRepository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IRepository<ClaimApprovalEntity> approvalRepository,
        IKafkaPublisher kafkaPublisher,
        ILogger<ApproveClaimCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _approvalRepository = approvalRepository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<ApproveClaimResponse> Handle(ApproveClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new ApproveClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != "SUBMITTED" && claim.Status != "UNDER_REVIEW")
                return new ApproveClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be approved from status '{claim.Status}'" } };

            // FR-086: Tiered Approval Matrix
            var approvalLevel = DetermineApprovalLevel(request.ApprovedAmount);

            claim.Status = "APPROVED";
            claim.ApprovedAmount = (long)(request.ApprovedAmount * 100);
            claim.ApprovedAt = DateTime.UtcNow;
            claim.UpdatedAt = DateTime.UtcNow;
            await _claimRepository.UpdateAsync(claim, cancellationToken);

            // Record approval
            await _approvalRepository.AddAsync(new ClaimApprovalEntity
            {
                ApprovalId = Guid.NewGuid(),
                ClaimId = claim.ClaimId,
                ApproverId = Guid.Parse(request.ApproverId),
                ApproverRole = "ClaimsOfficer",
                ApprovalLevel = approvalLevel,
                Decision = "APPROVED",
                ApprovedAmount = (long)(request.ApprovedAmount * 100),
                ApprovedCurrency = "BDT",
                Notes = request.Notes,
                ApprovedAt = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow
            }, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.approved", new { ClaimId = claim.ClaimId, ApprovedAmount = claim.ApprovedAmount });

            _logger.LogInformation("Claim approved: {ClaimNumber}, Amount: {ApprovedAmount}", claim.ClaimNumber, request.ApprovedAmount);

            return new ApproveClaimResponse { Message = "Claim approved successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve claim {ClaimId}", request.ClaimId);
            return new ApproveClaimResponse { Error = new Error { Code = "APPROVE_FAILED", Message = ex.Message } };
        }
    }

    private static int DetermineApprovalLevel(decimal amount) => amount switch
    {
        <= 50_000 => 1,   // L1: Claims Officer
        <= 500_000 => 2,  // L2: Claims Manager
        <= 2_000_000 => 3, // L3: Business Admin
        _ => 4             // Board level
    };
}

// ===== RejectClaim =====
public sealed class RejectClaimCommandHandler : IRequestHandler<RejectClaimCommand, RejectClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IRepository<ClaimApprovalEntity> _approvalRepository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RejectClaimCommandHandler> _logger;

    public RejectClaimCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IRepository<ClaimApprovalEntity> approvalRepository,
        IKafkaPublisher kafkaPublisher,
        ILogger<RejectClaimCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _approvalRepository = approvalRepository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<RejectClaimResponse> Handle(RejectClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new RejectClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status == "SETTLED" || claim.Status == "REJECTED")
                return new RejectClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be rejected from status '{claim.Status}'" } };

            claim.Status = "REJECTED";
            claim.RejectionReason = request.Reason;
            claim.AppealOptionAvailable = true; // Customer can dispute
            claim.UpdatedAt = DateTime.UtcNow;
            await _claimRepository.UpdateAsync(claim, cancellationToken);

            await _approvalRepository.AddAsync(new ClaimApprovalEntity
            {
                ApprovalId = Guid.NewGuid(),
                ClaimId = claim.ClaimId,
                ApproverId = Guid.Parse(request.ApproverId),
                ApproverRole = "ClaimsOfficer",
                ApprovalLevel = 1,
                Decision = "REJECTED",
                Notes = request.Reason,
                ApprovedAt = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow,
                ApprovedCurrency = "BDT"
            }, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.rejected", new { ClaimId = claim.ClaimId, Reason = request.Reason });

            _logger.LogInformation("Claim rejected: {ClaimNumber}, Reason: {Reason}", claim.ClaimNumber, request.Reason);

            return new RejectClaimResponse { Message = "Claim rejected" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject claim {ClaimId}", request.ClaimId);
            return new RejectClaimResponse { Error = new Error { Code = "REJECT_FAILED", Message = ex.Message } };
        }
    }
}

// ===== SettleClaim =====
public sealed class SettleClaimCommandHandler : IRequestHandler<SettleClaimCommand, SettleClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<SettleClaimCommandHandler> _logger;

    public SettleClaimCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IKafkaPublisher kafkaPublisher,
        ILogger<SettleClaimCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<SettleClaimResponse> Handle(SettleClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new SettleClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != "APPROVED")
                return new SettleClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim must be APPROVED to settle, current: '{claim.Status}'" } };

            claim.Status = "SETTLED";
            claim.SettledAmount = claim.ApprovedAmount; // Settle for approved amount
            claim.SettledAt = DateTime.UtcNow;
            claim.UpdatedAt = DateTime.UtcNow;
            await _claimRepository.UpdateAsync(claim, cancellationToken);

            var paymentId = $"PAY-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            await _kafkaPublisher.PublishAsync("insurance.claims.settled", new
            {
                ClaimId = claim.ClaimId,
                SettledAmount = claim.SettledAmount,
                PaymentMethod = request.PaymentMethod,
                PaymentId = paymentId
            });

            _logger.LogInformation("Claim settled: {ClaimNumber}, Amount: {Amount}", claim.ClaimNumber, claim.SettledAmount);

            return new SettleClaimResponse
            {
                Message = "Claim settled successfully",
                SettledAmount = new Money { Amount = claim.SettledAmount ?? 0, Currency = "BDT" },
                PaymentId = paymentId
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to settle claim {ClaimId}", request.ClaimId);
            return new SettleClaimResponse { Error = new Error { Code = "SETTLE_FAILED", Message = ex.Message } };
        }
    }
}

// ===== UploadDocument =====
public sealed class UploadDocumentCommandHandler : IRequestHandler<UploadDocumentCommand, UploadDocumentResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IRepository<ClaimDocumentEntity> _documentRepository;
    private readonly ILogger<UploadDocumentCommandHandler> _logger;

    public UploadDocumentCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IRepository<ClaimDocumentEntity> documentRepository,
        ILogger<UploadDocumentCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _documentRepository = documentRepository;
        _logger = logger;
    }

    public async Task<UploadDocumentResponse> Handle(UploadDocumentCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new UploadDocumentResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            // Compute SHA-256 hash for integrity verification
            var fileHash = Convert.ToHexStringLower(SHA256.HashData(request.FileData));

            // Simulate S3 upload — in production, upload to actual S3 bucket
            var documentUrl = $"https://storage.insuretech.labaid.com/claims/{claim.ClaimNumber}/{request.FileName}";

            var document = new ClaimDocumentEntity
            {
                DocumentId = Guid.NewGuid(),
                ClaimId = claim.ClaimId,
                DocumentType = request.DocumentType,
                FileUrl = documentUrl,
                FileHash = fileHash,
                Verified = false,
                UploadedAt = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _documentRepository.AddAsync(document, cancellationToken);

            // If claim was pending documents, transition back to under_review
            if (claim.Status == "PENDING_DOCUMENTS")
            {
                claim.Status = "UNDER_REVIEW";
                claim.UpdatedAt = DateTime.UtcNow;
                await _claimRepository.UpdateAsync(claim, cancellationToken);
            }

            _logger.LogInformation("Document uploaded: {DocumentId} for Claim: {ClaimNumber}", document.DocumentId, claim.ClaimNumber);

            return new UploadDocumentResponse
            {
                DocumentId = document.DocumentId.ToString(),
                DocumentUrl = documentUrl,
                FileHash = fileHash
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to upload document for claim {ClaimId}", request.ClaimId);
            return new UploadDocumentResponse { Error = new Error { Code = "UPLOAD_FAILED", Message = ex.Message } };
        }
    }
}

// ===== RequestMoreDocuments =====
public sealed class RequestMoreDocumentsCommandHandler : IRequestHandler<RequestMoreDocumentsCommand, RequestMoreDocumentsResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RequestMoreDocumentsCommandHandler> _logger;

    public RequestMoreDocumentsCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IKafkaPublisher kafkaPublisher,
        ILogger<RequestMoreDocumentsCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<RequestMoreDocumentsResponse> Handle(RequestMoreDocumentsCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new RequestMoreDocumentsResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            claim.Status = "PENDING_DOCUMENTS";
            claim.UpdatedAt = DateTime.UtcNow;

            // Store required document types as in-app message
            var docTypes = string.Join(", ", request.RequiredDocumentTypes);
            var messageText = request.Message ?? $"Please upload the following documents: {docTypes}";
            claim.InAppMessages = System.Text.Json.JsonSerializer.Serialize(new[]
            {
                new { timestamp = DateTime.UtcNow.ToString("o"), type = "document_request", message = messageText, requiredTypes = request.RequiredDocumentTypes }
            });

            await _claimRepository.UpdateAsync(claim, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.documents_requested", new
            {
                ClaimId = claim.ClaimId,
                RequiredTypes = request.RequiredDocumentTypes,
                Message = messageText
            });

            _logger.LogInformation("Document request sent for Claim: {ClaimNumber}", claim.ClaimNumber);

            return new RequestMoreDocumentsResponse { Message = messageText };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to request documents for claim {ClaimId}", request.ClaimId);
            return new RequestMoreDocumentsResponse { Error = new Error { Code = "REQUEST_FAILED", Message = ex.Message } };
        }
    }
}

// ===== DisputeClaim =====
public sealed class DisputeClaimCommandHandler : IRequestHandler<DisputeClaimCommand, DisputeClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<DisputeClaimCommandHandler> _logger;

    public DisputeClaimCommandHandler(
        IRepository<ClaimEntity> claimRepository,
        IKafkaPublisher kafkaPublisher,
        ILogger<DisputeClaimCommandHandler> logger)
    {
        _claimRepository = claimRepository;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<DisputeClaimResponse> Handle(DisputeClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (claim == null)
                return new DisputeClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != "REJECTED")
                return new DisputeClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = "Only rejected claims can be disputed" } };

            if (!claim.AppealOptionAvailable)
                return new DisputeClaimResponse { Error = new Error { Code = "APPEAL_NOT_AVAILABLE", Message = "Appeal option is not available for this claim" } };

            // Verify customer owns the claim
            if (claim.CustomerId != Guid.Parse(request.CustomerId))
                return new DisputeClaimResponse { Error = new Error { Code = "UNAUTHORIZED", Message = "You are not authorized to dispute this claim" } };

            claim.Status = "DISPUTED";
            claim.AppealOptionAvailable = false; // One-time appeal
            claim.UpdatedAt = DateTime.UtcNow;
            await _claimRepository.UpdateAsync(claim, cancellationToken);

            var disputeId = $"DSP-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            await _kafkaPublisher.PublishAsync("insurance.claims.disputed", new
            {
                ClaimId = claim.ClaimId,
                DisputeId = disputeId,
                Reason = request.DisputeReason,
                SupportingDocs = request.SupportingDocumentUrls
            });

            _logger.LogInformation("Claim disputed: {ClaimNumber}, DisputeId: {DisputeId}", claim.ClaimNumber, disputeId);

            return new DisputeClaimResponse
            {
                DisputeId = disputeId,
                Message = "Dispute submitted successfully. Your claim will be reviewed again."
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to dispute claim {ClaimId}", request.ClaimId);
            return new DisputeClaimResponse { Error = new Error { Code = "DISPUTE_FAILED", Message = ex.Message } };
        }
    }
}
