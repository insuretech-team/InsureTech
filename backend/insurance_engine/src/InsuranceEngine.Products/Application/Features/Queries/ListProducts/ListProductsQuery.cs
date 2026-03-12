using System;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.Products.Domain.Enums;

namespace InsuranceEngine.Products.Application.Features.Queries.ListProducts;

public record ListProductsQuery(
    ProductCategory? Category = null,
    int Page = 1,
    int PageSize = 20
) : IRequest<PaginatedResponse<ProductListDto>>;
