using FluentValidation;
using Microsoft.Extensions.Logging;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Products.Application.Queries;

public record BreakdownItem(
    string Description,
    long DeltaPaisa
);

public record PricingEvaluationResult(
    long FinalPremiumPaisa,
    long BasePremiumPaisa,
    List<string> AppliedRules,
    List<BreakdownItem> Breakdown
);

public record EvaluatePricingRulesQuery(
    Guid ProductId,
    Dictionary<string, string> InputFactors
) : IQuery<PricingEvaluationResult>;

public sealed class EvaluatePricingRulesQueryHandler : IQueryHandler<EvaluatePricingRulesQuery, PricingEvaluationResult>
{
    private readonly IProductRepository _productRepository;
    private readonly ILogger<EvaluatePricingRulesQueryHandler> _logger;

    public EvaluatePricingRulesQueryHandler(
        IProductRepository productRepository,
        ILogger<EvaluatePricingRulesQueryHandler> logger)
    {
        _productRepository = productRepository;
        _logger = logger;
    }

    public async Task<Result<PricingEvaluationResult>> Handle(EvaluatePricingRulesQuery request, CancellationToken ct)
    {
        // Load product
        var product = await _productRepository.GetByIdAsync(request.ProductId, ct);
        if (product is null)
        {
            _logger.LogWarning("Product not found: {ProductId}", request.ProductId);
            return Result<PricingEvaluationResult>.NotFound($"Product '{request.ProductId}' not found");
        }

        // Check if pricing config exists
        if (product.PricingConfig is null)
        {
            _logger.LogWarning("Pricing config not found for product {ProductId}", request.ProductId);
            return Result<PricingEvaluationResult>.NotFound($"Pricing config not found for product '{request.ProductId}'");
        }

        // Evaluate pricing rules
        var evaluationResult = PricingEngine.Evaluate(
            product.PricingConfig.Rules,
            product.BasePremiumPaisa,
            request.InputFactors
        );

        if (!evaluationResult.IsSuccess)
        {
            _logger.LogWarning("Failed to evaluate pricing rules: {Error}", evaluationResult.Error?.Message);
            return Result<PricingEvaluationResult>.Fail(
                evaluationResult.Error!.Code,
                evaluationResult.Error.Message
            );
        }

        var evaluation = evaluationResult.Value!;

        // Build breakdown items
        var breakdownItems = new List<BreakdownItem>
        {
            new BreakdownItem("Base Premium", product.BasePremiumPaisa)
        };

        foreach (var rule in evaluation.AppliedRules)
        {
            var delta = rule.Action.Value > 0
                ? (long)(product.BasePremiumPaisa * (decimal)rule.Action.Value / 100)
                : 0;

            if (delta != 0)
            {
                breakdownItems.Add(new BreakdownItem($"Rule: {rule.Name}", delta));
            }
        }

        var result = new PricingEvaluationResult(
            FinalPremiumPaisa: evaluation.FinalPremiumPaisa,
            BasePremiumPaisa: product.BasePremiumPaisa,
            AppliedRules: evaluation.AppliedRules.Select(r => r.Name).ToList(),
            Breakdown: breakdownItems
        );

        _logger.LogInformation("Pricing rules evaluated for product {ProductId}: {FinalPremium} paisa",
            request.ProductId, result.FinalPremiumPaisa);

        return Result<PricingEvaluationResult>.Ok(result);
    }
}

public sealed class EvaluatePricingRulesQueryValidator : AbstractValidator<EvaluatePricingRulesQuery>
{
    public EvaluatePricingRulesQueryValidator()
    {
        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.InputFactors)
            .NotNull().WithMessage("Input factors are required");
    }
}
