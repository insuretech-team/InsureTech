using InsuranceEngine.Api.RequestModels;
using InsuranceEngine.Api.ResponseModels;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Domain.Entities;
using ApiError = InsuranceEngine.Api.DTOs.Error;
using InsuranceEngine.Domain.Enums;

namespace InsuranceEngine.Application.Services.Beneficiaries;

public class BeneficiaryService(IBeneficiaryRepository beneficiaryRepository) : IBeneficiaryService
{
    private readonly IBeneficiaryRepository _beneficiaryRepository = beneficiaryRepository;

    public async Task<ListBeneficiariesResponseV1> GetBeneficiariesAsync(
        int page,
        int pageSize,
        CancellationToken cancellationToken)
    {
        (IReadOnlyList<Beneficiary> beneficiaries, int totalCount) =
            await _beneficiaryRepository.GetBeneficiariesPageAsync(page, pageSize, cancellationToken);

        return new ListBeneficiariesResponseV1
        {
            Beneficiaries = beneficiaries.Select(MapToSummary).ToList(),
            TotalCount = totalCount
        };
    }

    public async Task<(GetBeneficiaryResponseV1 Response, bool NotFound)> GetBeneficiaryByIdAsync(
        string beneficiaryId,
        CancellationToken cancellationToken)
    {
        if (!Guid.TryParse(beneficiaryId, out Guid parsedId))
        {
            return (new GetBeneficiaryResponseV1
            {
                Error = BuildError("INVALID_BENEFICIARY_ID", "beneficiary_id must be a valid GUID.", 400)
            }, false);
        }

        Beneficiary? beneficiary = await _beneficiaryRepository.GetBeneficiaryByIdAsync(parsedId, cancellationToken);
        if (beneficiary is null)
        {
            return (new GetBeneficiaryResponseV1
            {
                Error = BuildError("BENEFICIARY_NOT_FOUND", "Resource does not exist.", 404)
            }, true);
        }

        return (MapToDetailResponse(beneficiary), false);
    }

    public async Task<CreateBeneficiaryResponseV1> CreateIndividualBeneficiaryAsync(
        CreateIndividualBeneficiaryRequestV1 request,
        CancellationToken cancellationToken)
    {
        Beneficiary beneficiary = BuildIndividualBeneficiary(request);
        await _beneficiaryRepository.AddBeneficiaryAsync(beneficiary, cancellationToken);
        await _beneficiaryRepository.SaveChangesAsync(cancellationToken);

        return new CreateBeneficiaryResponseV1
        {
            BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
            BeneficiaryCode = beneficiary.BeneficiaryCode,
            Message = "Beneficiary created successfully for Bangladesh coverage."
        };
    }

    public async Task<CreateBeneficiaryResponseV1> CreateBusinessBeneficiaryAsync(
        CreateBusinessBeneficiaryRequestV1 request,
        CancellationToken cancellationToken)
    {
        Beneficiary beneficiary = BuildBusinessBeneficiary(request);
        await _beneficiaryRepository.AddBeneficiaryAsync(beneficiary, cancellationToken);
        await _beneficiaryRepository.SaveChangesAsync(cancellationToken);

        return new CreateBeneficiaryResponseV1
        {
            BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
            BeneficiaryCode = beneficiary.BeneficiaryCode,
            Message = "Beneficiary created successfully for Bangladesh coverage."
        };
    }

    public async Task<(UpdateBeneficiaryResponseV1 Response, bool NotFound)> UpdateBeneficiaryAsync(
        string beneficiaryId,
        UpdateBeneficiaryRequestV1 request,
        CancellationToken cancellationToken)
    {
        if (!Guid.TryParse(beneficiaryId, out Guid parsedId))
        {
            return (new UpdateBeneficiaryResponseV1
            {
                Error = BuildError("INVALID_BENEFICIARY_ID", "beneficiary_id must be a valid GUID.", 400)
            }, false);
        }

        Beneficiary? beneficiary = await _beneficiaryRepository.GetTrackedBeneficiaryByIdAsync(parsedId, cancellationToken);

        if (beneficiary is null)
        {
            return (new UpdateBeneficiaryResponseV1
            {
                Error = BuildError("BENEFICIARY_NOT_FOUND", "Resource does not exist.", 404)
            }, true);
        }

        ApplyUpdate(beneficiary, request);
        await _beneficiaryRepository.SaveChangesAsync(cancellationToken);

        return (new UpdateBeneficiaryResponseV1
        {
            Message = "Beneficiary updated successfully."
        }, false);
    }

    private static Beneficiary BuildIndividualBeneficiary(CreateIndividualBeneficiaryRequestV1 request)
    {
        DateTime now = DateTime.UtcNow;
        string? userId = request.UserId?.Trim();
        string? partnerId = request.PartnerId?.Trim();
        Beneficiary beneficiary = new()
        {
            BeneficiaryId = Guid.NewGuid(),
            UserId = userId ?? string.Empty,
            PartnerId = partnerId,
            BeneficiaryCode = GenerateBeneficiaryCode(),
            Type = BeneficiaryType.BeneficiaryTypeIndividual,
            Status = BeneficiaryStatus.BeneficiaryStatusPendingKyc,
            KycStatus = "KYC_STATUS_UNSPECIFIED",
            AuditInfo = BuildAuditInfo(request.AuditInfo, userId ?? string.Empty, now),
            IndividualDetails = new BeneficiaryIndividual
            {
                Id = Guid.NewGuid(),
                FullName = request.FullName?.Trim() ?? string.Empty,
                FullNameBn = request.FullNameBn?.Trim(),
                DateOfBirth = request.DateOfBirth ?? now.Date,
                Gender = request.Gender ?? BeneficiaryGender.GenderUnspecified,
                NidNumber = request.NidNumber?.Trim(),
                PassportNumber = request.PassportNumber?.Trim(),
                BirthCertificateNumber = request.BirthCertificateNumber?.Trim(),
                TinNumber = request.TinNumber?.Trim(),
                MaritalStatus = request.MaritalStatus,
                Occupation = request.Occupation?.Trim(),
                ContactInfo = BuildContactInfo(request.ContactInfo),
                PermanentAddress = BuildAddress(request.PermanentAddress),
                PresentAddress = BuildAddress(request.PresentAddress),
                NomineeName = request.NomineeName?.Trim(),
                NomineeRelationship = request.NomineeRelationship?.Trim(),
                AuditInfo = BuildAuditInfo(request.AuditInfo, userId ?? string.Empty, now)
            }
        };

        beneficiary.IndividualDetails.BeneficiaryId = beneficiary.BeneficiaryId;
        return beneficiary;
    }

    private static Beneficiary BuildBusinessBeneficiary(CreateBusinessBeneficiaryRequestV1 request)
    {
        DateTime now = DateTime.UtcNow;
        string? userId = request.UserId?.Trim();
        string? partnerId = request.PartnerId?.Trim();
        Beneficiary beneficiary = new()
        {
            BeneficiaryId = Guid.NewGuid(),
            UserId = userId ?? string.Empty,
            PartnerId = partnerId,
            BeneficiaryCode = GenerateBeneficiaryCode(),
            Type = BeneficiaryType.BeneficiaryTypeBusiness,
            Status = BeneficiaryStatus.BeneficiaryStatusPendingKyc,
            KycStatus = "KYC_STATUS_UNSPECIFIED",
            AuditInfo = BuildAuditInfo(request.AuditInfo, userId ?? string.Empty, now),
            BusinessDetails = new BeneficiaryBusiness
            {
                Id = Guid.NewGuid(),
                BusinessName = request.BusinessName?.Trim() ?? string.Empty,
                BusinessNameBn = request.BusinessNameBn?.Trim(),
                TradeLicenseNumber = request.TradeLicenseNumber?.Trim() ?? string.Empty,
                TradeLicenseIssueDate = request.TradeLicenseIssueDate,
                TradeLicenseExpiryDate = request.TradeLicenseExpiryDate,
                TinNumber = request.TinNumber?.Trim() ?? string.Empty,
                BinNumber = request.BinNumber?.Trim(),
                BusinessType = request.BusinessType ?? BusinessType.BusinessTypeUnspecified,
                IndustrySector = request.IndustrySector?.Trim(),
                EmployeeCount = request.EmployeeCount,
                IncorporationDate = request.IncorporationDate,
                ContactInfo = BuildContactInfo(request.ContactInfo),
                RegisteredAddress = BuildAddress(request.RegisteredAddress),
                BusinessAddress = BuildAddress(request.BusinessAddress),
                FocalPersonName = request.FocalPersonName?.Trim() ?? string.Empty,
                FocalPersonDesignation = request.FocalPersonDesignation?.Trim(),
                FocalPersonNid = request.FocalPersonNid?.Trim(),
                FocalPersonContact = BuildContactInfo(request.FocalPersonContact),
                AuditInfo = BuildAuditInfo(request.AuditInfo, userId ?? string.Empty, now)
            }
        };

        beneficiary.BusinessDetails.BeneficiaryId = beneficiary.BeneficiaryId;
        return beneficiary;
    }

    private static void ApplyUpdate(Beneficiary beneficiary, UpdateBeneficiaryRequestV1 request)
    {
        DateTime now = DateTime.UtcNow;

        if (!string.IsNullOrWhiteSpace(request.KycStatus))
        {
            beneficiary.KycStatus = request.KycStatus.Trim();
        }

        if (request.KycCompletedAt.HasValue)
        {
            beneficiary.KycCompletedAt = request.KycCompletedAt;
        }

        if (!string.IsNullOrWhiteSpace(request.RiskScore))
        {
            beneficiary.RiskScore = request.RiskScore.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.ReferralCode))
        {
            beneficiary.ReferralCode = request.ReferralCode.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.ReferredBy))
        {
            beneficiary.ReferredBy = request.ReferredBy.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.PartnerId))
        {
            beneficiary.PartnerId = request.PartnerId.Trim();
        }

        if (beneficiary.Type == BeneficiaryType.BeneficiaryTypeIndividual && beneficiary.IndividualDetails is not null)
        {
            ApplyIndividualUpdate(beneficiary.IndividualDetails, request.IndividualDetails);
            UpdateAuditInfo(beneficiary.IndividualDetails.AuditInfo, request.AuditInfo, now, beneficiary.UserId);
        }

        if (beneficiary.Type == BeneficiaryType.BeneficiaryTypeBusiness && beneficiary.BusinessDetails is not null)
        {
            ApplyBusinessUpdate(beneficiary.BusinessDetails, request.BusinessDetails);
            UpdateAuditInfo(beneficiary.BusinessDetails.AuditInfo, request.AuditInfo, now, beneficiary.UserId);
        }

        UpdateAuditInfo(beneficiary.AuditInfo, request.AuditInfo, now, beneficiary.UserId);
    }

    private static void ApplyIndividualUpdate(BeneficiaryIndividual individual, IndividualBeneficiaryUpdateV1? request)
    {
        if (request is null)
        {
            return;
        }

        if (!string.IsNullOrWhiteSpace(request.FullName))
        {
            individual.FullName = request.FullName.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.FullNameBn))
        {
            individual.FullNameBn = request.FullNameBn.Trim();
        }

        if (request.DateOfBirth.HasValue)
        {
            individual.DateOfBirth = request.DateOfBirth.Value;
        }

        if (request.Gender.HasValue)
        {
            individual.Gender = request.Gender.Value;
        }

        if (!string.IsNullOrWhiteSpace(request.NidNumber))
        {
            individual.NidNumber = request.NidNumber.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.PassportNumber))
        {
            individual.PassportNumber = request.PassportNumber.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.BirthCertificateNumber))
        {
            individual.BirthCertificateNumber = request.BirthCertificateNumber.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.TinNumber))
        {
            individual.TinNumber = request.TinNumber.Trim();
        }

        if (request.MaritalStatus.HasValue)
        {
            individual.MaritalStatus = request.MaritalStatus.Value;
        }

        if (!string.IsNullOrWhiteSpace(request.Occupation))
        {
            individual.Occupation = request.Occupation.Trim();
        }

        if (request.ContactInfo is not null)
        {
            individual.ContactInfo ??= new ContactInfo();
            UpdateContactInfo(individual.ContactInfo, request.ContactInfo);
        }

        if (request.PermanentAddress is not null)
        {
            individual.PermanentAddress ??= new Address();
            UpdateAddress(individual.PermanentAddress, request.PermanentAddress);
        }

        if (request.PresentAddress is not null)
        {
            individual.PresentAddress ??= new Address();
            UpdateAddress(individual.PresentAddress, request.PresentAddress);
        }

        if (!string.IsNullOrWhiteSpace(request.NomineeName))
        {
            individual.NomineeName = request.NomineeName.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.NomineeRelationship))
        {
            individual.NomineeRelationship = request.NomineeRelationship.Trim();
        }
    }

    private static void ApplyBusinessUpdate(BeneficiaryBusiness business, BusinessBeneficiaryUpdateV1? request)
    {
        if (request is null)
        {
            return;
        }

        if (!string.IsNullOrWhiteSpace(request.BusinessName))
        {
            business.BusinessName = request.BusinessName.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.BusinessNameBn))
        {
            business.BusinessNameBn = request.BusinessNameBn.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.TradeLicenseNumber))
        {
            business.TradeLicenseNumber = request.TradeLicenseNumber.Trim();
        }

        if (request.TradeLicenseIssueDate.HasValue)
        {
            business.TradeLicenseIssueDate = request.TradeLicenseIssueDate;
        }

        if (request.TradeLicenseExpiryDate.HasValue)
        {
            business.TradeLicenseExpiryDate = request.TradeLicenseExpiryDate;
        }

        if (!string.IsNullOrWhiteSpace(request.TinNumber))
        {
            business.TinNumber = request.TinNumber.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.BinNumber))
        {
            business.BinNumber = request.BinNumber.Trim();
        }

        if (request.BusinessType.HasValue)
        {
            business.BusinessType = request.BusinessType.Value;
        }

        if (!string.IsNullOrWhiteSpace(request.IndustrySector))
        {
            business.IndustrySector = request.IndustrySector.Trim();
        }

        if (request.EmployeeCount.HasValue)
        {
            business.EmployeeCount = request.EmployeeCount;
        }

        if (request.IncorporationDate.HasValue)
        {
            business.IncorporationDate = request.IncorporationDate;
        }

        if (request.ContactInfo is not null)
        {
            business.ContactInfo ??= new ContactInfo();
            UpdateContactInfo(business.ContactInfo, request.ContactInfo);
        }

        if (request.RegisteredAddress is not null)
        {
            business.RegisteredAddress ??= new Address();
            UpdateAddress(business.RegisteredAddress, request.RegisteredAddress);
        }

        if (request.BusinessAddress is not null)
        {
            business.BusinessAddress ??= new Address();
            UpdateAddress(business.BusinessAddress, request.BusinessAddress);
        }

        if (!string.IsNullOrWhiteSpace(request.FocalPersonName))
        {
            business.FocalPersonName = request.FocalPersonName.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.FocalPersonDesignation))
        {
            business.FocalPersonDesignation = request.FocalPersonDesignation.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.FocalPersonNid))
        {
            business.FocalPersonNid = request.FocalPersonNid.Trim();
        }

        if (request.FocalPersonContact is not null)
        {
            business.FocalPersonContact ??= new ContactInfo();
            UpdateContactInfo(business.FocalPersonContact, request.FocalPersonContact);
        }
    }

    private static void UpdateContactInfo(ContactInfo contactInfo, ContactInfoUpdateRequestV1 request)
    {
        if (!string.IsNullOrWhiteSpace(request.MobileNumber))
        {
            contactInfo.MobileNumber = request.MobileNumber.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.Email))
        {
            contactInfo.Email = request.Email.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.AlternateMobile))
        {
            contactInfo.AlternateMobile = request.AlternateMobile.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.Landline))
        {
            contactInfo.Landline = request.Landline.Trim();
        }
    }

    private static void UpdateAddress(Address address, AddressUpdateRequestV1 request)
    {
        if (!string.IsNullOrWhiteSpace(request.AddressLine1))
        {
            address.AddressLine1 = request.AddressLine1.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.AddressLine2))
        {
            address.AddressLine2 = request.AddressLine2.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.City))
        {
            address.City = request.City.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.District))
        {
            address.District = request.District.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.Division))
        {
            address.Division = request.Division.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.PostalCode))
        {
            address.PostalCode = request.PostalCode.Trim();
        }

        if (!string.IsNullOrWhiteSpace(request.Country))
        {
            address.Country = request.Country.Trim();
        }

        if (request.Latitude.HasValue)
        {
            address.Latitude = request.Latitude.Value;
        }

        if (request.Longitude.HasValue)
        {
            address.Longitude = request.Longitude.Value;
        }
    }

    private static BeneficiarySummaryV1 MapToSummary(Beneficiary beneficiary)
    {
        return new BeneficiarySummaryV1
        {
            BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
            BeneficiaryCode = beneficiary.BeneficiaryCode,
            UserId = beneficiary.UserId,
            PartnerId = beneficiary.PartnerId,
            Type = beneficiary.Type,
            Status = beneficiary.Status
        };
    }

    private static GetBeneficiaryResponseV1 MapToDetailResponse(Beneficiary beneficiary)
    {
        GetBeneficiaryResponseV1 response = new()
        {
            Beneficiary = new BeneficiaryProfileV1
            {
                BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
                UserId = beneficiary.UserId,
                Type = beneficiary.Type,
                Code = beneficiary.BeneficiaryCode,
                Status = beneficiary.Status,
                KycStatus = beneficiary.KycStatus,
                KycCompletedAt = beneficiary.KycCompletedAt,
                RiskScore = beneficiary.RiskScore,
                ReferralCode = beneficiary.ReferralCode,
                ReferredBy = beneficiary.ReferredBy,
                PartnerId = beneficiary.PartnerId,
                AuditInfo = BuildAuditInfoResponse(beneficiary.AuditInfo)
            }
        };

        if (beneficiary.Type == BeneficiaryType.BeneficiaryTypeIndividual && beneficiary.IndividualDetails is not null)
        {
            response.IndividualDetails = new IndividualDetailsV1
            {
                Id = beneficiary.IndividualDetails.Id.ToString(),
                BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
                FullName = beneficiary.IndividualDetails.FullName,
                FullNameBn = beneficiary.IndividualDetails.FullNameBn,
                DateOfBirth = beneficiary.IndividualDetails.DateOfBirth,
                Gender = beneficiary.IndividualDetails.Gender,
                NidNumber = beneficiary.IndividualDetails.NidNumber,
                PassportNumber = beneficiary.IndividualDetails.PassportNumber,
                BirthCertificateNumber = beneficiary.IndividualDetails.BirthCertificateNumber,
                TinNumber = beneficiary.IndividualDetails.TinNumber,
                MaritalStatus = beneficiary.IndividualDetails.MaritalStatus,
                Occupation = beneficiary.IndividualDetails.Occupation,
                ContactInfo = BuildContactInfoResponse(beneficiary.IndividualDetails.ContactInfo),
                PermanentAddress = BuildAddressInfoResponse(beneficiary.IndividualDetails.PermanentAddress),
                PresentAddress = BuildAddressInfoResponse(beneficiary.IndividualDetails.PresentAddress),
                NomineeName = beneficiary.IndividualDetails.NomineeName,
                NomineeRelationship = beneficiary.IndividualDetails.NomineeRelationship,
                AuditInfo = BuildAuditInfoResponse(beneficiary.IndividualDetails.AuditInfo)
            };
        }

        if (beneficiary.Type == BeneficiaryType.BeneficiaryTypeBusiness && beneficiary.BusinessDetails is not null)
        {
            response.BusinessDetails = new BusinessDetailsV1
            {
                Id = beneficiary.BusinessDetails.Id.ToString(),
                BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
                BusinessName = beneficiary.BusinessDetails.BusinessName,
                BusinessNameBn = beneficiary.BusinessDetails.BusinessNameBn,
                TradeLicenseNumber = beneficiary.BusinessDetails.TradeLicenseNumber,
                TradeLicenseIssueDate = beneficiary.BusinessDetails.TradeLicenseIssueDate,
                TradeLicenseExpiryDate = beneficiary.BusinessDetails.TradeLicenseExpiryDate,
                TinNumber = beneficiary.BusinessDetails.TinNumber,
                BinNumber = beneficiary.BusinessDetails.BinNumber,
                BusinessType = beneficiary.BusinessDetails.BusinessType,
                IndustrySector = beneficiary.BusinessDetails.IndustrySector,
                EmployeeCount = beneficiary.BusinessDetails.EmployeeCount,
                IncorporationDate = beneficiary.BusinessDetails.IncorporationDate,
                ContactInfo = BuildContactInfoResponse(beneficiary.BusinessDetails.ContactInfo),
                RegisteredAddress = BuildAddressInfoResponse(beneficiary.BusinessDetails.RegisteredAddress),
                BusinessAddress = BuildAddressInfoResponse(beneficiary.BusinessDetails.BusinessAddress),
                FocalPersonName = beneficiary.BusinessDetails.FocalPersonName,
                FocalPersonDesignation = beneficiary.BusinessDetails.FocalPersonDesignation,
                FocalPersonNid = beneficiary.BusinessDetails.FocalPersonNid,
                FocalPersonContact = BuildContactInfoResponse(beneficiary.BusinessDetails.FocalPersonContact),
                AuditInfo = BuildAuditInfoResponse(beneficiary.BusinessDetails.AuditInfo)
            };
        }

        return response;
    }

    private static ApiError BuildError(string code, string message, int statusCode)
    {
        return ApiError.Create(code, message, statusCode);
    }

    private static AuditInfo BuildAuditInfo(AuditInfoRequestV1? request, string defaultUserId, DateTime now)
    {
        return new AuditInfo
        {
            CreatedAt = request?.CreatedAt ?? now,
            UpdatedAt = request?.UpdatedAt ?? now,
            CreatedBy = request?.CreatedBy?.Trim() ?? defaultUserId,
            UpdatedBy = request?.UpdatedBy?.Trim() ?? defaultUserId,
            DeletedAt = request?.DeletedAt,
            DeletedBy = request?.DeletedBy?.Trim()
        };
    }

    private static void UpdateAuditInfo(AuditInfo auditInfo, AuditInfoRequestV1? request, DateTime now, string defaultUserId)
    {
        if (request is not null)
        {
            if (request.CreatedAt.HasValue)
            {
                auditInfo.CreatedAt = request.CreatedAt.Value;
            }

            if (request.UpdatedAt.HasValue)
            {
                auditInfo.UpdatedAt = request.UpdatedAt.Value;
            }

            if (!string.IsNullOrWhiteSpace(request.CreatedBy))
            {
                auditInfo.CreatedBy = request.CreatedBy.Trim();
            }

            if (!string.IsNullOrWhiteSpace(request.UpdatedBy))
            {
                auditInfo.UpdatedBy = request.UpdatedBy.Trim();
            }

            auditInfo.DeletedAt = request.DeletedAt;
            auditInfo.DeletedBy = request.DeletedBy?.Trim();

            return;
        }

        auditInfo.UpdatedAt = now;
        auditInfo.UpdatedBy = defaultUserId;
    }

    private static ContactInfo? BuildContactInfo(ContactInfoRequestV1? request)
    {
        if (request is null)
        {
            return null;
        }

        return new ContactInfo
        {
            MobileNumber = request.MobileNumber?.Trim(),
            Email = request.Email?.Trim(),
            AlternateMobile = request.AlternateMobile?.Trim(),
            Landline = request.Landline?.Trim()
        };
    }

    private static Address? BuildAddress(AddressRequestV1? request)
    {
        if (request is null)
        {
            return null;
        }

        return new Address
        {
            AddressLine1 = request.AddressLine1?.Trim(),
            AddressLine2 = request.AddressLine2?.Trim(),
            City = request.City?.Trim(),
            District = request.District?.Trim(),
            Division = request.Division?.Trim(),
            PostalCode = request.PostalCode?.Trim(),
            Country = string.IsNullOrWhiteSpace(request.Country) ? "Bangladesh" : request.Country.Trim(),
            Latitude = request.Latitude,
            Longitude = request.Longitude
        };
    }

    private static AuditInfoV1? BuildAuditInfoResponse(AuditInfo? auditInfo)
    {
        if (auditInfo is null)
        {
            return null;
        }

        return new AuditInfoV1
        {
            CreatedAt = auditInfo.CreatedAt,
            UpdatedAt = auditInfo.UpdatedAt,
            CreatedBy = auditInfo.CreatedBy,
            UpdatedBy = auditInfo.UpdatedBy,
            DeletedAt = auditInfo.DeletedAt,
            DeletedBy = auditInfo.DeletedBy
        };
    }

    private static ContactInfoV1? BuildContactInfoResponse(ContactInfo? contactInfo)
    {
        if (contactInfo is null)
        {
            return null;
        }

        return new ContactInfoV1
        {
            MobileNumber = contactInfo.MobileNumber,
            Email = contactInfo.Email,
            AlternateMobile = contactInfo.AlternateMobile,
            Landline = contactInfo.Landline
        };
    }

    private static AddressV1? BuildAddressInfoResponse(Address? address)
    {
        if (address is null)
        {
            return null;
        }

        return new AddressV1
        {
            AddressLine1 = address.AddressLine1,
            AddressLine2 = address.AddressLine2,
            City = address.City,
            District = address.District,
            Division = address.Division,
            PostalCode = address.PostalCode,
            Country = address.Country,
            Latitude = address.Latitude,
            Longitude = address.Longitude
        };
    }

    private static string GenerateBeneficiaryCode()
    {
        string randomSegment = Guid.NewGuid().ToString("N")[..8].ToUpperInvariant();
        return $"BD-BEN-{randomSegment}";
    }
}
