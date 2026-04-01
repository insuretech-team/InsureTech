using System;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.DTOs;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;

public record ListQuotesQuery(
    Guid? BeneficiaryId,
    QuoteStatus? Status,
    int Page = 1,
    int PageSize = 20
) : IRequest<PaginatedResponse<QuoteDto>>;
