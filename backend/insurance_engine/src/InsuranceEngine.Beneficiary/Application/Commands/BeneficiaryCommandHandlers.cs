using MediatR;
using Microsoft.Extensions.Logging;
using Microsoft.EntityFrameworkCore;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using Google.Protobuf.WellKnownTypes;

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

public sealed class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, CompleteKYCResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repository;
    public CompleteKYCCommandHandler(IRepository<BeneficiaryEntity> repository) => _repository = repository;

    public async Task<CompleteKYCResponse> Handle(CompleteKYCCommand request, CancellationToken cancellationToken)
    {
        var e = await _repository.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
        if (e == null) return new CompleteKYCResponse { Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

        e.KycStatus = "COMPLETED";
        e.KycCompletedAt = DateTime.UtcNow;
        if (e.Status == "PENDING_KYC") e.Status = "ACTIVE";
        e.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(e, cancellationToken);
        return new CompleteKYCResponse { Message = "KYC completed successfully" };
    }
}

public sealed class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, UpdateRiskScoreResponse>
{
    private readonly IRepository<BeneficiaryEntity> _repository;
    public UpdateRiskScoreCommandHandler(IRepository<BeneficiaryEntity> repository) => _repository = repository;

    public async Task<UpdateRiskScoreResponse> Handle(UpdateRiskScoreCommand request, CancellationToken cancellationToken)
    {
        var e = await _repository.GetByIdAsync(Guid.Parse(request.BeneficiaryId), cancellationToken);
        if (e == null) return new UpdateRiskScoreResponse { Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" } };

        e.RiskScore = request.RiskScore;
        e.UpdatedAt = DateTime.UtcNow;

        await _repository.UpdateAsync(e, cancellationToken);
        return new UpdateRiskScoreResponse { Message = "Risk score updated" };
    }
}

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
