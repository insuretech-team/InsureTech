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
    Guid PartnerId,
    string? BusinessName = null,
    string? TradeLicenseNumber = null,
    string? TinNumber = null,
    string? FocalPersonName = null,
    string? FocalPersonMobile = null
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
        var beneficiary = Domain.Entities.Beneficiary.CreateBusinessEmpty(request.UserId, request.PartnerId);

        if (beneficiary.Business != null)
        {
            if (!string.IsNullOrEmpty(request.BusinessName))
                beneficiary.Business.BusinessName = request.BusinessName;
            if (!string.IsNullOrEmpty(request.TradeLicenseNumber))
                beneficiary.Business.TradeLicenseNumber = request.TradeLicenseNumber;
            if (!string.IsNullOrEmpty(request.TinNumber))
                beneficiary.Business.TinNumber = request.TinNumber;
            if (!string.IsNullOrEmpty(request.FocalPersonName))
                beneficiary.Business.FocalPersonName = request.FocalPersonName;
            if (!string.IsNullOrEmpty(request.FocalPersonMobile))
                beneficiary.Business.FocalPersonContact = new SharedKernel.Domain.ValueObjects.ContactInfo(request.FocalPersonMobile);
        }

        await _repository.AddAsync(beneficiary);

        return Result.Ok(beneficiary.ToDto());
    }

}
