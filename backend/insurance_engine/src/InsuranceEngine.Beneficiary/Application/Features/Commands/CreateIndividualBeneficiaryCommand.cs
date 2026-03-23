using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.Beneficiary.Domain.Entities;
using InsuranceEngine.Beneficiary.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Commands;

public record CreateIndividualBeneficiaryCommand(
    Guid UserId,
    string FullName,
    DateTime DateOfBirth,
    string Gender,
    string NidNumber,
    string MobileNumber,
    string? Email = null,
    Guid? PartnerId = null
) : IRequest<Result<BeneficiaryDto>>;

public class CreateIndividualBeneficiaryCommandHandler : IRequestHandler<CreateIndividualBeneficiaryCommand, Result<BeneficiaryDto>>
{
    private readonly IBeneficiaryRepository _repository;
    private readonly IEncryptionService _encryptionService;

    public CreateIndividualBeneficiaryCommandHandler(IBeneficiaryRepository repository, IEncryptionService encryptionService)
    {
        _repository = repository;
        _encryptionService = encryptionService;
    }

    public async Task<Result<BeneficiaryDto>> Handle(CreateIndividualBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var gender = Enum.Parse<BeneficiaryGender>(request.Gender, true);
        
        var beneficiary = Domain.Entities.Beneficiary.CreateIndividual(
            request.UserId,
            request.FullName,
            request.DateOfBirth,
            gender,
            request.MobileNumber,
            request.Email);

        if (!string.IsNullOrEmpty(request.NidNumber) && beneficiary.Individual != null)
        {
            beneficiary.Individual.NidNumber = await _encryptionService.EncryptAsync(request.NidNumber);
        }

        await _repository.AddAsync(beneficiary);

        return Result.Ok(beneficiary.ToDto());
    }
}
