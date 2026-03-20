using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.Underwriting.Domain.Enums;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.CreateBeneficiary;

public record CreateBeneficiaryCommand(
    string Name,
    BeneficiaryType Type,
    string ContactNumber,
    string Email,
    string Address,
    Guid UserId,
    string Code = null,
    string Nid = null,
    string Tin = null,
    IndividualBeneficiaryDto IndividualDetails = null,
    BusinessBeneficiaryDto BusinessDetails = null
) : IRequest<Result<Guid>>;

public record IndividualBeneficiaryDto(
    string FatherName,
    string MotherName,
    DateTime DateOfBirth,
    string Occupation,
    double MonthlyIncome
);

public record BusinessBeneficiaryDto(
    string RegistrationNumber,
    string Industry,
    string FocalPersonName,
    string FocalPersonContact
);
