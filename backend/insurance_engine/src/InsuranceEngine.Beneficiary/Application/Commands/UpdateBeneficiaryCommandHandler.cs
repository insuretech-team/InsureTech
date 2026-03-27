using System.Data;
using Dapper;
using Insuretech.Beneficiary.Services.V1;
using MediatR;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public class UpdateBeneficiaryCommandHandler : IRequestHandler<UpdateBeneficiaryCommand, UpdateBeneficiaryResponse>
{
    private readonly DbContext _context;

    public UpdateBeneficiaryCommandHandler(DbContext context)
    {
        _context = context;
    }

    public async Task<UpdateBeneficiaryResponse> Handle(UpdateBeneficiaryCommand command, CancellationToken cancellationToken)
    {
        var request = command.Request;
        var connection = _context.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        // 1. Get beneficiary type
        const string typeSql = "SELECT type FROM insurance_schema.beneficiaries WHERE beneficiary_id = @Id::uuid";
        var type = await connection.ExecuteScalarAsync<string>(typeSql, new { Id = request.BeneficiaryId });

        if (string.IsNullOrEmpty(type))
        {
            return new UpdateBeneficiaryResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "BENEFICIARY_NOT_FOUND",
                    Message = $"Beneficiary with ID {request.BeneficiaryId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            string updateSql = "";
            object parameters = null;

            if (type == "INDIVIDUAL")
            {
                var getContactInfoSql = "SELECT contact_info FROM insurance_schema.individual_beneficiaries WHERE beneficiary_id = @Id::uuid";
                var existingContactInfo = await connection.QueryFirstOrDefaultAsync<string>(getContactInfoSql, new { Id = request.BeneficiaryId });
                
                var contactInfoObj = string.IsNullOrEmpty(existingContactInfo) || existingContactInfo == "{}" 
                    ? new Dictionary<string, object>() 
                    : System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(existingContactInfo) ?? new Dictionary<string, object>();

                if (!string.IsNullOrEmpty(request.MobileNumber))
                    contactInfoObj["mobile_number"] = request.MobileNumber;
                if (!string.IsNullOrEmpty(request.Email))
                    contactInfoObj["email"] = request.Email;

                var newContactInfo = System.Text.Json.JsonSerializer.Serialize(contactInfoObj);

                var getAddressSql = "SELECT permanent_address FROM insurance_schema.individual_beneficiaries WHERE beneficiary_id = @Id::uuid";
                var existingAddress = await connection.QueryFirstOrDefaultAsync<string>(getAddressSql, new { Id = request.BeneficiaryId });
                
                var addressObj = string.IsNullOrEmpty(existingAddress) || existingAddress == "{}"
                    ? new Dictionary<string, object>()
                    : System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(existingAddress) ?? new Dictionary<string, object>();

                if (!string.IsNullOrEmpty(request.Address))
                {
                    addressObj["address_line"] = request.Address;
                }

                var newAddress = System.Text.Json.JsonSerializer.Serialize(addressObj);

                updateSql = @"
                    UPDATE insurance_schema.individual_beneficiaries 
                    SET 
                        contact_info = @ContactInfo::jsonb,
                        permanent_address = @PermanentAddress::jsonb
                    WHERE beneficiary_id = @Id::uuid";
                
                parameters = new 
                { 
                    Id = request.BeneficiaryId, 
                    ContactInfo = newContactInfo,
                    PermanentAddress = newAddress
                };
            }
            else if (type == "BUSINESS")
            {
                var getContactInfoSql = "SELECT contact_info FROM insurance_schema.business_beneficiaries WHERE beneficiary_id = @Id::uuid";
                var existingContactInfo = await connection.QueryFirstOrDefaultAsync<string>(getContactInfoSql, new { Id = request.BeneficiaryId });
                
                var contactInfoObj = string.IsNullOrEmpty(existingContactInfo) || existingContactInfo == "{}" 
                    ? new Dictionary<string, object>() 
                    : System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(existingContactInfo) ?? new Dictionary<string, object>();

                if (!string.IsNullOrEmpty(request.MobileNumber))
                    contactInfoObj["focal_person_mobile"] = request.MobileNumber;
                if (!string.IsNullOrEmpty(request.Email))
                    contactInfoObj["focal_person_email"] = request.Email;

                var newContactInfo = System.Text.Json.JsonSerializer.Serialize(contactInfoObj);

                var getAddressSql = "SELECT business_address FROM insurance_schema.business_beneficiaries WHERE beneficiary_id = @Id::uuid";
                var existingAddress = await connection.QueryFirstOrDefaultAsync<string>(getAddressSql, new { Id = request.BeneficiaryId });
                
                var addressObj = string.IsNullOrEmpty(existingAddress) || existingAddress == "{}"
                    ? new Dictionary<string, object>()
                    : System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(existingAddress) ?? new Dictionary<string, object>();

                if (!string.IsNullOrEmpty(request.Address))
                {
                    addressObj["address_line"] = request.Address;
                }

                var newAddress = System.Text.Json.JsonSerializer.Serialize(addressObj);

                updateSql = @"
                    UPDATE insurance_schema.business_beneficiaries 
                    SET 
                        contact_info = @ContactInfo::jsonb,
                        business_address = @BusinessAddress::jsonb
                    WHERE beneficiary_id = @Id::uuid";
                
                parameters = new 
                { 
                    Id = request.BeneficiaryId, 
                    ContactInfo = newContactInfo,
                    BusinessAddress = newAddress
                };
            }

            if (!string.IsNullOrEmpty(updateSql))
            {
                await connection.ExecuteAsync(updateSql, parameters, transaction);
            }

            await transaction.CommitAsync(cancellationToken);

            return new UpdateBeneficiaryResponse
            {
                Message = "Beneficiary updated successfully"
            };
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            return new UpdateBeneficiaryResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "UPDATE_FAILED",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }
}
