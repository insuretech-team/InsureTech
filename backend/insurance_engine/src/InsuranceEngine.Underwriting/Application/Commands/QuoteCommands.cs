using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Underwriting.Application.Commands;

public sealed record CreateQuoteCommand(
    string? BeneficiaryId,
    string ProductId,
    decimal SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    int ApplicantAge,
    string? ApplicantOccupation,
    bool Smoker,
    string? SelectedRiders) : ICommand<string>;

public sealed record SubmitQuoteForUnderwritingCommand(string QuoteId) : ICommand<bool>;
public sealed record ApproveQuoteCommand(string QuoteId) : ICommand<bool>;
public sealed record RejectQuoteCommand(string QuoteId, string Reason) : ICommand<bool>;
public sealed record ConvertQuoteToPolicyCommand(string QuoteId) : ICommand<string>;
