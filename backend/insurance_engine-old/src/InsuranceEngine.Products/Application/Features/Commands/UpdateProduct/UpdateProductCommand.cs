using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.UpdateProduct;

public record UpdateProductCommand(
    Guid Id,
    string ProductName,
    string? Description,
    ProductCategory Category,
    long BasePremiumAmount,
    long MinSumInsuredAmount,
    long MaxSumInsuredAmount,
    int MinTenureMonths,
    int MaxTenureMonths,
    List<string>? Exclusions
) : IRequest<Result>;
