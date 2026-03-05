using Microsoft.EntityFrameworkCore;
using PoliSync.Infrastructure.Persistence;
using PoliSync.Products.Domain;

namespace PoliSync.Products.Persistence;

/// <summary>
/// EF Core implementation of IProductRepository.
/// </summary>
public class ProductRepository : IProductRepository
{
    private readonly PoliSyncDbContext _db;

    public ProductRepository(PoliSyncDbContext db) => _db = db;

    public async Task<Product?> GetByIdAsync(Guid productId, CancellationToken ct = default)
    {
        return await _db.Set<Product>()
            .Include(p => p.Plans)
            .Include(p => p.AvailableRiders)
            .Include(p => p.PricingConfig)
            .FirstOrDefaultAsync(p => p.ProductId == productId, ct);
    }

    public async Task<Product?> GetByCodeAsync(string productCode, CancellationToken ct = default)
    {
        return await _db.Set<Product>()
            .FirstOrDefaultAsync(p => p.ProductCode == productCode, ct);
    }

    public async Task<List<Product>> ListAsync(ProductCategory? category, int page, int pageSize, CancellationToken ct = default)
    {
        var query = _db.Set<Product>()
            .Include(p => p.Plans)
            .AsQueryable();

        if (category.HasValue && category.Value != ProductCategory.Unspecified)
            query = query.Where(p => p.Category == category.Value);

        return await query
            .OrderByDescending(p => p.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(ct);
    }

    public async Task<int> CountAsync(ProductCategory? category, CancellationToken ct = default)
    {
        var query = _db.Set<Product>().AsQueryable();

        if (category.HasValue && category.Value != ProductCategory.Unspecified)
            query = query.Where(p => p.Category == category.Value);

        return await query.CountAsync(ct);
    }

    public async Task<List<Product>> SearchAsync(string? queryText, ProductCategory? category, long? minPremium, long? maxPremium, CancellationToken ct = default)
    {
        var query = _db.Set<Product>()
            .Include(p => p.Plans)
            .AsQueryable();

        if (!string.IsNullOrWhiteSpace(queryText))
            query = query.Where(p =>
                p.ProductName.Contains(queryText) ||
                (p.Description != null && p.Description.Contains(queryText)));

        if (category.HasValue && category.Value != ProductCategory.Unspecified)
            query = query.Where(p => p.Category == category.Value);

        if (minPremium.HasValue)
            query = query.Where(p => p.BasePremium >= minPremium.Value);

        if (maxPremium.HasValue)
            query = query.Where(p => p.BasePremium <= maxPremium.Value);

        return await query
            .OrderByDescending(p => p.CreatedAt)
            .Take(50) // max search results
            .ToListAsync(ct);
    }

    public async Task AddAsync(Product product, CancellationToken ct = default)
    {
        await _db.Set<Product>().AddAsync(product, ct);
    }

    public void Update(Product product)
    {
        _db.Set<Product>().Update(product);
    }
}
