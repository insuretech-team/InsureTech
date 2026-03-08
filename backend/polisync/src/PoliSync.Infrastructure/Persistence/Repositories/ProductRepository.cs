using Microsoft.EntityFrameworkCore;
using PoliSync.Products.Domain;

namespace PoliSync.Infrastructure.Persistence.Repositories;

/// <summary>
/// Repository implementation for Product aggregate root.
/// Handles all data access operations for products, plans, riders, and pricing configs.
/// </summary>
public class ProductRepository : IProductRepository
{
    private readonly PoliSyncDbContext _dbContext;

    public ProductRepository(PoliSyncDbContext dbContext)
    {
        _dbContext = dbContext;
    }

    /// <summary>
    /// Get a product by ID with all related entities (Plans, Riders, PricingConfig).
    /// </summary>
    public async Task<Product?> GetByIdAsync(Guid id, CancellationToken ct = default)
    {
        return await _dbContext.Products
            .Include(p => p.Plans)
            .Include(p => p.Riders)
            .Include(p => p.PricingConfig)
            .FirstOrDefaultAsync(p => p.Id == id, ct);
    }

    /// <summary>
    /// Get a product by tenant ID and product code.
    /// </summary>
    public async Task<Product?> GetByCodeAsync(Guid tenantId, string code, CancellationToken ct = default)
    {
        return await _dbContext.Products
            .Include(p => p.Plans)
            .Include(p => p.Riders)
            .Include(p => p.PricingConfig)
            .FirstOrDefaultAsync(p => p.TenantId == tenantId && p.ProductCode == code, ct);
    }

    /// <summary>
    /// List products for a tenant with optional category filtering and pagination.
    /// </summary>
    public async Task<(List<Product> Items, int Total)> ListAsync(
        Guid tenantId,
        string? category = null,
        int page = 1,
        int pageSize = 20,
        CancellationToken ct = default)
    {
        var query = _dbContext.Products
            .Where(p => p.TenantId == tenantId);

        if (!string.IsNullOrWhiteSpace(category))
        {
            query = query.Where(p => p.Category == category);
        }

        var total = await query.CountAsync(ct);

        var items = await query
            .OrderBy(p => p.ProductCode)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .Include(p => p.Plans)
            .Include(p => p.Riders)
            .Include(p => p.PricingConfig)
            .ToListAsync(ct);

        return (items, total);
    }

    /// <summary>
    /// Check if a product code exists for a tenant.
    /// </summary>
    public async Task<bool> ExistsByCodeAsync(Guid tenantId, string code, CancellationToken ct = default)
    {
        return await _dbContext.Products
            .AnyAsync(p => p.TenantId == tenantId && p.ProductCode == code, ct);
    }

    /// <summary>
    /// Add a new product to the database.
    /// </summary>
    public async Task AddAsync(Product product, CancellationToken ct = default)
    {
        await _dbContext.Products.AddAsync(product, ct);
    }

    /// <summary>
    /// Update an existing product.
    /// </summary>
    public Task UpdateAsync(Product product, CancellationToken ct = default)
    {
        _dbContext.Products.Update(product);
        return Task.CompletedTask;
    }

    /// <summary>
    /// Delete a product by ID.
    /// </summary>
    public async Task DeleteAsync(Guid id, CancellationToken ct = default)
    {
        var product = await GetByIdAsync(id, ct);
        if (product != null)
        {
            _dbContext.Products.Remove(product);
        }
    }

    /// <summary>
    /// Check if a product exists by ID.
    /// </summary>
    public async Task<bool> ExistsAsync(Guid id, CancellationToken ct = default)
    {
        return await _dbContext.Products.AnyAsync(p => p.Id == id, ct);
    }
}
