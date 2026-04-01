using MediatR;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, CompleteKYCResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repository;
    public CompleteKYCCommandHandler(IRepository<BeneficiaryEntity> repository) => _repository = repository;

    public async Task<CompleteKYCResponse> Handle(CompleteKYCCommand request, CancellationToken cancellationToken)
    {
        var e = await _repository.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
        if (e == null) return new CompleteKYCResponse { Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

        e.KycStatus = "COMPLETED";
        e.KycCompletedAt = DateTime.UtcNow;
        if (e.Status == "PENDING_KYC") e.Status = "ACTIVE";
        e.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(e, cancellationToken);
        return new CompleteKYCResponse { Message = "KYC completed successfully" };
    }
}
