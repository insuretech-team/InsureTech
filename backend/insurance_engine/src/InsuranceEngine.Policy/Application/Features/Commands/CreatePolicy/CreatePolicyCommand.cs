using System;
using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;

public record CreatePolicyCommand(
    Guid ProductId,
    Guid CustomerId,
    Guid? PartnerId,
    Guid? AgentId,
    ApplicantDto Applicant,
    List<NomineeDto>? Nominees,
    List<PolicyRiderDto>? Riders,
    long PremiumAmount,
    long SumInsuredAmount,
    int TenureMonths,
    DateTime StartDate
) : IRequest<Result<CreatePolicyResponse>>;

public record CreatePolicyResponse(Guid PolicyId, string PolicyNumber);
