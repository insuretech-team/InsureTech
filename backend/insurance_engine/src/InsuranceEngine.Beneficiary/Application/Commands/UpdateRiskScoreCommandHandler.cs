using MediatR;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, UpdateRiskScoreResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repository;
    public UpdateRiskScoreCommandHandler(IRepository<BeneficiaryEntity> repository) => _repository = repository;

    public async Task<UpdateRiskScoreResponse> Handle(UpdateRiskScoreCommand request, CancellationToken cancellationToken)
    {
        var e = await _repository.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
        if (e == null) return new UpdateRiskScoreResponse { Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

        e.RiskScore = request.RiskScore;
        e.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(e, cancellationToken);
        return new UpdateRiskScoreResponse { Message = "Risk score updated" };
    }
}
