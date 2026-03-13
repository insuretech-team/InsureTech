using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Underwriting.Infrastructure.Repositories;

public class BeneficiaryRepository : IBeneficiaryRepository
{
    private readonly UnderwritingDbContext _context;

    public BeneficiaryRepository(UnderwritingDbContext context)
    {
        _context = context;
    }

    public async Task<Beneficiary?> GetByIdAsync(Guid id)
    {
        return await _context.Beneficiaries
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.Id == id);
    }

    public async Task<Beneficiary?> GetByCodeAsync(string code)
    {
        return await _context.Beneficiaries
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.Code == code);
    }

    public async Task<IEnumerable<Beneficiary>> ListAsync(string? type = null, string? status = null, int page = 1, int pageSize = 10)
    {
        var query = _context.Beneficiaries
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .AsQueryable();

        if (!string.IsNullOrEmpty(type))
        {
            query = query.Where(b => b.Type.ToString().ToUpper() == type.ToUpper());
        }

        if (!string.IsNullOrEmpty(status))
        {
            query = query.Where(b => b.Status.ToString().ToUpper() == status.ToUpper());
        }

        return await query
            .OrderByDescending(b => b.CreatedAt)
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
            query = query.Where(b => b.Status.ToString().ToUpper() == status.ToUpper());
        }

        return await query.CountAsync();
    }

    public async Task AddAsync(Beneficiary beneficiary)
    {
        await _context.Beneficiaries.AddAsync(beneficiary);
        await _context.SaveChangesAsync();
    }

    public async Task UpdateAsync(Beneficiary beneficiary)
    {
        _context.Beneficiaries.Update(beneficiary);
        await _context.SaveChangesAsync();
    }

    public async Task<string> GetNextSequenceAsync()
    {
        var count = await _context.Beneficiaries.IgnoreQueryFilters().CountAsync();
        return $"BEN-{(count + 1):D6}";
    }
}
