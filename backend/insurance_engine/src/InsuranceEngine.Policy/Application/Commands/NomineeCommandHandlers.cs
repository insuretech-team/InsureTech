using MediatR;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class AddNomineeCommandHandler : IRequestHandler<AddNomineeCommand, Result<string>>
{
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly IRepository<PolicyEntity> _policyRepository;
    private readonly ILogger<AddNomineeCommandHandler> _logger;

    public AddNomineeCommandHandler(
        IRepository<PolicyNomineeEntity> nomineeRepository,
        IRepository<PolicyEntity> policyRepository,
        ILogger<AddNomineeCommandHandler> logger)
    {
        _nomineeRepository = nomineeRepository;
        _policyRepository = policyRepository;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(AddNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            // Validate policy exists
            var policy = await _policyRepository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (policy == null)
                return Result<string>.NotFound("POLICY_NOT_FOUND", "Policy not found");

            var entity = new PolicyNomineeEntity
            {
                NomineeId = Guid.NewGuid(),
                PolicyId = Guid.Parse(request.PolicyId),
                FullName = request.FullName,
                Relationship = request.Relationship,
                SharePercentage = request.SharePercentage,
                DateOfBirth = request.DateOfBirth ?? DateTime.UtcNow,
                NidNumber = request.NidNumber,
                PhoneNumber = request.PhoneNumber,
                NomineeDobText = request.NomineeDobText,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _nomineeRepository.AddAsync(entity, cancellationToken);

            _logger.LogInformation("Nominee added: {NomineeId} to Policy: {PolicyId}", entity.NomineeId, request.PolicyId);
            return Result<string>.Ok(entity.NomineeId.ToString());
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
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly ILogger<UpdateNomineeCommandHandler> _logger;

    public UpdateNomineeCommandHandler(IRepository<PolicyNomineeEntity> nomineeRepository, ILogger<UpdateNomineeCommandHandler> logger)
    {
        _nomineeRepository = nomineeRepository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(UpdateNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var nominee = await _nomineeRepository.GetByIdAsync(Guid.Parse(request.NomineeId), cancellationToken);
            if (nominee == null || nominee.PolicyId != Guid.Parse(request.PolicyId))
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found");

            if (request.FullName != null) nominee.FullName = request.FullName;
            if (request.Relationship != null) nominee.Relationship = request.Relationship;
            if (request.SharePercentage.HasValue) nominee.SharePercentage = request.SharePercentage.Value;
            if (request.DateOfBirth.HasValue) nominee.DateOfBirth = request.DateOfBirth.Value;
            if (request.NidNumber != null) nominee.NidNumber = request.NidNumber;
            if (request.PhoneNumber != null) nominee.PhoneNumber = request.PhoneNumber;
            
            nominee.UpdatedAt = DateTime.UtcNow;
            await _nomineeRepository.UpdateAsync(nominee, cancellationToken);

            _logger.LogInformation("Nominee updated: {NomineeId}", request.NomineeId);
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
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly ILogger<DeleteNomineeCommandHandler> _logger;

    public DeleteNomineeCommandHandler(IRepository<PolicyNomineeEntity> nomineeRepository, ILogger<DeleteNomineeCommandHandler> logger)
    {
        _nomineeRepository = nomineeRepository;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(DeleteNomineeCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var nominee = await _nomineeRepository.GetByIdAsync(Guid.Parse(request.NomineeId), cancellationToken);
            if (nominee == null || nominee.PolicyId != Guid.Parse(request.PolicyId))
                return Result<bool>.NotFound("NOMINEE_NOT_FOUND", "Nominee not found");

            await _nomineeRepository.DeleteAsync(nominee, cancellationToken);

            _logger.LogInformation("Nominee deleted: {NomineeId}", request.NomineeId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to delete nominee {NomineeId}", request.NomineeId);
            return Result<bool>.Fail("NOMINEE_DELETE_FAILED", ex.Message);
        }
    }
}
