using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Domain.Entities;

namespace InsuranceEngine.Application.Interfaces;

public interface IProductRepository
{
    Task<Product?> GetByIdAsync(Guid id);
    Task<List<Product>> ListAsync();
    Task<Guid> AddAsync(Product product);
    Task UpdateAsync(Product product);
    Task DeleteAsync(Guid id);
    Task<List<Product>> SearchAsync(string? query, decimal? minPremium, decimal? maxPremium);
}
