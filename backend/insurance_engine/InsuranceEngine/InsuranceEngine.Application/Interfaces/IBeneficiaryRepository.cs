using InsuranceEngine.Domain.Entities;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.Interfaces;

public interface IBeneficiaryRepository
{
    Task<IReadOnlyList<Beneficiary>> GetBeneficiariesAsync(
        Guid? policyId,
        BeneficiaryStatus? status,
        BeneficiaryType? type,
        CancellationToken cancellationToken = default);

    Task<(IReadOnlyList<Beneficiary> Beneficiaries, int TotalCount)> GetBeneficiariesPageAsync(
        int page,
        int pageSize,
        CancellationToken cancellationToken = default);

    Task<IReadOnlyList<Beneficiary>> GetBeneficiariesByPolicyAsync(
        Guid policyId,
        CancellationToken cancellationToken = default);

    Task<Beneficiary?> GetBeneficiaryByIdAsync(Guid id, CancellationToken cancellationToken = default);

    Task<Beneficiary?> GetBeneficiaryByIdForPolicyAsync(
        Guid policyId,
        Guid id,
        CancellationToken cancellationToken = default);

    Task<Beneficiary?> GetTrackedBeneficiaryByIdAsync(Guid id, CancellationToken cancellationToken = default);

    Task<Beneficiary?> GetTrackedBeneficiaryByIdForPolicyAsync(
        Guid policyId,
        Guid id,
        CancellationToken cancellationToken = default);

    Task AddBeneficiaryAsync(Beneficiary beneficiary, CancellationToken cancellationToken = default);

    Task SaveChangesAsync(CancellationToken cancellationToken = default);
}
