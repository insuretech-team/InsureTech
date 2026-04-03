using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CreateIndividualBeneficiaryCommandHandler : IRequestHandler<CreateIndividualBeneficiaryCommand, CreateIndividualBeneficiaryResponse>
{
    private readonly IBeneficiaryDataGateway _gateway;
    private readonly ILogger<CreateIndividualBeneficiaryCommandHandler> _logger;

    public CreateIndividualBeneficiaryCommandHandler(
        IBeneficiaryDataGateway gateway,
        ILogger<CreateIndividualBeneficiaryCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<CreateIndividualBeneficiaryResponse> Handle(CreateIndividualBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Creating individual beneficiary for user: {UserId}", request.UserId);

            var grpcRequest = new CreateIndividualBeneficiaryRequest
            {
                UserId = request.UserId,
                FullName = request.FullName,
                DateOfBirth = request.DateOfBirth.ToString("yyyy-MM-dd"),
                Gender = request.Gender,
                NidNumber = request.NidNumber,
                MobileNumber = request.MobileNumber,
                Email = request.Email ?? "",
                PartnerId = request.PartnerId ?? ""
            };

            var response = await _gateway.CreateIndividualBeneficiaryAsync(grpcRequest, cancellationToken);
            
            if (response.Error != null)
            {
                _logger.LogWarning("Individual beneficiary creation failed: {ErrorCode} - {ErrorMessage}", response.Error.Code, response.Error.Message);
            }
            else
            {
                _logger.LogInformation("Individual beneficiary created successfully: {BeneficiaryId} ({BeneficiaryCode})", 
                    response.BeneficiaryId, response.BeneficiaryCode);
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create individual beneficiary via gateway");
            return new CreateIndividualBeneficiaryResponse { Error = new Error { Code = "GATEWAY_ERROR", Message = ex.Message } };
        }
    }
}
