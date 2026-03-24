using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public record UpdateRiskScoreCommand(UpdateRiskScoreRequest Request) : IRequest<UpdateRiskScoreResponse>;
