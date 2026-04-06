using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.Endorsements.Infrastructure;

namespace InsuranceEngine.Endorsements.Application.Commands;

public sealed class UpdatePolicyCommandHandler : IRequestHandler<UpdatePolicyCommand, UpdatePolicyResponse>
{
    private readonly ISqlEndorsementDataGateway _gateway;
    private readonly ILogger<UpdatePolicyCommandHandler> _logger;

    public UpdatePolicyCommandHandler(
        ISqlEndorsementDataGateway gateway,
        ILogger<UpdatePolicyCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<UpdatePolicyResponse> Handle(UpdatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            _logger.LogInformation("Processing policy endorsement: {PolicyId}", request.PolicyId);

            var endorsementId = Guid.NewGuid().ToString();
            var endorsementNumber = $"END-{DateTime.UtcNow:yyyyMMdd}-{endorsementId[..8].ToUpper()}";

            var changes = new Dictionary<string, object>();
            if (request.Nominees != null)
            {
                changes["nominees"] = request.Nominees.Select(n => new { 
                    fullName = n.FullName, 
                    relationship = n.Relationship,
                    sharePercentage = n.SharePercentage 
                }).ToList();
            }
            if (!string.IsNullOrEmpty(request.Address))
            {
                changes["address"] = request.Address;
            }

            var endorsement = new InsuranceEngine.SharedKernel.Persistence.Entities.EndorsementEntity
            {
                EndorsementId = endorsementId,
                EndorsementNumber = endorsementNumber,
                PolicyId = request.PolicyId,
                Type = request.Nominees != null ? EndorsementType.NomineeChange : EndorsementType.ContactChange,
                Status = "PROCESSED",
                Changes = System.Text.Json.JsonSerializer.Serialize(changes),
                EffectiveDate = DateTime.UtcNow,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _gateway.CreateAsync(endorsement, cancellationToken);

            _logger.LogInformation("Policy endorsement created successfully: {EndorsementNumber}", endorsementNumber);

            return new UpdatePolicyResponse();
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process policy endorsement: {PolicyId}", request.PolicyId);
            return new UpdatePolicyResponse { Error = new Error { Code = "ENDORSEMENT_ERROR", Message = ex.Message } };
        }
    }
}
