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
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string NidNumber,
    string MobileNumber,
    string? Email = null,
    Guid? PartnerId = null
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
        var gender = Enum.Parse<Gender>(request.Gender, true);
        
        var beneficiary = Domain.Entities.Beneficiary.CreateIndividual(
            request.UserId,
            request.FullName,
            request.DateOfBirth,
            gender,
            request.MobileNumber,
            request.Email);

        if (!string.IsNullOrEmpty(request.NidNumber) && beneficiary.IndividualDetails != null)
        {
            beneficiary.IndividualDetails.NidNumber = await _encryptionService.EncryptAsync(request.NidNumber);
        }

        await _repository.AddAsync(beneficiary);

        return Result.Ok(MapToDto(beneficiary));
    }

    private BeneficiaryDto MapToDto(Domain.Entities.Beneficiary b)
    {
        return new BeneficiaryDto(
            b.Id,
            b.UserId,
            b.Type.ToString(),
            b.Code,
            b.Status.ToString(),
            b.KycStatus.ToString(),
            b.KycCompletedAt,
            b.RiskScore,
            null,
            b.IndividualDetails != null ? new IndividualBeneficiaryDto(
                b.IndividualDetails.FullName,
                b.IndividualDetails.FullNameBn,
                b.IndividualDetails.DateOfBirth,
                b.IndividualDetails.Gender.ToString(),
                b.IndividualDetails.NidNumber, // In real app, might want to mask or keep encrypted
                b.IndividualDetails.PassportNumber,
                b.IndividualDetails.BirthCertificateNumber,
                b.IndividualDetails.TinNumber,
                b.IndividualDetails.MaritalStatus.ToString(),
                b.IndividualDetails.Occupation,
                b.IndividualDetails.ContactInfoJson,
                b.IndividualDetails.PermanentAddressJson,
                b.IndividualDetails.PresentAddressJson,
                null,
                null
            ) : null
        );
    }
}
