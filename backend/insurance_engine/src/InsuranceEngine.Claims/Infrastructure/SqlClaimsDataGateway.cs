using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using ClaimEntity = InsuranceEngine.SharedKernel.Persistence.Entities.ClaimEntity;

namespace InsuranceEngine.Claims.Infrastructure;

public class SqlClaimsDataGateway : IClaimDataGateway
{
    private readonly ClaimsDbContext _context;
    private readonly ILogger<SqlClaimsDataGateway> _logger;

    public SqlClaimsDataGateway(ClaimsDbContext context, ILogger<SqlClaimsDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<GetClaimResponse> GetClaimAsync(string claimId, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims
            .Include(c => c.Documents)
            .Include(c => c.Approvals)
            .FirstOrDefaultAsync(c => c.ClaimId == id, ct);

        if (claim == null)
        {
            return new GetClaimResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        return new GetClaimResponse { Claim = MapToProto(claim) };
    }

    public async Task<SubmitClaimResponse> SubmitClaimAsync(SubmitClaimRequest request, CancellationToken ct = default)
    {
        var claimId = Guid.NewGuid();
        var claimNumber = $"CLM-{DateTime.UtcNow.Year}-{DateTime.UtcNow:MMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        var now = DateTime.UtcNow;

        var claim = new ClaimEntity
        {
            ClaimId = claimId,
            ClaimNumber = claimNumber,
            PolicyId = Guid.TryParse(request.PolicyId, out var pid) ? pid : Guid.Empty,
            CustomerId = Guid.TryParse(request.CustomerId, out var cid) ? cid : Guid.Empty,
            Status = "SUBMITTED",
            Type = request.Type.ToString(),
            ClaimedAmount = request.ClaimedAmount?.Amount ?? 0,
            ClaimedCurrency = request.ClaimedAmount?.Currency ?? "BDT",
            IncidentDate = DateTime.UtcNow,
            IncidentDescription = request.IncidentDescription ?? "",
            SubmittedAt = now,
            AppealOptionAvailable = true,
            ProcessingType = "MANUAL",
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.Claims.Add(claim);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Submitted claim {ClaimNumber}", claimNumber);

        return new SubmitClaimResponse
        {
            ClaimId = claimId.ToString(),
            ClaimNumber = claimNumber
        };
    }

    public async Task<ApproveClaimResponse> ApproveClaimAsync(string claimId, string notes, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims
            .Include(c => c.Approvals)
            .FirstOrDefaultAsync(c => c.ClaimId == id, ct);

        if (claim == null)
        {
            return new ApproveClaimResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        claim.Status = "APPROVED";
        claim.ApprovedAt = DateTime.UtcNow;
        claim.UpdatedAt = DateTime.UtcNow;
        claim.ProcessorNotes = notes;

        var approval = new ClaimApprovalEntity
        {
            ApprovalId = Guid.NewGuid(),
            ClaimId = id,
            ApproverId = Guid.Empty,
            ApproverRole = "SYSTEM",
            ApprovalLevel = 1,
            Decision = "APPROVED",
            ApprovedAmount = claim.ClaimedAmount,
            ApprovedCurrency = claim.ClaimedCurrency,
            Notes = notes,
            ApprovedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow
        };

        _context.ClaimApprovals.Add(approval);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Approved claim {ClaimId}", claimId);

        return new ApproveClaimResponse { Message = "Claim approved successfully" };
    }

    public async Task<RejectClaimResponse> RejectClaimAsync(string claimId, string reason, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims.FindAsync([id], ct);

        if (claim == null)
        {
            return new RejectClaimResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        claim.Status = "REJECTED";
        claim.RejectionReason = reason;
        claim.AppealOptionAvailable = true;
        claim.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Rejected claim {ClaimId}", claimId);

        return new RejectClaimResponse { Message = "Claim rejected" };
    }

    public async Task<SettleClaimResponse> SettleClaimAsync(string claimId, string paymentMethod, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims.FindAsync([id], ct);

        if (claim == null)
        {
            return new SettleClaimResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        var settledAmount = claim.ApprovedAmount ?? claim.ClaimedAmount;
        claim.Status = "SETTLED";
        claim.SettledAmount = settledAmount;
        claim.SettledCurrency = claim.ApprovedCurrency ?? claim.ClaimedCurrency;
        claim.SettledAt = DateTime.UtcNow;
        claim.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Settled claim {ClaimId}, Amount: {Amount}", claimId, settledAmount);

        return new SettleClaimResponse
        {
            SettledAmount = new Money { Amount = settledAmount, Currency = claim.SettledCurrency }
        };
    }

    public async Task<UploadDocumentResponse> UploadDocumentAsync(string claimId, string fileName, string documentType, string documentUrl, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims.FindAsync([id], ct);

        if (claim == null)
        {
            return new UploadDocumentResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        var documentId = Guid.NewGuid();
        var document = new ClaimDocumentEntity
        {
            DocumentId = documentId,
            ClaimId = id,
            DocumentType = documentType,
            FileUrl = documentUrl,
            UploadedAt = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        _context.ClaimDocuments.Add(document);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Uploaded document {DocumentId} for claim {ClaimId}", documentId, claimId);

        return new UploadDocumentResponse
        {
            DocumentId = documentId.ToString(),
            DocumentUrl = documentUrl
        };
    }

    public async Task<RequestMoreDocumentsResponse> RequestMoreDocumentsAsync(string claimId, string message, List<string> requiredDocumentTypes, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims.FindAsync([id], ct);

        if (claim == null)
        {
            return new RequestMoreDocumentsResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        claim.InAppMessages = message;
        claim.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Requested more documents for claim {ClaimId}", claimId);

        return new RequestMoreDocumentsResponse { Message = message };
    }

    public async Task<DisputeClaimResponse> DisputeClaimAsync(string claimId, string disputeReason, string customerId, CancellationToken ct = default)
    {
        var id = Guid.TryParse(claimId, out var cid) ? cid : Guid.Empty;
        var claim = await _context.Claims.FindAsync([id], ct);

        if (claim == null)
        {
            return new DisputeClaimResponse { Error = new Error { Code = "NOT_FOUND", Message = "Claim not found" } };
        }

        claim.Status = "UNDER_REVIEW";
        claim.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Dispute filed for claim {ClaimId}", claimId);

        return new DisputeClaimResponse
        {
            DisputeId = $"DSP-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}",
            Message = "Dispute submitted successfully"
        };
    }

    private static Claim MapToProto(ClaimEntity entity)
    {
        var proto = new Claim
        {
            ClaimId = entity.ClaimId.ToString(),
            ClaimNumber = entity.ClaimNumber,
            PolicyId = entity.PolicyId.ToString(),
            CustomerId = entity.CustomerId.ToString(),
            Status = Enum.TryParse<ClaimStatus>(entity.Status, true, out var status) ? status : ClaimStatus.Unspecified,
            Type = Enum.TryParse<ClaimType>(entity.Type, true, out var type) ? type : ClaimType.Unspecified,
            ClaimedAmount = new Money { Amount = entity.ClaimedAmount, Currency = entity.ClaimedCurrency },
            IncidentDate = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.IncidentDate),
            IncidentDescription = entity.IncidentDescription,
            SubmittedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.SubmittedAt)
        };

        if (entity.ApprovedAmount.HasValue)
        {
            proto.ApprovedAmount = new Money { Amount = entity.ApprovedAmount.Value, Currency = entity.ApprovedCurrency };
        }

        if (entity.SettledAmount.HasValue)
        {
            proto.SettledAmount = new Money { Amount = entity.SettledAmount.Value, Currency = entity.SettledCurrency };
        }

        if (entity.ApprovedAt.HasValue)
        {
            proto.ApprovedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.ApprovedAt.Value);
        }

        if (entity.SettledAt.HasValue)
        {
            proto.SettledAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.SettledAt.Value);
        }

        proto.AppealOptionAvailable = entity.AppealOptionAvailable;

        if (!string.IsNullOrEmpty(entity.RejectionReason))
        {
            proto.RejectionReason = entity.RejectionReason;
        }

        return proto;
    }
}
