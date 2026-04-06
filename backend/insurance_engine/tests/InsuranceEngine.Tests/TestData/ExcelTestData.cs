using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Tests.TestData;

public static class ExcelTestData
{
    private static readonly string ExcelPath = Path.Combine(
        Directory.GetCurrentDirectory(), 
        "..", "..", "..", "..", "..", 
        "insurance_engine", "Test.xlsx");

    public static ExcelPolicyData GetOverseasMediclaimPolicy()
    {
        return new ExcelPolicyData
        {
            ProductId = "OMC-PRODUCT-001",
            CustomerId = "CUST-PRAGATI-001",
            ProductType = "OVERSEAS_MEDICLAIM",
            PolicyType = "INDIVIDUAL",
            PaymentMethod = "ONLINE",
            PremiumAmount = 1239,
            SumInsured = 50000,
            TenureMonths = 1,
            StartDate = new DateTime(2026, 4, 15),
            ProposerDetails = "Name: Md. Zubayed Ur Rahman, Address: N.B Tower Level-5, 40/7 North Avenue, Gulshan-2, Dhaka-1212, Mobile: 01985700011, Email: Zubayer@ymail.com, Occupation: Service at Medland Bank Plc, Passport: GA-18-6525, Plan: Business & Holiday (14-180 days)"
        };
    }

    public static ExcelPolicyData GetVehicleInsurancePolicy()
    {
        return new ExcelPolicyData
        {
            ProductId = "VEHICLE-PRODUCT-001",
            CustomerId = "CUST-PRAGATI-002",
            ProductType = "VEHICLE_INSURANCE",
            PolicyType = "INDIVIDUAL",
            PaymentMethod = "ONLINE",
            PremiumAmount = 5000,
            SumInsured = 500000,
            TenureMonths = 12,
            StartDate = new DateTime(2026, 4, 15),
            ProposerDetails = "Vehicle: Dhaka Metro GA-18-6525, Chassis: NKE165-7216292, Engine: G4NAEM48921, Make: Hyundai, Model: Tucson, Year: 2024"
        };
    }

    public static NomineeData GetBeneficiaryNominee()
    {
        return new NomineeData
        {
            FullName = "Fatema Begum",
            Relationship = "Wife",
            SharePercentage = 50,
            DateOfBirth = new DateTime(1985, 6, 15),
            NidNumber = "198515678900001",
            PhoneNumber = "+88017111234567",
            NomineeDobText = "15-06-1985"
        };
    }

    public static List<PremiumRateData> GetPremiumRates()
    {
        return new List<PremiumRateData>
        {
            new() { AgeGroup = "0-40", PeriodDays = 14, Plan = "A", Premium = 1239 },
            new() { AgeGroup = "41-50", PeriodDays = 14, Plan = "A", Premium = 2499 },
            new() { AgeGroup = "51-55", PeriodDays = 14, Plan = "A", Premium = 2499 },
            new() { AgeGroup = "56-60", PeriodDays = 14, Plan = "A", Premium = 4999 },
            new() { AgeGroup = "61-65", PeriodDays = 14, Plan = "A", Premium = 8131 },
        };
    }
}

public class ExcelPolicyData
{
    public string ProductId { get; set; } = "";
    public string CustomerId { get; set; } = "";
    public string? ProductType { get; set; }
    public string? PolicyType { get; set; }
    public string? PaymentMethod { get; set; }
    public decimal PremiumAmount { get; set; }
    public decimal SumInsured { get; set; }
    public int TenureMonths { get; set; }
    public DateTime StartDate { get; set; }
    public string? ProposerDetails { get; set; }
    public List<NomineeData>? Nominees { get; set; }
}

public class NomineeData
{
    public string FullName { get; set; } = "";
    public string Relationship { get; set; } = "";
    public int SharePercentage { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public string? NidNumber { get; set; }
    public string? PhoneNumber { get; set; }
    public string? NomineeDobText { get; set; }
}

public class PremiumRateData
{
    public string AgeGroup { get; set; } = "";
    public int PeriodDays { get; set; }
    public string Plan { get; set; } = "";
    public decimal Premium { get; set; }
}
