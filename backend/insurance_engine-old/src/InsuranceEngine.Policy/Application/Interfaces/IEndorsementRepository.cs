using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;

namespace InsuranceEngine.Policy.Application.Interfaces;

public interface IEndorsementRepository
{
    Task<Endorsement?> GetByIdAsync(Guid id);
    Task<(List<Endorsement> Items, int TotalCount)> ListAsync(Guid? policyId, EndorsementStatus? status, int page, int pageSize);
    Task<Guid> AddAsync(Endorsement endorsement);
    Task UpdateAsync(Endorsement endorsement);
    Task<long> GetNextSequenceNumberAsync();
}
