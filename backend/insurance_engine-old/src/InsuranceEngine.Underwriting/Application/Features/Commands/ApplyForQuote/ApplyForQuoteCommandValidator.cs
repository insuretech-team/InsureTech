using FluentValidation;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.ApplyForQuote;

public class ApplyForQuoteCommandValidator : AbstractValidator<ApplyForQuoteCommand>
{
    public ApplyForQuoteCommandValidator()
    {
        RuleFor(x => x.BeneficiaryId).NotEmpty();
        RuleFor(x => x.ProductId).NotEmpty();
        RuleFor(x => x.SumAssuredAmount).GreaterThan(0);
        RuleFor(x => x.TermYears).InclusiveBetween(1, 50);
        RuleFor(x => x.PremiumPaymentMode).NotEmpty();
        RuleFor(x => x.ApplicantAge).InclusiveBetween(0, 100);
        
        RuleFor(x => x.HealthDeclaration).NotNull();
        RuleFor(x => x.HealthDeclaration.HeightCm).InclusiveBetween(50, 300);
        RuleFor(x => x.HealthDeclaration.WeightKg).InclusiveBetween(10, 500);
    }
}
