using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Infrastructure.Persistence;

namespace InsuranceEngine.Products.Infrastructure;

public class ProductRepository : IProductRepository
{
    private readonly ProductsDbContext _context;

    public ProductRepository(ProductsDbContext context)
    {
        _context = context;
    }

    public async Task<Product?> GetByIdAsync(Guid id) => await _context.Products.Include(p => p.Plans).Include(p => p.Questions).FirstOrDefaultAsync(p => p.Id == id);
    
    public async Task<List<Product>> ListAsync() => await _context.Products.ToListAsync();
    
    public async Task<Guid> AddAsync(Product product)
    {
        _context.Products.Add(product);
        await _context.SaveChangesAsync();
        return product.Id;
    }
    
    public async Task UpdateAsync(Product product)
    {
        _context.Products.Update(product);
        await _context.SaveChangesAsync();
    }
    
    public async Task DeleteAsync(Guid id)
    {
        var product = await _context.Products.FindAsync(id);
        if (product != null)
        {
            _context.Products.Remove(product);
            await _context.SaveChangesAsync();
        }
    }
    
    public async Task<List<Product>> SearchAsync(string? query, decimal? minPremium, decimal? maxPremium)
    {
        var dbQuery = _context.Products.AsQueryable();
        if (!string.IsNullOrEmpty(query))
            dbQuery = dbQuery.Where(p => p.ProductName.Contains(query) || p.ProductCode.Contains(query));
        
        return await dbQuery.ToListAsync();
    }
    
    
}

