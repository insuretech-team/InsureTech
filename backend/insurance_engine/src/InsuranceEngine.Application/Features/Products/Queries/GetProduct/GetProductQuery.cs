using System;
using MediatR;
using InsuranceEngine.Application.DTOs;

namespace InsuranceEngine.Application.Features.Products.Queries.GetProduct;

public record GetProductQuery(Guid Id) : IRequest<ProductDto?>;
