using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Underwriting.Application.Queries;

public sealed record ListQuotesQuery(
    string? BeneficiaryId,
    string? Status,
    int Page = 1,
    int PageSize = 20) : IQuery<ListQuotesResult>;

public sealed record ListQuotesResult(
    IReadOnlyList<QuoteDto> Items,
    int TotalCount);

public sealed record QuoteDto(
    string QuoteId,
    string QuoteNumber,
    string? BeneficiaryId,
    string ProductId,
    string Status,
    decimal SumAssured,
    int TermYears,
    decimal BasePremium,
    decimal TotalPremium,
    int ApplicantAge,
    bool Smoker,
    DateTime? ValidUntil,
    DateTime? CreatedAt)
{
    public QuoteDto() : this("", "", null, "", "", 0, 0, 0, 0, 0, false, null, null) { }
}

public sealed record GetQuoteQuery(string QuoteId) : IQuery<QuoteDto?>;
