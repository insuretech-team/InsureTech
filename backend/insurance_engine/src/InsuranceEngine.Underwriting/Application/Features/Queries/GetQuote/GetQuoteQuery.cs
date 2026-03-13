using System;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;

public record GetQuoteQuery(Guid QuoteId) : IRequest<Result<QuoteDto>>;
