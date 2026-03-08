using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Policy.Application.Commands;

public sealed record RenewPolicyCommand(
    string PolicyId,
    long NewPremiumAmountPaisa,
    int NewTenureMonths
) : ICommand<string>; // Returns new policy_id
