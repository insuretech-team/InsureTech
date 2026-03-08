using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Underwriting.Entity.V1;
using Insuretech.Underwriting.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.Events;
using PoliSync.Underwriting.Infrastructure;
using PolicyQuotationStatus = Insuretech.Policy.Entity.V1.QuotationStatus;

namespace PoliSync.Underwriting.GrpcServices;

public sealed class UnderwritingGrpcService : UnderwritingService.UnderwritingServiceBase
{
    private const string UnderwritingDecisionMadeTopic = "insuretech.underwriting.decision_made.v1";

    private readonly IMediator _mediator;
    private readonly ILogger<UnderwritingGrpcService> _logger;
    private readonly IUnderwritingDataGateway _dataGateway;
    private readonly IEventBus _eventBus;
    private readonly IUnderwritingRiskScorer _riskScorer;

    public UnderwritingGrpcService(
        IMediator mediator,
        ILogger<UnderwritingGrpcService> logger,
        IUnderwritingDataGateway dataGateway,
        IEventBus eventBus,
        IUnderwritingRiskScorer riskScorer)
    {
        _mediator = mediator;
        _logger = logger;
        _dataGateway = dataGateway;
        _eventBus = eventBus;
        _riskScorer = riskScorer;
    }

    public override async Task<RequestQuoteResponse> RequestQuote(RequestQuoteRequest request, ServerCallContext context)
    {
        if (string.IsNullOrWhiteSpace(request.BeneficiaryId) || string.IsNullOrWhiteSpace(request.InsurerProductId))
        {
            return new RequestQuoteResponse
            {
                Error = BuildError("VALIDATION_ERROR", "BeneficiaryId and InsurerProductId are required")
            };
        }

        var sumAssuredAmount = request.SumAssured?.Amount > 0 ? request.SumAssured.Amount : 500_000;
        var termYears = request.TermYears <= 0 ? 1 : request.TermYears;
        var riderPremiumAmount = request.RiderCodes.Count * 5_000L;
        var baseRate = request.Smoker ? 0.06 : 0.04;
        var basePremiumAmount = (long)Math.Round(sumAssuredAmount * baseRate / termYears, MidpointRounding.AwayFromZero);
        var taxAmount = (long)Math.Round((basePremiumAmount + riderPremiumAmount) * 0.15, MidpointRounding.AwayFromZero);
        var totalPremiumAmount = basePremiumAmount + riderPremiumAmount + taxAmount;

        var quoteId = Guid.NewGuid().ToString("N");
        var validUntil = Timestamp.FromDateTime(DateTime.UtcNow.AddDays(15));
        var quote = new Quote
        {
            Id = quoteId,
            QuoteNumber = $"QTE-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}",
            BeneficiaryId = request.BeneficiaryId,
            InsurerProductId = request.InsurerProductId,
            Status = QuoteStatus.PendingUnderwriting,
            SumAssured = NormalizeMoney(request.SumAssured, sumAssuredAmount),
            TermYears = termYears,
            PremiumPaymentMode = request.PremiumPaymentMode,
            BasePremium = NewMoney(basePremiumAmount),
            RiderPremium = NewMoney(riderPremiumAmount),
            TaxAmount = NewMoney(taxAmount),
            TotalPremium = NewMoney(totalPremiumAmount),
            PremiumCalculation = $"sum_assured={sumAssuredAmount};base_rate={baseRate:F2};term={termYears}",
            SelectedRiders = string.Join(",", request.RiderCodes),
            ApplicantAge = request.ApplicantAge,
            ApplicantOccupation = "UNKNOWN",
            Smoker = request.Smoker,
            ValidUntil = validUntil
        };

        try
        {
            var created = await _dataGateway.CreateQuoteAsync(quote, GetCancellationToken(context));
            _logger.LogInformation("Underwriting quote created via Go insurance service: {QuoteId}", created.Id);

            return new RequestQuoteResponse
            {
                QuoteId = created.Id,
                QuoteNumber = created.QuoteNumber,
                BasePremium = created.BasePremium,
                TotalPremium = created.TotalPremium,
                ValidUntil = (created.ValidUntil?.ToDateTime() ?? DateTime.UtcNow.AddDays(15)).ToString("O"),
                Message = "Quote requested successfully"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to request underwriting quote for beneficiary {BeneficiaryId}", request.BeneficiaryId);
            return new RequestQuoteResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetQuoteResponse> GetQuote(GetQuoteRequest request, ServerCallContext context)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (quote is null)
            {
                return new GetQuoteResponse
                {
                    Error = BuildError("NOT_FOUND", "Quote not found")
                };
            }

            var response = new GetQuoteResponse { Quote = quote };
            var declaration = await _dataGateway.GetHealthDeclarationByQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (declaration is not null)
            {
                response.HealthDeclaration = declaration;
            }

            var decision = await _dataGateway.GetLatestDecisionByQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (decision is not null)
            {
                response.Decision = decision;
            }

            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get quote {QuoteId}", request.QuoteId);
            return new GetQuoteResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ListQuotesResponse> ListQuotes(ListQuotesRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;
        var status = ParseQuoteStatus(request.Status);

        try
        {
            var fetchSize = Math.Max(page * pageSize * 5, 500);
            var query = await _dataGateway.ListQuotesAsync(request.BeneficiaryId, 1, fetchSize, GetCancellationToken(context));
            var filtered = query.AsEnumerable();

            if (!string.IsNullOrWhiteSpace(request.BeneficiaryId))
            {
                filtered = filtered.Where(q => q.BeneficiaryId == request.BeneficiaryId);
            }

            if (status != QuoteStatus.Unspecified)
            {
                filtered = filtered.Where(q => q.Status == status);
            }

            var ordered = filtered.OrderByDescending(q => q.ValidUntil?.Seconds ?? 0).ToList();
            var pageItems = ordered.Skip((page - 1) * pageSize).Take(pageSize).ToList();

            var response = new ListQuotesResponse { TotalCount = ordered.Count };
            response.Quotes.AddRange(pageItems);
            return response;
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to list quotes for beneficiary {BeneficiaryId}", request.BeneficiaryId);
            return new ListQuotesResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<SubmitHealthDeclarationResponse> SubmitHealthDeclaration(SubmitHealthDeclarationRequest request, ServerCallContext context)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (quote is null)
            {
                return new SubmitHealthDeclarationResponse
                {
                    Error = BuildError("NOT_FOUND", "Quote not found")
                };
            }

            var declarationResult = HealthDeclarationAggregate.Create(
                quoteId: request.QuoteId,
                applicantAge: quote.ApplicantAge,
                heightCm: request.HeightCm,
                weightKg: request.WeightKg,
                hasPreExistingConditions: request.HasPreExistingConditions,
                preExistingConditions: request.PreExistingConditions,
                smoker: request.Smoker,
                alcoholConsumer: request.AlcoholConsumer,
                occupationRiskLevel: request.OccupationRiskLevel);

            if (declarationResult.IsFailure)
            {
                return new SubmitHealthDeclarationResponse
                {
                    Error = BuildError(
                        declarationResult.Error?.Code ?? "VALIDATION_ERROR",
                        declarationResult.Error?.Message ?? "Health declaration validation failed")
                };
            }

            var persisted = await _dataGateway.UpsertHealthDeclarationAsync(
                declarationResult.Value!.Declaration,
                GetCancellationToken(context));

            if (quote.Smoker != request.Smoker)
            {
                quote.Smoker = request.Smoker;
                await _dataGateway.UpdateQuoteAsync(quote, GetCancellationToken(context));
            }

            return new SubmitHealthDeclarationResponse
            {
                Message = "Health declaration submitted",
                MedicalExamRequired = persisted.MedicalExamRequired,
                AutoApprovalPossible = !persisted.MedicalExamRequired
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to submit health declaration for quote {QuoteId}", request.QuoteId);
            return new SubmitHealthDeclarationResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetHealthDeclarationResponse> GetHealthDeclaration(GetHealthDeclarationRequest request, ServerCallContext context)
    {
        try
        {
            var declaration = await _dataGateway.GetHealthDeclarationByQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (declaration is null)
            {
                return new GetHealthDeclarationResponse
                {
                    Error = BuildError("NOT_FOUND", "Health declaration not found")
                };
            }

            return new GetHealthDeclarationResponse
            {
                HealthDeclaration = declaration
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get health declaration for quote {QuoteId}", request.QuoteId);
            return new GetHealthDeclarationResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ApproveUnderwritingResponse> ApproveUnderwriting(ApproveUnderwritingRequest request, ServerCallContext context)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (quote is null)
            {
                return new ApproveUnderwritingResponse
                {
                    Error = BuildError("NOT_FOUND", "Quote not found")
                };
            }

            var now = Timestamp.FromDateTime(DateTime.UtcNow);
            var riskAssessment = await BuildRiskAssessmentAsync(quote, request.QuoteId, GetCancellationToken(context));
            var requestedRiskLevel = ParseRiskLevel(request.RiskLevel);
            var riskLevel = requestedRiskLevel == RiskLevel.Unspecified
                ? riskAssessment.RiskLevel
                : requestedRiskLevel;
            var adjustedPremiumInput = request.AdjustedPremium;
            var shouldAdjustPremium = request.PremiumAdjusted && adjustedPremiumInput is not null && adjustedPremiumInput.Amount > 0;
            var adjustedPremium = shouldAdjustPremium && adjustedPremiumInput is not null
                ? NormalizeMoney(adjustedPremiumInput, adjustedPremiumInput.Amount)
                : quote.TotalPremium;
            var conditionsJson = request.Conditions?.Fields.Count > 0 && request.Conditions is not null
                ? request.Conditions.ToString()
                : string.Empty;
            var decisionResult = UnderwritingDecisionAggregate.CreateApproved(
                quoteId: request.QuoteId,
                underwriterId: request.UnderwriterId,
                comments: request.Comments,
                conditionsJson: conditionsJson,
                premiumAdjusted: shouldAdjustPremium,
                adjustedPremium: adjustedPremium,
                riskAssessment: riskAssessment,
                riskLevel: riskLevel);

            if (decisionResult.IsFailure)
            {
                return new ApproveUnderwritingResponse
                {
                    Error = BuildError(
                        decisionResult.Error?.Code ?? "VALIDATION_ERROR",
                        decisionResult.Error?.Message ?? "Underwriting decision validation failed")
                };
            }

            var persistedDecision = await _dataGateway.UpsertUnderwritingDecisionAsync(
                decisionResult.Value!.Decision,
                GetCancellationToken(context));

            quote.Status = QuoteStatus.Approved;
            if (shouldAdjustPremium)
            {
                quote.TotalPremium = adjustedPremium;
                quote.PremiumCalculation = AppendPremiumCalculation(quote.PremiumCalculation, "adjusted_by_underwriter=true");
            }

            await _dataGateway.UpdateQuoteAsync(quote, GetCancellationToken(context));
            await TryApplyApprovedDecisionToLinkedQuotationAsync(
                quote,
                riskLevel,
                shouldAdjustPremium ? adjustedPremium : null,
                riskAssessment.LoadingPercentage,
                GetCancellationToken(context));
            await TryPublishDecisionMadeEventAsync(
                quoteId: quote.Id,
                decisionId: persistedDecision.Id,
                decision: "APPROVED",
                riskLevel: riskLevel.ToString(),
                premiumAdjusted: shouldAdjustPremium,
                quotedAmount: adjustedPremium.Amount,
                currency: adjustedPremium.Currency,
                reason: string.Empty,
                cancellationToken: GetCancellationToken(context));

            return new ApproveUnderwritingResponse
            {
                DecisionId = persistedDecision.Id,
                Message = "Underwriting approved"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to approve underwriting for quote {QuoteId}", request.QuoteId);
            return new ApproveUnderwritingResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<RejectUnderwritingResponse> RejectUnderwriting(RejectUnderwritingRequest request, ServerCallContext context)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (quote is null)
            {
                return new RejectUnderwritingResponse
                {
                    Error = BuildError("NOT_FOUND", "Quote not found")
                };
            }

            var riskAssessment = await BuildRiskAssessmentAsync(quote, request.QuoteId, GetCancellationToken(context));
            var requestedRiskLevel = ParseRiskLevel(request.RiskLevel);
            var riskLevel = requestedRiskLevel == RiskLevel.Unspecified
                ? riskAssessment.RiskLevel
                : requestedRiskLevel;
            var decisionResult = UnderwritingDecisionAggregate.CreateRejected(
                quoteId: request.QuoteId,
                underwriterId: request.UnderwriterId,
                reason: request.Reason,
                comments: request.Comments,
                riskAssessment: riskAssessment,
                riskLevel: riskLevel);

            if (decisionResult.IsFailure)
            {
                return new RejectUnderwritingResponse
                {
                    Error = BuildError(
                        decisionResult.Error?.Code ?? "VALIDATION_ERROR",
                        decisionResult.Error?.Message ?? "Underwriting decision validation failed")
                };
            }

            var persistedDecision = await _dataGateway.UpsertUnderwritingDecisionAsync(
                decisionResult.Value!.Decision,
                GetCancellationToken(context));
            quote.Status = QuoteStatus.Rejected;
            await _dataGateway.UpdateQuoteAsync(quote, GetCancellationToken(context));
            await TryApplyRejectedDecisionToLinkedQuotationAsync(
                quote.Id,
                request.Reason,
                GetCancellationToken(context));
            await TryPublishDecisionMadeEventAsync(
                quoteId: quote.Id,
                decisionId: persistedDecision.Id,
                decision: "REJECTED",
                riskLevel: riskLevel.ToString(),
                premiumAdjusted: false,
                quotedAmount: 0,
                currency: "BDT",
                reason: request.Reason,
                cancellationToken: GetCancellationToken(context));

            return new RejectUnderwritingResponse
            {
                DecisionId = persistedDecision.Id,
                Message = "Underwriting rejected"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to reject underwriting for quote {QuoteId}", request.QuoteId);
            return new RejectUnderwritingResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicy(ConvertQuoteToPolicyRequest request, ServerCallContext context)
    {
        try
        {
            var quote = await _dataGateway.GetQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (quote is null)
            {
                return new ConvertQuoteToPolicyResponse
                {
                    Error = BuildError("NOT_FOUND", "Quote not found")
                };
            }

            if (quote.Status != QuoteStatus.Approved)
            {
                return new ConvertQuoteToPolicyResponse
                {
                    Error = BuildError("INVALID_STATE", "Only approved quotes can be converted")
                };
            }

            var policyId = $"POL-{Guid.NewGuid():N}"[..16];
            quote.Status = QuoteStatus.Converted;
            quote.ConvertedPolicyId = policyId;
            quote.ConvertedAt = Timestamp.FromDateTime(DateTime.UtcNow);
            await _dataGateway.UpdateQuoteAsync(quote, GetCancellationToken(context));

            return new ConvertQuoteToPolicyResponse
            {
                PolicyId = policyId,
                PolicyNumber = $"LP-{DateTime.UtcNow:yyyy}-{Random.Shared.Next(100000, 999999)}",
                Message = "Quote converted to policy"
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to convert quote {QuoteId} to policy", request.QuoteId);
            return new ConvertQuoteToPolicyResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    public override async Task<GetUnderwritingDecisionResponse> GetUnderwritingDecision(GetUnderwritingDecisionRequest request, ServerCallContext context)
    {
        try
        {
            var decision = await _dataGateway.GetLatestDecisionByQuoteAsync(request.QuoteId, GetCancellationToken(context));
            if (decision is null)
            {
                return new GetUnderwritingDecisionResponse
                {
                    Error = BuildError("NOT_FOUND", "Underwriting decision not found")
                };
            }

            return new GetUnderwritingDecisionResponse
            {
                Decision = decision
            };
        }
        catch (RpcException ex)
        {
            _logger.LogError(ex, "Failed to get underwriting decision for quote {QuoteId}", request.QuoteId);
            return new GetUnderwritingDecisionResponse
            {
                Error = BuildError("UPSTREAM_ERROR", ex.Status.Detail)
            };
        }
    }

    private static QuoteStatus ParseQuoteStatus(string value)
    {
        return ParseEnum(value, QuoteStatus.Unspecified);
    }

    private static RiskLevel ParseRiskLevel(string value)
    {
        return ParseEnum(value, RiskLevel.Unspecified);
    }

    private static TEnum ParseEnum<TEnum>(string value, TEnum fallback) where TEnum : struct, System.Enum
    {
        if (string.IsNullOrWhiteSpace(value))
        {
            return fallback;
        }

        var token = value.Trim();
        if (System.Enum.TryParse<TEnum>(token, true, out var direct))
        {
            return direct;
        }

        var parts = token.Split('_', StringSplitOptions.RemoveEmptyEntries);
        for (var i = 0; i < parts.Length; i++)
        {
            var candidate = string.Concat(parts.Skip(i).Select(ToPascalCase));
            if (System.Enum.TryParse<TEnum>(candidate, true, out var parsed))
            {
                return parsed;
            }
        }

        return fallback;
    }

    private static string ToPascalCase(string segment)
    {
        if (string.IsNullOrWhiteSpace(segment)) return string.Empty;
        var lower = segment.ToLowerInvariant();
        return char.ToUpperInvariant(lower[0]) + lower[1..];
    }

    private async Task TryPublishDecisionMadeEventAsync(
        string quoteId,
        string decisionId,
        string decision,
        string riskLevel,
        bool premiumAdjusted,
        long quotedAmount,
        string currency,
        string reason,
        CancellationToken cancellationToken)
    {
        try
        {
            await _eventBus.PublishAsync(
                new UnderwritingDecisionMadeEvent
                {
                    QuoteId = quoteId ?? string.Empty,
                    DecisionId = decisionId ?? string.Empty,
                    Decision = decision,
                    RiskLevel = riskLevel,
                    PremiumAdjusted = premiumAdjusted,
                    QuotedAmount = quotedAmount,
                    Currency = string.IsNullOrWhiteSpace(currency) ? "BDT" : currency,
                    Reason = reason ?? string.Empty
                },
                UnderwritingDecisionMadeTopic,
                cancellationToken);
        }
        catch (Exception ex)
        {
            _logger.LogWarning(ex, "Failed to publish underwriting decision event for quote {QuoteId}", quoteId);
        }
    }

    private async Task TryApplyApprovedDecisionToLinkedQuotationAsync(
        Quote quote,
        RiskLevel riskLevel,
        Money? manualAdjustedPremium,
        decimal recommendedLoadingPercentage,
        CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(quote.Id))
        {
            return;
        }

        var quotation = await _dataGateway.GetQuotationAsync(quote.Id, cancellationToken);
        if (quotation is null)
        {
            return;
        }

        var quoteAmount = quote.TotalPremium?.Amount ?? 0;
        var baseAmount = quotation.EstimatedPremium?.Amount > 0
            ? quotation.EstimatedPremium.Amount
            : quoteAmount;
        var currency = !string.IsNullOrWhiteSpace(quotation.EstimatedPremium?.Currency)
            ? quotation.EstimatedPremium.Currency
            : (quote.TotalPremium?.Currency ?? "BDT");

        var quotedAmount = manualAdjustedPremium is not null && manualAdjustedPremium.Amount > 0
            ? manualAdjustedPremium.Amount
            : ApplyLoadingFactor(baseAmount, riskLevel, recommendedLoadingPercentage);

        quotation.EstimatedPremium ??= NewMoney(baseAmount, currency);
        quotation.QuotedAmount = NewMoney(quotedAmount, currency);
        if (quotation.Status == PolicyQuotationStatus.Unspecified)
        {
            quotation.Status = PolicyQuotationStatus.Received;
        }

        quotation.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        await _dataGateway.UpdateQuotationAsync(quotation, cancellationToken);
    }

    private async Task TryApplyRejectedDecisionToLinkedQuotationAsync(
        string quoteId,
        string reason,
        CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(quoteId))
        {
            return;
        }

        var quotation = await _dataGateway.GetQuotationAsync(quoteId, cancellationToken);
        if (quotation is null)
        {
            return;
        }

        quotation.Status = PolicyQuotationStatus.Rejected;
        if (!string.IsNullOrWhiteSpace(reason))
        {
            quotation.RejectionReason = reason;
        }

        quotation.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        await _dataGateway.UpdateQuotationAsync(quotation, cancellationToken);
    }

    private static long ApplyLoadingFactor(long baseAmount, RiskLevel riskLevel, decimal recommendedLoadingPercentage)
    {
        if (recommendedLoadingPercentage > 0)
        {
            return (long)Math.Round(baseAmount * (1m + recommendedLoadingPercentage / 100m), MidpointRounding.AwayFromZero);
        }

        var factor = riskLevel switch
        {
            RiskLevel.Low => 0.00m,
            RiskLevel.Medium => 0.10m,
            RiskLevel.High => 0.25m,
            RiskLevel.VeryHigh => 0.40m,
            _ => 0.05m
        };

        return (long)Math.Round(baseAmount * (1m + factor), MidpointRounding.AwayFromZero);
    }

    private async Task<UnderwritingRiskAssessment> BuildRiskAssessmentAsync(
        Quote quote,
        string quoteId,
        CancellationToken cancellationToken)
    {
        var declaration = await _dataGateway.GetHealthDeclarationByQuoteAsync(quoteId, cancellationToken);

        return _riskScorer.Evaluate(new UnderwritingRiskProfile(
            ApplicantAge: quote.ApplicantAge,
            HeightCm: declaration?.HeightCm ?? 0,
            WeightKg: declaration?.WeightKg ?? string.Empty,
            Smoker: declaration?.Smoker ?? quote.Smoker,
            PreExistingConditions: declaration?.PreExistingConditions ?? string.Empty,
            FamilyHistory: declaration?.FamilyHistory ?? string.Empty));
    }

    private static string AppendPremiumCalculation(string existing, string token)
    {
        if (string.IsNullOrWhiteSpace(existing))
        {
            return token;
        }

        return $"{existing};{token}";
    }

    private static Money NormalizeMoney(Money? source, long fallbackAmount)
        => new()
        {
            Amount = source?.Amount > 0 ? source.Amount : fallbackAmount,
            Currency = string.IsNullOrWhiteSpace(source?.Currency) ? "BDT" : source.Currency
        };

    private static Money NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;

    private static Error BuildError(string code, string message)
        => new() { Code = code, Message = message };
}
