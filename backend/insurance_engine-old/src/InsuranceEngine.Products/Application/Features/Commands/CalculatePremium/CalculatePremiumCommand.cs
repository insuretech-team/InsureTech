using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Features.Commands.CalculatePremium;

public record CalculatePremiumCommand(
    Guid ProductId,
    long SumInsuredAmount,
    int TenureMonths,
    List<Guid>? RiderIds,
    Dictionary<string, string> ApplicantData
) : IRequest<Result<CalculatePremiumResponse>>;
