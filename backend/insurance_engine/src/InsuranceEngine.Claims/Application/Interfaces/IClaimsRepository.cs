using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Claims.Domain.Entities;

namespace InsuranceEngine.Claims.Application.Interfaces;

public interface IClaimsRepository
{
    Task<Claim> CreateAsync(Claim claim, CancellationToken cancellationToken = default);
    Task<Claim?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);
    Task<Claim?> GetByClaimNumberAsync(string claimNumber, CancellationToken cancellationToken = default);
    Task UpdateAsync(Claim claim, CancellationToken cancellationToken = default);
    Task<string> GetNextClaimNumberAsync(CancellationToken cancellationToken = default);
    Task<List<Claim>> ListByCustomerAsync(Guid customerId, int page, int pageSize, CancellationToken cancellationToken = default);
    Task<int> GetTotalCountByCustomerAsync(Guid customerId, CancellationToken cancellationToken = default);
}
