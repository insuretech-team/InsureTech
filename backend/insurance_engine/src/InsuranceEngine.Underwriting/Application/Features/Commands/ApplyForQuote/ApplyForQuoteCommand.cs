using System;
using System.Collections.Generic;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.ApplyForQuote;

public record ApplyForQuoteCommand(
    Guid BeneficiaryId,
    Guid ProductId,
    long SumAssuredAmount,
    int TermYears,
    string PremiumPaymentMode,
    List<Guid>? SelectedRiderIds,
    int ApplicantAge,
    string? ApplicantOccupation,
    bool IsSmoker,
    UnderwritingHealthDeclarationDto HealthDeclaration
) : IRequest<Result<QuoteDto>>;
