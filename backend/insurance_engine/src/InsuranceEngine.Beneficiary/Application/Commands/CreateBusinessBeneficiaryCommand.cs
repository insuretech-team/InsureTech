using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public record CreateBusinessBeneficiaryCommand(
    string UserId,
    string BusinessName,
    string TradeLicenseNumber,
    string TinNumber,
    string FocalPersonName,
    string FocalPersonMobile,
    string? PartnerId = null
) : IRequest<Result<string>>;
