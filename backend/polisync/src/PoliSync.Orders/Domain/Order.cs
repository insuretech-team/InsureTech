using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Orders.Domain;

/// <summary>
/// Order aggregate root - manages the order lifecycle from quotation to policy issuance
/// </summary>
public sealed class Order : Entity
{
    public Guid TenantId { get; private set; }
    public string OrderNumber { get; private set; } = string.Empty;
    public Guid QuotationId { get; private set; }
    public Guid CustomerId { get; private set; }
    public Guid ProductId { get; private set; }
    public Guid PlanId { get; private set; }
    public OrderStatus Status { get; private set; }
    
    // Financial
    public long TotalPayable { get; private set; } // in paisa
    public string Currency { get; private set; } = "BDT";
    
    // Payment tracking
    public string? PaymentId { get; private set; }
    public string? PaymentGatewayRef { get; private set; }
    public OrderPaymentStatus PaymentStatus { get; private set; }
    
    // Policy tracking
    public string? PolicyId { get; private set; }
    
    // Cancellation/Failure
    public string? CancellationReason { get; private set; }
    public string? FailureReason { get; private set; }
    
    // Timestamps
    public DateTime? PaymentDueAt { get; private set; }
    public DateTime? CoverageStartAt { get; private set; }
    public DateTime? CoverageEndAt { get; private set; }

    private Order() { } // EF Core

    /// <summary>
    /// Creates a new order from an approved quotation
    /// </summary>
    public static Result<Order> Create(
        Guid tenantId,
        Guid quotationId,
        Guid customerId,
        Guid productId,
        Guid planId,
        long totalPayable,
        string currency = "BDT",
        DateTime? paymentDueAt = null,
        DateTime? coverageStartAt = null,
        DateTime? coverageEndAt = null)
    {
        if (tenantId == Guid.Empty)
            return Result.Fail<Order>("INVALID_TENANT", "Tenant ID is required");

        if (quotationId == Guid.Empty)
            return Result.Fail<Order>("INVALID_QUOTATION", "Quotation ID is required");

        if (customerId == Guid.Empty)
            return Result.Fail<Order>("INVALID_CUSTOMER", "Customer ID is required");

        if (productId == Guid.Empty)
            return Result.Fail<Order>("INVALID_PRODUCT", "Product ID is required");

        if (planId == Guid.Empty)
            return Result.Fail<Order>("INVALID_PLAN", "Plan ID is required");

        if (totalPayable <= 0)
            return Result.Fail<Order>("INVALID_AMOUNT", "Total payable must be positive");

        var order = new Order
        {
            Id = Guid.NewGuid(),
            TenantId = tenantId,
            OrderNumber = GenerateOrderNumber(),
            QuotationId = quotationId,
            CustomerId = customerId,
            ProductId = productId,
            PlanId = planId,
            Status = OrderStatus.Pending,
            TotalPayable = totalPayable,
            Currency = currency,
            PaymentStatus = OrderPaymentStatus.Unpaid,
            PaymentDueAt = paymentDueAt ?? DateTime.UtcNow.AddMinutes(30),
            CoverageStartAt = coverageStartAt,
            CoverageEndAt = coverageEndAt,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        order.RaiseDomainEvent(new OrderCreatedEvent(
            order.Id,
            order.OrderNumber,
            order.QuotationId,
            order.CustomerId,
            order.TotalPayable));

        return Result.Ok(order);
    }

    /// <summary>
    /// Initiates payment for the order
    /// </summary>
    public Result InitiatePayment(string paymentId, string paymentGatewayRef)
    {
        if (Status != OrderStatus.Pending)
            return Result.Fail("INVALID_TRANSITION", $"Cannot initiate payment from {Status} status");

        if (string.IsNullOrWhiteSpace(paymentId))
            return Result.Fail("INVALID_PAYMENT_ID", "Payment ID is required");

        PaymentId = paymentId;
        PaymentGatewayRef = paymentGatewayRef;
        Status = OrderStatus.PaymentInitiated;
        PaymentStatus = OrderPaymentStatus.PaymentInProgress;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new OrderPaymentInitiatedEvent(Id, OrderNumber, paymentId));

        return Result.Ok();
    }

    /// <summary>
    /// Confirms payment completion
    /// </summary>
    public Result ConfirmPayment()
    {
        if (Status != OrderStatus.PaymentInitiated)
            return Result.Fail("INVALID_TRANSITION", $"Cannot confirm payment from {Status} status");

        if (string.IsNullOrWhiteSpace(PaymentId))
            return Result.Fail("NO_PAYMENT", "No payment has been initiated");

        Status = OrderStatus.Paid;
        PaymentStatus = OrderPaymentStatus.Paid;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new OrderPaymentConfirmedEvent(Id, OrderNumber, PaymentId!));

        return Result.Ok();
    }

    /// <summary>
    /// Sets the policy ID after policy issuance
    /// </summary>
    public Result SetPolicyId(string policyId)
    {
        if (Status != OrderStatus.Paid)
            return Result.Fail("INVALID_TRANSITION", $"Cannot set policy from {Status} status");

        if (string.IsNullOrWhiteSpace(policyId))
            return Result.Fail("INVALID_POLICY_ID", "Policy ID is required");

        PolicyId = policyId;
        Status = OrderStatus.PolicyIssued;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new OrderPolicyIssuedEvent(Id, OrderNumber, policyId));

        return Result.Ok();
    }

    /// <summary>
    /// Cancels the order
    /// </summary>
    public Result Cancel(string reason)
    {
        if (Status is OrderStatus.Paid or OrderStatus.PolicyIssued)
            return Result.Fail("INVALID_TRANSITION", $"Cannot cancel order in {Status} status");

        if (string.IsNullOrWhiteSpace(reason))
            return Result.Fail("INVALID_REASON", "Cancellation reason is required");

        Status = OrderStatus.Cancelled;
        CancellationReason = reason;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new OrderCancelledEvent(Id, OrderNumber, reason));

        return Result.Ok();
    }

    /// <summary>
    /// Marks the order as failed
    /// </summary>
    public Result Fail(string reason)
    {
        if (Status is OrderStatus.Paid or OrderStatus.PolicyIssued or OrderStatus.Cancelled)
            return Result.Fail("INVALID_TRANSITION", $"Cannot fail order in {Status} status");

        if (string.IsNullOrWhiteSpace(reason))
            return Result.Fail("INVALID_REASON", "Failure reason is required");

        Status = OrderStatus.Failed;
        FailureReason = reason;
        PaymentStatus = OrderPaymentStatus.Failed;
        UpdatedAt = DateTime.UtcNow;

        RaiseDomainEvent(new OrderFailedEvent(Id, OrderNumber, reason));

        return Result.Ok();
    }

    private static string GenerateOrderNumber()
    {
        return $"ORD-{Guid.NewGuid().ToString()[..8].ToUpper()}";
    }
}

/// <summary>
/// Order status enum
/// </summary>
public enum OrderStatus
{
    Pending = 1,
    PaymentInitiated = 2,
    Paid = 3,
    PolicyIssued = 4,
    Cancelled = 5,
    Failed = 6
}

/// <summary>
/// Order payment status enum
/// </summary>
public enum OrderPaymentStatus
{
    Unpaid = 1,
    PaymentInProgress = 2,
    Paid = 3,
    Failed = 4
}
