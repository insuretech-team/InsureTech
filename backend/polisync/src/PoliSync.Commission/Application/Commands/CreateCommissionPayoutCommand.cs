using Insuretech.Common.V1;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Commission.Application.Commands;

public sealed record CreateCommissionPayoutCommand(
    string RecipientType,
    string RecipientId,
    string PeriodStart,
    string PeriodEnd,
    List<string> CommissionIds
) : ICommand<CreateCommissionPayoutResult>;

public sealed record CreateCommissionPayoutResult(
    string PayoutId,
    string PayoutNumber,
    Money TotalAmount,
    int CommissionCount);
