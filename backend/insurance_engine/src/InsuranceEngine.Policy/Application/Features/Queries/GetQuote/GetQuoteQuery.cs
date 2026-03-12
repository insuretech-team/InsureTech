using System;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.GetQuote;

public record GetQuoteQuery(Guid QuoteId) : IRequest<Result<QuoteDto>>;
