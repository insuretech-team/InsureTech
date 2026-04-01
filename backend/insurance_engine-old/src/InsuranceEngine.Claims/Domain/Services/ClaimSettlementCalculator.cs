using System;

namespace InsuranceEngine.Claims.Domain.Services;

/// <summary>
/// Calculates deductible, co-pay, and net settlement amount for a claim.
/// Product-level configuration drives the rates.
/// </summary>
public class ClaimSettlementCalculator
{
    /// <summary>
    /// Calculate the breakdown of a claim settlement.
    /// </summary>
    /// <param name="claimedAmount">Amount claimed in paisa</param>
    /// <param name="approvedAmount">Amount approved by claim officer in paisa</param>
    /// <param name="deductiblePercentage">Product-level deductible rate (0-100)</param>
    /// <param name="coPayPercentage">Product-level co-pay rate (0-100)</param>
    /// <param name="maxDeductibleAmount">Product-level cap on deductible in paisa (0 = no cap)</param>
    /// <returns>Settlement breakdown</returns>
    public ClaimSettlementResult Calculate(
        long claimedAmount,
        long approvedAmount,
        double deductiblePercentage = 0,
        double coPayPercentage = 0,
        long maxDeductibleAmount = 0)
    {
        if (approvedAmount <= 0)
        {
            return new ClaimSettlementResult
            {
                ApprovedAmount = 0,
                DeductibleAmount = 0,
                CoPayAmount = 0,
                NetSettlementAmount = 0
            };
        }

        // 1. Calculate deductible (applied first)
        long deductible = (long)Math.Round(approvedAmount * (deductiblePercentage / 100.0), MidpointRounding.AwayFromZero);
        
        // Cap the deductible if a maximum is configured
        if (maxDeductibleAmount > 0 && deductible > maxDeductibleAmount)
        {
            deductible = maxDeductibleAmount;
        }

        // 2. Amount after deductible
        long afterDeductible = approvedAmount - deductible;

        // 3. Calculate co-pay (insured's share of the remainder)
        long coPay = (long)Math.Round(afterDeductible * (coPayPercentage / 100.0), MidpointRounding.AwayFromZero);

        // 4. Net settlement = what insurer pays
        long netSettlement = afterDeductible - coPay;

        return new ClaimSettlementResult
        {
            ApprovedAmount = approvedAmount,
            DeductibleAmount = deductible,
            CoPayAmount = coPay,
            NetSettlementAmount = Math.Max(0, netSettlement)
        };
    }
}

public class ClaimSettlementResult
{
    public long ApprovedAmount { get; set; }
    public long DeductibleAmount { get; set; }
    public long CoPayAmount { get; set; }
    public long NetSettlementAmount { get; set; }
}
