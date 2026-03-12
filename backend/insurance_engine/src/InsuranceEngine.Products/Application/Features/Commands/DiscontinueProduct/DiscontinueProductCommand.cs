using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.DiscontinueProduct;

public record DiscontinueProductCommand(Guid ProductId, string? Reason) : IRequest<Result>;
