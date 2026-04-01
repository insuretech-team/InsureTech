using Insuretech.Policy.Services.V1;
using MediatR;

namespace InsuranceEngine.Cancellations.Application.Commands;

public sealed record CancelPolicyCommand(string PolicyId, string Reason) : IRequest<CancelPolicyResponse>;
