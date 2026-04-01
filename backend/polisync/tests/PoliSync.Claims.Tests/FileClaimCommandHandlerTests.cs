using FluentAssertions;
using Insuretech.Claims.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.Claims.Application.Commands;
using PoliSync.Claims.Domain;
using PoliSync.Infrastructure.Persistence;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Workflow.Application.Commands;
using PoliSync.Workflow.Domain;
using Xunit;

namespace PoliSync.Claims.Tests;

/// <summary>
/// Unit tests for FileClaimCommandHandler.
///
/// Pattern (matches all other PoliSync test projects):
///   - Pure unit tests — no live DB, no network, no EF InMemory
///   - DB interaction (SaveClaimToDatabase) is tested via PoliSync.DbTest separately
///   - Moq for IEventBus + IMediator
///   - Domain logic tested directly on ClaimAggregate (no handler needed)
/// </summary>
public class FileClaimCommandHandlerTests
{
    // ── Test doubles ─────────────────────────────────────────────────────────

    private readonly Mock<IEventBus> _eventBusMock = new();
    private readonly Mock<IMediator> _mediatorMock = new();

    // Spy handler that overrides SaveClaimToDatabase to avoid live DB
    private sealed class TestableFileClaimCommandHandler : FileClaimCommandHandler
    {
        public string? LastSavedClaimId { get; private set; }

        public TestableFileClaimCommandHandler(
            IEventBus eventBus,
            IMediator mediator)
            : base(null!, eventBus, mediator, NullLogger<FileClaimCommandHandler>.Instance)
        {
        }

        protected override Task SaveClaimToDatabase(
            Insuretech.Claims.Entity.V1.Claim claim,
            CancellationToken cancellationToken)
        {
            LastSavedClaimId = claim.ClaimId;
            return Task.CompletedTask; // no-op — DB tested in PoliSync.DbTest
        }
    }

    private TestableFileClaimCommandHandler CreateHandler()
    {
        // Default: workflow trigger returns success
        _mediatorMock
            .Setup(m => m.Send(It.IsAny<TriggerWorkflowCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(SharedKernel.CQRS.Result.Ok(new TriggerWorkflowResult(
                WorkflowInstanceId: Guid.NewGuid().ToString(),
                TemplateName: "claim.standard-approval",
                EntityType: "CLAIM",
                EntityId: Guid.NewGuid().ToString(),
                WasTriggered: true)));

        return new TestableFileClaimCommandHandler(
            _eventBusMock.Object,
            _mediatorMock.Object);
    }

    private static FileClaimCommand StandardClaim(long amountPaisa = 5_000_000) => new(
        PolicyId: $"pol-{Guid.NewGuid():N}",
        CustomerId: $"cust-{Guid.NewGuid():N}",
        ClaimType: ClaimType.HealthHospitalization,
        ClaimedAmountPaisa: amountPaisa,
        IncidentDate: DateTime.UtcNow.AddDays(-3),
        IncidentDescription: "Hospital admission for appendicitis surgery",
        PlaceOfIncident: "Dhaka Medical College Hospital");

    // ── Handler tests ─────────────────────────────────────────────────────────

    [Fact]
    public async Task FileClaim_ReturnsClaimId()
    {
        var handler = CreateHandler();
        var result = await handler.Handle(StandardClaim(), CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeNullOrWhiteSpace();
        result.Value.Should().MatchRegex(@"^[0-9a-f\-]{36}$", "claim ID should be a UUID");
    }

    [Fact]
    public async Task FileClaim_SavesToDatabaseWithCorrectClaimId()
    {
        var handler = CreateHandler();
        var result = await handler.Handle(StandardClaim(), CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        handler.LastSavedClaimId.Should().Be(result.Value, "saved claim ID should match returned ID");
    }

    [Fact]
    public async Task FileClaim_PublishesDomainEvent()
    {
        var handler = CreateHandler();
        await handler.Handle(StandardClaim(), CancellationToken.None);

        // Handler iterates claimAggregate.DomainEvents and calls _eventBus.PublishAsync per event.
        // At minimum, the ClaimFiledEvent is published.
        _eventBusMock.Verify(
            e => e.PublishAsync(It.IsAny<SharedKernel.Domain.DomainEvent>(), It.IsAny<CancellationToken>()),
            Times.AtLeastOnce,
            "at least one domain event (ClaimFiledEvent) should be published");
    }

    [Fact]
    public async Task FileClaim_TriggersWorkflowWithCorrectContext()
    {
        var handler = CreateHandler();
        var command = StandardClaim(amountPaisa: 5_000_000);

        await handler.Handle(command, CancellationToken.None);

        _mediatorMock.Verify(
            m => m.Send(
                It.Is<TriggerWorkflowCommand>(c =>
                    c.Context.EntityType == "CLAIM" &&
                    c.Context.AmountPaisa == 5_000_000 &&
                    c.Context.Portal == "B2C" &&
                    c.Context.InitiatedBy == command.CustomerId &&
                    c.Context.Metadata.ContainsKey("policy_id") &&
                    c.Context.Metadata["policy_id"] == command.PolicyId),
                It.IsAny<CancellationToken>()),
            Times.Once);
    }

    [Fact]
    public async Task FileClaim_HighValue_TriggersWorkflowWithHighAmount()
    {
        var handler = CreateHandler();
        var result = await handler.Handle(StandardClaim(amountPaisa: 15_000_000), CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
        _mediatorMock.Verify(
            m => m.Send(
                It.Is<TriggerWorkflowCommand>(c => c.Context.AmountPaisa == 15_000_000),
                It.IsAny<CancellationToken>()),
            Times.Once);
    }

    [Fact]
    public async Task FileClaim_WorkflowTriggerFails_StillSucceeds()
    {
        // Workflow is a side-effect — claim filing must not fail if workflow unavailable
        _mediatorMock
            .Setup(m => m.Send(It.IsAny<TriggerWorkflowCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(SharedKernel.CQRS.Result.Fail<TriggerWorkflowResult>(
                "WORKFLOW_START_FAILED", "Go engine unavailable"));

        var handler = CreateHandler();
        var result = await handler.Handle(StandardClaim(), CancellationToken.None);

        result.IsSuccess.Should().BeTrue("claim should be saved even if workflow trigger fails");
    }

    [Fact]
    public async Task FileClaim_WorkflowNotTriggered_StillSucceeds()
    {
        _mediatorMock
            .Setup(m => m.Send(It.IsAny<TriggerWorkflowCommand>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync(SharedKernel.CQRS.Result.Ok(new TriggerWorkflowResult(
                string.Empty, string.Empty, "CLAIM", string.Empty, WasTriggered: false)));

        var handler = CreateHandler();
        var result = await handler.Handle(StandardClaim(), CancellationToken.None);

        result.IsSuccess.Should().BeTrue();
    }

    [Fact]
    public async Task FileClaim_GeneratesUniqueClaimIds()
    {
        var handler = CreateHandler();
        var r1 = await handler.Handle(StandardClaim(), CancellationToken.None);
        var r2 = await handler.Handle(StandardClaim(), CancellationToken.None);

        r1.IsSuccess.Should().BeTrue();
        r2.IsSuccess.Should().BeTrue();
        r1.Value.Should().NotBe(r2.Value, "each claim must have a unique ID");
    }

    // ── ClaimAggregate domain tests (no handler, no DB needed) ───────────────

    [Fact]
    public void ClaimAggregate_FileClaim_SetsCorrectInitialStatus()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.HealthSurgery,
            2_500_000, DateTime.UtcNow.AddDays(-5), "Surgery", "Dhaka");

        agg.ClaimId.Should().NotBeNullOrWhiteSpace();
        agg.ClaimNumber.Should().StartWith("CLM-");
        agg.Status.Should().Be(ClaimStatus.Submitted);
        agg.Claim.Type.Should().Be(ClaimType.HealthSurgery);
        agg.Claim.ClaimedAmount.Amount.Should().Be(2_500_000);
        agg.Claim.ClaimedAmount.Currency.Should().Be("BDT");
        agg.DomainEvents.Should().ContainSingle(e => e is ClaimFiledEvent);
    }

    [Theory]
    [InlineData(500_000L,    0)] // 5,000 BDT — ZHTC auto
    [InlineData(3_000_000L,  1)] // 30,000 BDT — L1 Officer
    [InlineData(10_000_000L, 2)] // 1,00,000 BDT — L2 Manager
    [InlineData(30_000_000L, 3)] // 3,00,000 BDT — L3 Director
    [InlineData(60_000_000L, 4)] // 6,00,000 BDT — Board
    public void ClaimAggregate_ApprovalMatrix_CorrectLevels(long amountPaisa, int expectedLevel)
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.Death,
            amountPaisa, DateTime.UtcNow, "incident", "place");

        agg.GetRequiredApprovalLevel().Should().Be(expectedLevel,
            $"amount {amountPaisa} paisa should require approval level {expectedLevel}");
    }

    [Fact]
    public void ClaimAggregate_AddApproval_Approved_UpdatesStatus()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.HealthHospitalization,
            3_000_000, DateTime.UtcNow, "Hospitalised", "Dhaka");

        agg.AddApproval("officer-001", "claims_officer", 1,
            ApprovalDecision.Approved, 2_800_000, "Verified, approving reduced amount");

        agg.Status.Should().Be(ClaimStatus.Approved);
        agg.Claim.ApprovedAmount.Amount.Should().Be(2_800_000);
        agg.DomainEvents.Should().Contain(e => e is PoliSync.Claims.Domain.ClaimApprovedEvent);
    }

    [Fact]
    public void ClaimAggregate_AddApproval_Rejected_FailsClaim()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.MotorTheft,
            5_000_000, DateTime.UtcNow, "Car stolen", "Mirpur");

        agg.AddApproval("manager-001", "claims_manager", 2,
            ApprovalDecision.Rejected, 0, "Policy exclusion applies");

        agg.Status.Should().Be(ClaimStatus.Rejected);
        agg.Claim.RejectionReason.Should().Contain("Policy exclusion applies");
        agg.DomainEvents.Should().Contain(e => e is PoliSync.Claims.Domain.ClaimRejectedEvent);
    }

    [Fact]
    public void ClaimAggregate_Settle_AfterApproval_UpdatesToSettled()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.HealthHospitalization,
            500_000, DateTime.UtcNow, "Admitted", "Dhaka");

        agg.AddApproval("officer-001", "claims_officer", 0,
            ApprovalDecision.Approved, 500_000, "ZHTC auto-approve");

        agg.Settle(500_000, "bank_transfer", "TXN-12345");

        agg.Status.Should().Be(ClaimStatus.Settled);
        agg.Claim.SettledAmount.Amount.Should().Be(500_000);
        agg.DomainEvents.Should().Contain(e => e is PoliSync.Claims.Domain.ClaimSettledEvent);
    }

    [Fact]
    public void ClaimAggregate_Settle_BeforeApproval_Throws()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.HealthHospitalization,
            500_000, DateTime.UtcNow, "Admitted", "Dhaka");

        var act = () => agg.Settle(500_000, "bank_transfer", "TXN-00001");

        act.Should().Throw<InvalidOperationException>()
            .WithMessage("*Cannot settle claim*");
    }

    [Fact]
    public void ClaimAggregate_FraudCheck_HighScore_SetsUnderReview()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.MotorAccident,
            10_000_000, DateTime.UtcNow, "Accident", "Highway");

        agg.ApplyFraudCheck(0.85, ["duplicate_claim", "suspicious_timing", "high_value"]);

        agg.Status.Should().Be(ClaimStatus.UnderReview);
        agg.Claim.FraudCheck.Should().NotBeNull();
        agg.Claim.FraudCheck!.RiskFactors.Should().HaveCount(3);
        agg.DomainEvents.Should().Contain(e => e is ClaimFlaggedForFraudEvent);
    }

    [Fact]
    public void ClaimAggregate_FraudCheck_LowScoreZhtc_SetsAutoAdjudicated()
    {
        var agg = ClaimAggregate.FileClaim(
            "pol-001", "cust-001", ClaimType.HealthHospitalization,
            800_000, DateTime.UtcNow, "Checkup", "Clinic");

        agg.ApplyFraudCheck(0.10, []);

        agg.Claim.ProcessingType.Should().Be(ClaimProcessingType.AutoAdjudicated);
    }
}
