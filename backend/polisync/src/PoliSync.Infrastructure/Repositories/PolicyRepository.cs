using Google.Protobuf.WellKnownTypes;
using Insuretech.Policy.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IPolicyRepository
{
    Task<Policy> CreateAsync(Policy policy, CancellationToken cancellationToken = default);
    Task<Policy?> GetByIdAsync(string policyId, CancellationToken cancellationToken = default);
    Task<Policy?> GetByNumberAsync(string policyNumber, CancellationToken cancellationToken = default);
    Task<List<Policy>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default);
    Task<List<Policy>> GetByProductIdAsync(string productId, CancellationToken cancellationToken = default);
    Task<List<Policy>> GetByStatusAsync(PolicyStatus status, CancellationToken cancellationToken = default);
    Task<List<Policy>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<Policy> UpdateAsync(Policy policy, CancellationToken cancellationToken = default);
    Task DeleteAsync(string policyId, CancellationToken cancellationToken = default);
}

public class PolicyRepository : IPolicyRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<PolicyRepository> _logger;

    public PolicyRepository(PoliSyncDbContext context, ILogger<PolicyRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Policy> CreateAsync(Policy policy, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        policy.CreatedAt = now;
        policy.UpdatedAt = now;

        if (string.IsNullOrEmpty(policy.PolicyId))
        {
            policy.PolicyId = Guid.NewGuid().ToString();
        }

        _context.Policies.Add(policy);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created policy {PolicyId} - {PolicyNumber}", policy.PolicyId, policy.PolicyNumber);
        return policy;
    }

    public async Task<Policy?> GetByIdAsync(string policyId, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .FirstOrDefaultAsync(p => p.PolicyId == policyId, cancellationToken);
    }

    public async Task<Policy?> GetByNumberAsync(string policyNumber, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .FirstOrDefaultAsync(p => p.PolicyNumber == policyNumber, cancellationToken);
    }

    public async Task<List<Policy>> GetByCustomerIdAsync(string customerId, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .Where(p => p.CustomerId == customerId)
            .OrderByDescending(p => p.CreatedAt)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Policy>> GetByProductIdAsync(string productId, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .Where(p => p.ProductId == productId)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Policy>> GetByStatusAsync(PolicyStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .Where(p => p.Status == status)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Policy>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Policies
            .OrderByDescending(p => p.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<Policy> UpdateAsync(Policy policy, CancellationToken cancellationToken = default)
    {
        policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        _context.Policies.Update(policy);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated policy {PolicyId}", policy.PolicyId);
        return policy;
    }

    public async Task DeleteAsync(string policyId, CancellationToken cancellationToken = default)
    {
        var policy = await GetByIdAsync(policyId, cancellationToken);
        if (policy != null)
        {
            policy.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await UpdateAsync(policy, cancellationToken);
            _logger.LogInformation("Soft deleted policy {PolicyId}", policyId);
        }
    }
}
