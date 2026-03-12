using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

public record CompleteKYCCommand(
    Guid BeneficiaryId,
    string NidFrontUrl,
    string NidBackUrl,
    string SelfieUrl,
    string? PorichoyVerificationId = null
) : IRequest<Result>;

public class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, Result>
{
    private readonly IBeneficiaryRepository _repository;

    public CompleteKYCCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result> Handle(CompleteKYCCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null)
            return Result.Fail(Error.NotFound("Beneficiary", request.BeneficiaryId.ToString()));

        // Simulate KYC success
        beneficiary.KycStatus = KYCStatus.Completed;
        beneficiary.Status = BeneficiaryStatus.Active;
        beneficiary.KycCompletedAt = DateTime.UtcNow;
        beneficiary.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(beneficiary);

        return Result.Ok();
    }
}
