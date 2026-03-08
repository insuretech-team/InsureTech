using FluentAssertions;
using Insuretech.Claims.Entity.V1;
using Insuretech.Claims.Services.V1;
using Insuretech.Common.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Claims.GrpcServices;
using PoliSync.Claims.Infrastructure;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;
using Xunit;

namespace PoliSync.Claims.Tests;

public class ClaimGrpcServiceTests
{
    private sealed class FakeClaimDataGateway : IClaimDataGateway
    {
        private readonly Dictionary<string, ClaimEntity> _store = new();

        public Task<ClaimEntity> CreateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default)
        {
            _store[claim.ClaimId] = claim;
            return Task.FromResult(claim);
        }

        public Task<ClaimEntity?> GetClaimAsync(string claimId, CancellationToken cancellationToken = default)
            => Task.FromResult(_store.TryGetValue(claimId, out var claim) ? claim : null);

        public Task<ClaimEntity> UpdateClaimAsync(ClaimEntity claim, CancellationToken cancellationToken = default)
        {
            _store[claim.ClaimId] = claim;
            return Task.FromResult(claim);
        }

        public Task<IReadOnlyList<ClaimEntity>> ListClaimsAsync(string customerId, string policyId, int page, int pageSize, CancellationToken cancellationToken = default)
        {
            var items = _store.Values
                .Where(x => string.IsNullOrWhiteSpace(customerId) || x.CustomerId == customerId)
                .Where(x => string.IsNullOrWhiteSpace(policyId) || x.PolicyId == policyId)
                .Skip((page - 1) * pageSize)
                .Take(pageSize)
                .ToList();
            return Task.FromResult<IReadOnlyList<ClaimEntity>>(items);
        }
    }

    private static ClaimGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<ClaimGrpcService>.Instance, new FakeClaimDataGateway());

    [Fact]
    public async Task SubmitClaim_ThenGetClaim_ReturnsStoredClaim()
    {
        var service = CreateService();
        var customerId = $"cust-{Guid.NewGuid():N}";

        var submit = await service.SubmitClaim(new SubmitClaimRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            CustomerId = customerId,
            Type = ClaimType.HealthHospitalization,
            ClaimedAmount = new Money { Amount = 150_000, Currency = "BDT" },
            IncidentDate = DateTime.UtcNow.ToString("yyyy-MM-dd"),
            IncidentDescription = "Hospital admission"
        }, null!);

        var get = await service.GetClaim(new GetClaimRequest { ClaimId = submit.ClaimId }, null!);

        submit.ClaimId.Should().NotBeNullOrWhiteSpace();
        get.Claim.ClaimId.Should().Be(submit.ClaimId);
        get.Claim.CustomerId.Should().Be(customerId);
        get.Claim.Status.Should().Be(ClaimStatus.Submitted);
    }

    [Fact]
    public async Task ApproveAndSettleClaim_UpdatesClaimStatusAndAmounts()
    {
        var service = CreateService();

        var submit = await service.SubmitClaim(new SubmitClaimRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            CustomerId = $"cust-{Guid.NewGuid():N}",
            Type = ClaimType.MotorAccident,
            ClaimedAmount = new Money { Amount = 200_000, Currency = "BDT" },
            IncidentDate = DateTime.UtcNow.ToString("yyyy-MM-dd"),
            IncidentDescription = "Road accident"
        }, null!);

        await service.ApproveClaim(new ApproveClaimRequest
        {
            ClaimId = submit.ClaimId,
            ApproverId = "approver-1",
            ApprovedAmount = new Money { Amount = 175_000, Currency = "BDT" },
            Notes = "Within policy terms"
        }, null!);

        var settle = await service.SettleClaim(new SettleClaimRequest
        {
            ClaimId = submit.ClaimId,
            PaymentMethod = "bank_transfer",
            PaymentReference = "ref-001"
        }, null!);

        var get = await service.GetClaim(new GetClaimRequest { ClaimId = submit.ClaimId }, null!);

        settle.PaymentId.Should().NotBeNullOrWhiteSpace();
        settle.SettledAmount.Amount.Should().Be(175_000);
        get.Claim.Status.Should().Be(ClaimStatus.Settled);
        get.Claim.SettledAmount.Amount.Should().Be(175_000);
    }
}
