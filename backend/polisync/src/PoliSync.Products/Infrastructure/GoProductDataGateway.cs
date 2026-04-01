using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using Insuretech.Products.Entity.V1;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Clients;

namespace PoliSync.Products.Infrastructure;

public interface IProductRepository
{
    Task<Product> CreateAsync(Product product, CancellationToken ct = default);
    Task<Product?> GetByIdAsync(string productId, CancellationToken ct = default);
    Task<Product?> GetByCodeAsync(string productCode, CancellationToken ct = default);
    Task<List<Product>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken ct = default);
    Task<List<Product>> GetByCategoryAsync(ProductCategory category, CancellationToken ct = default);
    Task<List<Product>> GetActiveProductsAsync(CancellationToken ct = default);
    Task<Product> UpdateAsync(Product product, CancellationToken ct = default);
    Task DeleteAsync(string productId, CancellationToken ct = default);
}

/// <summary>
/// gRPC gateway: routes all product/plan/rider/pricing calls to the Go insurance service.
/// PoliSync never touches the DB directly. Proto-generated types are the source of truth.
/// </summary>
public sealed class GoProductDataGateway : IProductRepository
{
    private readonly InsuranceServiceClient _client;
    private readonly ILogger<GoProductDataGateway> _logger;

    public GoProductDataGateway(InsuranceServiceClient client, ILogger<GoProductDataGateway> logger)
    {
        _client = client;
        _logger = logger;
    }

    // ===== PRODUCTS =====

    public async Task<Product> CreateAsync(Product product, CancellationToken ct = default)
    {
        var resp = await _client.Client.CreateProductAsync(
            new CreateProductRequest { Product = product },
            _client.BuildCallOptions(ct));
        return resp.Product;
    }

    public async Task<Product?> GetByIdAsync(string productId, CancellationToken ct = default)
    {
        try
        {
            var resp = await _client.Client.GetProductAsync(
                new GetProductRequest { ProductId = productId },
                _client.BuildCallOptions(ct));
            return resp.Product;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            _logger.LogWarning("Product not found: {ProductId}", productId);
            return null;
        }
    }

    public async Task<Product?> GetByCodeAsync(string productCode, CancellationToken ct = default)
    {
        var all = await GetAllAsync(1, 200, ct);
        return all.FirstOrDefault(p => p.ProductCode == productCode);
    }

    public async Task<List<Product>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken ct = default)
    {
        var resp = await _client.Client.ListProductsAsync(
            new ListProductsRequest { Page = page, PageSize = pageSize },
            _client.BuildCallOptions(ct));
        return [.. resp.Products];
    }

    public async Task<List<Product>> GetByCategoryAsync(ProductCategory category, CancellationToken ct = default)
    {
        var all = await GetAllAsync(1, 200, ct);
        return all.Where(p => p.Category == category).ToList();
    }

    public async Task<List<Product>> GetActiveProductsAsync(CancellationToken ct = default)
    {
        var all = await GetAllAsync(1, 200, ct);
        return all.Where(p => p.Status == ProductStatus.Active).ToList();
    }

    public async Task<Product> UpdateAsync(Product product, CancellationToken ct = default)
    {
        var resp = await _client.Client.UpdateProductAsync(
            new UpdateProductRequest { Product = product },
            _client.BuildCallOptions(ct));
        return resp.Product;
    }

    public async Task DeleteAsync(string productId, CancellationToken ct = default)
    {
        await _client.Client.DeleteProductAsync(
            new DeleteProductRequest { ProductId = productId },
            _client.BuildCallOptions(ct));
    }

    // ===== PLANS =====

    public async Task<List<ProductPlan>> ListPlansByProductAsync(string productId, CancellationToken ct = default)
    {
        var resp = await _client.Client.ListProductPlansAsync(
            new ListProductPlansRequest { ProductId = productId },
            _client.BuildCallOptions(ct));
        return [.. resp.Plans];
    }

    public async Task<ProductPlan?> GetPlanAsync(string planId, CancellationToken ct = default)
    {
        try
        {
            var resp = await _client.Client.GetProductPlanAsync(
                new GetProductPlanRequest { PlanId = planId },
                _client.BuildCallOptions(ct));
            return resp.Plan;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            _logger.LogWarning("Plan not found: {PlanId}", planId);
            return null;
        }
    }

    public async Task<ProductPlan> CreatePlanAsync(ProductPlan plan, CancellationToken ct = default)
    {
        var resp = await _client.Client.CreateProductPlanAsync(
            new CreateProductPlanRequest { Plan = plan },
            _client.BuildCallOptions(ct));
        return resp.Plan;
    }

    // ===== RIDERS =====

    public async Task<List<Rider>> ListRidersByProductAsync(string productId, CancellationToken ct = default)
    {
        var resp = await _client.Client.ListRidersAsync(
            new ListRidersRequest { ProductId = productId },
            _client.BuildCallOptions(ct));
        return [.. resp.Riders];
    }

    public async Task<Rider?> GetRiderAsync(string riderId, CancellationToken ct = default)
    {
        try
        {
            var resp = await _client.Client.GetRiderAsync(
                new GetRiderRequest { RiderId = riderId },
                _client.BuildCallOptions(ct));
            return resp.Rider;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            _logger.LogWarning("Rider not found: {RiderId}", riderId);
            return null;
        }
    }

    public async Task<Rider> CreateRiderAsync(Rider rider, CancellationToken ct = default)
    {
        var resp = await _client.Client.CreateRiderAsync(
            new CreateRiderRequest { Rider = rider },
            _client.BuildCallOptions(ct));
        return resp.Rider;
    }

    // ===== PRICING CONFIG =====

    public async Task<List<PricingConfig>> ListPricingConfigsByProductAsync(string productId, CancellationToken ct = default)
    {
        var resp = await _client.Client.ListPricingConfigsAsync(
            new ListPricingConfigsRequest { ProductId = productId },
            _client.BuildCallOptions(ct));
        return [.. resp.Configs];
    }

    public async Task<PricingConfig> CreatePricingConfigAsync(PricingConfig config, CancellationToken ct = default)
    {
        var resp = await _client.Client.CreatePricingConfigAsync(
            new CreatePricingConfigRequest { Config = config },
            _client.BuildCallOptions(ct));
        return resp.Config;
    }
}
