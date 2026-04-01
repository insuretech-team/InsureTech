using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Domain.Entities;

namespace InsuranceEngine.Partners.Domain.Interfaces;

public interface IPartnerRepository
{
    Task<Partner?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);
    Task<Partner?> GetByCodeAsync(string code, CancellationToken cancellationToken = default);
    Task<List<Partner>> ListAsync(CancellationToken cancellationToken = default);
    Task CreateAsync(Partner partner, CancellationToken cancellationToken = default);
    Task UpdateAsync(Partner partner, CancellationToken cancellationToken = default);
}
