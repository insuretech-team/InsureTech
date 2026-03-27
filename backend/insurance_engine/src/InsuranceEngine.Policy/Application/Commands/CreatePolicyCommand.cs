using MediatR;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed record CreatePolicyCommand(
    string ProductId,
    string CustomerId,
    string? PartnerId,
    string? AgentId,
    string? QuoteId,
    decimal PremiumAmount,
    decimal SumInsured,
    int TenureMonths,
    DateTime StartDate,
    string? ProposerDetails,
    List<Insuretech.Policy.Entity.V1.Nominee>? Nominees) : IRequest<CreatePolicyResponse>;
