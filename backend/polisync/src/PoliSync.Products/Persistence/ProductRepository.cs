using Insuretech.Products.Entity.V1;
using Insuretech.Products.Services.V1;
using Microsoft.EntityFrameworkCore;
using PoliSync.Infrastructure.Persistence;

namespace PoliSync.Products.Persistence;

/// <summary>
/// Repository for ProductRecord EF POCOs.
/// All reads/writes go through here. Proto types are mapped at the gRPC service layer.
/// </summary>
public sealed class ProductRepository
{
    private readonly PoliSyncDbContext _db;

    public ProductRepository(PoliSyncDbContext db) => _db = db;

    public async Task<ProductRecord?> GetByIdAsync(Guid productId, CancellationToken ct = default)
        => await _db.Set<ProductRecord>()
            .Include(p => p.Plans)
            .Include(p => p.Riders)
            .Include(p => p.PricingConfig)
            .FirstOrDefaultAsync(p => p.ProductId == productId, ct);

    public async Task<ProductRecord?> GetByCodeAsync(string productCode, CancellationToken ct = default)
        => await _db.Set<ProductRecord>()
            .FirstOrDefaultAsync(p => p.ProductCode == productCode, ct);

    public async Task<(List<ProductRecord> Items, int Total)> ListAsync(
        string? category, int page, int pageSize, CancellationToken ct = default)
    {
        var query = _db.Set<ProductRecord>().AsQueryable();
        if (!string.IsNullOrEmpty(category))
            query = query.Where(p => p.Category == category);

        var total = await query.CountAsync(ct);
        var items = await query
            .OrderBy(p => p.ProductName)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .Include(p => p.Plans)
            .Include(p => p.Riders)
            .ToListAsync(ct);

        return (items, total);
    }

    public async Task<List<ProductRecord>> SearchAsync(
        string? query, string? category, CancellationToken ct = default)
    {
        var q = _db.Set<ProductRecord>().AsQueryable();
        if (!string.IsNullOrEmpty(query))
            q = q.Where(p => EF.Functions.ILike(p.ProductName, $"%{query}%")
                           || EF.Functions.ILike(p.ProductCode, $"%{query}%"));
        if (!string.IsNullOrEmpty(category))
            q = q.Where(p => p.Category == category);

        return await q.Include(p => p.Riders).Take(50).ToListAsync(ct);
    }

    public async Task AddAsync(ProductRecord record, CancellationToken ct = default)
    {
        record.CreatedAt = DateTime.UtcNow;
        record.UpdatedAt = DateTime.UtcNow;
        await _db.Set<ProductRecord>().AddAsync(record, ct);
    }

    public Task UpdateAsync(ProductRecord record, CancellationToken ct = default)
    {
        record.UpdatedAt = DateTime.UtcNow;
        _db.Set<ProductRecord>().Update(record);
        return Task.CompletedTask;
    }

    public async Task<bool> ExistsByCodeAsync(string code, CancellationToken ct = default)
        => await _db.Set<ProductRecord>().AnyAsync(p => p.ProductCode == code, ct);
}
