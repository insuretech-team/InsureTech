using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Interfaces;

public interface IPolicyRepository
{
    Task<PolicyAggregate?> GetByIdAsync(Guid id, CancellationToken ct = default);
    Task<PolicyAggregate?> GetByPolicyNumberAsync(string policyNumber, CancellationToken ct = default);
    Task<string> GetNextPolicyNumberAsync(CancellationToken ct = default);

    // Original methods, updated to use PolicyAggregate
    Task<PolicyAggregate?> GetByIdWithNomineesAsync(Guid id);
    Task<(List<PolicyAggregate> Items, int TotalCount)> ListAsync(
        Guid? customerId, PolicyStatus? status, Guid? productId, int page, int pageSize);
    Task<Guid> AddAsync(PolicyAggregate policy);
    Task UpdateAsync(PolicyAggregate policy);
    Task<long> GetNextSequenceNumberAsync();
    Task<string?> GetProductCodeAsync(Guid productId);
    Task<bool> ExistsByCustomerAndProductAsync(Guid customerId, Guid productId, DateTime sinceDate);
    Task<bool> ExistsByNidAsync(string encryptedNid, Guid? excludePolicyId = null);
}
