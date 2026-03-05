namespace PoliSync.Products.Domain;

/// <summary>
/// Repository interface for Product aggregate.
/// </summary>
public interface IProductRepository
{
    Task<Product?> GetByIdAsync(Guid productId, CancellationToken ct = default);
    Task<Product?> GetByCodeAsync(string productCode, CancellationToken ct = default);
    Task<List<Product>> ListAsync(ProductCategory? category, int page, int pageSize, CancellationToken ct = default);
    Task<int> CountAsync(ProductCategory? category, CancellationToken ct = default);
    Task<List<Product>> SearchAsync(string? query, ProductCategory? category, long? minPremium, long? maxPremium, CancellationToken ct = default);
    Task AddAsync(Product product, CancellationToken ct = default);
    void Update(Product product);
}
