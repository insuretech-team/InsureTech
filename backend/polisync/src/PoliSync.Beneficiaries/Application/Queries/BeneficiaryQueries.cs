using PoliSync.Beneficiaries.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Beneficiaries.Application.Queries;

public record GetBeneficiaryQuery(Guid BeneficiaryId) : IQuery<Beneficiary?>;

public class GetBeneficiaryHandler : IQueryHandler<GetBeneficiaryQuery, Beneficiary?>
{
    private readonly IBeneficiaryRepository _repo;

    public GetBeneficiaryHandler(IBeneficiaryRepository repo) => _repo = repo;

    public async Task<Result<Beneficiary?>> Handle(GetBeneficiaryQuery query, CancellationToken ct)
    {
        var beneficiary = await _repo.GetByIdAsync(query.BeneficiaryId, ct);
        return Result<Beneficiary?>.Ok(beneficiary);
    }
}

public record ListBeneficiariesQuery(
    BeneficiaryType? Type = null,
    BeneficiaryStatus? Status = null,
    int Page = 1,
    int PageSize = 20
) : IQuery<BeneficiaryListResult>;

public record BeneficiaryListResult(IEnumerable<Beneficiary> Items, int TotalCount, int Page, int PageSize);

public class ListBeneficiariesHandler : IQueryHandler<ListBeneficiariesQuery, BeneficiaryListResult>
{
    private readonly IBeneficiaryRepository _repo;

    public ListBeneficiariesHandler(IBeneficiaryRepository repo) => _repo = repo;

    public async Task<Result<BeneficiaryListResult>> Handle(ListBeneficiariesQuery query, CancellationToken ct)
    {
        var (items, total) = await _repo.ListAsync(query.Type, query.Status, query.Page, query.PageSize, ct);
        return Result<BeneficiaryListResult>.Ok(new BeneficiaryListResult(items, total, query.Page, query.PageSize));
    }
}
