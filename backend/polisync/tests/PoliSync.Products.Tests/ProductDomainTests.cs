using FluentAssertions;
using PoliSync.Products.Domain;
using PoliSync.SharedKernel.CQRS;
using Xunit;

namespace PoliSync.Products.Tests;

/// <summary>
/// Unit tests for Product domain model and PricingEngine.
/// Covers state machine transitions, validation rules, and pricing calculations.
/// </summary>
public class ProductDomainTests
{
    private readonly Guid _tenantId = Guid.NewGuid();
    private readonly Guid _partnerId = Guid.NewGuid();

    // ════════════════════════════════════════════════════════════════════════════
    // Product State Machine Tests
    // ════════════════════════════════════════════════════════════════════════════

    [Fact]
    public void Create_WithValidData_ReturnsSuccess()
    {
        // Arrange
        var tenantId = _tenantId;
        var partnerId = _partnerId;

        // Act
        var result = Product.Create(
            tenantId, partnerId, "HLTH-001", "Health Basic", "Basic health plan",
            "Health", 10000L, 500000L, 5000000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeNull();
        result.Value!.ProductCode.Should().Be("HLTH-001");
        result.Value.ProductName.Should().Be("Health Basic");
        result.Value.Status.Should().Be("Draft");
        result.Value.Version.Should().Be(1);
    }

    [Fact]
    public void Create_WithEmptyCode_ReturnsFail()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, string.Empty, "Health Basic", "Basic health plan",
            "Health", 10000L, 500000L, 5000000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_CODE");
    }

    [Fact]
    public void Create_WithNegativePremium_ReturnsFail()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, "HLTH-001", "Health Basic", "Basic health plan",
            "Health", -1000L, 500000L, 5000000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_PREMIUM");
    }

    [Fact]
    public void Create_WithInvalidSumInsuredRange_ReturnsFail()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, "HLTH-001", "Health Basic", "Basic health plan",
            "Health", 10000L, 5000000L, 500000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_SUM_INSURED");
    }

    [Fact]
    public void Create_WithTravelCategory_ValidTenure_ReturnsSuccess()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, "TRAVEL-001", "Travel Basic", "Basic travel plan",
            "Travel", 5000L, 100000L, 1000000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeTrue();
        result.Value!.Category.Should().Be("Travel");
        result.Value.MaxTenureMonths.Should().Be(12);
    }

    [Fact]
    public void Create_WithTravelCategory_TenureOver12_ReturnsFail()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, "TRAVEL-001", "Travel Basic", "Basic travel plan",
            "Travel", 5000L, 100000L, 1000000L, 1, 13, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_TENURE");
    }

    [Fact]
    public void Activate_FromDraft_ReturnsSuccess()
    {
        // Arrange
        var product = CreateValidProduct();

        // Act
        var result = product.Activate();

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be("Active");
    }

    [Fact]
    public void Activate_FromActive_ReturnsFail()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();

        // Act
        var result = product.Activate();

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_STATE");
    }

    [Fact]
    public void Activate_FromDiscontinued_ReturnsFail()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();
        product.Discontinue("End of life");

        // Act
        var result = product.Activate();

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_STATE");
    }

    [Fact]
    public void Deactivate_FromActive_ReturnsSuccess()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();

        // Act
        var result = product.Deactivate("Testing");

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be("Inactive");
    }

    [Fact]
    public void Deactivate_FromDraft_ReturnsFail()
    {
        // Arrange
        var product = CreateValidProduct();

        // Act
        var result = product.Deactivate("Testing");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_STATE");
    }

    [Fact]
    public void Discontinue_FromActive_ReturnsSuccess()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();

        // Act
        var result = product.Discontinue("End of life");

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be("Discontinued");
    }

    [Fact]
    public void Discontinue_FromInactive_ReturnsSuccess()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();
        product.Deactivate("Testing");

        // Act
        var result = product.Discontinue("End of life");

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be("Discontinued");
    }

    [Fact]
    public void Discontinue_FromDiscontinued_ReturnsFail()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();
        product.Discontinue("End of life");

        // Act
        var result = product.Discontinue("Another reason");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_STATE");
    }

    [Fact]
    public void Update_InDraftStatus_IncrementsVersion()
    {
        // Arrange
        var product = CreateValidProduct();
        var initialVersion = product.Version;

        // Act
        var result = product.Update("Updated Name", "Updated Description", 15000L, 500000L, 5000000L, new List<string>(), "admin");

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Version.Should().Be(initialVersion + 1);
        product.ProductName.Should().Be("Updated Name");
    }

    [Fact]
    public void Update_InActiveStatus_ReturnsFail()
    {
        // Arrange
        var product = CreateValidProduct();
        product.Activate();

        // Act
        var result = product.Update("Updated Name", "Updated Description", 15000L, 500000L, 5000000L, new List<string>(), "admin");

        // Assert
        result.IsSuccess.Should().BeFalse();
        result.Error!.Code.Should().Be("INVALID_STATE");
    }

    [Fact]
    public void ValidateSumInsured_WithinRange_ReturnsTrue()
    {
        // Arrange
        var product = CreateValidProduct();

        // Act
        var isValid = product.ValidateSumInsured(2000000L);

        // Assert
        isValid.Should().BeTrue();
    }

    [Fact]
    public void ValidateSumInsured_BelowMin_ReturnsFalse()
    {
        // Arrange
        var product = CreateValidProduct();

        // Act
        var isValid = product.ValidateSumInsured(100000L);

        // Assert
        isValid.Should().BeFalse();
    }

    [Fact]
    public void ValidateSumInsured_AboveMax_ReturnsFalse()
    {
        // Arrange
        var product = CreateValidProduct();

        // Act
        var isValid = product.ValidateSumInsured(10000000L);

        // Assert
        isValid.Should().BeFalse();
    }

    [Fact]
    public void Create_RaisesProductCreatedDomainEvent()
    {
        // Act
        var result = Product.Create(
            _tenantId, _partnerId, "HLTH-001", "Health Basic", "Basic health plan",
            "Health", 10000L, 500000L, 5000000L, 1, 12, new List<string>(), "BDT", "admin");

        // Assert
        result.IsSuccess.Should().BeTrue();
        var product = result.Value!;
        var events = product.DomainEvents;
        events.Should().NotBeEmpty();
        events.Should().Contain(e => e.GetType().Name == "ProductCreatedDomainEvent");
    }

    [Fact]
    public void Activate_RaisesProductActivatedDomainEvent()
    {
        // Arrange
        var product = CreateValidProduct();
        product.DomainEvents.Clear();

        // Act
        product.Activate();

        // Assert
        var events = product.DomainEvents;
        events.Should().NotBeEmpty();
        events.Should().Contain(e => e.GetType().Name == "ProductActivatedDomainEvent");
    }

    // ════════════════════════════════════════════════════════════════════════════
    // PricingEngine Tests
    // ════════════════════════════════════════════════════════════════════════════

    [Fact]
    public void Evaluate_WithMultiplyRule_AppliesCorrectly()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("multiply_age", "MULTIPLY", 1.1M, new List<RuleCondition> { new("age", "GTE", "30"), new("age", "LTE", "40") })
        };

        var inputFactors = new Dictionary<string, string> { { "age", "35" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(11000L); // 10000 * 1.1
    }

    [Fact]
    public void Evaluate_WithAddRule_AddsFixedAmount()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("add_surcharge", "ADD", 1000M, new List<RuleCondition> { new("rider_count", "GT", "0") })
        };

        var inputFactors = new Dictionary<string, string> { { "rider_count", "1" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(11000L);
    }

    [Fact]
    public void Evaluate_WithDiscountRule_SubtractsPercentage()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("loyalty_discount", "DISCOUNT", 10M, new List<RuleCondition> { new("years_customer", "GT", "5") })
        };

        var inputFactors = new Dictionary<string, string> { { "years_customer", "7" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(9000L); // 10000 - (10000 * 0.1)
    }

    [Fact]
    public void Evaluate_WithSetRule_SetsExactAmount()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("fixed_premium", "SET", 15000M, new List<RuleCondition> { new("product_type", "EQ", "special") })
        };

        var inputFactors = new Dictionary<string, string> { { "product_type", "special" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(15000L);
    }

    [Fact]
    public void Evaluate_ConditionEQ_Matches()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("gender_female", "MULTIPLY", 0.95M, new List<RuleCondition> { new("gender", "EQ", "F") })
        };

        var inputFactors = new Dictionary<string, string> { { "gender", "F" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(9500L); // 10000 * 0.95
    }

    [Fact]
    public void Evaluate_ConditionGT_FiltersCorrectly()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("high_age", "MULTIPLY", 1.5M, new List<RuleCondition> { new("age", "GT", "50") })
        };

        var inputFactors = new Dictionary<string, string> { { "age", "55" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(15000L); // 10000 * 1.5
    }

    [Fact]
    public void Evaluate_ConditionIN_MatchesList()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("city_surcharge", "ADD", 2000M, new List<RuleCondition> { new("city", "IN", "Dhaka,Chattogram") })
        };

        var inputFactors = new Dictionary<string, string> { { "city", "Dhaka" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(12000L);
    }

    [Fact]
    public void Evaluate_ConditionBETWEEN_MatchesRange()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("age_bracket", "MULTIPLY", 1.2M, new List<RuleCondition> { new("age", "BETWEEN", "30,50") })
        };

        var inputFactors = new Dictionary<string, string> { { "age", "40" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(12000L); // 10000 * 1.2
    }

    [Fact]
    public void Evaluate_MultipleRules_ApplyAllFalse_StopsAtFirstMatch()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("rule1", "MULTIPLY", 1.5M, new List<RuleCondition> { new("age", "GT", "50") }, applyAll: false),
            CreatePricingRule("rule2", "MULTIPLY", 1.2M, new List<RuleCondition> { new("age", "GT", "40") }, applyAll: false)
        };

        var inputFactors = new Dictionary<string, string> { { "age", "55" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(15000L); // Only first rule applied
        result.AppliedRules.Should().HaveCount(1);
    }

    [Fact]
    public void Evaluate_MultipleRules_ApplyAllTrue_AppliesAll()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("rule1", "MULTIPLY", 1.2M, new List<RuleCondition> { new("age", "GT", "40") }, applyAll: true),
            CreatePricingRule("rule2", "ADD", 1000M, new List<RuleCondition> { new("has_rider", "EQ", "true") }, applyAll: true)
        };

        var inputFactors = new Dictionary<string, string> { { "age", "55" }, { "has_rider", "true" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.AppliedRules.Should().HaveCount(2);
    }

    [Fact]
    public void Evaluate_WithNoMatchingConditions_ReturnsBasePremium()
    {
        // Arrange
        var basePremium = 10000L;
        var rules = new List<PricingRule>
        {
            CreatePricingRule("age_rule", "MULTIPLY", 1.5M, new List<RuleCondition> { new("age", "GT", "60") })
        };

        var inputFactors = new Dictionary<string, string> { { "age", "35" } };

        // Act
        var result = PricingEngine.Evaluate(rules, basePremium, inputFactors);

        // Assert
        result.FinalPremiumPaisa.Should().Be(basePremium);
    }

    // ════════════════════════════════════════════════════════════════════════════
    // Helper Methods
    // ════════════════════════════════════════════════════════════════════════════

    /// <summary>
    /// Creates a valid product for testing.
    /// </summary>
    private static Product CreateValidProduct() =>
        Product.Create(
            Guid.NewGuid(), Guid.NewGuid(), "HLTH-001", "Health Basic", "Basic health plan",
            "Health", 10000L, 500000L, 5000000L, 1, 12, new List<string>(), "BDT", "admin"
        ).Value!;

    /// <summary>
    /// Creates a pricing rule for testing.
    /// </summary>
    private static PricingRule CreatePricingRule(
        string name,
        string actionType,
        decimal value,
        List<RuleCondition> conditions,
        bool applyAll = false) =>
        new(
            RuleId: Guid.NewGuid(),
            RuleName: name,
            RuleType: "PRICING",
            Conditions: conditions,
            Action: new RuleAction(actionType, value),
            ApplyAll: applyAll,
            Priority: 1
        );
}
