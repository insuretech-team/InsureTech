using MediatR;
using Microsoft.Extensions.Logging;
using Microsoft.EntityFrameworkCore;
using Insuretech.Beneficiary.Services.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Beneficiary.Application.Queries;

public sealed class GetBeneficiaryQueryHandler : IRequestHandler<GetBeneficiaryQuery, GetBeneficiaryResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repo;
    private readonly ILogger<GetBeneficiaryQueryHandler> _logger;

    public GetBeneficiaryQueryHandler(IRepository<BeneficiaryEntity> repo, ILogger<GetBeneficiaryQueryHandler> logger)
    {
        _repo = repo;
        _logger = logger;
    }

    public async Task<GetBeneficiaryResponse> Handle(GetBeneficiaryQuery request, CancellationToken cancellationToken)
    {
        try
        {
            // BeneficiaryEntity includes navigation properties if registered correctly
            var e = await _repo.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
            if (e == null) return new GetBeneficiaryResponse { Error = new Insuretech.Common.V1.Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

            var response = new GetBeneficiaryResponse { Beneficiary = MapToProto(e) };
            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "GetBeneficiary failed");
            throw;
        }
    }

    internal static Insuretech.Beneficiary.Entity.V1.Beneficiary MapToProto(BeneficiaryEntity e)
    {
        var b = new Insuretech.Beneficiary.Entity.V1.Beneficiary
        {
            BeneficiaryId = e.BeneficiaryId.ToString(),
            UserId = e.UserId.ToString(),
            Code = e.Code,
            RiskScore = e.RiskScore ?? "",
            ReferralCode = e.ReferralCode ?? "",
            ReferredBy = e.ReferredBy?.ToString() ?? "",
            PartnerId = e.PartnerId?.ToString() ?? ""
        };

        if (System.Enum.TryParse<Insuretech.Beneficiary.Entity.V1.BeneficiaryType>(e.Type, true, out var bt)) b.Type = bt;
        if (System.Enum.TryParse<Insuretech.Beneficiary.Entity.V1.BeneficiaryStatus>(e.Status, true, out var bs)) b.Status = bs;
        if (System.Enum.TryParse<Insuretech.Beneficiary.Entity.V1.KYCStatus>(e.KycStatus, true, out var ks)) b.KycStatus = ks;
        if (e.KycCompletedAt.HasValue) b.KycCompletedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.KycCompletedAt.Value, DateTimeKind.Utc));

        return b;
    }
}

public sealed class ListBeneficiariesQueryHandler : IRequestHandler<ListBeneficiariesQuery, ListBeneficiariesResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repo;

    public ListBeneficiariesQueryHandler(IRepository<BeneficiaryEntity> repo) => _repo = repo;

    public async Task<ListBeneficiariesResponse> Handle(ListBeneficiariesQuery request, CancellationToken cancellationToken)
    {
        Expression<Func<BeneficiaryEntity, bool>> predicate = x => x.DeletedAt == null;
        
        if (!string.IsNullOrEmpty(request.Type))
            predicate = Combine(predicate, x => x.Type == request.Type);
        
        if (!string.IsNullOrEmpty(request.Status))
            predicate = Combine(predicate, x => x.Status == request.Status);

        var (items, total) = await _repo.GetPagedAsync(
            request.Page, request.PageSize, predicate, 
            x => x.CreatedAt, true, cancellationToken);

        var response = new ListBeneficiariesResponse { TotalCount = total };
        foreach (var item in items)
            response.Beneficiaries.Add(GetBeneficiaryQueryHandler.MapToProto(item));
        
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
