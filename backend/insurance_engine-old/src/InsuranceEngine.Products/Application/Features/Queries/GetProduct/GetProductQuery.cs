using System;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;

namespace InsuranceEngine.Products.Application.Features.Queries.GetProduct;

public record GetProductQuery(Guid Id) : IRequest<ProductDto?>;
