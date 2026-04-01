using System.Text.Json;
using Insuretech.Life.Entity.V1;

namespace PoliSync.LifeInsurance.Services;

public class LifeProductService : ILifeProductService
{
    private readonly ILifeProductRepository _productRepository;
    private readonly ILogger<LifeProductService> _logger;

    public LifeProductService(
        ILifeProductRepository productRepository,
        ILogger<LifeProductService> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public Task<LifeProduct?> GetProductAsync(string productId, CancellationToken cancellationToken = default)
    {
        return _productRepository.GetByIdAsync(productId, cancellationToken);
    }

    public Task<IEnumerable<LifeProduct>> ListProductsAsync(LifeProductType? productType, bool onlyActive, CancellationToken cancellationToken = default)
    {
        return _productRepository.GetByFilterAsync(productType, onlyActive, cancellationToken);
    }

    public Task<LifeProduct> CreateProductAsync(LifeProduct product, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Creating life product: {ProductCode}", product.ProductCode);
        return _productRepository.CreateAsync(product, cancellationToken);
    }

    public Task<LifeProduct?> UpdateProductAsync(string productId, LifeProduct product, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Updating life product: {ProductId}", productId);
        product.ProductId = productId;
        return _productRepository.UpdateAsync(product, cancellationToken);
    }

    public Task<bool> DeleteProductAsync(string productId, bool permanent = false, CancellationToken cancellationToken = default)
    {
        _logger.LogInformation("Deleting life product: {ProductId} (permanent: {Permanent})", productId, permanent);
        return _productRepository.DeleteAsync(productId, permanent, cancellationToken);
    }

    public Task<IEnumerable<ConditionMultiplier>> GetHealthConditionsAsync(string productId, CancellationToken cancellationToken = default)
    {
        var product = _productRepository.GetByIdAsync(productId, cancellationToken).Result;
        if (product == null)
        {
            return Task.FromResult(Enumerable.Empty<ConditionMultiplier>());
        }

        var conditions = JsonSerializer.Deserialize<List<ConditionMultiplier>>(product.ConditionMultipliersJson) ?? new List<ConditionMultiplier>();
        return Task.FromResult(conditions.AsEnumerable());
    }
}
