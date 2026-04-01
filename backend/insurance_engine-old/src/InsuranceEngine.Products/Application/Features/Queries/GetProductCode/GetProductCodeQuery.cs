using System;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Products.Application.Features.Queries.GetProductCode;

public record GetProductCodeQuery(Guid ProductId) : IRequest<Result<string>>;
