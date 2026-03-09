using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;

namespace InsuranceEngine.Products.Application.Features.Queries.ListProducts;

public record ListProductsQuery : IRequest<List<ProductDto>>;
