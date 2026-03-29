using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using InsuranceEngine.Claims.Application.Queries;
using InsuranceEngine.Claims.Application.Commands;

namespace InsuranceEngine.Claims.GrpcServices;

public sealed class ClaimGrpcService : ClaimService.ClaimServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ClaimGrpcService> _logger;

    public ClaimGrpcService(IMediator mediator, ILogger<ClaimGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    // RPC 1: SubmitClaim
    public override async Task<SubmitClaimResponse> SubmitClaim(
        SubmitClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.PolicyId) || string.IsNullOrEmpty(request.CustomerId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Policy ID and Customer ID are required"));
        }

        var command = new SubmitClaimCommand(
            PolicyId: request.PolicyId,
            CustomerId: request.CustomerId,
            ClaimType: request.Type.ToString(),
            ClaimAmount: request.ClaimedAmount != null ? (decimal)request.ClaimedAmount.Amount / 100m : 0m,
            IncidentDate: request.IncidentDate,
            Description: request.IncidentDescription,
            DocumentUrls: request.DocumentUrls.Count > 0 ? request.DocumentUrls.ToList() : null
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 2: GetClaim
    public override async Task<GetClaimResponse> GetClaim(
        GetClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID is required"));
        }

        return await _mediator.Send(new GetClaimQuery(request.ClaimId), context.CancellationToken);
    }

    // RPC 3: ListUserClaims
    public override async Task<ListUserClaimsResponse> ListUserClaims(
        ListUserClaimsRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.CustomerId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Customer ID is required"));
        }

        var query = new ListUserClaimsQuery(
            CustomerId: request.CustomerId,
            PolicyId: null,
            Status: request.Status != Insuretech.Claims.Entity.V1.ClaimStatus.Unspecified ? request.Status.ToString() : null,
            Page: request.Page <= 0 ? 1 : request.Page,
            PageSize: request.PageSize <= 0 ? 10 : request.PageSize
        );

        return await _mediator.Send(query, context.CancellationToken);
    }

    // RPC 4: UploadDocument
    public override async Task<UploadDocumentResponse> UploadDocument(
        UploadDocumentRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID is required"));
        }
        if (string.IsNullOrEmpty(request.DocumentType) || string.IsNullOrEmpty(request.FileName))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Document type and file name are required"));
        }

        var command = new UploadDocumentCommand(
            ClaimId: request.ClaimId,
            DocumentType: request.DocumentType,
            FileData: request.FileData.ToByteArray(),
            FileName: request.FileName,
            MimeType: request.MimeType
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 5: ApproveClaim
    public override async Task<ApproveClaimResponse> ApproveClaim(
        ApproveClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId) || string.IsNullOrEmpty(request.ApproverId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID and Approver ID are required"));
        }

        var command = new ApproveClaimCommand(
            ClaimId: request.ClaimId,
            ApproverId: request.ApproverId,
            ApprovedAmount: request.ApprovedAmount != null ? (decimal)request.ApprovedAmount.Amount / 100m : 0m,
            Notes: request.Notes
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 6: RejectClaim
    public override async Task<RejectClaimResponse> RejectClaim(
        RejectClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId) || string.IsNullOrEmpty(request.ApproverId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID and Approver ID are required"));
        }
        if (string.IsNullOrEmpty(request.Reason))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Rejection reason is required"));
        }

        return await _mediator.Send(new RejectClaimCommand(request.ClaimId, request.ApproverId, request.Reason), context.CancellationToken);
    }

    // RPC 7: SettleClaim
    public override async Task<SettleClaimResponse> SettleClaim(
        SettleClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID is required"));
        }
        if (string.IsNullOrEmpty(request.PaymentMethod))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Payment method is required"));
        }

        return await _mediator.Send(new SettleClaimCommand(request.ClaimId, request.PaymentMethod, request.PaymentReference), context.CancellationToken);
    }

    // RPC 8: RequestMoreDocuments
    public override async Task<RequestMoreDocumentsResponse> RequestMoreDocuments(
        RequestMoreDocumentsRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID is required"));
        }
        if (request.RequiredDocumentTypes.Count == 0)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "At least one document type is required"));
        }

        var command = new RequestMoreDocumentsCommand(
            ClaimId: request.ClaimId,
            RequiredDocumentTypes: request.RequiredDocumentTypes.ToList(),
            Message: request.Message
        );

        return await _mediator.Send(command, context.CancellationToken);
    }

    // RPC 9: DisputeClaim
    public override async Task<DisputeClaimResponse> DisputeClaim(
        DisputeClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.ClaimId) || string.IsNullOrEmpty(request.CustomerId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Claim ID and Customer ID are required"));
        }
        if (string.IsNullOrEmpty(request.DisputeReason))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Dispute reason is required"));
        }

        var command = new DisputeClaimCommand(
            ClaimId: request.ClaimId,
            CustomerId: request.CustomerId,
            DisputeReason: request.DisputeReason,
            SupportingDocumentUrls: request.SupportingDocumentUrls.Count > 0 ? request.SupportingDocumentUrls.ToList() : null
        );

        return await _mediator.Send(command, context.CancellationToken);
    }
}
