using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Domain.Entities;

namespace InsuranceEngine.Policy.Application.Interfaces;

public interface IBeneficiaryRepository
{
    Task<Beneficiary?> GetByIdAsync(Guid id);
    Task<Beneficiary?> GetByCodeAsync(string code);
    Task<IEnumerable<Beneficiary>> ListAsync(string? type = null, string? status = null, int page = 1, int pageSize = 10);
    Task<int> GetTotalCountAsync(string? type = null, string? status = null);
    Task AddAsync(Beneficiary beneficiary);
    Task UpdateAsync(Beneficiary beneficiary);
    Task<string> GetNextSequenceAsync();
}
