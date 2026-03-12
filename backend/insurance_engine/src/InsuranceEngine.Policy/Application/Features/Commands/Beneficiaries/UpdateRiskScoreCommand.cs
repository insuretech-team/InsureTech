using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

public record UpdateRiskScoreCommand(
    Guid BeneficiaryId,
    string RiskScore,
    string Reason
) : IRequest<Result>;

public class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, Result>
{
    private readonly IBeneficiaryRepository _repository;

    public UpdateRiskScoreCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result> Handle(UpdateRiskScoreCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null)
            return Result.Fail(Error.NotFound("Beneficiary", request.BeneficiaryId.ToString()));

        beneficiary.RiskScore = request.RiskScore;
        beneficiary.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(beneficiary);

        return Result.Ok();
    }
}
