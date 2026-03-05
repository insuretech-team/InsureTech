using PoliSync.SharedKernel.Domain;

namespace PoliSync.Beneficiaries.Domain;

public class BusinessBeneficiary : Entity
{
    public Guid BeneficiaryId { get; private set; } // PK and FK to Beneficiary
    public string BusinessName { get; private set; } = string.Empty;
    public string? TradeLicenseNumber { get; private set; }
    public string? TinNumber { get; private set; }
    public string? BusinessType { get; private set; }
    public string? FocalPersonName { get; private set; }
    public string? FocalPersonMobile { get; private set; }
    public Guid? PartnerId { get; private set; }
    public string? AuditInfo { get; private set; } // JSONB

    private BusinessBeneficiary() { }

    public static BusinessBeneficiary Create(
        Guid beneficiaryId,
        string businessName,
        string tradeLicenseNumber,
        string tinNumber,
        string focalPersonName,
        string focalPersonMobile,
        Guid? partnerId = null)
    {
        return new BusinessBeneficiary
        {
            BeneficiaryId = beneficiaryId,
            BusinessName = businessName,
            TradeLicenseNumber = tradeLicenseNumber,
            TinNumber = tinNumber,
            FocalPersonName = focalPersonName,
            FocalPersonMobile = focalPersonMobile,
            PartnerId = partnerId,
            AuditInfo = "{}"
        };
    }
}
