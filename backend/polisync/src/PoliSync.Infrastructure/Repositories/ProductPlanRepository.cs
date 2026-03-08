using Google.Protobuf.WellKnownTypes;
using Insuretech.Products.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Infrastructure.Repositories;

public interface IProductPlanRepository
{
    Task<ProductPlan> CreateAsync(ProductPlan entity, CancellationToken cancellationToken = default);
    Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default);
    Task<List<ProductPlan>> GetByStatusAsync(PlanStatus status, CancellationToken cancellationToken = default);    Task<List<ProductPlan>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<ProductPlan> UpdateAsync(ProductPlan entity, CancellationToken cancellationToken = default);
    Task DeleteAsync(string id, CancellationToken cancellationToken = default);
}

public class ProductPlanRepository : IProductPlanRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<ProductPlanRepository> _logger;

    public ProductPlanRepository(PoliSyncDbContext context, ILogger<ProductPlanRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<ProductPlan> CreateAsync(ProductPlan entity, CancellationToken cancellationToken = default)
    {
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        
        if (string.IsNullOrEmpty(entity.PlanId))
        {
            entity.PlanId = Guid.NewGuid().ToString();
        }

        _context.ProductPlans.Add(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created ProductPlan {Id}", entity.PlanId);
        return entity;
    }

    public async Task<> GetByIdAsync(string id, CancellationToken cancellationToken = default)
    {
        return await _context.ProductPlans
            .FirstOrDefaultAsync(e => e.PlanId == id, cancellationToken);
    }

    public async Task<List<ProductPlan>> GetByStatusAsync(PlanStatus status, CancellationToken cancellationToken = default)
    {
        return await _context.ProductPlans
            .Where(e => e.Status == status)
            .ToListAsync(cancellationToken);
    }
    public async Task<List<ProductPlan>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.ProductPlans
            .OrderByDescending(e => e.PlanId)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<ProductPlan> UpdateAsync(ProductPlan entity, CancellationToken cancellationToken = default)
    {
        _context.ProductPlans.Update(entity);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Updated ProductPlan {Id}", entity.PlanId);
        return entity;
    }

    public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
    {
        var entity = await GetByIdAsync(id, cancellationToken);
        if (entity != null)
        {
            _context.ProductPlans.Remove(entity);
            await _context.SaveChangesAsync(cancellationToken);
            _logger.LogInformation("Deleted ProductPlan {Id}", id);
        }
    }
}
