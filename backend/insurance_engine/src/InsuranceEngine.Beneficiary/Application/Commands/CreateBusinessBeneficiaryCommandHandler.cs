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

public sealed class CreateBusinessBeneficiaryCommandHandler : IRequestHandler<CreateBusinessBeneficiaryCommand, CreateBusinessBeneficiaryResponse>
{
    private readonly IRepository<BeneficiaryEntity> _beneficiaryRepo;
    private readonly IRepository<BusinessBeneficiaryEntity> _businessRepo;
    private readonly ILogger<CreateBusinessBeneficiaryCommandHandler> _logger;

    public CreateBusinessBeneficiaryCommandHandler(
        IRepository<BeneficiaryEntity> beneficiaryRepo,
        IRepository<BusinessBeneficiaryEntity> businessRepo,
        ILogger<CreateBusinessBeneficiaryCommandHandler> logger)
    {
        _beneficiaryRepo = beneficiaryRepo;
        _businessRepo = businessRepo;
        _logger = logger;
    }

    public async Task<CreateBusinessBeneficiaryResponse> Handle(CreateBusinessBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var exists = await _businessRepo.ExistsAsync(x => x.TradeLicenseNumber == request.TradeLicenseNumber, cancellationToken);
            if (exists)
                return new CreateBusinessBeneficiaryResponse { Error = new Error { Code = "DUPLICATE_BUSINESS", Message = "Business with this Trade License already exists" } };

            var beneficiaryId = Guid.NewGuid();
            var code = $"BEN-B-{Guid.NewGuid().ToString()[..8].ToUpper()}";

            var beneficiary = new BeneficiaryEntity
            {
                BeneficiaryId = beneficiaryId,
                UserId = Guid.Parse(request.UserId),
                Type = "BUSINESS",
                Code = code,
                Status = "PENDING_KYC",
                KycStatus = "NOT_STARTED",
                PartnerId = string.IsNullOrEmpty(request.PartnerId) ? null : Guid.Parse(request.PartnerId),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            var business = new BusinessBeneficiaryEntity
            {
                BeneficiaryId = Guid.NewGuid(),
                ParentBeneficiaryId = beneficiaryId,
                BusinessName = request.BusinessName,
                TradeLicenseNumber = request.TradeLicenseNumber,
                TinNumber = request.TinNumber,
                FocalPersonName = request.FocalPersonName,
                FocalPersonContact = $"{{\"mobile\":\"{request.FocalPersonMobile}\"}}",
                BusinessType = "PRIVATE_LIMITED", // Default
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _beneficiaryRepo.AddAsync(beneficiary, cancellationToken);
            await _businessRepo.AddAsync(business, cancellationToken);

            return new CreateBusinessBeneficiaryResponse
            {
                BeneficiaryId = beneficiaryId.ToString(),
                BeneficiaryCode = code,
                Message = "Business beneficiary created successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create business beneficiary");
            return new CreateBusinessBeneficiaryResponse { Error = new Error { Code = "INTERNAL_ERROR", Message = ex.Message } };
        }
    }
}
