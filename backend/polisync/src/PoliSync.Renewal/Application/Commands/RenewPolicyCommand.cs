using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Renewal.Application.Commands;

public sealed record RenewPolicyCommand(string PolicyId) : ICommand<string>; // Returns new policy id
