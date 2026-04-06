using Microsoft.Extensions.Logging;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Underwriting.Infrastructure;

public class SqlUnderwritingDataGateway : IUnderwritingDataGateway
{
    private readonly ILogger<SqlUnderwritingDataGateway> _logger;

    public SqlUnderwritingDataGateway(ILogger<SqlUnderwritingDataGateway> logger)
    {
        _logger = logger;
    }

    public Task<RequestQuoteResponse> RequestQuoteAsync(RequestQuoteRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Requesting quote for beneficiary");
        return Task.FromResult(new RequestQuoteResponse { QuoteId = Guid.NewGuid().ToString(), QuoteNumber = $"QT-{Guid.NewGuid().ToString()[..8].ToUpper()}" });
    }

    public Task<GetQuoteResponse> GetQuoteAsync(string quoteId, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Getting quote {QuoteId}", quoteId);
        return Task.FromResult(new GetQuoteResponse());
    }

    public Task<SubmitHealthDeclarationResponse> SubmitHealthDeclarationAsync(SubmitHealthDeclarationRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Submitting health declaration for quote {QuoteId}", request.QuoteId);
        return Task.FromResult(new SubmitHealthDeclarationResponse { Message = "Health declaration submitted" });
    }

    public Task<ApproveUnderwritingResponse> ApproveUnderwritingAsync(ApproveUnderwritingRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Approving underwriting for quote {QuoteId}", request.QuoteId);
        return Task.FromResult(new ApproveUnderwritingResponse { Message = "Underwriting approved" });
    }

    public Task<RejectUnderwritingResponse> RejectUnderwritingAsync(RejectUnderwritingRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Rejecting underwriting for quote {QuoteId}", request.QuoteId);
        return Task.FromResult(new RejectUnderwritingResponse { Message = "Underwriting rejected" });
    }

    public Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicyAsync(ConvertQuoteToPolicyRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Converting quote {QuoteId} to policy", request.QuoteId);
        return Task.FromResult(new ConvertQuoteToPolicyResponse { PolicyId = Guid.NewGuid().ToString(), PolicyNumber = $"POL-{Guid.NewGuid().ToString()[..8].ToUpper()}" });
    }
}
