using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.Grpc.Gateways;
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
            var policy = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policy == null)
                return Result<string>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var newNominee = new Nominee
            {
                NomineeId = Guid.NewGuid().ToString(),
                FullName = request.FullName,
                Relationship = request.Relationship,
                SharePercentage = (double)request.SharePercentage,
                DateOfBirth = request.DateOfBirth.HasValue ? Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(request.DateOfBirth.Value.ToUniversalTime()) : null,
                NidNumber = request.NidNumber,
                PhoneNumber = request.PhoneNumber,
                NomineeDobText = request.NomineeDobText
            };

            policy.Nominees.Add(newNominee);
            await _gateway.UpdatePolicyAsync(policy, cancellationToken);

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
            var policy = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policy == null)
                return Result<bool>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var nominee = policy.Nominees.FirstOrDefault(n => n.NomineeId == request.NomineeId);
            if (nominee == null)
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found in policy");

            if (request.FullName != null) nominee.FullName = request.FullName;
            if (request.Relationship != null) nominee.Relationship = request.Relationship;
            if (request.SharePercentage.HasValue) nominee.SharePercentage = (double)request.SharePercentage.Value;
            if (request.DateOfBirth.HasValue) 
                nominee.DateOfBirth = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(request.DateOfBirth.Value.ToUniversalTime());
            if (request.NidNumber != null) nominee.NidNumber = request.NidNumber;
            if (request.PhoneNumber != null) nominee.PhoneNumber = request.PhoneNumber;
            
            await _gateway.UpdatePolicyAsync(policy, cancellationToken);

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
            var policy = await _gateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (policy == null)
                return Result<bool>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var nominee = policy.Nominees.FirstOrDefault(n => n.NomineeId == request.NomineeId);
            if (nominee == null)
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found in policy");

            policy.Nominees.Remove(nominee);
            await _gateway.UpdatePolicyAsync(policy, cancellationToken);

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
