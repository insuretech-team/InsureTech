using System;
using System.Collections.Generic;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Queries.GetUnderwritingHistory;

public record GetUnderwritingHistoryQuery(Guid QuoteId) : IRequest<Result<List<UnderwritingDecisionDto>>>;
