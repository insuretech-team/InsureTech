using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.Beneficiary.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Beneficiary.Infrastructure.Repositories;

public class BeneficiaryRepository : IBeneficiaryRepository
{
    private readonly BeneficiaryDbContext _context;

    public BeneficiaryRepository(BeneficiaryDbContext context)
    {
        _context = context;
    }

    public async Task<Domain.Entities.Beneficiary?> GetByIdAsync(Guid id)
    {
        return await _context.Beneficiaries
            .Include(b => b.Individual)
            .Include(b => b.Business)
            .FirstOrDefaultAsync(b => b.Id == id);
    }

    public async Task<Domain.Entities.Beneficiary?> GetByCodeAsync(string code)
    {
        return await _context.Beneficiaries
            .Include(b => b.Individual)
            .Include(b => b.Business)
            .FirstOrDefaultAsync(b => b.Code == code);
    }

    public async Task<IEnumerable<Domain.Entities.Beneficiary>> ListAsync(string? type = null, string? status = null, int page = 1, int pageSize = 10)
    {
        var query = _context.Beneficiaries
            .Include(b => b.Individual)
            .Include(b => b.Business)
            .AsQueryable();

        if (!string.IsNullOrEmpty(type))
        {
            query = query.Where(b => b.Type.ToString().ToUpper() == type.ToUpper());
        }

        if (!string.IsNullOrEmpty(status))
        {
            query = query.Where(b => b.Status.Value.ToUpper() == status.ToUpper());
        }

        return await query
            .OrderByDescending(b => b.AuditInfo.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync();
    }

    public async Task<int> GetTotalCountAsync(string? type = null, string? status = null)
    {
        var query = _context.Beneficiaries.AsQueryable();

        if (!string.IsNullOrEmpty(type))
        {
            query = query.Where(b => b.Type.ToString().ToUpper() == type.ToUpper());
        }

        if (!string.IsNullOrEmpty(status))
        {
            query = query.Where(b => b.Status.Value.ToUpper() == status.ToUpper());
        }

        return await query.CountAsync();
    }

    public async Task AddAsync(Domain.Entities.Beneficiary beneficiary)
    {
        await _context.Beneficiaries.AddAsync(beneficiary);
        await _context.SaveChangesAsync();
    }

    public async Task UpdateAsync(Domain.Entities.Beneficiary beneficiary)
    {
        _context.Beneficiaries.Update(beneficiary);
        await _context.SaveChangesAsync();
    }
}
