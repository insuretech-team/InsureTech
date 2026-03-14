using FluentAssertions;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using Xunit;

namespace InsuranceEngine.Policy.Tests;

public class PolicyEntityTests
{
    [Fact]
    public void Issue_WhenStatusIsPendingPayment_ShoudSucceed()
    {
        // Arrange
        var policy = new PolicyEntity { Status = PolicyStatus.PendingPayment };
        var issuedAt = DateTime.UtcNow;

        // Act
        var result = policy.Issue(issuedAt);

        // Assert
        result.IsSuccess.Should().BeTrue();
        policy.Status.Should().Be(PolicyStatus.Active);
        policy.IssuedAt.Should().Be(issuedAt);
    }

    [Fact]
    public void Issue_WhenStatusIsNotPendingPayment_ShouldFail()
    {
        // Arrange
        var policy = new PolicyEntity { Status = PolicyStatus.Draft };

        // Act
        var result = policy.Issue(DateTime.UtcNow);

        // Assert
        result.IsSuccess.Should().BeFalse();
        policy.Status.Should().Be(PolicyStatus.Draft);
    }

    [Fact]
    public void AddNominee_WhenTotalShareIs100_ShouldSucceed()
    {
        // Arrange
        var policy = new PolicyEntity { Id = Guid.NewGuid() };

        // Act
        var result = policy.AddNominee(Guid.NewGuid(), "Spouse", 100);

        // Assert
        result.IsSuccess.Should().BeTrue();
        policy.Nominees.Should().HaveCount(1);
    }

    [Fact]
    public void AddNominee_WhenTotalShareIsNot100_ShouldFail()
    {
        // Arrange
        var policy = new PolicyEntity { Id = Guid.NewGuid() };

        // Act
        var result = policy.AddNominee(Guid.NewGuid(), "Spouse", 50);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("sum to 100");
    }

    [Fact]
    public void RemoveNominee_WhenRemainingShareIsNot100_ShouldFail()
    {
        // Arrange
        var policy = new PolicyEntity { Id = Guid.NewGuid() };
        var nominee1Id = Guid.NewGuid();
        var nominee2Id = Guid.NewGuid();
        
        policy.Nominees.Add(new Nominee { Id = nominee1Id, SharePercentage = 50 });
        policy.Nominees.Add(new Nominee { Id = nominee2Id, SharePercentage = 50 });

        // Act
        var result = policy.RemoveNominee(nominee1Id);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("sum to 100");
    }
}
