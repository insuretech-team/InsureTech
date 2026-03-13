using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.UpdateQuoteStatus;

public record UpdateQuoteStatusCommand(Guid QuoteId, QuoteStatus Status) : IRequest<Result>;
