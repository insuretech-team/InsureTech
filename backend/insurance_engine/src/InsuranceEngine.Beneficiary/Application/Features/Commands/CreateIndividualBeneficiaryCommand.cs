using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record CreateIndividualBeneficiaryCommand(
    Guid UserId,
    string Email,
    Guid PartnerId,
    string? FullName = null,
    DateTime? DateOfBirth = null,
    string? Gender = null,
    string? NidNumber = null,
    string? MobileNumber = null
) : IRequest<Result<BeneficiaryDto>>;

public class CreateIndividualBeneficiaryCommandHandler : IRequestHandler<CreateIndividualBeneficiaryCommand, Result<BeneficiaryDto>>
{
    private readonly IBeneficiaryRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public CreateIndividualBeneficiaryCommandHandler(IBeneficiaryRepository repository, IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<Result<BeneficiaryDto>> Handle(CreateIndividualBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = new Domain.Entities.Beneficiary(
            Guid.NewGuid(),
            request.UserId,
            BeneficiaryType.Individual);

        if (beneficiary.Individual != null)
        {
            if (!string.IsNullOrEmpty(request.FullName))
                beneficiary.Individual.FullName = request.FullName;
            if (request.DateOfBirth.HasValue)
                beneficiary.Individual.DateOfBirth = request.DateOfBirth.Value;
            if (!string.IsNullOrEmpty(request.Gender))
                beneficiary.Individual.Gender = Enum.Parse<BeneficiaryGender>(request.Gender, true);
            if (!string.IsNullOrEmpty(request.MobileNumber))
                beneficiary.Individual.ContactInfo = new SharedKernel.Domain.ValueObjects.ContactInfo(request.MobileNumber);
            if (!string.IsNullOrEmpty(request.Email))
                beneficiary.Individual.ContactInfo = beneficiary.Individual.ContactInfo with { Email = request.Email };
            if (!string.IsNullOrEmpty(request.NidNumber))
                beneficiary.Individual.NidNumber = await _encryptionService.EncryptAsync(request.NidNumber);
        }

        beneficiary.PartnerId = request.PartnerId;
        beneficiary.UpdateCode();

        await _repository.AddAsync(beneficiary);

        return Result.Ok(beneficiary.ToDto());
    }
}
