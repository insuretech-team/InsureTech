using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Claims.Services.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Claims.Application.Queries;

public sealed class ListUserClaimsQueryHandler : IRequestHandler<ListUserClaimsQuery, ListUserClaimsResponse>
{
    private readonly IRepository<ClaimEntity> _repository;
    private readonly ILogger<ListUserClaimsQueryHandler> _logger;

    public ListUserClaimsQueryHandler(IRepository<ClaimEntity> repository, ILogger<ListUserClaimsQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<ListUserClaimsResponse> Handle(ListUserClaimsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            Expression<Func<ClaimEntity, bool>> predicate = c => c.DeletedAt == null;

            if (!string.IsNullOrEmpty(request.CustomerId))
            {
                var customerId = Guid.Parse(request.CustomerId);
                predicate = Combine(predicate, c => c.CustomerId == customerId);
            }

            if (!string.IsNullOrEmpty(request.PolicyId))
            {
                var policyId = Guid.Parse(request.PolicyId);
                predicate = Combine(predicate, c => c.PolicyId == policyId);
            }

            if (!string.IsNullOrEmpty(request.Status))
            {
                var status = request.Status;
                predicate = Combine(predicate, c => c.Status == status);
            }

            var (items, totalCount) = await _repository.GetPagedAsync(
                page: request.Page,
                pageSize: request.PageSize,
                predicate: predicate,
                orderBy: c => c.CreatedAt,
                descending: true,
                cancellationToken: cancellationToken
            );

            var response = new ListUserClaimsResponse { TotalCount = totalCount };
            foreach (var entity in items)
            {
                response.Claims.Add(MapToProto(entity));
            }
            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list claims");
            throw;
        }
    }

    private static Expression<Func<T, bool>> Combine<T>(Expression<Func<T, bool>> expr1, Expression<Func<T, bool>> expr2)
    {
        var parameter = Expression.Parameter(typeof(T));
        var leftVisitor = new ReplaceExpressionVisitor(expr1.Parameters[0], parameter);
        var left = leftVisitor.Visit(expr1.Body);
        var rightVisitor = new ReplaceExpressionVisitor(expr2.Parameters[0], parameter);
        var right = rightVisitor.Visit(expr2.Body);
        return Expression.Lambda<Func<T, bool>>(Expression.AndAlso(left!, right!), parameter);
    }

    private class ReplaceExpressionVisitor(Expression oldValue, Expression newValue) : ExpressionVisitor
    {
        public override Expression Visit(Expression? node) => node == oldValue ? newValue : base.Visit(node)!;
    }

    private static Insuretech.Claims.Entity.V1.Claim MapToProto(ClaimEntity e) => MapClaimToProto(e);

    internal static Insuretech.Claims.Entity.V1.Claim MapClaimToProto(ClaimEntity e)
    {
        var c = new Insuretech.Claims.Entity.V1.Claim
        {
            ClaimId = e.ClaimId.ToString(),
            ClaimNumber = e.ClaimNumber,
            PolicyId = e.PolicyId.ToString(),
            CustomerId = e.CustomerId.ToString(),
            IncidentDescription = e.IncidentDescription,
            RejectionReason = e.RejectionReason ?? "",
            PlaceOfIncident = e.PlaceOfIncident ?? "",
            AppealOptionAvailable = e.AppealOptionAvailable,
            ProcessorNotes = e.ProcessorNotes ?? "",
            ClaimedAmount = new Money { Amount = e.ClaimedAmount, Currency = e.ClaimedCurrency }
        };

        if (System.Enum.TryParse<ClaimStatus>(e.Status, true, out var s)) c.Status = s;
        if (System.Enum.TryParse<Insuretech.Claims.Entity.V1.ClaimType>(e.Type, true, out var t)) c.Type = t;
        if (System.Enum.TryParse<ClaimProcessingType>(e.ProcessingType, true, out var pt)) c.ProcessingType = pt;

        if (e.ApprovedAmount.HasValue) c.ApprovedAmount = new Money { Amount = e.ApprovedAmount.Value, Currency = e.ApprovedCurrency };
        if (e.SettledAmount.HasValue) c.SettledAmount = new Money { Amount = e.SettledAmount.Value, Currency = e.SettledCurrency };
        if (e.DeductibleAmount.HasValue) c.DeductibleAmount = new Money { Amount = e.DeductibleAmount.Value, Currency = "BDT" };
        if (e.CoPayAmount.HasValue) c.CoPayAmount = new Money { Amount = e.CoPayAmount.Value, Currency = "BDT" };

        c.IncidentDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.IncidentDate, DateTimeKind.Utc));
        c.SubmittedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.SubmittedAt, DateTimeKind.Utc));
        c.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.CreatedAt, DateTimeKind.Utc));
        c.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.UpdatedAt, DateTimeKind.Utc));
        if (e.ApprovedAt.HasValue) c.ApprovedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.ApprovedAt.Value, DateTimeKind.Utc));
        if (e.SettledAt.HasValue) c.SettledAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.SettledAt.Value, DateTimeKind.Utc));

        return c;
    }
}

public sealed class GetClaimQueryHandler : IRequestHandler<GetClaimQuery, GetClaimResponse>
{
    private readonly IRepository<ClaimEntity> _claimRepository;
    private readonly IRepository<ClaimDocumentEntity> _documentRepository;
    private readonly IRepository<ClaimApprovalEntity> _approvalRepository;
    private readonly ILogger<GetClaimQueryHandler> _logger;

    public GetClaimQueryHandler(
        IRepository<ClaimEntity> claimRepository,
        IRepository<ClaimDocumentEntity> documentRepository,
        IRepository<ClaimApprovalEntity> approvalRepository,
        ILogger<GetClaimQueryHandler> logger)
    {
        _claimRepository = claimRepository;
        _documentRepository = documentRepository;
        _approvalRepository = approvalRepository;
        _logger = logger;
    }

    public async Task<GetClaimResponse> Handle(GetClaimQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var entity = await _claimRepository.GetByIdAsync(Guid.Parse(request.ClaimId), cancellationToken);
            if (entity == null)
            {
                return new GetClaimResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "CLAIM_NOT_FOUND", Message = "Claim not found" }
                };
            }

            var claim = ListUserClaimsQueryHandler.MapClaimToProto(entity);

            // Load documents
            var documents = await _documentRepository.FindAsync(d => d.ClaimId == entity.ClaimId, cancellationToken);
            foreach (var d in documents)
            {
                claim.Documents.Add(new ClaimDocument
                {
                    DocumentId = d.DocumentId.ToString(),
                    ClaimId = d.ClaimId.ToString(),
                    DocumentType = d.DocumentType,
                    FileUrl = d.FileUrl,
                    FileHash = d.FileHash,
                    Verified = d.Verified,
                    VerifiedBy = d.VerifiedBy?.ToString() ?? "",
                    UploadedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(d.UploadedAt, DateTimeKind.Utc)),
                    CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(d.CreatedAt, DateTimeKind.Utc))
                });
            }

            // Load approvals
            var approvals = await _approvalRepository.FindAsync(a => a.ClaimId == entity.ClaimId, cancellationToken);
            foreach (var a in approvals)
            {
                var approval = new ClaimApproval
                {
                    ApprovalId = a.ApprovalId.ToString(),
                    ClaimId = a.ClaimId.ToString(),
                    ApproverId = a.ApproverId.ToString(),
                    ApproverRole = a.ApproverRole,
                    ApprovalLevel = a.ApprovalLevel,
                    Notes = a.Notes ?? ""
                };
                if (System.Enum.TryParse<ApprovalDecision>(a.Decision, true, out var dec)) approval.Decision = dec;
                if (a.ApprovedAmount.HasValue) approval.ApprovedAmount = new Money { Amount = a.ApprovedAmount.Value, Currency = a.ApprovedCurrency };
                if (a.ApprovedAt.HasValue) approval.ApprovedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(a.ApprovedAt.Value, DateTimeKind.Utc));
                approval.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(a.CreatedAt, DateTimeKind.Utc));
                claim.Approvals.Add(approval);
            }

            return new GetClaimResponse { Claim = claim };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get claim {ClaimId}", request.ClaimId);
            throw;
        }
    }
}
