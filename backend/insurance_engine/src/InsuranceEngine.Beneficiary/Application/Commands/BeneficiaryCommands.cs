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
    string? Status,
    string? RiskScore,
    string? FocalPersonName,
    string? FocalPersonMobile
) : IRequest<UpdateBeneficiaryResponse>;

public record CompleteKYCCommand(
    string BeneficiaryId,
    string IdType,
    string IdNumber,
    string IdUrl
) : IRequest<CompleteKYCResponse>;

public record UpdateRiskScoreCommand(
    string BeneficiaryId,
    string RiskScore,
    string Reason
) : IRequest<UpdateRiskScoreResponse>;
