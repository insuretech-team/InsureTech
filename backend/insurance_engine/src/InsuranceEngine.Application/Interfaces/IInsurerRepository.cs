using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Domain.Entities;

namespace InsuranceEngine.Application.Interfaces;

public interface IInsurerRepository
{
    Task<Insurer?> GetByIdAsync(Guid id);
    Task<List<Insurer>> ListAsync();
    Task<List<Product>> GetProductsByInsurerAsync(Guid insurerId);
}
