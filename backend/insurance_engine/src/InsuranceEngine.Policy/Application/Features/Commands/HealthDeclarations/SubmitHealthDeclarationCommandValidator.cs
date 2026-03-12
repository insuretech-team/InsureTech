using FluentValidation;

namespace InsuranceEngine.Policy.Application.Features.Commands.HealthDeclarations;

public class SubmitHealthDeclarationCommandValidator : AbstractValidator<SubmitHealthDeclarationCommand>
{
    public SubmitHealthDeclarationCommandValidator()
    {
        RuleFor(x => x.QuoteId).NotEmpty();
        RuleFor(x => x.HealthDeclaration).NotNull();
        RuleFor(x => x.HealthDeclaration.HeightCm).InclusiveBetween(50, 300);
        RuleFor(x => x.HealthDeclaration.WeightKg).InclusiveBetween(10, 500);
        RuleFor(x => x.HealthDeclaration.Bmi).GreaterThan(0);
    }
}
