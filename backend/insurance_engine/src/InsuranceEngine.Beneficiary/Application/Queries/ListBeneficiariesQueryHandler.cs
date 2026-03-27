using System.Data;
using Dapper;
using Insuretech.Beneficiary.Services.V1;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Beneficiary.Application.Queries;

public class ListBeneficiariesQueryHandler : IRequestHandler<ListBeneficiariesQuery, ListBeneficiariesResponse>
{
    private readonly DbContext _context;

    public ListBeneficiariesQueryHandler(DbContext context)
    {
        _context = context;
    }

    public async Task<ListBeneficiariesResponse> Handle(ListBeneficiariesQuery request, CancellationToken cancellationToken)
    {
        var connection = _context.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        string sql = @"
            SELECT 
                beneficiary_id, user_id, code, type, status, kyc_status, risk_score, kyc_completed_at, partner_id
            FROM insurance_schema.beneficiaries 
            WHERE 1=1";

        var parameters = new DynamicParameters();
        if (!string.IsNullOrEmpty(request.Type))
        {
            sql += " AND type = @Type";
            parameters.Add("Type", request.Type);
        }
        if (!string.IsNullOrEmpty(request.Status))
        {
            sql += " AND status = @Status";
            parameters.Add("Status", request.Status);
        }

        sql += " ORDER BY beneficiary_id DESC LIMIT @Limit OFFSET @Offset";
        parameters.Add("Limit", request.PageSize);
        parameters.Add("Offset", (request.Page - 1) * request.PageSize);

        var results = await connection.QueryAsync<dynamic>(sql, parameters);
        
        var countSql = "SELECT COUNT(*) FROM insurance_schema.beneficiaries WHERE 1=1";
        if (!string.IsNullOrEmpty(request.Type)) countSql += " AND type = @Type";
        if (!string.IsNullOrEmpty(request.Status)) countSql += " AND status = @Status";
        var totalCount = await connection.ExecuteScalarAsync<int>(countSql, parameters);

        var response = new ListBeneficiariesResponse
        {
            TotalCount = totalCount
        };

        foreach (var item in results)
        {
            response.Beneficiaries.Add(new Insuretech.Beneficiary.Entity.V1.Beneficiary
            {
                BeneficiaryId = item.beneficiary_id.ToString(),
                UserId = item.user_id?.ToString() ?? "",
                Code = item.code,
                Type = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.BeneficiaryType>(item.type.ToString().Replace("BENEFICIARY_TYPE_", "").Replace("_", ""), true),
                Status = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.BeneficiaryStatus>(item.status.ToString().Replace("BENEFICIARY_STATUS_", "").Replace("_", ""), true),
                KycStatus = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.KYCStatus>(item.kyc_status.ToString().Replace("KYC_STATUS_", "").Replace("_", ""), true),
                RiskScore = item.risk_score ?? "",
                KycCompletedAt = item.kyc_completed_at != null ? Timestamp.FromDateTime(DateTime.SpecifyKind(item.kyc_completed_at, DateTimeKind.Utc)) : null,
                PartnerId = item.partner_id?.ToString() ?? ""
            });
        }

        return response;
    }
}
