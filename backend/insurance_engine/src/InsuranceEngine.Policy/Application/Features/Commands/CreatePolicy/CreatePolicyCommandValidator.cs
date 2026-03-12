using FluentValidation;
using InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;

namespace InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;

public class CreatePolicyCommandValidator : AbstractValidator<CreatePolicyCommand>
{
    public CreatePolicyCommandValidator()
    {
        RuleFor(x => x.ProductId).NotEmpty().WithMessage("Product ID is required.");
        RuleFor(x => x.CustomerId).NotEmpty().WithMessage("Customer ID is required.");

        RuleFor(x => x.PremiumAmount).GreaterThan(0).WithMessage("Premium amount must be greater than 0.");
        RuleFor(x => x.SumInsuredAmount).GreaterThan(0).WithMessage("Sum insured must be greater than 0.");

        RuleFor(x => x.TenureMonths)
            .GreaterThanOrEqualTo(1).WithMessage("Tenure must be at least 1 month.")
            .LessThanOrEqualTo(1200).WithMessage("Tenure cannot exceed 1200 months.");

        RuleFor(x => x.StartDate)
            .GreaterThanOrEqualTo(DateTime.UtcNow.Date)
            .WithMessage("Start date must be today or in the future.");

        RuleFor(x => x.Applicant).NotNull().WithMessage("Applicant details are required.");
        When(x => x.Applicant != null, () =>
        {
            RuleFor(x => x.Applicant.FullName).NotEmpty().WithMessage("Applicant full name is required.");

            RuleFor(x => x.Applicant.NidNumber)
                .Matches(@"^\d{10}$|^\d{13}$|^\d{17}$")
                .When(x => !string.IsNullOrEmpty(x.Applicant?.NidNumber))
                .WithMessage("NID must be 10, 13, or 17 digits.");

            RuleFor(x => x.Applicant.PhoneNumber)
                .Matches(@"^\+880\s?1[3-9]\d{2}\s?\d{6}$")
                .When(x => !string.IsNullOrEmpty(x.Applicant?.PhoneNumber))
                .WithMessage("Phone must match Bangladesh format: +880 1XXX XXXXXX.");
        });
    }
}
