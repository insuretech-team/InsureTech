using Insuretech.Commission.Services.V1;
using MediatR;

namespace InsuranceEngine.Commission.Application.Commands;

public sealed record CalculateCommissionCommand(
    string PolicyId,
    string CommissionType,
    string RecipientType,
    string RecipientId) : IRequest<CalculateCommissionResponse>;

public sealed record CreatePayoutCommand(
    string RecipientType,
    string RecipientId,
    string PeriodStart,
    string PeriodEnd,
    List<string>? CommissionIds) : IRequest<CreatePayoutResponse>;

public sealed record ProcessPayoutCommand(
    string PayoutId,
    string PaymentMethod,
    string? PaymentReference) : IRequest<ProcessPayoutResponse>;
