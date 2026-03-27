using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using Dapper;
using System.Data;
using InsuranceEngine.Beneficiary.Domain;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public sealed class CreateIndividualBeneficiaryCommandHandler : IRequestHandler<CreateIndividualBeneficiaryCommand, Result<string>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CreateIndividualBeneficiaryCommandHandler> _logger;

    public CreateIndividualBeneficiaryCommandHandler(
        DbContext dbContext,
        ILogger<CreateIndividualBeneficiaryCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(CreateIndividualBeneficiaryCommand request, CancellationToken cancellationToken)
    {
        var connection = _dbContext.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        // Domain Validation: Check for uniqueness (FR-033)
        const string checkSql = "SELECT COUNT(1) FROM insurance_schema.individual_beneficiaries WHERE nid_number = @Nid";
        var exists = await connection.ExecuteScalarAsync<int>(checkSql, new { Nid = request.NidNumber });
        if (exists > 0) return Result<string>.Fail("DUPLICATE_BENEFICIARY", "Beneficiary with this NID already exists");

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Domain Logic: Create Aggregate (DDD)
            var partnerId = string.IsNullOrEmpty(request.PartnerId) ? (Guid?)null : Guid.Parse(request.PartnerId);
            var beneficiary = BeneficiaryAggregate.CreateIndividual(
                userId: string.IsNullOrEmpty(request.UserId) ? (Guid?)null : Guid.Parse(request.UserId),
                fullName: request.FullName,
                partnerId: partnerId
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

            var insertIndividualSql = @"
                INSERT INTO insurance_schema.individual_beneficiaries (
                    id, beneficiary_id, full_name, date_of_birth, gender, nid_number, contact_info, permanent_address, present_address, marital_status, occupation, audit_info
                ) VALUES (
                    @IndividualId, @BeneficiaryId, @FullName, @DateOfBirth, @Gender, @NidNumber, @ContactInfo::jsonb, @PermanentAddress::jsonb, @PresentAddress::jsonb, @MaritalStatus, @Occupation, @AuditInfo::jsonb
                )";

            var contactInfo = $"{{\"mobile_number\": \"{request.MobileNumber}\", \"email\": \"{request.Email ?? ""}\"}}";

            await connection.ExecuteAsync(insertIndividualSql, new
            {
                IndividualId = Guid.NewGuid(),
                BeneficiaryId = beneficiary.Id,
                FullName = request.FullName,
                DateOfBirth = DateTime.SpecifyKind(request.DateOfBirth, DateTimeKind.Utc),
                Gender = request.Gender,
                NidNumber = request.NidNumber,
                ContactInfo = contactInfo,
                PermanentAddress = "{}",
                PresentAddress = "{}",
                MaritalStatus = "MARITAL_STATUS_SINGLE",
                Occupation = "OTHER",
                AuditInfo = "{}"
            }, transaction);

            await transaction.CommitAsync(cancellationToken);
            
            _logger.LogInformation("Beneficiary created successfully: {BeneficiaryId} ({BeneficiaryCode})", beneficiary.Id, beneficiary.Code);

            return Result<string>.Ok(beneficiary.Id.ToString());
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            _logger.LogError(ex, "Failed to create beneficiary");
            return Result<string>.Fail("BENEFICIARY_CREATION_FAILED", ex.Message);
        }
    }
}
