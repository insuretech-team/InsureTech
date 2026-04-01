using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Domain.Entities;

namespace InsuranceEngine.Beneficiary.Application.Interfaces;

public interface IBeneficiaryRepository
{
    Task<Domain.Entities.Beneficiary?> GetByIdAsync(Guid id);
    Task<Domain.Entities.Beneficiary?> GetByCodeAsync(string code);
    Task<IEnumerable<Domain.Entities.Beneficiary>> ListAsync(string? type = null, string? status = null, int page = 1, int pageSize = 10);
    Task<int> GetTotalCountAsync(string? type = null, string? status = null);
    Task AddAsync(Domain.Entities.Beneficiary beneficiary);
    Task UpdateAsync(Domain.Entities.Beneficiary beneficiary);
}
