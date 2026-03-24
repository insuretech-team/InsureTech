using System.Data;
using Dapper;
using Insuretech.Beneficiary.Services.V1;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Google.Protobuf.WellKnownTypes;

namespace InsuranceEngine.Beneficiary.Application.Queries;

public class GetBeneficiaryQueryHandler : IRequestHandler<GetBeneficiaryQuery, GetBeneficiaryResponse>
{
    private readonly DbContext _context;

    public GetBeneficiaryQueryHandler(DbContext context)
    {
        _context = context;
    }

    public async Task<GetBeneficiaryResponse> Handle(GetBeneficiaryQuery request, CancellationToken cancellationToken)
    {
        var connection = _context.Database.GetDbConnection();
        if (connection.State != ConnectionState.Open) await connection.OpenAsync(cancellationToken);

        const string sql = @"
            SELECT 
                b.beneficiary_id, b.code, b.type, b.status, b.kyc_status, b.risk_score, b.kyc_completed_at, b.referral_code, b.partner_id,
                i.full_name, i.date_of_birth, i.gender, i.nid_number, 
                biz.business_name, biz.trade_license_number, biz.tin_number, biz.focal_person_name, biz.business_type
            FROM insurance_schema.beneficiaries b
            LEFT JOIN insurance_schema.individual_beneficiaries i ON b.beneficiary_id = i.beneficiary_id
            LEFT JOIN insurance_schema.business_beneficiaries biz ON b.beneficiary_id = biz.beneficiary_id
            WHERE b.beneficiary_id = @Id::uuid";

        var result = await connection.QueryFirstOrDefaultAsync<dynamic>(sql, new { Id = request.BeneficiaryId });

        if (result == null)
        {
            return new GetBeneficiaryResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "BENEFICIARY_NOT_FOUND",
                    Message = $"Beneficiary with ID {request.BeneficiaryId} not found",
                    HttpStatusCode = 404
                }
            };
        }

        var response = new GetBeneficiaryResponse
        {
            Beneficiary = new Insuretech.Beneficiary.Entity.V1.Beneficiary
            {
                BeneficiaryId = result.beneficiary_id.ToString(),
                Code = result.code,
                Type = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.BeneficiaryType>(result.type.ToString().Replace("BENEFICIARY_TYPE_", "").Replace("_", ""), true),
                Status = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.BeneficiaryStatus>(result.status.ToString().Replace("BENEFICIARY_STATUS_", "").Replace("_", ""), true),
                KycStatus = System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.KYCStatus>(result.kyc_status.ToString().Replace("KYC_STATUS_", "").Replace("_", ""), true),
                RiskScore = result.risk_score ?? "",
                KycCompletedAt = result.kyc_completed_at != null ? Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)result.kyc_completed_at, DateTimeKind.Utc)) : null,
                PartnerId = result.partner_id?.ToString() ?? ""
            }
        };

        if (result.type == "INDIVIDUAL")
        {
            response.IndividualDetails = new Insuretech.Beneficiary.Entity.V1.IndividualBeneficiary
            {
                BeneficiaryId = result.beneficiary_id.ToString(),
                FullName = result.full_name ?? "",
                DateOfBirth = result.date_of_birth != null ? Timestamp.FromDateTime(((DateOnly)result.date_of_birth).ToDateTime(TimeOnly.MinValue, DateTimeKind.Utc)) : null,
                Gender = result.gender != null ? System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.Gender>(result.gender.ToString().Replace("GENDER_", "").Replace("_", ""), true) : Insuretech.Beneficiary.Entity.V1.Gender.Unspecified,
                NidNumber = result.nid_number ?? ""
            };
        }
        else if (result.type == "BUSINESS")
        {
            response.BusinessDetails = new Insuretech.Beneficiary.Entity.V1.BusinessBeneficiary
            {
                Id = result.beneficiary_id.ToString(),
                BusinessName = result.business_name ?? "",
                TradeLicenseNumber = result.trade_license_number ?? "",
                TinNumber = result.tin_number ?? "",
                FocalPersonName = result.focal_person_name ?? "",
                BusinessType = result.business_type != null ? System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.BusinessType>(result.business_type.ToString().Replace("_", ""), true) : Insuretech.Beneficiary.Entity.V1.BusinessType.Unspecified
            };
        }

        return response;
    }
}
