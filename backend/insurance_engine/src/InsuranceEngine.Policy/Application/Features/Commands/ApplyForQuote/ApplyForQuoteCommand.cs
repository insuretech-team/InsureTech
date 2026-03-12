using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.ApplyForQuote;

public record ApplyForQuoteCommand(
    Guid BeneficiaryId,
    Guid ProductId,
    long SumAssuredAmount,
    int TermYears,
    string PremiumPaymentMode, // YEARLY, MONTHLY, etc.
    List<Guid>? SelectedRiderIds,
    int ApplicantAge,
    string? ApplicantOccupation,
    bool IsSmoker,
    UnderwritingHealthDeclarationDto HealthDeclaration
) : IRequest<Result<QuoteDto>>;
