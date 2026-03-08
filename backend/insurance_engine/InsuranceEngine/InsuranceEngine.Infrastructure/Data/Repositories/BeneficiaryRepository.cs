using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;
using InsuranceEngine.Infrastructure.Data;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Infrastructure.Data.Repositories;

public class BeneficiaryRepository(InsuranceEngineDbContext dbContext) : IBeneficiaryRepository
{
    private readonly InsuranceEngineDbContext _dbContext = dbContext;

    public async Task<IReadOnlyList<Beneficiary>> GetBeneficiariesAsync(
        Guid? policyId,
        BeneficiaryStatus? status,
        BeneficiaryType? type,
        CancellationToken cancellationToken = default)
    {
        IQueryable<Beneficiary> query = BuildBeneficiariesQuery()
            .AsNoTracking();

        if (policyId.HasValue)
        {
            query = query.Where(b => b.PolicyId == policyId.Value);
        }

        if (status.HasValue)
        {
            query = query.Where(b => b.Status == status.Value);
        }

        if (type.HasValue)
        {
            query = query.Where(b => b.Type == type.Value);
        }

        return await query
            .OrderByDescending(b => b.AuditInfo!.CreatedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task<(IReadOnlyList<Beneficiary> Beneficiaries, int TotalCount)> GetBeneficiariesPageAsync(
        int page,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        IQueryable<Beneficiary> query = BuildBeneficiariesQuery()
            .AsNoTracking();

        int totalCount = await query.CountAsync(cancellationToken);
        List<Beneficiary> beneficiaries = await query
            .OrderByDescending(b => b.AuditInfo!.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);

        return (beneficiaries, totalCount);
    }

    public async Task<IReadOnlyList<Beneficiary>> GetBeneficiariesByPolicyAsync(
        Guid policyId,
        CancellationToken cancellationToken = default)
    {
        return await BuildBeneficiariesQuery()
            .AsNoTracking()
            .Where(b => b.PolicyId == policyId)
            .OrderByDescending(b => b.AuditInfo!.CreatedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task<Beneficiary?> GetBeneficiaryByIdAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return await BuildBeneficiariesQuery()
            .AsNoTracking()
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id, cancellationToken);
    }

    public async Task<Beneficiary?> GetBeneficiaryByIdForPolicyAsync(
        Guid policyId,
        Guid id,
        CancellationToken cancellationToken = default)
    {
        return await BuildBeneficiariesQuery()
            .AsNoTracking()
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id && b.PolicyId == policyId, cancellationToken);
    }

    public async Task<Beneficiary?> GetTrackedBeneficiaryByIdAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return await BuildBeneficiariesQuery()
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id, cancellationToken);
    }

    public async Task<Beneficiary?> GetTrackedBeneficiaryByIdForPolicyAsync(
        Guid policyId,
        Guid id,
        CancellationToken cancellationToken = default)
    {
        return await BuildBeneficiariesQuery()
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id && b.PolicyId == policyId, cancellationToken);
    }

    public Task AddBeneficiaryAsync(Beneficiary beneficiary, CancellationToken cancellationToken = default)
    {
        return _dbContext.Beneficiaries.AddAsync(beneficiary, cancellationToken).AsTask();
    }

    public Task SaveChangesAsync(CancellationToken cancellationToken = default)
    {
        return _dbContext.SaveChangesAsync(cancellationToken);
    }

    private IQueryable<Beneficiary> BuildBeneficiariesQuery()
    {
        return _dbContext.Beneficiaries
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails);
    }
}
