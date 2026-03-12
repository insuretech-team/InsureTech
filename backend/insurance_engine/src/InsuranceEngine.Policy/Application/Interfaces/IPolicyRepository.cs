using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Application.Interfaces;

public interface IPolicyRepository
{
    Task<PolicyEntity?> GetByIdAsync(Guid id);
    Task<PolicyEntity?> GetByIdWithNomineesAsync(Guid id);
    Task<(List<PolicyEntity> Items, int TotalCount)> ListAsync(
        Guid? customerId, PolicyStatus? status, Guid? productId, int page, int pageSize);
    Task<Guid> AddAsync(PolicyEntity policy);
    Task UpdateAsync(PolicyEntity policy);
    Task<long> GetNextSequenceNumberAsync();
    Task<string?> GetProductCodeAsync(Guid productId);
}
