using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Policy.Services.V1;
using MediatR;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed record IssuePolicyCommand(string PolicyId) : IRequest<IssuePolicyResponse>;
public sealed record CancelPolicyCommand(string PolicyId, string Reason) : IRequest<bool>;
public sealed record RenewPolicyCommand(string PolicyId, int TenureMonths) : IRequest<RenewPolicyResult>;

public sealed record RenewPolicyResult(string NewPolicyId, string NewPolicyNumber);
