using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Policy.Application.Commands;

public sealed record CancelPolicyCommand(
    string PolicyId,
    string Reason,
    string? RequestedBy = null,  // user UUID who requested cancellation
    string? Portal = "B2C"       // portal context for workflow routing
) : ICommand;
