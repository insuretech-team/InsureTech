using FluentValidation;
using PoliSync.Products.Application.Commands;

namespace PoliSync.Products.Application.Validators;

public class CreateProductValidator : AbstractValidator<CreateProductCommand>
{
    public CreateProductValidator()
    {
        RuleFor(x => x.ProductCode)
            .NotEmpty().WithMessage("Product code is required")
            .Matches(@"^[A-Z]{3}-[0-9]{3}$").WithMessage("Product code must match pattern XXX-999 (e.g., MOT-001)");

        RuleFor(x => x.ProductName)
            .NotEmpty().WithMessage("Product name is required")
            .MaximumLength(255);

        RuleFor(x => x.Category)
            .IsInEnum().WithMessage("Invalid product category")
            .NotEqual(Domain.ProductCategory.Unspecified).WithMessage("Product category is required");

        RuleFor(x => x.BasePremium)
            .GreaterThan(0).WithMessage("Base premium must be greater than 0");

        RuleFor(x => x.MinSumInsured)
            .GreaterThan(0).WithMessage("Minimum sum insured must be greater than 0");

        RuleFor(x => x.MaxSumInsured)
            .GreaterThanOrEqualTo(x => x.MinSumInsured)
            .WithMessage("Maximum sum insured must be >= minimum sum insured");

        RuleFor(x => x.MinTenureMonths)
            .GreaterThan(0).WithMessage("Minimum tenure must be greater than 0");

        RuleFor(x => x.MaxTenureMonths)
            .GreaterThanOrEqualTo(x => x.MinTenureMonths)
            .WithMessage("Maximum tenure must be >= minimum tenure");

        RuleFor(x => x.CreatedBy)
            .NotEmpty().WithMessage("Created by user ID is required");
    }
}

public class UpdateProductValidator : AbstractValidator<UpdateProductCommand>
{
    public UpdateProductValidator()
    {
        RuleFor(x => x.ProductId)
            .NotEmpty().WithMessage("Product ID is required");

        RuleFor(x => x.ProductName)
            .MaximumLength(255)
            .When(x => x.ProductName is not null);

        RuleFor(x => x.BasePremium)
            .GreaterThan(0).WithMessage("Base premium must be greater than 0")
            .When(x => x.BasePremium.HasValue);

        RuleFor(x => x.MinSumInsured)
            .GreaterThan(0).WithMessage("Minimum sum insured must be greater than 0")
            .When(x => x.MinSumInsured.HasValue);

        RuleFor(x => x.MaxSumInsured)
            .GreaterThanOrEqualTo(x => x.MinSumInsured ?? 0)
            .WithMessage("Maximum sum insured must be >= minimum sum insured")
            .When(x => x.MaxSumInsured.HasValue);
    }
}
