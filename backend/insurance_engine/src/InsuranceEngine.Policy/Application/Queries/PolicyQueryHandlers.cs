using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Products.Domain.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Policy.Application.Queries;

public sealed class ListUserPoliciesQueryHandler : IRequestHandler<ListUserPoliciesQuery, ListUserPoliciesResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly ILogger<ListUserPoliciesQueryHandler> _logger;

    public ListUserPoliciesQueryHandler(IRepository<PolicyEntity> repository, ILogger<ListUserPoliciesQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<ListUserPoliciesResponse> Handle(ListUserPoliciesQuery request, CancellationToken cancellationToken)
    {
        try
        {
            Expression<Func<PolicyEntity, bool>>? predicate = null;

            if (!string.IsNullOrEmpty(request.CustomerId))
            {
                var customerId = Guid.Parse(request.CustomerId);
                predicate = p => p.CustomerId == customerId && p.DeletedAt == null;
            }
            else
            {
                predicate = p => p.DeletedAt == null;
            }

            if (!string.IsNullOrEmpty(request.Status))
            {
                var status = request.Status;
                predicate = Combine(predicate, p => p.Status == status);
            }

            if (!string.IsNullOrEmpty(request.ProductId))
            {
                var productId = Guid.Parse(request.ProductId);
                predicate = Combine(predicate, p => p.ProductId == productId);
            }

            var (items, totalCount) = await _repository.GetPagedAsync(
                page: request.Page,
                pageSize: request.PageSize,
                predicate: predicate,
                orderBy: p => p.CreatedAt,
                descending: true,
                cancellationToken: cancellationToken
            );

            var response = new ListUserPoliciesResponse { TotalCount = totalCount };

            foreach (var entity in items)
            {
                response.Policies.Add(MapToProto(entity));
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list policies for customer {CustomerId}", request.CustomerId);
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

    private static Insuretech.Policy.Entity.V1.Policy MapToProto(PolicyEntity e)
    {
        var p = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = e.PolicyId.ToString(),
            PolicyNumber = e.PolicyNumber,
            CustomerId = e.CustomerId.ToString(),
            ProductId = e.ProductId.ToString(),
            PartnerId = e.PartnerId?.ToString() ?? "",
            AgentId = e.AgentId?.ToString() ?? "",
            TenureMonths = e.TenureMonths,
            PremiumAmount = new Money { Amount = e.PremiumAmount, Currency = e.PremiumCurrency },
            SumInsured = new Money { Amount = e.SumInsured, Currency = e.SumInsuredCurrency },
            PolicyDocumentUrl = e.PolicyDocumentUrl ?? ""
        };

        if (System.Enum.TryParse<Insuretech.Policy.Entity.V1.PolicyStatus>(e.Status, true, out var s)) p.Status = s;

        p.StartDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.StartDate, DateTimeKind.Utc));
        p.EndDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.EndDate, DateTimeKind.Utc));
        p.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.CreatedAt, DateTimeKind.Utc));
        if (e.IssuedAt.HasValue) p.IssuedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.IssuedAt.Value, DateTimeKind.Utc));

        return p;
    }
}

public sealed class GetPolicyQueryHandler : IRequestHandler<GetPolicyQuery, GetPolicyResponse>
{
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly IRepository<PolicyRiderEntity> _riderRepository;
    private readonly ILogger<GetPolicyQueryHandler> _logger;

    public GetPolicyQueryHandler(
        IRepository<PolicyEntity> policyRepository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        IRepository<PolicyRiderEntity> riderRepository,
        ILogger<GetPolicyQueryHandler> logger)
    {
        _policyRepository = policyRepository;
        _nomineeRepository = nomineeRepository;
        _riderRepository = riderRepository;
        _logger = logger;
    }

    public async Task<GetPolicyResponse> Handle(GetPolicyQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var entity = await _policyRepository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (entity == null)
            {
                return new GetPolicyResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            var policy = MapToProto(entity);

            // Load nominees
            var nominees = await _nomineeRepository.FindAsync(n => n.PolicyId == entity.PolicyId, cancellationToken);
            foreach (var n in nominees)
            {
                policy.Nominees.Add(new Insuretech.Policy.Entity.V1.Nominee
                {
                    NomineeId = n.NomineeId.ToString(),
                    PolicyId = n.PolicyId.ToString(),
                    FullName = n.FullName,
                    Relationship = n.Relationship,
                    SharePercentage = n.SharePercentage,
                    DateOfBirth = Timestamp.FromDateTime(DateTime.SpecifyKind(n.DateOfBirth, DateTimeKind.Utc)),
                    NidNumber = n.NidNumber ?? "",
                    PhoneNumber = n.PhoneNumber ?? ""
                });
            }

            // Load riders
            var riders = await _riderRepository.FindAsync(r => r.PolicyId == entity.PolicyId, cancellationToken);
            foreach (var r in riders)
            {
                policy.Riders.Add(new Insuretech.Policy.Entity.V1.Rider
                {
                    RiderId = r.RiderId.ToString(),
                    PolicyId = r.PolicyId.ToString(),
                    RiderName = r.RiderName,
                    PremiumAmount = new Money { Amount = r.PremiumAmount, Currency = r.PremiumCurrency },
                    CoverageAmount = new Money { Amount = r.CoverageAmount, Currency = r.CoverageCurrency }
                });
            }

            return new GetPolicyResponse { Policy = policy };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get policy {PolicyId}", request.PolicyId);
            throw;
        }
    }

    private static Insuretech.Policy.Entity.V1.Policy MapToProto(PolicyEntity e)
    {
        var p = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = e.PolicyId.ToString(),
            PolicyNumber = e.PolicyNumber,
            CustomerId = e.CustomerId.ToString(),
            ProductId = e.ProductId.ToString(),
            PartnerId = e.PartnerId?.ToString() ?? "",
            AgentId = e.AgentId?.ToString() ?? "",
            QuoteId = e.QuoteId?.ToString() ?? "",
            TenureMonths = e.TenureMonths,
            PremiumAmount = new Money { Amount = e.PremiumAmount, Currency = e.PremiumCurrency },
            SumInsured = new Money { Amount = e.SumInsured, Currency = e.SumInsuredCurrency },
            PolicyDocumentUrl = e.PolicyDocumentUrl ?? "",
            PaymentFrequency = e.PaymentFrequency ?? "",
            ProviderName = e.ProviderName ?? "",
            OccupationRiskClass = e.OccupationRiskClass ?? "",
            HasExistingPolicies = e.HasExistingPolicies,
            ClaimsHistorySummary = e.ClaimsHistorySummary ?? ""
        };

        if (System.Enum.TryParse<Insuretech.Policy.Entity.V1.PolicyStatus>(e.Status, true, out var s)) p.Status = s;

        p.StartDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.StartDate, DateTimeKind.Utc));
        p.EndDate = Timestamp.FromDateTime(DateTime.SpecifyKind(e.EndDate, DateTimeKind.Utc));
        p.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.CreatedAt, DateTimeKind.Utc));
        p.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.UpdatedAt, DateTimeKind.Utc));
        if (e.IssuedAt.HasValue) p.IssuedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.IssuedAt.Value, DateTimeKind.Utc));

        if (e.VatTax.HasValue) p.VatTax = new Money { Amount = e.VatTax.Value, Currency = "BDT" };
        if (e.ServiceFee.HasValue) p.ServiceFee = new Money { Amount = e.ServiceFee.Value, Currency = "BDT" };
        if (e.TotalPayable.HasValue) p.TotalPayable = new Money { Amount = e.TotalPayable.Value, Currency = "BDT" };

        return p;
    }
}
