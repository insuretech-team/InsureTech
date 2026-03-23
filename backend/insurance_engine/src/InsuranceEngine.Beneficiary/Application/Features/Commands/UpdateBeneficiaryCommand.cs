using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record UpdateBeneficiaryCommand(
    Guid BeneficiaryId,
    string? MobileNumber = null,
    string? Email = null,
    string? Address = null
) : IRequest<Result<bool>>;

public class UpdateBeneficiaryCommandHandler : IRequestHandler<UpdateBeneficiaryCommand, Result<bool>>
{
    private readonly IBeneficiaryRepository _repository;

    public UpdateBeneficiaryCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<bool>> Handle(UpdateBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null) return Result<bool>.Fail(Error.NotFound("Beneficiary", request.BeneficiaryId.ToString()));

        if (!string.IsNullOrEmpty(request.MobileNumber))
        {
            if (beneficiary.Type == Domain.Enums.BeneficiaryType.Individual && beneficiary.Individual != null)
            {
                beneficiary.Individual.ContactInfo = beneficiary.Individual.ContactInfo with { MobileNumber = request.MobileNumber };
            }
            else if (beneficiary.Type == Domain.Enums.BeneficiaryType.Business && beneficiary.Business != null)
            {
                beneficiary.Business.ContactInfo = beneficiary.Business.ContactInfo with { MobileNumber = request.MobileNumber };
            }
        }

        if (!string.IsNullOrEmpty(request.Email))
        {
            if (beneficiary.Type == Domain.Enums.BeneficiaryType.Individual && beneficiary.Individual != null)
            {
                beneficiary.Individual.ContactInfo = beneficiary.Individual.ContactInfo with { Email = request.Email };
            }
            else if (beneficiary.Type == Domain.Enums.BeneficiaryType.Business && beneficiary.Business != null)
            {
                beneficiary.Business.ContactInfo = beneficiary.Business.ContactInfo with { Email = request.Email };
            }
        }

        if (!string.IsNullOrEmpty(request.Address))
        {
            if (beneficiary.Type == Domain.Enums.BeneficiaryType.Individual && beneficiary.Individual != null)
            {
                beneficiary.Individual.PermanentAddress = beneficiary.Individual.PermanentAddress with { AddressLine1 = request.Address };
            }
            else if (beneficiary.Type == Domain.Enums.BeneficiaryType.Business && beneficiary.Business != null)
            {
                beneficiary.Business.RegisteredAddress = beneficiary.Business.RegisteredAddress with { AddressLine1 = request.Address };
            }
        }

        await _repository.UpdateAsync(beneficiary);
        return Result.Ok(true);
    }
}
