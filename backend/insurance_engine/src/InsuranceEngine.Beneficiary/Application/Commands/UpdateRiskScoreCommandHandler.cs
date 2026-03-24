using System.Data;
using Dapper;
using Insuretech.Beneficiary.Services.V1;
using MediatR;
using Microsoft.EntityFrameworkCore;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public class UpdateRiskScoreCommandHandler : IRequestHandler<UpdateRiskScoreCommand, UpdateRiskScoreResponse>
{
    private readonly DbContext _context;

    public UpdateRiskScoreCommandHandler(DbContext context)
    {
        _context = context;
    }

    public async Task<UpdateRiskScoreResponse> Handle(UpdateRiskScoreCommand command, CancellationToken cancellationToken)
    {
        var request = command.Request;
        var connection = _context.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        const string sql = @"
            UPDATE insurance_schema.beneficiaries 
            SET 
                risk_score = @Score,
                risk_score_reason = @Reason,
                updated_at = NOW()
            WHERE beneficiary_id = @Id::uuid";

        var affected = await connection.ExecuteAsync(sql, new { 
            Id = request.BeneficiaryId,
            Score = request.RiskScore,
            Reason = request.Reason
        });

        if (affected == 0)
        {
            return new UpdateRiskScoreResponse
            {
                Error = new Insuretech.Common.V1.Error { Code = "BENEFICIARY_NOT_FOUND", Message = "Beneficiary not found", HttpStatusCode = 404 }
            };
        }

        return new UpdateRiskScoreResponse
        {
            Message = "Risk score updated successfully"
        };
    }
}
