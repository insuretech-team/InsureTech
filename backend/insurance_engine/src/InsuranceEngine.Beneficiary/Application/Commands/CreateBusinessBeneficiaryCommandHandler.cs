using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using Dapper;
using System.Data;
using InsuranceEngine.Beneficiary.Domain;

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
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Domain Logic: Create Aggregate (DDD)
            var beneficiary = BeneficiaryAggregate.CreateBusiness(
                userId: Guid.Parse(request.UserId),
                businessName: request.BusinessName,
                partnerId: string.IsNullOrEmpty(request.PartnerId) ? (Guid?)null : Guid.Parse(request.PartnerId)
            );

            // 2. Persist Aggregate State
            var insertBeneficiarySql = @"
                INSERT INTO insurance_schema.beneficiaries (
                    beneficiary_id, user_id, type, code, status, kyc_status, audit_info, partner_id
                ) VALUES (
                    @BeneficiaryId, @UserId, @Type, @Code, @Status, @KycStatus, @AuditInfo::jsonb, @PartnerId
                )";

            await connection.ExecuteAsync(insertBeneficiarySql, new
            {
                BeneficiaryId = beneficiary.Id,
                UserId = beneficiary.UserId,
                Type = beneficiary.Type,
                Code = beneficiary.Code,
                Status = beneficiary.Status,
                KycStatus = beneficiary.KycStatus,
                AuditInfo = "{}",
                PartnerId = beneficiary.PartnerId
            }, transaction);

            var insertBusinessSql = @"
                INSERT INTO insurance_schema.business_beneficiaries (
                    id, beneficiary_id, parent_beneficiary_id, business_name, trade_license_number, tin_number, 
                    business_type, contact_info, registered_address, business_address, 
                    focal_person_name, focal_person_contact, primary_contact,
                    total_employees_covered, active_policies_count, total_premium_amount, pending_actions_count,
                    audit_info
                ) VALUES (
                    @Id, @BeneficiaryId, @ParentBeneficiaryId, @BusinessName, @TradeLicenseNumber, @TinNumber, 
                    @BusinessType, @ContactInfo::jsonb, @RegisteredAddress::jsonb, @BusinessAddress::jsonb, 
                    @FocalPersonName, @FocalPersonContact::jsonb, @PrimaryContact::jsonb,
                    @TotalEmployees, @ActivePolicies, @TotalPremium, @PendingActions,
                    @AuditInfo::jsonb
                )";

            var contactInfo = $"{{\"mobile_number\": \"{request.FocalPersonMobile}\"}}";

            await connection.ExecuteAsync(insertBusinessSql, new
            {
                Id = Guid.NewGuid(),
                BeneficiaryId = beneficiary.Id,
                ParentBeneficiaryId = beneficiary.Id,
                BusinessName = request.BusinessName,
                TradeLicenseNumber = request.TradeLicenseNumber,
                TinNumber = request.TinNumber,
                BusinessType = "BUSINESS_TYPE_SOLE_PROPRIETORSHIP", 
                ContactInfo = contactInfo,
                RegisteredAddress = "{}", 
                BusinessAddress = "{}", 
                FocalPersonName = request.FocalPersonName,
                FocalPersonContact = contactInfo,
                PrimaryContact = "{}",
                TotalEmployees = 0,
                ActivePolicies = 0,
                TotalPremium = 0L,
                PendingActions = 0,
                AuditInfo = "{}"
            }, transaction);

            await transaction.CommitAsync(cancellationToken);

            _logger.LogInformation("Business Beneficiary created successfully: {ParentId} ({Code})", beneficiary.Id, beneficiary.Code);

            return Result<string>.Ok(beneficiary.Id.ToString());
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to create business beneficiary");
            return Result<string>.Fail("BUSINESS_BENEFICIARY_CREATION_FAILED", ex.Message);
        }
    }
}
