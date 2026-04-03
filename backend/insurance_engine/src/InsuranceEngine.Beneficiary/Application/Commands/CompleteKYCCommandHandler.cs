using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Grpc.Gateways;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, CompleteKYCResponse>
{
    private readonly IBeneficiaryDataGateway _gateway;
    private readonly ILogger<CompleteKYCCommandHandler> _logger;

    public CompleteKYCCommandHandler(
        IBeneficiaryDataGateway gateway,
        ILogger<CompleteKYCCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CompleteKYCResponse> Handle(CompleteKYCCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Completing KYC for beneficiary: {BeneficiaryId}", request.BeneficiaryId);

            var grpcRequest = new CompleteKYCRequest
            {
                BeneficiaryId = request.BeneficiaryId,
                NidFrontUrl = request.NidFrontUrl,
                NidBackUrl = request.NidBackUrl,
                SelfieUrl = request.SelfieUrl,
                PorichoyVerificationId = request.PorichoyVerificationId ?? ""
            };

            var response = await _gateway.CompleteKYCAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("KYC completion failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("KYC completed successfully for beneficiary: {BeneficiaryId}", request.BeneficiaryId);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to complete KYC via gateway");
            return new CompleteKYCResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
