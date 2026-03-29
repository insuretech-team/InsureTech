using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Products.Services.V1;
using Insuretech.Products.Entity.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Products.Application.Queries;

public sealed class SearchProductsQueryHandler : IRequestHandler<SearchProductsQuery, SearchProductsResponse>
{
    private readonly IRepository<ProductEntity> _repository;
    private readonly ILogger<SearchProductsQueryHandler> _logger;

    public SearchProductsQueryHandler(IRepository<ProductEntity> repository, ILogger<SearchProductsQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<SearchProductsResponse> Handle(SearchProductsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            // Build base predicate
            Expression<Func<ProductEntity, bool>> predicate = p => p.Status == "ACTIVE";

            // Fuzzy search on name or code (FR-022)
            if (!string.IsNullOrWhiteSpace(request.Query))
            {
                var query = request.Query.ToLower();
                predicate = Combine(predicate, p => 
                    p.ProductName.ToLower().Contains(query) || 
                    p.ProductCode.ToLower().Contains(query));
            }

            // Category filter
            if (request.Category != ProductCategory.Unspecified)
            {
                var categoryStr = request.Category.ToString();
                predicate = Combine(predicate, p => p.Category == categoryStr);
            }

            // Premium range filter (FR-022)
            if (request.MinPremium != null)
            {
                var minAmount = request.MinPremium.Amount;
                predicate = Combine(predicate, p => p.BasePremium >= minAmount);
            }

            if (request.MaxPremium != null)
            {
                var maxAmount = request.MaxPremium.Amount;
                predicate = Combine(predicate, p => p.BasePremium <= maxAmount);
            }

            // For search, we typically want all matching results up to a reasonable limit
            // Using GetPagedAsync for consistency even if page 1 is the default
            var (items, totalCount) = await _repository.GetPagedAsync(
                page: 1,
                pageSize: 50, // Top 50 results
                predicate: predicate,
                orderBy: p => p.ProductName,
                cancellationToken: cancellationToken);

            var response = new SearchProductsResponse
            {
                TotalCount = totalCount
            };

            foreach (var entity in items)
            {
                response.Products.Add(MapToProto(entity));
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to search products with query {Query}", request.Query);
            throw;
        }
    }

    private static Expression<Func<T, bool>> Combine<T>(Expression<Func<T, bool>> expr1, Expression<Func<T, bool>> expr2)
    {
        var parameter = Expression.Parameter(typeof(T));
        var leftVisitor = new ReplaceExpressionVisitor(expr1.Parameters[0], parameter);
        var left = leftVisitor.Visit(expr1.Body);
        var rightVisitor = new ReplaceExpressionVisitor(expr2.Parameters[0], parameter);
        var right = rightVisitor.Visit(expr2.Body);
        return Expression.Lambda<Func<T, bool>>(Expression.AndAlso(left, right), parameter);
    }

    private class ReplaceExpressionVisitor : ExpressionVisitor
    {
        private readonly Expression _oldValue;
        private readonly Expression _newValue;
        public ReplaceExpressionVisitor(Expression oldValue, Expression newValue)
        {
            _oldValue = oldValue;
            _newValue = newValue;
        }
        public override Expression Visit(Expression node) => node == _oldValue ? _newValue : base.Visit(node);
    }

    private static Insuretech.Products.Entity.V1.Product MapToProto(ProductEntity entity)
    {
        var product = new Insuretech.Products.Entity.V1.Product
        {
            ProductId = entity.ProductId.ToString(),
            ProductCode = entity.ProductCode,
            ProductName = entity.ProductName,
            Description = entity.Description ?? ""
        };

        if (System.Enum.TryParse<ProductCategory>(entity.Category, true, out var cat)) product.Category = cat;
        if (System.Enum.TryParse<ProductStatus>(entity.Status, true, out var stat)) product.Status = stat;

        product.BasePremium = new Money { Amount = entity.BasePremium, Currency = entity.BasePremiumCurrency };
        product.MinSumInsured = new Money { Amount = entity.MinSumInsured, Currency = entity.MinSumInsuredCurrency };
        product.MaxSumInsured = new Money { Amount = entity.MaxSumInsured, Currency = entity.MaxSumInsuredCurrency };
        product.MinTenureMonths = entity.MinTenureMonths;
        product.MaxTenureMonths = entity.MaxTenureMonths;
        product.MinAge = entity.MinAge;
        product.MaxAge = entity.MaxAge;
        product.TermsUrl = entity.TermsUrl ?? "";

        product.CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.CreatedAt, DateTimeKind.Utc));
        product.UpdatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(entity.UpdatedAt, DateTimeKind.Utc));

        return product;
    }
}
