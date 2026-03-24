using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed record CreateIndividualBeneficiaryCommand(
    string UserId,
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string NidNumber,
    string MobileNumber,
    string? Email = null,
    string? PartnerId = null) : ICommand<string>;
