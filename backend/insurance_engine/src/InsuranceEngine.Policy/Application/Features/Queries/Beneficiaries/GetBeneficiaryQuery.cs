using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using InsuranceEngine.SharedKernel.Services;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.Beneficiaries;

public record GetBeneficiaryQuery(Guid BeneficiaryId) : IRequest<Result<BeneficiaryDto>>;

public class GetBeneficiaryQueryHandler : IRequestHandler<GetBeneficiaryQuery, Result<BeneficiaryDto>>
{
    private readonly IBeneficiaryRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public GetBeneficiaryQueryHandler(IBeneficiaryRepository repository, IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<Result<BeneficiaryDto>> Handle(GetBeneficiaryQuery request, CancellationToken cancellationToken)
    {
        var beneficiary = await _repository.GetByIdAsync(request.BeneficiaryId);
        if (beneficiary == null)
            return Result<BeneficiaryDto>.Fail(Error.NotFound("Beneficiary", request.BeneficiaryId.ToString()));

        IndividualBeneficiaryDto? individualDto = null;
        if (beneficiary.IndividualDetails != null)
        {
            var i = beneficiary.IndividualDetails;
            
            // Decrypt and mask PII
            var nid = !string.IsNullOrEmpty(i.NidNumber) ? PiiMasker.MaskNid(_encryptionService.Decrypt(i.NidNumber)) : null;
            var passport = !string.IsNullOrEmpty(i.PassportNumber) ? PiiMasker.MaskNid(_encryptionService.Decrypt(i.PassportNumber)) : null;
            var birthCert = !string.IsNullOrEmpty(i.BirthCertificateNumber) ? PiiMasker.MaskNid(_encryptionService.Decrypt(i.BirthCertificateNumber)) : null;

            individualDto = new IndividualBeneficiaryDto(
                i.FullName, i.FullNameBn, i.DateOfBirth, i.Gender.ToString(),
                nid, passport, birthCert, i.TinNumber,
                i.MaritalStatus.ToString(), i.Occupation, i.ContactInfoJson,
                i.PermanentAddressJson, i.PresentAddressJson, i.NomineeName, i.NomineeRelationship
            );
        }

        BusinessBeneficiaryDto? businessDto = null;
        if (beneficiary.BusinessDetails != null)
        {
            var b = beneficiary.BusinessDetails;
            businessDto = new BusinessBeneficiaryDto(
                b.BusinessName, b.BusinessNameBn, b.TradeLicenseNumber,
                b.TradeLicenseIssueDate, b.TradeLicenseExpiryDate, b.TinNumber,
                b.BinNumber, b.BusinessType.ToString(), b.IndustrySector,
                b.EmployeeCount, b.IncorporationDate, b.ContactInfoJson,
                b.RegisteredAddressJson, b.BusinessAddressJson, b.FocalPersonName,
                b.FocalPersonDesignation, b.FocalPersonNid, b.FocalPersonContactJson,
                b.RegistrationNumber, b.TaxId, b.TotalEmployeesCovered,
                b.ActivePoliciesCount, b.TotalPremiumAmount, b.PendingActionsCount
            );
        }

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
            individualDto,
            businessDto
        ));
    }
}
