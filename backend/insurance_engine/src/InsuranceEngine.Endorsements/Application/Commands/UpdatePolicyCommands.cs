using Insuretech.Policy.Services.V1;
using MediatR;
using System.Collections.Generic;

namespace InsuranceEngine.Endorsements.Application.Commands;

public sealed record UpdatePolicyCommand(
    string PolicyId, 
    List<Insuretech.Policy.Entity.V1.Nominee>? Nominees = null,
    string? Address = null) : IRequest<UpdatePolicyResponse>;
