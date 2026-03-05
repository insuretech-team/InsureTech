using Microsoft.EntityFrameworkCore;
using PoliSync.Beneficiaries.Domain;
using PoliSync.Infrastructure.Persistence;

namespace PoliSync.Beneficiaries.Persistence;

public class BeneficiaryRepository : IBeneficiaryRepository
{
    private readonly PoliSyncDbContext _db;

    public BeneficiaryRepository(PoliSyncDbContext db) => _db = db;

    public async Task<Beneficiary?> GetByIdAsync(Guid id, CancellationToken ct = default)
    {
        return await _db.Set<Beneficiary>()
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id, ct);
    }

    public async Task<Beneficiary?> GetByCodeAsync(string code, CancellationToken ct = default)
    {
        return await _db.Set<Beneficiary>()
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.Code == code, ct);
    }

    public async Task<Beneficiary?> GetByUserIdAsync(Guid userId, CancellationToken ct = default)
    {
        return await _db.Set<Beneficiary>()
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.UserId == userId, ct);
    }

    public async Task AddAsync(Beneficiary beneficiary, CancellationToken ct = default)
    {
        await _db.Set<Beneficiary>().AddAsync(beneficiary, ct);
    }

    public void Update(Beneficiary beneficiary)
    {
        _db.Set<Beneficiary>().Update(beneficiary);
    }

    public async Task<(IEnumerable<Beneficiary> Items, int TotalCount)> ListAsync(
        BeneficiaryType? type, 
        BeneficiaryStatus? status, 
        int page, 
        int pageSize, 
        CancellationToken ct = default)
    {
        var query = _db.Set<Beneficiary>().AsQueryable();

        if (type.HasValue && type.Value != BeneficiaryType.Unspecified)
            query = query.Where(b => b.Type == type.Value);

        if (status.HasValue && status.Value != BeneficiaryStatus.Unspecified)
            query = query.Where(b => b.Status == status.Value);

        var total = await query.CountAsync(ct);
        var items = await query
            .OrderByDescending(b => b.BeneficiaryId) // Simplified order
            .Skip((page - 1) * pageSize)
            .Take(pageSize)
            .ToListAsync(ct);

        return (items, total);
    }

    public async Task AddIndividualDetailsAsync(IndividualBeneficiary details, CancellationToken ct = default)
    {
        await _db.Set<IndividualBeneficiary>().AddAsync(details, ct);
    }

    public async Task AddBusinessDetailsAsync(BusinessBeneficiary details, CancellationToken ct = default)
    {
        await _db.Set<BusinessBeneficiary>().AddAsync(details, ct);
    }
}
