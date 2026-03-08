using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.ListProducts;

public record ListProductsQuery : IRequest<List<ProductDto>>;
