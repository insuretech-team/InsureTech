namespace PoliSync.Beneficiaries.Domain;

public interface IBeneficiaryRepository
{
    Task<Beneficiary?> GetByIdAsync(Guid id, CancellationToken ct = default);
    Task<Beneficiary?> GetByCodeAsync(string code, CancellationToken ct = default);
    Task<Beneficiary?> GetByUserIdAsync(Guid userId, CancellationToken ct = default);
    Task AddAsync(Beneficiary beneficiary, CancellationToken ct = default);
    void Update(Beneficiary beneficiary);
    Task<(IEnumerable<Beneficiary> Items, int TotalCount)> ListAsync(
        BeneficiaryType? type, 
        BeneficiaryStatus? status, 
        int page, 
        int pageSize, 
        CancellationToken ct = default);

    // Specific detail additions
    Task AddIndividualDetailsAsync(IndividualBeneficiary details, CancellationToken ct = default);
    Task AddBusinessDetailsAsync(BusinessBeneficiary details, CancellationToken ct = default);
}
