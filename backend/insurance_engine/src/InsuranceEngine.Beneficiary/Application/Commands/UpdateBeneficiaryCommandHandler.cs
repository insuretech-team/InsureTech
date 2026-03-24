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
                updateSql = @"
                    UPDATE insurance_schema.individual_beneficiaries 
                    SET 
                        mobile_number = COALESCE(NULLIF(@Mobile, ''), mobile_number),
                        email = COALESCE(NULLIF(@Email, ''), email),
                        address = COALESCE(NULLIF(@Address, ''), address),
                        updated_at = NOW()
                    WHERE beneficiary_id = @Id::uuid";
                
                parameters = new 
                { 
                    Id = request.BeneficiaryId, 
                    Mobile = request.MobileNumber, 
                    Email = request.Email, 
                    Address = request.Address 
                };
            }
            else if (type == "BUSINESS")
            {
                updateSql = @"
                    UPDATE insurance_schema.business_beneficiaries 
                    SET 
                        focal_person_mobile = COALESCE(NULLIF(@Mobile, ''), focal_person_mobile),
                        focal_person_email = COALESCE(NULLIF(@Email, ''), focal_person_email),
                        business_address = COALESCE(NULLIF(@Address, ''), business_address),
                        updated_at = NOW()
                    WHERE beneficiary_id = @Id::uuid";
                
                parameters = new 
                { 
                    Id = request.BeneficiaryId, 
                    Mobile = request.MobileNumber, 
                    Email = request.Email, 
                    Address = request.Address 
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
