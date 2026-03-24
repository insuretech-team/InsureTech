using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record UpdateRiskScoreCommand(
    Guid BeneficiaryId,
    string RiskScore,
    string? Reason = null
) : IRequest<Result<bool>>;

public class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, Result<bool>>
{
    private readonly IBeneficiaryRepository _repository;

    public UpdateRiskScoreCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<bool>> Handle(UpdateRiskScoreCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null) return Result<bool>.Failure("Beneficiary not found");

        beneficiary.UpdateRiskScore(request.RiskScore, request.Reason);

        await _repository.UpdateAsync(beneficiary);
        return Result.Ok(true);
    }
}
