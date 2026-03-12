using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.Interfaces;

public interface IProductRepository
{
    Task<Product?> GetByIdAsync(Guid id);
    Task<Product?> GetByIdWithRidersAsync(Guid id);
    Task<Product?> GetByCodeAsync(string productCode);
    Task<List<Product>> ListAsync();
    Task<(List<Product> Items, int TotalCount)> ListActiveAsync(ProductCategory? category, int page, int pageSize);
    Task<Guid> AddAsync(Product product);
    Task UpdateAsync(Product product);
    Task DeleteAsync(Guid id);
    Task<List<Product>> SearchAsync(string? query, decimal? minPremium, decimal? maxPremium);
    Task<List<Rider>> GetRidersByIdsAsync(List<Guid> riderIds);
}
