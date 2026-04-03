using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CreateBusinessBeneficiaryCommandHandler : IRequestHandler<CreateBusinessBeneficiaryCommand, CreateBusinessBeneficiaryResponse>
{
    private readonly IBeneficiaryDataGateway _gateway;
    private readonly ILogger<CreateBusinessBeneficiaryCommandHandler> _logger;

    public CreateBusinessBeneficiaryCommandHandler(
        IBeneficiaryDataGateway gateway,
        ILogger<CreateBusinessBeneficiaryCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CreateBusinessBeneficiaryResponse> Handle(CreateBusinessBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Creating business beneficiary for user: {UserId}", request.UserId);

            var grpcRequest = new CreateBusinessBeneficiaryRequest
            {
                UserId = request.UserId,
                BusinessName = request.BusinessName,
                TradeLicenseNumber = request.TradeLicenseNumber,
                TinNumber = request.TinNumber,
                FocalPersonName = request.FocalPersonName,
                FocalPersonMobile = request.FocalPersonMobile,
                PartnerId = request.PartnerId ?? ""
            };

            var response = await _gateway.CreateBusinessBeneficiaryAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Business beneficiary creation failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Business beneficiary created successfully: {BeneficiaryId} ({BeneficiaryCode})", 
                    response.BeneficiaryId, response.BeneficiaryCode);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create business beneficiary via gateway");
            return new CreateBusinessBeneficiaryResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
