using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Products.Application.Commands;

public sealed record UpdateProductCommand(
    string ProductId,
    string? ProductName,
    string? Description,
    decimal? BasePremium,
    decimal? MinSumInsured,
    decimal? MaxSumInsured,
    int? MinTenureMonths,
    int? MaxTenureMonths) : ICommand<bool>;

public sealed record ActivateProductCommand(string ProductId) : ICommand<bool>;
public sealed record DeactivateProductCommand(string ProductId) : ICommand<bool>;
public sealed record DiscontinueProductCommand(string ProductId) : ICommand<bool>;
