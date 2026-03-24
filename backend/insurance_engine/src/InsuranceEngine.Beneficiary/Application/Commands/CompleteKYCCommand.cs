using MediatR;
using Insuretech.Beneficiary.Services.V1;

namespace InsuranceEngine.Beneficiary.Application.Commands;

public record CompleteKYCCommand(CompleteKYCRequest Request) : IRequest<CompleteKYCResponse>;
