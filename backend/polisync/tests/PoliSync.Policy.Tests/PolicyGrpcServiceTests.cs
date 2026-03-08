using FluentAssertions;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Policy.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Policy.GrpcServices;
using PoliSync.Policy.Infrastructure;
using PolicyEntity = Insuretech.Policy.Entity.V1.Policy;
using Xunit;

namespace PoliSync.Policy.Tests;

public class PolicyGrpcServiceTests
{
    private sealed class FakePolicyDataGateway : IPolicyDataGateway
    {
        private readonly Dictionary<string, PolicyEntity> _store = new();

        public Task<PolicyEntity> CreatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
        {
            _store[policy.PolicyId] = policy;
            return Task.FromResult(policy);
        }

        public Task<PolicyEntity?> GetPolicyAsync(string policyId, CancellationToken cancellationToken = default)
            => Task.FromResult(_store.TryGetValue(policyId, out var policy) ? policy : null);

        public Task<PolicyEntity> UpdatePolicyAsync(PolicyEntity policy, CancellationToken cancellationToken = default)
        {
            _store[policy.PolicyId] = policy;
            return Task.FromResult(policy);
        }

        public Task<IReadOnlyList<PolicyEntity>> ListPoliciesAsync(string customerId, int page, int pageSize, CancellationToken cancellationToken = default)
        {
            var items = _store.Values
                .Where(x => string.IsNullOrWhiteSpace(customerId) || x.CustomerId == customerId)
                .Skip((page - 1) * pageSize)
                .Take(pageSize)
                .ToList();
            return Task.FromResult<IReadOnlyList<PolicyEntity>>(items);
        }

        public Task DeletePolicyAsync(string policyId, CancellationToken cancellationToken = default)
        {
            _store.Remove(policyId);
            return Task.CompletedTask;
        }
    }

    private static PolicyGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<PolicyGrpcService>.Instance, new FakePolicyDataGateway());

    [Fact]
    public async Task CreateIssueAndGetPolicy_CompletesLifecycle()
    {
        var service = CreateService();

        var created = await service.CreatePolicy(new CreatePolicyRequest
        {
            ProductId = $"prod-{Guid.NewGuid():N}",
            CustomerId = $"cust-{Guid.NewGuid():N}",
            PartnerId = $"partner-{Guid.NewGuid():N}",
            AgentId = $"agent-{Guid.NewGuid():N}",
            PremiumAmount = new Money { Amount = 120_000, Currency = "BDT" },
            SumInsured = new Money { Amount = 1_000_000, Currency = "BDT" },
            TenureMonths = 12,
            Applicant = new Applicant
            {
                FullName = "Test User",
                Address = "Dhaka",
                DateOfBirth = Timestamp.FromDateTime(DateTime.UtcNow.AddYears(-30))
            }
        }, null!);

        var issued = await service.IssuePolicy(new IssuePolicyRequest
        {
            PolicyId = created.PolicyId,
            QuoteId = $"quote-{Guid.NewGuid():N}",
            PaymentId = $"pay-{Guid.NewGuid():N}"
        }, null!);

        var get = await service.GetPolicy(new GetPolicyRequest { PolicyId = created.PolicyId }, null!);

        created.PolicyId.Should().NotBeNullOrWhiteSpace();
        issued.Policy.Status.Should().Be(PolicyStatus.Active);
        get.Policy.ReceiptNumber.Should().StartWith("RCPT-");
    }

    [Fact]
    public async Task UpdateAndListPolicy_ReturnsUpdatedNominees()
    {
        var service = CreateService();
        var customerId = $"cust-{Guid.NewGuid():N}";

        var created = await service.CreatePolicy(new CreatePolicyRequest
        {
            ProductId = $"prod-{Guid.NewGuid():N}",
            CustomerId = customerId,
            PremiumAmount = new Money { Amount = 100_000, Currency = "BDT" },
            SumInsured = new Money { Amount = 700_000, Currency = "BDT" },
            TenureMonths = 6
        }, null!);

        await service.UpdatePolicy(new UpdatePolicyRequest
        {
            PolicyId = created.PolicyId,
            Address = "New Address",
            Nominees =
            {
                new Nominee
                {
                    NomineeId = Guid.NewGuid().ToString("N"),
                    FullName = "Nominee One",
                    Relationship = "Spouse",
                    SharePercentage = 100
                }
            }
        }, null!);

        var list = await service.ListUserPolicies(new ListUserPoliciesRequest
        {
            CustomerId = customerId,
            Page = 1,
            PageSize = 10
        }, null!);

        list.Policies.Should().ContainSingle(x => x.PolicyId == created.PolicyId);
        list.Policies.Single(x => x.PolicyId == created.PolicyId).Nominees.Should().ContainSingle();
    }

    [Fact]
    public async Task CancelPolicy_ReturnsRefundAndCancelledStatus()
    {
        var service = CreateService();

        var created = await service.CreatePolicy(new CreatePolicyRequest
        {
            ProductId = $"prod-{Guid.NewGuid():N}",
            CustomerId = $"cust-{Guid.NewGuid():N}",
            PremiumAmount = new Money { Amount = 150_000, Currency = "BDT" },
            SumInsured = new Money { Amount = 800_000, Currency = "BDT" },
            TenureMonths = 12
        }, null!);

        var cancelled = await service.CancelPolicy(new CancelPolicyRequest
        {
            PolicyId = created.PolicyId,
            Reason = "Customer request"
        }, null!);

        var get = await service.GetPolicy(new GetPolicyRequest { PolicyId = created.PolicyId }, null!);

        cancelled.RefundAmount.Amount.Should().Be(90_000);
        get.Policy.Status.Should().Be(PolicyStatus.Cancelled);
    }
}
