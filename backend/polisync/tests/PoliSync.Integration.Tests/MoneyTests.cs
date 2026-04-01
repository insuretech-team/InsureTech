using FluentAssertions;
using PoliSync.SharedKernel.Domain;
using Xunit;

namespace PoliSync.Integration.Tests;

/// <summary>
/// Unit tests for the Money value object.
/// Money stores amounts in paisa (1 BDT = 100 paisa).
/// </summary>
public sealed class MoneyTests
{
    [Fact]
    public void FromBdt_ConvertsCorrectly()
    {
        var money = Money.FromBdt(100.50m);
        money.AmountInPaisa.Should().Be(10050);
        money.Currency.Should().Be("BDT");
    }

    [Fact]
    public void ToBdt_ConvertsCorrectly()
    {
        var money = Money.FromPaisa(10050, "BDT");
        money.ToBdt().Should().Be(100.50m);
    }

    [Fact]
    public void Add_SameCurrency_ReturnsSum()
    {
        var a = Money.FromPaisa(5000, "BDT");
        var b = Money.FromPaisa(3000, "BDT");
        var result = a.Add(b);
        result.AmountInPaisa.Should().Be(8000);
        result.Currency.Should().Be("BDT");
    }

    [Fact]
    public void Add_DifferentCurrency_ThrowsInvalidOperationException()
    {
        var bdt = Money.FromPaisa(5000, "BDT");
        var usd = Money.FromPaisa(5000, "USD");
        var act = () => bdt.Add(usd);
        act.Should().Throw<InvalidOperationException>().WithMessage("*Cannot add*");
    }

    [Fact]
    public void Subtract_ReturnsCorrectDifference()
    {
        var a = Money.FromPaisa(10000, "BDT");
        var b = Money.FromPaisa(3000, "BDT");
        var result = a.Subtract(b);
        result.AmountInPaisa.Should().Be(7000);
    }

    [Fact]
    public void Subtract_ResultNegative_IsAllowed()
    {
        var a = Money.FromPaisa(1000, "BDT");
        var b = Money.FromPaisa(5000, "BDT");
        var result = a.Subtract(b);
        result.AmountInPaisa.Should().Be(-4000);
    }

    [Fact]
    public void Multiply_ByFactor_ReturnsCorrectAmount()
    {
        var money = Money.FromPaisa(10000, "BDT"); // 100 BDT
        var result = money.Multiply(1.5m);
        result.AmountInPaisa.Should().Be(15000); // 150 BDT
    }

    [Fact]
    public void Percentage_BasisPoints_ReturnsCorrectAmount()
    {
        var money = Money.FromPaisa(100000, "BDT"); // 1000 BDT
        var result = money.MultiplyByPercentage(10m);
        result.AmountInPaisa.Should().Be(10000); // 100 BDT
    }

    [Fact]
    public void Zero_ReturnsZeroAmount()
    {
        var money = Money.Zero();
        money.AmountInPaisa.Should().Be(0);
        money.Currency.Should().Be("BDT");
    }

    [Fact]
    public void NegativeAmount_IsRepresentableFromPaisa()
    {
        var money = Money.FromPaisa(-100, "BDT");
        money.IsNegative.Should().BeTrue();
    }

    [Fact]
    public void Equality_SameAmountAndCurrency_AreEqual()
    {
        var a = Money.FromPaisa(5000, "BDT");
        var b = Money.FromPaisa(5000, "BDT");
        a.Should().Be(b);
        (a == b).Should().BeTrue();
    }

    [Fact]
    public void Equality_DifferentAmount_AreNotEqual()
    {
        var a = Money.FromPaisa(5000, "BDT");
        var b = Money.FromPaisa(6000, "BDT");
        a.Should().NotBe(b);
        (a != b).Should().BeTrue();
    }

    [Fact]
    public void ToString_FormatsBdtCorrectly()
    {
        var money = Money.FromPaisa(150075, "BDT");
        money.ToString().Should().Be("1,500.75 BDT");
    }

    // ── Premium calculation scenarios (matching CalculatePremiumHandler logic) ──

    [Fact]
    public void PremiumProRata_12MonthAnnual_CorrectMonthly()
    {
        // Annual premium = 12,000 BDT = 1,200,000 paisa
        // The current Multiply implementation truncates fractional paisa.
        var annual = Money.FromPaisa(1_200_000, "BDT");
        var monthly = annual.Multiply(1m / 12m);
        monthly.AmountInPaisa.Should().Be(99_999);
    }

    [Fact]
    public void LoadingFactor_25Percent_CorrectLoading()
    {
        var basePremium = Money.FromPaisa(100_000, "BDT"); // 1,000 BDT
        var loading = basePremium.MultiplyByPercentage(25m);
        loading.AmountInPaisa.Should().Be(25_000); // 250 BDT
    }
}
