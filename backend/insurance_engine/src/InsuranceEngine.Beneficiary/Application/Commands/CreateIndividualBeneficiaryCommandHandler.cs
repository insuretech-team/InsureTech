using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CreateIndividualBeneficiaryCommandHandler : IRequestHandler<CreateIndividualBeneficiaryCommand, CreateIndividualBeneficiaryResponse>
{
    private readonly IRepository<BeneficiaryEntity> _beneficiaryRepo;
    private readonly IRepository<IndividualBeneficiaryEntity> _individualRepo;
    private readonly ILogger<CreateIndividualBeneficiaryCommandHandler> _logger;

    public CreateIndividualBeneficiaryCommandHandler(
        IRepository<BeneficiaryEntity> beneficiaryRepo,
        IRepository<IndividualBeneficiaryEntity> individualRepo,
        ILogger<CreateIndividualBeneficiaryCommandHandler> logger)
    {
        _beneficiaryRepo = beneficiaryRepo;
        _individualRepo = individualRepo;
        _logger = logger;
    }

    public async Task<CreateIndividualBeneficiaryResponse> Handle(CreateIndividualBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var exists = await _individualRepo.ExistsAsync(x => x.NidNumber == request.NidNumber, cancellationToken);
            if (exists)
                return new CreateIndividualBeneficiaryResponse { Error = new Error { Code = "DUPLICATE_BENEFICIARY", Message = "Beneficiary with this NID already exists" } };

            var beneficiaryId = Guid.NewGuid();
            var code = $"BEN-I-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            var beneficiary = new BeneficiaryEntity
            {
                BeneficiaryId = beneficiaryId,
                UserId = Guid.Parse(request.UserId),
                Type = "INDIVIDUAL",
                Code = code,
                Status = "PENDING_KYC",
                KycStatus = "NOT_STARTED",
                PartnerId = string.IsNullOrEmpty(request.PartnerId) ? null : Guid.Parse(request.PartnerId),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            var individual = new IndividualBeneficiaryEntity
            {
                BeneficiaryId = beneficiaryId,
                FullName = request.FullName,
                DateOfBirth = DateTime.SpecifyKind(request.DateOfBirth, DateTimeKind.Utc),
                Gender = request.Gender,
                NidNumber = request.NidNumber,
                ContactInfo = $"{{\"mobile\":\"{request.MobileNumber}\",\"email\":\"{request.Email ?? ""}\"}}",
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _beneficiaryRepo.AddAsync(beneficiary, cancellationToken);
            await _individualRepo.AddAsync(individual, cancellationToken);

            return new CreateIndividualBeneficiaryResponse
            {
                BeneficiaryId = beneficiaryId.ToString(),
                BeneficiaryCode = code,
                Message = "Individual beneficiary created successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create individual beneficiary");
            return new CreateIndividualBeneficiaryResponse { Error = new Error { Code = "INTERNAL_ERROR", Message = ex.Message } };
        }
    }
}
