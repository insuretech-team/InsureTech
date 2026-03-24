using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public record UpdateBeneficiaryCommand(UpdateBeneficiaryRequest Request) : IRequest<UpdateBeneficiaryResponse>;
