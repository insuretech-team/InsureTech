using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.RenewPolicy;

public record RenewPolicyCommand(Guid PolicyId, int TenureMonths) : IRequest<Result<RenewPolicyResponse>>;
public record RenewPolicyResponse(Guid NewPolicyId, string NewPolicyNumber);
