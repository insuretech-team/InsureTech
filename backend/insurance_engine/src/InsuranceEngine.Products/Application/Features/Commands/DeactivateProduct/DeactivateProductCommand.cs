using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.DeactivateProduct;

public record DeactivateProductCommand(Guid ProductId, string? Reason) : IRequest<Result>;
