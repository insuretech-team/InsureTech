using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.ActivateProduct;

public record ActivateProductCommand(Guid ProductId) : IRequest<Result>;
