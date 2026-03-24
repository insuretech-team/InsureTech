using System.Data;
using Dapper;
using Insuretech.Beneficiary.Services.V1;
using MediatR;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public class CompleteKYCCommandHandler : IRequestHandler<CompleteKYCCommand, CompleteKYCResponse>
{
    private readonly DbContext _context;

    public CompleteKYCCommandHandler(DbContext context)
    {
        _context = context;
    }

    public async Task<CompleteKYCResponse> Handle(CompleteKYCCommand command, CancellationToken cancellationToken)
    {
        var request = command.Request;
        var connection = _context.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        using var transaction = await connection.BeginTransactionAsync(cancellationToken);
        try
        {
            // 1. Update status in base table
            const string updateBaseSql = @"
                UPDATE insurance_schema.beneficiaries 
                SET 
                    kyc_status = 'COMPLETED',
                    kyc_completed_at = NOW(),
                    status = 'ACTIVE',
                    updated_at = NOW()
                WHERE beneficiary_id = @Id::uuid";
            
            var affected = await connection.ExecuteAsync(updateBaseSql, new { Id = request.BeneficiaryId }, transaction);

            if (affected == 0)
            {
                return new CompleteKYCResponse
                {
                    Error = new Insuretech.Common.V1.Error { Code = "BENEFICIARY_NOT_FOUND", Message = "Beneficiary not found", HttpStatusCode = 404 }
                };
            }

            // 2. Store document URLs in audit_info (JSONB) since dedicated columns are missing in schema
            const string updateIndividualSql = @"
                UPDATE insurance_schema.individual_beneficiaries 
                SET 
                    audit_info = audit_info || jsonb_build_object('kyc_documents', jsonb_build_object(
                        'nid_front_url', @NidFront,
                        'nid_back_url', @NidBack,
                        'selfie_url', @Selfie,
                        'porichoy_id', @PorichoyId
                    )),
                    updated_at = NOW()
                WHERE beneficiary_id = @Id::uuid";

            await connection.ExecuteAsync(updateIndividualSql, new { 
                Id = request.BeneficiaryId,
                NidFront = request.NidFrontUrl,
                NidBack = request.NidBackUrl,
                Selfie = request.SelfieUrl,
                PorichoyId = request.PorichoyVerificationId
            }, transaction);

            await transaction.CommitAsync(cancellationToken);

            return new CompleteKYCResponse
            {
                KycStatus = "COMPLETED",
                Message = "KYC completed successfully"
            };
        }
        catch (Exception ex)
        {
            await transaction.RollbackAsync(cancellationToken);
            return new CompleteKYCResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "KYC_FAILED", Message = ex.Message, HttpStatusCode = 500 }
            };
        }
    }
}
