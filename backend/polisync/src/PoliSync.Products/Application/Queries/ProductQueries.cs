using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

// ── GetProduct ───────────────────────────────────────────────────────

public record GetProductQuery(Guid ProductId) : IQuery<Product>;

public class GetProductHandler : IQueryHandler<GetProductQuery, Product>
{
    private readonly IProductRepository _repo;

    public GetProductHandler(IProductRepository repo) => _repo = repo;

    public async Task<Result<Product>> Handle(GetProductQuery query, CancellationToken ct)
    {
        var product = await _repo.GetByIdAsync(query.ProductId, ct);
        return product is not null
            ? Result<Product>.Ok(product)
            : Result<Product>.NotFound($"Product '{query.ProductId}' not found");
    }
}

// ── ListProducts ─────────────────────────────────────────────────────

public record ListProductsQuery(
    ProductCategory? Category = null,
    int Page = 1,
    int PageSize = 20
) : IQuery<ListProductsResult>;

public record ListProductsResult(
    List<Product> Products,
    int TotalCount,
    int Page,
    int PageSize
);

public class ListProductsHandler : IQueryHandler<ListProductsQuery, ListProductsResult>
{
    private readonly IProductRepository _repo;

    public ListProductsHandler(IProductRepository repo) => _repo = repo;

    public async Task<Result<ListProductsResult>> Handle(ListProductsQuery query, CancellationToken ct)
    {
        var page = Math.Max(1, query.Page);
        var pageSize = Math.Clamp(query.PageSize, 1, 100);

        var products = await _repo.ListAsync(query.Category, page, pageSize, ct);
        var totalCount = await _repo.CountAsync(query.Category, ct);

        return Result<ListProductsResult>.Ok(new ListProductsResult(products, totalCount, page, pageSize));
    }
}

// ── SearchProducts ───────────────────────────────────────────────────

public record SearchProductsQuery(
    string? Query = null,
    ProductCategory? Category = null,
    long? MinPremium = null,
    long? MaxPremium = null
) : IQuery<SearchProductsResult>;

public record SearchProductsResult(
    List<Product> Products,
    int TotalCount
);

public class SearchProductsHandler : IQueryHandler<SearchProductsQuery, SearchProductsResult>
{
    private readonly IProductRepository _repo;

    public SearchProductsHandler(IProductRepository repo) => _repo = repo;

    public async Task<Result<SearchProductsResult>> Handle(SearchProductsQuery query, CancellationToken ct)
    {
        var products = await _repo.SearchAsync(query.Query, query.Category, query.MinPremium, query.MaxPremium, ct);
        return Result<SearchProductsResult>.Ok(new SearchProductsResult(products, products.Count));
    }
}
