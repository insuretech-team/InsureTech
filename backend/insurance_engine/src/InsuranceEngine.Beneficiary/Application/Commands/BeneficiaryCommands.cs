using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public record CreateIndividualBeneficiaryCommand(
    string UserId,
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string NidNumber,
    string MobileNumber,
    string? Email,
    string? PartnerId
) : IRequest<CreateIndividualBeneficiaryResponse>;

public record CreateBusinessBeneficiaryCommand(
    string UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    string? PartnerId
) : IRequest<CreateBusinessBeneficiaryResponse>;

public record UpdateBeneficiaryCommand(
    string BeneficiaryId,
    string? MobileNumber,
    string? Email,
    string? Address
) : IRequest<UpdateBeneficiaryResponse>;

public record CompleteKYCCommand(
    string BeneficiaryId,
    string NidFrontUrl,
    string NidBackUrl,
    string SelfieUrl,
    string? PorichoyVerificationId
) : IRequest<CompleteKYCResponse>;

public record UpdateRiskScoreCommand(
    string BeneficiaryId,
    string RiskScore,
    string Reason
) : IRequest<UpdateRiskScoreResponse>;
