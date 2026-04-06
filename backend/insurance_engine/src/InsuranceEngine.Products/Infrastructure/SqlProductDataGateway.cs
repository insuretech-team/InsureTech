using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Entity.V1;
using InsuranceEngine.Products.Domain.Entities;

namespace InsuranceEngine.Products.Infrastructure;

public class SqlProductDataGateway : IProductDataGateway
{
    private readonly ProductsDbContext _context;
    private readonly ILogger<SqlProductDataGateway> _logger;

    public SqlProductDataGateway(ProductsDbContext context, ILogger<SqlProductDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Product?> GetProductAsync(string productId, CancellationToken ct = default)
    {
        var id = Guid.TryParse(productId, out var pid) ? pid : Guid.Empty;
        var product = await _context.Products
            .FirstOrDefaultAsync(p => p.ProductId == id, ct);

        if (product == null)
        {
            return null;
        }

        return MapToProto(product);
    }

    public async Task<IReadOnlyList<Product>> ListProductsAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var pageActual = page > 0 ? page : 1;
        var pageSizeActual = pageSize > 0 ? pageSize : 10;

        var products = await _context.Products
            .Where(p => p.IsActive && p.DeletedAt == null)
            .OrderByDescending(p => p.CreatedAt)
            .Skip((pageActual - 1) * pageSizeActual)
            .Take(pageSizeActual)
            .ToListAsync(ct);

        return products.Select(MapToProto).ToList();
    }

    public async Task<string> CreateProductAsync(Product product, CancellationToken ct = default)
    {
        var productId = Guid.NewGuid();
        var now = DateTime.UtcNow;

        var entity = new ProductEntity
        {
            ProductId = productId,
            TenantId = "default",
            ProductCode = product.ProductCode,
            ProductName = product.ProductName,
            ProductType = "GENERAL",
            Category = product.Category.ToString(),
            Description = product.Description,
            Status = "DRAFT",
            IsActive = false,
            BasePremium = product.BasePremium?.Amount ?? 0,
            BasePremiumCurrency = product.BasePremium?.Currency ?? "BDT",
            MinSumInsured = product.MinSumInsured?.Amount ?? 0,
            MinSumInsuredCurrency = product.MinSumInsured?.Currency ?? "BDT",
            MaxSumInsured = product.MaxSumInsured?.Amount ?? 0,
            MaxSumInsuredCurrency = product.MaxSumInsured?.Currency ?? "BDT",
            UnitAmount = 100000,
            MinAge = 0,
            MaxAge = 0,
            MinTenureMonths = 1,
            MaxTenureMonths = 12,
            CreatedBy = "system",
            CreatedAt = now,
            UpdatedAt = now,
            Version = 1
        };

        _context.Products.Add(entity);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Created product {ProductCode}", product.ProductCode);

        return productId.ToString();
    }

    public async Task UpdateProductAsync(Product product, CancellationToken ct = default)
    {
        var id = Guid.TryParse(product.ProductId, out var pid) ? pid : Guid.Empty;
        var entity = await _context.Products.FindAsync([id], ct);

        if (entity == null)
        {
            _logger.LogWarning("SQL: Product {ProductId} not found for update", product.ProductId);
            return;
        }

        entity.ProductName = product.ProductName;
        entity.Description = product.Description;
        entity.BasePremium = product.BasePremium?.Amount ?? entity.BasePremium;
        entity.MinSumInsured = product.MinSumInsured?.Amount ?? entity.MinSumInsured;
        entity.MaxSumInsured = product.MaxSumInsured?.Amount ?? entity.MaxSumInsured;
        entity.UpdatedAt = DateTime.UtcNow;
        entity.Version++;

        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Updated product {ProductId}", product.ProductId);
    }

    private static Product MapToProto(ProductEntity entity)
    {
        var proto = new Product
        {
            ProductId = entity.ProductId.ToString(),
            ProductCode = entity.ProductCode,
            ProductName = entity.ProductName,
            Category = Enum.TryParse<ProductCategory>(entity.Category, true, out var cat) ? cat : ProductCategory.Unspecified,
            Description = entity.Description ?? "",
            Status = Enum.TryParse<ProductStatus>(entity.Status, true, out var ps) ? ps : ProductStatus.Draft,
            BasePremium = new Insuretech.Common.V1.Money { Amount = entity.BasePremium, Currency = entity.BasePremiumCurrency },
            MinSumInsured = new Insuretech.Common.V1.Money { Amount = entity.MinSumInsured, Currency = entity.MinSumInsuredCurrency },
            MaxSumInsured = new Insuretech.Common.V1.Money { Amount = entity.MaxSumInsured, Currency = entity.MaxSumInsuredCurrency },
            CreatedBy = entity.CreatedBy,
            CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(entity.CreatedAt)
        };

        return proto;
    }
}
