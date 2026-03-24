using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record CompleteKYCCommand(
    Guid BeneficiaryId,
    string Status,
    string? NidFrontUrl = null,
    string? NidBackUrl = null,
    string? SelfieUrl = null,
    string? PorichoyVerificationId = null
) : IRequest<Result<bool>>;

public class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, Result<bool>>
{
    private readonly IBeneficiaryRepository _repository;

    public CompleteKYCCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<bool>> Handle(CompleteKYCCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null) return Result<bool>.Failure("Beneficiary not found");

        var status = Enum.Parse<KYCStatus>(request.Status, true);
        beneficiary.CompleteKYC(status, request.NidFrontUrl, request.NidBackUrl, request.SelfieUrl, request.PorichoyVerificationId);

        await _repository.UpdateAsync(beneficiary);
        return Result.Ok(true);
    }
}
