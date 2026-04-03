using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Common.V1;
namespace InsuranceEngine.Underwriting.Application.Commands;

// ===== RequestQuote =====
public sealed class RequestQuoteCommandHandler : IRequestHandler<RequestQuoteCommand, RequestQuoteResponse>
{
    private readonly IUnderwritingDataGateway _gateway;
    private readonly ILogger<RequestQuoteCommandHandler> _logger;

    public RequestQuoteCommandHandler(IUnderwritingDataGateway gateway, ILogger<RequestQuoteCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<RequestQuoteResponse> Handle(RequestQuoteCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Requesting quote for beneficiary: {BeneficiaryId}", request.BeneficiaryId);

            var grpcRequest = new RequestQuoteRequest
            {
                BeneficiaryId = request.BeneficiaryId,
                InsurerProductId = request.InsurerProductId,
                SumAssured = new Money { Amount = request.SumAssured, Currency = "BDT" },
                TermYears = request.TermYears,
                PremiumPaymentMode = request.PremiumPaymentMode,
                ApplicantAge = request.ApplicantAge,
                Smoker = request.Smoker
            };

            if (request.RiderCodes != null)
            {
                grpcRequest.RiderCodes.AddRange(request.RiderCodes);
            }

            var response = await _gateway.RequestQuoteAsync(grpcRequest, cancellationToken);

            if (response.Error != null)
            {
                _logger.LogWarning("Quote request failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Quote generated successfully: {QuoteNumber}", response.QuoteNumber);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to request quote via gateway");
            return new RequestQuoteResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== SubmitHealthDeclaration =====
public sealed class SubmitHealthDeclarationCommandHandler : IRequestHandler<SubmitHealthDeclarationCommand, SubmitHealthDeclarationResponse>
{
    private readonly IUnderwritingDataGateway _gateway;
    private readonly ILogger<SubmitHealthDeclarationCommandHandler> _logger;

    public SubmitHealthDeclarationCommandHandler(IUnderwritingDataGateway gateway, ILogger<SubmitHealthDeclarationCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<SubmitHealthDeclarationResponse> Handle(SubmitHealthDeclarationCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Submitting health declaration for quote: {QuoteId}", request.QuoteId);

            var grpcRequest = new SubmitHealthDeclarationRequest
            {
                QuoteId = request.QuoteId,
                HeightCm = request.HeightCm,
                WeightKg = request.WeightKg.ToString(), // Proto weight_kg is string in definition based on view_file output
                HasPreExistingConditions = request.HasPreExistingConditions,
                PreExistingConditions = request.PreExistingConditions ?? "",
                Smoker = request.Smoker,
                AlcoholConsumer = request.AlcoholConsumer,
                OccupationRiskLevel = request.OccupationRiskLevel
            };

            var response = await _gateway.SubmitHealthDeclarationAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Health declaration failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to submit health declaration via gateway");
            return new SubmitHealthDeclarationResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== ApproveUnderwriting =====
public sealed class ApproveUnderwritingCommandHandler : IRequestHandler<ApproveUnderwritingCommand, ApproveUnderwritingResponse>
{
    private readonly IUnderwritingDataGateway _gateway;
    private readonly ILogger<ApproveUnderwritingCommandHandler> _logger;

    public ApproveUnderwritingCommandHandler(IUnderwritingDataGateway gateway, ILogger<ApproveUnderwritingCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<ApproveUnderwritingResponse> Handle(ApproveUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Approving underwriting for quote: {QuoteId}", request.QuoteId);

            var grpcRequest = new ApproveUnderwritingRequest
            {
                QuoteId = request.QuoteId,
                UnderwriterId = request.UnderwriterId,
                RiskLevel = request.RiskLevel,
                PremiumAdjusted = request.PremiumAdjusted,
                Comments = request.Comments ?? ""
            };

            if (request.AdjustedPremium.HasValue)
            {
                grpcRequest.AdjustedPremium = new Money { Amount = request.AdjustedPremium.Value, Currency = "BDT" };
            }

            var response = await _gateway.ApproveUnderwritingAsync(grpcRequest, cancellationToken);

            if (response.Error != null)
            {
                _logger.LogWarning("Underwriting approval failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to approve underwriting via gateway");
            return new ApproveUnderwritingResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== RejectUnderwriting =====
public sealed class RejectUnderwritingCommandHandler : IRequestHandler<RejectUnderwritingCommand, RejectUnderwritingResponse>
{
    private readonly IUnderwritingDataGateway _gateway;
    private readonly ILogger<RejectUnderwritingCommandHandler> _logger;

    public RejectUnderwritingCommandHandler(IUnderwritingDataGateway gateway, ILogger<RejectUnderwritingCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<RejectUnderwritingResponse> Handle(RejectUnderwritingCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Rejecting underwriting for quote: {QuoteId}", request.QuoteId);

            var grpcRequest = new RejectUnderwritingRequest
            {
                QuoteId = request.QuoteId,
                UnderwriterId = request.UnderwriterId,
                Reason = request.Reason,
                RiskLevel = request.RiskLevel,
                Comments = request.Comments ?? ""
            };

            var response = await _gateway.RejectUnderwritingAsync(grpcRequest, cancellationToken);

            if (response.Error != null)
            {
                _logger.LogWarning("Underwriting rejection failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to reject underwriting via gateway");
            return new RejectUnderwritingResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== ConvertQuoteToPolicy =====
public sealed class ConvertQuoteToPolicyCommandHandler : IRequestHandler<ConvertQuoteToPolicyCommand, ConvertQuoteToPolicyResponse>
{
    private readonly IUnderwritingDataGateway _gateway;
    private readonly ILogger<ConvertQuoteToPolicyCommandHandler> _logger;

    public ConvertQuoteToPolicyCommandHandler(IUnderwritingDataGateway gateway, ILogger<ConvertQuoteToPolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<ConvertQuoteToPolicyResponse> Handle(ConvertQuoteToPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Converting quote to policy: {QuoteId}", request.QuoteId);

            var grpcRequest = new ConvertQuoteToPolicyRequest
            {
                QuoteId = request.QuoteId,
                PaymentMethod = "NOT_SPECIFIED", // Defaulting as C# command doesn't provide this yet
                PaymentReference = ""
            };

            var response = await _gateway.ConvertQuoteToPolicyAsync(grpcRequest, cancellationToken);

            if (response.Error != null)
            {
                _logger.LogWarning("Quote conversion failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Quote successfully converted to Policy: {PolicyNumber}", response.PolicyNumber);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to convert quote to policy via gateway");
            return new ConvertQuoteToPolicyResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
