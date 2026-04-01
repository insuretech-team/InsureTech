using Google.Protobuf;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Common.V1;
using Insuretech.Quoting.Entity.V1;
using Insuretech.Quoting.Services.V1;
using PoliSync.Infrastructure.Clients;
using PoliSync.Quoting.Services;
using InsuranceGetProductRequest = Insuretech.Insurance.Services.V1.GetProductRequest;

namespace PoliSync.Quoting.GrpcServices;

public class QuotingGrpcService : QuotingService.QuotingServiceBase
{
    private readonly IQuoteService _quoteService;
    private readonly InsuranceServiceClient _insuranceClient;
    private readonly ILogger<QuotingGrpcService> _logger;

    public QuotingGrpcService(
        IQuoteService quoteService,
        InsuranceServiceClient insuranceClient,
        ILogger<QuotingGrpcService> logger)
    {
        _quoteService = quoteService;
        _insuranceClient = insuranceClient;
        _logger = logger;
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
                request.Parameters,
                request.AgentId,
                request.ValidityDays > 0 ? request.ValidityDays : 30,
                context.CancellationToken);

            return new GenerateQuoteResponse
            {
                Quote = quote
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error generating quote");
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
        var quote = await _quoteService.GetQuoteByNumberAsync(
            request.QuoteNumber, 
            context.CancellationToken);
        
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
            request.CustomerId,
            request.ProductId,
            request.Status != QuoteStatus.Unspecified ? request.Status : null,
            context.CancellationToken);

        var quoteList = quotes.ToList();

        return new ListQuotesResponse
        {
            Quotes = { quoteList },
            TotalCount = quoteList.Count
        };
    }

    public override async Task<ReviseQuoteResponse> ReviseQuote(
        ReviseQuoteRequest request,
        ServerCallContext context)
    {
        try
        {
            var newQuote = await _quoteService.ReviseQuoteAsync(
                request.QuoteId,
                request.NewParameters,
                request.RevisionReason,
                request.ValidityDays > 0 ? request.ValidityDays : 30,
                context.CancellationToken);

            var parentQuote = await _quoteService.GetQuoteAsync(request.QuoteId, context.CancellationToken);

            return new ReviseQuoteResponse
            {
                Quote = newQuote,
                ParentQuote = parentQuote
            };
        }
        catch (InvalidOperationException ex)
        {
            return new ReviseQuoteResponse
            {
                Error = new Error
                {
                    Code = "QUOTE_NOT_FOUND",
                    Message = ex.Message,
                    HttpStatusCode = 404
                }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error revising quote");
            return new ReviseQuoteResponse
            {
                Error = new Error
                {
                    Code = "REVISE_ERROR",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }

    public override async Task<CompareQuotesResponse> CompareQuotes(
        CompareQuotesRequest request,
        ServerCallContext context)
    {
        var quotes = await _quoteService.CompareQuotesAsync(
            request.QuoteIds.ToList(),
            context.CancellationToken);

        var comparisons = new List<QuoteComparison>();
        var productNames = await ResolveProductNamesAsync(quotes, context.CancellationToken);
        foreach (var quote in quotes)
        {
            var comparison = new QuoteComparison
            {
                QuoteId = quote.QuoteId,
                QuoteNumber = quote.QuoteNumber,
                ProductName = productNames.GetValueOrDefault(quote.ProductId, BuildFallbackProductLabel(quote.ProductId)),
                TotalPremium = quote.TotalPremium,
                ValidUntil = quote.ValidUntil,
                Status = quote.Status
            };
            comparisons.Add(comparison);
        }

        return new CompareQuotesResponse
        {
            Comparisons = { comparisons }
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

    public override async Task<ExpireQuoteResponse> ExpireQuote(
        ExpireQuoteRequest request,
        ServerCallContext context)
    {
        var result = await _quoteService.ExpireQuoteAsync(
            request.QuoteId,
            context.CancellationToken);

        return new ExpireQuoteResponse { Success = result };
    }

    public override async Task<DeleteQuoteResponse> DeleteQuote(
        DeleteQuoteRequest request,
        ServerCallContext context)
    {
        var result = await _quoteService.DeleteQuoteAsync(
            request.QuoteId,
            request.Permanent,
            context.CancellationToken);

        return new DeleteQuoteResponse { Success = result };
    }

    public override async Task<CalculatePremiumResponse> CalculatePremium(
        CalculatePremiumRequest request,
        ServerCallContext context)
    {
        try
        {
            var (calculation, coverages, discounts) = await _quoteService.CalculatePremiumAsync(
                request.ProductId,
                request.Parameters,
                context.CancellationToken);

            return new CalculatePremiumResponse
            {
                Calculation = calculation,
                Coverages = { coverages },
                Discounts = { discounts }
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error calculating premium");
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

    private async Task<Dictionary<string, string>> ResolveProductNamesAsync(
        IEnumerable<Quote> quotes,
        CancellationToken cancellationToken)
    {
        var productIds = quotes
            .Select(q => q.ProductId)
            .Where(id => !string.IsNullOrWhiteSpace(id))
            .Distinct(StringComparer.Ordinal)
            .ToList();

        var tasks = productIds.Select(async productId =>
        {
            try
            {
                var response = await _insuranceClient.Client.GetProductAsync(
                    new InsuranceGetProductRequest { ProductId = productId },
                    _insuranceClient.BuildCallOptions(cancellationToken));

                return new KeyValuePair<string, string>(
                    productId,
                    string.IsNullOrWhiteSpace(response.Product.ProductName)
                        ? response.Product.ProductCode
                        : response.Product.ProductName);
            }
            catch (RpcException ex) when (ex.StatusCode == StatusCode.NotFound)
            {
                _logger.LogWarning("Product not found while comparing quotes: {ProductId}", productId);
                return new KeyValuePair<string, string>(productId, BuildFallbackProductLabel(productId));
            }
            catch (Exception ex)
            {
                _logger.LogWarning(ex, "Failed to resolve product name while comparing quotes: {ProductId}", productId);
                return new KeyValuePair<string, string>(productId, BuildFallbackProductLabel(productId));
            }
        });

        var resolved = await Task.WhenAll(tasks);
        return resolved.ToDictionary(pair => pair.Key, pair => pair.Value, StringComparer.Ordinal);
    }

    private static string BuildFallbackProductLabel(string? productId)
        => string.IsNullOrWhiteSpace(productId) ? "Unknown product" : $"Unknown product ({productId})";
}
