using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Commission.Application.Commands;

public sealed record CalculateCommissionCommand(
    string PolicyId,
    string AgentId,
    decimal PremiumAmount) : ICommand<string>;

public sealed record ProcessPayoutCommand(string CommissionId) : ICommand<bool>;
