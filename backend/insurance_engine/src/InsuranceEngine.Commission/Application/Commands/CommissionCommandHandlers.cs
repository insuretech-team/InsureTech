using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Commission.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Commission.Application.Commands;

// ===== CalculateCommission =====
public sealed class CalculateCommissionCommandHandler : IRequestHandler<CalculateCommissionCommand, CalculateCommissionResponse>
{
    private readonly ICommissionDataGateway _gateway;
    private readonly ILogger<CalculateCommissionCommandHandler> _logger;

    public CalculateCommissionCommandHandler(
        ICommissionDataGateway gateway,
        ILogger<CalculateCommissionCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CalculateCommissionResponse> Handle(CalculateCommissionCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Calculating commission for policy: {PolicyId}", request.PolicyId);

            var grpcRequest = new CalculateCommissionRequest
            {
                PolicyId = request.PolicyId,
                RecipientId = request.RecipientId,
                RecipientType = request.RecipientType,
                CommissionType = request.CommissionType
            };

            var response = await _gateway.CalculateCommissionAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Commission calculation failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Commission calculated successfully: {CommissionNumber}, Amount: {Amount}", 
                    response.CommissionNumber, response.Amount?.Amount);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to calculate commission via gateway");
            return new CalculateCommissionResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== CreatePayout =====
public sealed class CreatePayoutCommandHandler : IRequestHandler<CreatePayoutCommand, CreatePayoutResponse>
{
    private readonly ICommissionDataGateway _gateway;
    private readonly ILogger<CreatePayoutCommandHandler> _logger;

    public CreatePayoutCommandHandler(
        ICommissionDataGateway gateway,
        ILogger<CreatePayoutCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CreatePayoutResponse> Handle(CreatePayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Creating payout batch for: {RecipientId}", request.RecipientId);

            var grpcRequest = new CreatePayoutRequest
            {
                RecipientId = request.RecipientId,
                RecipientType = request.RecipientType,
                PeriodStart = request.PeriodStart,
                PeriodEnd = request.PeriodEnd
            };

            if (request.CommissionIds != null)
            {
                grpcRequest.CommissionIds.AddRange(request.CommissionIds);
            }

            var response = await _gateway.CreatePayoutAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Payout creation failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Payout batch created successfully: {PayoutNumber}, Commissions: {Count}", 
                    response.PayoutNumber, response.CommissionCount);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create payout via gateway");
            return new CreatePayoutResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}

// ===== ProcessPayout =====
public sealed class ProcessPayoutCommandHandler : IRequestHandler<ProcessPayoutCommand, ProcessPayoutResponse>
{
    private readonly ICommissionDataGateway _gateway;
    private readonly ILogger<ProcessPayoutCommandHandler> _logger;

    public ProcessPayoutCommandHandler(
        ICommissionDataGateway gateway,
        ILogger<ProcessPayoutCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<ProcessPayoutResponse> Handle(ProcessPayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Processing payout: {PayoutId}", request.PayoutId);

            var grpcRequest = new ProcessPayoutRequest
            {
                PayoutId = request.PayoutId,
                PaymentMethod = request.PaymentMethod,
                PaymentReference = request.PaymentReference
            };

            var response = await _gateway.ProcessPayoutAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Payout processing failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Payout processed successfully for {PayoutId}", request.PayoutId);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process payout via gateway");
            return new ProcessPayoutResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
