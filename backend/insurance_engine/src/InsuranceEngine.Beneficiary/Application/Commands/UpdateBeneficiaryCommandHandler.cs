using MediatR;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class UpdateBeneficiaryCommandHandler : IRequestHandler<UpdateBeneficiaryCommand, UpdateBeneficiaryResponse>
{
    private readonly IRepository<BeneficiaryEntity> _beneficiaryRepo;
    private readonly IRepository<IndividualBeneficiaryEntity> _individualRepo;
    private readonly IRepository<BusinessBeneficiaryEntity> _businessRepo;

    public UpdateBeneficiaryCommandHandler(
        IRepository<BeneficiaryEntity> beneficiaryRepo,
        IRepository<IndividualBeneficiaryEntity> individualRepo,
        IRepository<BusinessBeneficiaryEntity> businessRepo)
    {
        _beneficiaryRepo = beneficiaryRepo;
        _individualRepo = individualRepo;
        _businessRepo = businessRepo;
    }

    public async Task<UpdateBeneficiaryResponse> Handle(UpdateBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var e = await _beneficiaryRepo.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
        if (e == null) return new UpdateBeneficiaryResponse { Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

        e.UpdatedAt = DateTime.UtcNow;

        if (e.Type == "INDIVIDUAL")
        {
            var individual = (await _individualRepo.FindAsync(x => x.BeneficiaryId == e.BeneficiaryId, cancellationToken)).FirstOrDefault();
            if (individual != null)
            {
                if (!string.IsNullOrEmpty(request.MobileNumber) || !string.IsNullOrEmpty(request.Email))
                {
                    individual.ContactInfo = $"{{\"mobile\":\"{request.MobileNumber ?? ""}\",\"email\":\"{request.Email ?? ""}\"}}";
                }
                if (!string.IsNullOrEmpty(request.Address))
                {
                    individual.PresentAddress = $"{{\"address\":\"{request.Address}\"}}";
                }
                individual.UpdatedAt = DateTime.UtcNow;
                await _individualRepo.UpdateAsync(individual, cancellationToken);
            }
        }
        else if (e.Type == "BUSINESS")
        {
            var business = (await _businessRepo.FindAsync(x => x.ParentBeneficiaryId == e.BeneficiaryId, cancellationToken)).FirstOrDefault();
            if (business != null)
            {
                if (!string.IsNullOrEmpty(request.MobileNumber))
                {
                    business.ContactInfo = $"{{\"mobile\":\"{request.MobileNumber}\"}}";
                }
                if (!string.IsNullOrEmpty(request.Address))
                {
                    business.BusinessAddress = $"{{\"address\":\"{request.Address}\"}}";
                }
                business.UpdatedAt = DateTime.UtcNow;
                await _businessRepo.UpdateAsync(business, cancellationToken);
            }
        }

        await _beneficiaryRepo.UpdateAsync(e, cancellationToken);
        return new UpdateBeneficiaryResponse { Message = "Beneficiary updated successfully" };
    }
}
