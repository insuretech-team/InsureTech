using PoliSync.SharedKernel.Persistence;

namespace PoliSync.Products.Domain;

/// <summary>
/// Repository contract for Product aggregate root.
/// </summary>
public interface IProductRepository : IRepository<Product>
{
    /// <summary>
    /// Get a product by tenant and product code.
    /// </summary>
    Task<Product?> GetByCodeAsync(Guid tenantId, string code, CancellationToken ct = default);

    /// <summary>
    /// List products for a tenant with optional filtering by category.
    /// </summary>
    Task<(List<Product> Items, int Total)> ListAsync(
        Guid tenantId,
        string? category = null,
        int page = 1,
        int pageSize = 20,
        CancellationToken ct = default);

    /// <summary>
    /// Check if a product code exists for a tenant.
    /// </summary>
    Task<bool> ExistsByCodeAsync(Guid tenantId, string code, CancellationToken ct = default);
}
