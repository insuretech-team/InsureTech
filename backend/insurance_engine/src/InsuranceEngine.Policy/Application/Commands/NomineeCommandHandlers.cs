using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Policy.Entity.V1;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class AddNomineeCommandHandler : IRequestHandler<AddNomineeCommand, Result<string>>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly ILogger<AddNomineeCommandHandler> _logger;

    public AddNomineeCommandHandler(
        IPolicyDataGateway gateway,
        ILogger<AddNomineeCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(AddNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
                return Result<string>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            _logger.LogDebug("Creating nominee: {FullName}, {Relationship}, {Share}", request.FullName, request.Relationship, request.SharePercentage);
            
            Nominee newNominee;
            newNominee = new Nominee
            {
                NomineeId = Guid.NewGuid().ToString(),
                PolicyId = request.PolicyId,
                FullName = request.FullName ?? "",
                Relationship = request.Relationship ?? "",
                SharePercentage = (double)request.SharePercentage,
                DateOfBirth = request.DateOfBirth.HasValue ? Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(request.DateOfBirth.Value.ToUniversalTime()) : null,
                NidNumber = request.NidNumber ?? "",
                PhoneNumber = request.PhoneNumber ?? "",
                NomineeDobText = request.NomineeDobText ?? ""
            };

            _logger.LogDebug("Getting existing nominees, count: {Count}", policyResponse.Policy.Nominees?.Count ?? 0);
            
            var nomineesToUpdate = new List<Nominee>();
            if (policyResponse.Policy.Nominees != null)
            {
                foreach (var existing in policyResponse.Policy.Nominees)
                {
                    nomineesToUpdate.Add(existing);
                }
            }
            nomineesToUpdate.Add(newNominee);
            
            _logger.LogDebug("After adding new nominee, count: {Count}", nomineesToUpdate.Count);

            var updateResponse = await _gateway.UpdatePolicyAsync(request.PolicyId, nomineesToUpdate, null, cancellationToken);
            
            if (updateResponse.Error != null)
                return Result<string>.Fail("UPDATE_FAILED", updateResponse.Error.Message);

            _logger.LogInformation("Nominee added via Go SSOT: {NomineeId} to Policy: {PolicyId}", newNominee.NomineeId, request.PolicyId);
            return Result<string>.Ok(newNominee.NomineeId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to add nominee to policy {PolicyId}", request.PolicyId);
            return Result<string>.Fail("NOMINEE_ADD_FAILED", ex.Message);
        }
    }
}

public sealed class UpdateNomineeCommandHandler : IRequestHandler<UpdateNomineeCommand, Result<bool>>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly ILogger<UpdateNomineeCommandHandler> _logger;

    public UpdateNomineeCommandHandler(IPolicyDataGateway gateway, ILogger<UpdateNomineeCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
                return Result<bool>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var nominees = policyResponse.Policy.Nominees.ToList();
            var nominee = nominees.FirstOrDefault(n => n.NomineeId == request.NomineeId);
            if (nominee == null)
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found in policy");

            if (request.FullName != null) nominee.FullName = request.FullName;
            if (request.Relationship != null) nominee.Relationship = request.Relationship;
            if (request.SharePercentage.HasValue) nominee.SharePercentage = (double)request.SharePercentage.Value;
            if (request.DateOfBirth.HasValue) 
                nominee.DateOfBirth = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(request.DateOfBirth.Value.ToUniversalTime());
            if (request.NidNumber != null) nominee.NidNumber = request.NidNumber;
            if (request.PhoneNumber != null) nominee.PhoneNumber = request.PhoneNumber;
            
            var updateResponse = await _gateway.UpdatePolicyAsync(request.PolicyId, nominees, null, cancellationToken);
            
            if (updateResponse.Error != null)
                return Result<bool>.Fail("UPDATE_FAILED", updateResponse.Error.Message);

            _logger.LogInformation("Nominee updated via Go SSOT: {NomineeId}", request.NomineeId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update nominee {NomineeId}", request.NomineeId);
            return Result<bool>.Fail("NOMINEE_UPDATE_FAILED", ex.Message);
        }
    }
}

public sealed class DeleteNomineeCommandHandler : IRequestHandler<DeleteNomineeCommand, Result<bool>>
{
    private readonly IPolicyDataGateway _gateway;
    private readonly ILogger<DeleteNomineeCommandHandler> _logger;

    public DeleteNomineeCommandHandler(IPolicyDataGateway gateway, ILogger<DeleteNomineeCommandHandler> logger)
    {
        _gateway = gateway;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeleteNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policyResponse = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policyResponse.Policy == null)
                return Result<bool>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var nominees = policyResponse.Policy.Nominees.ToList();
            var nominee = nominees.FirstOrDefault(n => n.NomineeId == request.NomineeId);
            if (nominee == null)
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found in policy");

            nominees.Remove(nominee);
            
            var updateResponse = await _gateway.UpdatePolicyAsync(request.PolicyId, nominees, null, cancellationToken);
            
            if (updateResponse.Error != null)
                return Result<bool>.Fail("UPDATE_FAILED", updateResponse.Error.Message);

            _logger.LogInformation("Nominee deleted via Go SSOT: {NomineeId}", request.NomineeId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to delete nominee {NomineeId}", request.NomineeId);
            return Result<bool>.Fail("NOMINEE_DELETE_FAILED", ex.Message);
        }
    }
}
