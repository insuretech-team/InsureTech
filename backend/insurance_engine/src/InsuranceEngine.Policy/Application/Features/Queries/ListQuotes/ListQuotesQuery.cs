using System;
using System.Collections.Generic;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Domain.Enums;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.ListQuotes;

/// <summary>
/// List quotes with optional filters and pagination.
/// </summary>
public record ListQuotesQuery(
    Guid? BeneficiaryId,
    QuoteStatus? Status,
    int Page = 1,
    int PageSize = 20
) : IRequest<PaginatedResponse<QuoteDto>>;
