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

        return Result.Ok(beneficiary.ToDto());
    }

}
