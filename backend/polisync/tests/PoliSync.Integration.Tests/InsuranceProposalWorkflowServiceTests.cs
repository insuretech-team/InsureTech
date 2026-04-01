using FluentAssertions;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Orders.Entity.V1;
using Insuretech.Orders.Services.V1;
using Insuretech.Payment.Services.V1;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging.Abstractions;
using Microsoft.Extensions.Options;
using PoliSync.ApiHost.BackgroundServices;
using PoliSync.ApiHost.Services;
using PoliSync.Infrastructure.Messaging;
using PoliSync.Orders.Infrastructure;
using PoliSync.Policy.Infrastructure;
using PoliSync.Quotes.Infrastructure;
using PoliSync.Refund.Infrastructure;
using PoliSync.SharedKernel.Messaging;
using System.Reflection;
using Xunit;
using DomainEventBase = PoliSync.SharedKernel.Domain.DomainEvent;
using DomainQuotation = PoliSync.Quotes.Domain.Quotation;
using DomainQuotationStatus = PoliSync.Quotes.Domain.QuotationStatus;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;
using ProposalEntity = Insuretech.Policy.Entity.V1.InsuranceProposal;
using ProposalStatus = Insuretech.Policy.Entity.V1.ProposalStatus;
using ProtoMoney = Insuretech.Common.V1.Money;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;
using QuotationStatus = Insuretech.Policy.Entity.V1.QuotationStatus;
using OrderInitiatePaymentResponse = Insuretech.Orders.Services.V1.InitiatePaymentResponse;

namespace PoliSync.Integration.Tests;

public sealed class InsuranceProposalWorkflowServiceTests
{
    [Fact]
    public async Task SubmitProposalForOrderAsync_PublishesSubmissionEventOnce()
    {
        var harness = CreateHarness();

        var created = await harness.Service.SubmitProposalForOrderAsync("order-1");
        var duplicate = await harness.Service.SubmitProposalForOrderAsync("order-1");

        created.ProposalId.Should().NotBeNullOrWhiteSpace();
        duplicate.ProposalId.Should().Be(created.ProposalId);
        harness.EventBus.Published.Should().ContainSingle();
        harness.EventBus.Published[0].Topic.Should().Be("insuretech.proposal.submitted.v1");
        harness.EventBus.Published[0].Event.EventType.Should().Be("InsuranceProposalSubmittedEvent");
    }

    [Fact]
    public async Task ApproveProposalAsync_PublishesApprovalAndPolicyProjectionEvents()
    {
        var harness = CreateHarness();
        var proposal = await harness.Service.SubmitProposalForOrderAsync("order-1");
        harness.EventBus.Clear();

        var result = await harness.Service.ApproveProposalAsync(
            proposal.ProposalId,
            reviewedByUserId: "reviewer-1",
            insurerResponsePayload: "{\"decision\":\"approved\"}",
            decisionReason: "Approved by insurer");

        result.Proposal.Status.Should().Be(ProposalStatus.Approved);
        result.Policy.PolicyId.Should().NotBeNullOrWhiteSpace();
        harness.EventBus.Published.Should().HaveCount(2);
        harness.EventBus.Published.Select(x => x.Topic).Should().BeEquivalentTo(
        [
            "insuretech.insurance.v1.proposal.approved",
            "insuretech.insurance.v1.policy.issued"
        ]);
        harness.EventBus.Published.Select(x => x.Event.EventType).Should().BeEquivalentTo(
        [
            "InsuranceProposalApprovedEvent",
            "PolicyIssuedProjectionEvent"
        ]);
    }

    [Fact]
    public async Task RejectProposalAsync_PublishesRejectionAndRefundEvents()
    {
        var harness = CreateHarness();
        var proposal = await harness.Service.SubmitProposalForOrderAsync("order-1");
        harness.EventBus.Clear();

        var result = await harness.Service.RejectProposalAsync(
            proposal.ProposalId,
            reviewedByUserId: "reviewer-2",
            insurerResponsePayload: "{\"decision\":\"rejected\"}",
            decisionReason: "Risk outside appetite");

        result.Proposal.Status.Should().Be(ProposalStatus.RefundInitiated);
        result.Proposal.RefundId.Should().Be("refund-1");
        harness.EventBus.Published.Should().HaveCount(2);
        harness.EventBus.Published.Select(x => x.Topic).Should().BeEquivalentTo(
        [
            "insuretech.insurance.v1.proposal.rejected",
            "insuretech.insurance.v1.proposal.refund_initiated"
        ]);
        harness.EventBus.Published.Select(x => x.Event.EventType).Should().BeEquivalentTo(
        [
            "InsuranceProposalRejectedEvent",
            "ProposalRefundInitiatedEvent"
        ]);
    }

    [Fact]
    public async Task RejectProposalAsync_IsIdempotentForExistingRefund()
    {
        var harness = CreateHarness();
        var proposal = await harness.Service.SubmitProposalForOrderAsync("order-1");

        await harness.Service.RejectProposalAsync(
            proposal.ProposalId,
            reviewedByUserId: "reviewer-2",
            insurerResponsePayload: "{\"decision\":\"rejected\"}",
            decisionReason: "Risk outside appetite");

        harness.EventBus.Clear();
        var duplicate = await harness.Service.RejectProposalAsync(
            proposal.ProposalId,
            reviewedByUserId: "reviewer-2",
            insurerResponsePayload: "{\"decision\":\"rejected\"}",
            decisionReason: "Risk outside appetite");

        duplicate.Proposal.RefundId.Should().Be("refund-1");
        duplicate.Proposal.Status.Should().Be(ProposalStatus.RefundInitiated);
        harness.RefundGateway.Calls.Should().Be(1);
        harness.EventBus.Published.Should().BeEmpty();
    }

    [Fact]
    public async Task InsuranceProposalDecisionConsumer_ProcessesApprovalEvents()
    {
        var harness = CreateHarness();
        var proposal = await harness.Service.SubmitProposalForOrderAsync("order-1");
        harness.EventBus.Clear();

        using var consumer = CreateDecisionConsumer(harness);
        var processed = await InvokeDecisionConsumerAsync(
            consumer,
            "insuretech.proposal.approved.v1",
            $$"""{"proposal_id":"{{proposal.ProposalId}}","reviewed_by_user_id":"insurer-user","decision_reason":"Approved asynchronously"}""");

        processed.Should().BeTrue();
        var updatedProposal = await harness.PolicyGateway.GetInsuranceProposalAsync(proposal.ProposalId);
        updatedProposal!.Status.Should().Be(ProposalStatus.Approved);
        harness.EventBus.Published.Select(x => x.Topic).Should().Contain(
        [
            "insuretech.insurance.v1.proposal.approved",
            "insuretech.insurance.v1.policy.issued"
        ]);
    }

    [Fact]
    public async Task InsuranceProposalDecisionConsumer_ProcessesRejectionEvents()
    {
        var harness = CreateHarness();
        var proposal = await harness.Service.SubmitProposalForOrderAsync("order-1");
        harness.EventBus.Clear();

        using var consumer = CreateDecisionConsumer(harness);
        var processed = await InvokeDecisionConsumerAsync(
            consumer,
            "insuretech.proposal.rejected.v1",
            $$"""{"proposal_id":"{{proposal.ProposalId}}","reviewed_by_user_id":"insurer-user","decision_reason":"Rejected asynchronously"}""");

        processed.Should().BeTrue();
        var updatedProposal = await harness.PolicyGateway.GetInsuranceProposalAsync(proposal.ProposalId);
        updatedProposal!.Status.Should().Be(ProposalStatus.RefundInitiated);
        harness.RefundGateway.Calls.Should().Be(1);
    }

    private static TestHarness CreateHarness()
    {
        var orderGateway = new FakeOrderGateway();
        var quotationGateway = new FakeQuotationGateway();
        var policyGateway = new FakePolicyGateway();
        var refundGateway = new FakeRefundGateway();
        var eventBus = new RecordingEventBus();

        var topics = Options.Create(new KafkaOptions
        {
            Topics = new Dictionary<string, string>
            {
                ["OrderProposalSubmitted"] = "insuretech.proposal.submitted.v1",
                ["OrderProposalApproved"] = "insuretech.insurance.v1.proposal.approved",
                ["OrderProposalRejected"] = "insuretech.insurance.v1.proposal.rejected",
                ["OrderProposalRefundInitiated"] = "insuretech.insurance.v1.proposal.refund_initiated",
                ["OrderPolicyProjectionIssued"] = "insuretech.insurance.v1.policy.issued"
            }
        });

        var service = new InsuranceProposalWorkflowService(
            orderGateway,
            quotationGateway,
            policyGateway,
            refundGateway,
            eventBus,
            topics,
            NullLogger<InsuranceProposalWorkflowService>.Instance);

        return new TestHarness(service, eventBus, refundGateway, policyGateway);
    }

    private static InsuranceProposalDecisionConsumer CreateDecisionConsumer(TestHarness harness)
    {
        var configuration = new ConfigurationBuilder()
            .AddInMemoryCollection(new Dictionary<string, string?>
            {
                ["Kafka:BootstrapServers"] = "localhost:9092",
                ["Kafka:Topics:InsurerProposalApprovedInbound"] = "insuretech.proposal.approved.v1",
                ["Kafka:Topics:InsurerProposalRejectedInbound"] = "insuretech.proposal.rejected.v1"
            })
            .Build();

        var services = new ServiceCollection();
        services.AddScoped(_ => harness.Service);

        var provider = services.BuildServiceProvider();
        return new InsuranceProposalDecisionConsumer(
            configuration,
            NullLogger<InsuranceProposalDecisionConsumer>.Instance,
            provider.GetRequiredService<IServiceScopeFactory>());
    }

    private static async Task<bool> InvokeDecisionConsumerAsync(
        InsuranceProposalDecisionConsumer consumer,
        string topic,
        string payload)
    {
        var method = typeof(InsuranceProposalDecisionConsumer).GetMethod(
            "TryProcessMessageAsync",
            BindingFlags.Instance | BindingFlags.NonPublic);

        method.Should().NotBeNull();
        var task = (Task<bool>)method!.Invoke(consumer, [topic, payload, CancellationToken.None])!;
        return await task;
    }
}

public sealed record TestHarness(
    InsuranceProposalWorkflowService Service,
    RecordingEventBus EventBus,
    FakeRefundGateway RefundGateway,
    FakePolicyGateway PolicyGateway);

public sealed class FakeOrderGateway : IOrderDataGateway
{
    private readonly Dictionary<string, OrderView> _orders = new()
    {
        ["order-1"] = new OrderView
        {
            Order = new Order
            {
                OrderId = "order-1",
                QuotationId = "quote-1",
                TenantId = "tenant-1",
                CustomerId = "customer-1",
                ProductId = "product-1",
                PlanId = "plan-1",
                InsurerId = "insurer-1",
                CorrelationId = "corr-1",
                PaymentId = "payment-1",
                TotalPayable = TestMoney.New(250_000),
                CoverageStartAt = Timestamp.FromDateTime(DateTime.UtcNow.Date.ToUniversalTime()),
                CoverageEndAt = Timestamp.FromDateTime(DateTime.UtcNow.Date.AddYears(1).ToUniversalTime())
            }
        }
    };

    public Task<CreateOrderResponse> CreateOrderAsync(
        string quotationId,
        string customerId,
        string paymentMethod,
        CancellationToken cancellationToken = default,
        string? productId = null,
        string? planId = null,
        long totalPayable = 0,
        string currency = "BDT")
    {
        var orderId = $"order-{_orders.Count + 1}";
        var created = new OrderView
        {
            Order = new Order
            {
                OrderId = orderId,
                QuotationId = quotationId,
                CustomerId = customerId,
                ProductId = productId ?? string.Empty,
                PlanId = planId ?? string.Empty,
                TotalPayable = TestMoney.New(totalPayable, currency),
                CoverageStartAt = Timestamp.FromDateTime(DateTime.UtcNow),
                CoverageEndAt = Timestamp.FromDateTime(DateTime.UtcNow.AddYears(1))
            }
        };

        _orders[orderId] = created;

        return Task.FromResult(new CreateOrderResponse
        {
            Order = created,
            Message = "created"
        });
    }

    public Task<OrderView?> GetOrderAsync(string orderId, CancellationToken cancellationToken = default)
        => Task.FromResult(_orders.TryGetValue(orderId, out var order) ? order : null);

    public Task<ListOrdersResponse> ListOrdersAsync(ListOrdersRequest request, CancellationToken cancellationToken = default)
    {
        var orders = _orders.Values
            .Where(x => string.IsNullOrWhiteSpace(request.CustomerId) || x.Order.CustomerId == request.CustomerId)
            .ToList();

        var response = new ListOrdersResponse { TotalCount = orders.Count };
        response.Orders.AddRange(orders);
        return Task.FromResult(response);
    }

    public Task<OrderInitiatePaymentResponse> InitiatePaymentAsync(
        string orderId,
        string paymentMethod,
        string callbackUrl,
        string idempotencyKey,
        CancellationToken cancellationToken = default)
        => Task.FromResult(new OrderInitiatePaymentResponse
        {
            OrderId = orderId,
            PaymentId = $"payment-{orderId}",
            PaymentUrl = callbackUrl,
            PaymentGatewayRef = idempotencyKey,
            ExpiresAt = Timestamp.FromDateTime(DateTime.UtcNow.AddMinutes(15))
        });

    public Task<ConfirmPaymentResponse> ConfirmPaymentAsync(
        string orderId,
        string paymentId,
        string transactionId,
        CancellationToken cancellationToken = default)
        => Task.FromResult(new ConfirmPaymentResponse
        {
            OrderId = orderId,
            Status = OrderStatus.Paid,
            Message = transactionId
        });

    public Task<CancelOrderResponse> CancelOrderAsync(string orderId, string reason, CancellationToken cancellationToken = default)
        => Task.FromResult(new CancelOrderResponse
        {
            OrderId = orderId,
            Status = OrderStatus.Cancelled,
            Message = reason
        });

    public Task<GetOrderStatusResponse?> GetOrderStatusAsync(string orderId, CancellationToken cancellationToken = default)
        => Task.FromResult<GetOrderStatusResponse?>(_orders.TryGetValue(orderId, out var order)
            ? new GetOrderStatusResponse
            {
                OrderId = orderId,
                Status = order.Order.Status,
                PaymentId = order.Order.PaymentId,
                PolicyId = order.Order.PolicyId
            }
            : null);
}

public sealed class FakeQuotationGateway : IQuotationDataGateway
{
    private readonly Dictionary<string, QuotationEntity> _quotations = new()
    {
        ["quote-1"] = new QuotationEntity
        {
            QuotationId = "quote-1",
            BusinessId = "tenant-1",
            CreatedByUserId = "customer-1",
            DepartmentId = "product-1",
            PlanId = "plan-1",
            Status = QuotationStatus.Approved,
            QuotedAmount = TestMoney.New(250_000),
            EstimatedPremium = TestMoney.New(250_000)
        }
    };

    private readonly Dictionary<Guid, DomainQuotation> _domainQuotations = new();

    public Task CreateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default)
    {
        _domainQuotations[quotation.Id] = quotation;
        return Task.CompletedTask;
    }

    public Task<DomainQuotation?> GetByIdAsync(Guid quotationId, CancellationToken cancellationToken = default)
        => Task.FromResult(_domainQuotations.TryGetValue(quotationId, out var quotation) ? quotation : null);

    public Task UpdateAsync(DomainQuotation quotation, CancellationToken cancellationToken = default)
    {
        _domainQuotations[quotation.Id] = quotation;
        return Task.CompletedTask;
    }

    public Task<(IReadOnlyList<DomainQuotation> Quotations, int TotalCount)> ListAsync(
        Guid tenantId,
        Guid? customerId,
        DomainQuotationStatus? status,
        int pageNumber,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        var quotations = _domainQuotations.Values
            .Where(q => tenantId == Guid.Empty || q.TenantId == tenantId)
            .Where(q => customerId is null || q.CustomerId == customerId.Value)
            .Where(q => status is null || q.Status == status.Value)
            .Skip(Math.Max(pageNumber - 1, 0) * pageSize)
            .Take(pageSize)
            .ToList();

        return Task.FromResult<(IReadOnlyList<DomainQuotation>, int)>((quotations, quotations.Count));
    }

    public Task<IReadOnlyList<DomainQuotation>> GetExpiredQuotationsAsync(CancellationToken cancellationToken = default)
        => Task.FromResult<IReadOnlyList<DomainQuotation>>(
            _domainQuotations.Values.Where(q => q.ExpiryDate < DateTime.UtcNow).ToList());

    public Task<QuotationEntity> CreateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrWhiteSpace(quotation.QuotationId))
        {
            quotation.QuotationId = Guid.NewGuid().ToString();
        }

        _quotations[quotation.QuotationId] = quotation;
        return Task.FromResult(quotation);
    }

    public Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
        => Task.FromResult(_quotations.TryGetValue(quotationId, out var quotation) ? quotation : null);

    public Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
    {
        _quotations[quotation.QuotationId] = quotation;
        return Task.FromResult(quotation);
    }

    public Task<IReadOnlyList<QuotationEntity>> ListQuotationsAsync(
        string businessId,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        var quotations = _quotations.Values
            .Where(q => string.IsNullOrWhiteSpace(businessId) || q.BusinessId == businessId)
            .Skip(Math.Max(page - 1, 0) * pageSize)
            .Take(pageSize)
            .ToList();

        return Task.FromResult<IReadOnlyList<QuotationEntity>>(quotations);
    }

    public Task DeleteQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
    {
        _quotations.Remove(quotationId);
        return Task.CompletedTask;
    }
}

public sealed class FakePolicyGateway : IPolicyDataGateway
{
    private readonly Dictionary<string, PolicyEntity> _policies = new();
    private readonly Dictionary<string, ProposalEntity> _proposals = new();

    public Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
    {
        _policies[policy.PolicyId] = policy;
        return Task.FromResult(policy);
    }

    public Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken cancellationToken = default)
        => Task.FromResult(_policies.TryGetValue(policyId, out var policy) ? policy : null);

    public Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
    {
        _policies[policy.PolicyId] = policy;
        return Task.FromResult(policy);
    }

    public Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken cancellationToken = default)
        => Task.FromResult<IReadOnlyList<PolicyEntity>>(_policies.Values.ToList());

    public Task DeletePolicyAsync(string policyId, CancellationToken cancellationToken = default)
    {
        _policies.Remove(policyId);
        return Task.CompletedTask;
    }

    public Task<ProposalEntity> CreateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default)
    {
        _proposals[proposal.ProposalId] = proposal;
        return Task.FromResult(proposal);
    }

    public Task<ProposalEntity?> GetInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default)
        => Task.FromResult(_proposals.TryGetValue(proposalId, out var proposal) ? proposal : null);

    public Task<ProposalEntity> UpdateInsuranceProposalAsync(ProposalEntity proposal, CancellationToken cancellationToken = default)
    {
        _proposals[proposal.ProposalId] = proposal;
        return Task.FromResult(proposal);
    }

    public Task<IReadOnlyList<ProposalEntity>> ListInsuranceProposalsAsync(
        string? orderId,
        string? insurerId,
        string? customerId,
        ProposalStatus? status,
        int page,
        int pageSize,
        CancellationToken cancellationToken = default)
    {
        var proposals = _proposals.Values
            .Where(x => string.IsNullOrWhiteSpace(orderId) || x.OrderId == orderId)
            .Where(x => string.IsNullOrWhiteSpace(insurerId) || x.InsurerId == insurerId)
            .Where(x => string.IsNullOrWhiteSpace(customerId) || x.CustomerId == customerId)
            .Where(x => !status.HasValue || status == ProposalStatus.Unspecified || x.Status == status)
            .Take(pageSize)
            .ToList();

        return Task.FromResult<IReadOnlyList<ProposalEntity>>(proposals);
    }

    public Task DeleteInsuranceProposalAsync(string proposalId, CancellationToken cancellationToken = default)
    {
        _proposals.Remove(proposalId);
        return Task.CompletedTask;
    }
}

public sealed class FakeRefundGateway : IRefundPaymentGateway
{
    public int Calls { get; private set; }

    public Task<InitiateRefundResponse> InitiateRefundAsync(
        string paymentId,
        ProtoMoney refundAmount,
        string reason,
        string initiatedBy,
        CancellationToken cancellationToken = default)
    {
        Calls++;
        return Task.FromResult(new InitiateRefundResponse
        {
            RefundId = "refund-1",
            Status = "REFUND_INITIATED",
            InitiatedAt = Timestamp.FromDateTime(DateTime.UtcNow)
        });
    }
}

public sealed class RecordingEventBus : IEventBus
{
    public List<PublishedEvent> Published { get; } = [];

    public Task PublishAsync<TEvent>(TEvent @event, string topic, CancellationToken cancellationToken = default)
        where TEvent : DomainEventBase
    {
        Published.Add(new PublishedEvent(topic, @event));
        return Task.CompletedTask;
    }

    public Task PublishAsync<TEvent>(TEvent @event, CancellationToken cancellationToken = default)
        where TEvent : DomainEventBase
    {
        Published.Add(new PublishedEvent(@event.EventType, @event));
        return Task.CompletedTask;
    }

    public Task PublishBatchAsync<TEvent>(IEnumerable<TEvent> events, string topic, CancellationToken cancellationToken = default)
        where TEvent : DomainEventBase
    {
        foreach (var @event in events)
        {
            Published.Add(new PublishedEvent(topic, @event));
        }

        return Task.CompletedTask;
    }

    public void Clear() => Published.Clear();
}

public sealed record PublishedEvent(string Topic, DomainEventBase Event);

public static class TestMoney
{
    public static ProtoMoney New(long amount, string currency = "BDT")
        => new() { Amount = amount, Currency = currency };
}
