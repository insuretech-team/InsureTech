using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Life.Entity.V1;
using Insuretech.Life.Services.V1;
using PoliSync.LifeInsurance.Services;

namespace PoliSync.LifeInsurance.GrpcServices;

public class LifeInsuranceGrpcService : LifeInsuranceService.LifeInsuranceServiceBase
{
    private readonly ILifeProductService _productService;
    private readonly ILifeQuoteService _quoteService;
    private readonly ILogger<LifeInsuranceGrpcService> _logger;

    public LifeInsuranceGrpcService(
        ILifeProductService productService,
        ILifeQuoteService quoteService,
        ILogger<LifeInsuranceGrpcService> logger)
    {
        _productService = productService;
        _quoteService = quoteService;
        _logger = logger;
    }

    public override async Task<GetLifeProductResponse> GetLifeProduct(
        GetLifeProductRequest request,
        ServerCallContext context)
    {
        var product = await _productService.GetProductAsync(request.ProductId, context.CancellationToken);
        
        if (product == null)
        {
            return new GetLifeProductResponse
            {
                Error = new Error
                {
                    Code = "PRODUCT_NOT_FOUND",
                    Message = $"Life product {request.ProductId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetLifeProductResponse { Product = product };
    }

    public override async Task<ListLifeProductsResponse> ListLifeProducts(
        ListLifeProductsRequest request,
        ServerCallContext context)
    {
        var products = await _productService.ListProductsAsync(
            request.ProductType != LifeProductType.Unspecified ? request.ProductType : null,
            request.OnlyActive,
            context.CancellationToken);

        var productList = products.ToList();

        return new ListLifeProductsResponse
        {
            Products = { productList },
            TotalCount = productList.Count
        };
    }

    public override async Task<CreateLifeProductResponse> CreateLifeProduct(
        CreateLifeProductRequest request,
        ServerCallContext context)
    {
        try
        {
            var product = new LifeProduct
            {
                ProductCode = request.ProductCode,
                ProductName = request.ProductName,
                ProductType = request.ProductType,
                Description = request.Description,
                BaseRate = request.BaseRate,
                AgeAdditionConfig = request.AgeAdditionConfig,
                ConditionMultipliersJson = request.ConditionMultipliers.ToString(),
                BonusConfigJson = request.Bonuses.ToString(),
                MinSumAssured = request.MinSumAssured,
                MaxSumAssured = request.MaxSumAssured,
                MinEntryAge = request.MinEntryAge,
                MaxEntryAge = request.MaxEntryAge,
                MinPolicyTerm = request.MinPolicyTerm,
                MaxPolicyTerm = request.MaxPolicyTerm,
                Metadata = { request.Metadata },
                IsActive = true
            };

            var created = await _productService.CreateProductAsync(product, context.CancellationToken);

            return new CreateLifeProductResponse { Product = created };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating life product");
            return new CreateLifeProductResponse
            {
                Error = new Error
                {
                    Code = "CREATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<CalculatePremiumResponse> CalculatePremium(
        CalculatePremiumRequest request,
        ServerCallContext context)
    {
        try
        {
            var (basePremium, ageAddition, conditionMultiplier, conditionAddition, bonusDiscount, totalPremium,
                 breakdown, appliedConditions, appliedBonuses) = await _quoteService.CalculatePremiumAsync(
                request.ProductId,
                request.InsuredPerson,
                request.AgeAtEntry,
                request.PolicyTermYears,
                request.SumAssured,
                request.BonusCodes.ToList(),
                context.CancellationToken);

            return new CalculatePremiumResponse
            {
                BasePremium = basePremium,
                AgeAddition = ageAddition,
                ConditionMultiplier = conditionMultiplier,
                ConditionAddition = conditionAddition,
                BonusDiscount = bonusDiscount,
                TotalPremium = totalPremium,
                Breakdown = { breakdown },
                AppliedConditions = { appliedConditions },
                AppliedBonuses = { appliedBonuses }
            };
        }
        catch (InvalidOperationException ex)
        {
            return new CalculatePremiumResponse
            {
                Error = new Error
                {
                    Code = "PRODUCT_NOT_FOUND",
                    Message = ex.Message,
                    HttpStatusCode = 404
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating life premium");
            return new CalculatePremiumResponse
            {
                Error = new Error
                {
                    Code = "CALCULATION_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GenerateQuoteResponse> GenerateQuote(
        GenerateQuoteRequest request,
        ServerCallContext context)
    {
        try
        {
            var quote = await _quoteService.GenerateQuoteAsync(
                request.ProductId,
                request.CustomerId,
                request.InsuredPerson,
                request.AgeAtEntry,
                request.PolicyTermYears,
                request.SumAssured,
                request.BonusCodes.ToList(),
                request.AgentId,
                request.ValidityDays > 0 ? request.ValidityDays : 30,
                context.CancellationToken);

            return new GenerateQuoteResponse { Quote = quote };
        }
        catch (InvalidOperationException ex)
        {
            return new GenerateQuoteResponse
            {
                Error = new Error
                {
                    Code = "NOT_FOUND",
                    Message = ex.Message,
                    HttpStatusCode = 404
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error generating life quote");
            return new GenerateQuoteResponse
            {
                Error = new Error
                {
                    Code = "GENERATE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<GetQuoteResponse> GetQuote(
        GetQuoteRequest request,
        ServerCallContext context)
    {
        var quote = await _quoteService.GetQuoteAsync(request.QuoteId, context.CancellationToken);
        
        if (quote == null)
        {
            return new GetQuoteResponse
            {
                Error = new Error
                {
                    Code = "QUOTE_NOT_FOUND",
                    Message = $"Quote {request.QuoteId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetQuoteResponse { Quote = quote };
    }

    public override async Task<GetQuoteResponse> GetQuoteByNumber(
        GetQuoteByNumberRequest request,
        ServerCallContext context)
    {
        var quote = await _quoteService.GetQuoteByNumberAsync(request.QuoteNumber, context.CancellationToken);
        
        if (quote == null)
        {
            return new GetQuoteResponse
            {
                Error = new Error
                {
                    Code = "QUOTE_NOT_FOUND",
                    Message = $"Quote {request.QuoteNumber} not found",
                    HttpStatusCode = 404
                }
            };
        }

        return new GetQuoteResponse { Quote = quote };
    }

    public override async Task<ListQuotesResponse> ListQuotes(
        ListQuotesRequest request,
        ServerCallContext context)
    {
        var quotes = await _quoteService.ListQuotesAsync(
            string.IsNullOrEmpty(request.CustomerId) ? null : request.CustomerId,
            string.IsNullOrEmpty(request.ProductId) ? null : request.ProductId,
            request.Status != LifeQuoteStatus.Unspecified ? request.Status : null,
            context.CancellationToken);

        var quoteList = quotes.ToList();

        return new ListQuotesResponse
        {
            Quotes = { quoteList },
            TotalCount = quoteList.Count
        };
    }

    public override async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicy(
        ConvertQuoteToPolicyRequest request,
        ServerCallContext context)
    {
        var result = await _quoteService.ConvertQuoteToPolicyAsync(
            request.QuoteId,
            request.PolicyId,
            request.ConvertedBy,
            context.CancellationToken);

        if (!result)
        {
            return new ConvertQuoteToPolicyResponse
            {
                Error = new Error
                {
                    Code = "CONVERT_ERROR",
                    Message = "Failed to convert quote to policy",
                    HttpStatusCode = 400
                }
            };
        }

        return new ConvertQuoteToPolicyResponse
        {
            QuoteId = request.QuoteId,
            PolicyId = request.PolicyId,
            ConvertedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        };
    }

    public override async Task<GetHealthConditionsResponse> GetHealthConditions(
        GetHealthConditionsRequest request,
        ServerCallContext context)
    {
        var conditions = await _productService.GetHealthConditionsAsync(
            request.ProductId,
            context.CancellationToken);

        return new GetHealthConditionsResponse
        {
            Conditions = { conditions }
        };
    }
}
