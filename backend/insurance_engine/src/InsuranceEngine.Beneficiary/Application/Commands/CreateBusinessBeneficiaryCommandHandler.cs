using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CreateBusinessBeneficiaryCommandHandler : IRequestHandler<CreateBusinessBeneficiaryCommand, Result<string>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CreateBusinessBeneficiaryCommandHandler> _logger;

    public CreateBusinessBeneficiaryCommandHandler(
        DbContext dbContext,
        ILogger<CreateBusinessBeneficiaryCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(CreateBusinessBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        using var transaction = await _dbContext.Database.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Generate IDs and Code
            var parentBeneficiaryId = Guid.NewGuid();
            var childBeneficiaryId = Guid.NewGuid();
            var beneficiaryCode = $"BEN-B-{Guid.NewGuid().ToString().Substring(0, 8).ToUpper()}";

            // 2. Insert into beneficiaries table
            var insertBeneficiarySql = @"
                INSERT INTO insurance_schema.beneficiaries (
                    beneficiary_id, user_id, type, code, status, kyc_status, audit_info, partner_id
                ) VALUES (
                    @p0, @p1, @p2, @p3, @p4, @p5, @p6::jsonb, @p7
                )";

            var auditInfo = "{}";

            await _dbContext.Database.ExecuteSqlRawAsync(insertBeneficiarySql,
                new object[] {
                    parentBeneficiaryId,
                    Guid.Parse(request.UserId),
                    "BUSINESS",
                    beneficiaryCode,
                    "PENDING_KYC",
                    "NOT_STARTED",
                    auditInfo,
                    string.IsNullOrEmpty(request.PartnerId) ? (Guid?)null : Guid.Parse(request.PartnerId)
                }, cancellationToken);

            // 3. Insert into business_beneficiaries table
            var insertBusinessSql = @"
                INSERT INTO insurance_schema.business_beneficiaries (
                    id, beneficiary_id, parent_beneficiary_id, business_name, trade_license_number, tin_number, 
                    business_type, contact_info, registered_address, business_address, 
                    focal_person_name, focal_person_contact, primary_contact,
                    total_employees_covered, active_policies_count, total_premium_amount, pending_actions_count,
                    audit_info
                ) VALUES (
                    @p0, @p1, @p2, @p3, @p4, @p5, @p6, @p7::jsonb, @p8::jsonb, @p9::jsonb, @p10, @p11::jsonb, @p12::jsonb, @p13, @p14, @p15, @p16, @p17::jsonb
                )";

            var contactInfo = $"{{\"mobile_number\": \"{request.FocalPersonMobile}\"}}";
            var emptyJson = "{}";

            await _dbContext.Database.ExecuteSqlRawAsync(insertBusinessSql,
                new object[] {
                    Guid.NewGuid(), // id (Child PK)
                    parentBeneficiaryId, // beneficiary_id (FK to base)
                    parentBeneficiaryId, // parent_beneficiary_id (Additional FK/Ref)
                    request.BusinessName,
                    request.TradeLicenseNumber,
                    request.TinNumber,
                    "BUSINESS_TYPE_SOLE_PROPRIETORSHIP", 
                    contactInfo,
                    emptyJson, 
                    emptyJson, 
                    request.FocalPersonName,
                    contactInfo,
                    emptyJson,
                    0,
                    0,
                    0L,
                    0,
                    auditInfo
                }, cancellationToken);

            await transaction.CommitAsync(cancellationToken);

            _logger.LogInformation("Business Beneficiary created successfully: {ParentId} ({Code})", parentBeneficiaryId, beneficiaryCode);

            return Result<string>.Ok(parentBeneficiaryId.ToString());
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to create business beneficiary for user {UserId}", request.UserId);
            return Result<string>.Fail("BUSINESS_BENEFICIARY_CREATION_FAILED", ex.Message);
        }
    }
}
