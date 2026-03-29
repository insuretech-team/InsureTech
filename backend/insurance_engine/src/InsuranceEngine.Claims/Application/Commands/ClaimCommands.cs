using Insuretech.Claims.Services.V1;
using MediatR;

namespace InsuranceEngine.Claims.Application.Commands;

public sealed record SubmitClaimCommand(
    string PolicyId,
    string CustomerId,
    string ClaimType,
    decimal ClaimAmount,
    string IncidentDate,
    string Description,
    List<string>? DocumentUrls = null) : IRequest<SubmitClaimResponse>;

public sealed record ApproveClaimCommand(
    string ClaimId, 
    string ApproverId, 
    decimal ApprovedAmount,
    string? Notes = null) : IRequest<ApproveClaimResponse>;

public sealed record RejectClaimCommand(
    string ClaimId, 
    string ApproverId, 
    string Reason) : IRequest<RejectClaimResponse>;

public sealed record SettleClaimCommand(
    string ClaimId, 
    string PaymentMethod,
    string? PaymentReference = null) : IRequest<SettleClaimResponse>;

public sealed record UploadDocumentCommand(
    string ClaimId,
    string DocumentType,
    byte[] FileData,
    string FileName,
    string? MimeType = null) : IRequest<UploadDocumentResponse>;

public sealed record RequestMoreDocumentsCommand(
    string ClaimId,
    List<string> RequiredDocumentTypes,
    string? Message = null) : IRequest<RequestMoreDocumentsResponse>;

public sealed record DisputeClaimCommand(
    string ClaimId,
    string CustomerId,
    string DisputeReason,
    List<string>? SupportingDocumentUrls = null) : IRequest<DisputeClaimResponse>;
