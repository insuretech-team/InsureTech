using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

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
        var code = await _repository.GetNextSequenceAsync();
        
        var beneficiary = new Beneficiary
        {
            Id = Guid.NewGuid(),
            UserId = request.UserId,
            Type = BeneficiaryType.Individual,
            Code = code,
            Status = BeneficiaryStatus.PendingKyc,
            KycStatus = KYCStatus.NotStarted,
            PartnerId = request.PartnerId,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        var encryptedNid = await _encryptionService.EncryptAsync(request.NidNumber);

        var individualDetails = new IndividualBeneficiary
        {
            BeneficiaryId = beneficiary.Id,
            FullName = request.FullName,
            DateOfBirth = request.DateOfBirth,
            Gender = Enum.Parse<Gender>(request.Gender, true),
            NidNumber = encryptedNid,
            ContactInfoJson = $"{{\"mobile\": \"{request.MobileNumber}\", \"email\": \"{request.Email}\"}}",
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        beneficiary.IndividualDetails = individualDetails;

        await _repository.AddAsync(beneficiary);

        return Result.Ok(new BeneficiaryDto(
            beneficiary.Id,
            beneficiary.UserId,
            beneficiary.Type.ToString(),
            beneficiary.Code,
            beneficiary.Status.ToString(),
            beneficiary.KycStatus.ToString(),
            beneficiary.KycCompletedAt,
            beneficiary.RiskScore,
            beneficiary.ReferralCode,
            new IndividualBeneficiaryDto(
                individualDetails.FullName,
                null,
                individualDetails.DateOfBirth,
                individualDetails.Gender.ToString(),
                request.NidNumber, // Return plaintext for UI if needed, or masked
                null,
                null,
                null,
                MaritalStatus.Unspecified.ToString(),
                null,
                individualDetails.ContactInfoJson,
                null,
                null,
                null,
                null
            )
        ));
    }
}
