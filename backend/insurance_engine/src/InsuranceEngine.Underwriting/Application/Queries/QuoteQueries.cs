using Insuretech.Underwriting.Services.V1;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Queries;

public sealed record GetQuoteQuery(string QuoteId) : IRequest<GetQuoteResponse>;
public sealed record ListQuotesQuery(string? BeneficiaryId, string? Status, int Page = 1, int PageSize = 20) : IRequest<ListQuotesResponse>;
public sealed record GetHealthDeclarationQuery(string QuoteId) : IRequest<GetHealthDeclarationResponse>;
public sealed record GetUnderwritingDecisionQuery(string QuoteId) : IRequest<GetUnderwritingDecisionResponse>;
