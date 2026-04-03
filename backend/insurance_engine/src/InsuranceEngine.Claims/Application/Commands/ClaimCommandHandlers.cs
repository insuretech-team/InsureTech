using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Domain.Events;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.Policy;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Claims.Application.Commands;

public sealed class SubmitClaimCommandHandler : IRequestHandler<SubmitClaimCommand, SubmitClaimResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IPolicyDataGateway _policyGateway;
    private readonly ILogger<SubmitClaimCommandHandler> _logger;
    private readonly IKafkaPublisher _kafkaPublisher;

    public SubmitClaimCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var policyResponse = await _policyGateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }
            
            var policy = policyResponse.Policy;

            if (policy.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Active)
            {
                return new SubmitClaimResponse
                {
                    Error = new Error { Code = "POLICY_NOT_ACTIVE", Message = $"Cannot submit claim for policy in status '{policy.Status}'" }
                };
            }

            var submitRequest = new SubmitClaimRequest
            {
                PolicyId = request.PolicyId,
                CustomerId = request.CustomerId,
                Type = MapToClaimType(request.ClaimType),
                ClaimedAmount = new Money { Amount = (long)(request.ClaimAmount * 100), Currency = "BDT" },
                IncidentDate = request.IncidentDate
            };
            
            if (request.DocumentUrls != null)
            {
                submitRequest.DocumentUrls.AddRange(request.DocumentUrls);
            }

            var response = await _claimsGateway.SubmitClaimAsync(submitRequest, cancellationToken);
            
            if (response.Error != null)
            {
                return response;
            }

            var evt = new ClaimSubmittedEvent(
                Guid.Parse(response.ClaimId), 
                response.ClaimNumber, 
                Guid.Parse(request.PolicyId), 
                Guid.Parse(request.CustomerId), 
                0,
                string.IsNullOrEmpty(policy.PartnerId) ? null : Guid.Parse(policy.PartnerId),
                string.IsNullOrEmpty(policy.AgentId) ? null : Guid.Parse(policy.AgentId)
            );
            await _kafkaPublisher.PublishAsync("insurance.claims.submitted", evt);

            _logger.LogInformation("Claim submitted via Go SSOT: {ClaimNumber} for Policy: {PolicyId}", 
                response.ClaimNumber, request.PolicyId);

            return response;
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

public sealed class ApproveClaimCommandHandler : IRequestHandler<ApproveClaimCommand, ApproveClaimResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<ApproveClaimCommandHandler> _logger;

    public ApproveClaimCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new ApproveClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;

            if (claim.Status != ClaimStatus.Submitted && claim.Status != ClaimStatus.UnderReview)
                return new ApproveClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be approved from status '{claim.Status}'" } };

            var approvalLevel = DetermineApprovalLevel(request.ApprovedAmount);
            var role = request.Role ?? "Unknown";

            var approvalResponse = await _claimsGateway.ApproveClaimAsync(request.ClaimId, request.Notes, cancellationToken);
            
            if (approvalResponse.Error != null)
                return new ApproveClaimResponse { Error = approvalResponse.Error };

            claim.Status = ClaimStatus.Approved;
            claim.ApprovedAmount = new Money { Amount = (long)(request.ApprovedAmount * 100), Currency = "BDT" };
            claim.ApprovedAt = Timestamp.FromDateTime(DateTime.UtcNow);

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

    private static int DetermineApprovalLevel(decimal amount) => amount switch
    {
        <= 10_000 => 1,
        <= 50_000 => 2,
        <= 200_000 => 3,
        _ => 4
    };
}

public sealed class RejectClaimCommandHandler : IRequestHandler<RejectClaimCommand, RejectClaimResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RejectClaimCommandHandler> _logger;

    public RejectClaimCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new RejectClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;

            if (claim.Status == ClaimStatus.Settled || claim.Status == ClaimStatus.Rejected)
                return new RejectClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim cannot be rejected from status '{claim.Status}'" } };

            var response = await _claimsGateway.RejectClaimAsync(request.ClaimId, request.Reason, cancellationToken);
            
            if (response.Error != null)
                return new RejectClaimResponse { Error = response.Error };

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

public sealed class SettleClaimCommandHandler : IRequestHandler<SettleClaimCommand, SettleClaimResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<SettleClaimCommandHandler> _logger;

    public SettleClaimCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new SettleClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;

            if (claim.Status != ClaimStatus.Approved)
                return new SettleClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = $"Claim must be APPROVED to settle, current: '{claim.Status}'" } };

            var response = await _claimsGateway.SettleClaimAsync(request.ClaimId, request.PaymentMethod, cancellationToken);
            
            if (response.Error != null)
                return new SettleClaimResponse { Error = response.Error };

            var paymentId = $"PAY-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            await _kafkaPublisher.PublishAsync("insurance.claims.settled", new
            {
                ClaimId = claim.ClaimId,
                SettledAmount = response.SettledAmount?.Amount,
                PaymentMethod = request.PaymentMethod,
                PaymentId = paymentId
            });

            _logger.LogInformation("Claim settled via Go SSOT: {ClaimNumber}, Amount: {Amount}", 
                claim.ClaimNumber, response.SettledAmount?.Amount);

            return new SettleClaimResponse
            {
                Message = "Claim settled successfully",
                SettledAmount = response.SettledAmount,
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

public sealed class UploadDocumentCommandHandler : IRequestHandler<UploadDocumentCommand, UploadDocumentResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly ILogger<UploadDocumentCommandHandler> _logger;

    public UploadDocumentCommandHandler(
        IClaimDataGateway claimsGateway,
        ILogger<UploadDocumentCommandHandler> logger)
    {
        _claimsGateway = claimsGateway;
        _logger = logger;
    }

    public async Task<UploadDocumentResponse> Handle(UploadDocumentCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new UploadDocumentResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;
            var documentUrl = $"https://storage.insuretech.labaid.com/claims/{claim.ClaimNumber}/{request.FileName}";

            var response = await _claimsGateway.UploadDocumentAsync(request.ClaimId, request.FileName, request.DocumentType, documentUrl, cancellationToken);
            
            if (response.Error != null)
                return new UploadDocumentResponse { Error = response.Error };

            _logger.LogInformation("Document uploaded via Go SSOT: {DocumentId} for Claim: {ClaimNumber}", response.DocumentId, claim.ClaimNumber);

            return new UploadDocumentResponse
            {
                DocumentId = response.DocumentId,
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

public sealed class RequestMoreDocumentsCommandHandler : IRequestHandler<RequestMoreDocumentsCommand, RequestMoreDocumentsResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<RequestMoreDocumentsCommandHandler> _logger;

    public RequestMoreDocumentsCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new RequestMoreDocumentsResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;
            var docTypes = string.Join(", ", request.RequiredDocumentTypes);
            var messageText = request.Message ?? $"Please upload the following documents: {docTypes}";

            var response = await _claimsGateway.RequestMoreDocumentsAsync(request.ClaimId, messageText, request.RequiredDocumentTypes, cancellationToken);
            
            if (response.Error != null)
                return new RequestMoreDocumentsResponse { Error = response.Error };

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

public sealed class DisputeClaimCommandHandler : IRequestHandler<DisputeClaimCommand, DisputeClaimResponse>
{
    private readonly IClaimDataGateway _claimsGateway;
    private readonly IKafkaPublisher _kafkaPublisher;
    private readonly ILogger<DisputeClaimCommandHandler> _logger;

    public DisputeClaimCommandHandler(
        IClaimDataGateway claimsGateway,
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
            var claimResponse = await _claimsGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claimResponse.Claim == null)
                return new DisputeClaimResponse { Error = new Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" } };

            var claim = claimResponse.Claim;

            if (claim.Status != ClaimStatus.Rejected)
                return new DisputeClaimResponse { Error = new Error { Code = "INVALID_STATUS", Message = "Only rejected claims can be disputed" } };

            if (!claim.AppealOptionAvailable)
                return new DisputeClaimResponse { Error = new Error { Code = "APPEAL_NOT_AVAILABLE", Message = "Appeal option is not available for this claim" } };

            if (claim.CustomerId != request.CustomerId)
                return new DisputeClaimResponse { Error = new Error { Code = "UNAUTHORIZED", Message = "You are not authorized to dispute this claim" } };

            var response = await _claimsGateway.DisputeClaimAsync(request.ClaimId, request.DisputeReason, request.CustomerId, cancellationToken);
            
            if (response.Error != null)
                return new DisputeClaimResponse { Error = response.Error };

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
