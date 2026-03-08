using System.Security.Cryptography;
using System.Text;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;
using Insuretech.Common.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;

namespace PoliSync.Claims.GrpcServices;

public sealed class ClaimGrpcService : ClaimService.ClaimServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<ClaimGrpcService> _logger;
    private readonly IClaimDataGateway _dataGateway;

    public ClaimGrpcService(
        IMediator mediator,
        ILogger<ClaimGrpcService> logger,
        IClaimDataGateway dataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _dataGateway = dataGateway;
    }

    public override async Task<SubmitClaimResponse> SubmitClaim(SubmitClaimRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.PolicyId) || string.IsNullOrWhiteSpace(request.CustomerId))
        {
            return new SubmitClaimResponse
            {
                Error = BuildError("VALIDATION_ERROR", "PolicyId and CustomerId are required")
            };
        }

        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        var incidentDate = ParseDateOrNow(request.IncidentDate);
        var claimId = Guid.NewGuid().ToString("N");
        var claim = new Claim
        {
            ClaimId = claimId,
            ClaimNumber = $"CLM-{DateTime.UtcNow:yyyy}-{Random.Shared.Next(100000, 999999)}",
            PolicyId = request.PolicyId,
            CustomerId = request.CustomerId,
            Type = request.Type,
            Status = ClaimStatus.Submitted,
            ClaimedAmount = request.ClaimedAmount ?? NewMoney(0),
            ApprovedAmount = NewMoney(0),
            SettledAmount = NewMoney(0),
            IncidentDate = Timestamp.FromDateTime(incidentDate),
            IncidentDescription = request.IncidentDescription,
            SubmittedAt = now,
            CreatedAt = now,
            UpdatedAt = now,
            ProcessingType = ClaimProcessingType.Manual,
            ClaimedCurrency = request.ClaimedAmount?.Currency ?? "BDT",
            ApprovedCurrency = "BDT",
            SettledCurrency = "BDT"
        };

        foreach (var url in request.DocumentUrls)
        {
            claim.Documents.Add(new ClaimDocument
            {
                DocumentId = Guid.NewGuid().ToString("N"),
                ClaimId = claimId,
                DocumentType = "SUPPORTING",
                FileUrl = url,
                FileHash = ComputeSha256(url),
                UploadedAt = now,
                CreatedAt = now,
                UpdatedAt = now
            });
        }

        try
        {
            var created = await _dataGateway.CreateClaimAsync(claim, GetCancellationToken(context));
            _logger.LogInformation("Claim submitted: {ClaimId}", created.ClaimId);

            return new SubmitClaimResponse
            {
                ClaimId = created.ClaimId,
                ClaimNumber = created.ClaimNumber,
                Message = "Claim submitted successfully"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to submit claim for policy {PolicyId}", request.PolicyId);
            return new SubmitClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetClaimResponse> GetClaim(GetClaimRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new GetClaimResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            return new GetClaimResponse { Claim = claim };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get claim {ClaimId}", request.ClaimId);
            return new GetClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ListUserClaimsResponse> ListUserClaims(ListUserClaimsRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;

        try
        {
            var claims = await _dataGateway.ListClaimsAsync(request.CustomerId, string.Empty, page, pageSize, GetCancellationToken(context));

            var filtered = request.Status == ClaimStatus.Unspecified
                ? claims
                : claims.Where(c => c.Status == request.Status).ToList();

            var ordered = filtered.OrderByDescending(c => c.CreatedAt?.Seconds ?? 0).ToList();
            var response = new ListUserClaimsResponse { TotalCount = ordered.Count };
            response.Claims.AddRange(ordered);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list claims for customer {CustomerId}", request.CustomerId);
            return new ListUserClaimsResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<UploadDocumentResponse> UploadDocument(UploadDocumentRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new UploadDocumentResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            var documentId = Guid.NewGuid().ToString("N");
            var contentHash = request.FileData.Length > 0
                ? ComputeSha256(request.FileData.ToByteArray())
                : ComputeSha256(request.FileName + request.MimeType + request.DocumentType);

            var documentUrl = $"memory://claims/{claim.ClaimId}/documents/{documentId}/{request.FileName}";
            claim.Documents.Add(new ClaimDocument
            {
                DocumentId = documentId,
                ClaimId = claim.ClaimId,
                DocumentType = request.DocumentType,
                FileUrl = documentUrl,
                FileHash = contentHash,
                UploadedAt = now,
                CreatedAt = now,
                UpdatedAt = now
            });
            claim.UpdatedAt = now;

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));

            return new UploadDocumentResponse
            {
                DocumentId = documentId,
                DocumentUrl = documentUrl,
                FileHash = contentHash
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to upload document for claim {ClaimId}", request.ClaimId);
            return new UploadDocumentResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ApproveClaimResponse> ApproveClaim(ApproveClaimRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new ApproveClaimResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            claim.Status = ClaimStatus.Approved;
            claim.ApprovedAmount = request.ApprovedAmount ?? claim.ClaimedAmount;
            claim.ApprovedAt = now;
            claim.UpdatedAt = now;
            claim.Approvals.Add(new ClaimApproval
            {
                ApprovalId = Guid.NewGuid().ToString("N"),
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = "L1",
                ApprovalLevel = 1,
                Decision = ApprovalDecision.Approved,
                ApprovedAmount = claim.ApprovedAmount,
                Notes = request.Notes,
                ApprovedAt = now,
                CreatedAt = now,
                ApprovedCurrency = claim.ApprovedAmount?.Currency ?? "BDT"
            });

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));
            return new ApproveClaimResponse { Message = "Claim approved" };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to approve claim {ClaimId}", request.ClaimId);
            return new ApproveClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RejectClaimResponse> RejectClaim(RejectClaimRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new RejectClaimResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            claim.Status = ClaimStatus.Rejected;
            claim.RejectionReason = request.Reason;
            claim.UpdatedAt = now;
            claim.Approvals.Add(new ClaimApproval
            {
                ApprovalId = Guid.NewGuid().ToString("N"),
                ClaimId = claim.ClaimId,
                ApproverId = request.ApproverId,
                ApproverRole = "L1",
                ApprovalLevel = 1,
                Decision = ApprovalDecision.Rejected,
                ApprovedAmount = NewMoney(0),
                Notes = request.Reason,
                ApprovedAt = now,
                CreatedAt = now,
                ApprovedCurrency = "BDT"
            });

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));
            return new RejectClaimResponse { Message = "Claim rejected" };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to reject claim {ClaimId}", request.ClaimId);
            return new RejectClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<SettleClaimResponse> SettleClaim(SettleClaimRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new SettleClaimResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            var paid = claim.ApprovedAmount?.Amount > 0 ? claim.ApprovedAmount : claim.ClaimedAmount;

            claim.Status = ClaimStatus.Settled;
            claim.SettledAmount = paid;
            claim.SettledAt = now;
            claim.UpdatedAt = now;

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));

            return new SettleClaimResponse
            {
                Message = "Claim settled",
                SettledAmount = paid,
                PaymentId = $"PAY-{Guid.NewGuid():N}"[..18]
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to settle claim {ClaimId}", request.ClaimId);
            return new SettleClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RequestMoreDocumentsResponse> RequestMoreDocuments(RequestMoreDocumentsRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new RequestMoreDocumentsResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            claim.Status = ClaimStatus.PendingDocuments;
            claim.ProcessorNotes = string.Join(",", request.RequiredDocumentTypes);
            claim.InAppMessages = request.Message;
            claim.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));
            return new RequestMoreDocumentsResponse
            {
                Message = "Additional documents requested"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to request more documents for claim {ClaimId}", request.ClaimId);
            return new RequestMoreDocumentsResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<DisputeClaimResponse> DisputeClaim(DisputeClaimRequest request, ServerCallContext context)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, GetCancellationToken(context));
            if (claim is null)
            {
                return new DisputeClaimResponse
                {
                    Error = BuildError("NOT_FOUND", "Claim not found")
                };
            }

            claim.Status = ClaimStatus.Disputed;
            claim.InAppMessages = request.DisputeReason;
            claim.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

            await _dataGateway.UpdateClaimAsync(claim, GetCancellationToken(context));
            return new DisputeClaimResponse
            {
                DisputeId = $"DSP-{Guid.NewGuid():N}"[..16],
                Message = "Claim dispute registered"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to dispute claim {ClaimId}", request.ClaimId);
            return new DisputeClaimResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static Money NewMoney(long amount, string currency = "BDT") => new() { Amount = amount, Currency = currency };

    private static Error BuildError(string code, string message) => new() { Code = code, Message = message };

    private static DateTime ParseDateOrNow(string input)
    {
        return DateTime.TryParse(input, out var dt)
            ? DateTime.SpecifyKind(dt, DateTimeKind.Utc)
            : DateTime.UtcNow;
    }

    private static string ComputeSha256(string text)
    {
        var bytes = Encoding.UTF8.GetBytes(text);
        return ComputeSha256(bytes);
    }

    private static string ComputeSha256(byte[] bytes)
    {
        using var sha = SHA256.Create();
        return Convert.ToHexString(sha.ComputeHash(bytes)).ToLowerInvariant();
    }
}
