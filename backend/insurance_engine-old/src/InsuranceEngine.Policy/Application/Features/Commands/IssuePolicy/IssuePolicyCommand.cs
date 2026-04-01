using System;
using MediatR;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.IssuePolicy;

public record IssuePolicyCommand(Guid PolicyId) : IRequest<Result>;
