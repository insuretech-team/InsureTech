using FluentAssertions;
using InsuranceEngine.Claims.Domain.Entities;
using InsuranceEngine.Claims.Domain.Enums;
using Xunit;

namespace InsuranceEngine.Claims.Tests;

public class ClaimTests
{
    [Theory]
    [InlineData(500_000, 1)]      // 5,000 BDT -> Level 1
    [InlineData(1_500_000, 2)]    // 15,000 BDT -> Level 2
    [InlineData(10_000_000, 3)]   // 100,000 BDT -> Level 3
    [InlineData(30_000_000, 4)]   // 300,000 BDT -> Level 4
    public void GetRequiredApprovalLevel_BasedOnAmount_ReturnsCorrectLevel(long amount, int expectedLevel)
    {
        // Arrange
        var claim = new Claim { ClaimedAmount = amount };

        // Act
        var level = claim.GetRequiredApprovalLevel();

        // Assert
        level.Should().Be(expectedLevel);
    }

    [Fact]
    public void AddApproval_WhenLevelIsLowerThanRequired_SetsStatusToUnderReview()
    {
        // Arrange
        var claim = new Claim { ClaimedAmount = 5_000_000 }; // Requires Level 2
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
        var claim = new Claim { ClaimedAmount = 5_000_000 }; // Requires Level 2
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
        var claim = new Claim { ClaimedAmount = 500_000 };
        var approverId = Guid.NewGuid();

        // Act
        var result = claim.AddApproval(approverId, "Officer", 1, ApprovalDecision.Rejected, 0, "Fake claim");

        // Assert
        result.IsSuccess.Should().BeTrue();
        claim.Status.Should().Be(ClaimStatus.Rejected);
        claim.RejectionReason.Should().Be("Fake claim");
    }
}
