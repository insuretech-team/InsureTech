using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.CancelPolicy;

public record CancelPolicyCommand(Guid PolicyId, string Reason) : IRequest<Result>;
