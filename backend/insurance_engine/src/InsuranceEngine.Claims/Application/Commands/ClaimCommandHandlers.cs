using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;
using InsuranceEngine.Grpc.Gateways;
using InsuranceEngine.Grpc.Gateways;
using InsuranceEngine.SharedKernel.Domain.Events;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Claims.Application.Commands;

// ===== SubmitClaim =====
public sealed class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, SubmitClaimResponse>
{
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IPolicyDataGateway _policyGateway;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;
    private readonly IKafkaPublisher _kafkaPublisher;

    public SubmitClaimCommandHandler(
        IClaimsDataGateway claimsGateway,
        IPolicyDataGateway policyGateway,
        ILogger<SubmitClaimCommandHandler> logger,
        IKafkaPublisher kafkaPublisher)
    {
        _claimsGateway = claimsGateway;
        _policyGateway = policyGateway;
        _logger = logger;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<SubmitClaimResponse> Handle(SubmitClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Validate policy exists and is ACTIVE via Go SSOT
            var policy = await _policyGateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policy == null)
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }
            if (policy.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Active)
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_ACTIVE", Message = $"Cannot submit claim for policy in status '{policy.Status}'" }
                };
            }

            // Note: Sequence numbers, ZHTC Auto-Approval, and status logic 
            // are now handled by the Go backend (SSOT).

            var claimReq = new Insuretech.Claims.Entity.V1.Claim
            {
                PolicyId = request.PolicyId,
                CustomerId = request.CustomerId,
                Type = MapToClaimType(request.ClaimType),
                ClaimedAmount = new Money { Amount = (long)(request.ClaimAmount * 100), Currency = "BDT" },
                IncidentDescription = request.Description,
                IncidentDate = Timestamp.FromDateTime(DateTime.TryParse(request.IncidentDate, out var dt) ? dt.ToUniversalTime() : DateTime.UtcNow),
                PlaceOfIncident = request.PlaceOfIncident ?? string.Empty
            };

            var createdClaim = await _claimsGateway.CreateClaimAsync(claimReq, cancellationToken);
            
            // Add initial documents via Go Gateway
            if (request.DocumentUrls != null)
            {
                foreach (var url in request.DocumentUrls)
                {
                    await _claimsGateway.CreateClaimDocumentAsync(new ClaimDocument
                    {
                        ClaimId = createdClaim.ClaimId,
                        DocumentType = "supporting_document",
                        FileUrl = url,
                        Verified = false
                    }, cancellationToken);
                }
            }

            // Kafka event: We still publish from C# as specified in architecture, 
            // but Go backend should ideally be the source.
            var evt = new ClaimSubmittedEvent(
                Guid.Parse(createdClaim.ClaimId), 
                createdClaim.ClaimNumber, 
                Guid.Parse(createdClaim.PolicyId), 
                Guid.Parse(createdClaim.CustomerId), 
                createdClaim.ClaimedAmount.Amount,
                policy.PartnerId,
                policy.AgentId
            );
            await _kafkaPublisher.PublishAsync("insurance.claims.submitted", evt);

            _logger.LogInformation("Claim submitted via Go SSOT: {ClaimNumber} for Policy: {PolicyId}", 
                createdClaim.ClaimNumber, request.PolicyId);

            return new SubmitClaimResponse
            {
                ClaimId = createdClaim.ClaimId,
                ClaimNumber = createdClaim.ClaimNumber,
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

    private static ClaimType MapToClaimType(string type) => type switch
    {
        "HEALTH" => ClaimType.HealthHospitalization,
        "MOTOR" => ClaimType.MotorAccident,
        "TRAVEL" => ClaimType.TravelMedical,
        "DEVICE" => ClaimType.DeviceDamage,
        "DEATH" => ClaimType.Death,
        _ => ClaimType.Unspecified
    };
}

// ===== ApproveClaim =====
public sealed class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, ApproveClaimResponse>
{
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(
        IClaimsDataGateway claimsGateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<ApproveClaimCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<ApproveClaimResponse> Handle(ApproveClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new ApproveClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != ClaimStatus.Submitted && claim.Status != ClaimStatus.UnderReview)
                return new ApproveClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be approved from status '{claim.Status}'" } };

            var approvalLevel = DetermineApprovalLevel(request.ApprovedAmount);
            var role = request.Role ?? "Unknown";

            // Record Approval via Go Gateway
            await _claimsGateway.CreateClaimApprovalAsync(new ClaimApproval
            {
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = role,
                ApprovalLevel = approvalLevel,
                Decision = ApprovalDecision.Approved,
                ApprovedAmount = new Money { Amount = (long)(request.ApprovedAmount * 100), Currency = "BDT" },
                Notes = request.Notes
            }, cancellationToken);

            // Level 3 Joint Approval Logic check
            if (approvalLevel == 3)
            {
                // In a pure SSOT, Go backend should decide if 'fully approved'
                // Here we simulate the logic: if both signatures exist, set to Approved.
                var totalSignatures = claim.Approvals.Count(a => a.ApprovalLevel == 3) + 1;
                if (totalSignatures < 2)
                {
                    claim.Status = ClaimStatus.UnderReview;
                    await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);
                    return new ApproveClaimResponse { Message = "Approval recorded. Waiting for remaining signature." };
                }
            }

            // Mark as Approved in Go SSOT
            claim.Status = ClaimStatus.Approved;
            claim.ApprovedAmount = new Money { Amount = (long)(request.ApprovedAmount * 100), Currency = "BDT" };
            claim.ApprovedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            
            await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.approved", new { ClaimId = claim.ClaimId, ApprovedAmount = claim.ApprovedAmount.Amount, Level = approvalLevel });

            _logger.LogInformation("Claim approved via Go SSOT: {ClaimNumber} by {Role}", claim.ClaimNumber, role);

            return new ApproveClaimResponse { Message = "Claim approved successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve claim {ClaimId}", request.ClaimId);
            return new ApproveClaimResponse { Error = new Error { Code = "APPROVE_FAILED", Message = ex.Message } };
        }
    }

    // FR-086/542-545: Tiered Approval Matrix:
    // BDT 0–10K: Officer (Level 1)
    // BDT 10K–50K: Manager (Level 2)
    // BDT 50K–2L: Joint (BA + FP) (Level 3)
    // BDT 2L+: Board (Level 4)
    private static int DetermineApprovalLevel(decimal amount) => amount switch
    {
        <= 10_000 => 1,
        <= 50_000 => 2,
        <= 200_000 => 3,
        _ => 4
    };
}

// ===== RejectClaim =====
public sealed class RejectClaimCommandHandler : IRequestHandler<RejectClaimCommand, RejectClaimResponse>
{
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RejectClaimCommandHandler> _logger;

    public RejectClaimCommandHandler(
        IClaimsDataGateway claimsGateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<RejectClaimCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<RejectClaimResponse> Handle(RejectClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new RejectClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status == ClaimStatus.Settled || claim.Status == ClaimStatus.Rejected)
                return new RejectClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be rejected from status '{claim.Status}'" } };

            // Reject in Go SSOT
            claim.Status = ClaimStatus.Rejected;
            claim.RejectionReason = request.Reason;
            claim.AppealOptionAvailable = true;
            
            await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);

            // Record rejection via Go Gateway
            await _claimsGateway.CreateClaimApprovalAsync(new ClaimApproval
            {
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = "ClaimsOfficer",
                ApprovalLevel = 1,
                Decision = ApprovalDecision.Rejected,
                Notes = request.Reason
            }, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.rejected", new { ClaimId = claim.ClaimId, Reason = request.Reason });

            _logger.LogInformation("Claim rejected via Go SSOT: {ClaimNumber}, Reason: {Reason}", claim.ClaimNumber, request.Reason);

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
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<SettleClaimCommandHandler> _logger;

    public SettleClaimCommandHandler(
        IClaimsDataGateway claimsGateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<SettleClaimCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<SettleClaimResponse> Handle(SettleClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new SettleClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != ClaimStatus.Approved)
                return new SettleClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim must be APPROVED to settle, current: '{claim.Status}'" } };

            // Settle in Go SSOT
            claim.Status = ClaimStatus.Settled;
            claim.SettledAmount = claim.ApprovedAmount;
            claim.SettledAt = Timestamp.FromDateTime(DateTime.UtcNow);
            
            await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);

            var paymentId = $"PAY-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            await _kafkaPublisher.PublishAsync("insurance.claims.settled", new
            {
                ClaimId = claim.ClaimId,
                SettledAmount = claim.SettledAmount.Amount,
                PaymentMethod = request.PaymentMethod,
                PaymentId = paymentId
            });

            _logger.LogInformation("Claim settled via Go SSOT: {ClaimNumber}, Amount: {Amount}", 
                claim.ClaimNumber, claim.SettledAmount.Amount);

            return new SettleClaimResponse
            {
                Message = "Claim settled successfully",
                SettledAmount = claim.SettledAmount,
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
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly ILogger<UploadDocumentCommandHandler> _logger;

    public UploadDocumentCommandHandler(
        IClaimsDataGateway claimsGateway,
        ILogger<UploadDocumentCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _logger = logger;
    }

    public async Task<UploadDocumentResponse> Handle(UploadDocumentCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new UploadDocumentResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            // Simulate S3 upload
            var documentUrl = $"https://storage.insuretech.labaid.com/claims/{claim.ClaimNumber}/{request.FileName}";

            var document = new ClaimDocument
            {
                ClaimId = claim.ClaimId,
                DocumentType = request.DocumentType,
                FileUrl = documentUrl,
                Verified = false
            };

            var createdDoc = await _claimsGateway.CreateClaimDocumentAsync(document, cancellationToken);

            // If claim was pending documents, transition back to under_review in Go SSOT
            if (claim.Status == ClaimStatus.PendingDocuments)
            {
                claim.Status = ClaimStatus.UnderReview;
                await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);
            }

            _logger.LogInformation("Document uploaded via Go SSOT: {DocumentId} for Claim: {ClaimNumber}", createdDoc.DocumentId, claim.ClaimNumber);

            return new UploadDocumentResponse
            {
                DocumentId = createdDoc.DocumentId,
                DocumentUrl = documentUrl
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
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RequestMoreDocumentsCommandHandler> _logger;

    public RequestMoreDocumentsCommandHandler(
        IClaimsDataGateway claimsGateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<RequestMoreDocumentsCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<RequestMoreDocumentsResponse> Handle(RequestMoreDocumentsCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new RequestMoreDocumentsResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            claim.Status = ClaimStatus.PendingDocuments;

            // Store required document types as in-app message via Go SSOT
            var docTypes = string.Join(", ", request.RequiredDocumentTypes);
            var messageText = request.Message ?? $"Please upload the following documents: {docTypes}";
            claim.InAppMessages = System.Text.Json.JsonSerializer.Serialize(new[]
            {
                new { timestamp = DateTime.UtcNow.ToString("o"), type = "document_request", message = messageText, requiredTypes = request.RequiredDocumentTypes }
            });

            await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);

            await _kafkaPublisher.PublishAsync("insurance.claims.documents_requested", new
            {
                ClaimId = claim.ClaimId,
                RequiredTypes = request.RequiredDocumentTypes,
                Message = messageText
            });

            _logger.LogInformation("Document request sent via Go SSOT for Claim: {ClaimNumber}", claim.ClaimNumber);

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
    private readonly IClaimsDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<DisputeClaimCommandHandler> _logger;

    public DisputeClaimCommandHandler(
        IClaimsDataGateway claimsGateway,
        IKafkaPublisher kafkaPublisher,
        ILogger<DisputeClaimCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _kafkaPublisher = kafkaPublisher;
        _logger = logger;
    }

    public async Task<DisputeClaimResponse> Handle(DisputeClaimCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim == null)
                return new DisputeClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            if (claim.Status != ClaimStatus.Rejected)
                return new DisputeClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = "Only rejected claims can be disputed" } };

            if (!claim.AppealOptionAvailable)
                return new DisputeClaimResponse { Error = new Error { Code = "APPEAL_NOT_AVAILABLE", Message = "Appeal option is not available for this claim" } };

            // Verify customer owns the claim via Go SSOT data
            if (claim.CustomerId != request.CustomerId)
                return new DisputeClaimResponse { Error = new Error { Code = "UNAUTHORIZED", Message = "You are not authorized to dispute this claim" } };

            claim.Status = ClaimStatus.Disputed;
            claim.AppealOptionAvailable = false; 

            await _claimsGateway.UpdateClaimAsync(claim, cancellationToken);

            var disputeId = $"DSP-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            await _kafkaPublisher.PublishAsync("insurance.claims.disputed", new
            {
                ClaimId = claim.ClaimId,
                DisputeId = disputeId,
                Reason = request.DisputeReason,
                SupportingDocs = request.SupportingDocumentUrls
            });

            _logger.LogInformation("Claim disputed via Go SSOT: {ClaimNumber}, DisputeId: {DisputeId}", claim.ClaimNumber, disputeId);

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
