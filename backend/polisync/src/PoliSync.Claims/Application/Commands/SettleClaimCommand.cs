using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Claims.Application.Commands;

public sealed record SettleClaimCommand(
    string ClaimId,
    string SettledBy,
    string PaymentReference
) : ICommand<string>; // Returns payment_id
