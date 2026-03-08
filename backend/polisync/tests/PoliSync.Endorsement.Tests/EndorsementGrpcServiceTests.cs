using FluentAssertions;
using Insuretech.Endorsement.Entity.V1;
using Insuretech.Endorsement.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Endorsement.GrpcServices;
using PoliSync.Endorsement.Infrastructure;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;
using Xunit;

namespace PoliSync.Endorsement.Tests;

public class EndorsementGrpcServiceTests
{
    private sealed class FakeEndorsementDataGateway : IEndorsementDataGateway
    {
        private readonly Dictionary<string, EndorsementEntity> _store = new();

        public Task<EndorsementEntity> CreateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default)
        {
            _store[endorsement.Id] = endorsement;
            return Task.FromResult(endorsement);
        }

        public Task<EndorsementEntity?> GetEndorsementAsync(string endorsementId, CancellationToken cancellationToken = default)
            => Task.FromResult(_store.TryGetValue(endorsementId, out var endorsement) ? endorsement : null);

        public Task<EndorsementEntity> UpdateEndorsementAsync(EndorsementEntity endorsement, CancellationToken cancellationToken = default)
        {
            _store[endorsement.Id] = endorsement;
            return Task.FromResult(endorsement);
        }

        public Task<IReadOnlyList<EndorsementEntity>> ListEndorsementsByPolicyAsync(string policyId, CancellationToken cancellationToken = default)
        {
            var items = _store.Values
                .Where(x => x.PolicyId == policyId)
                .ToList();
            return Task.FromResult<IReadOnlyList<EndorsementEntity>>(items);
        }
    }

    private static EndorsementGrpcService CreateService()
        => new(Mock.Of<IMediator>(), NullLogger<EndorsementGrpcService>.Instance, new FakeEndorsementDataGateway());

    [Fact]
    public async Task RequestAndGetEndorsement_ReturnsPendingEndorsement()
    {
        var service = CreateService();

        var requested = await service.RequestEndorsement(new RequestEndorsementRequest
        {
            PolicyId = $"pol-{Guid.NewGuid():N}",
            Type = "ENDORSEMENT_TYPE_CONTACT_CHANGE",
            Reason = "Phone updated",
            Changes = "phone:+8801xxxx",
            EffectiveDate = DateTime.UtcNow.AddDays(1).ToString("yyyy-MM-dd")
        }, null!);

        var get = await service.GetEndorsement(new GetEndorsementRequest
        {
            EndorsementId = requested.EndorsementId
        }, null!);

        requested.EndorsementId.Should().NotBeNullOrWhiteSpace();
        get.Endorsement.Status.Should().Be(EndorsementStatus.Pending);
        get.Endorsement.Type.Should().Be(EndorsementType.ContactChange);
    }

    [Fact]
    public async Task ApproveEndorsement_MarksAsAppliedAndListable()
    {
        var service = CreateService();
        var policyId = $"pol-{Guid.NewGuid():N}";

        var requested = await service.RequestEndorsement(new RequestEndorsementRequest
        {
            PolicyId = policyId,
            Type = "ENDORSEMENT_TYPE_SUM_ASSURED_CHANGE",
            Reason = "Coverage increase",
            Changes = "sum_assured:+200000",
            EffectiveDate = DateTime.UtcNow.AddDays(2).ToString("yyyy-MM-dd")
        }, null!);

        await service.ApproveEndorsement(new ApproveEndorsementRequest
        {
            EndorsementId = requested.EndorsementId,
            ApprovedBy = "ops-user",
            Comments = "Approved"
        }, null!);

        var list = await service.ListEndorsements(new ListEndorsementsRequest
        {
            PolicyId = policyId,
            Status = "ENDORSEMENT_STATUS_APPLIED",
            Page = 1,
            PageSize = 10
        }, null!);

        list.Endorsements.Should().ContainSingle(x => x.Id == requested.EndorsementId);
        list.Endorsements.Single().Status.Should().Be(EndorsementStatus.Applied);
    }
}
