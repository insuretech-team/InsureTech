using InsuranceEngine.SharedKernel.Domain;
using System.Security.Cryptography;
using System.Text;

namespace InsuranceEngine.Claims.Domain;

public sealed class ClaimAggregate : AggregateRoot<Guid>
{
    public string ClaimNumber { get; private set; } = string.Empty;
    public Guid PolicyId { get; private set; }
    public string ClaimType { get; private set; } = string.Empty;
    public Money ClaimAmount { get; private set; } = default!;
    public Money? ApprovedAmount { get; private set; }
    public string Status { get; private set; } = "SUBMITTED";
    public string Description { get; private set; } = string.Empty;
    public string DocumentHash { get; private set; } = string.Empty;
    public DateTime CreatedAt { get; private set; }

    private ClaimAggregate(Guid id, string claimNumber, Guid policyId, string type, Money amount, string description)
    {
        Id = id;
        ClaimNumber = claimNumber;
        PolicyId = policyId;
        ClaimType = type;
        ClaimAmount = amount;
        Description = description;
        Status = "SUBMITTED";
        CreatedAt = DateTime.UtcNow;
    }

    public static ClaimAggregate Submit(Guid policyId, string type, decimal amount, string description, long sequenceNumber, string? documentContent = null)
    {
        // FR-083: Format: CLM-YYYY-XXXX-NNNNNN (collision-safe via DB sequence)
        var year = DateTime.UtcNow.Year;
        var monthDay = DateTime.UtcNow.ToString("MMdd");
        var seq = sequenceNumber.ToString().PadLeft(6, '0');
        var claimNumber = $"CLM-{year}-{monthDay}-{seq}";

        var claim = new ClaimAggregate(Guid.NewGuid(), claimNumber, policyId, type, Money.FromDecimal(amount), description);
        
        if (!string.IsNullOrEmpty(documentContent))
        {
            claim.GenerateDocumentHash(documentContent);
        }

        return claim;
    }

    public void GenerateDocumentHash(string content)
    {
        var bytes = Encoding.UTF8.GetBytes(content);
        var hash = SHA256.HashData(bytes);
        DocumentHash = Convert.ToHexString(hash);
    }

    public void Approve(decimal approvedAmount)
    {
        ApprovedAmount = Money.FromDecimal(approvedAmount);
        Status = "APPROVED";
    }

    public void RequestApproval()
    {
        // FR-086 & FR-542: Claims Approval Matrix (Tiered)
        decimal amount = ClaimAmount.ToDecimal();

        if (amount < 10000)
        {
            Status = "AUTO_APPROVED";
            ApprovedAmount = ClaimAmount;
        }
        else if (amount < 50000)
        {
            Status = "PENDING_MANAGER_APPROVAL"; // L2
        }
        else if (amount < 200000)
        {
            Status = "PENDING_JOINT_APPROVAL"; // L3: Admin + Focal Person
        }
        else
        {
            Status = "PENDING_BOARD_APPROVAL"; // L4: Board + Insurer
        }
    }

    public Money CalculateCoPayment(decimal deductibleAmount, decimal coInsurancePercentage)
    {
        // FR-100: Co-payment Engine
        decimal claimValue = ClaimAmount.ToDecimal();
        
        // 1. Subtract Deductible
        decimal balanceAfterDeductible = Math.Max(0, claimValue - deductibleAmount);
        
        // 2. Calculate Co-insurance (Patient's share)
        decimal coInsuranceAmount = balanceAfterDeductible * (coInsurancePercentage / 100m);
        
        // 3. Final Approved Amount (Deductible + Co-insurance is paid by user, rest by insurer)
        decimal insurerShare = balanceAfterDeductible - coInsuranceAmount;
        
        return Money.FromDecimal(insurerShare);
    }

    public void Reject(string reason)
    {
        Status = "REJECTED";
    }
}
