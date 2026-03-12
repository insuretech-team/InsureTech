using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

public record UpdateBeneficiaryCommand(
    Guid BeneficiaryId,
    string? MobileNumber = null,
    string? Email = null,
    string? Address = null
) : IRequest<Result>;

public class UpdateBeneficiaryCommandHandler : IRequestHandler<UpdateBeneficiaryCommand, Result>
{
    private readonly IBeneficiaryRepository _repository;

    public UpdateBeneficiaryCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result> Handle(UpdateBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null)
            return Result.Fail(Error.NotFound("Beneficiary", request.BeneficiaryId.ToString()));

        // Update top-level info if needed (usually it's in individual/business details)
        // For simplicity, we'll assume we update the JSON fields in the specialized details
        if (beneficiary.Type == BeneficiaryType.Individual && beneficiary.IndividualDetails != null)
        {
            // Simple string replacement for now, or proper JSON parsing if needed
            beneficiary.IndividualDetails.ContactInfoJson = $"{{\"mobile\": \"{request.MobileNumber ?? ""}\", \"email\": \"{request.Email ?? ""}\"}}";
            if (!string.IsNullOrEmpty(request.Address))
            {
                beneficiary.IndividualDetails.PresentAddressJson = $"{{\"address\": \"{request.Address}\"}}";
            }
        }
        else if (beneficiary.Type == BeneficiaryType.Business && beneficiary.BusinessDetails != null)
        {
            beneficiary.BusinessDetails.ContactInfoJson = $"{{\"mobile\": \"{request.MobileNumber ?? ""}\", \"email\": \"{request.Email ?? ""}\"}}";
        }

        beneficiary.UpdatedAt = DateTime.UtcNow;
        await _repository.UpdateAsync(beneficiary);

        return Result.Ok();
    }
}
