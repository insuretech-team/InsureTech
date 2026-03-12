using FluentValidation;

namespace InsuranceEngine.Products.Application.Features.Commands.CreateProduct;

public class CreateProductCommandValidator : AbstractValidator<CreateProductCommand>
{
    public CreateProductCommandValidator()
    {
        RuleFor(x => x.ProductCode)
            .NotEmpty().WithMessage("Product code is required.")
            .Matches(@"^[A-Z]{3}-\d{3}$").WithMessage("Product code must match pattern: XXX-000 (e.g., HLT-001).");

        RuleFor(x => x.ProductName)
            .NotEmpty().WithMessage("Product name is required.")
            .MaximumLength(255);

        RuleFor(x => x.BasePremiumAmount)
            .GreaterThan(0).WithMessage("Base premium must be greater than 0.");

        RuleFor(x => x.MinSumInsuredAmount)
            .GreaterThan(0).WithMessage("Minimum sum insured must be greater than 0.");

        RuleFor(x => x.MaxSumInsuredAmount)
            .GreaterThan(0).WithMessage("Maximum sum insured must be greater than 0.");

        RuleFor(x => x)
            .Must(x => x.MinSumInsuredAmount < x.MaxSumInsuredAmount)
            .WithMessage("Minimum sum insured must be less than maximum sum insured.");

        RuleFor(x => x.MinTenureMonths)
            .GreaterThanOrEqualTo(1).WithMessage("Minimum tenure must be at least 1 month.");

        RuleFor(x => x.MaxTenureMonths)
            .LessThanOrEqualTo(360).WithMessage("Maximum tenure must be at most 360 months.");

        RuleFor(x => x)
            .Must(x => x.MinTenureMonths <= x.MaxTenureMonths)
            .WithMessage("Minimum tenure must be less than or equal to maximum tenure.");

        RuleFor(x => x.Category)
            .IsInEnum().WithMessage("Invalid product category.")
            .NotEqual(Domain.Enums.ProductCategory.Unspecified).WithMessage("Product category must be specified.");

        RuleFor(x => x.CreatedBy)
            .NotEmpty().WithMessage("CreatedBy is required.");
    }
}
