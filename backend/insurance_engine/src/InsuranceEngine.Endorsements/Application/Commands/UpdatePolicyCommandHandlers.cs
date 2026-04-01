using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Endorsements.Application.Commands;

public sealed class UpdatePolicyCommandHandler : IRequestHandler<UpdatePolicyCommand, UpdatePolicyResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly ILogger<UpdatePolicyCommandHandler> _logger;

    public UpdatePolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        ILogger<UpdatePolicyCommandHandler> logger)
    {
        _repository = repository;
        _nomineeRepository = nomineeRepository;
        _logger = logger;
    }

    public async Task<UpdatePolicyResponse> Handle(UpdatePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var policy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
            {
                return new UpdatePolicyResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            // Update nominees if provided
            if (request.Nominees != null && request.Nominees.Count > 0)
            {
                // Remove existing nominees
                var existingNominees = await _nomineeRepository.FindAsync(n => n.PolicyId == policy.PolicyId, cancellationToken);
                foreach (var existing in existingNominees)
                {
                    await _nomineeRepository.DeleteAsync(existing, cancellationToken);
                }

                // Add new nominees
                foreach (var nominee in request.Nominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = policy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth?.ToDateTime() ?? DateTime.UtcNow,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }

            policy.UpdatedAt = DateTime.UtcNow;
            await _repository.UpdateAsync(policy, cancellationToken);

            _logger.LogInformation("Policy updated: {PolicyId}", request.PolicyId);

            return new UpdatePolicyResponse { Message = "Policy updated successfully" };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update policy {PolicyId}", request.PolicyId);
            return new UpdatePolicyResponse
            {
                Error = new Error { Code = "UPDATE_FAILED", Message = ex.Message }
            };
        }
    }
}
