using System;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Newtonsoft.Json;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.Endorsements;

public record SubmitEndorsementCommand(
    Guid PolicyId,
    EndorsementType Type,
    string Reason,
    object Changes,
    Guid RequestedBy
) : IRequest<Result<Endorsement>>;

public class SubmitEndorsementCommandHandler : IRequestHandler<SubmitEndorsementCommand, Result<Endorsement>>
{
    private readonly IEndorsementRepository _endorsementRepo;
    private readonly IPolicyRepository _policyRepo;
    private readonly EndorsementNumberGenerator _numberGenerator;

    public SubmitEndorsementCommandHandler(
        IEndorsementRepository endorsementRepo,
        IPolicyRepository policyRepo,
        EndorsementNumberGenerator numberGenerator)
    {
        _endorsementRepo = endorsementRepo;
        _policyRepo = policyRepo;
        _numberGenerator = numberGenerator;
    }

    public async Task<Result<Endorsement>> Handle(SubmitEndorsementCommand request, CancellationToken cancellationToken)
    {
        var policy = await _policyRepo.GetByIdAsync(request.PolicyId);
        if (policy == null) return Result<Endorsement>.Failure("Policy not found");

        if (policy.Status != PolicyStatus.Active && policy.Status != PolicyStatus.GracePeriod)
            return Result<Endorsement>.Failure("Policy must be Active or in Grace Period for endorsements");

        var endorsement = new Endorsement
        {
            Id = Guid.NewGuid(),
            PolicyId = request.PolicyId,
            EndorsementNumber = await _numberGenerator.GenerateNumberAsync(policy.PolicyNumber),
            Type = request.Type,
            Reason = request.Reason,
            ChangesJson = JsonConvert.SerializeObject(request.Changes),
            Status = EndorsementStatus.Pending,
            RequestedBy = request.RequestedBy,
            EffectiveDate = DateTime.UtcNow,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        await _endorsementRepo.AddAsync(endorsement);
        return Result<Endorsement>.Success(endorsement);
    }
}
