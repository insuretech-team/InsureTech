using Insuretech.Common.V1;
using Insuretech.Products.Entity.V1;
using ActivateProductGrpcResponse = Insuretech.Products.Services.V1.ActivateProductResponse;
using DeactivateProductGrpcResponse = Insuretech.Products.Services.V1.DeactivateProductResponse;
using ProductPremiumBreakdown = Insuretech.Products.Services.V1.PremiumBreakdown;
using ProtoTimestamp = Google.Protobuf.WellKnownTypes.Timestamp;
using System.Globalization;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using PoliSync.Products.Infrastructure;
using PoliSync.Products.Application.Commands;
using PoliSync.Products.Application.Queries;

namespace PoliSync.ApiHost.Controllers;

/// <summary>
/// HTTP companion for the Product gRPC service (Kestrel port 50121).
/// The InScore gateway (Go) validates JWT, injects X-* identity headers,
/// then reverse-proxies /v1/products/* to this controller.
///
/// BUG-001 FIX: Added /v1/products HTTP routes — previously missing, causing gateway 404.
///
/// Plans, Riders, and Pricing endpoints now use direct repository injection instead of MediatR.
/// </summary>
[ApiController]
public sealed class ProductsController : ControllerBase
{
    private readonly IMediator _mediator;
    private readonly GoProductDataGateway _gateway;
    private readonly ILogger<ProductsController> _logger;

    public ProductsController(
        IMediator mediator,
        GoProductDataGateway gateway,
        ILogger<ProductsController> logger)
    {
        _mediator = mediator;
        _gateway = gateway;
        _logger = logger;
    }

    /// <summary>List all active products. B2C users see only ACTIVE products by default.</summary>
    [HttpGet("/v1/products")]
    public async Task<IActionResult> ListProducts(
        [FromQuery] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        [FromQuery] string? category = null,
        [FromQuery] string? status = null,
        CancellationToken cancellationToken = default)
    {
        ProductCategory? cat = null;
        if (!string.IsNullOrWhiteSpace(category) && System.Enum.TryParse<ProductCategory>(category, ignoreCase: true, out var parsedCat))
            cat = parsedCat;

        // B2C default: only active products
        ProductStatus? st = ProductStatus.Active;
        if (!string.IsNullOrWhiteSpace(status) && System.Enum.TryParse<ProductStatus>(status, ignoreCase: true, out var parsedSt))
            st = parsedSt;

        var query = new ListProductsQuery { Page = page, PageSize = pageSize, Category = cat, Status = st };
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess)
            return StatusCode(500, new { success = false, error = new { message = result.Error } });

        var products = result.Value?.Products ?? [];
        var totalCount = result.Value?.TotalCount ?? 0;

        return Ok(new
        {
            success = true,
            data = new
            {
                products,
                total_count = totalCount,
                page,
                page_size = pageSize
            }
        });
    }

    /// <summary>Get a single product by ID.</summary>
    [HttpGet("/v1/products/{productId}")]
    public async Task<IActionResult> GetProduct(string productId, CancellationToken cancellationToken = default)
    {
        var query = new GetProductQuery { ProductId = productId };
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess || result.Value is null)
            return NotFound(new { success = false, error = new { message = $"Product not found: {productId}" } });

        return Ok(new { success = true, data = result.Value });
    }

    /// <summary>Search products by keyword (in-memory filter over active products).</summary>
    [HttpGet("/v1/products:search")]
    public async Task<IActionResult> SearchProducts(
        [FromQuery] string? q = null,
        [FromQuery] int page = 1,
        [FromQuery(Name = "page_size")] int pageSize = 20,
        CancellationToken cancellationToken = default)
    {
        var query = new ListProductsQuery { Page = 1, PageSize = 200, Status = ProductStatus.Active };
        var result = await _mediator.Send(query, cancellationToken);

        if (!result.IsSuccess)
            return StatusCode(500, new { error = result.Error });

        var products = result.Value?.Products?.ToList() ?? [];
        if (!string.IsNullOrWhiteSpace(q))
        {
            var lower = q.ToLowerInvariant();
            products = products
                .Where(p =>
                    (p.ProductName?.Contains(lower, StringComparison.OrdinalIgnoreCase) ?? false) ||
                    (p.Description?.Contains(lower, StringComparison.OrdinalIgnoreCase) ?? false) ||
                    (p.ProductCode?.Contains(lower, StringComparison.OrdinalIgnoreCase) ?? false))
                .ToList();
        }

        return Ok(new
        {
            success = true,
            data = new { products, total_count = products.Count, page, page_size = pageSize, query = q }
        });
    }

    /// <summary>Create a new product (admin only — AuthZ enforced at gateway).</summary>
    [HttpPost("/v1/products")]
    public async Task<IActionResult> CreateProduct(
        [FromBody] CreateProductRequest request,
        CancellationToken cancellationToken = default)
    {
        if (!System.Enum.TryParse<ProductCategory>(request.Category, ignoreCase: true, out var category))
            return BadRequest(new { success = false, error = new { message = $"Invalid category: {request.Category}" } });

        var command = new CreateProductCommand
        {
            ProductCode = request.ProductCode,
            ProductName = request.Name,
            Category = category,
            Description = request.Description,
            BasePremiumAmount = request.BasePremiumAmount,
            MinSumInsuredAmount = request.MinSumInsuredAmount,
            MaxSumInsuredAmount = request.MaxSumInsuredAmount,
            MinTenureMonths = request.MinTenureMonths,
            MaxTenureMonths = request.MaxTenureMonths,
            Exclusions = request.Exclusions ?? [],
            CreatedBy = request.CreatedBy
                ?? HttpContext.Request.Headers["X-User-ID"].FirstOrDefault()
                ?? "system"
        };

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        var createdProduct = result.Value!;

        return Created($"/v1/products/{createdProduct.ProductId}",
            new { success = true, data = createdProduct });
    }

    /// <summary>Update an existing product (admin only).</summary>
    [HttpPatch("/v1/products/{productId}")]
    public async Task<IActionResult> UpdateProduct(
        string productId,
        [FromBody] UpdateProductRequest request,
        CancellationToken cancellationToken = default)
    {
        var command = new UpdateProductCommand
        {
            ProductId = productId,
            ProductName = request.Name,
            Description = request.Description,
            BasePremiumAmount = request.BasePremiumAmount,
            MinSumInsuredAmount = request.MinSumInsuredAmount,
            MaxSumInsuredAmount = request.MaxSumInsuredAmount,
            MinTenureMonths = request.MinTenureMonths,
            MaxTenureMonths = request.MaxTenureMonths,
            Exclusions = request.Exclusions
        };

        var result = await _mediator.Send(command, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        return Ok(new { success = true, data = result.Value });
    }

    /// <summary>Activate a product.</summary>
    [HttpPost("/v1/products/{productId}/activate")]
    public async Task<IActionResult> ActivateProduct(string productId, CancellationToken cancellationToken = default)
    {
        var result = await _mediator.Send(new UpdateProductCommand
        {
            ProductId = productId,
            Status = ProductStatus.Active
        }, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        return Ok(new
        {
            success = true,
            data = new ActivateProductGrpcResponse { Message = "Product activated" }
        });
    }

    /// <summary>Deactivate a product.</summary>
    [HttpPost("/v1/products/{productId}/deactivate")]
    public async Task<IActionResult> DeactivateProduct(
        string productId,
        CancellationToken cancellationToken = default)
    {
        var result = await _mediator.Send(new UpdateProductCommand
        {
            ProductId = productId,
            Status = ProductStatus.Inactive
        }, cancellationToken);

        if (!result.IsSuccess)
            return BadRequest(new { success = false, error = new { message = result.Error } });

        return Ok(new
        {
            success = true,
            data = new DeactivateProductGrpcResponse { Message = "Product deactivated" }
        });
    }

    /// <summary>Calculate premium.</summary>
    [HttpPost("/v1/premium:calculate")]
    public async Task<IActionResult> CalculatePremium(
        [FromBody] CalculatePremiumHttpRequest request,
        CancellationToken ct = default)
    {
        if (string.IsNullOrWhiteSpace(request.ProductId))
            return BadRequest(new { success = false, error = new { message = "product_id is required" } });

        var product = await _gateway.GetByIdAsync(request.ProductId, ct);
        if (product == null)
            return NotFound(new { success = false, error = new { message = $"Product not found: {request.ProductId}" } });

        if (product.BasePremium == null || product.BasePremium.Amount <= 0)
            return BadRequest(new { success = false, error = new { message = $"Product {request.ProductId} does not have a valid base premium" } });

        if (request.SumInsuredAmount <= 0)
            return BadRequest(new { success = false, error = new { message = "sum_insured_amount must be greater than zero" } });

        if (request.TenureMonths <= 0)
            return BadRequest(new { success = false, error = new { message = "tenure_months must be greater than zero" } });

        var minSumInsured = product.MinSumInsured?.Amount ?? 0;
        var maxSumInsured = product.MaxSumInsured?.Amount ?? long.MaxValue;
        if (request.SumInsuredAmount < minSumInsured || request.SumInsuredAmount > maxSumInsured)
        {
            return BadRequest(new
            {
                success = false,
                error = new
                {
                    message = $"sum_insured_amount must be between {minSumInsured} and {maxSumInsured}"
                }
            });
        }

        if (request.TenureMonths < product.MinTenureMonths || request.TenureMonths > product.MaxTenureMonths)
        {
            return BadRequest(new
            {
                success = false,
                error = new
                {
                    message = $"tenure_months must be between {product.MinTenureMonths} and {product.MaxTenureMonths}"
                }
            });
        }

        var currency = NormalizeCurrency(product.BasePremium.Currency);
        var selectedRiderIds = request.RiderIds?
            .Where(id => !string.IsNullOrWhiteSpace(id))
            .Select(id => id.Trim())
            .Distinct(StringComparer.OrdinalIgnoreCase)
            .ToList() ?? [];

        var selectedRiders = new List<Rider>();
        if (selectedRiderIds.Count > 0)
        {
            var allProductRiders = await _gateway.ListRidersByProductAsync(request.ProductId, ct);
            selectedRiders = allProductRiders
                .Where(rider => selectedRiderIds.Contains(rider.RiderId, StringComparer.OrdinalIgnoreCase))
                .ToList();

            var foundRiderIds = selectedRiders
                .Select(rider => rider.RiderId)
                .ToHashSet(StringComparer.OrdinalIgnoreCase);
            var missingRiderIds = selectedRiderIds
                .Where(id => !foundRiderIds.Contains(id))
                .ToList();

            if (missingRiderIds.Count > 0)
            {
                return BadRequest(new
                {
                    success = false,
                    error = new
                    {
                        message = $"Unknown rider ids for product {request.ProductId}: {string.Join(", ", missingRiderIds)}"
                    }
                });
            }
        }

        long adjustedBasePremium = product.BasePremium.Amount;
        long riderPremium = selectedRiders.Sum(rider => rider.PremiumAmount?.Amount ?? 0);
        var breakdown = new List<ProductPremiumBreakdown>
        {
            new()
            {
                Item = "base_premium",
                Amount = NewMoney(adjustedBasePremium, currency),
                Description = $"Configured base premium for product {product.ProductCode}"
            }
        };

        var now = DateTimeOffset.UtcNow;
        var applicantData = request.ApplicantData ?? new Dictionary<string, string>(StringComparer.OrdinalIgnoreCase);
        var pricingConfigs = await _gateway.ListPricingConfigsByProductAsync(request.ProductId, ct);

        foreach (var config in pricingConfigs.Where(config => IsConfigEffective(config, now)))
        {
            foreach (var rule in config.Rules)
            {
                if (!RuleMatches(rule, applicantData))
                    continue;

                if (rule.Action?.Type == ActionType.Reject)
                {
                    return UnprocessableEntity(new
                    {
                        success = false,
                        error = new
                        {
                            message = $"Application rejected by pricing rule '{rule.RuleName}'"
                        }
                    });
                }

                var delta = CalculateRuleDelta(rule.Action, adjustedBasePremium);
                if (delta == 0)
                    continue;

                adjustedBasePremium += delta;
                breakdown.Add(new ProductPremiumBreakdown
                {
                    Item = rule.RuleName,
                    Amount = NewMoney(delta, currency),
                    Description = $"Applied {rule.Action?.Type} from pricing config {config.PricingConfigId}"
                });
            }
        }

        if (selectedRiders.Count > 0)
        {
            breakdown.Add(new ProductPremiumBreakdown
            {
                Item = "rider_premium",
                Amount = NewMoney(riderPremium, currency),
                Description = $"Premium for {selectedRiders.Count} selected rider(s)"
            });
        }

        var response = new Insuretech.Products.Services.V1.CalculatePremiumResponse
        {
            BasePremium = NewMoney(adjustedBasePremium, currency),
            RiderPremium = NewMoney(riderPremium, currency),
            TotalPremium = NewMoney(adjustedBasePremium + riderPremium, currency)
        };
        response.Breakdown.AddRange(breakdown);

        return Ok(new { success = true, data = response });
    }

    /// <summary>List all plans for a product.</summary>
    [HttpGet("/v1/products/{productId}/plans")]
    public async Task<IActionResult> GetProductPlans(string productId, CancellationToken ct = default)
    {
        try
        {
            _logger.LogInformation("Fetching plans for product {ProductId}", productId);
            var plans = await _gateway.ListPlansByProductAsync(productId, ct);

            if (plans == null || plans.Count == 0)
                return NotFound(new { success = false, error = new { message = $"No plans found for product {productId}" } });

            return Ok(new { success = true, data = new { plans, total = plans.Count } });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching plans for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error fetching plans" } });
        }
    }

    /// <summary>Add plan.</summary>
    [HttpPost("/v1/products/{productId}/plans")]
    public async Task<IActionResult> CreateProductPlan(
        string productId,
        [FromBody] CreateProductPlanHttpRequest request,
        CancellationToken ct = default)
    {
        try
        {
            var plan = new ProductPlan
            {
                PlanId = string.IsNullOrWhiteSpace(request.PlanId) ? Guid.NewGuid().ToString() : request.PlanId,
                ProductId = productId,
                PlanName = request.PlanName,
                PlanDescription = request.PlanDescription ?? string.Empty,
                PremiumAmount = NewMoney(request.PremiumAmount, request.PremiumCurrency),
                MinSumInsured = NewMoney(request.MinSumInsuredAmount, request.MinSumInsuredCurrency),
                MaxSumInsured = NewMoney(request.MaxSumInsuredAmount, request.MaxSumInsuredCurrency),
                Attributes = request.Attributes ?? string.Empty
            };

            var createdPlan = await _gateway.CreatePlanAsync(plan, ct);

            return Created($"/v1/products/{productId}/plans/{createdPlan.PlanId}",
                new { success = true, data = createdPlan });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating plan for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error creating product plan" } });
        }
    }

    /// <summary>Get specific plan by ID.</summary>
    [HttpGet("/v1/products/{productId}/plans/{planId}")]
    public async Task<IActionResult> GetProductPlan(string productId, string planId, CancellationToken ct = default)
    {
        try
        {
            _logger.LogInformation("Fetching plan {PlanId} for product {ProductId}", planId, productId);
            var plan = await _gateway.GetPlanAsync(planId, ct);

            if (plan == null || plan.ProductId != productId)
                return NotFound(new { success = false, error = new { message = $"Plan {planId} not found for product {productId}" } });

            return Ok(new { success = true, data = plan });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching plan {PlanId}", planId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error fetching plan" } });
        }
    }

    /// <summary>Get riders for a product.</summary>
    [HttpGet("/v1/products/{productId}/riders")]
    public async Task<IActionResult> GetProductRiders(string productId, CancellationToken ct = default)
    {
        try
        {
            _logger.LogInformation("Fetching riders for product {ProductId}", productId);
            var riders = await _gateway.ListRidersByProductAsync(productId, ct);

            if (riders == null || riders.Count == 0)
                return NotFound(new { success = false, error = new { message = $"No riders found for product {productId}" } });

            return Ok(new { success = true, data = new { riders, total = riders.Count } });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching riders for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error fetching riders" } });
        }
    }

    /// <summary>Add rider.</summary>
    [HttpPost("/v1/products/{productId}/riders")]
    public async Task<IActionResult> AddRider(
        string productId,
        [FromBody] CreateRiderHttpRequest request,
        CancellationToken ct = default)
    {
        try
        {
            var rider = new Rider
            {
                RiderId = string.IsNullOrWhiteSpace(request.RiderId) ? Guid.NewGuid().ToString() : request.RiderId,
                ProductId = productId,
                RiderName = request.RiderName,
                Description = request.Description ?? string.Empty,
                PremiumAmount = NewMoney(request.PremiumAmount, request.PremiumCurrency),
                CoverageAmount = NewMoney(request.CoverageAmount, request.CoverageCurrency),
                IsMandatory = request.IsMandatory,
                PremiumCurrency = NormalizeCurrency(request.PremiumCurrency),
                CoverageCurrency = NormalizeCurrency(request.CoverageCurrency)
            };

            var createdRider = await _gateway.CreateRiderAsync(rider, ct);

            return Created($"/v1/products/{productId}/riders/{createdRider.RiderId}",
                new { success = true, data = createdRider });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating rider for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error creating rider" } });
        }
    }

    /// <summary>Get pricing config for a product.</summary>
    [HttpGet("/v1/products/{productId}/pricing")]
    public async Task<IActionResult> GetProductPricing(string productId, CancellationToken ct = default)
    {
        try
        {
            _logger.LogInformation("Fetching pricing config for product {ProductId}", productId);
            var pricingConfigs = await _gateway.ListPricingConfigsByProductAsync(productId, ct);

            if (pricingConfigs == null || pricingConfigs.Count == 0)
                return NotFound(new { success = false, error = new { message = $"No pricing config found for product {productId}" } });

            return Ok(new { success = true, data = new { pricing = pricingConfigs, total = pricingConfigs.Count } });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching pricing for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error fetching pricing" } });
        }
    }

    /// <summary>Create pricing config.</summary>
    [HttpPost("/v1/products/{productId}/pricing")]
    public async Task<IActionResult> CreatePricingConfig(
        string productId,
        [FromBody] CreatePricingConfigHttpRequest request,
        CancellationToken ct = default)
    {
        if (request.Rules == null || request.Rules.Count == 0)
            return BadRequest(new { success = false, error = new { message = "At least one pricing rule is required" } });

        if (!TryBuildPricingRules(request.Rules, out var rules, out var error))
            return BadRequest(new { success = false, error = new { message = error } });

        try
        {
            var config = new PricingConfig
            {
                PricingConfigId = string.IsNullOrWhiteSpace(request.PricingConfigId) ? Guid.NewGuid().ToString() : request.PricingConfigId,
                ProductId = productId,
                EffectiveFrom = ProtoTimestamp.FromDateTime(request.EffectiveFrom.UtcDateTime),
                Rules = { rules }
            };

            if (request.EffectiveTo.HasValue)
                config.EffectiveTo = ProtoTimestamp.FromDateTime(request.EffectiveTo.Value.UtcDateTime);

            var createdConfig = await _gateway.CreatePricingConfigAsync(config, ct);

            return Created($"/v1/products/{productId}/pricing/{createdConfig.PricingConfigId}",
                new { success = true, data = createdConfig });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating pricing config for product {ProductId}", productId);
            return StatusCode(500, new { success = false, error = new { message = "Internal server error creating pricing config" } });
        }
    }

    private static Money NewMoney(long amount, string? currency = null) =>
        new()
        {
            Amount = amount,
            Currency = NormalizeCurrency(currency)
        };

    private static string NormalizeCurrency(string? currency) =>
        string.IsNullOrWhiteSpace(currency) ? "BDT" : currency.Trim().ToUpperInvariant();

    private static bool TryBuildPricingRules(
        IEnumerable<PricingRuleHttpRequest> requests,
        out IEnumerable<PricingRule> rules,
        out string? error)
    {
        var builtRules = new List<PricingRule>();

        foreach (var request in requests)
        {
            if (!TryParseProtoEnum(request.Type, RuleType.Unspecified, "RULE_TYPE_", out RuleType ruleType))
            {
                rules = [];
                error = $"Invalid pricing rule type: {request.Type}";
                return false;
            }

            if (!TryParseProtoEnum(request.Action.Type, ActionType.Unspecified, "ACTION_TYPE_", out ActionType actionType))
            {
                rules = [];
                error = $"Invalid pricing action type: {request.Action.Type}";
                return false;
            }

            builtRules.Add(new PricingRule
            {
                RuleId = string.IsNullOrWhiteSpace(request.RuleId) ? Guid.NewGuid().ToString() : request.RuleId,
                RuleName = request.RuleName,
                Type = ruleType,
                Conditions = { request.Conditions.Select(condition => new RuleCondition
                {
                    Field = condition.Field,
                    Operator = condition.Operator,
                    Value = condition.Value
                }) },
                Action = new RuleAction
                {
                    Type = actionType,
                    Value = request.Action.Value
                }
            });
        }

        rules = builtRules;
        error = null;
        return true;
    }

    private static bool IsConfigEffective(PricingConfig config, DateTimeOffset now)
    {
        var effectiveFrom = config.EffectiveFrom?.ToDateTimeOffset() ?? DateTimeOffset.MinValue;
        if (effectiveFrom > now)
            return false;

        var effectiveTo = config.EffectiveTo?.ToDateTimeOffset();
        return !effectiveTo.HasValue || effectiveTo.Value >= now;
    }

    private static bool RuleMatches(PricingRule rule, IReadOnlyDictionary<string, string> applicantData) =>
        rule.Conditions.Count == 0 || rule.Conditions.All(condition => ConditionMatches(condition, applicantData));

    private static bool ConditionMatches(RuleCondition condition, IReadOnlyDictionary<string, string> applicantData)
    {
        if (!applicantData.TryGetValue(condition.Field, out var actualValue))
            return false;

        var op = (condition.Operator ?? string.Empty).Trim().ToUpperInvariant();
        var expectedValue = condition.Value ?? string.Empty;

        return op switch
        {
            "" or "EQ" or "=" => string.Equals(actualValue, expectedValue, StringComparison.OrdinalIgnoreCase),
            "NEQ" or "!=" => !string.Equals(actualValue, expectedValue, StringComparison.OrdinalIgnoreCase),
            "GT" => CompareAsNumbers(actualValue, expectedValue, (left, right) => left > right),
            "GTE" or ">=" => CompareAsNumbers(actualValue, expectedValue, (left, right) => left >= right),
            "LT" => CompareAsNumbers(actualValue, expectedValue, (left, right) => left < right),
            "LTE" or "<=" => CompareAsNumbers(actualValue, expectedValue, (left, right) => left <= right),
            "IN" => expectedValue.Split(',', StringSplitOptions.TrimEntries | StringSplitOptions.RemoveEmptyEntries)
                .Any(value => string.Equals(actualValue, value, StringComparison.OrdinalIgnoreCase)),
            "NOT_IN" => expectedValue.Split(',', StringSplitOptions.TrimEntries | StringSplitOptions.RemoveEmptyEntries)
                .All(value => !string.Equals(actualValue, value, StringComparison.OrdinalIgnoreCase)),
            "BETWEEN" => MatchesBetween(actualValue, expectedValue),
            "CONTAINS" => actualValue.Contains(expectedValue, StringComparison.OrdinalIgnoreCase),
            _ => false
        };
    }

    private static bool MatchesBetween(string actualValue, string expectedValue)
    {
        var bounds = expectedValue.Split(',', StringSplitOptions.TrimEntries | StringSplitOptions.RemoveEmptyEntries);
        return bounds.Length == 2
            && double.TryParse(bounds[0], NumberStyles.Float, CultureInfo.InvariantCulture, out var lower)
            && double.TryParse(bounds[1], NumberStyles.Float, CultureInfo.InvariantCulture, out var upper)
            && double.TryParse(actualValue, NumberStyles.Float, CultureInfo.InvariantCulture, out var actual)
            && actual >= lower
            && actual <= upper;
    }

    private static bool CompareAsNumbers(string actualValue, string expectedValue, Func<double, double, bool> comparer) =>
        double.TryParse(actualValue, NumberStyles.Float, CultureInfo.InvariantCulture, out var actual)
        && double.TryParse(expectedValue, NumberStyles.Float, CultureInfo.InvariantCulture, out var expected)
        && comparer(actual, expected);

    private static long CalculateRuleDelta(RuleAction? action, long currentPremium)
    {
        if (action == null)
            return 0;

        return action.Type switch
        {
            ActionType.IncreasePercentage => (long)Math.Round(currentPremium * (action.Value / 100d), MidpointRounding.AwayFromZero),
            ActionType.DecreasePercentage => -(long)Math.Round(currentPremium * (action.Value / 100d), MidpointRounding.AwayFromZero),
            ActionType.FixedAmount => (long)Math.Round(action.Value, MidpointRounding.AwayFromZero),
            _ => 0
        };
    }

    private static bool TryParseProtoEnum<TEnum>(
        string? input,
        TEnum unspecifiedValue,
        string prefix,
        out TEnum value)
        where TEnum : struct, System.Enum
    {
        value = unspecifiedValue;
        if (string.IsNullOrWhiteSpace(input))
            return false;

        var normalized = input.Trim().Replace('-', '_').Replace(' ', '_').ToUpperInvariant();
        if (System.Enum.TryParse<TEnum>(normalized, ignoreCase: true, out var directMatch))
        {
            value = directMatch;
            return !EqualityComparer<TEnum>.Default.Equals(value, unspecifiedValue);
        }

        var prefixed = normalized.StartsWith(prefix, StringComparison.Ordinal)
            ? normalized
            : prefix + normalized;

        if (!System.Enum.TryParse<TEnum>(prefixed, ignoreCase: true, out var prefixedMatch))
            return false;

        value = prefixedMatch;
        return !EqualityComparer<TEnum>.Default.Equals(value, unspecifiedValue);
    }
}

// ── Request DTOs ──────────────────────────────────────────────────────────────

public sealed record CreateProductRequest(
    string ProductCode,
    string Name,
    string Category,
    string Description,
    long BasePremiumAmount,
    long MinSumInsuredAmount,
    long MaxSumInsuredAmount,
    int MinTenureMonths,
    int MaxTenureMonths,
    List<string>? Exclusions = null,
    string? CreatedBy = null);

public sealed record UpdateProductRequest(
    string? Name = null,
    string? Description = null,
    long? BasePremiumAmount = null,
    long? MinSumInsuredAmount = null,
    long? MaxSumInsuredAmount = null,
    int? MinTenureMonths = null,
    int? MaxTenureMonths = null,
    List<string>? Exclusions = null);

public sealed record CreateProductPlanHttpRequest(
    string PlanName,
    long PremiumAmount,
    long MinSumInsuredAmount,
    long MaxSumInsuredAmount,
    string? PlanId = null,
    string? PlanDescription = null,
    string? PremiumCurrency = null,
    string? MinSumInsuredCurrency = null,
    string? MaxSumInsuredCurrency = null,
    string? Attributes = null);

public sealed record CreateRiderHttpRequest(
    string RiderName,
    long PremiumAmount,
    long CoverageAmount,
    bool IsMandatory = false,
    string? RiderId = null,
    string? Description = null,
    string? PremiumCurrency = null,
    string? CoverageCurrency = null);

public sealed record CreatePricingConfigHttpRequest(
    DateTimeOffset EffectiveFrom,
    List<PricingRuleHttpRequest> Rules,
    string? PricingConfigId = null,
    DateTimeOffset? EffectiveTo = null);

public sealed record PricingRuleHttpRequest(
    string RuleName,
    string Type,
    PricingActionHttpRequest Action,
    List<PricingConditionHttpRequest> Conditions,
    string? RuleId = null);

public sealed record PricingConditionHttpRequest(
    string Field,
    string Operator,
    string Value);

public sealed record PricingActionHttpRequest(
    string Type,
    double Value);

public sealed record CalculatePremiumHttpRequest(
    string ProductId,
    long SumInsuredAmount,
    int TenureMonths,
    List<string>? RiderIds = null,
    Dictionary<string, string>? ApplicantData = null);

