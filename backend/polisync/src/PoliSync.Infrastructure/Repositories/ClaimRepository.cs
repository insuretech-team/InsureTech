using Google.Protobuf.WellKnownTypes;
using Insuretech.Claims.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IClaimRepository
{
    Task<Claim> CreateAsync(Claim claim, CancellationToken cancellationToken = default);
    Task<Claim?> GetByIdAsync(string claimId, CancellationToken cancellationToken = default);
    Task<Claim?> GetByNumberAsync(string claimNumber, CancellationToken cancellationToken = default);
    Task<List<Claim>> GetByPolicyIdAsync(string policyId, CancellationToken cancellationToken = default);
    Task<List<Claim>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default);
    Task<List<Claim>> GetByStatusAsync(ClaimStatus status, CancellationToken cancellationToken = default);
    Task<List<Claim>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Claim> UpdateAsync(Claim claim, CancellationToken cancellationToken = default);
    Task DeleteAsync(string claimId, CancellationToken cancellationToken = default);
}

public class ClaimRepository : IClaimRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<ClaimRepository> _logger;

    public ClaimRepository(PoliSyncDbContext context, ILogger<ClaimRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Claim> CreateAsync(Claim claim, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        claim.CreatedAt = now;
        claim.UpdatedAt = now;
        claim.SubmittedAt = now;

        if (string.IsNullOrEmpty(claim.ClaimId))
        {
            claim.ClaimId = Guid.NewGuid().ToString();
        }

        if (claim.Status == ClaimStatus.Unspecified)
        {
            claim.Status = ClaimStatus.Submitted;
        }

        _context.Claims.Add(claim);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created claim {ClaimId} - {ClaimNumber}", claim.ClaimId, claim.ClaimNumber);
        return claim;
    }

    public async Task<Claim?> GetByIdAsync(string claimId, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .FirstOrDefaultAsync(c => c.ClaimId == claimId, cancellationToken);
    }

    public async Task<Claim?> GetByNumberAsync(string claimNumber, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .FirstOrDefaultAsync(c => c.ClaimNumber == claimNumber, cancellationToken);
    }

    public async Task<List<Claim>> GetByPolicyIdAsync(string policyId, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Where(c => c.PolicyId == policyId)
            .OrderByDescending(c => c.SubmittedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Claim>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Where(c => c.CustomerId == customerId)
            .OrderByDescending(c => c.SubmittedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Claim>> GetByStatusAsync(ClaimStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .Where(c => c.Status == status)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Claim>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Claims
            .OrderByDescending(c => c.SubmittedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Claim> UpdateAsync(Claim claim, CancellationToken cancellationToken = default)
    {
        claim.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _context.Claims.Update(claim);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated claim {ClaimId}", claim.ClaimId);
        return claim;
    }

    public async Task DeleteAsync(string claimId, CancellationToken cancellationToken = default)
    {
        var claim = await GetByIdAsync(claimId, cancellationToken);
        if (claim != null)
        {
            claim.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await UpdateAsync(claim, cancellationToken);
            _logger.LogInformation("Soft deleted claim {ClaimId}", claimId);
        }
    }
}
