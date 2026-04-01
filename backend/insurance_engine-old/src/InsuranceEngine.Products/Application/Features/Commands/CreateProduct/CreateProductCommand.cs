using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.CreateProduct;

public record CreateProductCommand(
    string ProductCode,
    string ProductName,
    string? Description,
    ProductCategory Category,
    long BasePremiumAmount,
    long MinSumInsuredAmount,
    long MaxSumInsuredAmount,
    int MinTenureMonths,
    int MaxTenureMonths,
    List<string>? Exclusions,
    Guid CreatedBy
) : IRequest<Result<Guid>>;
