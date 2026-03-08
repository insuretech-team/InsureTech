using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Infrastructure.Persistence;

namespace InsuranceEngine.Infrastructure.Repositories;

public class InsurerRepository : IInsurerRepository
{
    private readonly InsuranceDbContext _context;

    public InsurerRepository(InsuranceDbContext context)
    {
        _context = context;
    }

    public async Task<Insurer?> GetByIdAsync(Guid id) => await _context.Insurers.FindAsync(id);
    public async Task<List<Insurer>> ListAsync() => await _context.Insurers.ToListAsync();
    public async Task<List<Product>> GetProductsByInsurerAsync(Guid insurerId)
    {
        return await _context.Products.Where(p => p.InsurerId == insurerId).ToListAsync();
    }
}
