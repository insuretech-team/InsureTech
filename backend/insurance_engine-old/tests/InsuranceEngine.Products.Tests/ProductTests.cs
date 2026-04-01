using FluentAssertions;
using InsuranceEngine.Products.Domain;
using InsuranceEngine.Products.Domain.Enums;
using Xunit;

namespace InsuranceEngine.Products.Tests;

public class ProductTests
{
    [Fact]
    public void Activate_WhenStatusIsDraft_ShouldSucceed()
    {
        // Arrange
        var product = new Product { Status = ProductStatus.Draft };

        // Act
        var result = product.Activate();

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be(ProductStatus.Active);
    }

    [Fact]
    public void Activate_WhenStatusIsActive_ShouldFail()
    {
        // Arrange
        var product = new Product { Status = ProductStatus.Active };

        // Act
        var result = product.Activate();

        // Assert
        result.IsSuccess.Should().BeFalse();
        product.Status.Should().Be(ProductStatus.Active);
    }

    [Fact]
    public void Deactivate_WhenStatusIsActive_ShouldSucceed()
    {
        // Arrange
        var product = new Product { Status = ProductStatus.Active };

        // Act
        var result = product.Deactivate();

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be(ProductStatus.Inactive);
    }

    [Fact]
    public void Discontinue_ShouldSetStatusToDiscontinued()
    {
        // Arrange
        var product = new Product { Status = ProductStatus.Active };

        // Act
        var result = product.Discontinue();

        // Assert
        result.IsSuccess.Should().BeTrue();
        product.Status.Should().Be(ProductStatus.Discontinued);
    }
}
