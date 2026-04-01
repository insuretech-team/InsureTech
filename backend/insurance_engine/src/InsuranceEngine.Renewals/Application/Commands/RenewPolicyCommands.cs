using Insuretech.Policy.Services.V1;
using MediatR;
using System.Collections.Generic;

namespace InsuranceEngine.Renewals.Application.Commands;

public sealed record RenewPolicyCommand(
    string PolicyId, 
    int TenureMonths, 
    bool UpdateNominees = false, 
    List<Insuretech.Policy.Entity.V1.Nominee>? Nominees = null) : IRequest<RenewPolicyTenureResponse>;
