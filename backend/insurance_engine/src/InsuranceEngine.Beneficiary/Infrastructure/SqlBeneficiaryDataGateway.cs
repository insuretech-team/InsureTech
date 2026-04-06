using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence.Entities;

namespace InsuranceEngine.Beneficiary.Infrastructure;

public class SqlBeneficiaryDataGateway : IBeneficiaryDataGateway
{
    private readonly BeneficiaryDbContext _context;
    private readonly ILogger<SqlBeneficiaryDataGateway> _logger;

    public SqlBeneficiaryDataGateway(BeneficiaryDbContext context, ILogger<SqlBeneficiaryDataGateway> logger)
    {
        _context = context;
        _logger = logger;
    }

    public async Task<CreateIndividualBeneficiaryResponse> CreateIndividualBeneficiaryAsync(CreateIndividualBeneficiaryRequest request, CancellationToken ct = default)
    {
        var beneficiaryId = Guid.NewGuid();
        var code = $"IND-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        var now = DateTime.UtcNow;

        _logger.LogInformation("SQL: Creating individual beneficiary for user {UserId}", request.UserId);

        var beneficiary = new BeneficiaryEntity
        {
            BeneficiaryId = beneficiaryId,
            UserId = Guid.TryParse(request.UserId, out var uid) ? uid : Guid.Empty,
            Type = "INDIVIDUAL",
            Code = code,
            Status = "PENDING_KYC",
            KycStatus = "NOT_STARTED",
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.Beneficiaries.Add(beneficiary);

        var individual = new IndividualBeneficiaryEntity
        {
            BeneficiaryId = beneficiaryId,
            FullName = request.FullName ?? "",
            DateOfBirth = DateTime.UtcNow.AddYears(-30),
            Gender = request.Gender ?? "",
            NidNumber = request.NidNumber,
            ContactInfo = request.MobileNumber ?? "",
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.IndividualBeneficiaries.Add(individual);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Created individual beneficiary {BeneficiaryId} with code {Code}", beneficiaryId, code);

        return new CreateIndividualBeneficiaryResponse
        {
            BeneficiaryId = beneficiaryId.ToString(),
            BeneficiaryCode = code
        };
    }

    public async Task<CreateBusinessBeneficiaryResponse> CreateBusinessBeneficiaryAsync(CreateBusinessBeneficiaryRequest request, CancellationToken ct = default)
    {
        var beneficiaryId = Guid.NewGuid();
        var code = $"BUS-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid().ToString()[..8].ToUpper()}";
        var now = DateTime.UtcNow;

        _logger.LogInformation("SQL: Creating business beneficiary for user {UserId}", request.UserId);

        var beneficiary = new BeneficiaryEntity
        {
            BeneficiaryId = beneficiaryId,
            UserId = Guid.TryParse(request.UserId, out var uid) ? uid : Guid.Empty,
            Type = "BUSINESS",
            Code = code,
            Status = "PENDING_KYC",
            KycStatus = "NOT_STARTED",
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.Beneficiaries.Add(beneficiary);

        var business = new BusinessBeneficiaryEntity
        {
            BeneficiaryId = beneficiaryId,
            ParentBeneficiaryId = beneficiaryId,
            BusinessName = request.BusinessName ?? "",
            TradeLicenseNumber = request.TradeLicenseNumber ?? "",
            TinNumber = request.TinNumber ?? "",
            FocalPersonName = request.FocalPersonName ?? "",
            FocalPersonContact = request.FocalPersonMobile ?? "",
            CreatedAt = now,
            UpdatedAt = now
        };

        _context.BusinessBeneficiaries.Add(business);
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Created business beneficiary {BeneficiaryId} with code {Code}", beneficiaryId, code);

        return new CreateBusinessBeneficiaryResponse
        {
            BeneficiaryId = beneficiaryId.ToString(),
            BeneficiaryCode = code
        };
    }

    public async Task<GetBeneficiaryResponse> GetBeneficiaryAsync(GetBeneficiaryRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Getting beneficiary {BeneficiaryId}", request.BeneficiaryId);

        if (string.IsNullOrEmpty(request.BeneficiaryId))
        {
            return new GetBeneficiaryResponse
            {
                Error = new Error { Code = "INVALID_REQUEST", Message = "BeneficiaryId is required" }
            };
        }

        var id = Guid.TryParse(request.BeneficiaryId, out var bid) ? bid : Guid.Empty;
        var beneficiary = await _context.Beneficiaries
            .Include(b => b.IndividualDetails)
            .Include(b => b.BusinessDetails)
            .FirstOrDefaultAsync(b => b.BeneficiaryId == id, ct);

        if (beneficiary == null)
        {
            return new GetBeneficiaryResponse
            {
                Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" }
            };
        }

        var type = Enum.TryParse<Insuretech.Beneficiary.Entity.V1.BeneficiaryType>(beneficiary.Type, true, out var bt) ? bt : Insuretech.Beneficiary.Entity.V1.BeneficiaryType.Unspecified;
        var status = Enum.TryParse<Insuretech.Beneficiary.Entity.V1.BeneficiaryStatus>(beneficiary.Status, true, out var bs) ? bs : Insuretech.Beneficiary.Entity.V1.BeneficiaryStatus.Unspecified;
        var kycStatus = Enum.TryParse<Insuretech.Beneficiary.Entity.V1.KYCStatus>(beneficiary.KycStatus, true, out var ks) ? ks : Insuretech.Beneficiary.Entity.V1.KYCStatus.Unspecified;

        return new GetBeneficiaryResponse
        {
            Beneficiary = new Insuretech.Beneficiary.Entity.V1.Beneficiary
            {
                BeneficiaryId = beneficiary.BeneficiaryId.ToString(),
                UserId = beneficiary.UserId.ToString(),
                Type = type,
                Code = beneficiary.Code,
                Status = status,
                KycStatus = kycStatus
            }
        };
    }

    public async Task<UpdateBeneficiaryResponse> UpdateBeneficiaryAsync(UpdateBeneficiaryRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Updating beneficiary {BeneficiaryId}", request.BeneficiaryId);

        var id = Guid.TryParse(request.BeneficiaryId, out var bid) ? bid : Guid.Empty;
        var beneficiary = await _context.Beneficiaries.FindAsync([id], ct);

        if (beneficiary == null)
        {
            return new UpdateBeneficiaryResponse
            {
                Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" }
            };
        }

        beneficiary.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Updated beneficiary {BeneficiaryId}", request.BeneficiaryId);

        return new UpdateBeneficiaryResponse { Message = "Beneficiary updated" };
    }

    public async Task<CompleteKYCResponse> CompleteKYCAsync(CompleteKYCRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Completing KYC for beneficiary {BeneficiaryId}", request.BeneficiaryId);

        var id = Guid.TryParse(request.BeneficiaryId, out var bid) ? bid : Guid.Empty;
        var beneficiary = await _context.Beneficiaries.FindAsync([id], ct);

        if (beneficiary == null)
        {
            return new CompleteKYCResponse
            {
                Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" }
            };
        }

        beneficiary.KycStatus = "COMPLETED";
        beneficiary.KycCompletedAt = DateTime.UtcNow;
        beneficiary.Status = "ACTIVE";
        beneficiary.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Completed KYC for beneficiary {BeneficiaryId}", request.BeneficiaryId);

        return new CompleteKYCResponse { Message = "KYC completed" };
    }

    public async Task<UpdateRiskScoreResponse> UpdateRiskScoreAsync(UpdateRiskScoreRequest request, CancellationToken ct = default)
    {
        _logger.LogInformation("SQL: Updating risk score for beneficiary {BeneficiaryId}", request.BeneficiaryId);

        var id = Guid.TryParse(request.BeneficiaryId, out var bid) ? bid : Guid.Empty;
        var beneficiary = await _context.Beneficiaries.FindAsync([id], ct);

        if (beneficiary == null)
        {
            return new UpdateRiskScoreResponse
            {
                Error = new Error { Code = "NOT_FOUND", Message = "Beneficiary not found" }
            };
        }

        beneficiary.RiskScore = request.RiskScore;
        beneficiary.UpdatedAt = DateTime.UtcNow;
        await _context.SaveChangesAsync(ct);

        _logger.LogInformation("SQL: Updated risk score for beneficiary {BeneficiaryId} to {RiskScore}", request.BeneficiaryId, request.RiskScore);

        return new UpdateRiskScoreResponse { Message = "Risk score updated" };
    }
}
