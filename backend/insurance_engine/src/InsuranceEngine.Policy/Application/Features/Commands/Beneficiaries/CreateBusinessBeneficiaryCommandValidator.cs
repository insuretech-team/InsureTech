using FluentValidation;

namespace InsuranceEngine.Policy.Application.Features.Commands.Beneficiaries;

public class CreateBusinessBeneficiaryCommandValidator : AbstractValidator<CreateBusinessBeneficiaryCommand>
{
    public CreateBusinessBeneficiaryCommandValidator()
    {
        RuleFor(x => x.UserId).NotEmpty();
        RuleFor(x => x.BusinessName).NotEmpty().MaximumLength(255);
        RuleFor(x => x.TradeLicenseNumber).NotEmpty().MaximumLength(50);
        RuleFor(x => x.TinNumber).NotEmpty().Length(12);
        RuleFor(x => x.FocalPersonName).NotEmpty().MaximumLength(255);
        RuleFor(x => x.FocalPersonMobile)
            .NotEmpty()
            .Matches(@"^(?:\+88|88)?(01[3-9]\d{8})$")
            .WithMessage("Invalid Bangladesh mobile number format.");
    }
}
