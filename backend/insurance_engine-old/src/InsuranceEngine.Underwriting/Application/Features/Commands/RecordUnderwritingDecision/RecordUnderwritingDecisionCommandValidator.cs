using FluentValidation;

namespace InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;

public class RecordUnderwritingDecisionCommandValidator : AbstractValidator<RecordUnderwritingDecisionCommand>
{
    public RecordUnderwritingDecisionCommandValidator()
    {
        RuleFor(x => x.QuoteId).NotEmpty();
        RuleFor(x => x.Decision).NotEmpty();
        RuleFor(x => x.Method).NotEmpty();
        RuleFor(x => x.RiskScore).InclusiveBetween(0, 100);
    }
}
