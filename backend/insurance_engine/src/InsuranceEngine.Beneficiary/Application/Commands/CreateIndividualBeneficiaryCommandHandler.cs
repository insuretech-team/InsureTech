using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Beneficiary.Entity.V1;

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
        try
        {
            // 1. Validate Uniqueness (NID and Mobile)
            // (Skipped for simplicity in this clean build)

            // 2. Generate Beneficiary ID and Code
            var beneficiaryId = Guid.NewGuid();
            var individualId = Guid.NewGuid();
            var beneficiaryCode = $"BEN-{Guid.NewGuid().ToString().Substring(0, 8).ToUpper()}";

            // 3. Insert into beneficiaries table
            var insertBeneficiarySql = @"
                INSERT INTO insurance_schema.beneficiaries (
                    beneficiary_id, user_id, type, code, status, kyc_status, audit_info
                ) VALUES (
                    @p0, @p1, @p2, @p3, @p4, @p5, @p6::jsonb
                )";

            var auditInfo = "{}"; // Proto expects JSONB

            await _dbContext.Database.ExecuteSqlRawAsync(insertBeneficiarySql,
                new object[] {
                    beneficiaryId,
                    Guid.Parse(request.UserId),
                    "INDIVIDUAL",
                    beneficiaryCode,
                    "PENDING_KYC",
                    "NOT_STARTED",
                    auditInfo
                }, cancellationToken);

            // 4. Insert into individual_beneficiaries table
            var insertIndividualSql = @"
                INSERT INTO insurance_schema.individual_beneficiaries (
                    id, beneficiary_id, full_name, date_of_birth, gender, nid_number, contact_info, permanent_address, present_address, marital_status, occupation, audit_info
                ) VALUES (
                    @p0, @p1, @p2, @p3, @p4, @p5, @p6::jsonb, @p7::jsonb, @p8::jsonb, @p9, @p10, @p11::jsonb
                )";

            var contactInfo = $"{{\"mobile_number\": \"{request.MobileNumber}\", \"email\": \"{request.Email ?? ""}\"}}";
            var emptyAddress = "{}";

            await _dbContext.Database.ExecuteSqlRawAsync(insertIndividualSql,
                new object[] {
                    individualId,
                    beneficiaryId,
                    request.FullName,
                    request.DateOfBirth,
                    request.Gender,
                    request.NidNumber,
                    contactInfo,
                    emptyAddress,
                    emptyAddress,
                    "MARITAL_STATUS_SINGLE",
                    "OTHER",
                    auditInfo
                }, cancellationToken);

            _logger.LogInformation("Beneficiary created successfully: {BeneficiaryId} ({BeneficiaryCode})", beneficiaryId, beneficiaryCode);

            return Result<string>.Ok(beneficiaryId.ToString());
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create beneficiary for user {UserId}", request.UserId);
            return Result<string>.Fail("BENEFICIARY_CREATION_FAILED", ex.Message);
        }
    }
}
