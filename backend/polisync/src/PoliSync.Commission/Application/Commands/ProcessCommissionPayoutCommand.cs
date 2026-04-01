using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Commission.Application.Commands;

public sealed record ProcessCommissionPayoutCommand(
    string PayoutId,
    string PaymentMethod,
    string PaymentReference
) : ICommand<string>; // Returns paid_at timestamp
