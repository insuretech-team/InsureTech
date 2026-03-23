using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

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

    public CreateBusinessBeneficiaryCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<BeneficiaryDto>> Handle(CreateBusinessBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var beneficiary = Domain.Entities.Beneficiary.CreateBusiness(
            request.UserId,
            request.BusinessName,
            request.TradeLicenseNumber,
            request.TinNumber,
            request.FocalPersonName,
            request.FocalPersonMobile);

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
            null,
            b.BusinessDetails != null ? new BusinessBeneficiaryDto(
                b.BusinessDetails.BusinessName,
                b.BusinessDetails.BusinessNameBn,
                b.BusinessDetails.TradeLicenseNumber,
                b.BusinessDetails.TinNumber,
                null,
                b.BusinessDetails.BusinessType.ToString(),
                b.BusinessDetails.IndustrySector,
                b.BusinessDetails.FocalPersonName,
                b.BusinessDetails.FocalPersonMobile,
                null,
                null,
                null
            ) : null
        );
    }
}
