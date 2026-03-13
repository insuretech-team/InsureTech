using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Application.Interfaces;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Claims.Infrastructure.Repositories;

public class ClaimsRepository : IClaimsRepository
{
    private readonly ClaimsDbContext _context;

    public ClaimsRepository(ClaimsDbContext context)
    {
        _context = context;
    }

    public async Task<Claim> CreateAsync(Claim claim, CancellationToken cancellationToken = default)
    {
        _context.Claims.Add(claim);
        await _context.SaveChangesAsync(cancellationToken);
        return claim;
    }

    public async Task<Claim?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Include(c => c.Approvals)
            .Include(c => c.Documents)
            .FirstOrDefaultAsync(c => c.Id == id, cancellationToken);
    }

    public async Task<Claim?> GetByClaimNumberAsync(string claimNumber, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Include(c => c.Approvals)
            .Include(c => c.Documents)
            .FirstOrDefaultAsync(c => c.ClaimNumber == claimNumber, cancellationToken);
    }

    public async Task UpdateAsync(Claim claim, CancellationToken cancellationToken = default)
    {
        _context.Claims.Update(claim);
        await _context.SaveChangesAsync(cancellationToken);
    }

    public async Task<string> GetNextClaimNumberAsync(CancellationToken cancellationToken = default)
    {
        var year = DateTime.UtcNow.Year;
        var result = await _context.Database
            .SqlQueryRaw<long>("SELECT nextval('insurance_schema.claim_number_seq')")
            .ToListAsync(cancellationToken);
        
        var sequence = result.FirstOrDefault();
        var random = new Random().Next(1000, 9999);
        return $"CLM-{year}-{random:D4}-{sequence:D6}";
    }

    public async Task<List<Claim>> ListByCustomerAsync(Guid customerId, int page, int pageSize, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Where(c => c.CustomerId == customerId && !c.IsDeleted)
            .OrderByDescending(c => c.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }
}
