using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Commission.Application.Queries;

public sealed class GetCommissionQueryHandler : IRequestHandler<GetCommissionQuery, GetCommissionResponse>
{
    private readonly IRepository<CommissionEntity> _repository;

    public GetCommissionQueryHandler(IRepository<CommissionEntity> repository) => _repository = repository;

    public async Task<GetCommissionResponse> Handle(GetCommissionQuery request, CancellationToken cancellationToken)
    {
        var entity = await _repository.GetByIdAsync(Guid.Parse(request.CommissionId), cancellationToken);
        if (entity == null)
            return new GetCommissionResponse { Error = new Error { Code = "NOT_FOUND", Message = "Commission not found" } };

        return new GetCommissionResponse { Commission = MapToProto(entity) };
    }

    internal static Insuretech.Partner.Entity.V1.Commission MapToProto(CommissionEntity e)
    {
        var c = new Insuretech.Partner.Entity.V1.Commission
        {
            CommissionId = e.CommissionId.ToString(),
            PolicyId = e.PolicyId.ToString(),
            CommissionRate = (double)e.CommissionRate,
            CommissionAmount = new Money { Amount = e.CommissionAmount, Currency = e.CommissionCurrency }
        };

        if (e.PartnerId.HasValue) c.PartnerId = e.PartnerId.Value.ToString();
        if (e.AgentId.HasValue) c.AgentId = e.AgentId.Value.ToString();

        if (System.Enum.TryParse<Insuretech.Partner.Entity.V1.CommissionType>(e.CommissionType, true, out var ct)) c.Type = ct;
        if (System.Enum.TryParse<Insuretech.Partner.Entity.V1.CommissionStatus>(e.Status, true, out var cs)) c.Status = cs;

        if (e.PaidAt.HasValue) c.PaidAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.PaidAt.Value, DateTimeKind.Utc));
        c.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.CreatedAt, DateTimeKind.Utc));
        c.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.UpdatedAt, DateTimeKind.Utc));

        return c;
    }
}

public sealed class ListCommissionsQueryHandler : IRequestHandler<ListCommissionsQuery, ListCommissionsResponse>
{
    private readonly IRepository<CommissionEntity> _repository;

    public ListCommissionsQueryHandler(IRepository<CommissionEntity> repository) => _repository = repository;

    public async Task<ListCommissionsResponse> Handle(ListCommissionsQuery request, CancellationToken cancellationToken)
    {
        Expression<Func<CommissionEntity, bool>> predicate = c => c.DeletedAt == null;
        var recipientId = Guid.Parse(request.RecipientId);
        predicate = Combine(predicate, c => c.PartnerId == recipientId || c.AgentId == recipientId);

        if (!string.IsNullOrEmpty(request.Status))
        {
            var status = request.Status;
            predicate = Combine(predicate, c => c.Status == status);
        }

        var (items, totalCount) = await _repository.GetPagedAsync(
            page: request.Page, pageSize: request.PageSize,
            predicate: predicate, orderBy: c => c.CreatedAt, descending: true,
            cancellationToken: cancellationToken);

        var response = new ListCommissionsResponse
        {
            TotalCount = totalCount,
            TotalAmount = new Money { Amount = items.Sum(c => c.CommissionAmount), Currency = "BDT" }
        };
        foreach (var e in items)
            response.Commissions.Add(GetCommissionQueryHandler.MapToProto(e));
        return response;
    }

    private static Expression<Func<T, bool>> Combine<T>(Expression<Func<T, bool>> e1, Expression<Func<T, bool>> e2)
    {
        var p = Expression.Parameter(typeof(T));
        var v1 = new ReplaceVisitor(e1.Parameters[0], p);
        var v2 = new ReplaceVisitor(e2.Parameters[0], p);
        return Expression.Lambda<Func<T, bool>>(Expression.AndAlso(v1.Visit(e1.Body)!, v2.Visit(e2.Body)!), p);
    }
    private class ReplaceVisitor(Expression old, Expression @new) : ExpressionVisitor
    { public override Expression Visit(Expression? n) => n == old ? @new : base.Visit(n)!; }
}

public sealed class GetCommissionStatementQueryHandler : IRequestHandler<GetCommissionStatementQuery, GetCommissionStatementResponse>
{
    private readonly IRepository<CommissionEntity> _repository;

    public GetCommissionStatementQueryHandler(IRepository<CommissionEntity> repository) => _repository = repository;

    public async Task<GetCommissionStatementResponse> Handle(GetCommissionStatementQuery request, CancellationToken cancellationToken)
    {
        var recipientId = Guid.Parse(request.RecipientId);
        var all = await _repository.FindAsync(c => (c.PartnerId == recipientId || c.AgentId == recipientId) && c.DeletedAt == null, cancellationToken);

        var totalEarned = all.Sum(c => c.CommissionAmount);
        var totalPaid = all.Where(c => c.Status == "PAID").Sum(c => c.CommissionAmount);
        var pending = totalEarned - totalPaid;

        var byType = all.GroupBy(c => c.CommissionType).Select(g => new CommissionSummary
        {
            Type = g.Key,
            Count = g.Count(),
            TotalAmount = new Money { Amount = g.Sum(c => c.CommissionAmount), Currency = "BDT" }
        });

        var response = new GetCommissionStatementResponse
        {
            RecipientId = request.RecipientId,
            PeriodStart = request.PeriodStart,
            PeriodEnd = request.PeriodEnd,
            TotalEarned = new Money { Amount = totalEarned, Currency = "BDT" },
            TotalPaid = new Money { Amount = totalPaid, Currency = "BDT" },
            PendingAmount = new Money { Amount = pending, Currency = "BDT" }
        };
        response.ByType.AddRange(byType);
        return response;
    }
}

public sealed class GetRevenueShareReportQueryHandler : IRequestHandler<GetRevenueShareReportQuery, GetRevenueShareReportResponse>
{
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly IRepository<CommissionEntity> _commissionRepository;

    public GetRevenueShareReportQueryHandler(IRepository<PolicyEntity> policyRepository, IRepository<CommissionEntity> commissionRepository)
    {
        _policyRepository = policyRepository;
        _commissionRepository = commissionRepository;
    }

    public async Task<GetRevenueShareReportResponse> Handle(GetRevenueShareReportQuery request, CancellationToken cancellationToken)
    {
        // Simplified revenue share report
        var allCommissions = await _commissionRepository.FindAsync(c => c.DeletedAt == null, cancellationToken);
        var allPolicies = await _policyRepository.FindAsync(p => p.DeletedAt == null, cancellationToken);

        var totalGross = allPolicies.Sum(p => p.PremiumAmount);
        var totalCommissions = allCommissions.Sum(c => c.CommissionAmount);
        var platformShare = (long)(totalGross * 0.25m); // 25% platform share
        var insurerShare = totalGross - platformShare - totalCommissions;

        return new GetRevenueShareReportResponse
        {
            InsurerId = request.InsurerId,
            TotalGrossPremium = new Money { Amount = totalGross, Currency = "BDT" },
            TotalPlatformShare = new Money { Amount = platformShare, Currency = "BDT" },
            TotalInsurerShare = new Money { Amount = insurerShare, Currency = "BDT" },
            PolicyCount = allPolicies.Count
        };
    }
}
