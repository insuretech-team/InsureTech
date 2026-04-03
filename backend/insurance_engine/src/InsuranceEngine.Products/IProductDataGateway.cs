using Insuretech.Products.Entity.V1;

namespace InsuranceEngine.Products;

/// <summary>
/// Domain-centric data gateway for Product operations.
/// Moved to Application layer to allow handlers to reference it without circular dependency on Infrastructure.
/// </summary>
public interface IProductDataGateway
{
    Task<Product?> GetProductAsync(string productId, CancellationToken ct = default);
    Task<IReadOnlyList<Product>> ListProductsAsync(int page, int pageSize, CancellationToken ct = default);
    Task<string> CreateProductAsync(Product product, CancellationToken ct = default);
    Task UpdateProductAsync(Product product, CancellationToken ct = default);
}
