using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class UpdateBeneficiaryCommandHandler : IRequestHandler<UpdateBeneficiaryCommand, UpdateBeneficiaryResponse>
{
    private readonly IBeneficiaryDataGateway _gateway;
    private readonly ILogger<UpdateBeneficiaryCommandHandler> _logger;

    public UpdateBeneficiaryCommandHandler(
        IBeneficiaryDataGateway gateway,
        ILogger<UpdateBeneficiaryCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<UpdateBeneficiaryResponse> Handle(UpdateBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Updating beneficiary: {BeneficiaryId}", request.BeneficiaryId);

            var grpcRequest = new UpdateBeneficiaryRequest
            {
                BeneficiaryId = request.BeneficiaryId,
                MobileNumber = request.MobileNumber ?? "",
                Email = request.Email ?? "",
                Address = request.Address ?? ""
            };

            var response = await _gateway.UpdateBeneficiaryAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Beneficiary update failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Beneficiary updated successfully: {BeneficiaryId}", request.BeneficiaryId);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update beneficiary via gateway");
            return new UpdateBeneficiaryResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
