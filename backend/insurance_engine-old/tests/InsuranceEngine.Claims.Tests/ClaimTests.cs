using System;
using FluentAssertions;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using Xunit;

namespace InsuranceEngine.Claims.Tests;

public class ClaimTests
{
    private Claim CreateTestClaim(long amount)
    {
        return Claim.File(
            "CLM-2024-TEST-001",
            Guid.NewGuid(),
            Guid.NewGuid(),
            ClaimType.Health,
            amount,
            DateTime.UtcNow.AddDays(-1),
            "Test incident",
            "Dhaka"
        );
    }

    [Theory]
    [InlineData(500_000, 0)]      // 5,000 BDT -> Level 0 (ZHTC)
    [InlineData(1_500_000, 1)]    // 15,000 BDT -> Level 1
    [InlineData(10_000_000, 2)]   // 100,000 BDT -> Level 2
    [InlineData(30_000_000, 3)]   // 300,000 BDT -> Level 3
    [InlineData(60_000_000, 4)]   // 600,000 BDT -> Level 4
    public void GetRequiredApprovalLevel_BasedOnAmount_ReturnsCorrectLevel(long amount, int expectedLevel)
    {
        // Arrange
        var claim = CreateTestClaim(amount);

        // Act
        var level = claim.GetRequiredApprovalLevel();

        // Assert
        level.Should().Be(expectedLevel);
    }

    [Fact]
    public void AddApproval_WhenLevelIsLowerThanRequired_SetsStatusToUnderReview()
    {
        // Arrange
        var claim = CreateTestClaim(10_000_000); // 100,000 BDT -> Requires Level 2
        var approverId = Guid.NewGuid();

        // Act
        var result = claim.AddApproval(approverId, "Officer", 1, ApprovalDecision.Approved, 5_000_000, "Approved by L1");

        // Assert
        result.IsSuccess.Should().BeTrue();
        claim.Status.Should().Be(ClaimStatus.UnderReview);
    }

    [Fact]
    public void AddApproval_WhenLevelMatchesRequired_SetsStatusToApproved()
    {
        // Arrange
        var claim = CreateTestClaim(10_000_000); // 100,000 BDT -> Requires Level 2
        var approverId = Guid.NewGuid();

        // Act
        var result = claim.AddApproval(approverId, "Manager", 2, ApprovalDecision.Approved, 5_000_000, "Approved by L2");

        // Assert
        result.IsSuccess.Should().BeTrue();
        claim.Status.Should().Be(ClaimStatus.Approved);
        claim.ApprovedAt.Should().NotBeNull();
    }

    [Fact]
    public void AddApproval_WhenDecisionIsRejected_SetsStatusToRejected()
    {
        // Arrange
        var claim = CreateTestClaim(500_000);
        var approverId = Guid.NewGuid();

        // Act
        var result = claim.AddApproval(approverId, "Officer", 1, ApprovalDecision.Rejected, 0, "Fake claim");

        // Assert
        result.IsSuccess.Should().BeTrue();
        claim.Status.Should().Be(ClaimStatus.Rejected);
        claim.RejectionReason.Should().Be("Fake claim");
    }
}
