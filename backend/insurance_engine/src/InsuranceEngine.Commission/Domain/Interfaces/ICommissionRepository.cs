using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Domain.Entities;

namespace InsuranceEngine.Commission.Domain.Interfaces;

public interface ICommissionRepository
{
    Task<Domain.Entities.Commission?> GetByIdAsync(Guid id, CancellationToken cancellationToken = default);
    Task<List<Domain.Entities.Commission>> ListByRecipientAsync(Guid recipientId, CancellationToken cancellationToken = default);
    Task CreateAsync(Domain.Entities.Commission commission, CancellationToken cancellationToken = default);
    Task UpdateAsync(Domain.Entities.Commission commission, CancellationToken cancellationToken = default);
    Task CreatePayoutAsync(Payout payout, CancellationToken cancellationToken = default);
}
