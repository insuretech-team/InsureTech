using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Quotes.Domain;

/// <summary>
/// Quotation aggregate root - manages the quotation lifecycle and premium calculations
/// </summary>
public sealed class Quotation : Entity
{
    public Guid TenantId { get; private set; }
    public string QuotationNumber { get; private set; } = string.Empty;
    public Guid ProductId { get; private set; }
    public Guid PlanId { get; private set; }
    public Guid CustomerId { get; private set; }
    public QuotationStatus Status { get; private set; }
    public DateTime ExpiryDate { get; private set; }

    // Premium breakdown - all in paisa (int64)
    public long BasePremium { get; private set; }
    public long RiderPremium { get; private set; }
    public long LoadingAmount { get; private set; }
    public long DiscountAmount { get; private set; }
    public long VatTax { get; private set; }
    public long ServiceFee { get; private set; }
    public long TotalPayable { get; private set; }

    public string? RejectionReason { get; private set; }

    private Quotation() { } // EF Core

    internal static Quotation Rehydrate(
        Guid id,
        Guid tenantId,
        string quotationNumber,
        Guid productId,
        Guid planId,
        Guid customerId,
        QuotationStatus status,
        DateTime expiryDate,
        long basePremium,
        long riderPremium,
        long loadingAmount,
        long discountAmount,
        long vatTax,
        long serviceFee,
        long totalPayable,
        string? rejectionReason,
        DateTime createdAt,
        DateTime? updatedAt)
    {
        return new Quotation
        {
            Id = id,
            TenantId = tenantId,
            QuotationNumber = quotationNumber,
            ProductId = productId,
            PlanId = planId,
            CustomerId = customerId,
            Status = status,
            ExpiryDate = expiryDate,
            BasePremium = basePremium,
            RiderPremium = riderPremium,
            LoadingAmount = loadingAmount,
            DiscountAmount = discountAmount,
            VatTax = vatTax,
            ServiceFee = serviceFee,
            TotalPayable = totalPayable,
            RejectionReason = rejectionReason,
            CreatedAt = createdAt,
            UpdatedAt = updatedAt
        };
    }

    /// <summary>
    /// Creates a new quotation in DRAFT status
    /// </summary>
    public static Result<Quotation> Create(
        Guid tenantId,
        Guid productId,
        Guid planId,
        Guid customerId,
        long basePremium,
        long riderPremium,
        int expiryDays = 30)
    {
        if (productId == Guid.Empty)
            return Result.Fail<Quotation>("INVALID_PRODUCT", "Product ID is required");

        if (planId == Guid.Empty)
            return Result.Fail<Quotation>("INVALID_PLAN", "Plan ID is required");

        if (customerId == Guid.Empty)
            return Result.Fail<Quotation>("INVALID_CUSTOMER", "Customer ID is required");

        if (basePremium <= 0)
            return Result.Fail<Quotation>("INVALID_PREMIUM", "Base premium must be positive");

        if (expiryDays <= 0 || expiryDays > 90)
            return Result.Fail<Quotation>("INVALID_EXPIRY", "Expiry days must be between 1 and 90");

        var quotation = new Quotation
        {
            Id = Guid.NewGuid(),
            TenantId = tenantId,
            QuotationNumber = GenerateQuotationNumber(),
            ProductId = productId,
            PlanId = planId,
            CustomerId = customerId,
            Status = QuotationStatus.Draft,
            BasePremium = basePremium,
            RiderPremium = riderPremium,
            LoadingAmount = 0,
            DiscountAmount = 0,
            ServiceFee = 0,
            ExpiryDate = DateTime.UtcNow.AddDays(expiryDays),
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        quotation.RecalculateTotal();
        quotation.RaiseDomainEvent(new QuotationCreatedEvent(quotation.Id, quotation.QuotationNumber));

        return Result.Ok(quotation);
    }

    /// <summary>
    /// Submits the quotation for underwriting review
    /// </summary>
    public Result Submit()
    {
        if (Status != QuotationStatus.Draft)
            return Result.Fail("INVALID_TRANSITION", $"Cannot submit quotation from {Status} status");

        if (DateTime.UtcNow > ExpiryDate)
            return Result.Fail("EXPIRED", "Quotation has expired");

        Status = QuotationStatus.Submitted;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new QuotationSubmittedEvent(Id, QuotationNumber, ProductId, CustomerId));

        return Result.Ok();
    }

    /// <summary>
    /// Marks quotation as received by underwriting
    /// </summary>
    public Result MarkAsReceived()
    {
        if (Status != QuotationStatus.Submitted)
            return Result.Fail("INVALID_TRANSITION", $"Cannot mark as received from {Status} status");

        Status = QuotationStatus.Received;
        UpdatedAt = DateTime.UtcNow;

        return Result.Ok();
    }

    /// <summary>
    /// Approves the quotation after underwriting review
    /// </summary>
    public Result Approve()
    {
        if (Status is not (QuotationStatus.Submitted or QuotationStatus.Received))
            return Result.Fail("INVALID_TRANSITION", $"Cannot approve quotation from {Status} status");

        if (DateTime.UtcNow > ExpiryDate)
            return Result.Fail("EXPIRED", "Quotation has expired");

        Status = QuotationStatus.Approved;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new QuotationApprovedEvent(Id, QuotationNumber, TotalPayable));

        return Result.Ok();
    }

    /// <summary>
    /// Rejects the quotation with a reason
    /// </summary>
    public Result Reject(string reason)
    {
        if (string.IsNullOrWhiteSpace(reason))
            return Result.Fail("INVALID_REASON", "Rejection reason is required");

        if (Status is not (QuotationStatus.Submitted or QuotationStatus.Received))
            return Result.Fail("INVALID_TRANSITION", $"Cannot reject quotation from {Status} status");

        Status = QuotationStatus.Rejected;
        RejectionReason = reason;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new QuotationRejectedEvent(Id, QuotationNumber, reason));

        return Result.Ok();
    }

    /// <summary>
    /// Expires the quotation (called by background job)
    /// </summary>
    public Result Expire()
    {
        if (Status is QuotationStatus.Approved or QuotationStatus.Rejected or QuotationStatus.Expired)
            return Result.Fail("TERMINAL", "Quotation is already in a terminal state");

        Status = QuotationStatus.Expired;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new QuotationExpiredEvent(Id, QuotationNumber));

        return Result.Ok();
    }

    /// <summary>
    /// Applies underwriting loading to the premium
    /// </summary>
    public Result ApplyLoading(long loadingAmount)
    {
        if (loadingAmount < 0)
            return Result.Fail("INVALID_LOADING", "Loading amount cannot be negative");

        LoadingAmount = loadingAmount;
        RecalculateTotal();
        UpdatedAt = DateTime.UtcNow;

        return Result.Ok();
    }

    /// <summary>
    /// Applies a discount to the premium
    /// </summary>
    public Result ApplyDiscount(long discountAmount)
    {
        if (discountAmount < 0)
            return Result.Fail("INVALID_DISCOUNT", "Discount amount cannot be negative");

        if (discountAmount > (BasePremium + RiderPremium + LoadingAmount))
            return Result.Fail("EXCESSIVE_DISCOUNT", "Discount cannot exceed total premium");

        DiscountAmount = discountAmount;
        RecalculateTotal();
        UpdatedAt = DateTime.UtcNow;

        return Result.Ok();
    }

    /// <summary>
    /// Sets the service fee
    /// </summary>
    public Result SetServiceFee(long serviceFee)
    {
        if (serviceFee < 0)
            return Result.Fail("INVALID_FEE", "Service fee cannot be negative");

        ServiceFee = serviceFee;
        RecalculateTotal();
        UpdatedAt = DateTime.UtcNow;

        return Result.Ok();
    }

    /// <summary>
    /// Recalculates the total payable amount including VAT
    /// </summary>
    private void RecalculateTotal()
    {
        var subtotal = BasePremium + RiderPremium + LoadingAmount - DiscountAmount;
        
        // 15% VAT on subtotal (Bangladesh standard rate)
        VatTax = (long)(subtotal * 0.15m);
        
        TotalPayable = subtotal + VatTax + ServiceFee;
    }

    /// <summary>
    /// Generates a unique quotation number in format QT-XXXXXXXX
    /// </summary>
    private static string GenerateQuotationNumber()
    {
        return $"QT-{Guid.NewGuid().ToString()[..8].ToUpper()}";
    }
}

/// <summary>
/// Quotation status enum matching proto definition
/// </summary>
public enum QuotationStatus
{
    Draft = 1,
    Submitted = 2,
    Received = 3,
    Approved = 4,
    Rejected = 5,
    Expired = 6
}
