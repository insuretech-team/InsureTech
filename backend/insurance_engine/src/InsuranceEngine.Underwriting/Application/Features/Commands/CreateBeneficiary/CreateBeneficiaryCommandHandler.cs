using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Underwriting.Application.Interfaces;
using InsuranceEngine.Underwriting.Domain.Entities;
using InsuranceEngine.Underwriting.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.CreateBeneficiary;

public class CreateBeneficiaryCommandHandler : IRequestHandler<CreateBeneficiaryCommand, Result<Guid>>
{
    private readonly IBeneficiaryRepository _repository;

    public CreateBeneficiaryCommandHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<Guid>> Handle(CreateBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var code = await _repository.GetNextSequenceAsync();
        
        var beneficiary = new Beneficiary
        {
            Id = Guid.NewGuid(),
            Code = request.Code,
            UserId = request.UserId,
            Type = request.Type,
            Status = BeneficiaryStatus.PendingKyc,
            Name = request.Name,
            ContactNumber = request.ContactNumber,
            Email = request.Email,
            Address = request.Address,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        if (request.Type == BeneficiaryType.Individual && request.IndividualDetails != null)
        {
            beneficiary.IndividualDetails = new IndividualBeneficiary
            {
                Id = Guid.NewGuid(),
                BeneficiaryId = beneficiary.Id,
                FatherName = request.IndividualDetails.FatherName,
                MotherName = request.IndividualDetails.MotherName,
                DateOfBirth = request.IndividualDetails.DateOfBirth,
                Occupation = request.IndividualDetails.Occupation,
                MonthlyIncome = (decimal)request.IndividualDetails.MonthlyIncome
            };
        }
        else if (request.Type == BeneficiaryType.Business && request.BusinessDetails != null)
        {
            beneficiary.BusinessDetails = new BusinessBeneficiary
            {
                Id = Guid.NewGuid(),
                BeneficiaryId = beneficiary.Id,
                RegistrationNumber = request.BusinessDetails.RegistrationNumber,
                Industry = request.BusinessDetails.Industry,
                FocalPersonName = request.BusinessDetails.FocalPersonName,
                FocalPersonContact = request.BusinessDetails.FocalPersonContact
            };
        }

        await _repository.AddAsync(beneficiary);
        return Result<Guid>.Ok(beneficiary.Id);
    }
}
