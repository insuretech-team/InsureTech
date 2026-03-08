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
        money.AmountPaisa.Should().Be(10050);
        money.Currency.Should().Be("BDT");
    }

    [Fact]
    public void ToBdt_ConvertsCorrectly()
    {
        var money = new Money(10050, "BDT");
        money.ToBdt().Should().Be(100.50m);
    }

    [Fact]
    public void Add_SameCurrency_ReturnsSum()
    {
        var a = new Money(5000, "BDT");
        var b = new Money(3000, "BDT");
        var result = a.Add(b);
        result.AmountPaisa.Should().Be(8000);
        result.Currency.Should().Be("BDT");
    }

    [Fact]
    public void Add_DifferentCurrency_ThrowsInvalidOperationException()
    {
        var bdt = new Money(5000, "BDT");
        var usd = new Money(5000, "USD");
        var act = () => bdt.Add(usd);
        act.Should().Throw<InvalidOperationException>().WithMessage("*Currency mismatch*");
    }

    [Fact]
    public void Subtract_ReturnsCorrectDifference()
    {
        var a = new Money(10000, "BDT");
        var b = new Money(3000, "BDT");
        var result = a.Subtract(b);
        result.AmountPaisa.Should().Be(7000);
    }

    [Fact]
    public void Subtract_ResultNegative_ThrowsInvalidOperationException()
    {
        var a = new Money(1000, "BDT");
        var b = new Money(5000, "BDT");
        var act = () => a.Subtract(b);
        act.Should().Throw<InvalidOperationException>();
    }

    [Fact]
    public void Multiply_ByFactor_ReturnsCorrectAmount()
    {
        var money = new Money(10000, "BDT"); // 100 BDT
        var result = money.Multiply(1.5m);
        result.AmountPaisa.Should().Be(15000); // 150 BDT
    }

    [Fact]
    public void Percentage_BasisPoints_ReturnsCorrectAmount()
    {
        var money = new Money(100000, "BDT"); // 1000 BDT
        var result = money.Percentage(1000); // 10% = 1000 basis points
        result.AmountPaisa.Should().Be(10000); // 100 BDT
    }

    [Fact]
    public void Zero_ReturnsZeroAmount()
    {
        var money = Money.Zero();
        money.AmountPaisa.Should().Be(0);
        money.Currency.Should().Be("BDT");
    }

    [Fact]
    public void NegativeAmount_ThrowsArgumentException()
    {
        var act = () => new Money(-100, "BDT");
        act.Should().Throw<ArgumentException>();
    }

    [Fact]
    public void Equality_SameAmountAndCurrency_AreEqual()
    {
        var a = new Money(5000, "BDT");
        var b = new Money(5000, "BDT");
        a.Should().Be(b);
        (a == b).Should().BeTrue();
    }

    [Fact]
    public void Equality_DifferentAmount_AreNotEqual()
    {
        var a = new Money(5000, "BDT");
        var b = new Money(6000, "BDT");
        a.Should().NotBe(b);
        (a != b).Should().BeTrue();
    }

    [Fact]
    public void ToString_FormatsBdtCorrectly()
    {
        var money = new Money(150075, "BDT");
        money.ToString().Should().Be("1500.75 BDT");
    }

    // ── Premium calculation scenarios (matching CalculatePremiumHandler logic) ──

    [Fact]
    public void PremiumProRata_12MonthAnnual_CorrectMonthly()
    {
        // Annual premium = 12,000 BDT = 1,200,000 paisa
        // 1 month = 100,000 paisa = 1,000 BDT
        var annual = new Money(1_200_000, "BDT");
        var monthly = annual.Multiply(1m / 12m);
        monthly.AmountPaisa.Should().Be(100_000);
    }

    [Fact]
    public void LoadingFactor_25Percent_CorrectLoading()
    {
        var basePremium = new Money(100_000, "BDT"); // 1,000 BDT
        var loading = basePremium.Percentage(2500);   // 25% = 2500 basis points
        loading.AmountPaisa.Should().Be(25_000);       // 250 BDT
    }
}
