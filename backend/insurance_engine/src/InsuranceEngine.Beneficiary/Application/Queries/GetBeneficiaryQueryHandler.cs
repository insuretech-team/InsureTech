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
                i.full_name, i.date_of_birth, i.gender, i.nid_number, i.contact_info, i.permanent_address, i.present_address,
                biz.business_name, biz.trade_license_number, biz.tin_number, biz.focal_person_name, biz.business_type, biz.contact_info as biz_contact_info, biz.business_address
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

        try
        {
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
                var contactInfo = new Insuretech.Common.V1.ContactInfo();
                var permanentAddress = new Insuretech.Common.V1.Address();
                var presentAddress = new Insuretech.Common.V1.Address();

                try
                {
                    if (result.contact_info != null)
                    {
                        var contactDict = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(result.contact_info.ToString());
                        if (contactDict != null)
                        {
                            if (contactDict.TryGetValue("mobile_number", out object mobile)) contactInfo.MobileNumber = mobile?.ToString();
                            if (contactDict.TryGetValue("email", out object email)) contactInfo.Email = email?.ToString();
                            if (contactDict.TryGetValue("alternate_mobile", out object altMobile)) contactInfo.AlternateMobile = altMobile?.ToString();
                            if (contactDict.TryGetValue("landline", out object landline)) contactInfo.Landline = landline?.ToString();
                        }
                    }

                    if (result.permanent_address != null)
                    {
                        var permAddrDict = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(result.permanent_address.ToString());
                        if (permAddrDict != null)
                        {
                            if (permAddrDict.TryGetValue("address_line", out object addr)) permanentAddress.AddressLine1 = addr?.ToString();
                            if (permAddrDict.TryGetValue("address_line1", out object addr1)) permanentAddress.AddressLine1 = addr1?.ToString();
                            if (permAddrDict.TryGetValue("city", out object city)) permanentAddress.City = city?.ToString();
                            if (permAddrDict.TryGetValue("district", out object district)) permanentAddress.District = district?.ToString();
                            if (permAddrDict.TryGetValue("division", out object division)) permanentAddress.Division = division?.ToString();
                            if (permAddrDict.TryGetValue("postal_code", out object postal)) permanentAddress.PostalCode = postal?.ToString();
                            if (permAddrDict.TryGetValue("country", out object country)) permanentAddress.Country = country?.ToString();
                        }
                    }

                    if (result.present_address != null)
                    {
                        var presAddrDict = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(result.present_address.ToString());
                        if (presAddrDict != null)
                        {
                            if (presAddrDict.TryGetValue("address_line", out object addr)) presentAddress.AddressLine1 = addr?.ToString();
                            if (presAddrDict.TryGetValue("address_line1", out object addr1)) presentAddress.AddressLine1 = addr1?.ToString();
                            if (presAddrDict.TryGetValue("city", out object city)) presentAddress.City = city?.ToString();
                            if (presAddrDict.TryGetValue("district", out object district)) presentAddress.District = district?.ToString();
                            if (presAddrDict.TryGetValue("division", out object division)) presentAddress.Division = division?.ToString();
                            if (presAddrDict.TryGetValue("postal_code", out object postal)) presentAddress.PostalCode = postal?.ToString();
                            if (presAddrDict.TryGetValue("country", out object country)) presentAddress.Country = country?.ToString();
                        }
                    }
                }
                catch { }

                response.IndividualDetails = new Insuretech.Beneficiary.Entity.V1.IndividualBeneficiary
                {
                    BeneficiaryId = result.beneficiary_id.ToString(),
                    FullName = result.full_name ?? "",
                    DateOfBirth = result.date_of_birth != null ? Timestamp.FromDateTime(((DateOnly)result.date_of_birth).ToDateTime(TimeOnly.MinValue, DateTimeKind.Utc)) : null,
                    Gender = result.gender != null ? System.Enum.Parse<Insuretech.Beneficiary.Entity.V1.Gender>(result.gender.ToString().Replace("GENDER_", "").Replace("_", ""), true) : Insuretech.Beneficiary.Entity.V1.Gender.Unspecified,
                    NidNumber = result.nid_number ?? "",
                    ContactInfo = contactInfo,
                    PermanentAddress = permanentAddress,
                    PresentAddress = presentAddress
                };
            }
            else if (result.type == "BUSINESS")
            {
                var contactInfo = new Insuretech.Common.V1.ContactInfo();
                var businessAddress = new Insuretech.Common.V1.Address();

                try
                {
                    if (result.biz_contact_info != null)
                    {
                        var contactDict = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(result.biz_contact_info.ToString());
                        if (contactDict != null)
                        {
                            if (contactDict.TryGetValue("focal_person_mobile", out object mobile)) contactInfo.MobileNumber = mobile?.ToString();
                            if (contactDict.TryGetValue("focal_person_email", out object email)) contactInfo.Email = email?.ToString();
                        }
                    }

                    if (result.business_address != null)
                    {
                        var addrDict = System.Text.Json.JsonSerializer.Deserialize<Dictionary<string, object>>(result.business_address.ToString());
                        if (addrDict != null)
                        {
                            if (addrDict.TryGetValue("address_line", out object addr)) businessAddress.AddressLine1 = addr?.ToString();
                            if (addrDict.TryGetValue("address_line1", out object addr1)) businessAddress.AddressLine1 = addr1?.ToString();
                            if (addrDict.TryGetValue("city", out object city)) businessAddress.City = city?.ToString();
                            if (addrDict.TryGetValue("district", out object district)) businessAddress.District = district?.ToString();
                        }
                    }
                }
                catch { }

                var businessType = Insuretech.Beneficiary.Entity.V1.BusinessType.Unspecified;
                if (result.business_type != null)
                {
                    var typeStr = result.business_type.ToString().Replace("BUSINESS_TYPE_", "").Replace("_", "");
                    try
                    {
                        businessType = (Insuretech.Beneficiary.Entity.V1.BusinessType)System.Enum.Parse(typeof(Insuretech.Beneficiary.Entity.V1.BusinessType), typeStr, true);
                    }
                    catch
                    {
                        businessType = Insuretech.Beneficiary.Entity.V1.BusinessType.Unspecified;
                    }
                }

                response.BusinessDetails = new Insuretech.Beneficiary.Entity.V1.BusinessBeneficiary
                {
                    Id = result.beneficiary_id.ToString(),
                    BusinessName = result.business_name ?? "",
                    TradeLicenseNumber = result.trade_license_number ?? "",
                    TinNumber = result.tin_number ?? "",
                    FocalPersonName = result.focal_person_name ?? "",
                    BusinessType = businessType,
                    ContactInfo = contactInfo,
                    BusinessAddress = businessAddress
                };
            }

            return response;
        }
        catch (Exception ex)
        {
            return new GetBeneficiaryResponse
            {
                Error = new Insuretech.Common.V1.Error
                {
                    Code = "GET_BENEFICIARY_FAILED",
                    Message = ex.Message,
                    HttpStatusCode = 500
                }
            };
        }
    }
}
