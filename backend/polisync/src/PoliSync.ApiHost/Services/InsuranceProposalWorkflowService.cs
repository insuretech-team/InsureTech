using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Payment.Services.V1;
using Microsoft.Extensions.Options;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Orders.Infrastructure;
using PoliSync.Policy.Domain;
using PoliSync.Policy.Infrastructure;
using PoliSync.Quotes.Infrastructure;
using PoliSync.Refund.Infrastructure;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;
using System.Text.Json.Serialization;
using ProtoMoney = Insuretech.Common.V1.Money;
using ProposalEntity = Insuretech.Policy.Entity.V1.InsuranceProposal;
using ProposalStatus = Insuretech.Policy.Entity.V1.ProposalStatus;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;

namespace PoliSync.ApiHost.Services;

public sealed class InsuranceProposalWorkflowService
{
    private readonly IOrderDataGateway _orderGateway;
    private readonly IQuotationDataGateway _quotationGateway;
    private readonly IPolicyDataGateway _policyGateway;
    private readonly IRefundPaymentGateway _refundGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<InsuranceProposalWorkflowService> _logger;
    private readonly string _proposalSubmittedTopic;
    private readonly string _proposalApprovedTopic;
    private readonly string _proposalRejectedTopic;
    private readonly string _proposalRefundInitiatedTopic;
    private readonly string _policyIssuedProjectionTopic;

    public InsuranceProposalWorkflowService(
        IOrderDataGateway orderGateway,
        IQuotationDataGateway quotationGateway,
        IPolicyDataGateway policyGateway,
        IRefundPaymentGateway refundGateway,
        IEventBus eventBus,
        IOptions<KafkaOptions> kafkaOptions,
        ILogger<InsuranceProposalWorkflowService> logger)
    {
        _orderGateway = orderGateway;
        _quotationGateway = quotationGateway;
        _policyGateway = policyGateway;
        _refundGateway = refundGateway;
        _eventBus = eventBus;
        _logger = logger;

        var topics = kafkaOptions.Value.Topics ?? [];
        _proposalSubmittedTopic = ResolveTopic(topics, "OrderProposalSubmitted", "insuretech.proposal.submitted.v1");
        _proposalApprovedTopic = ResolveTopic(topics, "OrderProposalApproved", "insuretech.insurance.v1.proposal.approved");
        _proposalRejectedTopic = ResolveTopic(topics, "OrderProposalRejected", "insuretech.insurance.v1.proposal.rejected");
        _proposalRefundInitiatedTopic = ResolveTopic(topics, "OrderProposalRefundInitiated", "insuretech.insurance.v1.proposal.refund_initiated");
        _policyIssuedProjectionTopic = ResolveTopic(
            topics,
            "OrderPolicyProjectionIssued",
            ResolveTopic(topics, "PolicyIssued", "insuretech.insurance.v1.policy.issued"));
    }

    public async Task<ProposalEntity> SubmitProposalForOrderAsync(
        string orderId,
        string? insurerId = null,
        string? correlationId = null,
        string? submissionPayload = null,
        long? totalPayableAmount = null,
        string? totalPayableCurrency = null,
        CancellationToken cancellationToken = default)
    {
        var existing = await TryGetProposalByOrderAsync(orderId, cancellationToken);
        if (existing is not null)
        {
            return existing;
        }

        var orderView = await _orderGateway.GetOrderAsync(orderId, cancellationToken)
            ?? throw new KeyNotFoundException($"Order {orderId} was not found.");
        var order = orderView.Order ?? throw new KeyNotFoundException($"Order {orderId} payload was empty.");

        if (!string.IsNullOrWhiteSpace(order.PolicyId))
        {
            throw new InvalidOperationException($"Order {orderId} already has an issued policy.");
        }

        if (string.IsNullOrWhiteSpace(order.QuotationId))
        {
            throw new InvalidOperationException($"Order {orderId} is missing quotation_id.");
        }

        var quotation = await _quotationGateway.GetQuotationAsync(order.QuotationId, cancellationToken)
            ?? throw new KeyNotFoundException($"Quotation {order.QuotationId} was not found.");

        if (quotation.Status != Insuretech.Policy.Entity.V1.QuotationStatus.Approved)
        {
            throw new InvalidOperationException(
                $"Quotation {quotation.QuotationId} must be approved before a proposal can be submitted.");
        }

        var resolvedInsurerId = FirstNonEmpty(insurerId, order.InsurerId);
        if (string.IsNullOrWhiteSpace(resolvedInsurerId))
        {
            throw new InvalidOperationException($"Order {orderId} is missing insurer_id.");
        }

        var tenantId = FirstNonEmpty(order.TenantId, quotation.BusinessId);
        var customerId = FirstNonEmpty(order.CustomerId, quotation.CreatedByUserId);
        var productId = FirstNonEmpty(order.ProductId, quotation.DepartmentId);
        var planId = FirstNonEmpty(order.PlanId, quotation.PlanId);

        if (string.IsNullOrWhiteSpace(tenantId) ||
            string.IsNullOrWhiteSpace(customerId) ||
            string.IsNullOrWhiteSpace(productId) ||
            string.IsNullOrWhiteSpace(planId))
        {
            throw new InvalidOperationException($"Order {orderId} is missing tenant/customer/product/plan data.");
        }

        var premium = ResolveMoney(totalPayableAmount, totalPayableCurrency, order.TotalPayable, quotation);
        var sumInsured = ResolveSumInsured(quotation, premium.Currency);
        var now = DateTime.UtcNow;

        var proposal = new ProposalEntity
        {
            ProposalId = Guid.NewGuid().ToString(),
            ProposalNumber = GenerateProposalNumber(),
            TenantId = tenantId,
            OrderId = order.OrderId,
            QuotationId = quotation.QuotationId,
            CustomerId = customerId,
            InsurerId = resolvedInsurerId,
            ProductId = productId,
            PlanId = planId,
            ProposedPremium = premium,
            ProposedSumInsured = sumInsured,
            Status = ProposalStatus.Submitted,
            SubmissionPayload = submissionPayload ?? string.Empty,
            CorrelationId = correlationId ?? order.CorrelationId ?? string.Empty,
            SubmittedAt = Timestamp.FromDateTime(now),
            CreatedAt = Timestamp.FromDateTime(now),
            UpdatedAt = Timestamp.FromDateTime(now)
        };

        var created = await _policyGateway.CreateInsuranceProposalAsync(proposal, cancellationToken);
        await _eventBus.PublishAsync(
            new InsuranceProposalSubmittedEvent(
                created.ProposalId,
                created.ProposalNumber,
                created.OrderId,
                created.QuotationId,
                created.CustomerId,
                created.InsurerId,
                created.ProductId,
                created.PlanId,
                created.CorrelationId,
                created.Status.ToString()),
            _proposalSubmittedTopic,
            cancellationToken);

        _logger.LogInformation(
            "Submitted insurance proposal {ProposalId} for order {OrderId} to insurer {InsurerId}",
            created.ProposalId,
            created.OrderId,
            created.InsurerId);

        return created;
    }

    public async Task<ProposalApprovalResult> ApproveProposalAsync(
        string proposalId,
        string reviewedByUserId,
        string? insurerResponsePayload,
        string? decisionReason,
        CancellationToken cancellationToken = default)
    {
        var proposal = await _policyGateway.GetInsuranceProposalAsync(proposalId, cancellationToken)
            ?? throw new KeyNotFoundException($"Proposal {proposalId} was not found.");

        if (!string.IsNullOrWhiteSpace(proposal.ApprovedPolicyId))
        {
            var existingPolicy = await _policyGateway.GetPolicyAsync(proposal.ApprovedPolicyId, cancellationToken);
            if (existingPolicy is not null)
            {
                return new ProposalApprovalResult(proposal, existingPolicy);
            }
        }

        if (proposal.Status == ProposalStatus.Rejected ||
            proposal.Status == ProposalStatus.RefundInitiated ||
            proposal.Status == ProposalStatus.Refunded ||
            proposal.Status == ProposalStatus.Cancelled)
        {
            throw new InvalidOperationException($"Proposal {proposalId} is already closed with status {proposal.Status}.");
        }

        var order = await _orderGateway.GetOrderAsync(proposal.OrderId, cancellationToken);
        var coverageStart = order?.Order?.CoverageStartAt?.ToDateTime() ?? DateTime.UtcNow.Date;
        var coverageEnd = order?.Order?.CoverageEndAt?.ToDateTime() ?? coverageStart.AddMonths(12);

        var aggregate = PolicyAggregate.Create(
            proposal.CustomerId,
            proposal.ProductId,
            proposal.QuotationId,
            proposal.ProposedPremium?.Amount ?? 0,
            proposal.ProposedSumInsured?.Amount ?? 0,
            12,
            coverageStart,
            coverageEnd);

        aggregate.IssuePolicy();
        aggregate.Policy.PartnerId = proposal.TenantId;
        aggregate.Policy.PaymentGatewayReference = order?.Order?.PaymentId ?? string.Empty;
        aggregate.Policy.PaymentFrequency = "ANNUAL";
        aggregate.Policy.ReceiptNumber = $"RCPT-{DateTime.UtcNow:yyyyMMdd}-{Random.Shared.Next(100000, 999999)}";
        aggregate.Policy.PremiumCurrency = proposal.ProposedPremium?.Currency ?? "BDT";
        aggregate.Policy.SumInsuredCurrency = proposal.ProposedSumInsured?.Currency ?? "BDT";
        aggregate.Policy.PremiumAmount = proposal.ProposedPremium ?? NewMoney(0);
        aggregate.Policy.SumInsured = proposal.ProposedSumInsured ?? NewMoney(0);
        aggregate.Policy.TotalPayable = proposal.ProposedPremium ?? NewMoney(0);
        aggregate.Policy.VatTax = NewMoney(0, aggregate.Policy.PremiumCurrency);
        aggregate.Policy.ServiceFee = NewMoney(0, aggregate.Policy.PremiumCurrency);
        aggregate.Policy.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        var createdPolicy = await _policyGateway.CreatePolicyAsync(aggregate.Policy, cancellationToken);

        var updatedProposal = proposal.Clone();
        updatedProposal.Status = ProposalStatus.Approved;
        updatedProposal.ApprovedPolicyId = createdPolicy.PolicyId;
        updatedProposal.DecisionReason = decisionReason ?? updatedProposal.DecisionReason;
        updatedProposal.InsurerResponsePayload = insurerResponsePayload ?? updatedProposal.InsurerResponsePayload;
        updatedProposal.ReviewedByUserId = reviewedByUserId;
        updatedProposal.ReviewedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        updatedProposal.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        updatedProposal = await _policyGateway.UpdateInsuranceProposalAsync(updatedProposal, cancellationToken);
        await _eventBus.PublishAsync(
            new InsuranceProposalApprovedEvent(
                updatedProposal.ProposalId,
                updatedProposal.OrderId,
                createdPolicy.PolicyId,
                updatedProposal.CustomerId,
                updatedProposal.InsurerId,
                reviewedByUserId,
                updatedProposal.DecisionReason ?? string.Empty,
                updatedProposal.Status.ToString()),
            _proposalApprovedTopic,
            cancellationToken);
        await _eventBus.PublishAsync(
            new PolicyIssuedProjectionEvent(
                updatedProposal.OrderId,
                createdPolicy.PolicyId,
                updatedProposal.ProposalId,
                updatedProposal.CustomerId,
                updatedProposal.ProductId),
            _policyIssuedProjectionTopic,
            cancellationToken);

        _logger.LogInformation(
            "Approved proposal {ProposalId} and issued policy {PolicyId}",
            updatedProposal.ProposalId,
            createdPolicy.PolicyId);

        return new ProposalApprovalResult(updatedProposal, createdPolicy);
    }

    public async Task<ProposalRejectionResult> RejectProposalAsync(
        string proposalId,
        string reviewedByUserId,
        string? insurerResponsePayload,
        string? decisionReason,
        CancellationToken cancellationToken = default)
    {
        var proposal = await _policyGateway.GetInsuranceProposalAsync(proposalId, cancellationToken)
            ?? throw new KeyNotFoundException($"Proposal {proposalId} was not found.");

        if (proposal.Status == ProposalStatus.Approved && !string.IsNullOrWhiteSpace(proposal.ApprovedPolicyId))
        {
            throw new InvalidOperationException($"Proposal {proposalId} has already been approved into a policy.");
        }

        if ((proposal.Status == ProposalStatus.RefundInitiated || proposal.Status == ProposalStatus.Refunded) &&
            !string.IsNullOrWhiteSpace(proposal.RefundId))
        {
            return new ProposalRejectionResult(
                proposal,
                new InitiateRefundResponse
                {
                    RefundId = proposal.RefundId,
                    Status = proposal.Status == ProposalStatus.Refunded ? "REFUNDED" : "REFUND_INITIATED"
                });
        }

        if (proposal.Status == ProposalStatus.Cancelled)
        {
            throw new InvalidOperationException($"Proposal {proposalId} has already been cancelled.");
        }

        var order = await _orderGateway.GetOrderAsync(proposal.OrderId, cancellationToken)
            ?? throw new KeyNotFoundException($"Order {proposal.OrderId} was not found for proposal {proposalId}.");
        var paymentId = order.Order?.PaymentId;
        if (string.IsNullOrWhiteSpace(paymentId))
        {
            throw new InvalidOperationException($"Order {proposal.OrderId} is missing payment_id for refund processing.");
        }

        var reason = string.IsNullOrWhiteSpace(decisionReason)
            ? "Proposal rejected by insurer"
            : decisionReason;

        var refundResponse = await _refundGateway.InitiateRefundAsync(
            paymentId,
            proposal.ProposedPremium ?? NewMoney(0),
            reason,
            reviewedByUserId,
            cancellationToken);

        var updatedProposal = proposal.Clone();
        updatedProposal.Status = IsRefundCompleted(refundResponse.Status)
            ? ProposalStatus.Refunded
            : ProposalStatus.RefundInitiated;
        updatedProposal.DecisionReason = reason;
        updatedProposal.InsurerResponsePayload = insurerResponsePayload ?? updatedProposal.InsurerResponsePayload;
        updatedProposal.ReviewedByUserId = reviewedByUserId;
        updatedProposal.ReviewedAt = Timestamp.FromDateTime(DateTime.UtcNow);
        updatedProposal.RefundId = refundResponse.RefundId;
        updatedProposal.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        updatedProposal = await _policyGateway.UpdateInsuranceProposalAsync(updatedProposal, cancellationToken);
        await _eventBus.PublishAsync(
            new InsuranceProposalRejectedEvent(
                updatedProposal.ProposalId,
                updatedProposal.OrderId,
                updatedProposal.CustomerId,
                updatedProposal.InsurerId,
                reviewedByUserId,
                updatedProposal.DecisionReason ?? string.Empty,
                updatedProposal.RefundId,
                updatedProposal.Status.ToString()),
            _proposalRejectedTopic,
            cancellationToken);
        await _eventBus.PublishAsync(
            new ProposalRefundInitiatedEvent(
                updatedProposal.ProposalId,
                updatedProposal.OrderId,
                refundResponse.RefundId,
                paymentId,
                updatedProposal.CustomerId,
                updatedProposal.InsurerId,
                proposal.ProposedPremium?.Amount ?? 0,
                proposal.ProposedPremium?.Currency ?? "BDT",
                refundResponse.Status ?? string.Empty),
            _proposalRefundInitiatedTopic,
            cancellationToken);

        _logger.LogInformation(
            "Rejected proposal {ProposalId} and initiated refund {RefundId}",
            updatedProposal.ProposalId,
            refundResponse.RefundId);

        return new ProposalRejectionResult(updatedProposal, refundResponse);
    }

    public async Task<ProposalEntity?> TryGetProposalByOrderAsync(string orderId, CancellationToken cancellationToken = default)
    {
        var proposals = await _policyGateway.ListInsuranceProposalsAsync(
            orderId,
            insurerId: null,
            customerId: null,
            status: null,
            page: 1,
            pageSize: 1,
            cancellationToken);

        return proposals.FirstOrDefault();
    }

    private static ProtoMoney ResolveMoney(
        long? explicitAmount,
        string? explicitCurrency,
        ProtoMoney? orderAmount,
        QuotationEntity quotation)
    {
        if (explicitAmount.GetValueOrDefault() > 0)
        {
            return NewMoney(explicitAmount!.Value, string.IsNullOrWhiteSpace(explicitCurrency) ? "BDT" : explicitCurrency);
        }

        if (orderAmount is not null && orderAmount.Amount > 0)
        {
            return NewMoney(orderAmount.Amount, string.IsNullOrWhiteSpace(orderAmount.Currency) ? "BDT" : orderAmount.Currency);
        }

        if (quotation.QuotedAmount is not null && quotation.QuotedAmount.Amount > 0)
        {
            return NewMoney(
                quotation.QuotedAmount.Amount,
                string.IsNullOrWhiteSpace(quotation.QuotedAmount.Currency) ? "BDT" : quotation.QuotedAmount.Currency);
        }

        if (quotation.EstimatedPremium is not null && quotation.EstimatedPremium.Amount > 0)
        {
            return NewMoney(
                quotation.EstimatedPremium.Amount,
                string.IsNullOrWhiteSpace(quotation.EstimatedPremium.Currency) ? "BDT" : quotation.EstimatedPremium.Currency);
        }

        return NewMoney(0);
    }

    private static ProtoMoney ResolveSumInsured(QuotationEntity quotation, string currency)
    {
        if (quotation.QuotedAmount is not null && quotation.QuotedAmount.Amount > 0)
        {
            return NewMoney(Math.Max(quotation.QuotedAmount.Amount * 10, 1_000_000), currency);
        }

        if (quotation.EstimatedPremium is not null && quotation.EstimatedPremium.Amount > 0)
        {
            return NewMoney(Math.Max(quotation.EstimatedPremium.Amount * 10, 1_000_000), currency);
        }

        return NewMoney(1_000_000, currency);
    }

    private static string GenerateProposalNumber()
        => $"PRP-{DateTime.UtcNow:yyyyMMdd}-{Guid.NewGuid():N}"[..24];

    private static string? FirstNonEmpty(params string?[] values)
        => values.FirstOrDefault(value => !string.IsNullOrWhiteSpace(value));

    private static ProtoMoney NewMoney(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };

    private static bool IsRefundCompleted(string? status)
    {
        if (string.IsNullOrWhiteSpace(status))
        {
            return false;
        }

        var normalized = status.Trim().Replace("-", "_", StringComparison.Ordinal).ToUpperInvariant();
        return normalized is "COMPLETED" or "COMPLETE" or "SUCCESS" or "SUCCEEDED" or "REFUNDED" or "REFUND_COMPLETED";
    }

    private static string ResolveTopic(
        IReadOnlyDictionary<string, string> topics,
        string key,
        string fallback)
        => topics.TryGetValue(key, out var topic) && !string.IsNullOrWhiteSpace(topic)
            ? topic
            : fallback;

    private sealed record InsuranceProposalSubmittedEvent(
        [property: JsonPropertyName("proposal_id")] string ProposalId,
        [property: JsonPropertyName("proposal_number")] string ProposalNumber,
        [property: JsonPropertyName("order_id")] string OrderId,
        [property: JsonPropertyName("quotation_id")] string QuotationId,
        [property: JsonPropertyName("customer_id")] string CustomerId,
        [property: JsonPropertyName("insurer_id")] string InsurerId,
        [property: JsonPropertyName("product_id")] string ProductId,
        [property: JsonPropertyName("plan_id")] string PlanId,
        [property: JsonPropertyName("correlation_id")] string CorrelationId,
        [property: JsonPropertyName("status")] string Status) : DomainEvent;

    private sealed record InsuranceProposalApprovedEvent(
        [property: JsonPropertyName("proposal_id")] string ProposalId,
        [property: JsonPropertyName("order_id")] string OrderId,
        [property: JsonPropertyName("policy_id")] string PolicyId,
        [property: JsonPropertyName("customer_id")] string CustomerId,
        [property: JsonPropertyName("insurer_id")] string InsurerId,
        [property: JsonPropertyName("reviewed_by_user_id")] string ReviewedByUserId,
        [property: JsonPropertyName("decision_reason")] string DecisionReason,
        [property: JsonPropertyName("status")] string Status) : DomainEvent;

    private sealed record InsuranceProposalRejectedEvent(
        [property: JsonPropertyName("proposal_id")] string ProposalId,
        [property: JsonPropertyName("order_id")] string OrderId,
        [property: JsonPropertyName("customer_id")] string CustomerId,
        [property: JsonPropertyName("insurer_id")] string InsurerId,
        [property: JsonPropertyName("reviewed_by_user_id")] string ReviewedByUserId,
        [property: JsonPropertyName("decision_reason")] string DecisionReason,
        [property: JsonPropertyName("refund_id")] string RefundId,
        [property: JsonPropertyName("status")] string Status) : DomainEvent;

    private sealed record ProposalRefundInitiatedEvent(
        [property: JsonPropertyName("proposal_id")] string ProposalId,
        [property: JsonPropertyName("order_id")] string OrderId,
        [property: JsonPropertyName("refund_id")] string RefundId,
        [property: JsonPropertyName("payment_id")] string PaymentId,
        [property: JsonPropertyName("customer_id")] string CustomerId,
        [property: JsonPropertyName("insurer_id")] string InsurerId,
        [property: JsonPropertyName("amount")] long Amount,
        [property: JsonPropertyName("currency")] string Currency,
        [property: JsonPropertyName("status")] string Status) : DomainEvent;

    private sealed record PolicyIssuedProjectionEvent(
        [property: JsonPropertyName("order_id")] string OrderId,
        [property: JsonPropertyName("policy_id")] string PolicyId,
        [property: JsonPropertyName("proposal_id")] string ProposalId,
        [property: JsonPropertyName("customer_id")] string CustomerId,
        [property: JsonPropertyName("product_id")] string ProductId) : DomainEvent;
}

public sealed record ProposalApprovalResult(ProposalEntity Proposal, Insuretech.Policy.Entity.V1.Policy Policy);

public sealed record ProposalRejectionResult(ProposalEntity Proposal, InitiateRefundResponse Refund);
