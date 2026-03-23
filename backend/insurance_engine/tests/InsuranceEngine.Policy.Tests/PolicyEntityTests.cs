using System.Reflection;
using FluentAssertions;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using Xunit;

namespace InsuranceEngine.Policy.Tests;

public class PolicyAggregateTests
{
    private void SetStatus(PolicyAggregate policy, PolicyStatus status)
    {
        var property = typeof(PolicyAggregate).GetProperty("Status", BindingFlags.Public | BindingFlags.Instance);
        property?.SetValue(policy, status);
    }

    private void SetProperty(PolicyAggregate policy, string propertyName, object value)
    {
        var property = typeof(PolicyAggregate).GetProperty(propertyName, BindingFlags.Public | BindingFlags.Instance | BindingFlags.NonPublic);
        property?.SetValue(policy, value);
    }

    [Fact]
    public void Issue_WhenStatusIsPendingPayment_ShoudSucceed()
    {
        // Arrange
        var policy = new PolicyAggregate();
        SetStatus(policy, PolicyStatus.PendingPayment);
        policy.SetUnderwritingDecision(Guid.NewGuid());
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
        var policy = new PolicyAggregate();
        SetStatus(policy, PolicyStatus.Draft);

        // Act
        var result = policy.Issue(DateTime.UtcNow);

        // Assert
        result.IsSuccess.Should().BeFalse();
        policy.Status.Should().Be(PolicyStatus.Draft);
    }

    [Fact]
    public void Issue_WhenUnderwritingDecisionIsMissing_ShouldFail()
    {
        // Arrange
        var policy = new PolicyAggregate();
        SetStatus(policy, PolicyStatus.PendingPayment);
        // No SetUnderwritingDecision call

        // Act
        var result = policy.Issue(DateTime.UtcNow);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("Underwriting");
    }

    [Fact]
    public void AddNominee_WhenTotalShareIs100_ShouldSucceed()
    {
        // Arrange
        var policy = new PolicyAggregate();

        // Act
        var result = policy.AddNominee(Guid.NewGuid(), "Test Spouse", "Spouse", 100);

        // Assert
        result.IsSuccess.Should().BeTrue();
        policy.Nominees.Should().HaveCount(1);
    }

    [Fact]
    public void AddNominee_WhenTotalShareIsNot100_ShouldFail()
    {
        // Arrange
        var policy = new PolicyAggregate();

        // Act
        var result = policy.AddNominee(Guid.NewGuid(), "Test Spouse", "Spouse", 50);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("sum to 100");
    }

    [Fact]
    public void RemoveNominee_WhenRemainingShareIsNot100_ShouldFail()
    {
        // Arrange
        var policy = new PolicyAggregate();
        var nominee1Id = Guid.NewGuid();
        var nominee2Id = Guid.NewGuid();
        
        // Use reflection to add nominees directly to avoid share percentage validation in AddNominee
        var nomineesField = typeof(PolicyAggregate).GetField("<Nominees>k__BackingField", BindingFlags.Instance | BindingFlags.NonPublic);
        var nominees = (List<Nominee>)nomineesField.GetValue(policy);
        nominees.Add(new Nominee(nominee1Id) { SharePercentage = 50 });
        nominees.Add(new Nominee(nominee2Id) { SharePercentage = 50 });

        // Act
        var result = policy.RemoveNominee(nominee1Id);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("sum to 100");
    }

    [Fact]
    public void UpdateNominee_WhenTotalShareIsNot100_ShouldFail()
    {
        // Arrange
        var policy = new PolicyAggregate();
        var nomineeId = Guid.NewGuid();
        
        // Add valid 100% nominee first (via reflection to bypass initial validation)
        var nomineesField = typeof(PolicyAggregate).GetField("<Nominees>k__BackingField", BindingFlags.Instance | BindingFlags.NonPublic);
        var nominees = (List<Nominee>)nomineesField?.GetValue(policy);
        nominees?.Add(new Nominee(nomineeId) { SharePercentage = 100 });

        // Act - Try updating to 50%
        var result = policy.UpdateNominee(nomineeId, null, null, 50);

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error.Message.Should().Contain("sum to 100");
    }

    [Fact]
    public void Renew_WithoutClaims_ShouldApplyNCB()
    {
        // Arrange
        var oldPolicy = new PolicyAggregate();
        SetProperty(oldPolicy, "PremiumAmount", 10000L);
        SetProperty(oldPolicy, "ClaimsHistorySummary", null);
        
        // Act
        var newPolicy = PolicyAggregate.Renew(oldPolicy, "NEW-123", 12);

        // Assert
        newPolicy.PremiumAmount.Should().Be((long)(10000 * 0.90)); // 10% discount
    }

    [Fact]
    public void Renew_WithClaims_ShouldApplyPenalty()
    {
        // Arrange
        var oldPolicy = new PolicyAggregate();
        SetProperty(oldPolicy, "PremiumAmount", 10000L);
        SetProperty(oldPolicy, "ClaimsHistorySummary", "1 claim");
        
        // Act
        var newPolicy = PolicyAggregate.Renew(oldPolicy, "NEW-123", 12);

        // Assert
        newPolicy.PremiumAmount.Should().Be((long)(10000 * 1.15)); // 15% penalty
    }

    [Fact]
    public void Renew_InGracePeriod_WithoutClaims_ShouldApplyBoth()
    {
        // Arrange
        var oldPolicy = new PolicyAggregate();
        SetProperty(oldPolicy, "PremiumAmount", 10000L);
        SetProperty(oldPolicy, "ClaimsHistorySummary", null);
        SetStatus(oldPolicy, PolicyStatus.GracePeriod);
        
        // Act
        var newPolicy = PolicyAggregate.Renew(oldPolicy, "NEW-123", 12);

        // Assert
        // First NCB: 10000 * 0.9 = 9000
        // Then Grace: 9000 * 1.05 = 9450
        newPolicy.PremiumAmount.Should().Be(9450);
    }
}
