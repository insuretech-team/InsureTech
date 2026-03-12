using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Policy.Infrastructure;

public class BeneficiaryRepository : IBeneficiaryRepository
{
    private readonly PolicyDbContext _context;

    public BeneficiaryRepository(PolicyDbContext context)
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
        // Simple mock sequence for now
        var count = await _context.Beneficiaries.IgnoreQueryFilters().CountAsync();
        return $"BEN-{(count + 1):D6}";
    }
}
