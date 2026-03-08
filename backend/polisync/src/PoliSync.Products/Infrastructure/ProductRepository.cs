using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using PoliSync.Infrastructure.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace PoliSync.Products.Infrastructure;

public interface IProductRepository
{
    Task<Product> CreateAsync(Product product, CancellationToken cancellationToken = default);
    Task<Product?> GetByIdAsync(string productId, CancellationToken cancellationToken = default);
    Task<Product?> GetByCodeAsync(string productCode, CancellationToken cancellationToken = default);
    Task<List<Product>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default);
    Task<List<Product>> GetByCategoryAsync(ProductCategory category, CancellationToken cancellationToken = default);
    Task<List<Product>> GetActiveProductsAsync(CancellationToken cancellationToken = default);
    Task<Product> UpdateAsync(Product product, CancellationToken cancellationToken = default);
    Task DeleteAsync(string productId, CancellationToken cancellationToken = default);
}

public class ProductRepository : IProductRepository
{
    private readonly PoliSyncDbContext _context;
    private readonly ILogger<ProductRepository> _logger;

    public ProductRepository(PoliSyncDbContext context, ILogger<ProductRepository> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<Product> CreateAsync(Product product, CancellationToken cancellationToken = default)
    {
        // Set timestamps
        var now = Timestamp.FromDateTime(DateTime.UtcNow);
        product.CreatedAt = now;
        product.UpdatedAt = now;

        // Generate UUID if not set
        if (string.IsNullOrEmpty(product.ProductId))
        {
            product.ProductId = Guid.NewGuid().ToString();
        }

        // Set default status if not set
        if (product.Status == ProductStatus.Unspecified)
        {
            product.Status = ProductStatus.Draft;
        }

        // Set default currency if not set
        if (string.IsNullOrEmpty(product.BasePremiumCurrency))
        {
            product.BasePremiumCurrency = "BDT";
        }
        if (string.IsNullOrEmpty(product.MinSumInsuredCurrency))
        {
            product.MinSumInsuredCurrency = "BDT";
        }
        if (string.IsNullOrEmpty(product.MaxSumInsuredCurrency))
        {
            product.MaxSumInsuredCurrency = "BDT";
        }

        _context.Products.Add(product);
        await _context.SaveChangesAsync(cancellationToken);
        
        _logger.LogInformation("Created product {ProductId} - {ProductName}", 
            product.ProductId, product.ProductName);
        
        return product;
    }

    public async Task<Product?> GetByIdAsync(string productId, CancellationToken cancellationToken = default)
    {
        return await _context.Products
            .FirstOrDefaultAsync(p => p.ProductId == productId, cancellationToken);
    }

    public async Task<Product?> GetByCodeAsync(string productCode, CancellationToken cancellationToken = default)
    {
        return await _context.Products
            .FirstOrDefaultAsync(p => p.ProductCode == productCode, cancellationToken);
    }

    public async Task<List<Product>> GetAllAsync(int page = 1, int pageSize = 50, CancellationToken cancellationToken = default)
    {
        return await _context.Products
            .OrderByDescending(p => p.CreatedAt)
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Product>> GetByCategoryAsync(ProductCategory category, CancellationToken cancellationToken = default)
    {
        return await _context.Products
            .Where(p => p.Category == category)
            .ToListAsync(cancellationToken);
    }

    public async Task<List<Product>> GetActiveProductsAsync(CancellationToken cancellationToken = default)
    {
        return await _context.Products
            .Where(p => p.Status == ProductStatus.Active && p.DeletedAt == null)
            .ToListAsync(cancellationToken);
    }

    public async Task<Product> UpdateAsync(Product product, CancellationToken cancellationToken = default)
    {
        // Update timestamp
        product.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        _context.Products.Update(product);
        await _context.SaveChangesAsync(cancellationToken);

        _logger.LogInformation("Updated product {ProductId}", product.ProductId);
        
        return product;
    }

    public async Task DeleteAsync(string productId, CancellationToken cancellationToken = default)
    {
        // Soft delete - update deleted_at timestamp
        var product = await GetByIdAsync(productId, cancellationToken);
        if (product != null)
        {
            product.DeletedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await UpdateAsync(product, cancellationToken);
            
            _logger.LogInformation("Soft deleted product {ProductId}", productId);
        }
    }
}
