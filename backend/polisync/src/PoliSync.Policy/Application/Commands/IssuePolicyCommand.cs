using MediatR;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Policy.Application.Commands;

public sealed record IssuePolicyCommand(
    string QuoteId,
    string CustomerId,
    string ProductId,
    long PremiumAmountPaisa,
    long SumInsuredPaisa,
    int TenureMonths,
    DateTime StartDate,
    DateTime EndDate,
    string? PartnerId = null,
    string? AgentId = null
) : ICommand<string>; // Returns policy_id
