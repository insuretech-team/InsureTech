using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record UpdateBeneficiaryCommand(
    Guid BeneficiaryId,
    string? Name = null,
    string? Email = null,
    string? ContactNumber = null,
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

        if (beneficiary.Type == Domain.Enums.BeneficiaryType.Individual && beneficiary.IndividualDetails != null)
        {
            // Placeholder: In a real app, we'd have methods on the detail entity
            // For now, aligning with the domain change
        }
        else if (beneficiary.Type == Domain.Enums.BeneficiaryType.Business && beneficiary.BusinessDetails != null)
        {
            // Placeholder: In a real app, we'd have methods on the detail entity
        }

        beneficiary.UpdatedAt = DateTime.UtcNow;
        await _repository.UpdateAsync(beneficiary);
        return Result.Ok(true);
    }
}
