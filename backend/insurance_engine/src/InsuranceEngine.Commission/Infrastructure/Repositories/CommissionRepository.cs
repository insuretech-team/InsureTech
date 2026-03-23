using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Domain.Entities;
using InsuranceEngine.Commission.Domain.Interfaces;
using InsuranceEngine.Commission.Infrastructure.Persistence;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Commission.Infrastructure.Repositories;

public class CommissionRepository : ICommissionRepository
{
    private readonly CommissionDbContext _context;

    public CommissionRepository(CommissionDbContext context)
    {
        _context = context;
    }

    public async Task<Domain.Entities.Commission?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return await _context.Commissions.FirstOrDefaultAsync(c => c.Id == id, cancellationToken);
    }

    public async Task<List<Domain.Entities.Commission>> ListByRecipientAsync(Guid recipientId, CancellationToken cancellationToken = default)
    {
        return await _context.Commissions
            .Where(c => c.PartnerId == recipientId || c.AgentId == recipientId)
            .OrderByDescending(c => c.CreatedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task CreateAsync(Domain.Entities.Commission commission, CancellationToken cancellationToken = default)
    {
        await _context.Commissions.AddAsync(commission, cancellationToken);
        await _context.SaveChangesAsync(cancellationToken);
    }

    public async Task UpdateAsync(Domain.Entities.Commission commission, CancellationToken cancellationToken = default)
    {
        _context.Commissions.Update(commission);
        await _context.SaveChangesAsync(cancellationToken);
    }

    public async Task CreatePayoutAsync(Payout payout, CancellationToken cancellationToken = default)
    {
        await _context.Payouts.AddAsync(payout, cancellationToken);
        await _context.SaveChangesAsync(cancellationToken);
    }
}
