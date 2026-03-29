using Insuretech.Underwriting.Services.V1;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Commands;

public sealed record RequestQuoteCommand(
    string BeneficiaryId,
    string InsurerProductId,
    long SumAssured,
    int TermYears,
    string PremiumPaymentMode,
    List<string>? RiderCodes,
    int ApplicantAge,
    bool Smoker) : IRequest<RequestQuoteResponse>;

public sealed record SubmitHealthDeclarationCommand(
    string QuoteId,
    int HeightCm,
    string WeightKg,
    bool HasPreExistingConditions,
    string? PreExistingConditions,
    bool Smoker,
    bool AlcoholConsumer,
    string? OccupationRiskLevel) : IRequest<SubmitHealthDeclarationResponse>;

public sealed record ApproveUnderwritingCommand(
    string QuoteId,
    string UnderwriterId,
    string? RiskLevel,
    bool PremiumAdjusted,
    long? AdjustedPremium,
    string? Comments) : IRequest<ApproveUnderwritingResponse>;

public sealed record RejectUnderwritingCommand(
    string QuoteId,
    string UnderwriterId,
    string Reason,
    string? RiskLevel,
    string? Comments) : IRequest<RejectUnderwritingResponse>;

public sealed record ConvertQuoteToPolicyCommand(
    string QuoteId,
    string? PaymentMethod,
    string? PaymentReference) : IRequest<ConvertQuoteToPolicyResponse>;
