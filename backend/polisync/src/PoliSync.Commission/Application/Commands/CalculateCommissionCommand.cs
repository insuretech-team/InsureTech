using Insuretech.Common.V1;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Commission.Application.Commands;

public sealed record CalculateCommissionCommand(
    string PolicyId,
    string CommissionType,
    string RecipientType,
    string RecipientId
) : ICommand<CalculateCommissionResult>;

public sealed record CalculateCommissionResult(
    string CommissionId,
    string CommissionNumber,
    Money Amount,
    string CalculationBreakdown);
