using FluentValidation;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

public class CreateIndividualBeneficiaryCommandValidator : AbstractValidator<CreateIndividualBeneficiaryCommand>
{
    public CreateIndividualBeneficiaryCommandValidator()
    {
        RuleFor(x => x.UserId).NotEmpty();
        RuleFor(x => x.FullName).NotEmpty().MaximumLength(255);
        RuleFor(x => x.DateOfBirth).NotEmpty().LessThan(DateTime.UtcNow);
        RuleFor(x => x.Gender).NotEmpty();
        RuleFor(x => x.NidNumber)
            .NotEmpty()
            .Matches(@"^\d{10}$|^\d{13}$|^\d{17}$")
            .WithMessage("NID must be 10, 13, or 17 digits.");
        RuleFor(x => x.MobileNumber)
            .NotEmpty()
            .Matches(@"^(?:\+88|88)?(01[3-9]\d{8})$")
            .WithMessage("Invalid Bangladesh mobile number format.");
    }
}
