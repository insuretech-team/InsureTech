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

public record CreateBusinessBeneficiaryCommand(
    Guid UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    Guid? PartnerId = null
) : IRequest<Result<BeneficiaryDto>>;

public class CreateBusinessBeneficiaryCommandHandler : IRequestHandler<CreateBusinessBeneficiaryCommand, Result<BeneficiaryDto>>
{
    private readonly IBeneficiaryRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public CreateBusinessBeneficiaryCommandHandler(IBeneficiaryRepository repository, IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<Result<BeneficiaryDto>> Handle(CreateBusinessBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var code = await _repository.GetNextSequenceAsync();

        var beneficiary = new Beneficiary
        {
            Id = Guid.NewGuid(),
            UserId = request.UserId,
            Type = BeneficiaryType.Business,
            Code = code,
            Status = BeneficiaryStatus.PendingKyc,
            KycStatus = KYCStatus.NotStarted,
            PartnerId = request.PartnerId,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        var businessDetails = new BusinessBeneficiary
        {
            Id = Guid.NewGuid(),
            BeneficiaryId = beneficiary.Id,
            BusinessName = request.BusinessName,
            TradeLicenseNumber = request.TradeLicenseNumber,
            TinNumber = request.TinNumber,
            FocalPersonName = request.FocalPersonName,
            FocalPersonContactJson = $"{{\"mobile\": \"{request.FocalPersonMobile}\"}}",
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        beneficiary.BusinessDetails = businessDetails;

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
            null,
            new BusinessBeneficiaryDto(
                businessDetails.BusinessName,
                null,
                businessDetails.TradeLicenseNumber,
                null,
                null,
                businessDetails.TinNumber,
                null,
                businessDetails.BusinessType.ToString(),
                null,
                0,
                null,
                null,
                null,
                null,
                businessDetails.FocalPersonName,
                null,
                null,
                businessDetails.FocalPersonContactJson,
                null,
                null,
                0,
                0,
                0,
                0
            )
        ));
    }
}
