using Insuretech.Products.Entity.V1;
using Insuretech.Products.Services.V1;
using PoliSync.Products.Persistence;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public sealed record SearchProductsQuery(SearchProductsRequest Request) : IQuery<SearchProductsResponse>;

public sealed class SearchProductsHandler : IQueryHandler<SearchProductsQuery, SearchProductsResponse>
{
    private readonly ProductRepository _repo;

    public SearchProductsHandler(ProductRepository repo) => _repo = repo;

    public async Task<Result<SearchProductsResponse>> Handle(SearchProductsQuery query, CancellationToken ct)
    {
        var req      = query.Request;
        var category = req.Category != ProductCategory.Unspecified ? req.Category.ToString() : null;

        var items = await _repo.SearchAsync(req.Query, category, ct);

        var response = new SearchProductsResponse { TotalCount = items.Count };
        response.Products.AddRange(items.Select(r => r.ToProto()));
        return Result<SearchProductsResponse>.Ok(response);
    }
}
