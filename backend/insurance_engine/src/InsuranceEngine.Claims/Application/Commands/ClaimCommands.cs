using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Claims.Services.V1;
using MediatR;

namespace InsuranceEngine.Claims.Application.Commands;

public sealed record SubmitClaimCommand(
    string PolicyId,
    string ClaimType,
    decimal ClaimAmount,
    string Description,
    string? BeneficiaryId,
    string? DocumentContent = null) : IRequest<SubmitClaimResponse>;

public sealed record ApproveClaimCommand(string ClaimId, decimal ApprovedAmount) : IRequest<ApproveClaimResponse>;
public sealed record RejectClaimCommand(string ClaimId, string Reason) : IRequest<bool>;
public sealed record SettleClaimCommand(string ClaimId, decimal SettlementAmount) : IRequest<bool>;
