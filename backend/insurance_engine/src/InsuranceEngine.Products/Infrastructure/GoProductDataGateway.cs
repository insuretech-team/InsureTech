using Grpc.Core;
using Insuretech.Products.Entity.V1;
using InsuranceEngine.Grpc.Clients;
using InsuranceEngine.Products;

namespace InsuranceEngine.Products.Infrastructure;

/// <summary>
/// Implementation of IProductDataGateway using gRPC calls to the Go backend's ProductService.
/// </summary>
public sealed class GoProductDataGateway : IProductDataGateway
{
    private readonly InsuranceServiceClient _client;

    public GoProductDataGateway(InsuranceServiceClient client)
    {
        _client = client;
    }

    public async Task<Product?> GetProductAsync(string productId, CancellationToken ct = default)
    {
        try
        {
            var response = await _client.Products.GetProductAsync(
                new Insuretech.Products.Services.V1.GetProductRequest { ProductId = productId }, 
                _client.BuildCallOptions(ct));
            
            return response.Product;
        }
        catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
        {
            return null;
        }
    }

    public async Task<IReadOnlyList<Product>> ListProductsAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var response = await _client.Products.ListProductsAsync(
            new Insuretech.Products.Services.V1.ListProductsRequest { Page = page, PageSize = pageSize }, 
            _client.BuildCallOptions(ct));
            
        return response.Products.ToList();
    }

    public async Task<string> CreateProductAsync(Product product, CancellationToken ct = default)
    {
        var response = await _client.Products.CreateProductAsync(
            new Insuretech.Products.Services.V1.CreateProductRequest { Product = product }, 
            _client.BuildCallOptions(ct));
            
        return response.ProductId;
    }

    public async Task UpdateProductAsync(Product product, CancellationToken ct = default)
    {
        await _client.Products.UpdateProductAsync(
            new Insuretech.Products.Services.V1.UpdateProductRequest { Product = product }, 
            _client.BuildCallOptions(ct));
    }
}
