using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Policy.Application.Commands;

public sealed record CancelPolicyCommand(string PolicyId, string Reason) : ICommand;
