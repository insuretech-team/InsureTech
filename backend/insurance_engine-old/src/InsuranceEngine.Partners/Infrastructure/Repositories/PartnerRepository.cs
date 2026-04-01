using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Domain.Entities;
using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.Partners.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Partners.Infrastructure.Repositories;

public class PartnerRepository : IPartnerRepository
{
    private readonly PartnerDbContext _context;

    public PartnerRepository(PartnerDbContext context)
    {
        _context = context;
    }

    public async Task<Partner?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return await _context.Partners
            .Include(p => p.Agents)
            .FirstOrDefaultAsync(p => p.Id == id, cancellationToken);
    }

    public async Task<Partner?> GetByCodeAsync(string code, CancellationToken cancellationToken = default)
    {
        return await _context.Partners
            .Include(p => p.Agents)
            .FirstOrDefaultAsync(p => p.Code == code, cancellationToken);
    }

    public async Task<List<Partner>> ListAsync(CancellationToken cancellationToken = default)
    {
        return await _context.Partners
            .Include(p => p.Agents)
            .OrderByDescending(p => p.CreatedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task CreateAsync(Partner partner, CancellationToken cancellationToken = default)
    {
        await _context.Partners.AddAsync(partner, cancellationToken);
        await _context.SaveChangesAsync(cancellationToken);
    }

    public async Task UpdateAsync(Partner partner, CancellationToken cancellationToken = default)
    {
        _context.Partners.Update(partner);
        await _context.SaveChangesAsync(cancellationToken);
    }
}
