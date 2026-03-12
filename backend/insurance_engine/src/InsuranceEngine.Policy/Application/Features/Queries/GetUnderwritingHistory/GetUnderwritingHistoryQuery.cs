using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.GetUnderwritingHistory;

public record GetUnderwritingHistoryQuery(Guid QuoteId) : IRequest<Result<List<UnderwritingDecisionDto>>>;
